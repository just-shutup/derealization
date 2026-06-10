// internal/bot/events.go
package bot

import "time"

type EventType string

const (
	StatusChanged EventType = "StatusChanged"
	Disconnected  EventType = "Disconnected"
	LogMessage    EventType = "LogMessage"
)

type BotEvent struct {
	BotName   string
	Type      EventType
	Data      interface{}
	Timestamp time.Time
}