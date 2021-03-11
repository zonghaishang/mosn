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

package wasm

import (
	"context"
	"fmt"

	mosnctx "mosn.io/mosn/pkg/context"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

func (proto *wasmProtocol) encodeRequest(context context.Context, request *Request) (types.IoBuffer, error) {
	ctx := mosnctx.Get(context, types.ContextKeyWasmContext)
	if ctx == nil {
		log.DefaultLogger.Errorf("[protocol] wasm %s encode request failed, wasm context not found.", proto.name)
		return nil, fmt.Errorf("wasm %s encode request failed, wasm context not found", proto.name)
	}

	wasmCtx := ctx.(*Context)
	wasmCtx.instance.Lock(wasmCtx.abi)
	wasmCtx.abi.SetABIImports(wasmCtx)
	// only for debug
	wasmCtx.SetEncodeCmd(request)
	// invoke plugin encode impl
	err := wasmCtx.exports.ProxyEncodeRequestBufferBytes(wasmCtx.contextId, request)
	wasmCtx.instance.Unlock()

	// clean plugin context
	proto.finishWasmContext(context)

	return wasmCtx.GetEncodeBuffer(), err
}

func (proto *wasmProtocol) encodeResponse(context context.Context, response *Response) (types.IoBuffer, error) {
	ctx := mosnctx.Get(context, types.ContextKeyWasmContext)
	if ctx == nil {
		log.DefaultLogger.Errorf("[protocol] wasm %s encode response failed, wasm context not found.", proto.name)
		return nil, fmt.Errorf("wasm %s encode response failed, wasm context not found", proto.name)
	}

	wasmCtx := ctx.(*Context)
	wasmCtx.instance.Lock(wasmCtx.abi)
	wasmCtx.abi.SetABIImports(wasmCtx)
	// only for debug
	wasmCtx.SetEncodeCmd(response)
	// invoke plugin encode impl
	err := wasmCtx.exports.ProxyEncodeResponseBufferBytes(wasmCtx.contextId, response)
	wasmCtx.instance.Unlock()

	// clean plugin context
	proto.finishWasmContext(context)

	return wasmCtx.GetEncodeBuffer(), err
}
