package proxywasm

// #include <stdlib.h>
//
// extern int proxy_log(void *context, int, int, int);
// extern int proxy_get_property(void *context, int, int, int, int);
// extern int proxy_set_effective_context(void *context, int);
//
// extern int proxy_get_buffer_bytes(void *context, int, int, int, int, int);
// extern int proxy_get_header_map_pairs(void *context, int, int, int);
// extern int proxy_get_header_map_value(void *context, int, int, int, int, int);
// extern int proxy_replace_header_map_value(void *context, int, int, int,int, int);
// extern int proxy_add_header_map_value(void *context, int, int, int, int, int);
// extern int proxy_remove_header_map_value(void *context, int, int, int);
import "C"

import (
	"encoding/binary"
	"unsafe"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)


func ProxyWasmImports() *wasm.Imports {
	im := wasm.NewImports()
	im, _ = im.AppendFunction("proxy_log", proxy_log, C.proxy_log)
	im, _ = im.AppendFunction("proxy_get_property", proxy_get_property, C.proxy_get_property)
	im, _ = im.AppendFunction("proxy_set_effective_context", proxy_set_effective_context, C.proxy_set_effective_context)

	im, _ = im.AppendFunction("proxy_get_buffer_bytes", proxy_get_buffer_bytes, C.proxy_get_buffer_bytes)
	im, _ = im.AppendFunction("proxy_get_header_map_pairs", proxy_get_header_map_pairs, C.proxy_get_header_map_pairs)
	im, _ = im.AppendFunction("proxy_get_header_map_value", proxy_get_header_map_value, C.proxy_get_header_map_value)
	im, _ = im.AppendFunction("proxy_replace_header_map_value", proxy_replace_header_map_value, C.proxy_replace_header_map_value)
	im, _ = im.AppendFunction("proxy_add_header_map_value", proxy_add_header_map_value, C.proxy_add_header_map_value)
	im, _ = im.AppendFunction("proxy_remove_header_map_value", proxy_remove_header_map_value, C.proxy_remove_header_map_value)

	return im
}

//export proxy_get_buffer_bytes
func proxy_get_buffer_bytes(context unsafe.Pointer, bufferType int32, start int32, length int32, returnBufferData int32, returnBufferSize int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if BufferType(bufferType) > BufferTypeMax {
		return WasmResultBadArgument.Int32()
	}

	body := ctx.GetBuffer(BufferType(bufferType))
	if body == nil {
		return WasmResultNotFound.Int32()
	}

	// Check for overflow.
	if start > start+length {
		return WasmResultBadArgument.Int32()
	}

	if start+length > int32(len(body)) {
		length = int32(len(body)) - start
	}

	if length > 0 {
		addr, err := ctx.malloc(length)
		if err != nil {
			log.DefaultLogger.Errorf("wasm malloc error: %v", err)
			return WasmResultInternalFailure.Int32()
		}
		copy(memory[addr:], body[start:start+length])

		binary.LittleEndian.PutUint32(memory[returnBufferData:], uint32(addr))
		binary.LittleEndian.PutUint32(memory[returnBufferSize:], uint32(length))
	}
	return WasmResultOk.Int32()
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

//export proxy_get_header_map_value
func proxy_get_header_map_value(context unsafe.Pointer, mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_get_header_map_value")
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	key := string(memory[keyDataPtr : keyDataPtr+keySize])
	if key == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	value := ctx.GetHeaderMapValue(MapType(mapType), key)

	addr, err := ctx.malloc(int32(len(value)))
	if err != nil {
		log.DefaultLogger.Errorf("wasm malloc error: %v", err)
		return WasmResultInternalFailure.Int32()
	}
	copy(memory[addr:], value)

	binary.LittleEndian.PutUint32(memory[valueDataPtr:], uint32(addr))
	binary.LittleEndian.PutUint32(memory[valueSize:], uint32(len(value)))

	return WasmResultOk.Int32()
}

//export proxy_replace_header_map_value
func proxy_replace_header_map_value(context unsafe.Pointer, mapType int32, keyData int32, keySize int32, valueData int32, valueSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_replace_header_map_value")
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	key := string(memory[keyData : keyData+keySize])
	if key == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	// TODO: need to check null for value?
	value := string(memory[valueData : valueData+valueSize])

	ctx.SetHeaderMapValue(MapType(mapType), key, value)

	return WasmResultOk.Int32()
}

//export proxy_add_header_map_value
func proxy_add_header_map_value(context unsafe.Pointer, mapType int32, keyData int32, keySize int32, valueData int32, valueSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_add_header_map_value")
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	key := string(memory[keyData : keyData+keySize])
	if key == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	value := string(memory[valueData : valueData+valueSize])
	if value == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	ctx.SetHeaderMapValue(MapType(mapType), key, value)

	return WasmResultOk.Int32()
}

//export proxy_remove_header_map_value
func proxy_remove_header_map_value(context unsafe.Pointer, mapType int32, keyDataPtr int32, keySize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_remove_header_map_value")
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	key := string(memory[keyDataPtr : keyDataPtr+keySize])
	if key == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	ctx.DelHeaderMapValue(MapType(mapType), key)

	return WasmResultOk.Int32()
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