package proxywasm

import (
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)

var root_id = 100

type wasmContext struct {
	filter   *streamProxyWasmFilter
	instance *wasm.Instance
}

var wasmCode []byte
var wasmModule wasm.Module
var wasiVersion wasm.WasiVersion
var wasmImportObj *wasm.ImportObject

func initWasm(path string) {
	wasmCode, _ = wasm.ReadBytes(path)
	wasmModule, _ = wasm.Compile(wasmCode)
	wasiVersion = wasm.WasiGetVersion(wasmModule)

	wasmImportObj = wasm.NewDefaultWasiImportObjectForVersion(wasiVersion)

	im := ProxyWasmImports()

	wasmImportObj.Extend(*im)
}

func NewWasmInstance() *wasm.Instance {
	instance, err := wasmModule.InstantiateWithImportObject(wasmImportObj)
	if err != nil {
		log.DefaultLogger.Errorf("wasm instance error :%v", err)
		return nil
	}

	if _, err := instance.Exports["_start"](); err != nil {
		log.DefaultLogger.Errorf("wasm start err: %v\n", err)
		return nil
	}

	instance.SetContextData(&wasmContext{
		nil, &instance,
	})

	if _, err := instance.Exports["proxy_on_context_create"](root_id, 0); err != nil {
		log.DefaultLogger.Errorf("root err %v\n", err)
		return nil
	}

	if _, err := instance.Exports["proxy_on_vm_start"](root_id, 1000); err != nil {
		log.DefaultLogger.Errorf("start err %v\n", err)
		return nil
	}

	return &instance
}
