package eosabi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackAction(t *testing.T) {
	abi, _ := ioutil.ReadFile("fixtures/eosio.json")
	action := "newaccount"
	data, _ := hex.DecodeString("0000000000ea3055a01861fc499b89690100000001000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f010000000100000001000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f01000000")

	unpacked, err := UnpackAction(abi, action, data)
	assert.Nil(t, err)
	actual := unpacked.(map[string]interface{})

	assert.Equal(t, "ha4tqmjwg4ge", actual["name"])
	assert.Equal(t, "eosio", actual["creator"])

	owner := actual["owner"].(map[string]interface{})
	assert.Equal(t, uint32(1), owner["threshold"])

	keys := owner["keys"].([]interface{})
	assert.Len(t, keys, 1)

	key := keys[0].(map[string]interface{})
	assert.Equal(t, "EOS7agSKkiM1bUz4vJZ5DB6eMNZUjridajqggr8hrPXaL63mTuL5E", key["key"])
	assert.Equal(t, uint16(1), key["weight"])

	assert.Len(t, owner["accounts"].([]interface{}), 0)
	assert.Len(t, owner["waits"].([]interface{}), 0)

	active := actual["active"].(map[string]interface{})
	assert.Equal(t, uint32(1), active["threshold"])

	keys = active["keys"].([]interface{})
	assert.Len(t, keys, 1)

	key = keys[0].(map[string]interface{})
	assert.Equal(t, "EOS7agSKkiM1bUz4vJZ5DB6eMNZUjridajqggr8hrPXaL63mTuL5E", key["key"])
	assert.Equal(t, uint16(1), key["weight"])

	assert.Len(t, active["accounts"].([]interface{}), 0)
	assert.Len(t, active["waits"].([]interface{}), 0)
}

func TestUnpackABIDef(t *testing.T) {
	expectedRaw, _ := ioutil.ReadFile("fixtures/abidef.json")
	expected := string(expectedRaw)

	abiRaw, _ := hex.DecodeString("0e656f73696f3a3a6162692f312e300110657468657265756d5f6164647265737306737472696e6702076164647265737300030269640675696e74363410657468657265756d5f6164647265737310657468657265756d5f616464726573730762616c616e636505617373657403616464000210657468657265756d5f6164647265737310657468657265756d5f616464726573730762616c616e63650561737365740100000000000052320361646400010000c00a637553320369363401026964010675696e7436340761646472657373000000")

	abiMap, err := UnpackABIDef(abiRaw)
	assert.Nil(t, err)

	abiBytes, err := json.Marshal(abiMap)
	assert.Nil(t, err)

	actual := string(abiBytes)

	assert.JSONEq(t, expected, actual)
}

func TestUnpackArray(t *testing.T) {
	abiData, _ := ioutil.ReadFile("fixtures/eosio.json")
	var abi abi
	json.Unmarshal(abiData, &abi)

	tpe := "key_weight"
	data, _ := hex.DecodeString("01000362a6a7e46c62856973506a0c9cd9311b7829c563a1f39f8ebcb0d1618e527b0f0100")
	stream := bytes.NewBuffer(data)

	unpacked, err := unpackArray(tpe, stream, &abi)
	assert.Nil(t, err)

	assert.Len(t, unpacked, 1)
	el := unpacked[0].(map[string]interface{})
	assert.Equal(t,
		"EOS7agSKkiM1bUz4vJZ5DB6eMNZUjridajqggr8hrPXaL63mTuL5E",
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
