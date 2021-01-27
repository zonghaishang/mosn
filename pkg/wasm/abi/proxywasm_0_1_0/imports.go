package proxywasm_0_1_0

import (
	"encoding/binary"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/buffer"
)

func (a *abiImpl) proxyLog(level int32, logDataPtr int32, logDataSize int32) int32 {
	logContent, err := a.instance.GetMemory(uint64(logDataPtr), uint64(logDataSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	a.callback.Log(toMosnLogLevel(LogLevel(level)), string(logContent))

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxyGetBufferBytes(bufferType int32, start int32, length int32, returnBufferData int32, returnBufferSize int32) int32 {
	if BufferType(bufferType) > BufferTypeMax {
		return WasmResultBadArgument.Int32()
	}

	buf := a.GetBuffer(BufferType(bufferType))
	if buf == nil {
		return WasmResultNotFound.Int32()
	}

	if  start > start+length {
		return WasmResultBadArgument.Int32()
	}

	if start + length > int32(buf.Len()) {
		length = int32(buf.Len()) - start
	}

	addr, err := a.instance.Malloc(int32(length))
	if err != nil {
		return WasmResultInternalFailure.Int32()
	}

	err = a.instance.PutMemory(addr, uint64(length), buf.Bytes())
	if err != nil {
		return WasmResultInternalFailure.Int32()
	}

	_ = a.instance.PutUint32(uint64(returnBufferData), uint32(addr))
	_ = a.instance.PutUint32(uint64(returnBufferSize), uint32(length))

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxySetBufferBytes(bufferType int32, start int32, length int32, dataPtr int32, dataSize int32) int32 {
	if BufferType(bufferType) > BufferTypeMax {
		return WasmResultBadArgument.Int32()
	}

	buf := a.GetBuffer(BufferType(bufferType))
	if buf == nil {
		return WasmResultNotFound.Int32()
	}

	content, err := a.instance.GetMemory(uint64(dataPtr), uint64(dataSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	if start + length > int32(buf.Len()) {
		length = int32(buf.Len()) - start
	}

	copy(buf.Bytes()[start:start+length], content)

	return WasmResultOk.Int32()
}

func (a *abiImpl) GetBuffer(bufferType BufferType) buffer.IoBuffer {
	switch bufferType {
	case BufferTypeHttpRequestBody:
		return a.callback.GetHttpRequestBody()
	case BufferTypeHttpResponseBody:
		return a.callback.GetHttpResponseBody()
	case BufferTypePluginConfiguration:
		return a.callback.GetPluginConfig()
	case BufferTypeVmConfiguration:
		return a.callback.GetVmConfig()
	}
	return nil
}

func (a *abiImpl) GetMap(mapType MapType) api.HeaderMap {
	switch mapType {
	case MapTypeHttpRequestHeaders:
		return a.callback.GetHttpRequestHeader()
	case MapTypeHttpRequestTrailers:
		return a.callback.GetHttpRequestTrailer()
	case MapTypeHttpResponseHeaders:
		return a.callback.GetHttpResponseHeader()
	case MapTypeHttpResponseTrailers:
		return a.callback.GetHttpResponseTrailer()
	}
	return nil
}

func (a *abiImpl) proxyGetHeaderMapPairs(mapType int32, returnDataPtr int32, returnDataSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	header := a.GetMap(MapType(mapType))
	if header == nil {
		return WasmResultNotFound.Int32()
	}

	cloneMap := make(map[string]string)
	totalBytesLen := 4
	header.Range(func(key, value string) bool {
		cloneMap[key] = value
		totalBytesLen += 4 + 4 // keyLen + valueLen
		totalBytesLen += len(key) + 1 + len(value) + 1 // key + \0 + value + \0
		return true
	})

	addr, err := a.instance.Malloc(int32(totalBytesLen))
	if err != nil {
		log.DefaultLogger.Errorf("wasm malloc error: %v", err)
		return WasmResultInternalFailure.Int32()
	}

	err = a.instance.PutUint32(addr, uint32(len(cloneMap)))

	lenPtr := addr + 4
	dataPtr := lenPtr + uint64(8 * len(cloneMap))
	for k, v := range cloneMap {
		err = a.instance.PutUint32(lenPtr, uint32(len(k)))
		lenPtr += 4
		err = a.instance.PutUint32(lenPtr, uint32(len(v)))
		lenPtr += 4

		err = a.instance.PutMemory(dataPtr, uint64(len(k)), []byte(k))
		dataPtr += uint64(len(k))
		err = a.instance.PutByte(dataPtr, 0)
		dataPtr++

		err = a.instance.PutMemory(dataPtr, uint64(len(v)), []byte(v))
		dataPtr += uint64(len(v))
		err = a.instance.PutByte(dataPtr, 0)
		dataPtr++
	}

	err = a.instance.PutUint32(uint64(returnDataPtr), uint32(addr))
	err = a.instance.PutUint32(uint64(returnDataSize), uint32(totalBytesLen))

	return WasmResultOk.Int32()
}

// unmarshal map from rawData
func (a *abiImpl) unmarshalMap(rawData []byte) map[string]string {
	if len(rawData) < 4 {
		return nil
	}
	res := make(map[string]string)
	headerSize := binary.LittleEndian.Uint32(rawData[0:4])
	p := 4 + (4+4)*headerSize // headerSize + (key1_size + value1_size) * headerSize
	if int(p) >= len(rawData) {
		return nil
	}
	for i := 0; i < int(headerSize); i++ {
		lenIndex := 4 + (4+4)*i
		keySize := binary.LittleEndian.Uint32(rawData[lenIndex : lenIndex+4])
		valueSize := binary.LittleEndian.Uint32(rawData[lenIndex+4 : lenIndex+8])
		key := string(rawData[p : p+keySize])
		p += keySize
		p++ // 0
		value := string(rawData[p : p+valueSize])
		p += valueSize
		p++ // 0
		res[key] = value
	}
	return res
}

func (a *abiImpl) proxySetHeaderMapPairs(mapType int32, ptr int32, size int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := a.GetMap(MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	newMapContent, err := a.instance.GetMemory(uint64(ptr), uint64(size))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	newMap := a.unmarshalMap(newMapContent)

	for k, v := range newMap {
		headerMap.Set(k, v)
	}

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxyGetHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := a.GetMap(MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := a.instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	value, ok := headerMap.Get(string(key))
	if !ok {
		return WasmResultNotFound.Int32()
	}

	addr, err := a.instance.Malloc(int32(len(value)))
	if err != nil {
		return WasmResultInternalFailure.Int32()
	}

	err = a.instance.PutMemory(addr, uint64(len(value)), []byte(value))

	err = a.instance.PutUint32(uint64(valueDataPtr), uint32(addr))
	err = a.instance.PutUint32(uint64(valueSize), uint32(len(value)))

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxyReplaceHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := a.GetMap(MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := a.instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	value, err := a.instance.GetMemory(uint64(valueDataPtr), uint64(valueSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(value) == 0 {
		return WasmResultBadArgument.Int32()
	}

	headerMap.Set(string(key), string(value))

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxyAddHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := a.GetMap(MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := a.instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	value, err := a.instance.GetMemory(uint64(valueDataPtr), uint64(valueSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(value) == 0 {
		return WasmResultBadArgument.Int32()
	}

	headerMap.Add(string(key), string(value))

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxyRemoveHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := a.GetMap(MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := a.instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	headerMap.Del(string(key))

	return WasmResultOk.Int32()
}

func (a *abiImpl) proxyGetProperty(keyPtr int32, keySize int32, returnValueData int32, returnValueSize int32) int32 {
	return WasmResultOk.Int32()
}

func (a *abiImpl) proxySetProperty(keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxySetEffectiveContext(contextID int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxySetTickPeriodMilliseconds(tickPeriodMilliseconds int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGetCurrentTimeNanoseconds(resultUint64Ptr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGrpcCall(servicePtr int32, serviceSize int32, serviceNamePtr int32, serviceNameSize int32,
	methodNamePtr int32, methodNameSize int32,
	initialMetadataPtr int32, initialMetadataSize int32,
	requestPtr int32, requestSize int32,
	timeoutMilliseconds int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGrpcStream(servicePtr int32, serviceSize int32, serviceNamePtr int32, serviceNameSize int32,
	methodNamePtr int32, methodNameSize int32,
	initialMetadataPtr int32, initialMetadataSize int32,
	tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGrpcCancel(token int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGrpcClose(token int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGrpcSend(token int32, messagePtr int32, messageSize int32, endStream int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyHttpCall(uriPtr int32, uriSize int32,
	headerPairsPtr int32, headerPairsSize int32,
	bodyPtr int32, bodySize int32,
	trailerPairsPtr int32, trailerPairsSize int32,
	timeoutMilliseconds int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyDefineMetric(metricType int32, namePtr int32, nameSize int32, resultPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyIncrementMetric(metricId int32, offset int64) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyRecordMetric(metricId int32, value int64) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGetMetric(metricId int32, resultUint64Ptr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyRegisterSharedQueue(queueNamePtr int32, queueNameSize int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyResolveSharedQueue(vmIdPtr int32, vmIdSize int32, queueNamePtr int32, queueNameSize int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyDequeueSharedQueue(token int32, dataPtr int32, dataSize int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyEnqueueSharedQueue(token int32, dataPtr int32, dataSize int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxyGetSharedData(keyPtr int32, keySize int32, valuePtr int32, valueSizePtr int32, casPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func (a *abiImpl) proxySetSharedData(keyPtr int32, keySize int32, valuePtr int32, valueSize int32, cas int32) int32 {
	return WasmResultUnimplemented.Int32()
}
