package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeMap(t *testing.T) {
	testMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "",
	}

	bytes := encodeMap(testMap)
	m := decodeMap(bytes)

	assert.Equal(t, m["key1"], "value1")
	assert.Equal(t, m["key2"], "value2")
	assert.Equal(t, m["key3"], "")
}