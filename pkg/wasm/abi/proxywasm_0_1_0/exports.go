/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proxywasm_0_1_0

import "mosn.io/mosn/pkg/log"

func (a *abiImpl) ProxyOnContextCreate(contextId int32, parentContextId int32) error {
	log.DefaultLogger.Infof("[proxywasm_0_1_0][export] ProxyOnContextCreate contextID: %v, parentContextId: %v", contextId, parentContextId)

	ff, err := a.instance.GetExportsFunc("proxy_on_context_create")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_context_create, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, parentContextId)
	if err != nil {
		return err
	}
	return nil
}
func (a *abiImpl) ProxyOnDone(contextId int32) (int32, error) {
	log.DefaultLogger.Infof("[proxywasm_0_1_0][export] ProxyOnDone contextID: %v", contextId)

	ff, err := a.instance.GetExportsFunc("proxy_on_done")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_done, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnLog(contextId int32) error {
	log.DefaultLogger.Debugf("[wasmer][instance] ProxyOnLog")

	ff, err := a.instance.GetExportsFunc("proxy_on_log")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_log, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId)

	return err
}

func (a *abiImpl) ProxyOnVmStart(rootContextId int32, vmConfigurationSize int32) (int32, error) {
	log.DefaultLogger.Infof("[proxywasm_0_1_0][export] ProxyOnVmStart rootContextId: %v, vmConfigurationSize: %v", rootContextId, vmConfigurationSize)

	ff, err := a.instance.GetExportsFunc("proxy_on_vm_start")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_vm_start, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(rootContextId, vmConfigurationSize)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}
func (a *abiImpl) ProxyOnDelete(contextId int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnDelete")

	ff, err := a.instance.GetExportsFunc("proxy_on_delete")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_delete, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId)

	return err
}

func (a *abiImpl) ProxyOnConfigure(rootContextId int32, configurationSize int32) (int32, error) {
	log.DefaultLogger.Infof("[proxywasm_0_1_0][export] ProxyOnConfigure rootContextId: %v, configurationSize: %v", rootContextId, configurationSize)

	ff, err := a.instance.GetExportsFunc("proxy_on_configure")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_configure, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(rootContextId, configurationSize)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnTick(rootContextId int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnTick")

	ff, err := a.instance.GetExportsFunc("proxy_on_tick")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_tick, err: %v", err)
		return err
	}

	_, err = ff.Call(rootContextId)

	return err
}

func (a *abiImpl) ProxyOnNewConnection(contextId int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnNewConnection")

	ff, err := a.instance.GetExportsFunc("proxy_on_new_connection")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_new_connection, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId)

	return err
}

func (a *abiImpl) ProxyOnDownstreamData(contextId int32, dataLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnDownstreamData")

	ff, err := a.instance.GetExportsFunc("proxy_on_downstream_data")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_downstream_data, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, dataLength, endOfStream)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnDownstreamConnectionClose(contextId int32, closeType int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnDownstreamConnectionClose")

	ff, err := a.instance.GetExportsFunc("proxy_on_downstream_connection_close")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_downstream_connection_close, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, closeType)

	return err
}

func (a *abiImpl) ProxyOnUpstreamData(contextId int32, dataLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnUpstreamData")

	ff, err := a.instance.GetExportsFunc("proxy_on_upstream_data")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_upstream_data, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, dataLength, endOfStream)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnUpstreamConnectionClose(contextId int32, closeType int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnUpstreamConnectionClose")

	ff, err := a.instance.GetExportsFunc("proxy_on_upstream_connection_close")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_upstream_connection_close, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, closeType)

	return err
}

func (a *abiImpl) ProxyOnRequestHeaders(contextID int32, numHeaders int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnRequestHeaders")

	ff, err := a.instance.GetExportsFunc("proxy_on_request_headers")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_request_headers, err: %v", err)
		return 0, err
	}

	ret, err := ff.Call(contextID, numHeaders, endOfStream)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to call proxy_on_request_headers func, err: %v", err)
		return 0, err
	}

	return ret.(int32), nil
}

func (a *abiImpl) ProxyOnRequestBody(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnRequestBody")

	ff, err := a.instance.GetExportsFunc("proxy_on_request_body")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_request_body, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, bodyBufferLength, endOfStream)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnRequestTrailers(contextId int32, trailers int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnRequestTrailers")

	ff, err := a.instance.GetExportsFunc("proxy_on_request_trailers")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_request_trailers, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, trailers)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnRequestMetadata(contextId int32, nElements int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnRequestMetadata")

	ff, err := a.instance.GetExportsFunc("proxy_on_request_metadata")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_request_metadata, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, nElements)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnResponseHeaders(contextId int32, headers int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnResponseHeaders")

	ff, err := a.instance.GetExportsFunc("proxy_on_response_headers")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_response_headers, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, headers, endOfStream)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnResponseBody(contextId int32, bodyBufferLength int32, endOfStream int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnResponseBody")

	ff, err := a.instance.GetExportsFunc("proxy_on_response_body")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_response_body, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, bodyBufferLength, endOfStream)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnResponseTrailers(contextId int32, trailers int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnResponseTrailers")

	ff, err := a.instance.GetExportsFunc("proxy_on_response_trailers")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_response_trailers, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, trailers)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnResponseMetadata(contextId int32, nElements int32) (int32, error) {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnResponseMetadata")

	ff, err := a.instance.GetExportsFunc("proxy_on_response_metadata")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_response_metadata, err: %v", err)
		return 0, err
	}

	res, err := ff.Call(contextId, nElements)
	if err != nil {
		return 0, err
	}

	return res.(int32), nil
}

func (a *abiImpl) ProxyOnHttpCallResponse(contextId int32, token int32, headers int32, bodySize int32, trailers int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnHttpCallResponse")

	ff, err := a.instance.GetExportsFunc("proxy_on_http_call_response")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_http_call_response, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, token, headers, bodySize, trailers)

	return err
}

func (a *abiImpl) ProxyOnGrpcReceiveInitialMetadata(contextId int32, token int32, headers int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnGrpcReceiveInitialMetadata")

	ff, err := a.instance.GetExportsFunc("proxy_on_grpc_receive_initial_metadata")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_grpc_receive_initial_metadata, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, token, headers)

	return err
}

func (a *abiImpl) ProxyOnGrpcTrailingMetadata(contextId int32, token int32, trailers int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnGrpcTrailingMetadata")

	ff, err := a.instance.GetExportsFunc("proxy_on_grpc_trailing_metadata")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_grpc_trailing_metadata, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, token, trailers)

	return err
}

func (a *abiImpl) ProxyOnGrpcReceive(contextId int32, token int32, responseSize int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnGrpcReceive")

	ff, err := a.instance.GetExportsFunc("proxy_on_grpc_receive")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_grpc_receive, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, token, responseSize)

	return err
}

func (a *abiImpl) ProxyOnGrpcClose(contextId int32, token int32, statusCode int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnGrpcClose")

	ff, err := a.instance.GetExportsFunc("proxy_on_grpc_close")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_grpc_close, err: %v", err)
		return err
	}

	_, err = ff.Call(contextId, token, statusCode)

	return err
}

func (a *abiImpl) ProxyOnQueueReady(rootContextId int32, token int32) error {
	log.DefaultLogger.Debugf("[proxywasm_0_1_0][export] WasmerInstance ProxyOnQueueReady")

	ff, err := a.instance.GetExportsFunc("proxy_on_queue_ready")
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm_0_1_0][export] WasmerInstance fail to get export func: proxy_on_queue_ready, err: %v", err)
		return err
	}

	_, err = ff.Call(rootContextId, token)

	return err
}
