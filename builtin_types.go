package eosabi

import (
	"bytes"
	"encoding/binary"
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

	// port for floating point types. For now this is good enough.
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

func unpackBuiltin(t string, stream *bytes.Buffer) interface{} {
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

	return nil
}

func unpackBool(stream *bytes.Buffer) interface{} {
	return unpackUint8(stream)
}

func unpackInt8(stream *bytes.Buffer) interface{} {
	b, err := stream.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}

	return int8(b)
}

func unpackUint8(stream *bytes.Buffer) interface{} {
	b, err := stream.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}

	return uint8(b)
}

func unpackInt16(stream *bytes.Buffer) interface{} {
	return int16(unpackUint16(stream).(uint16))
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

func unpackInt32(stream *bytes.Buffer) interface{} {
	return int32(unpackUint32(stream).(uint32))
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

func unpackInt64(stream *bytes.Buffer) interface{} {
	return int64(unpackUint64(stream).(uint64))
}

func unpackUint64(stream *bytes.Buffer) interface{} {
	size := 8
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return binary.LittleEndian.Uint64(data)
}

func unpackInt128(stream *bytes.Buffer) interface{} {
	return unpackUint128(stream).(string)
}

func unpackUint128(stream *bytes.Buffer) interface{} {
	size := 16
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	bigN := big.NewInt(0)
	bigN.SetBytes(data)

	return bigN.String()
}

func unpackVarInt32(stream *bytes.Buffer) interface{} {
	return int32(unpackVarUint32(stream).(uint32))
}

func unpackVarUint32(stream *bytes.Buffer) interface{} {
	size := 4
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	v, n := binary.Varint(data)
	if n <= 0 {
		log.Fatalln("varint error")
	}

	return uint32(v)
}

func unpackFloat32(stream *bytes.Buffer) interface{} {
	v := unpackUint32(stream).(uint32)
	return math.Float32frombits(v)
}

func unpackFloat64(stream *bytes.Buffer) interface{} {
	v := unpackUint64(stream).(uint64)
	return math.Float64frombits(v)
}

func unpackFloat128(stream *bytes.Buffer) interface{} {
	return unpackUint128(stream)
}

func unpackTimePoint(stream *bytes.Buffer) interface{} {
	microseconds := unpackInt64(stream).(int64)
	timestamp := time.Unix(0, microseconds*1000)

	return timestamp.Format(time.RFC3339)
}

func unpackTimePointSec(stream *bytes.Buffer) interface{} {
	seconds := unpackInt64(stream).(int64)
	timestamp := time.Unix(seconds, 0)

	return timestamp.Format(time.RFC3339)
}

func unpackBlockTimestamp(stream *bytes.Buffer) interface{} {
	slots := unpackUint32(stream).(uint32)
	milliseconds := int64(slots*blockIntervalMs) + epochMs
	timestamp := time.Unix(0, milliseconds*1000000)

	return timestamp.Format(time.RFC3339)
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

func unpackBytes(stream *bytes.Buffer) interface{} {
	var v uint64
	var err error
	var b byte
	var by uint

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

	data := make([]byte, v)
	_, err = stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}

	return data
}

func unpackString(stream *bytes.Buffer) interface{} {
	data := unpackBytes(stream).([]byte)
	return string(data)
}

func unpackChecksum160(stream *bytes.Buffer) interface{} {
	size := 20
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return data
}

func unpackChecksum256(stream *bytes.Buffer) interface{} {
	size := 32
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return data
}

func unpackChecksum512(stream *bytes.Buffer) interface{} {
	size := 64
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	return data
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

func unpackSignature(stream *bytes.Buffer) interface{} {
	size := 65
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	b58 := base58.CheckEncode(data)

	return "SIG_" + string(b58)
}

func unpackSymbol(stream *bytes.Buffer) interface{} {
	size := 8
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}
	p := uint8(data[0])
	s := string(bytes.Trim(data[1:], "\x00"))

	symbol := symbol{Precision: p, Token: s}

	return symbol.String()
}

func unpackSymbolCode(stream *bytes.Buffer) interface{} {
	size := 8
	data := make([]byte, size)
	n, err := stream.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	if n != size {
		log.Fatalln(err)
	}

	s := string(data[1:])

	return s
}

func unpackAsset(stream *bytes.Buffer) interface{} {
	a := unpackUint64(stream).(uint64)
	s := unpackSymbol(stream).(string)

	asset := asset{Amount: a, Symbol: parseSymbol(s)}

	return asset.String()
}
