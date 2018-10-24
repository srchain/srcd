package vm

import "encoding/binary"

func PushdataBytes(in []byte) []byte {
	l := len(in)
	if l == 0 {
		return []byte{byte(OP_0)}
	}
	if l <= 75 {
		return append([]byte{byte(OP_DATA_1) + uint8(l) - 1}, in...)
	}
	if l < 1<<8 {
		return append([]byte{byte(OP_PUSHDATA1), uint8(l)}, in...)
	}
	if l < 1<<16 {
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], uint16(l))
		return append([]byte{byte(OP_PUSHDATA2), b[0], b[1]}, in...)
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(l))
	return append([]byte{byte(OP_PUSHDATA4), b[0], b[1], b[2], b[3]}, in...)
}

func PushdataInt64(n int64) []byte {
	if n == 0 {
		return []byte{byte(OP_0)}
	}
	if n >= 1 && n <= 16 {
		return []byte{uint8(OP_1) + uint8(n) - 1}
	}
	return PushdataBytes(Int64Bytes(n))
}

func Int64Bytes(n int64) []byte {
	if n == 0 {
		return []byte{}
	}
	res := make([]byte, 8)
	// converting int64 to uint64 is a safe operation that
	// preserves all data
	binary.LittleEndian.PutUint64(res, uint64(n))
	for len(res) > 0 && res[len(res)-1] == 0 {
		res = res[:len(res)-1]
	}
	return res
}