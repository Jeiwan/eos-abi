package eosabi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackWithABI(t *testing.T) {
	abi, _ := ioutil.ReadFile("fixtures/eosio.json")
	action := "newaccount"
	data, _ := hex.DecodeString("0000000000ea3055a01861fc499b89690100000001000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f010000000100000001000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f01000000")

	unpacked := unpackWithABI(abi, action, data)
	actual := unpacked.(map[string]interface{})

	assert.Equal(t, "ha4tqmjwg4ge", actual["name"])
	assert.Equal(t, "eosio", actual["creator"])

	owner := actual["owner"].(map[string]interface{})
	assert.Equal(t, uint32(1), owner["threshold"])

	keys := owner["keys"].([]interface{})
	assert.Len(t, keys, 1)

	key := keys[0].(map[string]interface{})
	assert.Equal(t, "000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f", key["key"])
	// assert.Equal(t, "EOS7agSKkiM1bUz4vJZ5DB6eMNZUjridajqggr8hrPXaL63mTuL5E", key["key"])
	assert.Equal(t, uint16(1), key["weight"])

	assert.Len(t, owner["accounts"].([]interface{}), 0)
	assert.Len(t, owner["waits"].([]interface{}), 0)

	active := actual["active"].(map[string]interface{})
	assert.Equal(t, uint32(1), active["threshold"])

	keys = active["keys"].([]interface{})
	assert.Len(t, keys, 1)

	key = keys[0].(map[string]interface{})
	assert.Equal(t, "000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f", key["key"])
	// assert.Equal(t, "EOS7agSKkiM1bUz4vJZ5DB6eMNZUjridajqggr8hrPXaL63mTuL5E", key["key"])
	assert.Equal(t, uint16(1), key["weight"])

	assert.Len(t, active["accounts"].([]interface{}), 0)
	assert.Len(t, active["waits"].([]interface{}), 0)
}

func TestUnpackArray(t *testing.T) {
	abiData, _ := ioutil.ReadFile("fixtures/eosio.json")
	var abi abi
	json.Unmarshal(abiData, &abi)

	tpe := "key_weight"
	data, _ := hex.DecodeString("01000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f0100")
	stream := bytes.NewBuffer(data)

	unpacked := unpackArray(tpe, stream, &abi)

	assert.Len(t, unpacked, 1)
	el := unpacked[0].(map[string]interface{})
	assert.Equal(t,
		"000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f",
		el["key"],
	)
	assert.Equal(t,
		uint16(1),
		el["weight"],
	)
}

func TestFundamentalType(t *testing.T) {
	assert.Equal(t,
		"string",
		fundamentalType("string[]"),
	)
	assert.Equal(t,
		"string",
		fundamentalType("string?"),
	)
	assert.Equal(t,
		"string",
		fundamentalType("string"),
	)
}

func TestIsArray(t *testing.T) {
	assert.True(t, isArray("uint32[]"))
	assert.False(t, isArray("uint32"))
}

func TestIsOptional(t *testing.T) {
	assert.True(t, isOptional("uint32?"))
	assert.False(t, isOptional("uint32"))
}
