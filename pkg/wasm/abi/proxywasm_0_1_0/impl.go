package proxywasm_0_1_0

import (
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
)

func init() {
	abi.RegisterABI("proxy_abi_version_0_1_0", &abiImpl{})
}

type abiImpl struct {
	instance types.WasmInstance
	callback InstanceCallback
}

func (a *abiImpl) SetInstance(instance types.WasmInstance) {
	a.instance = instance
}

func (a *abiImpl) SetInstanceCallBack(callback interface{}) {
	cb, ok := callback.(InstanceCallback)
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][impl] SetInstanceCallBack invalid callback type")
		return
	}
	a.callback = cb
}

func (a *abiImpl) OnInstanceCreate(instance types.WasmInstance) {

	instance.RegisterFunc("env", "proxy_log", a.proxyLog)

	instance.RegisterFunc("env", "proxy_set_effective_context", a.proxySetEffectiveContext)

	instance.RegisterFunc("env", "proxy_get_property", a.proxyGetProperty)
	instance.RegisterFunc("env", "proxy_set_property", a.proxySetProperty)

	instance.RegisterFunc("env", "proxy_get_buffer_bytes", a.proxyGetBufferBytes)
	instance.RegisterFunc("env", "proxy_set_buffer_bytes", a.proxySetBufferBytes)

	instance.RegisterFunc("env", "proxy_get_header_map_pairs", a.proxyGetHeaderMapPairs)
	instance.RegisterFunc("env", "proxy_set_header_map_pairs", a.proxySetHeaderMapPairs)

	instance.RegisterFunc("env", "proxy_get_header_map_value", a.proxyGetHeaderMapValue)
	instance.RegisterFunc("env", "proxy_replace_header_map_value", a.proxyReplaceHeaderMapValue)
	instance.RegisterFunc("env", "proxy_add_header_map_value", a.proxyAddHeaderMapValue)
	instance.RegisterFunc("env", "proxy_remove_header_map_value", a.proxyRemoveHeaderMapValue)

	instance.RegisterFunc("env", "proxy_set_tick_period_milliseconds", a.proxySetTickPeriodMilliseconds)
	instance.RegisterFunc("env", "proxy_get_current_time_nanoseconds", a.proxyGetCurrentTimeNanoseconds)

	instance.RegisterFunc("env", "proxy_grpc_call", a.proxyGrpcCall)
	instance.RegisterFunc("env", "proxy_grpc_stream", a.proxyGrpcStream)
	instance.RegisterFunc("env", "proxy_grpc_cancel", a.proxyGrpcCancel)
	instance.RegisterFunc("env", "proxy_grpc_close", a.proxyGrpcClose)
	instance.RegisterFunc("env", "proxy_grpc_send", a.proxyGrpcSend)

	instance.RegisterFunc("env", "proxy_http_call", a.proxyHttpCall)

	instance.RegisterFunc("env", "proxy_define_metric", a.proxyDefineMetric)
	instance.RegisterFunc("env", "proxy_increment_metric", a.proxyIncrementMetric)
	instance.RegisterFunc("env", "proxy_record_metric", a.proxyRecordMetric)
	instance.RegisterFunc("env", "proxy_get_metric", a.proxyGetMetric)

	instance.RegisterFunc("env", "proxy_register_shared_queue", a.proxyRegisterSharedQueue)
	instance.RegisterFunc("env", "proxy_resolve_shared_queue", a.proxyResolveSharedQueue)
	instance.RegisterFunc("env", "proxy_dequeue_shared_queue", a.proxyDequeueSharedQueue)
	instance.RegisterFunc("env", "proxy_enqueue_shared_queue", a.proxyEnqueueSharedQueue)

	instance.RegisterFunc("env", "proxy_get_shared_data", a.proxyGetSharedData)
	instance.RegisterFunc("env", "proxy_set_shared_data", a.proxySetSharedData)

	a.SetInstance(instance)

	return
}

func (a *abiImpl) OnStart(instance types.WasmInstance) {
	err := a.ProxyOnContextCreate(int32(100), int32(0))
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][impl] new Wasmer Instance, fail to call ProxyOnContextCreate(), err: %v", err)
		return
	}

	_, err = a.ProxyOnVmStart(int32(100), int32(0))
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][impl] new Wasmer Instance, fail to call ProxyOnVMStart(), err: %v", err)
		return
	}
}

func (a *abiImpl) OnInstanceDestroy(instance types.WasmInstance) {
	return
}