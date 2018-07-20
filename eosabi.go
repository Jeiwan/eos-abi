package eosabi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// UnpackAction unpacks an action
func UnpackAction(abiBytes []byte, action string, data []byte) (interface{}, error) {
	var abi abi
	err := json.Unmarshal(abiBytes, &abi)
	if err != nil {
		return nil, fmt.Errorf("unpackAction %s: %s", action, err)
	}

	stream := bytes.NewBuffer(data)
	return unpack(action, stream, &abi)
}

// UnpackABIDef unpacks ABI definition as passed to 'setabi' action
func UnpackABIDef(data []byte) (interface{}, error) {
	return UnpackAction([]byte(abiDef), "abi_def", data)
}

func unpack(t string, stream *bytes.Buffer, abi *abi) (interface{}, error) {
	rType := abi.resolveType(t)
	fType := fundamentalType(rType)
	bType := builtinTypes[fType]

	if bType != nil {
		return unpackBuiltin(rType, stream)
	}

	if isArray(rType) {
		return unpackArray(fType, stream, abi)
	}

	return unpackStruct(fType, stream, abi)
}

func unpackArray(t string, stream *bytes.Buffer, abi *abi) ([]interface{}, error) {
	var v uint64
	var err error
	var b byte
	var by uint
	result := []interface{}{}

	for {
		b, err = stream.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("unpackArray %s: %s", t, err)
		}

		v = v | uint64(b&0x7f)<<by
		by += 7

		if !(b&0x80 != 0 && by < 32) {
			break
		}
	}

	for i := uint64(0); i < v; i++ {
		element, err := unpack(t, stream, abi)
		if err != nil {
			return nil, fmt.Errorf("unpackArray %s: %s", t, err)
		}
		result = append(result, element)
	}

	return result, nil
}

func unpackStruct(t string, stream *bytes.Buffer, abi *abi) (interface{}, error) {
	result := make(map[string]interface{})

	var strct abiStruct
	for _, s := range abi.Structs {
		if s.Name == t {
			strct = s
			break
		}
	}
	if strct.Name == "" {
		return nil, fmt.Errorf("unpackStruct %s: empty name", t)
	}

	for _, f := range strct.Fields {
		r, err := unpack(f.Type, stream, abi)
		if err != nil {
			return nil, fmt.Errorf("unpackStruct %s: %s", t, err)
		}
		result[f.Name] = r
	}

	return result, nil
}

func fundamentalType(t string) string {
	if isArray(t) {
		return t[0 : len(t)-2]
	}

	if isOptional(t) {
		return t[0 : len(t)-1]
	}

	return t
}

func isArray(t string) bool {
	return strings.HasSuffix(t, "[]")
}

func isOptional(t string) bool {
	return strings.HasSuffix(t, "?")
}
