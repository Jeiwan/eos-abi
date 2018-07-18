package eosabi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameToString(t *testing.T) {
	name, _ := hex.DecodeString("a01861fc499b8969")
	str := nameToString(name)
	assert.Equal(t, "ha4tqmjwg4ge", str)

	name, _ = hex.DecodeString("0000000000ea3055")
	str = nameToString(name)
	assert.Equal(t, "eosio", str)
}
