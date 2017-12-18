package bele

import (
	"encoding/binary"
	"io"
	"math"
)

// ----- 反序列化 -----

func BEUint16(p []byte) uint16 {
	return binary.BigEndian.Uint16(p)
}

func BEUint24(p []byte) uint32 {
	return uint32(p[2]) | uint32(p[1])<<8 | uint32(p[0])<<16
}

func BEUint32(p []byte) (ret uint32) {
	return binary.BigEndian.Uint32(p)
}

func BEUint64(p []byte) (ret uint64) {
	return binary.BigEndian.Uint64(p)
}

func BEFloat64(p []byte) (ret float64) {
	a := binary.BigEndian.Uint64(p)
	return math.Float64frombits(a)
}

func LEUint32(p []byte) (ret uint32) {
	return binary.LittleEndian.Uint32(p)
}

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	// 原生Read函数，读不够时，会在第一次调用时读入剩余的数据，并返回读入数据的真实长度，以及nil值的error
	// 在下一次Read时，才返回EOF
	// 这里我们在第一次读不够时，就直接返回EOF。（但是也会把剩余的数据读取到）
	nn, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if nn != n {
		return b, io.EOF
	}
	return b, nil
}

func ReadString(r io.Reader, n int) (string, error) {
	b, err := ReadBytes(r, n)
	return string(b), err
}

func ReadUint8(r io.Reader) (uint8, error) {
	b, err := ReadBytes(r, 1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func ReadBEUint16(r io.Reader) (uint16, error) {
	b, err := ReadBytes(r, 2)
	if err != nil {
		return 0, err
	}
	return BEUint16(b), nil
}

func ReadBEUint24(r io.Reader) (uint32, error) {
	b, err := ReadBytes(r, 3)
	if err != nil {
		return 0, err
	}
	return BEUint24(b), nil
}

func ReadBEUint32(r io.Reader) (uint32, error) {
	b, err := ReadBytes(r, 4)
	if err != nil {
		return 0, err
	}
	return BEUint32(b), nil
}

func ReadBEUint64(r io.Reader) (uint64, error) {
	b, err := ReadBytes(r, 8)
	if err != nil {
		return 0, err
	}
	return BEUint64(b), nil
}

func ReadLEUint32(r io.Reader) (uint32, error) {
	b, err := ReadBytes(r, 4)
	if err != nil {
		return 0, err
	}
	return LEUint32(b), nil
}

// ----- 序列化 -----

func BEPutUint24(out []byte, in uint32) {
	out[0] = byte(in >> 16)
	out[1] = byte(in >> 8)
	out[2] = byte(in)
}

func BEPutUint32(out []byte, in uint32) {
	binary.BigEndian.PutUint32(out, in)
}

func LEPutUint32(out []byte, in uint32) {
	binary.LittleEndian.PutUint32(out, in)
}

func WriteBEUint24(writer io.Writer, in uint32) error {
	_, err := writer.Write([]byte{uint8(in >> 16), uint8(in >> 8), uint8(in & 0xFF)})
	return err
}

func WriteBE(writer io.Writer, in interface{}) error {
	return binary.Write(writer, binary.BigEndian, in)
}

func WriteLE(writer io.Writer, in interface{}) error {
	return binary.Write(writer, binary.LittleEndian, in)
}
