package abi

import (
	"mosn.io/mosn/pkg/types"
)

var abiMap = make(map[string]ABI)

type ABI interface {
	OnInstanceCreate(instance types.WasmInstance)
	OnStart(instance types.WasmInstance)
	OnInstanceDestroy(instance types.WasmInstance)

	SetInstance(instance types.WasmInstance)
	SetInstanceCallBack(callback interface{})
}

func RegisterABI(abiName string, abiImpl ABI) {
	abiMap[abiName] = abiImpl
}

func GetABI(abiName string) ABI {
	return abiMap[abiName]
}
