package eosabi

import (
	"encoding/binary"
	"strings"
)

func nameToString(bytes []byte) string {
	charmap := ".12345abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, 13)

	value := binary.LittleEndian.Uint64(bytes)
	for i := 0; i <= 12; i++ {
		b := uint64(0x1f)
		if i == 0 {
			b = 0x0f
		}

		idx := value & b
		c := charmap[idx]
		result[12-i] = c

		shift := uint(5)
		if i == 0 {
			shift = 4
		}

		value = value >> shift
	}

	return strings.TrimRight(string(result), ".")
}
