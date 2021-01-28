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
	"encoding/binary"
	"errors"
	"reflect"

	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

var (
	ErrAddrOverflow = errors.New("addr overflow")
)

type Instance struct {
	vm           *VM
	module       *wasmerGo.Module
	importObject *wasmerGo.ImportObject
	instance     *wasmerGo.Instance
	memory       *wasmerGo.Memory
	funcCache    map[string]*wasmerGo.Function
}

func newWasmerInstance(vm *VM, module *wasmerGo.Module) *Instance {
	instance := &Instance{
		vm:           vm,
		module:       module,
		importObject: wasmerGo.NewImportObject(),
		funcCache:    make(map[string]*wasmerGo.Function),
	}

	instance.ensureWASIimports()

	return instance
}

func (w *Instance) RegisterFunc(namespace string, funcName string, f interface{}) {
	ftype := reflect.TypeOf(f)

	argsNum := ftype.NumIn()
	argsKind := make([]*wasmerGo.ValueType, argsNum)
	for i := 0; i < argsNum; i++ {
		argsKind[i] = convertFromGoType(ftype.In(i))
	}

	retsNum := ftype.NumOut()
	retsKind := make([]*wasmerGo.ValueType, retsNum)
	for i := 0; i < retsNum; i++ {
		retsKind[i] = convertFromGoType(ftype.Out(i))
	}

	fwasmer := wasmerGo.NewFunction(
		w.vm.store,
		wasmerGo.NewFunctionType(argsKind, retsKind),
		func(args []wasmerGo.Value) ([]wasmerGo.Value, error) {
			aa := convertToGoTypes(args)

			callResult := reflect.ValueOf(f).Call(aa)

			ret := convertFromGoValue(callResult[0])

			return []wasmerGo.Value{ret}, nil
		},
	)

	w.importObject.Register(namespace, map[string]wasmerGo.IntoExtern{
		funcName: fwasmer,
	})
}

func (w *Instance) ensureWASIimports() {
	importList := w.module.Imports()

	for _, im := range importList {
		if im.Type().Kind() == wasmerGo.FUNCTION && im.Module() == "wasi_unstable" {

			fType := im.Type().IntoFunctionType()

			f := wasmerGo.NewFunction(w.vm.store, wasmerGo.NewFunctionType(fType.Params(), fType.Results()),
				func(values []wasmerGo.Value) ([]wasmerGo.Value, error) {
					return nil, nil
				},
			)

			w.importObject.Register("wasi_unstable", map[string]wasmerGo.IntoExtern{im.Name(): f})
		}
	}
}

func (w *Instance) Start() error {
	f, err := w.instance.Exports.GetFunction("_start")
	if err != nil {
		log.DefaultLogger.Errorf("[wasmer][instance] WasmerInstance fail to get export func: _start, err: %v", err)
		return err
	}

	_, err = f()
	if err != nil {
		log.DefaultLogger.Errorf("[wasmer][instance] WasmerInstance fail to call _start func, err: %v", err)
		return err
	}

	return nil
}

func (w *Instance) Malloc(size int32) (uint64, error) {
	malloc, err := w.GetExportsFunc("malloc")
	if err != nil {
		return 0, err
	}
	addr, err := malloc.Call(size)
	if err != nil {
		return 0, err
	}
	return uint64(addr.(int32)), nil
}

func (w *Instance) GetExportsFunc(funcName string) (types.WasmFunction, error) {
	if f, ok := w.funcCache[funcName]; ok {
		return f, nil
	}

	f, err := w.instance.Exports.GetRawFunction(funcName)
	if err != nil {
		return nil, err
	}

	w.funcCache[funcName] = f

	return f, nil
}

func (w *Instance) GetExportsMem(memName string) ([]byte, error) {
	if w.memory == nil {
		m, err := w.instance.Exports.GetMemory(memName)
		if err != nil {
			return nil, err
		}
		w.memory = m
	}
	return w.memory.Data(), nil
}

func (w *Instance) GetMemory(addr uint64, size uint64) ([]byte, error) {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return nil, err
	}
	if int(addr) > len(mem) || int(addr+size) > len(mem) {
		return nil, ErrAddrOverflow
	}
	return mem[addr : addr+size], nil
}

func (w *Instance) PutMemory(addr uint64, size uint64, content []byte) error {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return err
	}
	if int(addr) > len(mem) || int(addr+size) > len(mem) {
		return ErrAddrOverflow
	}
	copySize := uint64(len(content))
	if size < copySize {
		copySize = size
	}
	copy(mem[addr:], content[:copySize])
	return nil
}

func (w *Instance) GetByte(addr uint64) (byte, error) {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}
	if int(addr) > len(mem) {
		return 0, ErrAddrOverflow
	}
	return mem[addr], nil
}

func (w *Instance) PutByte(addr uint64, b byte) error {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return err
	}
	if int(addr) > len(mem) {
		return ErrAddrOverflow
	}
	mem[addr] = b
	return nil
}

func (w *Instance) GetUint32(addr uint64) (uint32, error) {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}
	if int(addr) > len(mem) || int(addr+4) > len(mem) {
		return 0, ErrAddrOverflow
	}
	return binary.LittleEndian.Uint32(mem[addr:]), nil
}

func (w *Instance) PutUint32(addr uint64, value uint32) error {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return err
	}
	if int(addr) > len(mem) || int(addr+4) > len(mem) {
		return ErrAddrOverflow
	}
	binary.LittleEndian.PutUint32(mem[addr:], value)
	return nil
}
