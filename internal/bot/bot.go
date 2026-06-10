// internal/bot/bot.go
package bot

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"
	"bytes"

	"derealization/internal/game"
	"derealization/internal/protocol"
	"derealization/internal/tasks"

	"github.com/rs/zerolog"
)

type Bot struct {
	config     BotConfig
	conn       *protocol.Conn
	state      *game.GameState
	task       tasks.Task
	logger     zerolog.Logger

	inPackets  chan *protocol.Packet
	outPackets chan *protocol.Packet
	events     chan BotEvent

	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewBot(cfg BotConfig, events chan BotEvent, logger zerolog.Logger) *Bot {
	return &Bot{
		config:     cfg,
		state:      game.NewGameState(),
		inPackets:  make(chan *protocol.Packet, 100),
		outPackets: make(chan *protocol.Packet, 100),
		events:     events,
		logger:     logger.With().Str("bot", cfg.Name).Logger(),
	}
}

func (b *Bot) Start() error {
	b.ctx, b.cancel = context.WithCancel(context.Background())
	
	// Подключение к серверу
	netConn, err := net.Dial("tcp", b.config.Server)
	if err != nil {
		return err
	}
	b.conn = protocol.NewConn(netConn)

	// Парсинг порта для хендшейка
	port := uint16(25565)
	// TODO: полноценный парсинг host:port

	// Login Sequence
	if err := protocol.SendHandshake(b.conn, strings.Split(b.config.Server, ":")[0], port, 2); err != nil {
		return err
	}
	if err := protocol.SendLoginStart(b.conn, b.config.Username); err != nil {
		return err
	}

	b.wg.Add(3)
	go b.readLoop()
	go b.writeLoop()
	go b.gameLoop()

	b.logger.Info().Msg("Bot started and connected")
	return nil
}

func (b *Bot) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
	if b.conn != nil {
		_ = b.conn.Close()
	}
	b.wg.Wait()
	b.logger.Info().Msg("Bot stopped")
}

func (b *Bot) readLoop() {
	defer b.wg.Done()
	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			pkt, err := b.conn.ReadPacket()
			if err != nil {
				b.logger.Error().Err(err).Msg("Disconnect during read")
				b.cancel()
				return
			}
			b.inPackets <- pkt
		}
	}
}

func (b *Bot) writeLoop() {
	defer b.wg.Done()
	for {
		select {
		case <-b.ctx.Done():
			return
		case pkt := <-b.outPackets:
			if err := b.conn.WritePacket(pkt); err != nil {
				b.logger.Error().Err(err).Msg("Disconnect during write")
				b.cancel()
				return
			}
		}
	}
}

func (b *Bot) gameLoop() {
	defer b.wg.Done()
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case pkt := <-b.inPackets:
			b.handlePacket(pkt)
		case <-tick.C:
			// Каждый тик (50ms) вызываем логику задания
			if b.task != nil {
				// ПЕРЕДАЕМ b (самого бота) В МЕТОД Tick!
				if err := b.task.Tick(b); err != nil {
					b.logger.Warn().Err(err).Msg("Task tick failed")
				}
			}
		}
	}
}

func (b *Bot) handlePacket(pkt *protocol.Packet) {
	reader := protocol.NewReader(bytes.NewReader(pkt.Data))

	switch pkt.ID {
	case protocol.SKeepAlive:
		// Сервер шлет Keep-Alive (ID: 0x1F), payload: Int64
		keepAliveID, err := reader.ReadInt64()
		if err == nil {
			writer := protocol.NewWriter()
			_ = writer.WriteInt64(keepAliveID)
			b.outPackets <- &protocol.Packet{
				ID:   protocol.CKeepAlive,
				Data: writer.Bytes(),
			}
			// Опционально: можно закомментировать лог, чтобы не спамил каждые 15 сек
			b.logger.Debug().Int64("id", keepAliveID).Msg("Keep-Alive pong sent")
		}

	case protocol.SChatMessage:
		// Чат (ID: 0x0F), payload: JSON String, Byte (Position)
		chatJSON, err := reader.ReadString()
		if err != nil {
			return
		}
		position, _ := reader.ReadByte() // 0: chat, 1: system, 2: hotbar

		// Игнорируем спам над панелью быстрого доступа (hotbar)
		if position != 2 {
			cleanText := game.ParseChat(chatJSON)
			b.logger.Info().Msgf("[CHAT] %s", cleanText)

			// Отправляем событие в TUI
			b.events <- BotEvent{
				BotName:   b.config.Name,
				Type:      LogMessage,
				Data:      cleanText,
				Timestamp: time.Now(),
			}
		}

	case protocol.SPlayerPositionLook:
		// Сервер телепортирует нас или впервые спавнит (ID: 0x2F)
		x, _ := reader.ReadFloat64()
		y, _ := reader.ReadFloat64()
		z, _ := reader.ReadFloat64()
		yaw, _ := reader.ReadFloat32()
		pitch, _ := reader.ReadFloat32()
		_, _ = reader.ReadByte()
		teleportID, _ := reader.ReadVarInt()

		// Обновляем состояние игрока
		b.state.Player.SetPosition(x, y, z)
		b.logger.Info().Msgf("Spawned/Teleported to %.2f, %.2f, %.2f (Yaw: %.1f, Pitch: %.1f)", x, y, z, yaw, pitch)

		// Обязательное подтверждение телепортации (ID: 0x00 CTeleportConfirm)
		writer := protocol.NewWriter()
		_ = writer.WriteVarInt(teleportID)
		b.outPackets <- &protocol.Packet{
			ID:   protocol.CTeleportConfirm,
			Data: writer.Bytes(),
		}
	}
}

// Возвращает true, если бот запущен и подключен
func (b *Bot) IsRunning() bool {
	return b.ctx != nil && b.ctx.Err() == nil
}