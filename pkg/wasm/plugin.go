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
	"runtime"

	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

var (
	ErrEngineInit     = errors.New("fail to init wasm engine")
	ErrInstanceCreate = errors.New("fail to create wasm instance")
	ErrModuleCreate   = errors.New("fail to create wasm module")
)

type wasmPluginImpl struct {
	config v2.WasmPluginConfig

	instances      chan types.WasmInstance
	maxInstanceNum int

	vm     types.WasmVM
	module types.WasmModule
}

func NewWasmPlugin(wasmConfig v2.WasmPluginConfig) (types.WasmPlugin, error) {
	maxInstanceNum := wasmConfig.InstanceNum
	if maxInstanceNum <= 0 {
		// by default, we create P * 2 wasm instances
		p := runtime.GOMAXPROCS(0)
		maxInstanceNum = p * 2
	}

	plugin := &wasmPluginImpl{
		config:         wasmConfig,
		maxInstanceNum: maxInstanceNum,
		instances:      make(chan types.WasmInstance, maxInstanceNum),
	}

	err := plugin.init()
	if err != nil {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin init failed, err: %v", err)
		return nil, err
	}

	return plugin, nil
}

func (w *wasmPluginImpl) GetConfig() v2.WasmPluginConfig {
	return w.config
}

func (w *wasmPluginImpl) GetVmConfig() v2.WasmVmConfig {
	return w.config.VmConfig
}

func (w *wasmPluginImpl) GetInstance() types.WasmInstance {
	log.DefaultLogger.Debugf("[wasm][plugin] GetInstance got called")
	return <-w.instances
}

func (w *wasmPluginImpl) ReleaseInstance(instance types.WasmInstance) {
	log.DefaultLogger.Debugf("[wasm][plugin] ReleaseInstance got called")
	w.instances <- instance
}

func (w *wasmPluginImpl) init() error {
	w.vm = GetWasmEngine(w.config.VmConfig.Engine)
	if w.vm == nil {
		log.DefaultLogger.Errorf("[wasm][plugin] init fail to create vm, engine: %v", w.config.VmConfig.Engine)
		return ErrEngineInit
	}

	var path string
	if w.config.VmConfig.Url != "" {
		// TODO: load wasm file from URL
		path = ""
	} else {
		path = w.config.VmConfig.Path
	}

	w.module = w.vm.NewModule(path)
	if w.module == nil {
		log.DefaultLogger.Errorf("[wasm][plugin] init fail to create module, engine: %v, path: %v", w.config.VmConfig.Engine, path)
		return ErrModuleCreate
	}

	// create wasm instance
	for i := 0; i < w.maxInstanceNum; i++ {
		instance := w.module.NewInstance()
		if instance == nil {
			log.DefaultLogger.Errorf("[wasm][plugin] init fail to create instance")
			break
		}

		w.instances <- instance
	}

	if len(w.instances) == 0 {
		log.DefaultLogger.Errorf("[wasm][plugin] init fail to create wasm instance, engine: %v, path: %v", w.config.VmConfig.Engine, path)
		return ErrInstanceCreate
	}

	return nil
}
