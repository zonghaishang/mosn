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

package x_proxy_wasm

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/mosn/pkg/mock"
	_ "mosn.io/mosn/pkg/wasm/runtime/wasmer"
	"mosn.io/pkg/buffer"
)

func mockHeaderMap(ctrl *gomock.Controller) api.HeaderMap {
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

	return h
}

func TestProxyWasmStreamFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configMap := map[string]interface{}{
		"type": "x-proxy-wasm",
		"config": map[string]interface{}{
			"instance_num": 1,
			"vm_config": map[string]interface{}{
				"engine": "wasmer",
				"path":   "./data/test.wasm",
				"cpu":    50,
				"mem":    50,
			},
			"root_context_id": 1,
			"user_config1":    "user_value1",
			"user_config2":    "user_value2",
		},
	}

	factory, err := createProxyWasmFilterFactory(configMap)
	if err != nil || factory == nil {
		t.Errorf("fail to create filter factory")
		return
	}

	var rFilter api.StreamReceiverFilter
	var sFilter api.StreamSenderFilter

	cb := mock.NewMockStreamFilterChainFactoryCallbacks(ctrl)
	cb.EXPECT().AddStreamReceiverFilter(gomock.Any(), gomock.Any()).Do(func(receiverFilter api.StreamReceiverFilter, p api.ReceiverFilterPhase) {
		assert.Equal(t, p, api.BeforeRoute, "add receiver filter at wrong phase")
		rFilter = receiverFilter
	}).AnyTimes()
	cb.EXPECT().AddStreamSenderFilter(gomock.Any(), gomock.Any()).Do(func(senderFilter api.StreamSenderFilter, p api.SenderFilterPhase) {
		assert.Equal(t, p, api.BeforeSend, "add sender filter at wrong phase")
		sFilter = senderFilter
	}).AnyTimes()

	factory.CreateFilterChain(context.TODO(), cb)

	// for coverage
	rFilter.SetReceiveFilterHandler(nil)
	sFilter.SetSenderFilterHandler(nil)

	reqHeaderMap := mockHeaderMap(ctrl)

	rFilter.OnReceive(context.TODO(), reqHeaderMap, buffer.NewIoBufferString("request body"), nil)
	sFilter.Append(context.TODO(), nil, buffer.NewIoBufferString("response body"), nil)

	rFilter.OnDestroy()
	sFilter.OnDestroy()

}
