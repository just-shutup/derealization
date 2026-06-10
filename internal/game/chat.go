// internal/game/chat.go
package game

import (
	"encoding/json"
	"strings"
)

// ChatComponent представляет JSON-структуру чата Minecraft 1.12.2.
type ChatComponent struct {
	Text  string          `json:"text,omitempty"`
	Extra []ChatComponent `json:"extra,omitempty"`
	// Здесь можно добавить поля Color, Bold, Italic для продвинутого парсинга
}

// ParseChat извлекает чистый текст из JSON-объекта сервера.
func ParseChat(rawJSON string) string {
	var comp ChatComponent
	if err := json.Unmarshal([]byte(rawJSON), &comp); err != nil {
		// Некоторые кастомные серверы или плагины могут слать простой текст
		return rawJSON
	}

	var sb strings.Builder
	extractText(&comp, &sb)
	return sb.String()
}

func extractText(comp *ChatComponent, sb *strings.Builder) {
	sb.WriteString(comp.Text)
	for _, extra := range comp.Extra {
		extractText(&extra, sb)
	}
}