package proxywasm

import (
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)

var root_id = 100
var id int32

var RootContext *rootContext
var rootWasmInstance *wasmContext

func initWasmVM(config *StreamProxyWasmConfig) {
	RootContext = &rootContext{
		config: config,
	}

	RootContext.wasmCode, _ = wasm.ReadBytes(config.Path)
	RootContext.wasmModule, _ = wasm.Compile(RootContext.wasmCode)
	RootContext.wasiVersion = wasm.WasiGetVersion(RootContext.wasmModule)

	RootContext.wasmImportObj = wasm.NewDefaultWasiImportObjectForVersion(RootContext.wasiVersion)

	im := ProxyWasmImports()

	RootContext.wasmImportObj.Extend(*im)
}

func NewWasmInstance() *wasmContext {
	instance, err := RootContext.wasmModule.InstantiateWithImportObject(RootContext.wasmImportObj)
	if err != nil {
		log.DefaultLogger.Errorf("wasm instance error :%v", err)
		return nil
	}

	id++
	instanceCtx := &wasmContext{
		contextId: id,
		instance:  &instance,
	}

	// _start must be in the front of SetContextData, don't ask me why
	if err := instanceCtx._start(); err != nil {
		log.DefaultLogger.Errorf("wasm start err: %v\n", err)
		return nil
	}

	instance.SetContextData(instanceCtx)

	if err := instanceCtx.proxy_on_context_create(int32(root_id), 0); err != nil {
		log.DefaultLogger.Errorf("root err %v\n", err)
		return nil
	}

	return instanceCtx
}
