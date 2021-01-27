package test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/mock"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm_0_1_0"
	_ "mosn.io/mosn/pkg/wasm/runtime/wasmer"
	"mosn.io/pkg/buffer"
)

type mockInstanceCallback struct {
	ctrl           *gomock.Controller
	requestHeader  api.HeaderMap
	requestBody    buffer.IoBuffer
	responseHeader api.HeaderMap
	responseBody   buffer.IoBuffer
	vmConfig       buffer.IoBuffer
	pluginConfig   buffer.IoBuffer
}

func newMockInstanceCallback(ctrl *gomock.Controller) *mockInstanceCallback {
	var m = map[string]string{
		"requestHeaderKey1": "requestHeaderValue1",
		"requestHeaderKey2": "requestHeaderValue2",
		"requestHeaderKey3": "requestHeaderValue3",
	}
	h := mock.NewMockHeaderMap(ctrl)
	h.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	})
	h.EXPECT().Del(gomock.Any()).AnyTimes().DoAndReturn(func(key string) { delete(m, key) })
	h.EXPECT().Add(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(key string, val string) { m[key] = val })
	h.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(key string, val string) { m[key] = val })
	h.EXPECT().Range(gomock.Any()).AnyTimes().Do(func(f func(key, value string) bool) {
		for k, v := range m {
			if !f(k, v) {
				break
			}
		}
	})
	h.EXPECT().ByteSize().AnyTimes().DoAndReturn(func() uint64 {
		var size uint64
		for k, v := range m {
			size += uint64(len(k) + len(v))
		}
		return size
	})

	return &mockInstanceCallback{
		ctrl:           ctrl,
		requestHeader:  h,
		requestBody:    buffer.NewIoBufferString("request body"),
		responseHeader: nil,
		responseBody:   buffer.NewIoBufferString("response body"),
		vmConfig:       buffer.NewIoBufferString("vm config"),
		pluginConfig:   buffer.NewIoBufferString("plugin config"),
	}
}

func (i *mockInstanceCallback) GetVmConfig() buffer.IoBuffer {
	return i.vmConfig
}

func (i *mockInstanceCallback) GetPluginConfig() buffer.IoBuffer {
	return i.pluginConfig
}

func (i *mockInstanceCallback) GetHttpRequestHeader() api.HeaderMap {
	return i.requestHeader
}

func (i *mockInstanceCallback) GetHttpRequestBody() buffer.IoBuffer {
	return i.requestBody
}

func (i *mockInstanceCallback) GetHttpRequestTrailer() api.HeaderMap {
	return nil
}

func (i *mockInstanceCallback) GetHttpResponseHeader() api.HeaderMap {
	return i.responseHeader
}

func (i *mockInstanceCallback) GetHttpResponseBody() buffer.IoBuffer {
	return i.responseBody
}

func (i *mockInstanceCallback) GetHttpResponseTrailer() api.HeaderMap {
	return nil
}

func (i *mockInstanceCallback) Log(level log.Level, msg string) {
	logFunc := log.DefaultLogger.Infof
	switch level {
	case log.TRACE:
		logFunc = log.DefaultLogger.Tracef
	case log.DEBUG:
		logFunc = log.DefaultLogger.Debugf
	case log.INFO:
		logFunc = log.DefaultLogger.Debugf // TODO: info -> debug
	case log.WARN:
		logFunc = log.DefaultLogger.Warnf
	case log.ERROR:
		logFunc = log.DefaultLogger.Errorf
	case log.FATAL:
		logFunc = log.DefaultLogger.Fatalf
	}
	logFunc(msg)
}

func testCommon(t *testing.T, pluginName string, engine string, path string) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	manager := wasm.GetWasmManager()
	_ = manager.AddOrUpdateWasm(v2.WasmPluginConfig{
		PluginName: pluginName,
		VmConfig: v2.WasmVmConfig{
			Engine: engine,
			Path:   path,
		},
		InstanceNum: 1,
	})

	plugin := manager.GetWasmPluginWrapperByName(pluginName).GetPlugin()
	v := plugin.GetInstance()

	cb := newMockInstanceCallback(ctrl)

	proxyWasm := abi.GetABI("proxy_abi_version_0_1_0")
	proxyWasm.SetInstance(v)
	proxyWasm.SetInstanceCallBack(cb)

	exports := proxyWasm.(proxywasm_0_1_0.Exports)

	rootContextID := 100
	contextID := 101

	if err := exports.ProxyOnContextCreate(int32(contextID), int32(rootContextID)); err != nil {
		t.Errorf("fail to call proxyOnContextCreate, err: %v", err)
	}

	_, err := exports.ProxyOnRequestHeaders(int32(contextID), 0, 1)
	if err != nil {
		t.Errorf("on request headers fail, err: %v", err)
	}

	_, err = exports.ProxyOnDone(int32(contextID))
	if err != nil {
		t.Errorf("on done err: %v", err)
	}
}

func TestWasmProxyLog(t *testing.T) {
	testCommon(t, "testWasmProxyLog", "wasmer", "./data/log.wasm")
}

func TestWasmHttp(t *testing.T) {
	testCommon(t, "testWasmHttpFull", "wasmer", "./data/httpFull.wasm")
}

func benchCommon(b *testing.B, pluginName string, engine string, path string) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	manager := wasm.GetWasmManager()
	_ = manager.AddOrUpdateWasm(v2.WasmPluginConfig{
		PluginName: pluginName,
		VmConfig: v2.WasmVmConfig{
			Engine: engine,
			Path:   path,
		},
		InstanceNum: 1,
	})

	cb := newMockInstanceCallback(ctrl)

	for i := 0; i < b.N; i++ {
		plugin := manager.GetWasmPluginWrapperByName(pluginName).GetPlugin()
		v := plugin.GetInstance()

		proxyWasm := abi.GetABI("proxy_abi_version_0_1_0")
		proxyWasm.SetInstance(v)
		proxyWasm.SetInstanceCallBack(cb)

		exports := proxyWasm.(proxywasm_0_1_0.Exports)

		contextID := 101 + i

		_ = exports.ProxyOnContextCreate(int32(contextID), int32(100))

		_, _ = exports.ProxyOnRequestHeaders(int32(contextID), 0, 1)

		_, _ = exports.ProxyOnDone(int32(contextID))

		plugin.ReleaseInstance(v)
	}
}

func BenchmarkWasmEmptyCall(b *testing.B) {
	benchCommon(b, "benchPluginEmptyCall", "wasmer", "./data/emptyCall.wasm")
}

func BenchmarkWasmProxyHttp(b *testing.B) {
	benchCommon(b, "benchPluginProxyHttp", "wasmer", "./data/httpFull.wasm")
}
