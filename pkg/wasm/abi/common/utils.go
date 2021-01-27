package common

import "encoding/binary"

// By default, we use LittleEndian to decode/encode data.
var byteOrder = binary.LittleEndian

func EncodeMap(pairs map[string]string) []byte {
	// calculate buf len
	size :=  4
	for key, value := range pairs {
		size += 4 + 4 + len(key) + 1 + len(value) + 1
	}

	// TODO: check size
	buf := make([]byte, size)

	// write pairs size
	byteOrder.PutUint32(buf, uint32(len(pairs)))

	lenBase := 4
	dataBase := 4 + len(pairs) * (4 + 4)

	for key, value := range pairs {
		// write key len
		byteOrder.PutUint32(buf[lenBase:], uint32(len(key)))
		lenBase += 4
		// write value len
		byteOrder.PutUint32(buf[lenBase:], uint32(len(value)))
		lenBase += 4

		// write key data
		copy(buf[dataBase:], key)
		dataBase += len(key)
		dataBase++ // for nil byte

		// write value data
		copy(buf[dataBase:], value)
		dataBase += len(value)
		dataBase++ // for nil byte
	}

	return buf
}

func DecodeMap(buf []byte) map[string]string {
	if len(buf) < 4 {
		return nil
	}

	res := make(map[string]string)

	// TODO check pairSize
	pairSize := byteOrder.Uint32(buf)

	lenBase := 4
	dataBase := 4 + pairSize * (4 + 4)

	for i := 0; i < int(pairSize); i++ {
		// decode key len
		keyLen := byteOrder.Uint32(buf[lenBase:])
		lenBase += 4

		// decode value len
		ValueLen := byteOrder.Uint32(buf[lenBase:])
		lenBase += 4

		// decode key
		key := string(buf[dataBase:dataBase+keyLen])
		dataBase += keyLen
		dataBase++ // for nil byte

		// decode value
		value := string(buf[dataBase:dataBase+ValueLen])
		dataBase += ValueLen
		dataBase++ // for nil byte

		res[key] = value
	}

	return res
}