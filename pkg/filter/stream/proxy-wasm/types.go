package proxywasm

type MapType int32

const (
	MapTypeHttpRequestHeaders       MapType = 0
	MapTypeHttpRequestTrailers      MapType = 1
	MapTypeHttpResponseHeaders      MapType = 2
	MapTypeHttpResponseTrailers     MapType = 3
	MapTypeHttpCallResponseHeaders  MapType = 7
	MapTypeHttpCallResponseTrailers MapType = 8
	MapTypeMax                      MapType = 8
)

type BufferType int32

const (
	BufferTypeHttpRequestBody     BufferType = 0
	BufferTypeHttpResponseBody    BufferType = 1
	BufferTypeVmConfiguration     BufferType = 6
	BufferTypePluginConfiguration BufferType = 7
	BufferTypeMax                 BufferType = 7
)

type MetricType int32

const (
	MetricTypeCounter   MetricType = 0
	MetricTypeGauge     MetricType = 1
	MetricTypeHistogram MetricType = 2
	MetricTypeMax       MetricType = 2
)

type WasmResult int32

const (
	WasmResultOk                   WasmResult = 0
	WasmResultNotFound             WasmResult = 1  // The result could not be found, e.g. a provided key did not appear in a table.
	WasmResultBadArgument          WasmResult = 2  // An argument was bad, e.g. did not not conform to the required range.
	WasmResultSerializationFailure WasmResult = 3  // A protobuf could not be serialized.
	WasmResultParseFailure         WasmResult = 4  // A protobuf could not be parsed.
	WasmResultBadExpression        WasmResult = 5  // A provided expression (e.g. "foo.bar") was illegal or unrecognized.
	WasmResultInvalidMemoryAccess  WasmResult = 6  // A provided memory range was not legal.
	WasmResultEmpty                WasmResult = 7  // Data was requested from an empty container.
	WasmResultCasMismatch          WasmResult = 8  // The provided CAS did not match that of the stored data.
	WasmResultResultMismatch       WasmResult = 9  // Returned result was unexpected, e.g. of the incorrect size.
	WasmResultInternalFailure      WasmResult = 10 // Internal failure: trying check logs of the surrounding system.
	WasmResultBrokenConnection     WasmResult = 11 // The connection/stream/pipe was broken/closed unexpectedly.
	WasmResultUnimplemented        WasmResult = 12 // Feature not implemented.
)

func (wasmResult WasmResult) Int32() int32 {
	return int32(wasmResult)
}
