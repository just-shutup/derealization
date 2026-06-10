// internal/bot/manager.go
// Управление множеством ботов, агрегация логов и событий.

package bot

import (
	"errors"
	"sync"

	"github.com/rs/zerolog"
)

var ErrBotNotFound = errors.New("bot not found")
var ErrBotAlreadyExists = errors.New("bot already exists")

type BotManager struct {
	bots   map[string]*Bot
	events chan BotEvent
	mu     sync.RWMutex
	logger zerolog.Logger
}

func NewBotManager(logger zerolog.Logger) *BotManager {
	return &BotManager{
		bots:   make(map[string]*Bot),
		events: make(chan BotEvent, 1000),
		logger: logger,
	}
}

func (m *BotManager) AddBot(config BotConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.bots[config.Name]; exists {
		return ErrBotAlreadyExists
	}

	bot := NewBot(config, m.events, m.logger)
	m.bots[config.Name] = bot
	return nil
}

func (m *BotManager) RemoveBot(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	bot, exists := m.bots[name]
	if !exists {
		return ErrBotNotFound
	}

	bot.Stop()
	delete(m.bots, name)
	return nil
}

func (m *BotManager) StartBot(name string) error {
	m.mu.RLock()
	bot, exists := m.bots[name]
	m.mu.RUnlock()

	if !exists {
		return ErrBotNotFound
	}

	return bot.Start()
}

func (m *BotManager) StopBot(name string) error {
	m.mu.RLock()
	bot, exists := m.bots[name]
	m.mu.RUnlock()

	if !exists {
		return ErrBotNotFound
	}

	bot.Stop()
	return nil
}

func (m *BotManager) StopAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, bot := range m.bots { // Убрали лишнее :=
		bot.Stop()
	}
	return nil
}

// Events возвращает канал событий для чтения из TUI
func (m *BotManager) Events() <-chan BotEvent {
	return m.events
}

type BotStatus struct {
	Name   string
	Status string
	Health string
	Task   string
}

// ListBots собирает статусы всех добавленных ботов для TUI
func (m *BotManager) ListBots() []BotStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var statuses []BotStatus
	for _, b := range m.bots {
		state := "🔴 Offline"
		if b.IsRunning() {
			state = "🟢 Online"
		}

		taskName := "None"
		if b.task != nil {
			taskName = b.task.Name()
		}

		statuses = append(statuses, BotStatus{
			Name:   b.config.Name,
			Status: state,
			Health: "—", // ХП добавим позже, когда расшифруем пакет UpdateHealth
			Task:   taskName,
		})
	}
	return statuses
}