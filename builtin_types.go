package eosabi

import (
	"bytes"
	"encoding/binary"

	base58 "github.com/Jeiwan/eos-b58"
	log "github.com/sirupsen/logrus"
)

var builtinTypes = map[string]interface{}{
	// "bool":uint8,
	// "int8":int8,
	// "uint8":uint8,
	// "int16": true,
	"uint16": true,
	// "int32":int32,
	"uint32": true,
	// "int64":int64,
	// "uint64":uint64,
	// "int128":int128,
	// "uint128":uint128,
	// "varint32":fc::signed_int,
	// "varuint32":fc::unsigned_int,

	// port for floating point types. For now this is good enough.
	// "float32":float,
	// "float64":double,
	// "float128":uint128,

	// "time_point":fc::time_point,
	// "time_point_sec":fc::time_point_sec,
	// "block_timestamp_type":block_timestamp_type,

	"name": true,

	// "bytes":bytes,
	// "string":string,

	// "checksum160":checksum160_type,
	// "checksum256":checksum256_type,
	// "checksum512":checksum512_type,

	"public_key": true,
	// "signature":signature_type,

	// "symbol":symbol,
	// "symbol_code":symbol_code,
	// "asset":asset,
	// "extended_asset":extended_asset,
}

func unpackBuiltin(t string, stream *bytes.Buffer) interface{} {
	switch t {
	case "name":
		return unpackName(stream)
	case "public_key":
		return unpackPublicKey(stream)
	case "uint16":
		return unpackUint16(stream)
	case "uint32":
		return unpackUint32(stream)
	}

	return nil
}

func unpackName(stream *bytes.Buffer) interface{} {
	size := 8
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return nameToString(data)
}

func unpackUint32(stream *bytes.Buffer) interface{} {
	size := 4
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return binary.LittleEndian.Uint32(data)
}

func unpackUint16(stream *bytes.Buffer) interface{} {
	size := 2
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return binary.LittleEndian.Uint16(data)
}

func unpackPublicKey(stream *bytes.Buffer) interface{} {
	size := 34
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	b58 := base58.CheckEncode(data[1:])

	return "EOS" + string(b58)
}
