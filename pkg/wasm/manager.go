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
	ErrEmptyPluginName = errors.New("wasm config without plugin name")
	ErrUnexpectedType  = errors.New("unexpected object type in map")
	ErrPluginNotFound  = errors.New("wasm plugin not found")
)

var wasmManagerInstance types.WasmManager = &wasmMangerImpl{}

func GetWasmManager() types.WasmManager {
	return wasmManagerInstance
}

type pluginWrapper struct {
	mu     sync.RWMutex
	plugin types.WasmPlugin
	config v2.WasmPluginConfig
}

func (w *pluginWrapper) GetPlugin() types.WasmPlugin {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.plugin
}

func (w *pluginWrapper) GetConfig() v2.WasmPluginConfig {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.config
}

func (w *pluginWrapper) Update(plugin types.WasmPlugin, config v2.WasmPluginConfig) {
	if plugin == nil || config.PluginName == "" || config.PluginName != w.config.PluginName {
		return
	}

	w.mu.Lock()
	oldPlugin := w.plugin

	w.plugin = plugin
	w.config = config

	w.mu.Unlock()

	oldPlugin.Clear()
}

type wasmMangerImpl struct {
	pluginMap sync.Map
}

func (w *wasmMangerImpl) AddOrUpdateWasm(config v2.WasmPluginConfig) error {
	if config.PluginName == "" {
		log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm empty plugin name")
		return ErrEmptyPluginName
	}

	if v, ok := w.pluginMap.Load(config.PluginName); ok {
		pluginWrapper, ok := v.(*pluginWrapper)
		if !ok {
			log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm unexpected type in map")
			return ErrUnexpectedType
		}

		plugin, err := NewWasmPlugin(config)
		if err != nil {
			log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm fail to create wasm plugin: %v, err: %v", config.PluginName, err)
			return err
		}

		pluginWrapper.Update(plugin, config)

		log.DefaultLogger.Infof("[wasm][manager] AddOrUpdateWasm update wasm plugin: %v", config.PluginName)
	} else {
		// add new wasm plugin
		plugin, err := NewWasmPlugin(config)
		if err != nil {
			log.DefaultLogger.Errorf("[wasm][manager] AddOrUpdateWasm fail to create wasm plugin: %v, err: %v", config.PluginName, err)
			return err
		}

		pw := &pluginWrapper{
			plugin: plugin,
			config: config,
		}

		w.pluginMap.LoadOrStore(config.PluginName, pw)

		log.DefaultLogger.Infof("[wasm][manager] AddOrUpdateWasm add new wasm plugin: %v", config.PluginName)
	}

	return nil
}

func (w *wasmMangerImpl) GetWasmPluginWrapperByName(pluginName string) types.WasmPluginWrapper {
	if pluginName == "" {
		return nil
	}

	if v, ok := w.pluginMap.Load(pluginName); ok {
		pw, ok := v.(*pluginWrapper)
		if !ok {
			log.DefaultLogger.Errorf("[wasm][manager] GetWasmPluginWrapperByName unexpected object type in map")
			return nil
		}
		return pw
	} else {
		log.DefaultLogger.Errorf("[wasm][manager] GetWasmPluginWrapperByName not found in map, plugin name: %v", pluginName)
	}

	return nil
}

func (w *wasmMangerImpl) UninstallWasmPluginByName(pluginName string) error {
	v, ok := w.pluginMap.Load(pluginName)
	if !ok {
		log.DefaultLogger.Errorf("[wasm][manager] UninstallWasmPluginByName plugin not found, name: %v", pluginName)
		return ErrPluginNotFound
	}

	w.pluginMap.Delete(pluginName)

	pw := v.(*pluginWrapper)
	pw.GetPlugin().Clear()

	log.DefaultLogger.Infof("[wasm][manager] UninstallWasmPluginByName uninstall wasm plugin: %v", pluginName)

	return nil
}
