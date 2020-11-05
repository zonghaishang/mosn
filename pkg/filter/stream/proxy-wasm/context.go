package proxywasm

import (
	"errors"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)

type rootContext struct {
	config *StreamProxyWasmConfig

	//vmConfig     string
	//pluginConfig string
	//
	//contextId    uint32

	wasmCode      []byte
	wasmModule    wasm.Module
	wasiVersion   wasm.WasiVersion
	wasmImportObj *wasm.ImportObject
}

type wasmContext struct {
	contextId int32
	filter    *streamProxyWasmFilter
	instance  *wasm.Instance
}

func (wasm *wasmContext) _start() error {
	log.DefaultLogger.Debugf("wasm call exported func: _start")
	ff := wasm.instance.Exports["_start"]
	if ff == nil {
		return errors.New("func _start not found")
	}
	_, err := ff()
	return err
}

func (wasm *wasmContext) proxy_on_context_create(contextId int32, parentContextId int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_context_create")
	ff := wasm.instance.Exports["proxy_on_context_create"]
	if ff == nil {
		return errors.New("func proxy_on_context_create not found")
	}
	_, err := ff(contextId, parentContextId)
	return err
}

func (wasm *wasmContext) proxy_on_vm_start(rootContextId int32, configurationSize int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_vm_start")
	ff := wasm.instance.Exports["proxy_on_vm_start"]
	if ff == nil {
		return 0, errors.New("func proxy_on_vm_start not found")
	}
	res, err := ff(rootContextId, configurationSize)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}
