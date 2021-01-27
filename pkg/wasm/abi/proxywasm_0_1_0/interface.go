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
