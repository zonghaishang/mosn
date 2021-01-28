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

type WasmManager interface {
	AddOrUpdateWasm(wasmConfig v2.WasmPluginConfig) error
	GetWasmPluginWrapperByName(pluginName string) WasmPluginWrapper
	UninstallWasmPluginByName(pluginName string) error
}

type WasmPluginWrapper interface {
	GetPlugin() WasmPlugin
	GetConfig() v2.WasmPluginConfig

	Update(plugin WasmPlugin, config v2.WasmPluginConfig)
}

type WasmPlugin interface {
	PluginName() string

	GetPluginConfig() v2.WasmPluginConfig
	GetVmConfig() v2.WasmVmConfig

	GetInstance() WasmInstanceWrapper
	ReleaseInstance(instanceWrapper WasmInstanceWrapper)

	Exec(func(instanceWrapper WasmInstanceWrapper)bool)

	Clear()
}

type WasmInstanceWrapper interface {
	Acquire()
	Release()

	WasmInstance
}

type WasmVM interface {
	Init()
	NewModule(wasmBytes []byte) WasmModule
}

type WasmModule interface {
	Init()
	NewInstance() WasmInstance
}

type WasmInstance interface {
	GetExportsFunc(funcName string) (WasmFunction, error)
	GetExportsMem(memName string) ([]byte, error)

	GetMemory(addr uint64, size uint64) ([]byte, error)
	PutMemory(addr uint64, size uint64, content []byte) error

	GetByte(addr uint64) (byte, error)
	PutByte(addr uint64, b byte) error

	GetUint32(addr uint64) (uint32, error)
	PutUint32(addr uint64, value uint32) error

	Malloc(size int32) (uint64, error)

	RegisterFunc(namespace string, funcName string, f interface{})
}

type WasmFunction interface {
	Call(args ...interface{}) (interface{}, error)
}