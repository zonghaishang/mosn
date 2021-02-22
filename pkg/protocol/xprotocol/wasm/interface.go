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
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/types"
	v1 "mosn.io/mosn/pkg/wasm/abi/proxywasm_0_1_0"
	"mosn.io/pkg/buffer"
)

type Exports interface {
	v1.Exports
	ProxyDecodeBufferBytes(contextId int32, buf types.IoBuffer) error
	ProxyEncodeRequestBufferBytes(contextId int32, cmd *Request) error
	ProxyEncodeResponseBufferBytes(contextId int32, cmd *Response) error

	ProxyKeepAliveBufferBytes(contextId int32, id uint64) error
	ProxyReplyKeepAliveBufferBytes(contextId int32, cmd *Request) error
	ProxyHijackBufferBytes(contextId int32, cmd *Request, statusCode uint32) error
}

type ContextCallback interface {
	// extension for abi 0_1_0
	v1.InstanceCallback

	//DrainLength() uint32
	GetDecodeCmd() xprotocol.XFrame
	SetDecodeCmd(cmd xprotocol.XFrame)
	GetDecodeBuffer() buffer.IoBuffer

	GetEncodeCmd() xprotocol.XFrame
	SetEncodeBuffer(buf buffer.IoBuffer)
	GetEncodeBuffer() buffer.IoBuffer
}

var contextId int32

type Context struct {
	v1.DefaultInstanceCallback
	decodeCmd     xprotocol.XFrame
	decodeBuffer  buffer.IoBuffer
	encodeCmd     xprotocol.XFrame
	encodeBuffer  buffer.IoBuffer
	proto         *wasmRpcProtocol
	keepaliveReq  *Request
	keepaliveResp *Response
	contextId     int32
	exports       Exports
}

func (c *Context) GetDecodeCmd() xprotocol.XFrame {
	return c.decodeCmd
}

func (c *Context) SetDecodeCmd(cmd xprotocol.XFrame) {
	c.decodeCmd = cmd
}

func (c *Context) GetDecodeBuffer() buffer.IoBuffer {
	return c.decodeBuffer
}

func (c *Context) SetDecodeBuffer(buf buffer.IoBuffer) {
	c.decodeBuffer = buf
}

func (c *Context) GetEncodeCmd() xprotocol.XFrame {
	return c.encodeCmd
}

func (c *Context) SetEncodeBuffer(buf buffer.IoBuffer) {
	c.encodeBuffer = buf
}

func (c *Context) GetEncodeBuffer() buffer.IoBuffer {
	return c.encodeBuffer
}

func (c *Context) ContextId() int32 {
	return c.contextId
}
