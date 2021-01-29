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
	"encoding/json"
	"errors"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm_0_1_0"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/utils"
)

const ProxyWasm = "x-proxy-wasm"

func init() {
	api.RegisterStream(ProxyWasm, createProxyWasmFilterFactory)
}

type FilterConfigFactory struct {
	proxywasm_0_1_0.DefaultInstanceCallback

	pluginName string
	config     *filterConfig
}

func createProxyWasmFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	config, err := parseFilterConfig(conf)
	if err != nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] createProxyWasmFilterFactory fail to parse config, err: %v", err)
		return nil, err
	}

	var pluginName string

	if config.FromWasmPlugin == "" {
		pluginName = utils.GenerateUUID()

		v2Config := v2.WasmPluginConfig{
			PluginName:  pluginName,
			VmConfig:    config.VmConfig,
			InstanceNum: config.InstanceNum,
		}

		err = wasm.GetWasmManager().AddOrUpdateWasm(v2Config)
		if err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] createProxyWasmFilterFactory fail to add plugin, err: %v", err)
			return nil, err
		}
	} else {
		pluginName = config.FromWasmPlugin
	}

	pw := wasm.GetWasmManager().GetWasmPluginWrapperByName(pluginName)
	if pw == nil {
		return nil, errors.New("plugin not found")
	}

	if config.FromWasmPlugin == "" {
		config.VmConfig = pw.GetConfig().VmConfig
	}

	factory := &FilterConfigFactory{
		pluginName: pluginName,
		config:     config,
	}

	pw.RegisterPluginHandler(factory)

	return factory, nil
}

func configSize(config interface{}) int {
	b, err := json.Marshal(config)
	if err != nil {
		return 0
	}
	return len(b)
}

func (f *FilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewFilter(context, f.pluginName, f.config.RootContextID)

	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

func (f *FilterConfigFactory) GetVmConfig() buffer.IoBuffer {
	b, err := json.Marshal(f.config.VmConfig)
	if err != nil {
		return nil
	}
	return buffer.NewIoBufferBytes(b)
}

func (f *FilterConfigFactory) GetPluginConfig() buffer.IoBuffer {
	b, err := json.Marshal(f.config.pluginConfig)
	if err != nil {
		return nil
	}
	return buffer.NewIoBufferBytes(b)
}

func (f *FilterConfigFactory) OnConfigUpdate(config v2.WasmPluginConfig) {
	f.config.InstanceNum = config.InstanceNum
	f.config.VmConfig = config.VmConfig
}

func (f *FilterConfigFactory) OnPluginStart(plugin types.WasmPlugin) {
	abiVersion := abi.GetABI("proxy_abi_version_0_1_0")
	if abiVersion == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] NewFilter abi version not found")
		return
	}

	plugin.Exec(func(instanceWrapper types.WasmInstanceWrapper) bool {
		instanceWrapper.Acquire()
		defer instanceWrapper.Release()

		abiVersion.SetInstance(instanceWrapper)
		abiVersion.SetInstanceCallBack(f)

		exports := abiVersion.(proxywasm_0_1_0.Exports)

		_ = exports.ProxyOnContextCreate(f.config.RootContextID, 0)
		_, _ = exports.ProxyOnConfigure(f.config.RootContextID, int32(configSize(f.config.pluginConfig)))
		_, _ = exports.ProxyOnVmStart(f.config.RootContextID, int32(configSize(f.config.VmConfig)))

		return true
	})
}

func (f *FilterConfigFactory) OnPluginDestroy(plugin types.WasmPlugin) {
	return
}
