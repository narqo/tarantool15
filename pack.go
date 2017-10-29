package tarantool15

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

func PackRawBytes(data []byte, v []byte) []byte {
	if v == nil {
		return data
	}
	return append(data, v...)
}

func PackString(data []byte, v string) []byte {
	if v == "" {
		return data
	}
	return append(data, v...)
}

func PackUint(data []byte, v uint) []byte {
	return PackUint32(data, uint32(v))
}

func PackUint8(data []byte, v uint8) []byte {
	return append(data, byte(v))
}

func PackUint16(data []byte, v uint16) []byte {
	return append(data,
		byte(v),
		byte(v>>8),
		0,
		0,
	)
}

func PackUint32(data []byte, v uint32) []byte {
	return append(data,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

func PackUint64(data []byte, v uint64) []byte {
	return append(data,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

func PackBool(data []byte, v bool) []byte {
	panic("not implemented")
}

func Pack(data []byte, args ...interface{}) (ret []byte, err error) {
	ret = data
	for _, v := range args {
		ret, err = packAny(ret, v)
		if err != nil {
			return
		}
	}
	return ret, err
}

func packAny(data []byte, v interface{}) ([]byte, error) {
	if v == nil {
		return data, nil
	}
	switch vt := v.(type) {
	case int:
		data = PackUint32(data, uint32(vt))
	case int8:
		data = PackUint8(data, uint8(vt))
	case int16:
		data = PackUint32(data, uint32(vt))
	case int32:
		data = PackUint32(data, uint32(vt))
	case int64:
		data = PackUint64(data, uint64(vt))
	case uint:
		data = PackUint(data, vt)
	case uint8:
		data = PackUint8(data, vt)
	case uint16:
		data = PackUint16(data, vt)
	case uint32:
		data = PackUint32(data, vt)
	case uint64:
		data = PackUint64(data, vt)
	case string:
		data = PackString(data, vt)
	case []byte:
		data = PackRawBytes(data, vt)
	case bool:
		data = PackBool(data, vt)
	default:
		return nil, fmt.Errorf("unsupported type %T", vt)
	}
	return data, nil
}

func UnpackUint(data []byte) (v uint64, err error) {
	switch len(data) {
	case 8:
		v = binary.LittleEndian.Uint64(data)
	case 4:
		v = uint64(binary.LittleEndian.Uint32(data))
	case 2:
		v = uint64(binary.LittleEndian.Uint16(data))
	case 0:
		break
	default:
		err = fmt.Errorf("bad data size %d", len(data))
	}
	return v, err
}

func UnpackString(data []byte) (string, error) {
	return bytesToStr(data), nil
}

func bytesToStr(data []byte) string {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	shdr := reflect.StringHeader{h.Data, h.Len}
	return *(*string)(unsafe.Pointer(&shdr))
}

func PackTupleFieldUint32(data []byte, v uint32) []byte {
	data = packUint64BER(data, uint64(4))
	data = PackUint32(data, v)
	return data
}

func PackTupleFieldUint64(data []byte, v uint64) []byte {
	data = packUint64BER(data, uint64(8))
	data = PackUint64(data, v)
	return data
}

func PackTupleFieldString(data []byte, v string) []byte {
	data = packUint64BER(data, uint64(len(v)))
	data = PackString(data, v)
	return data
}

func PackTupleFieldBytes(data []byte, v []byte) []byte {
	data = packUint64BER(data, uint64(len(v)))
	data = PackRawBytes(data, v)
	return data
}

// See https://github.com/tarantool/tarantool/blob/stable/doc/box-protocol.txt
// See https://en.wikipedia.org/wiki/Variable-length_quantity
func packUint64BER(data []byte, v uint64) []byte {
	if v == 0 {
		return append(data, 0)
	}

	buf := make([]byte, binary.MaxVarintLen64)
	maxLen := len(buf)
	n := maxLen - 1

	for ; n >= 0 && v > 0; n-- {
		buf[n] = byte(v & 0x7f)
		v >>= 7

		if n != (maxLen - 1) {
			buf[n] |= 0x80
		}
	}
	return append(data, buf[n+1:]...)
}
