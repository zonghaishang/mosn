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

package wasm

import (
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
	v1 "mosn.io/mosn/pkg/wasm/abi/proxywasm_0_1_0"
)

func init() {
	abi.RegisterABI("proxy_abi_version_0_2_0", abiImplFactory)
}

func abiImplFactory() abi.ABI {
	return &AbiV2Impl{}
}

// easy for extension
type AbiV2Impl struct {
	v1.AbiImpl
}

func (a *AbiV2Impl) OnInstanceCreate(instance types.WasmInstance) {

	// pre register abi 0_1_0 version
	a.AbiImpl.OnInstanceCreate(instance)

	instance.RegisterFunc("env", "proxy_set_buffer_bytes", proxySetBufferBytes)

	return
}
