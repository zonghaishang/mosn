package proxywasm

import (
	"errors"
	"sync"

	"mosn.io/mosn/pkg/log"
)

var (
	ErrNilWasmConfig  = errors.New("wasm config is nil")
	ErrUnexpectedType = errors.New("unexpected wasm config type")
)

type WasmConfig interface {
	GetName() string
	GetRuntime() string
	GetPath() string
	GetData() interface{}
}

type WasmConfigBase struct {
	Name         string `json:"name"`
	WasmVmConfig
}

func (base *WasmConfigBase) GetName() string {
	return base.Name
}

type WasmVmConfig struct {
	Runtime string `json:"runtime"`
	Path    string `json:"path"`
}

func (vmConfig *WasmVmConfig) GetRuntime() string {
	return vmConfig.Runtime
}

func (vmConfig *WasmVmConfig) GetPath() string {
	return vmConfig.Path
}

type WasmConfigWrapper struct {
	wasmModule WasmModule
	config     WasmConfig
}

type WasmConfigManager interface {
	AddOrUpdateWasmConfig(config WasmConfig) error
	GetWasmConfigByName(name string) *WasmConfigWrapper
}

type wasmConfigManagerImpl struct {
	configMap sync.Map
}

func (w *wasmConfigManagerImpl) AddOrUpdateWasmConfig(config WasmConfig) error {
	if config == nil {
		log.DefaultLogger.Errorf("wasmConfigManagerImpl.AddOrUpdateWasmConfig nil config")
		return ErrNilWasmConfig
	}
	newConfig := config
	if _, ok := w.configMap.Load(config.GetName()); ok {
		//oldConfig, ok := v.(WasmConfig)
		//if !ok {
		//	log.DefaultLogger.Errorf("wasmConfigManagerImpl.AddOrUpdateWasmConfig unexpected type in map, name: %v", config.Name())
		//	return ErrUnexpectedType
		//}
	} else {
		vm := NewWasmVM(newConfig.GetRuntime())
		module := vm.NewModule(newConfig.GetPath())
		w.configMap.Store(newConfig.GetName(), &WasmConfigWrapper{
			wasmModule: module,
			config:     newConfig,
		})
	}
	return nil
}

func (w *wasmConfigManagerImpl) GetWasmConfigByName(name string) *WasmConfigWrapper {
	if v, ok := w.configMap.Load(name); ok {
		config, ok := v.(*WasmConfigWrapper)
		if !ok {
			log.DefaultLogger.Errorf("wasmConfigManagerImpl.GetWasmConfigByName unexpected object in routers map, name: %v", name)
			return nil
		}
		return config
	}
	return nil
}

var (
	wasmConfigSingleton sync.Mutex
	wasmConfigManager   *wasmConfigManagerImpl
)

func GetWasmConfigManger() WasmConfigManager {
	wasmConfigSingleton.Lock()
	defer wasmConfigSingleton.Unlock()

	if wasmConfigManager != nil {
		return wasmConfigManager
	}

	wasmConfigManager := &wasmConfigManagerImpl{
		configMap: sync.Map{},
	}
	return wasmConfigManager
}
