package proxywasm

//var root_id = 100
//var id int32

//var RootContext *rootContext
//var rootWasmInstance *wasmContext

//func initWasmVM(config *StreamProxyWasmConfig) {
//	RootContext = &rootContext{
//		config:      config,
//		propertyMap: make(map[string]string),
//		metrics:     newWasmMetricsManager(config.name),
//	}
//	var err error
//
//	RootContext.wasmCode, err = wasm.ReadBytes(config.path)
//	if err != nil {
//		log.DefaultLogger.Errorf("wasm fail to read code bytes, err: %v", err)
//		return
//	}
//
//	RootContext.wasmModule, err = wasm.Compile(RootContext.wasmCode)
//	if err != nil {
//		log.DefaultLogger.Errorf("wasm fail to compile wasm code, err: %v", err)
//		return
//	}
//	RootContext.wasiVersion = wasm.WasiGetVersion(RootContext.wasmModule)
//
//	RootContext.wasmImportObj = wasm.NewDefaultWasiImportObjectForVersion(RootContext.wasiVersion)
//
//	im := ProxyWasmImports()
//
//	err = RootContext.wasmImportObj.Extend(*im)
//	if err != nil {
//		log.DefaultLogger.Errorf("wasm fail to extend import obj, err: %v", err)
//	}
//}
//
//func NewWasmInstance() *wasmContext {
//	id++
//	instanceCtx := &wasmContext{
//		rootContext: RootContext,
//		contextId:   id,
//	}
//
//	instance, err := RootContext.wasmModule.InstantiateWithImportObject(RootContext.wasmImportObj)
//	if err != nil {
//		log.DefaultLogger.Errorf("wasm instance error :%v", err)
//		instanceCtx.ProxyWasmInstance = &NilWasmInstance{}
//		return instanceCtx
//	} else {
//		instanceCtx.ProxyWasmInstance = &WasmerInstance{instance}
//	}
//
//	// _start must be in the front of SetContextData, don't ask me why
//	if err := instanceCtx._start(); err != nil {
//		log.DefaultLogger.Errorf("wasm start err: %v\n", err)
//		instanceCtx.ProxyWasmInstance = &NilWasmInstance{}
//		return instanceCtx
//	}
//
//	instance.SetContextData(instanceCtx)
//
//	if err := instanceCtx.proxy_on_context_create(int32(root_id), 0); err != nil {
//		log.DefaultLogger.Errorf("root err %v\n", err)
//		instanceCtx.ProxyWasmInstance = &NilWasmInstance{}
//		return instanceCtx
//	}
//
//	return instanceCtx
//}
