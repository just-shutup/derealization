// cmd/derealization/main.go
package main

import (
	"fmt"
	"os"

	"derealization/internal/bot"
	"derealization/internal/config"
	"derealization/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
)

func main() {
	logFile, err := os.OpenFile("derealization.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Не удалось открыть файл логов: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logFile.Close() }()

	logger := zerolog.New(logFile).With().Timestamp().Logger()
	logger.Info().Msg("Запуск Derealization Manager...")

	cfgPath := "config.yaml"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		createDefaultConfig(cfgPath)
	}

	appCfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
	}

	manager := bot.NewBotManager(logger)

	for _, botCfg := range appCfg.Bots {
		if err := manager.AddBot(botCfg); err != nil {
			logger.Error().Err(err).Str("bot", botCfg.Name).Msg("Не удалось добавить бота")
		}
	}

	appModel := tui.NewAppModel(manager)
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		logger.Fatal().Err(err).Msg("Ошибка запуска TUI")
	}

	logger.Info().Msg("Остановка всех ботов...")
	_ = manager.StopAll()
	logger.Info().Msg("Программа завершена.")
}

func createDefaultConfig(path string) {
	defaultYaml := `
app:
  log_level: "info"
  log_file: "derealization.log"
  max_bots: 10
bots:
  - name: "MinerBot"
    username: "Bot_Alex"
    server: "localhost:25565"
    auth_mode: "offline"
`
	_ = os.WriteFile(path, []byte(defaultYaml), 0644)
}