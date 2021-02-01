/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proxywasm_0_1_0

import (
	"encoding/binary"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
)

func getInstanceCallback(instance types.WasmInstance) InstanceCallback {
	v := instance.GetData()
	if v == nil {
		return &DefaultInstanceCallback{}
	}

	cb, ok := v.(InstanceCallback)
	if !ok {
		return &DefaultInstanceCallback{}
	}

	return cb
}

func proxyLog(instance types.WasmInstance, level int32, logDataPtr int32, logDataSize int32) int32 {
	callback := getInstanceCallback(instance)

	logContent, err := instance.GetMemory(uint64(logDataPtr), uint64(logDataSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	callback.Log(toMosnLogLevel(LogLevel(level)), string(logContent))

	return WasmResultOk.Int32()
}

func proxyGetBufferBytes(instance types.WasmInstance, bufferType int32, start int32, length int32, returnBufferData int32, returnBufferSize int32) int32 {
	if BufferType(bufferType) > BufferTypeMax {
		return WasmResultBadArgument.Int32()
	}

	buf := GetBuffer(instance, BufferType(bufferType))
	if buf == nil {
		return WasmResultNotFound.Int32()
	}

	if start > start+length {
		return WasmResultBadArgument.Int32()
	}

	if start+length > int32(buf.Len()) {
		length = int32(buf.Len()) - start
	}

	addr, err := instance.Malloc(int32(length))
	if err != nil {
		return WasmResultInternalFailure.Int32()
	}

	err = instance.PutMemory(addr, uint64(length), buf.Bytes())
	if err != nil {
		return WasmResultInternalFailure.Int32()
	}

	_ = instance.PutUint32(uint64(returnBufferData), uint32(addr))
	_ = instance.PutUint32(uint64(returnBufferSize), uint32(length))

	return WasmResultOk.Int32()
}

func proxySetBufferBytes(instance types.WasmInstance, bufferType int32, start int32, length int32, dataPtr int32, dataSize int32) int32 {
	if BufferType(bufferType) > BufferTypeMax {
		return WasmResultBadArgument.Int32()
	}

	buf := GetBuffer(instance, BufferType(bufferType))
	if buf == nil {
		return WasmResultNotFound.Int32()
	}

	content, err := instance.GetMemory(uint64(dataPtr), uint64(dataSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	if start+length > int32(buf.Len()) {
		length = int32(buf.Len()) - start
	}

	copy(buf.Bytes()[start:start+length], content)

	return WasmResultOk.Int32()
}

func GetBuffer(instance types.WasmInstance, bufferType BufferType) buffer.IoBuffer {
	callback := getInstanceCallback(instance)

	switch bufferType {
	case BufferTypeHttpRequestBody:
		return callback.GetHttpRequestBody()
	case BufferTypeHttpResponseBody:
		return callback.GetHttpResponseBody()
	case BufferTypePluginConfiguration:
		return callback.GetPluginConfig()
	case BufferTypeVmConfiguration:
		return callback.GetVmConfig()
	}
	return nil
}

func GetMap(instance types.WasmInstance, mapType MapType) api.HeaderMap {
	callback := getInstanceCallback(instance)

	switch mapType {
	case MapTypeHttpRequestHeaders:
		return callback.GetHttpRequestHeader()
	case MapTypeHttpRequestTrailers:
		return callback.GetHttpRequestTrailer()
	case MapTypeHttpResponseHeaders:
		return callback.GetHttpResponseHeader()
	case MapTypeHttpResponseTrailers:
		return callback.GetHttpResponseTrailer()
	}
	return nil
}

func proxyGetHeaderMapPairs(instance types.WasmInstance, mapType int32, returnDataPtr int32, returnDataSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	header := GetMap(instance, MapType(mapType))
	if header == nil {
		return WasmResultNotFound.Int32()
	}

	cloneMap := make(map[string]string)
	totalBytesLen := 4
	header.Range(func(key, value string) bool {
		cloneMap[key] = value
		totalBytesLen += 4 + 4                         // keyLen + valueLen
		totalBytesLen += len(key) + 1 + len(value) + 1 // key + \0 + value + \0
		return true
	})

	addr, err := instance.Malloc(int32(totalBytesLen))
	if err != nil {
		log.DefaultLogger.Errorf("wasm malloc error: %v", err)
		return WasmResultInternalFailure.Int32()
	}

	err = instance.PutUint32(addr, uint32(len(cloneMap)))

	lenPtr := addr + 4
	dataPtr := lenPtr + uint64(8*len(cloneMap))
	for k, v := range cloneMap {
		err = instance.PutUint32(lenPtr, uint32(len(k)))
		lenPtr += 4
		err = instance.PutUint32(lenPtr, uint32(len(v)))
		lenPtr += 4

		err = instance.PutMemory(dataPtr, uint64(len(k)), []byte(k))
		dataPtr += uint64(len(k))
		err = instance.PutByte(dataPtr, 0)
		dataPtr++

		err = instance.PutMemory(dataPtr, uint64(len(v)), []byte(v))
		dataPtr += uint64(len(v))
		err = instance.PutByte(dataPtr, 0)
		dataPtr++
	}

	err = instance.PutUint32(uint64(returnDataPtr), uint32(addr))
	err = instance.PutUint32(uint64(returnDataSize), uint32(totalBytesLen))

	return WasmResultOk.Int32()
}

// unmarshal map from rawData
func unmarshalMap(rawData []byte) map[string]string {
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

func proxySetHeaderMapPairs(instance types.WasmInstance, mapType int32, ptr int32, size int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := GetMap(instance, MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	newMapContent, err := instance.GetMemory(uint64(ptr), uint64(size))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	newMap := unmarshalMap(newMapContent)

	for k, v := range newMap {
		headerMap.Set(k, v)
	}

	return WasmResultOk.Int32()
}

func proxyGetHeaderMapValue(instance types.WasmInstance, mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := GetMap(instance, MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
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

	addr, err := instance.Malloc(int32(len(value)))
	if err != nil {
		return WasmResultInternalFailure.Int32()
	}

	err = instance.PutMemory(addr, uint64(len(value)), []byte(value))

	err = instance.PutUint32(uint64(valueDataPtr), uint32(addr))
	err = instance.PutUint32(uint64(valueSize), uint32(len(value)))

	return WasmResultOk.Int32()
}

func proxyReplaceHeaderMapValue(instance types.WasmInstance, mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := GetMap(instance, MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	value, err := instance.GetMemory(uint64(valueDataPtr), uint64(valueSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(value) == 0 {
		return WasmResultBadArgument.Int32()
	}

	headerMap.Set(string(key), string(value))

	return WasmResultOk.Int32()
}

func proxyAddHeaderMapValue(instance types.WasmInstance, mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := GetMap(instance, MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	value, err := instance.GetMemory(uint64(valueDataPtr), uint64(valueSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(value) == 0 {
		return WasmResultBadArgument.Int32()
	}

	headerMap.Add(string(key), string(value))

	return WasmResultOk.Int32()
}

func proxyRemoveHeaderMapValue(instance types.WasmInstance, mapType int32, keyDataPtr int32, keySize int32) int32 {
	if MapType(mapType) > MapTypeMax {
		return WasmResultBadArgument.Int32()
	}

	headerMap := GetMap(instance, MapType(mapType))
	if headerMap == nil {
		return WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return WasmResultBadArgument.Int32()
	}

	headerMap.Del(string(key))

	return WasmResultOk.Int32()
}

func proxyGetProperty(instance types.WasmInstance, keyPtr int32, keySize int32, returnValueData int32, returnValueSize int32) int32 {
	return WasmResultOk.Int32()
}

func proxySetProperty(instance types.WasmInstance, keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxySetEffectiveContext(instance types.WasmInstance, contextID int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxySetTickPeriodMilliseconds(instance types.WasmInstance, tickPeriodMilliseconds int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGetCurrentTimeNanoseconds(instance types.WasmInstance, resultUint64Ptr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGrpcCall(instance types.WasmInstance, servicePtr int32, serviceSize int32, serviceNamePtr int32, serviceNameSize int32,
	methodNamePtr int32, methodNameSize int32,
	initialMetadataPtr int32, initialMetadataSize int32,
	requestPtr int32, requestSize int32,
	timeoutMilliseconds int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGrpcStream(instance types.WasmInstance, servicePtr int32, serviceSize int32, serviceNamePtr int32, serviceNameSize int32,
	methodNamePtr int32, methodNameSize int32,
	initialMetadataPtr int32, initialMetadataSize int32,
	tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGrpcCancel(instance types.WasmInstance, token int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGrpcClose(instance types.WasmInstance, token int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGrpcSend(instance types.WasmInstance, token int32, messagePtr int32, messageSize int32, endStream int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyHttpCall(instance types.WasmInstance, uriPtr int32, uriSize int32,
	headerPairsPtr int32, headerPairsSize int32,
	bodyPtr int32, bodySize int32,
	trailerPairsPtr int32, trailerPairsSize int32,
	timeoutMilliseconds int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyDefineMetric(instance types.WasmInstance, metricType int32, namePtr int32, nameSize int32, resultPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyIncrementMetric(instance types.WasmInstance, metricId int32, offset int64) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyRecordMetric(instance types.WasmInstance, metricId int32, value int64) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGetMetric(instance types.WasmInstance, metricId int32, resultUint64Ptr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyRegisterSharedQueue(instance types.WasmInstance, queueNamePtr int32, queueNameSize int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyResolveSharedQueue(instance types.WasmInstance, vmIdPtr int32, vmIdSize int32, queueNamePtr int32, queueNameSize int32, tokenPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyDequeueSharedQueue(instance types.WasmInstance, token int32, dataPtr int32, dataSize int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyEnqueueSharedQueue(instance types.WasmInstance, token int32, dataPtr int32, dataSize int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxyGetSharedData(instance types.WasmInstance, keyPtr int32, keySize int32, valuePtr int32, valueSizePtr int32, casPtr int32) int32 {
	return WasmResultUnimplemented.Int32()
}

func proxySetSharedData(instance types.WasmInstance, keyPtr int32, keySize int32, valuePtr int32, valueSize int32, cas int32) int32 {
	return WasmResultUnimplemented.Int32()
}
