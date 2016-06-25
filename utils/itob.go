package utils

import "encoding/binary"

// Itob converts an int into a big endian byte slice
// for sequential ids in the bolt database
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
