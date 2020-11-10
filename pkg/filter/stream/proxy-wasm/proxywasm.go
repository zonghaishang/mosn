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
		config:      config,
		propertyMap: make(map[string]string),
		metrics:     newWasmMetricsManager(config.Name),
	}
	var err error

	RootContext.wasmCode, err = wasm.ReadBytes(config.Path)
	if err != nil {
		log.DefaultLogger.Errorf("wasm fail to read code bytes, err: %v", err)
	}

	RootContext.wasmModule, err = wasm.Compile(RootContext.wasmCode)
	if err != nil {
		log.DefaultLogger.Errorf("wasm fail to compile wasm code, err: %v", err)
	}
	RootContext.wasiVersion = wasm.WasiGetVersion(RootContext.wasmModule)

	RootContext.wasmImportObj = wasm.NewDefaultWasiImportObjectForVersion(RootContext.wasiVersion)

	im := ProxyWasmImports()

	err = RootContext.wasmImportObj.Extend(*im)
	if err != nil {
		log.DefaultLogger.Errorf("wasm fail to extend import obj, err: %v", err)
	}
}

func NewWasmInstance() *wasmContext {
	instance, err := RootContext.wasmModule.InstantiateWithImportObject(RootContext.wasmImportObj)
	if err != nil {
		log.DefaultLogger.Errorf("wasm instance error :%v", err)
		return nil
	}

	id++
	instanceCtx := &wasmContext{
		rootContext: RootContext,
		contextId:   id,
		instance:    &instance,
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
