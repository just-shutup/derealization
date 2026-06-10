// internal/protocol/types.go
// Содержит базовые типы данных протокола Minecraft (VarInt, Position, String).

package protocol

import (
	"errors"
	"io"
)

// ErrVarIntTooBig возвращается, если VarInt превышает допустимый размер (5 байт).
var ErrVarIntTooBig = errors.New("VarInt is too big")

// ReadVarInt читает VarInt из потока.
func ReadVarInt(r io.ByteReader) (int32, error) {
	var num uint32
	var shift uint
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		num |= uint32(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 35 {
			return 0, ErrVarIntTooBig
		}
	}
	return int32(num), nil
}

// WriteVarInt записывает VarInt в поток.
func WriteVarInt(w io.ByteWriter, val int32) error {
	uval := uint32(val)
	for {
		b := byte(uval & 0x7F)
		uval >>= 7
		if uval != 0 {
			b |= 0x80
		}
		if err := w.WriteByte(b); err != nil {
			return err
		}
		if uval == 0 {
			break
		}
	}
	return nil
}

// Position представляет координаты блока.
type Position struct {
	X, Y, Z int
}

// Encode упаковывает Position в uint64.
func (p Position) Encode() uint64 {
	return ((uint64(p.X) & 0x3FFFFFF) << 38) | ((uint64(p.Y) & 0xFFF) << 26) | (uint64(p.Z) & 0x3FFFFFF)
}

// Decode распаковывает uint64 в Position.
func DecodePosition(val uint64) Position {
	x := int(val >> 38)
	y := int((val >> 26) & 0xFFF)
	z := int(val & 0x3FFFFFF)

	// Обработка отрицательных чисел (26 бит для X и Z, 12 для Y)
	if x >= 1<<25 {
		x -= 1 << 26
	}
	if y >= 1<<11 {
		y -= 1 << 12
	}
	if z >= 1<<25 {
		z -= 1 << 26
	}
	return Position{X: x, Y: y, Z: z}
}