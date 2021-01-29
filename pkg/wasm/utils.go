package wasm

import (
	"io/ioutil"

	"mosn.io/mosn/pkg/log"
)

func loadWasmBytesFromPath(path string) []byte {
	if path == "" {
		return nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.DefaultLogger.Errorf("[wasm][utils] loadWasmBytesFromPath read file err: %v", err)
		return nil
	}

	return bytes
}

func loadWasmBytesFromUrl(url string) []byte {
	if url == "" {
		return nil
	}

	return nil
}
