package abi

// #cgo linux CFLAGS: -I${SRCDIR}/abieos/src
// #cgo darwin CFLAGS: -I${SRCDIR}/abieos/src
// #cgo linux CXXFLAGS: -I${SRCDIR}/abieos/src
// #cgo darwin CXXFLAGS: -I${SRCDIR}/abieos/src
// #cgo linux LDFLAGS: -L${SRCDIR}/abieos/lib/linux
// #cgo darwin LDFLAGS: -L${SRCDIR}/abieos/lib/darwin
// #cgo LDFLAGS: -labieos
// #include "abieos/src/abieos.h"
import "C"

import (
	"errors"
)

// BinToJSON converts binary to JSON (decodes ABI).
func BinToJSON(abiDef []byte, bin []byte, typeName string) (string, error) {
	ctx := C.abieos_create()

	contract := C.abieos_string_to_name(ctx, C.CString(""))

	if C.abieos_set_abi(ctx, contract, C.CString(string(abiDef))) != 1 {
		return "", errors.New(C.GoString(C.abieos_get_error(ctx)))
	}

	t := C.abieos_get_type_for_action(ctx, contract, C.abieos_string_to_name(ctx, C.CString(typeName)))
	if len(C.GoString(t)) == 0 {
		return "", errors.New(C.GoString(C.abieos_get_error(ctx)))
	}

	json := C.abieos_bin_to_json(ctx, contract, t, C.CString(string(bin)), C.ulong(len(bin)))
	if len(C.GoString(json)) == 0 {
		return "", errors.New(C.GoString(C.abieos_get_error(ctx)))
	}

	return C.GoString(json), nil
}

// JSONToBin converts JSON to binary (encodes ABI).
func JSONToBin(abiDef []byte, json []byte, typeName string) (string, error) {
	ctx := C.abieos_create()

	contract := C.abieos_string_to_name(ctx, C.CString(""))

	if C.abieos_set_abi(ctx, contract, C.CString(string(abiDef))) != 1 {
		return "", errors.New(C.GoString(C.abieos_get_error(ctx)))
	}

	t := C.abieos_get_type_for_action(ctx, contract, C.abieos_string_to_name(ctx, C.CString(typeName)))
	if len(C.GoString(t)) == 0 {
		return "", errors.New(C.GoString(C.abieos_get_error(ctx)))
	}

	if C.abieos_json_to_bin_reorderable(ctx, contract, t, C.CString(string(json))) != 1 {
		return "", errors.New(C.GoString(C.abieos_get_error(ctx)))
	}

	return C.GoString(C.abieos_get_bin_hex(ctx)), nil
}
