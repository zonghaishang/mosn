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
	"encoding/json"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/buffer"
)

type proxyWasmInstanceCallback struct {
	filter *Filter
}

func (p *proxyWasmInstanceCallback) GetVmConfig() buffer.IoBuffer {
	vmConfig := p.filter.plugin.GetVmConfig()
	vmConfigBytes, err := json.Marshal(vmConfig)
	if err != nil {
		return nil
	}
	return buffer.NewIoBufferBytes(vmConfigBytes)
}

func (p *proxyWasmInstanceCallback) GetPluginConfig() buffer.IoBuffer {
	pluginConfig := p.filter.plugin.GetPluginConfig()
	pluginConfigBytes, err := json.Marshal(pluginConfig)
	if err != nil {
		return nil
	}
	return buffer.NewIoBufferBytes(pluginConfigBytes)
}

func (p *proxyWasmInstanceCallback) Log(level log.Level, msg string) {
	logFunc := log.DefaultLogger.Infof
	switch level {
	case log.TRACE:
		logFunc = log.DefaultLogger.Tracef
	case log.DEBUG:
		logFunc = log.DefaultLogger.Debugf
	case log.INFO:
		logFunc = log.DefaultLogger.Infof
	case log.WARN:
		logFunc = log.DefaultLogger.Warnf
	case log.ERROR:
		logFunc = log.DefaultLogger.Errorf
	case log.FATAL:
		logFunc = log.DefaultLogger.Fatalf
	}
	logFunc(msg)
}

func (p *proxyWasmInstanceCallback) GetHttpRequestHeader() api.HeaderMap {
	return p.filter.receiverFilterHandler.GetRequestHeaders()
}

func (p *proxyWasmInstanceCallback) GetHttpRequestBody() buffer.IoBuffer {
	return p.filter.receiverFilterHandler.GetRequestData()
}

func (p *proxyWasmInstanceCallback) GetHttpRequestTrailer() api.HeaderMap {
	return p.filter.receiverFilterHandler.GetRequestTrailers()
}

func (p *proxyWasmInstanceCallback) GetHttpResponseHeader() api.HeaderMap {
	return p.filter.senderFilterHandler.GetResponseHeaders()
}

func (p *proxyWasmInstanceCallback) GetHttpResponseBody() buffer.IoBuffer {
	return p.filter.senderFilterHandler.GetResponseData()
}

func (p *proxyWasmInstanceCallback) GetHttpResponseTrailer() api.HeaderMap {
	return p.filter.senderFilterHandler.GetResponseTrailers()
}
