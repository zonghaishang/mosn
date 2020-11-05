package proxywasm

import wasm "github.com/wasmerio/go-ext-wasm/wasmer"

type rootContext struct {
	config       *StreamProxyWasmConfig

	//vmConfig     string
	//pluginConfig string
	//
	//contextId    uint32

	wasmCode []byte
	wasmModule wasm.Module
	wasiVersion wasm.WasiVersion
	wasmImportObj *wasm.ImportObject
}


type wasmContext struct {
	contextId int32
	filter   *streamProxyWasmFilter
	instance *wasm.Instance
}