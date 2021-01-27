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

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/utils"
)

func init() {
	api.RegisterStream(ProxyWasm, CreateProxyWasmFilterFactory)
}

const ProxyWasm = "x-proxy-wasm"

type filterConfig struct {
	FromWasmPlugin string          `json:"from_wasm_plugin,omitempty"`
	VmConfig       v2.WasmVmConfig `json:"vm_config,omitempty"`
	pluginConfig
}

type pluginConfig struct {
	UserConfig string `json:"user_config,omitempty"`
}

type FilterConfigFactory struct {
	pluginName   string
	vmConfig     v2.WasmVmConfig
	pluginConfig pluginConfig
}

func CreateProxyWasmFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	config, err := ParseStreamProxyWasmFilter(conf)
	if err != nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] CreateProxyWasmFilterFactory fail to parse config, err: %v", err)
		return nil, err
	}

	var pluginName string

	if config.FromWasmPlugin == "" {
		pluginName = utils.GenerateUUID()
		v2Config := v2.WasmPluginConfig{
			PluginName: pluginName,
			VmConfig:   config.VmConfig,
		}
		err = wasm.GetWasmManager().AddOrUpdateWasm(v2Config)
		if err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] CreateProxyWasmFilterFactory fail to add plugin, err: %v", err)
			return nil, err
		}
	} else {
		pluginName = config.FromWasmPlugin
	}

	return &FilterConfigFactory{
		pluginName:   pluginName,
		vmConfig:     config.VmConfig,
		pluginConfig: config.pluginConfig,
	}, nil
}

func ParseStreamProxyWasmFilter(cfg map[string]interface{}) (*filterConfig, error) {
	filterConfig := &filterConfig{}
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, filterConfig); err != nil {
		return nil, err
	}

	return filterConfig, nil
}

func (f *FilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewFilter(context, f.pluginName, f.vmConfig, f.pluginConfig)
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type Filter struct {
	ctx context.Context

	pluginName   string
	vmConfig     v2.WasmVmConfig
	pluginConfig pluginConfig

	// TODO
	contextId int32

	receiverFilterHandler api.StreamReceiverFilterHandler
	senderFilterHandler   api.StreamSenderFilterHandler
}

func NewFilter(ctx context.Context, pluginName string, vmConfig v2.WasmVmConfig, pluginConfig pluginConfig) *Filter {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[x-proxy-wasm][filter] create a new proxy wasm filter")
	}

	filter := &Filter{
		ctx:          ctx,
		pluginName:   pluginName,
		vmConfig:     vmConfig,
		pluginConfig: pluginConfig,
	}

	return filter
}

func (f *Filter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.receiverFilterHandler = handler
}

func (f *Filter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.senderFilterHandler = handler
}

func (f *Filter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[x-proxy-wasm][filter] OnReceive, id = %d", f.contextId)
	}

	// get wasm plugin instance
	pluginWrapper := wasm.GetWasmManager().GetWasmPluginWrapperByName(f.pluginName)
	if pluginWrapper == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive wasm plugin not exists, plugin name: %v", f.pluginName)
		return api.StreamFilterContinue
	}

	plugin := pluginWrapper.GetPlugin()

	instance := plugin.GetInstance()
	defer plugin.ReleaseInstance(instance)

	if instance == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive fail to get wasm instance, plugin name: %v", f.pluginName)
		return api.StreamFilterContinue
	}

	instance.SetInstanceCallback(f)

	// do filter
	if buf != nil && buf.Len() > 0 {
		if _, err := instance.ProxyOnRequestHeaders(f.contextId, 0, 0); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestHeaders err: %v", err)
		}
		if _, err := instance.ProxyOnRequestBody(f.contextId, int32(buf.Len()), 1); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestBody err: %v", err)
		}
	} else {
		if _, err := instance.ProxyOnRequestHeaders(f.contextId, 0, 1); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestHeaders err: %v", err)
		}
	}

	return api.StreamFilterContinue
}

