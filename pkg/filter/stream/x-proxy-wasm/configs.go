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
	"encoding/json"

	v2 "mosn.io/mosn/pkg/config/v2"
)

type filterConfig struct {
	FromWasmPlugin string           `json:"from_wasm_plugin,omitempty"`
	VmConfig       *v2.WasmVmConfig `json:"vm_config,omitempty"`
	InstanceNum    int              `json:"instance_num,omitempty"`
	RootContextID  int32            `json:"root_context_id, omitempty"`
	pluginConfig
}

type pluginConfig struct {
	UserConfig1 string `json:"user_config1,omitempty"`
	UserConfig2 string `json:"user_config2,omitempty"`
}

func parseFilterConfig(cfg map[string]interface{}) (*filterConfig, error) {
	var filterConfig filterConfig

	data, err := json.Marshal(cfg["config"])
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &filterConfig); err != nil {
		return nil, err
	}

	if filterConfig.FromWasmPlugin != "" {
		filterConfig.VmConfig = nil
	}

	return &filterConfig, nil
}
