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

package abi

import (
	"sync"

	"mosn.io/mosn/pkg/types"
)

// Factory is the ABI factory func
type Factory func() ABI

// string -> Factory
var abiMap = sync.Map{}

// ABI represents the abi between the host and wasm, which consists of two parts: exports and life-cycle handler
// *exports* represents the exports elements of the wasm module
// *life-cycle handler* is used to manage the life-cycle of an abi
type ABI interface {
	// life-cycle: OnInstanceCreate got called when instantiating the wasm instance
	OnInstanceCreate(instance types.WasmInstance)

	// life-cycle: OnStart got called when starting the wasm instance
	OnStart(instance types.WasmInstance)

	// life-cycle: OnInstanceDestroy got called when destroying the wasm instance
	OnInstanceDestroy(instance types.WasmInstance)

	// life-cycle: SetInstance sets the wasm instance to the abi
	SetInstance(instance types.WasmInstance)
}

// RegisterABI registers an abi factory
func RegisterABI(abiName string, factory Factory) {
	abiMap.Store(abiName, factory)
}

// GetABI returns an abi by name
func GetABI(abiName string) ABI {
	if v, ok := abiMap.Load(abiName); ok {
		factory := v.(Factory)
		return factory()
	}
	return nil
}
