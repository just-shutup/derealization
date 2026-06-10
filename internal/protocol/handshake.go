// internal/protocol/handshake.go
// Обработка Login flow и формирования первых пакетов.

package protocol

// SendHandshake формирует и отправляет пакет Handshake.
func SendHandshake(conn *Conn, serverAddr string, port uint16, nextState int32) error {
	w := NewWriter()
	
	// Protocol Version 340 (1.12.2)
	if err := w.WriteVarInt(340); err != nil {
		return err
	}
	if err := w.WriteString(serverAddr); err != nil {
		return err
	}
	if err := w.WriteUint16(port); err != nil {
		return err
	}
	if err := w.WriteVarInt(nextState); err != nil {
		return err
	}

	return conn.WritePacket(&Packet{
		ID:   0x00, // Handshake ID всегда 0x00
		Data: w.Bytes(),
	})
}

// SendLoginStart отправляет пакет старта логина с именем пользователя.
func SendLoginStart(conn *Conn, username string) error {
	w := NewWriter()
	if err := w.WriteString(username); err != nil {
		return err
	}

	return conn.WritePacket(&Packet{
		ID:   CLoginStart,
		Data: w.Bytes(),
	})
}