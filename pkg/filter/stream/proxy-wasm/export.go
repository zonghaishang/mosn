package proxywasm

// #include <stdlib.h>
//
// extern int proxy_log(void *context, int, int, int);
// extern int proxy_get_property(void *context, int, int, int, int);
// extern int proxy_set_effective_context(void *context, int);
//
// extern int proxy_get_buffer_bytes(void *context, int, int, int, int, int);
// extern int proxy_set_buffer_bytes(void *context, int, int, int, int, int);
//
// extern int proxy_get_header_map_pairs(void *context, int, int, int);
// extern int proxy_get_header_map_value(void *context, int, int, int, int, int);
// extern int proxy_replace_header_map_value(void *context, int, int, int,int, int);
// extern int proxy_add_header_map_value(void *context, int, int, int, int, int);
// extern int proxy_remove_header_map_value(void *context, int, int, int);
//
// extern int proxy_set_tick_period_milliseconds(void *context, int);
// extern int proxy_get_current_time_nanoseconds(void *context, int);
//
// extern int proxy_grpc_call(void *context, int, int, int, int, int, int, int, int, int, int, int, int);
// extern int proxy_grpc_stream(void *context, int, int, int, int, int, int, int, int, int);
// extern int proxy_grpc_cancel(void *context, int);
// extern int proxy_grpc_close(void *context, int);
// extern int proxy_grpc_send(void *context, int, int, int, int);
//
// extern int proxy_http_call(void *context, int, int, int, int, int, int, int, int, int, int);
//
// extern int proxy_define_metric(void *context, int, int, int, int);
// extern int proxy_increment_metric(void *context, int, int);
// extern int proxy_record_metric(void *context, int, int);
// extern int proxy_get_metric(void *context, int, int);
//
// extern int proxy_register_shared_queue(void *context, int, int, int);
// extern int proxy_resolve_shared_queue(void *context, int, int, int, int, int);
// extern int proxy_dequeue_shared_queue(void *context, int, int, int);
// extern int proxy_enqueue_shared_queue(void *context, int, int, int);
//
// extern int proxy_get_shared_data(void *context, int, int, int, int, int);
// extern int proxy_set_shared_data(void *context, int, int, int, int, int);
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
	im, _ = im.AppendFunction("proxy_set_buffer_bytes", proxy_set_buffer_bytes, C.proxy_set_buffer_bytes)

	im, _ = im.AppendFunction("proxy_get_header_map_pairs", proxy_get_header_map_pairs, C.proxy_get_header_map_pairs)
	im, _ = im.AppendFunction("proxy_get_header_map_value", proxy_get_header_map_value, C.proxy_get_header_map_value)
	im, _ = im.AppendFunction("proxy_replace_header_map_value", proxy_replace_header_map_value, C.proxy_replace_header_map_value)
	im, _ = im.AppendFunction("proxy_add_header_map_value", proxy_add_header_map_value, C.proxy_add_header_map_value)
	im, _ = im.AppendFunction("proxy_remove_header_map_value", proxy_remove_header_map_value, C.proxy_remove_header_map_value)

	im, _ = im.AppendFunction("proxy_set_tick_period_milliseconds", proxy_set_tick_period_milliseconds, C.proxy_set_tick_period_milliseconds)
	im, _ = im.AppendFunction("proxy_get_current_time_nanoseconds", proxy_get_current_time_nanoseconds, C.proxy_get_current_time_nanoseconds)

	im, _ = im.AppendFunction("proxy_grpc_call", proxy_grpc_call, C.proxy_grpc_call)
	im, _ = im.AppendFunction("proxy_grpc_stream", proxy_grpc_stream, C.proxy_grpc_stream)
	im, _ = im.AppendFunction("proxy_grpc_cancel", proxy_grpc_cancel, C.proxy_grpc_cancel)
	im, _ = im.AppendFunction("proxy_grpc_close", proxy_grpc_close, C.proxy_grpc_close)
	im, _ = im.AppendFunction("proxy_grpc_send", proxy_grpc_send, C.proxy_grpc_send)

	im, _ = im.AppendFunction("proxy_http_call", proxy_http_call, C.proxy_http_call)

	im, _ = im.AppendFunction("proxy_define_metric", proxy_define_metric, C.proxy_define_metric)
	im, _ = im.AppendFunction("proxy_increment_metric", proxy_increment_metric, C.proxy_increment_metric)
	im, _ = im.AppendFunction("proxy_record_metric", proxy_record_metric, C.proxy_record_metric)
	im, _ = im.AppendFunction("proxy_get_metric", proxy_get_metric, C.proxy_get_metric)

	im, _ = im.AppendFunction("proxy_register_shared_queue", proxy_register_shared_queue, C.proxy_register_shared_queue)
	im, _ = im.AppendFunction("proxy_resolve_shared_queue", proxy_resolve_shared_queue, C.proxy_resolve_shared_queue)
	im, _ = im.AppendFunction("proxy_dequeue_shared_queue", proxy_dequeue_shared_queue, C.proxy_dequeue_shared_queue)
	im, _ = im.AppendFunction("proxy_enqueue_shared_queue", proxy_enqueue_shared_queue, C.proxy_enqueue_shared_queue)

	im, _ = im.AppendFunction("proxy_get_shared_data", proxy_get_shared_data, C.proxy_get_shared_data)
	im, _ = im.AppendFunction("proxy_set_shared_data", proxy_set_shared_data, C.proxy_set_shared_data)

	return im
}

