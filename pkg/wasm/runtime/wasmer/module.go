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
