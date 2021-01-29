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
	"sync"
	"sync/atomic"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm_0_1_0"
	"mosn.io/pkg/buffer"
)

type Filter struct {
	ctx context.Context

	pluginName string
	plugin     types.WasmPlugin
	instance   types.WasmInstanceWrapper

	abi     abi.ABI
	exports proxywasm_0_1_0.Exports

	rootContextID int32
	contextID     int32

	receiverFilterHandler api.StreamReceiverFilterHandler
	senderFilterHandler   api.StreamSenderFilterHandler

	reqHeader api.HeaderMap
	reqBody buffer.IoBuffer
	reqTrailer api.HeaderMap

	respHeader api.HeaderMap
	respBody buffer.IoBuffer
	respTrailer api.HeaderMap

	destroyOnce sync.Once
}

var contextIDGenerator int32

func NewFilter(ctx context.Context, pluginName string, rootContextID int32) *Filter {
	abiVersion := abi.GetABI("proxy_abi_version_0_1_0")
	if abiVersion == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] NewFilter abi version not found")
		return nil
	}

	// get wasm plugin instance
	pluginWrapper := wasm.GetWasmManager().GetWasmPluginWrapperByName(pluginName)
	if pluginWrapper == nil {
		log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive wasm plugin not exists, plugin name: %v", pluginName)
		return nil
	}

	plugin := pluginWrapper.GetPlugin()
	instance := plugin.GetInstance()

	filter := &Filter{
		ctx:           ctx,
		pluginName:    pluginName,
		plugin:        plugin,
		instance:      instance,
		abi:           abiVersion,
		exports:       abiVersion.(proxywasm_0_1_0.Exports),
		rootContextID: rootContextID,
		contextID:     atomic.AddInt32(&contextIDGenerator, 1),
	}

	filter.abi.SetInstance(filter.instance)
	filter.abi.SetInstanceCallBack(&proxyWasmInstanceCallback{filter: filter})

	filter.instance.Acquire()
	_ = filter.exports.ProxyOnContextCreate(filter.contextID, filter.rootContextID)
	filter.instance.Release()

	return filter
}

func (f *Filter) OnDestroy() {
	f.destroyOnce.Do(func() {
		f.instance.Acquire()
		_, _ = f.exports.ProxyOnDone(f.contextID)
		f.instance.Release()

		f.plugin.ReleaseInstance(f.instance)
	})
}

func (f *Filter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.receiverFilterHandler = handler
}

func (f *Filter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.senderFilterHandler = handler
}

func (f *Filter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.reqHeader = headers
	f.reqBody = buf
	f.reqTrailer = trailers

	f.instance.Acquire()
	defer f.instance.Release()

	// do filter
	if buf != nil && buf.Len() > 0 {
		if _, err := f.exports.ProxyOnRequestHeaders(f.contextID, 0, 0); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestHeaders err: %v", err)
		}
		if _, err := f.exports.ProxyOnRequestBody(f.contextID, int32(buf.Len()), 1); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestBody err: %v", err)
		}
	} else {
		if _, err := f.exports.ProxyOnRequestHeaders(f.contextID, 0, 1); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestHeaders err: %v", err)
		}
	}

	return api.StreamFilterContinue
}

func (f *Filter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.respHeader = headers
	f.respBody = buf
	f.respTrailer = trailers

	f.instance.Acquire()
	defer f.instance.Release()

	// do filter
	if buf != nil && buf.Len() > 0 {
		if _, err := f.exports.ProxyOnResponseHeaders(f.contextID, 0, 0); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestHeaders err: %v", err)
		}
		if _, err := f.exports.ProxyOnResponseBody(f.contextID, int32(buf.Len()), 1); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestBody err: %v", err)
		}
	} else {
		if _, err := f.exports.ProxyOnResponseHeaders(f.contextID, 0, 1); err != nil {
			log.DefaultLogger.Errorf("[x-proxy-wasm][filter] OnReceive ProxyOnRequestHeaders err: %v", err)
		}
	}

	return api.StreamFilterContinue
}
