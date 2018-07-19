package eosabi

import (
	"bytes"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
)

func unpackWithABI(abiBytes []byte, t string, data []byte) interface{} {
	var abi abi
	err := json.Unmarshal(abiBytes, &abi)
	if err != nil {
		log.Fatalln(err)
	}

	stream := bytes.NewBuffer(data)
	return unpack(t, stream, &abi)
}

func unpack(t string, stream *bytes.Buffer, abi *abi) interface{} {
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

func unpackArray(t string, stream *bytes.Buffer, abi *abi) []interface{} {
	var v uint64
	var err error
	var b byte
	var by uint
	result := []interface{}{}

	for {
		b, err = stream.ReadByte()
		if err != nil {
			log.Fatalln(err)
		}

		v = v | uint64(b&0x7f)<<by
		by += 7

		if !(b&0x80 != 0 && by < 32) {
			break
		}
	}

	for i := uint64(0); i < v; i++ {
		element := unpack(t, stream, abi)
		result = append(result, element)
	}

	return result
}

func unpackStruct(t string, stream *bytes.Buffer, abi *abi) interface{} {
	result := make(map[string]interface{})

	var strct abiStruct
	for _, s := range abi.Structs {
		if s.Name == t {
			strct = s
			break
		}
	}
	if strct.Name == "" {
		log.Fatalf("struct not found: %s", t)
	}

	for _, f := range strct.Fields {
		result[f.Name] = unpack(f.Type, stream, abi)
	}

	return result
}

func fundamentalType(t string) string {
	// TODO: implement type_name?
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
