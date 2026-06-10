// internal/protocol/reader.go
// Обертка над io.Reader для удобного чтения примитивов протокола.

package protocol

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// ReadByte для поддержки io.ByteReader.
func (r *Reader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r.r, b[:])
	return b[0], err
}

func (r *Reader) ReadVarInt() (int32, error) {
	return ReadVarInt(r)
}

func (r *Reader) ReadString() (string, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (r *Reader) ReadUint16() (uint16, error) {
	var val uint16
	err := binary.Read(r.r, binary.BigEndian, &val)
	return val, err
}

func (r *Reader) ReadInt64() (int64, error) {
	var val int64
	err := binary.Read(r.r, binary.BigEndian, &val)
	return val, err
}

func (r *Reader) ReadFloat32() (float32, error) {
	var val float32
	err := binary.Read(r.r, binary.BigEndian, &val)
	return val, err
}

func (r *Reader) ReadFloat64() (float64, error) {
	var val float64
	err := binary.Read(r.r, binary.BigEndian, &val)
	return val, err
}