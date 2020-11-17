package proxywasm

import (
	"sync"

	"github.com/bytecodealliance/wasmtime-go"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)

type WasmVM interface {
	Init()
	NewModule(path string) WasmModule
}

type WasmModule interface {
	NewInstance() *wasmContext
}

var (
	vmMapSingleton sync.Mutex
	vmMap          map[string]WasmVM
)

func init() {
	vmMap = make(map[string]WasmVM)
}

func NewWasmVM(name string) WasmVM {
	vmMapSingleton.Lock()
	defer vmMapSingleton.Unlock()

	if vm, ok := vmMap[name]; ok {
		return vm
	}

	switch name {
	case "wasmer":
		vm := &WasmerVM{}
		vm.Init()
		vmMap[name] = vm
		return vm
	default:
		log.DefaultLogger.Errorf("Unsupported wasm vm: %v", name)
	}
	return &NilVM{}
}

type NilVM struct {
}

func (n *NilVM) Init() {
	log.DefaultLogger.Errorf("NilWasmVM.Init")
}

func (n *NilVM) NewModule(path string) WasmModule {
	log.DefaultLogger.Errorf("NilWasmVM.NewModule return nil")
	return &NilWasmModule{}
}

type WasmerVM struct {
}

func (vm *WasmerVM) Init() {
	log.DefaultLogger.Infof("init Wasmer VM")
}

func (vm *WasmerVM) NewModule(path string) WasmModule {
	module := &WasmerModule{
		path: path,
		vm:   vm,
	}
	var err error
	module.code, err = wasm.ReadBytes(path)
	if err != nil {
		log.DefaultLogger.Errorf("WasmerVM.NewModule fail to read code bytes, err: %v", err)
		return &NilWasmModule{}
	}

	module.module, err = wasm.Compile(module.code)
	if err != nil {
		log.DefaultLogger.Errorf("WasmerVM.NewModule fail to compile wasm code, err: %v", err)
		return &NilWasmModule{}
	}
	module.wasiVersion = wasm.WasiGetVersion(module.module)

	module.wasmImportObj = wasm.NewDefaultWasiImportObjectForVersion(module.wasiVersion)

	im := ProxyWasmImports()

	err = module.wasmImportObj.Extend(*im)
	if err != nil {
		log.DefaultLogger.Errorf("WasmerVM.NewModule fail to extend import obj, err: %v", err)
		return &NilWasmModule{}
	}

	module.Init()

	return module
}

type WasmerModule struct {
	path          string
	code          []byte
	module        wasm.Module
	wasiVersion   wasm.WasiVersion
	wasmImportObj *wasm.ImportObject
	vm            *WasmerVM

	rootCtx *rootContext
}

func (m *WasmerModule) Init() {
	m.rootCtx = &rootContext{
		propertyMap: make(map[string]string),
		metrics:     newWasmMetricsManager("wasm"),
	}
}

var root_id = 100
var id int32

func (m *WasmerModule) NewInstance() *wasmContext {
	id++
	instanceCtx := &wasmContext{
		rootContext: m.rootCtx,
		contextId:   id,
	}

	instance, err := m.module.InstantiateWithImportObject(m.wasmImportObj)
	if err != nil {
		log.DefaultLogger.Errorf("wasm instance error :%v", err)
		instanceCtx.ProxyWasmInstance = &NilWasmInstance{}
		return instanceCtx
	} else {
		instanceCtx.ProxyWasmInstance = &WasmerInstance{instance}
	}

	// _start must be in the front of SetContextData, don't ask me why
	if err := instanceCtx._start(); err != nil {
		log.DefaultLogger.Errorf("wasm start err: %v\n", err)
		instanceCtx.ProxyWasmInstance = &NilWasmInstance{}
		return instanceCtx
	}

	instance.SetContextData(instanceCtx)

	if err := instanceCtx.proxy_on_context_create(int32(root_id), 0); err != nil {
		log.DefaultLogger.Errorf("root err %v\n", err)
		instanceCtx.ProxyWasmInstance = &NilWasmInstance{}
		return instanceCtx
	}

	return instanceCtx
}

type NilWasmModule struct{}

func (w *NilWasmModule) NewInstance() *wasmContext {
	log.DefaultLogger.Errorf("NilWasmModule.NewInstance return nil")
	return &wasmContext{
		ProxyWasmInstance: &NilWasmInstance{},
	}
}

type WasmtimeVM struct {
	engine *wasmtime.Engine
	store *wasmtime.Store
	rootCtx *rootContext
}

func (w *WasmtimeVM) Init() {
	w.engine = wasmtime.NewEngine()
	w.store = wasmtime.NewStore(w.engine)
	w.rootCtx = &rootContext{
		propertyMap: make(map[string]string),
		metrics:     newWasmMetricsManager("wasm"),
	}
}

func (w *WasmtimeVM) NewModule(path string) WasmModule {
	m := &WasmtimeModule{
		path:   path,
		vm:     w,
	}
	var err error
	m.module, err = wasmtime.NewModuleFromFile(w.engine, path)
	if err != nil {
		log.DefaultLogger.Errorf("WasmtimeVM.NewModule call NewModuleFromFile failed, err: %v", err)
		return &NilWasmModule{}
	}
	return m
}

type WasmtimeModule struct {
	path string
	vm *WasmtimeVM
	module *wasmtime.Module
}

func (w *WasmtimeModule) NewInstance() *wasmContext {

	f1 := wasmtime.NewFunc(w.vm.store,
		wasmtime.NewFuncType([]*wasmtime.ValType{}, []*wasmtime.ValType{}),
		)

	instance, err := wasmtime.NewInstance(w.vm.store, w.module, imports)
}
