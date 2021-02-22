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
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
)

func init() {
	abi.RegisterABI("proxy_abi_version_0_1_0", abiImplFactory)
}

func abiImplFactory() abi.ABI {
	return &AbiImpl{}
}

// easy for extension
type AbiImpl struct {
	instance types.WasmInstance
}

func (a *AbiImpl) SetInstance(instance types.WasmInstance) {
	a.instance = instance
}

func (a *AbiImpl) GetInstance() types.WasmInstance {
	return a.instance
}

func (a *AbiImpl) OnStart(instance types.WasmInstance) {
	return
}

func (a *AbiImpl) OnInstanceDestroy(instance types.WasmInstance) {
	return
}

func (a *AbiImpl) OnInstanceCreate(instance types.WasmInstance) {
	instance.RegisterFunc("env", "proxy_log", proxyLog)

	instance.RegisterFunc("env", "proxy_set_effective_context", proxySetEffectiveContext)

	instance.RegisterFunc("env", "proxy_get_property", proxyGetProperty)
	instance.RegisterFunc("env", "proxy_set_property", proxySetProperty)

	instance.RegisterFunc("env", "proxy_get_buffer_bytes", ProxyGetBufferBytes)
	instance.RegisterFunc("env", "proxy_set_buffer_bytes", ProxySetBufferBytes)

	instance.RegisterFunc("env", "proxy_get_header_map_pairs", proxyGetHeaderMapPairs)
	instance.RegisterFunc("env", "proxy_set_header_map_pairs", proxySetHeaderMapPairs)

	instance.RegisterFunc("env", "proxy_get_header_map_value", proxyGetHeaderMapValue)
	instance.RegisterFunc("env", "proxy_replace_header_map_value", proxyReplaceHeaderMapValue)
	instance.RegisterFunc("env", "proxy_add_header_map_value", proxyAddHeaderMapValue)
	instance.RegisterFunc("env", "proxy_remove_header_map_value", proxyRemoveHeaderMapValue)

	instance.RegisterFunc("env", "proxy_set_tick_period_milliseconds", proxySetTickPeriodMilliseconds)
	instance.RegisterFunc("env", "proxy_get_current_time_nanoseconds", proxyGetCurrentTimeNanoseconds)

	instance.RegisterFunc("env", "proxy_grpc_call", proxyGrpcCall)
	instance.RegisterFunc("env", "proxy_grpc_stream", proxyGrpcStream)
	instance.RegisterFunc("env", "proxy_grpc_cancel", proxyGrpcCancel)
	instance.RegisterFunc("env", "proxy_grpc_close", proxyGrpcClose)
	instance.RegisterFunc("env", "proxy_grpc_send", proxyGrpcSend)

	instance.RegisterFunc("env", "proxy_http_call", proxyHttpCall)

	instance.RegisterFunc("env", "proxy_define_metric", proxyDefineMetric)
	instance.RegisterFunc("env", "proxy_increment_metric", proxyIncrementMetric)
	instance.RegisterFunc("env", "proxy_record_metric", proxyRecordMetric)
	instance.RegisterFunc("env", "proxy_get_metric", proxyGetMetric)

	instance.RegisterFunc("env", "proxy_register_shared_queue", proxyRegisterSharedQueue)
	instance.RegisterFunc("env", "proxy_resolve_shared_queue", proxyResolveSharedQueue)
	instance.RegisterFunc("env", "proxy_dequeue_shared_queue", proxyDequeueSharedQueue)
	instance.RegisterFunc("env", "proxy_enqueue_shared_queue", proxyEnqueueSharedQueue)

	instance.RegisterFunc("env", "proxy_get_shared_data", proxyGetSharedData)
	instance.RegisterFunc("env", "proxy_set_shared_data", proxySetSharedData)

	return
}
