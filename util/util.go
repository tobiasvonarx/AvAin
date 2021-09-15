package util

import (
	"bytes"
	"encoding/binary"
)

// converts a int64 into a slice of bytes
func ToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, num); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
