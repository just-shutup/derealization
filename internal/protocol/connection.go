// internal/protocol/connection.go
// Управляет TCP соединением, фреймингом пакетов (Length + PacketID + Data).

package protocol

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type Packet struct {
	ID   int32
	Data []byte
}

type Conn struct {
	netConn net.Conn
	// TODO: поля для Compression и Encryption (zlib / aes)
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{netConn: conn}
}

// ReadPacket читает пакет с учетом его длины.
func (c *Conn) ReadPacket() (*Packet, error) {
	reader := NewReader(c.netConn)
	
	// Читаем длину пакета
	length, err := reader.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read packet length: %w", err)
	}

	// Читаем тело пакета (ID + Data)
	payload := make([]byte, length)
	if _, err := io.ReadFull(c.netConn, payload); err != nil {
		return nil, fmt.Errorf("failed to read packet payload: %w", err)
	}

	payloadReader := NewReader(bytes.NewReader(payload))
	
	// Читаем Packet ID
	packetID, err := payloadReader.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %w", err)
	}

	// Оставшиеся байты - это Data
	// Вычисляем размер: length - (размер_закодированного_PacketID)
	// Для простоты читаем весь буфер до конца
	data, err := io.ReadAll(payloadReader.r)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet data: %w", err)
	}

	return &Packet{
		ID:   packetID,
		Data: data,
	}, nil
}

// WritePacket оборачивает Packet ID и Data длиной, затем отправляет в сокет.
func (c *Conn) WritePacket(packet *Packet) error {
	payloadWriter := NewWriter()
	
	// Пишем PacketID
	if err := payloadWriter.WriteVarInt(packet.ID); err != nil {
		return err
	}
	
	// Пишем данные
	if _, err := payloadWriter.buf.Write(packet.Data); err != nil {
		return err
	}

	payload := payloadWriter.Bytes()

	// Формируем финальный пакет: Длина + payload
	finalWriter := NewWriter()
	if err := finalWriter.WriteVarInt(int32(len(payload))); err != nil {
		return err
	}
	if _, err := finalWriter.buf.Write(payload); err != nil {
		return err
	}

	_, err := c.netConn.Write(finalWriter.Bytes())
	return err
}

func (c *Conn) Close() error {
	return c.netConn.Close()
}