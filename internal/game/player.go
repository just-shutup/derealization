// internal/game/player.go
// Хранение и потокобезопасное обновление состояния игрока.

package game

import "sync"

// GameState объединяет глобальное состояние мира и игрока.
type GameState struct {
	Player *Player
	// TODO: World, Entities, Inventory
}

func NewGameState() *GameState {
	return &GameState{
		Player: &Player{Health: 20.0},
	}
}

// Player хранит координаты и статус персонажа.
type Player struct {
	mu sync.RWMutex

	EntityID int32
	X, Y, Z  float64
	Yaw      float32
	Pitch    float32
	OnGround bool

	Health float32
	Food   int32
}

func (p *Player) SetPosition(x, y, z float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.X = x
	p.Y = y
	p.Z = z
}

func (p *Player) GetPosition() (float64, float64, float64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.X, p.Y, p.Z
}

func (p *Player) SetHealth(health float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Health = health
}