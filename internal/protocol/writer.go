// internal/protocol/writer.go
// Обертка над io.Writer для удобной записи примитивов протокола.

package protocol

import (
	"bytes"
	"encoding/binary"
)

type Writer struct {
	buf *bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{buf: new(bytes.Buffer)}
}

func (w *Writer) WriteByte(b byte) error {
	return w.buf.WriteByte(b)
}

func (w *Writer) WriteVarInt(val int32) error {
	return WriteVarInt(w, val)
}

func (w *Writer) WriteString(s string) error {
	b := []byte(s)
	if err := w.WriteVarInt(int32(len(b))); err != nil {
		return err
	}
	_, err := w.buf.Write(b)
	return err
}

func (w *Writer) WriteUint16(val uint16) error {
	return binary.Write(w.buf, binary.BigEndian, val)
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Writer) WriteInt64(val int64) error {
	return binary.Write(w.buf, binary.BigEndian, val)
}