//export proxy_get_buffer_bytes
func proxy_get_buffer_bytes(context unsafe.Pointer, bufferType int32, start int32, length int32, returnBufferData int32, returnBufferSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_get_buffer_bytes")

	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if BufferType(bufferType) > BufferTypeMax {
		return WasmResultBadArgument.Int32()
	}

	body, result := ctx.GetBuffer(BufferType(bufferType))
	if result != WasmResultOk {
		return result.Int32()
	}
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

//export proxy_set_buffer_bytes
func proxy_set_buffer_bytes(context unsafe.Pointer, bufferType int32, start int32, length int32, dataPtr int32, dataSize int32) int32 {
	log.DefaultLogger.Debugf("wasm call host.proxy_set_buffer_bytes")

	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if BufferType(bufferType) > BufferTypeMax ||
		(BufferType(bufferType) != BufferTypeHttpRequestBody && BufferType(bufferType) != BufferTypeHttpResponseBody) {
		return WasmResultBadArgument.Int32()
	}

	buffer, result := ctx.GetBuffer(BufferType(bufferType))
	if result != WasmResultOk {
		return result.Int32()
	}

	// Check for overflow.
	if start > start+length {
		return WasmResultBadArgument.Int32()
	}
	if start + length > int32(len(buffer)) {
		length = int32(len(buffer)) - start
	}
	if dataSize > length {
		dataSize = length
	}

	copy(buffer[start:], memory[dataPtr:dataPtr+dataSize])

	result = ctx.SetBuffer(BufferType(bufferType), buffer)

	return result.Int32()
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

	header, result := ctx.GetHeaderMap(MapType(mapType))
	if result != WasmResultOk {
		return result.Int32()
	}
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

	value, result := ctx.GetHeaderMapValue(MapType(mapType), key)
	if result != WasmResultOk {
		return result.Int32()
	}

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

	result := ctx.SetHeaderMapValue(MapType(mapType), key, value)

	return result.Int32()
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

	result := ctx.SetHeaderMapValue(MapType(mapType), key, value)

	return result.Int32()
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

	result := ctx.DelHeaderMapValue(MapType(mapType), key)

	return result.Int32()
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
	copy(memory[addr:], value)

	binary.LittleEndian.PutUint32(memory[returnValueData:], uint32(addr))
	binary.LittleEndian.PutUint32(memory[returnValueSize:], uint32(len(value)))
	return WasmResultOk.Int32()
}

//export proxy_set_effective_context
func proxy_set_effective_context(context unsafe.Pointer, context_id int32) int32 {
	return 0
}

// Set timer period. Once set, the host environment will call proxy_on_tick every tick_period_milliseconds.
//export proxy_set_tick_period_milliseconds
func proxy_set_tick_period_milliseconds(context unsafe.Pointer, tickPeriodMilliseconds int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	return ctx.SetTimerPeriod(int64(tickPeriodMilliseconds), 0).Int32()
}

//export proxy_get_current_time_nanoseconds
func proxy_get_current_time_nanoseconds(context unsafe.Pointer, resultUint64Ptr int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	nanoSecond, result := ctx.GetCurrentTimeNanoseconds()
	if result != WasmResultOk {
		return result.Int32()
	}

	binary.LittleEndian.PutUint64(memory[resultUint64Ptr:], uint64(nanoSecond))

	return WasmResultOk.Int32()
}

//export proxy_grpc_call
func proxy_grpc_call(context unsafe.Pointer,
	servicePtr int32, serviceSize int32, serviceNamePtr int32, serviceNameSize int32,
	methodNamePtr int32, methodNameSize int32,
	initialMetadataPtr int32, initialMetadataSize int32,
	requestPtr int32, requestSize int32,
	timeoutMilliseconds int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

//export proxy_grpc_stream
func proxy_grpc_stream(context unsafe.Pointer,
	servicePtr int32, serviceSize int32, serviceNamePtr int32, serviceNameSize int32,
	methodNamePtr int32, methodNameSize int32,
	initialMetadataPtr int32, initialMetadataSize int32,
	tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

//export proxy_grpc_cancel
func proxy_grpc_cancel(context unsafe.Pointer, token int32) int32 {
	return WasmResultUnimplemented.Int32()
}

//export proxy_grpc_close
func proxy_grpc_close(context unsafe.Pointer, token int32) int32 {
	return WasmResultUnimplemented.Int32()
}

//export proxy_grpc_send
func proxy_grpc_send(context unsafe.Pointer, token int32, messagePtr int32, messageSize int32, endStream int32) int32 {
	return WasmResultUnimplemented.Int32()
}

//export proxy_http_call
func proxy_http_call(context unsafe.Pointer,
	uriPtr int32, uriSize int32,
	headerPairsPtr int32, headerPairsSize int32,
	bodyPtr int32, bodySize int32,
	trailerPairsPtr int32, trailerPairsSize int32,
	timeoutMilliseconds int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

//export proxy_define_metric
func proxy_define_metric(context unsafe.Pointer, metricType int32, namePtr int32, nameSize int32, resultPtr int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if MetricType(metricType) > MetricTypeMax {
		return WasmResultBadArgument.Int32()
	}
	if namePtr > namePtr+nameSize {
		return WasmResultBadArgument.Int32()
	}

	name := string(memory[namePtr : namePtr+nameSize])
	if name == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	metricId, result := ctx.DefineMetric(MetricType(metricType), name)
	if result != WasmResultOk {
		return result.Int32()
	}

	binary.LittleEndian.PutUint32(memory[resultPtr:], metricId)

	return WasmResultOk.Int32()
}

//export proxy_increment_metric
func proxy_increment_metric(context unsafe.Pointer, metricId int32, offset int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)

	result := ctx.IncrementMetric(uint32(metricId), int64(offset))
	if result != WasmResultOk {
		return result.Int32()
	}

	return WasmResultOk.Int32()
}

//export proxy_record_metric
func proxy_record_metric(context unsafe.Pointer, metricId int32, value int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)

	result := ctx.RecordMetric(uint32(metricId), uint64(value))
	if result != WasmResultOk {
		return result.Int32()
	}

	return WasmResultOk.Int32()
}

//export proxy_get_metric
func proxy_get_metric(context unsafe.Pointer, metricId int32, resultUint64Ptr int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	value, result := ctx.GetMetric(uint32(metricId))
	if result != WasmResultOk {
		return result.Int32()
	}

	binary.LittleEndian.PutUint64(memory[resultUint64Ptr:], uint64(value))

	return WasmResultOk.Int32()
}

//export proxy_register_shared_queue
func proxy_register_shared_queue(context unsafe.Pointer, queueNamePtr int32, queueNameSize int32, tokenPtr int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if queueNamePtr > queueNamePtr+queueNameSize {
		return WasmResultBadArgument.Int32()
	}

	queueName := string(memory[queueNamePtr : queueNamePtr+queueNameSize])
	if queueName == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	queueToken, result := ctx.RegisterSharedQueue(queueName)
	if result != WasmResultOk {
		return result.Int32()
	}

	binary.LittleEndian.PutUint32(memory[tokenPtr:], uint32(queueToken))

	return WasmResultOk.Int32()
}

// TODO: currently we ignore vmId
//export proxy_resolve_shared_queue
func proxy_resolve_shared_queue(context unsafe.Pointer, vmIdPtr int32, vmIdSize int32, queueNamePtr int32, queueNameSize int32, tokenPtr int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if queueNamePtr > queueNamePtr+queueNameSize {
		return WasmResultBadArgument.Int32()
	}

	queueName := string(memory[queueNamePtr : queueNamePtr+queueNameSize])
	if queueName == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	queueToken, result := ctx.LookupSharedQueue(queueName)
	if result != WasmResultOk {
		return result.Int32()
	}

	binary.LittleEndian.PutUint32(memory[tokenPtr:], uint32(queueToken))

	return WasmResultOk.Int32()
}

//export proxy_dequeue_shared_queue
func proxy_dequeue_shared_queue(context unsafe.Pointer, token int32, dataPtr int32, dataSize int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if dataPtr > dataPtr+dataSize {
		return WasmResultBadArgument.Int32()
	}
	data := string(memory[dataPtr : dataPtr+dataSize])
	if data == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	result := ctx.DequeueSharedQueue(uint32(token), data)

	return result.Int32()
}

//export proxy_enqueue_shared_queue
func proxy_enqueue_shared_queue(context unsafe.Pointer, token int32, dataPtr int32, dataSize int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	if dataPtr > dataPtr+dataSize {
		return WasmResultBadArgument.Int32()
	}
	data := string(memory[dataPtr : dataPtr+dataSize])
	if data == "" {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	result := ctx.DequeueSharedQueue(uint32(token), data)

	return result.Int32()
}

// Get proxy-wide key-value data shared between VMs.
// TODO: currently we ignore the 'cas' field
//export proxy_get_shared_data
func proxy_get_shared_data(context unsafe.Pointer, keyPtr int32, keySize int32, valuePtr int32, valueSizePtr int32, casPtr int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	key := string(memory[keyPtr : keyPtr+keySize])

	value, cas, result := ctx.GetSharedData(key)
	if result != WasmResultOk {
		return result.Int32()
	}

	addr, err := ctx.malloc(int32(len(value)))
	if err != nil {
		log.DefaultLogger.Errorf("wasm malloc error: %v", err)
		return WasmResultInternalFailure.Int32()
	}

	copy(memory[addr:], value)
	binary.LittleEndian.PutUint32(memory[valuePtr:], uint32(addr))
	binary.LittleEndian.PutUint32(memory[valueSizePtr:], uint32(len(value)))

	binary.LittleEndian.PutUint32(memory[casPtr:], uint32(cas))

	return WasmResultOk.Int32()
}

// Set a key-value data shared between VMs.
// TODO: currently we ignore the 'cas' field
//export proxy_set_shared_data
func proxy_set_shared_data(context unsafe.Pointer, keyPtr int32, keySize int32, valuePtr int32, valueSize int32, cas int32) int32 {
	var instanceCtx = wasm.IntoInstanceContext(context)
	ctx := instanceCtx.Data().(*wasmContext)
	memory := ctx.instance.Memory.Data()

	key := string(memory[keyPtr : keyPtr+keySize])
	value := string(memory[valuePtr : valuePtr+valueSize])

	result := ctx.SetSharedData(key, value, uint32(cas))

	return result.Int32()
}