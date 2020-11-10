package proxywasm

import (
	"context"
	"encoding/json"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/buffer"
)

func init() {
	api.RegisterStream(ProxyWasm, CreateProxyWasmFilterFactory)
}

const ProxyWasm = "proxy-wasm"

type StreamProxyWasmConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type FilterConfigFactory struct {
	Config *StreamProxyWasmConfig
}

func CreateProxyWasmFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	log.DefaultLogger.Debugf("create proxy wasm stream filter factory")
	cfg, err := ParseStreamProxyWasmFilter(conf)
	if err != nil {
		return nil, err
	}

	initWasmVM(cfg)
	rootWasmInstance = NewWasmInstance()
	if rootWasmInstance == nil {
		log.DefaultLogger.Errorf("wasm init error")
		return &FilterConfigFactory{cfg}, nil
	}

	if _, err := rootWasmInstance.proxy_on_vm_start(int32(root_id), 1000); err != nil {
		log.DefaultLogger.Errorf("start err %v\n", err)
	}

	log.DefaultLogger.Debugf("wasm init %+v", rootWasmInstance)

	return &FilterConfigFactory{cfg}, nil
}

// ParseStreamPayloadLimitFilter
func ParseStreamProxyWasmFilter(cfg map[string]interface{}) (*StreamProxyWasmConfig, error) {
	filterConfig := &StreamProxyWasmConfig{}
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, filterConfig); err != nil {
		return nil, err
	}
	return filterConfig, nil
}

func (f *FilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {

	filter := NewFilter(context, f.Config)
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter)

	log.DefaultLogger.Debugf("wasm filter init success")
}

// streamProxyWasmFilter is an implement of StreamReceiverFilter
type streamProxyWasmFilter struct {
	ctx      context.Context
	rhandler api.StreamReceiverFilterHandler
	shandler api.StreamSenderFilterHandler
	config   *StreamProxyWasmConfig
	*wasmContext
}

func NewFilter(ctx context.Context, config *StreamProxyWasmConfig) *streamProxyWasmFilter {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("create a new proxy wasm filter")
	}
	filter := &streamProxyWasmFilter{
		ctx:         ctx,
		config:      config,
		wasmContext: NewWasmInstance(),
	}
	filter.wasmContext.filter = filter

	if err := filter.proxy_on_context_create(filter.contextId, int32(root_id)); err != nil {
		log.DefaultLogger.Errorf("root err %v\n", err)
		return nil
	}

	return filter
}

func (f *streamProxyWasmFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.rhandler = handler
}

func (f *streamProxyWasmFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.shandler = handler
}

func (f *streamProxyWasmFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("proxy wasm stream do receive headers, id = %d", f.contextId)
	}
	if buf != nil && buf.Len() > 0 {
		if _, err := f.proxy_on_request_headers(f.contextId, 0, 0); err != nil {
			log.DefaultLogger.Errorf("wasm proxy_on_request_headers err: %v", err)
		}
		if _, err := f.proxy_on_request_body(f.contextId, int32(buf.Len()), 1); err != nil {
			log.DefaultLogger.Errorf("wasm proxy_on_request_body err: %v", err)
		}
	} else {
		if _, err := f.proxy_on_request_headers(f.contextId, 0, 1); err != nil {
			log.DefaultLogger.Errorf("wasm proxy_on_request_headers err: %v", err)
		}
	}

	return api.StreamFilterContinue
}

func (f *streamProxyWasmFilter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("proxy wasm stream do receive headers")
	}

	if _, err := f.proxy_on_response_headers(f.contextId, 1, 0); err != nil {
		log.DefaultLogger.Errorf("wasm proxy_on_response_headers err: %v", err)
	}

	return api.StreamFilterContinue
}

func (f *streamProxyWasmFilter) OnDestroy() {
	_ = f.proxy_on_log(f.contextId)
	_, _ = f.proxy_on_done(f.contextId)
	_ = f.proxy_on_delete(f.contextId)
}
