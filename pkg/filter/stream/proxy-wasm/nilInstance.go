package proxywasm

type NilWasmInstance struct {}

func (instance *NilWasmInstance) GetMemory() []byte {
	return nil
}

func (instance *NilWasmInstance) _start() error {
	return nil
}

func (instance *NilWasmInstance) malloc(size int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_context_create(contextId int32, parentContextId int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_done(contextId int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_log(contextId int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_delete(contextId int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_vm_start(rootContextId int32, configurationSize int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_configure(rootContextId int32, configurationSize int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_tick(rootContextId int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_new_connection(contextId int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_downstream_data(contextId int32, dataLength int32, endOfStream int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_downstream_connection_close(contextId int32, closeType int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_upstream_data(contextId int32, dataLength int32, endOfStream int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_upstream_connection_close(contextId int32, closeType int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_request_headers(contextId int32, headers int32, endOfStream int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_request_body(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_request_trailers(contextId int32, trailers int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_request_metadata(contextId int32, nElements int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_response_headers(contextId int32, headers int32, endOfStream int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_response_body(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_response_trailers(contextId int32, trailers int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_response_metadata(contextId int32, nElements int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_http_call_response(contextId int32, token int32, headers int32, bodySize int32, trailers int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_grpc_receive_initial_metadata(contextId int32, token int32, headers int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_grpc_trailing_metadata(contextId int32, token int32, trailers int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_grpc_receive(contextId int32, token int32, responseSize int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_grpc_close(contextId int32, token int32, statusCode int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_on_queue_ready(rootContextId int32, token int32) error {
	return nil
}

func (instance *NilWasmInstance) proxy_validate_configuration(rootContextId int32, configurationSize int32) (int32, error) {
	return 0, nil
}

func (instance *NilWasmInstance) proxy_on_foreign_function(rootContextId int32, functionId int32, dataSize int32) error {
	return nil
}



