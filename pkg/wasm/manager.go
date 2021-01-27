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
	"errors"
	"sync"

	"mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

var (
	ErrWasmConfig     = errors.New("wasm config error")
	ErrUnexpectedType = errors.New("unexpected object type in map")
	ErrPluginNotFound = errors.New("wasm plugin not found")
)

var wasmManagerInstance types.WasmManager = &wasmMangerImpl{}

func GetWasmManager() types.WasmManager {
	return wasmManagerInstance
}

type pluginWrapper struct {
	mu         sync.RWMutex
	plugin     types.WasmPlugin
	wasmConfig v2.WasmPluginConfig
}

func (w *pluginWrapper) GetPlugin() types.WasmPlugin {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.plugin
}

func (w *pluginWrapper) GetConfig() v2.WasmPluginConfig {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.wasmConfig
}

type wasmMangerImpl struct {
	pluginMap sync.Map
}

func (w *wasmMangerImpl) AddOrUpdateWasm(wasmConfig v2.WasmPluginConfig) error {
	if wasmConfig.PluginName == "" {
		log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm err: %v", ErrWasmConfig)
		return ErrWasmConfig
	}
	if v, ok := w.pluginMap.Load(wasmConfig.PluginName); ok {
		pw, ok := v.(*pluginWrapper)
		if !ok {
			log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm err: %v", ErrUnexpectedType)
			return ErrUnexpectedType
		}
		wasm, err := NewWasmPlugin(wasmConfig)
		if err != nil {
			log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm fail to create wasm plugin: %v, err: %v", wasmConfig.PluginName, err)
			return err
		}
		pw.mu.Lock()
		pw.plugin = wasm
		pw.wasmConfig = wasmConfig
		pw.mu.Unlock()
		log.DefaultLogger.Infof("[wasm][manager] AddOrUpdateWasm update wasm plugin: %v", wasmConfig.PluginName)
	} else {
		// add new wasm plugin
		plugin, err := NewWasmPlugin(wasmConfig)
		if err != nil {
			log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm fail to create wasm plugin: %v, err: %v", wasmConfig.PluginName, err)
			return err
		}
		pw := &pluginWrapper{
			plugin:     plugin,
			wasmConfig: wasmConfig,
		}
		w.pluginMap.LoadOrStore(wasmConfig.PluginName, pw)
		log.DefaultLogger.Infof("[wasm][manager] AddOrUpdateWasm add new wasm plugin: %v", wasmConfig.PluginName)
	}

	return nil
}

func (w *wasmMangerImpl) GetWasmPluginWrapperByName(wasmPluginName string) types.WasmPluginWrapper {
	if v, ok := w.pluginMap.Load(wasmPluginName); ok {
		pw, ok := v.(*pluginWrapper)
		if !ok {
			log.DefaultLogger.Errorf("[wasm][manager] GetWasmPluginWrapperByName unexpected object type in map")
			return nil
		}
		return pw
	}

	log.DefaultLogger.Errorf("[wasm][manager] GetWasmPluginWrapperByName not found in map, plugin name: %v", wasmPluginName)
	return nil
}

func (w *wasmMangerImpl) UninstallWasmPluginByName(wasmPluginName string) error {
	if _, ok := w.pluginMap.Load(wasmPluginName); !ok {
		log.DefaultLogger.Errorf("[wasm][manager] UninstallWasmPluginByName plugin not found, name: %v", wasmPluginName)
		return ErrPluginNotFound
	}
	w.pluginMap.Delete(wasmPluginName)
	log.DefaultLogger.Infof("[wasm][manager] UninstallWasmPluginByName uninstall wasm plugin: %v", wasmPluginName)

	return nil
}