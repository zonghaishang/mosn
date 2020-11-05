package proxywasm

import (
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)

var root_id = 100
var id int32

var RootContext *rootContext

func initWasmVM(config *StreamProxyWasmConfig) {
	RootContext = &rootContext{
		config:        config,
	}

	RootContext.wasmCode, _ = wasm.ReadBytes(config.Path)
	RootContext.wasmModule, _ = wasm.Compile(RootContext.wasmCode)
	RootContext.wasiVersion = wasm.WasiGetVersion(RootContext.wasmModule)

	RootContext.wasmImportObj = wasm.NewDefaultWasiImportObjectForVersion(RootContext.wasiVersion)

	im := ProxyWasmImports()

	RootContext.wasmImportObj.Extend(*im)
}

func NewWasmInstance() *wasm.Instance {
	instance, err := RootContext.wasmModule.InstantiateWithImportObject(RootContext.wasmImportObj)
	if err != nil {
		log.DefaultLogger.Errorf("wasm instance error :%v", err)
		return nil
	}

	if _, err := instance.Exports["_start"](); err != nil {
		log.DefaultLogger.Errorf("wasm start err: %v\n", err)
		return nil
	}

	id++
	instanceCtx := &wasmContext{
		contextId : id,
		filter: nil,
		instance: &instance,
	}

	instance.SetContextData(instanceCtx)

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