func (f *Filter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[x-proxy-wasm][filter] Append, id = %d", f.contextId)
	}

	// get wasm plugin instance
	pluginWrapper := wasm.GetWasmManager().GetWasmPluginWrapperByName(f.pluginName)
	if pluginWrapper == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] Append wasm plugin not exists, plugin name: %v", f.pluginName)
		return api.StreamFilterContinue
	}

	plugin := pluginWrapper.GetPlugin()

	instance := plugin.GetInstance()
	defer plugin.ReleaseInstance(instance)

	instance.SetInstanceCallback(f)

	if instance == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] Append fail to get wasm instance, plugin name: %v", f.pluginName)
		return api.StreamFilterContinue
	}

	// do filter
	if _, err := instance.ProxyOnResponseHeaders(f.contextId, 1, 0); err != nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] Append ProxyOnResponseHeaders err: %v", err)
	}

	return api.StreamFilterContinue
}

func (f *Filter) OnDestroy() {
}

func (f *Filter) GetVmConfig() []byte {
	bytes, err := json.Marshal(f.vmConfig)
	if err != nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] GetVmConfig fail to marshal vm config, err: %v", err)
		return nil
	}

	return bytes
}

func (f *Filter) GetPluginConfig() []byte {
	bytes, err := json.Marshal(f.pluginConfig)
	if err != nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] GetPluginConfig fail to marshal plugin config, err: %v", err)
		return nil
	}

	return bytes
}

func (f *Filter) Log(msg string) {
	log.DefaultLogger.Errorf(msg)
}

func (f *Filter) GetHttpRequestHeader() api.HeaderMap {
	if f.receiverFilterHandler == nil {
		return nil
	}
	return f.receiverFilterHandler.GetRequestHeaders()
}

func (f *Filter) SetHttpRequestHeader(header api.HeaderMap) {
	if f.receiverFilterHandler == nil {
		return
	}
	f.receiverFilterHandler.SetRequestHeaders(header)
}

func (f *Filter) GetHttpRequestBody() buffer.IoBuffer {
	if f.receiverFilterHandler == nil {
		return nil
	}
	return f.receiverFilterHandler.GetRequestData()
}

func (f *Filter) SetHttpRequestBody(body buffer.IoBuffer) {
	if f.receiverFilterHandler == nil {
		return
	}
	f.receiverFilterHandler.SetRequestData(body)
}

func (f *Filter) GetHttpRequestTrailer() api.HeaderMap {
	if f.receiverFilterHandler == nil {
		return nil
	}
	return f.receiverFilterHandler.GetRequestTrailers()
}

func (f *Filter) SetHttpRequestTrailer(trailer api.HeaderMap) {
	if f.receiverFilterHandler == nil {
		return
	}
	f.receiverFilterHandler.SetRequestTrailers(trailer)
}

func (f *Filter) GetHttpResponseHeader() api.HeaderMap {
	if f.senderFilterHandler == nil {
		return nil
	}
	return f.senderFilterHandler.GetResponseHeaders()
}

func (f *Filter) SetHttpResponseHeader(header api.HeaderMap) {
	if f.senderFilterHandler == nil {
		return
	}
	f.senderFilterHandler.SetResponseHeaders(header)
}

func (f *Filter) GetHttpResponseBody() buffer.IoBuffer {
	if f.senderFilterHandler == nil {
		return nil
	}
	return f.senderFilterHandler.GetResponseData()
}

func (f *Filter) SetHttpResponseBody(body buffer.IoBuffer) {
	if f.senderFilterHandler == nil {
		return
	}
	f.senderFilterHandler.SetResponseData(body)
}

func (f *Filter) GetHttpResponseTrailer() api.HeaderMap {
	if f.senderFilterHandler == nil {
		return nil
	}
	return f.senderFilterHandler.GetResponseTrailers()
}

func (f *Filter) SetHttpResponseTrailer(trailer api.HeaderMap) {
	if f.senderFilterHandler == nil {
		return
	}
	f.senderFilterHandler.SetResponseTrailers(trailer)
}
