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

package types

import (
	v2 "mosn.io/mosn/pkg/config/v2"
)

// WasmManager managers all wasm configs
type WasmManager interface {
	// AddOrUpdateWasm add or update wasm plugin config
	AddOrUpdateWasm(wasmConfig v2.WasmPluginConfig) error

	// GetWasmPluginWrapperByName returns wasm plugin by name
	GetWasmPluginWrapperByName(pluginName string) WasmPluginWrapper

	// UninstallWasmPluginByName remove wasm plugin by name
	UninstallWasmPluginByName(pluginName string) error
}

// WasmPluginWrapper wraps wasm plugin with its config and plugin handler
type WasmPluginWrapper interface {
	// GetPlugin returns the wasm plugin
	GetPlugin() WasmPlugin

	// GetConfig returns the config of wasm plugin
	GetConfig() v2.WasmPluginConfig

	// RegisterPluginHandler registers a plugin handler for the wasm plugin
	RegisterPluginHandler(pluginHandler WasmPluginHandler)

	// Update updates the plugin
	Update(config v2.WasmPluginConfig, plugin WasmPlugin)
}

// WasmPluginHandler provides callbacks to manager the life cycle of the wasm plugin
type WasmPluginHandler interface {
	// OnConfigUpdate got called when updating the config of the wasm plugin
	OnConfigUpdate(config v2.WasmPluginConfig)

	// OnPluginStart got called when starting the wasm plugin
	OnPluginStart(plugin WasmPlugin)

	// OnPluginDestroy got called when destroying the wasm plugin
	OnPluginDestroy(plugin WasmPlugin)
}

// WasmPlugin manages the collection of wasm plugin instances
type WasmPlugin interface {
	// PluginName returns the name of wasm plugin
	PluginName() string

	// GetPluginConfig returns the config of wasm plugin
	GetPluginConfig() v2.WasmPluginConfig

	// GetVmConfig returns the vm config of wasm plugin
	GetVmConfig() v2.WasmVmConfig

	// EnsureInstanceNum tries to expand/shrink the num of instance to 'num'
	// and returns the actual instance num
	EnsureInstanceNum(num int) int

	// InstanceNum returns the current number of wasm instance
	InstanceNum() int

	// GetInstance returns one plugin instance of the plugin
	GetInstance() WasmInstanceWrapper

	// ReleaseInstance releases the instance to the plugin
	ReleaseInstance(instanceWrapper WasmInstanceWrapper)

	// Exec execute the f for each instance
	Exec(f func(instanceWrapper WasmInstanceWrapper)bool)

	// SetCpuLimit set cpu limit of the plugin, not supported currently
	SetCpuLimit(cpu int)

	// SetCpuLimit set cpu limit of the plugin, not supported currently
	SetMemLimit(mem int)

	// Clear got called when the plugin is destroyed
	Clear()
}

// WasmInstanceWrapper wraps the wasm instance with an exclusive lock
type WasmInstanceWrapper interface {
	WasmInstance

	// Acquire gets the exclusive ownership of the wasm instance
	// and sets the user-defined data
	Acquire(data interface{})

	// Release releases the exclusive ownership of the wasm instance
	// and sets the users-defined data to nil
	Release()
}

// WasmVM represents the wasm vm(engine)
type WasmVM interface {
	// Init got called when creating a new wasm vm(engine)
	Init()

	// NewModule compiles the 'wasmBytes' into a wasm module
	NewModule(wasmBytes []byte) WasmModule
}

// WasmModule represents the wasm module
type WasmModule interface {
	// Init got called when creating a new wasm module
	Init()

	// NewInstance instantiates and returns a new wasm instance
	NewInstance() WasmInstance
}

// WasmInstance represents the wasm instance
type WasmInstance interface {
	// GetExportsFunc returns the exported func of the wasm instance
	GetExportsFunc(funcName string) (WasmFunction, error)

	// GetExportsMem returns the exported mem of the wasm instance
	GetExportsMem(memName string) ([]byte, error)

	// GetMemory returns wasm mem bytes from specified addr and size
	GetMemory(addr uint64, size uint64) ([]byte, error)

	// PutMemory sets wasm mem bytes to specified addr and size
	PutMemory(addr uint64, size uint64, content []byte) error

	// GetByte returns one wasm byte from specified addr
	GetByte(addr uint64) (byte, error)

	// PutByte sets one wasms bytes to specified addr
	PutByte(addr uint64, b byte) error

	// GetUint32 returns uint32 from specified addr
	GetUint32(addr uint64) (uint32, error)

	// PutUint32 set uint32 to specified addr
	PutUint32(addr uint64, value uint32) error

	// Malloc allocates size of mem from wasm default memory
	Malloc(size int32) (uint64, error)

	// RegisterFunc registers a func to the wasm instance
	RegisterFunc(namespace string, funcName string, f interface{})

	// GetData returns user-defined data
	GetData() interface{}

	// SetData sets user-defined data into the wasm instance
	SetData(data interface{})
}

// WasmFunction is the func exported by wasm module
type WasmFunction interface {
	// Call invokes the wasm func
	Call(args ...interface{}) (interface{}, error)
}