package proxywasm

import (
	"errors"

	"github.com/wasmerio/go-ext-wasm/wasmer"
	"mosn.io/mosn/pkg/log"
)

type WasmerInstance struct {
	wasmer.Instance
}

func (instance *WasmerInstance) GetMemory() []byte {
	return instance.Memory.Data()
}

func (instance *WasmerInstance) _start() error {
	log.DefaultLogger.Debugf("wasm call exported func: _start")
	ff := instance.Exports["_start"]
	if ff == nil {
		return errors.New("func _start not found")
	}
	_, err := ff()
	return err
}

func (instance *WasmerInstance) malloc(size int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: malloc")
	ff := instance.Exports["malloc"]
	if ff == nil {
		return 0, errors.New("func malloc not found")
	}
	addr, err := ff(size)
	if err != nil {
		return 0, err
	}
	return addr.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_context_create(contextId int32, parentContextId int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_context_create")
	ff := instance.Exports["proxy_on_context_create"]
	if ff == nil {
		return errors.New("func proxy_on_context_create not found")
	}
	_, err := ff(contextId, parentContextId)
	return err
}

func (instance *WasmerInstance) proxy_on_vm_start(rootContextId int32, configurationSize int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_vm_start")
	ff := instance.Exports["proxy_on_vm_start"]
	if ff == nil {
		return 0, errors.New("func proxy_on_vm_start not found")
	}
	res, err := ff(rootContextId, configurationSize)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_done(contextId int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_done")
	ff := instance.Exports["proxy_on_done"]
	if ff == nil {
		return 0, errors.New("func proxy_on_done not found")
	}
	res, err := ff(contextId)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_log(contextId int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_log")
	ff := instance.Exports["proxy_on_log"]
	if ff == nil {
		return errors.New("func proxy_on_log not found")
	}
	_, err := ff(contextId)
	return err
}

func (instance *WasmerInstance) proxy_on_delete(contextId int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_delete")
	ff := instance.Exports["proxy_on_delete"]
	if ff == nil {
		return errors.New("func proxy_on_delete not found")
	}
	_, err := ff(contextId)
	return err
}

func (instance *WasmerInstance) proxy_on_configure(rootContextId int32, configurationSize int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_configure")
	ff := instance.Exports["proxy_on_configure"]
	if ff == nil {
		return 0, errors.New("func proxy_on_configure not found")
	}
	res, err := ff(rootContextId, configurationSize)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_tick(rootContextId int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_tick")
	ff := instance.Exports["proxy_on_tick"]
	if ff == nil {
		return errors.New("func proxy_on_tick not found")
	}
	_, err := ff(rootContextId)
	return err
}

func (instance *WasmerInstance) proxy_on_new_connection(contextId int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_new_connection")
	ff := instance.Exports["proxy_on_new_connection"]
	if ff == nil {
		return errors.New("func proxy_on_new_connection not found")
	}
	_, err := ff(contextId)
	return err
}

func (instance *WasmerInstance) proxy_on_downstream_data(contextId int32, dataLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_downstream_data")
	ff := instance.Exports["proxy_on_downstream_data"]
	if ff == nil {
		return 0, errors.New("func proxy_on_downstream_data not found")
	}
	res, err := ff(contextId, dataLength, endOfStream)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_downstream_connection_close(contextId int32, closeType int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_downstream_connection_close")
	ff := instance.Exports["proxy_on_downstream_connection_close"]
	if ff == nil {
		return errors.New("func proxy_on_downstream_connection_close not found")
	}
	_, err := ff(contextId, closeType)
	return err
}

func (instance *WasmerInstance) proxy_on_upstream_data(contextId int32, dataLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_upstream_data")
	ff := instance.Exports["proxy_on_upstream_data"]
	if ff == nil {
		return 0, errors.New("func proxy_on_upstream_data not found")
	}
	res, err := ff(contextId, dataLength, endOfStream)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_upstream_connection_close(contextId int32, closeType int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_upstream_connection_close")
	ff := instance.Exports["proxy_on_upstream_connection_close"]
	if ff == nil {
		return errors.New("func proxy_on_upstream_connection_close not found")
	}
	_, err := ff(contextId, closeType)
	return err
}

func (instance *WasmerInstance) proxy_on_request_headers(contextId int32, headers int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_request_headers")
	ff := instance.Exports["proxy_on_request_headers"]
	if ff == nil {
		return 0, errors.New("func proxy_on_request_headers not found")
	}
	res, err := ff(contextId, headers, endOfStream)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_request_body(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_request_body")
	ff := instance.Exports["proxy_on_request_body"]
	if ff == nil {
		return 0, errors.New("func proxy_on_request_body not found")
	}
	res, err := ff(contextId, bodyBufferLength, endOfStream)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_request_trailers(contextId int32, trailers int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_request_trailers")
	ff := instance.Exports["proxy_on_request_trailers"]
	if ff == nil {
		return 0, errors.New("func proxy_on_request_trailers not found")
	}
	res, err := ff(contextId, trailers)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_request_metadata(contextId int32, nElements int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_request_metadata")
	ff := instance.Exports["proxy_on_request_metadata"]
	if ff == nil {
		return 0, errors.New("func proxy_on_request_metadata not found")
	}
	res, err := ff(contextId, nElements)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_response_headers(contextId int32, headers int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_response_headers")
	ff := instance.Exports["proxy_on_response_headers"]
	if ff == nil {
		return 0, errors.New("func proxy_on_response_headers not found")
	}
	res, err := ff(contextId, headers, endOfStream)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_response_body(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_response_body")
	ff := instance.Exports["proxy_on_response_body"]
	if ff == nil {
		return 0, errors.New("func proxy_on_response_body not found")
	}
	res, err := ff(contextId, bodyBufferLength, endOfStream)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_response_trailers(contextId int32, trailers int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_response_trailers")
	ff := instance.Exports["proxy_on_response_trailers"]
	if ff == nil {
		return 0, errors.New("func proxy_on_response_trailers not found")
	}
	res, err := ff(contextId, trailers)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_response_metadata(contextId int32, nElements int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_response_metadata")
	ff := instance.Exports["proxy_on_response_metadata"]
	if ff == nil {
		return 0, errors.New("func proxy_on_response_metadata not found")
	}
	res, err := ff(contextId, nElements)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_http_call_response(contextId int32, token int32, headers int32, bodySize int32, trailers int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_http_call_response")
	ff := instance.Exports["proxy_on_http_call_response"]
	if ff == nil {
		return errors.New("func proxy_on_http_call_response not found")
	}
	_, err := ff(contextId, token, headers, bodySize, trailers)
	return err
}

func (instance *WasmerInstance) proxy_on_grpc_receive_initial_metadata(contextId int32, token int32, headers int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_grpc_receive_initial_metadata")
	ff := instance.Exports["proxy_on_grpc_receive_initial_metadata"]
	if ff == nil {
		return errors.New("func proxy_on_grpc_receive_initial_metadata not found")
	}
	_, err := ff(contextId, token, headers)
	return err
}

func (instance *WasmerInstance) proxy_on_grpc_trailing_metadata(contextId int32, token int32, trailers int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_grpc_trailing_metadata")
	ff := instance.Exports["proxy_on_grpc_trailing_metadata"]
	if ff == nil {
		return errors.New("func proxy_on_grpc_trailing_metadata not found")
	}
	_, err := ff(contextId, token, trailers)
	return err
}

func (instance *WasmerInstance) proxy_on_grpc_receive(contextId int32, token int32, responseSize int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_grpc_receive")
	ff := instance.Exports["proxy_on_grpc_receive"]
	if ff == nil {
		return errors.New("func proxy_on_grpc_receive not found")
	}
	_, err := ff(contextId, token, responseSize)
	return err
}

func (instance *WasmerInstance) proxy_on_grpc_close(contextId int32, token int32, statusCode int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_grpc_close")
	ff := instance.Exports["proxy_on_grpc_close"]
	if ff == nil {
		return errors.New("func proxy_on_grpc_close not found")
	}
	_, err := ff(contextId, token, statusCode)
	return err
}

func (instance *WasmerInstance) proxy_on_queue_ready(rootContextId int32, token int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_queue_ready")
	ff := instance.Exports["proxy_on_queue_ready"]
	if ff == nil {
		return errors.New("func proxy_on_queue_ready not found")
	}
	_, err := ff(rootContextId, token)
	return err
}

func (instance *WasmerInstance) proxy_validate_configuration(rootContextId int32, configurationSize int32) (int32, error) {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_validate_configuration")
	ff := instance.Exports["proxy_validate_configuration"]
	if ff == nil {
		return 0, errors.New("func proxy_validate_configuration not found")
	}
	res, err := ff(rootContextId, configurationSize)
	if err != nil {
		return 0, err
	}
	return res.ToI32(), nil
}

func (instance *WasmerInstance) proxy_on_foreign_function(rootContextId int32, functionId int32, dataSize int32) error {
	log.DefaultLogger.Debugf("wasm call exported func: proxy_on_foreign_function")
	ff := instance.Exports["proxy_on_foreign_function"]
	if ff == nil {
		return errors.New("func proxy_on_foreign_function not found")
	}
	_, err := ff(rootContextId, functionId, dataSize)
	return err
}
