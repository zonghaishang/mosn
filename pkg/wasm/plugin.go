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
	"io/ioutil"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"runtime"
	"sync/atomic"
)

var (
	ErrEngineNotFound = errors.New("fail to get wasm engine")
	ErrWasmBytesLoad  = errors.New("fail to load wasm bytes")
	ErrInstanceCreate = errors.New("fail to create wasm instance")
	ErrModuleCreate   = errors.New("fail to create wasm module")
)

type wasmPluginImpl struct {
	config v2.WasmPluginConfig

	instanceWrappers    []types.WasmInstanceWrapper
	instanceWrappersIdx int32
	occupy              int32

	instanceNum int

	vm     types.WasmVM
	module types.WasmModule
}

func NewWasmPlugin(wasmConfig v2.WasmPluginConfig) (types.WasmPlugin, error) {
	instanceNum := wasmConfig.InstanceNum
	if instanceNum <= 0 {
		instanceNum = runtime.NumCPU()
	}

	plugin := &wasmPluginImpl{
		config:           wasmConfig,
		instanceNum:      instanceNum,
		instanceWrappers: make([]types.WasmInstanceWrapper, 0),
	}

	err := plugin.createInstances()
	if err != nil {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin fail to create instance, err: %v", err)
		return nil, err
	}

	return plugin, nil
}

func (w *wasmPluginImpl) PluginName() string {
	return w.config.PluginName
}

func (w *wasmPluginImpl) Clear() {
	// do nothing
	return
}

func (w *wasmPluginImpl) Exec(f func(instanceWrapper types.WasmInstanceWrapper) bool) {
	for _, iw := range w.instanceWrappers {
		if !f(iw) {
			break
		}
	}
}

func (w *wasmPluginImpl) GetPluginConfig() v2.WasmPluginConfig {
	return w.config
}

func (w *wasmPluginImpl) GetVmConfig() v2.WasmVmConfig {
	return *w.config.VmConfig
}

func (w *wasmPluginImpl) GetInstance() types.WasmInstanceWrapper {
	idx := int(w.instanceWrappersIdx) % len(w.instanceWrappers)

	iw := w.instanceWrappers[idx]

	w.instanceWrappersIdx++
	atomic.AddInt32(&w.occupy, 1)

	return iw
}

func (w *wasmPluginImpl) ReleaseInstance(instanceWrapper types.WasmInstanceWrapper) {
	atomic.AddInt32(&w.occupy, -1)
}

func (w *wasmPluginImpl) loadWasmBytesFromPath(path string) []byte {
	if path == "" {
		return nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.DefaultLogger.Errorf("[wasm][plugin] loadWasmBytesFromPath read file err: %v", err)
		return nil
	}

	return bytes
}

func (w *wasmPluginImpl) loadWasmBytesFromUrl(url string) []byte {
	if url == "" {
		return nil
	}

	return nil
}

func (w *wasmPluginImpl) createInstances() error {
	// get wasm engine
	w.vm = GetWasmEngine(w.config.VmConfig.Engine)
	if w.vm == nil {
		log.DefaultLogger.Errorf("[wasm][plugin] init fail to get wasm engine: %v", w.config.VmConfig.Engine)
		return ErrEngineNotFound
	}

	// load wasm bytes
	var wasmBytes []byte
	if w.config.VmConfig.Path != "" {
		wasmBytes = w.loadWasmBytesFromPath(w.config.VmConfig.Path)
	} else {
		wasmBytes = w.loadWasmBytesFromUrl(w.config.VmConfig.Url)
	}

	if wasmBytes == nil || len(wasmBytes) == 0 {
		log.DefaultLogger.Errorf("[wasm][plugin] createInstances fail to load wasm bytes, config: %v", w.config)
		return ErrWasmBytesLoad
	}

	// create wasm module
	w.module = w.vm.NewModule(wasmBytes)
	if w.module == nil {
		log.DefaultLogger.Errorf("[wasm][plugin] createInstances fail to create module, config: %v", w.config)
		return ErrModuleCreate
	}

	// create wasm instance
	instanceCreated := 0
	for i := 0; i < w.instanceNum; i++ {
		instance := w.module.NewInstance()
		if instance == nil {
			log.DefaultLogger.Errorf("[wasm][plugin] createInstances fail to create instance, i: %v", i)
			continue
		}

		w.instanceWrappers = append(w.instanceWrappers, &wasmInstanceWrapperImpl{WasmInstance: instance})
		instanceCreated++
	}

	if instanceCreated == 0 {
		log.DefaultLogger.Errorf("[wasm][plugin] createInstances create 0 instance")
		return ErrInstanceCreate
	}

	return nil
}
