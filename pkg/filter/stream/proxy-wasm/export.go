package proxywasm

// #include <stdlib.h>
//
// extern int proxy_log(void *context, int, int, int);
// extern int proxy_get_property(void *context, int, int, int, int);
// extern int proxy_get_header_map_pairs(void *context, int, int, int);
// extern int proxy_get_buffer_bytes(void *context, int, int, int, int, int);
// extern int proxy_replace_header_map_value(void *context, int, int, int, int, int);
// extern int proxy_add_header_map_value(void *context, int, int, int, int, int);
// extern int proxy_set_effective_context(void *context, int);
import "C"

import (
	"encoding/binary"
	"unsafe"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
)

//export proxy_get_buffer_bytes
func proxy_get_buffer_bytes(context unsafe.Pointer, bufferType int32, start int32, maxSize int32, returnBufferData int32, returnBufferSize int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	ctx := instanceContext.Data().(*wasmContext)

	var body []byte
	switch BufferType(bufferType) {
	case BufferTypeHttpRequestBody:
		body = ctx.filter.rhandler.GetRequestData().Bytes()
	case BufferTypeHttpResponseBody:
		body = ctx.filter.shandler.GetResponseData().Bytes()
	}

	r, _ := ctx.instance.Exports["malloc"](len(body))
	p := r.ToI32()
	memory := ctx.instance.Memory.Data()
	copy(memory[p:], body)

	binary.LittleEndian.PutUint32(memory[returnBufferData:], uint32(p))
	binary.LittleEndian.PutUint32(memory[returnBufferSize:], uint32(len(body)))

	return 0
}

//export proxy_get_header_map_pairs
func proxy_get_header_map_pairs(context unsafe.Pointer, mapType int32, returnDataPtr int32, returnDataSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_get_header_map_pairs")
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	header := ctx.GetHeaderMap(MapType(mapType))
	if header == nil {
		return WasmResultNotFound.Int32()
	}

	headerSize := 0

	addr, err := ctx.malloc(int32(header.ByteSize() + 200))
	if err != nil {
		log.DefaultLogger.Errorf("wasm malloc error: %v", err)
		return WasmResultInternalFailure.Int32()
	}
	start := addr
	p := start

	p += 4
	header.Range(func(key string, value string) bool {
		headerSize++
		binary.LittleEndian.PutUint32(memory[p:], uint32(len(key)))
		p += 4
		binary.LittleEndian.PutUint32(memory[p:], uint32(len(value)))
		p += 4
		return true
	})

	header.Range(func(key string, value string) bool {
		copy(memory[p:], key)
		p += int32(len(key))
		memory[p] = 0
		p++

		copy(memory[p:], value)
		p += int32(len(value))
		memory[p] = 0
		p++

		return true
	})

	binary.LittleEndian.PutUint32(memory[start:], uint32(headerSize))

	binary.LittleEndian.PutUint32(memory[returnDataPtr:], uint32(start))
	binary.LittleEndian.PutUint32(memory[returnDataSize:], uint32(p-start))

	return WasmResultOk.Int32()
}

//export proxy_replace_header_map_value
func proxy_replace_header_map_value(context unsafe.Pointer, mapType int32, keyData int32, keySize int32, valueData int32, valueSize int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	ctx := instanceContext.Data().(*wasmContext)

	key := string(memory[keyData : keyData+keySize])
	value := string(memory[valueData : valueData+valueSize])

	var header api.HeaderMap
	switch MapType(mapType) {
	case MapTypeHttpRequestHeaders:
		header = ctx.filter.rhandler.GetRequestHeaders()
		header.Set(key, value)
	case MapTypeHttpResponseHeaders:
		header = ctx.filter.shandler.GetResponseHeaders()
		header.Set(key, value)
	}

	return 0
}

//export proxy_add_header_map_value
func proxy_add_header_map_value(context unsafe.Pointer, mapType int32, keyData int32, keySize int32, valueData int32, valueSize int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	ctx := instanceContext.Data().(*wasmContext)

	key := string(memory[keyData : keyData+keySize])
	value := string(memory[valueData : valueData+valueSize])

	var header api.HeaderMap
	switch MapType(mapType) {
	case MapTypeHttpRequestHeaders:
		header = ctx.filter.rhandler.GetRequestHeaders()
		header.Add(key, value)
	case MapTypeHttpResponseHeaders:
		header = ctx.filter.shandler.GetResponseHeaders()
		header.Add(key, value)
	}

	return 0
}

//export proxy_log
func proxy_log(context unsafe.Pointer, logLevel int32, messageData int32, messageSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_log")
	var instanceCtx = wasm.IntoInstanceContext(context)
	var memory = instanceCtx.Memory().Data()

	msg := string(memory[messageData : messageData+messageSize])
	log.DefaultLogger.Errorf("wasm log: %s", msg)
	return WasmResultOk.Int32()
}

//export proxy_get_property
func proxy_get_property(context unsafe.Pointer, pathData int32, pathSize int32, returnValueData int32, returnValueSize int32) int32 {

	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	//key := string(memory[pathData : pathData+pathSize])

	// TODO: property not exist
	//value := ctx.rootContext.GetProperty(key)
	//if value == "" {
	//	return 0
	//}
	value := "my_root_id"

	addr, err := ctx.malloc(int32(len(value)))
	if err != nil {
		log.DefaultLogger.Errorf("wasm malloc error: %v", err)
		return WasmResultInternalFailure.Int32()
	}
	p := addr
	copy(memory[p:], value)

	binary.LittleEndian.PutUint32(memory[returnValueData:], uint32(p))
	binary.LittleEndian.PutUint32(memory[returnValueSize:], uint32(len(value)))
	return WasmResultOk.Int32()
}

//export proxy_set_effective_context
func proxy_set_effective_context(context unsafe.Pointer, context_id int32) int32 {
	return 0
}


func ProxyWasmImports() *wasm.Imports {
	im := wasm.NewImports()
	im, _ = im.AppendFunction("proxy_log", proxy_log, C.proxy_log)
	im, _ = im.AppendFunction("proxy_get_property", proxy_get_property, C.proxy_get_property)
	im, _ = im.AppendFunction("proxy_get_header_map_pairs", proxy_get_header_map_pairs, C.proxy_get_header_map_pairs)
	im, _ = im.AppendFunction("proxy_get_buffer_bytes", proxy_get_buffer_bytes, C.proxy_get_buffer_bytes)
	im, _ = im.AppendFunction("proxy_replace_header_map_value", proxy_replace_header_map_value, C.proxy_replace_header_map_value)
	im, _ = im.AppendFunction("proxy_add_header_map_value", proxy_add_header_map_value, C.proxy_add_header_map_value)

	im, _ = im.AppendFunction("proxy_set_effective_context", proxy_set_effective_context, C.proxy_set_effective_context)

	return im
}