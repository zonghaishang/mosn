package proxywasm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	admin "mosn.io/mosn/pkg/admin/server"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

var (
	ErrNilWasmConfig  = errors.New("wasm config is nil")
	ErrUnexpectedType = errors.New("unexpected wasm config type")
)

func init() {
	admin.RegisterAdminHandleFunc("/api/v1/updateWasmConfig", updateWasmConfig)
}

func updateWasmConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.DefaultLogger.Alertf(types.ErrorKeyAdmin, "api: %s, error: invalid method: %s", "updateWasmConfig", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.DefaultLogger.Alertf(types.ErrorKeyAdmin, "api: %s, error: read body failed, %v", "updateWasmConfig", err)
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("read body error: %v", err)
		fmt.Fprint(w, msg)
		return
	}

	var config WasmConfigBase
	err = json.Unmarshal(body, &config)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		msg := fmt.Sprintf("json decode error: %v", err)
		fmt.Fprint(w, msg)
	} else {
		log.DefaultLogger.Infof("updateWasmConfig: %v", config)

		configManager := GetWasmConfigManger()
		err = configManager.AddOrUpdateWasmConfig(&config)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			msg := fmt.Sprintf("update wasm config error: %v", err)
			fmt.Fprint(w, msg)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

type WasmConfig interface {
	GetName() string
	GetRuntime() string
	SetRuntime(runtime string)
	GetPath() string
	SetPath(path string)
	GetData() interface{}
	SetData(data interface{})
}

type WasmConfigBase struct {
	Name string `json:"name"`
	WasmVmConfig
}

func (base *WasmConfigBase) GetName() string {
	return base.Name
}

func (base *WasmConfigBase) GetData() interface{} {
	return nil
}

func (base *WasmConfigBase) SetData(data interface{}) {
	return
}

type WasmVmConfig struct {
	Runtime string `json:"runtime"`
	Path    string `json:"path"`
}

func (vmConfig *WasmVmConfig) GetRuntime() string {
	return vmConfig.Runtime
}

func (vmConfig *WasmVmConfig) SetRuntime(runtime string) {
	vmConfig.Runtime = runtime
}

func (vmConfig *WasmVmConfig) GetPath() string {
	return vmConfig.Path
}

func (vmConfig *WasmVmConfig) SetPath(path string) {
	vmConfig.Path = path
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

	if v, ok := w.configMap.Load(config.GetName()); ok {
		oldConfig, ok := v.(*WasmConfigWrapper)
		if !ok {
			log.DefaultLogger.Errorf("wasmConfigManagerImpl.AddOrUpdateWasmConfig unexpected type in map, name: %v", config.GetName())
			return ErrUnexpectedType
		}

		if newConfig.GetRuntime() == "" {
			newConfig.SetRuntime(oldConfig.config.GetRuntime())
		}
		if newConfig.GetPath() == "" {
			newConfig.SetPath(oldConfig.config.GetPath())
		}
		if newConfig.GetData() == nil {
			newConfig.SetData(oldConfig.config.GetData())
		}

		// only update config data
		if newConfig.GetRuntime() == oldConfig.config.GetRuntime() &&
			newConfig.GetPath() == oldConfig.config.GetPath() {
			w.configMap.Store(newConfig.GetName(), &WasmConfigWrapper{
				wasmModule: oldConfig.wasmModule,
				config:     newConfig,
			})
			return nil
		}
	}

	vm := NewWasmVM(newConfig.GetRuntime())
	module := vm.NewModule(newConfig.GetPath())
	w.configMap.Store(newConfig.GetName(), &WasmConfigWrapper{
		wasmModule: module,
		config:     newConfig,
	})

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

	wasmConfigManager = &wasmConfigManagerImpl{
		configMap: sync.Map{},
	}
	return wasmConfigManager
}
