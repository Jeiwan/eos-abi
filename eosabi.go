package eosabi

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type abi struct {
	Version string      `json:"version"`
	Types   []abiType   `json:"types"`
	Structs []abiStruct `json:"structs"`
	Actions []abiAction `json:"actions"`
	// Tables
	// RicardianClauses
	// ErrorMessages
	// AbiExtensions
}

type abiType struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type abiStruct struct {
	Name   string           `json:"name"`
	Base   string           `json:"base"`
	Fields []abiStructField `json:"fields"`
}

type abiStructField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type abiAction struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	RicardianContract string `json:"ricardian_contract"`
}

func unpackWithABI(abiBytes []byte, action string, binary []byte) map[string]interface{} {
	var result map[string]interface{}

	var abi abi
	err := json.Unmarshal(abiBytes, &abi)
	if err != nil {
		log.Fatalln(err)
	}

	var strct abiStruct
	for _, s := range abi.Structs {
		if s.Name == action {
			strct = s
			break
		}
	}
	if strct.Name == "" {
		log.Fatalf("stuct not found: %s", action)
	}

	for _, f := range strct.Fields {
		t := resolveType(&abi, f.Type)
		if len(t) == 0 {
			log.Fatalf("unrecognized type: %s", t)
		}
	}

	return result
}

func resolveType(abi *abi, t string) string {
	var typeName string

	for _, tt := range abi.Types {
		if tt.NewTypeName == t {
			typeName = tt.Type
			break
		}
	}

	return typeName
}
