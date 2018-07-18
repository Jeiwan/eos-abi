package eosabi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestUnpackWithABI(t *testing.T) {
// 	abi, _ := ioutil.ReadFile("fixtures/eosio.json")
// 	action := "newaccount"
// 	data, _ := hex.DecodeString("0000000000ea3055a01861fc499b89690100000001000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f010000000100000001000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f01000000")

// 	unpacked := unpackWithABI(abi, action, data)

// 	// assert.NotNil(tt, unpacked)
// }

func TestNameToString(t *testing.T) {
	name, _ := hex.DecodeString("a01861fc499b8969")

	str := nameToString(name)

	assert.Equal(t, "ha4tqmjwg4ge", str)
}
