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
	"mosn.io/mosn/pkg/types"
)

var abiMap = make(map[string]ABI)

type ABI interface {
	OnInstanceCreate(instance types.WasmInstance)
	OnStart(instance types.WasmInstance)
	OnInstanceDestroy(instance types.WasmInstance)

	SetInstance(instance types.WasmInstance)
	SetInstanceCallBack(callback interface{})
}

func RegisterABI(abiName string, abiImpl ABI) {
	abiMap[abiName] = abiImpl
}

func GetABI(abiName string) ABI {
	return abiMap[abiName]
}
