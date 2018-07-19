package eosabi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"

	base58 "github.com/Jeiwan/eos-b58"
	log "github.com/sirupsen/logrus"
)

const epochMs = 946684800000
const blockIntervalMs = 500

var builtinTypes = map[string]interface{}{
	"bool":      true,
	"int8":      true,
	"uint8":     true,
	"int16":     true,
	"uint16":    true,
	"int32":     true,
	"uint32":    true,
	"int64":     true,
	"uint64":    true,
	"int128":    true,
	"uint128":   true,
	"varint32":  true,
	"varuint32": true,

	"float32":  true,
	"float64":  true,
	"float128": true,

	"time_point":           true,
	"time_point_sec":       true,
	"block_timestamp_type": true,

	"name": true,

	"bytes":  true,
	"string": true,

	"checksum160": true,
	"checksum256": true,
	"checksum512": true,

	"public_key": true,
	"signature":  true,

	"symbol":         true,
	"symbol_code":    true,
	"asset":          true,
	"extended_asset": true,
}

func unpackBuiltin(t string, stream *bytes.Buffer) (interface{}, error) {
	switch t {
	case "bool":
		return unpackBool(stream)
	case "int8":
		return unpackInt8(stream)
	case "uint8":
		return unpackUint8(stream)
	case "int16":
		return unpackInt16(stream)
	case "uint16":
		return unpackUint16(stream)
	case "int32":
		return unpackInt32(stream)
	case "uint32":
		return unpackUint32(stream)
	case "int64":
		return unpackInt64(stream)
	case "uint64":
		return unpackUint64(stream)
	case "int128":
		return unpackInt128(stream)
	case "uint128":
		return unpackUint128(stream)
	case "varint32":
		return unpackVarInt32(stream)
	case "varuint32":
		return unpackVarUint32(stream)

	case "float32":
		return unpackFloat32(stream)
	case "float64":
		return unpackFloat64(stream)
	case "float128":
		return unpackFloat128(stream)

	case "time_point":
		return unpackTimePoint(stream)
	case "time_point_sec":
		return unpackTimePointSec(stream)

	case "name":
		return unpackName(stream)

	case "bytes":
		return unpackBytes(stream)
	case "string":
		return unpackString(stream)

	case "checksum160":
		return unpackChecksum160(stream)
	case "checksum256":
		return unpackChecksum256(stream)
	case "checksum512":
		return unpackChecksum512(stream)

	case "public_key":
		return unpackPublicKey(stream)
	case "signature":
		return unpackSignature(stream)

	case "symbol":
		return unpackSymbol(stream)
	case "symbol_code":
		return unpackSymbolCode(stream)
	case "asset":
		return unpackAsset(stream)
	case "extended_asset":
		return unpackAsset(stream)
	}

	return nil, fmt.Errorf("builtin type not found: %s", t)
}

func unpackBool(stream *bytes.Buffer) (interface{}, error) {
	return unpackUint8(stream)
}

func unpackInt8(stream *bytes.Buffer) (interface{}, error) {
	b, err := stream.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("unpackInt8: %s", err)
	}

	return int8(b), nil
}

func unpackUint8(stream *bytes.Buffer) (interface{}, error) {
	b, err := stream.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("unpackUint8: %s", err)
	}

	return uint8(b), nil
}

func unpackInt16(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint16(stream)
	if err != nil {
		return nil, err
	}

	return int16(r.(uint16)), nil
}

func unpackUint16(stream *bytes.Buffer) (interface{}, error) {
	size := 2
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackUint16: %s", err)
	}

	return binary.LittleEndian.Uint16(data), nil
}

func unpackInt32(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint32(stream)
	if err != nil {
		return nil, err
	}

	return int32(r.(uint32)), nil
}

func unpackUint32(stream *bytes.Buffer) (interface{}, error) {
	size := 4
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackUint32: %s", err)
	}

	return binary.LittleEndian.Uint32(data), nil
}

func unpackInt64(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint64(stream)
	if err != nil {
		return nil, err
	}

	return int64(r.(uint64)), nil
}

func unpackUint64(stream *bytes.Buffer) (interface{}, error) {
	size := 8
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackUint64: %s", err)
	}

	return binary.LittleEndian.Uint64(data), nil
}

func unpackInt128(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint128(stream)
	if err != nil {
		return nil, err
	}

	return r.(string), nil
}

