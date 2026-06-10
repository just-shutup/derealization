// internal/tasks/task.go
package tasks


type TaskStatus string

const (
	StatusIdle      TaskStatus = "Idle"
	StatusRunning   TaskStatus = "Running"
	StatusCompleted TaskStatus = "Completed"
	StatusFailed    TaskStatus = "Failed"
)

// BotContext описывает, что задание (Task) может делать с ботом.
// Бот из пакета bot будет неявно реализовывать этот интерфейс.
type BotContext interface {
	// В будущем добавим сюда методы, нужные заданиям:
	// SendPacket(pkt *protocol.Packet) error
	// GetPosition() (float64, float64, float64)
}

// Task определяет контракт для любого поведения бота.
type Task interface {
	Name() string
	Description() string
	Status() TaskStatus
	Tick(b BotContext) error // Принимает абстрактный интерфейс, а не конкретную структуру
}