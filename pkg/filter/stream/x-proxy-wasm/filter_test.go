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

package x_proxy_wasm

import (
	"context"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/pkg/buffer"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/mock"
)

func TestCreateProxyWasmFilterFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filterConfigMap := map[string]interface{}{
		"vm_config": v2.WasmVmConfig{
			Engine: "wasmer",
		},
		"user_config": "user_value",
	}

	factory, err := CreateProxyWasmFilterFactory(filterConfigMap)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	addReceiverFilterCount := 0
	addSenderFilterCount := 0

	cb := mock.NewMockStreamFilterChainFactoryCallbacks(ctrl)
	cb.EXPECT().AddStreamReceiverFilter(gomock.Any(), gomock.Any()).Do(func(api.StreamReceiverFilter, api.ReceiverFilterPhase) {
		addReceiverFilterCount++
	}).AnyTimes()
	cb.EXPECT().AddStreamSenderFilter(gomock.Any(), gomock.Any()).Do(func(api.StreamSenderFilter, api.SenderFilterPhase) {
		addSenderFilterCount++
	}).AnyTimes()

	factory.CreateFilterChain(context.TODO(), cb)
	assert.Equal(t, addReceiverFilterCount, 1)
	assert.Equal(t, addSenderFilterCount, 1)
}

func TestOnReciever(t *testing.T) {
	wasmConfig := v2.WasmPluginConfig{
		PluginName: "test_plugin",
		VmConfig: v2.WasmVmConfig{
			Engine: "wasmer",
			Path:   "test.wasm",
		},
		InstanceNum: 2,
	}
	wasm.GetWasmManager().AddOrUpdateWasm(wasmConfig)

	filter := NewFilter(context.TODO(), "test_plugin", v2.WasmVmConfig{}, pluginConfig{})
	filter.OnReceive(context.TODO(), protocol.CommonHeader{}, buffer.NewIoBuffer(1), nil)
}