func unpackUint128(stream *bytes.Buffer) (interface{}, error) {
	size := 16
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackUint128: %s", err)
	}

	bigN := big.NewInt(0)
	bigN.SetBytes(data)

	return bigN.String(), nil
}

func unpackVarInt32(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackVarUint32(stream)
	if err != nil {
		return nil, err
	}

	return int32(r.(uint32)), nil
}

func unpackVarUint32(stream *bytes.Buffer) (interface{}, error) {
	size := 4
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackVarUint32: %s", err)
	}

	v, n := binary.Varint(data)
	if n <= 0 {
		log.Fatalln("varint error")
	}

	return uint32(v), err
}

func unpackFloat32(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint32(stream)
	if err != nil {
		return nil, err
	}

	return math.Float32frombits(r.(uint32)), nil
}

func unpackFloat64(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint64(stream)
	if err != nil {
		return nil, err
	}

	return math.Float64frombits(r.(uint64)), nil
}

func unpackFloat128(stream *bytes.Buffer) (interface{}, error) {
	return unpackUint128(stream)
}

func unpackTimePoint(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackInt64(stream)
	if err != nil {
		return nil, err
	}

	microseconds := r.(int64)
	timestamp := time.Unix(0, microseconds*1000)

	return timestamp.Format(time.RFC3339), nil
}

func unpackTimePointSec(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackInt64(stream)
	if err != nil {
		return nil, err
	}

	seconds := r.(int64)
	timestamp := time.Unix(seconds, 0)

	return timestamp.Format(time.RFC3339), nil
}

func unpackBlockTimestamp(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint32(stream)
	if err != nil {
		return nil, err
	}

	slots := r.(uint32)
	milliseconds := int64(slots*blockIntervalMs) + epochMs
	timestamp := time.Unix(0, milliseconds*1000000)

	return timestamp.Format(time.RFC3339), nil
}

func unpackName(stream *bytes.Buffer) (interface{}, error) {
	size := 8
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackName: %s", err)
	}

	return nameToString(data), nil
}

func unpackBytes(stream *bytes.Buffer) (interface{}, error) {
	var v uint64
	var err error
	var b byte
	var by uint

	for {
		b, err = stream.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("unpackBytes: %s", err)
		}

		v = v | uint64(b&0x7f)<<by
		by += 7

		if !(b&0x80 != 0 && by < 32) {
			break
		}
	}

	data := make([]byte, v)
	_, err = stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackBytes: %s", err)
	}

	return data, nil
}

func unpackString(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackBytes(stream)
	if err != nil {
		return nil, err
	}

	return string(r.([]byte)), nil
}

func unpackChecksum160(stream *bytes.Buffer) (interface{}, error) {
	size := 20
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackChecksum160: %s", err)
	}

	return data, nil
}

func unpackChecksum256(stream *bytes.Buffer) (interface{}, error) {
	size := 32
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackChecksum256: %s", err)
	}

	return data, nil
}

func unpackChecksum512(stream *bytes.Buffer) (interface{}, error) {
	size := 64
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackChecksum512: %s", err)
	}

	return data, nil
}

func unpackPublicKey(stream *bytes.Buffer) (interface{}, error) {
	size := 34
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackPublicKey: %s", err)
	}

	b58 := base58.CheckEncode(data[1:])

	return "EOS" + string(b58), nil
}

func unpackSignature(stream *bytes.Buffer) (interface{}, error) {
	size := 65
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackSignature: %s", err)
	}

	b58 := base58.CheckEncode(data)

	return "SIG_" + string(b58), nil
}

func unpackSymbol(stream *bytes.Buffer) (interface{}, error) {
	size := 8
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackSymbol: %s", err)
	}

	p := uint8(data[0])
	s := string(bytes.Trim(data[1:], "\x00"))

	symbol := symbol{Precision: p, Token: s}

	return symbol.String(), nil
}

func unpackSymbolCode(stream *bytes.Buffer) (interface{}, error) {
	size := 8
	data := make([]byte, size)
	_, err := stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unpackSymbolCode: %s", err)
	}

	s := string(data[1:])

	return s, nil
}

func unpackAsset(stream *bytes.Buffer) (interface{}, error) {
	r, err := unpackUint64(stream)
	if err != nil {
		return nil, err
	}
	a := r.(uint64)

	r, err = unpackSymbol(stream)
	if err != nil {
		return nil, err
	}
	s := r.(string)

	asset := asset{Amount: a, Symbol: parseSymbol(s)}

	return asset.String(), nil
}
