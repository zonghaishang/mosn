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
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/buffer"
)

// Exports contains ABI that exported by wasm module
type Exports interface {
	ProxyOnContextCreate(contextId int32, parentContextId int32) error
	ProxyOnDone(contextId int32) (int32, error)
	ProxyOnLog(contextId int32) error
	ProxyOnDelete(contextId int32) error

	ProxyOnVmStart(rootContextId int32, vmConfigurationSize int32) (int32, error)
	ProxyOnConfigure(rootContextId int32, pluginConfigurationSize int32) (int32, error)

	ProxyOnTick(rootContextId int32) error

	ProxyOnNewConnection(contextId int32) error
	ProxyOnDownstreamData(contextId int32, dataLength int32, endOfStream int32) (int32, error)
	ProxyOnDownstreamConnectionClose(contextId int32, closeType int32) error
	ProxyOnUpstreamData(contextId int32, dataLength int32, endOfStream int32) (int32, error)
	ProxyOnUpstreamConnectionClose(contextId int32, closeType int32) error

	ProxyOnRequestHeaders(contextId int32, headers int32, endOfStream int32) (int32, error)
	ProxyOnRequestBody(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error)
	ProxyOnRequestTrailers(contextId int32, trailers int32) (int32, error)
	ProxyOnRequestMetadata(contextId int32, nElements int32) (int32, error)

	ProxyOnResponseHeaders(contextId int32, headers int32, endOfStream int32) (int32, error)
	ProxyOnResponseBody(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error)
	ProxyOnResponseTrailers(contextId int32, trailers int32) (int32, error)
	ProxyOnResponseMetadata(contextId int32, nElements int32) (int32, error)

	ProxyOnHttpCallResponse(contextId int32, token int32, headers int32, bodySize int32, trailers int32) error

	ProxyOnGrpcReceiveInitialMetadata(contextId int32, token int32, headers int32) error
	ProxyOnGrpcTrailingMetadata(contextId int32, token int32, trailers int32) error
	ProxyOnGrpcReceive(contextId int32, token int32, responseSize int32) error
	ProxyOnGrpcClose(contextId int32, token int32, statusCode int32) error

	ProxyOnQueueReady(rootContextId int32, token int32) error
}

type InstanceCallback interface {
	GetVmConfig() buffer.IoBuffer
	GetPluginConfig() buffer.IoBuffer

	Log(level log.Level, msg string)

	GetHttpRequestHeader() api.HeaderMap
	GetHttpRequestBody() buffer.IoBuffer
	GetHttpRequestTrailer() api.HeaderMap

	GetHttpResponseHeader() api.HeaderMap
	GetHttpResponseBody() buffer.IoBuffer
	GetHttpResponseTrailer() api.HeaderMap
}

type DefaultInstanceCallback struct{}

func (d *DefaultInstanceCallback) GetVmConfig() buffer.IoBuffer {
	return nil
}

func (d *DefaultInstanceCallback) GetPluginConfig() buffer.IoBuffer {
	return nil
}

func (d *DefaultInstanceCallback) Log(level log.Level, msg string) {
	return
}

func (d *DefaultInstanceCallback) GetHttpRequestHeader() api.HeaderMap {
	return nil
}

func (d *DefaultInstanceCallback) GetHttpRequestBody() buffer.IoBuffer {
	return nil
}

func (d *DefaultInstanceCallback) GetHttpRequestTrailer() api.HeaderMap {
	return nil
}

func (d *DefaultInstanceCallback) GetHttpResponseHeader() api.HeaderMap {
	return nil
}

func (d *DefaultInstanceCallback) GetHttpResponseBody() buffer.IoBuffer {
	return nil
}

func (d *DefaultInstanceCallback) GetHttpResponseTrailer() api.HeaderMap {
	return nil
}
