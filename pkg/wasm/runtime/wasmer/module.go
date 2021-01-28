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

package wasmer

import (
	"strings"

	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
)

type Module struct {
	vm           *VM
	module       *wasmerGo.Module
	importObject *wasmerGo.ImportObject
}

func (w *Module) Init() {
	return
}

func (w *Module) GetABIList() []abi.ABI {
	abiList := make([]abi.ABI, 0)

	exportList := w.module.Exports()

	for _, export := range exportList {
		if export.Type().Kind() == wasmerGo.FUNCTION {
			if strings.HasPrefix(export.Name(), "proxy_abi") {
				v := abi.GetABI(export.Name())

				if v != nil {
					abiList = append(abiList, v)
				}
			}
		}
	}

	return abiList
}

func (w *Module) NewInstance() types.WasmInstance {
	abiList := w.GetABIList()

	res := newWasmerInstance(w.vm, w.module)

	for _, v := range abiList {
		v.OnInstanceCreate(res)
	}

	instance, err := wasmerGo.NewInstance(w.module, res.importObject)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmer][module] new Wasmer Instance, failed to instantiate the module, err: %v", err)
		return nil
	}

	res.instance = instance

	err = res.Start()
	if err != nil {
		log.DefaultLogger.Errorf("[wasmer][module] new Wasmer Instance, fail to call Start(), err: %v", err)
		return nil
	}

	for _, v := range abiList {
		v.OnStart(res)
	}

	return res
}
