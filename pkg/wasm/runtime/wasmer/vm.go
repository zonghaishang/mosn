package wasmer

import (
	"io/ioutil"

	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
)

func init() {
	wasm.RegisterWasmEngine("wasmer", NewWasmerVM())
}

type VM struct {
	engine *wasmerGo.Engine
	store  *wasmerGo.Store
}

func NewWasmerVM() types.WasmVM {
	vm := &VM{}
	vm.Init()
	return vm
}

func (w *VM) Init() {
	w.engine = wasmerGo.NewEngine()
	w.store = wasmerGo.NewStore(w.engine)
}

func (w *VM) NewModule(path string) types.WasmModule {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmer][vm] fail to read wasm file, path: %v, err: %v", path, err)
		return nil
	}

	m, err := wasmerGo.NewModule(w.store, bytes)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmer][vm] fail to new module, err: %v", err)
		return nil
	}

	module := &Module{
		vm:     w,
		module: m,
	}

	module.Init()

	return module
}

