package eosabi

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

func (a abi) resolveType(t string) string {
	// TODO: recursive
	for _, tt := range a.Types {
		if tt.NewTypeName == t {
			return tt.Type
		}
	}

	for _, tt := range a.Actions {
		if tt.Name == t {
			return tt.Type
		}
	}

	return t
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
