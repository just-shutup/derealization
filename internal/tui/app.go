// internal/tui/app.go
package tui

import (
	"fmt"
	"time"

	"derealization/internal/bot"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

// tickMsg используется для периодического обновления (перерисовки) таблицы
type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type AppModel struct {
	botManager *bot.BotManager
	table      table.Model
	width      int
	height     int
}

func NewAppModel(manager *bot.BotManager) AppModel {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 15},
		{Title: "Status", Width: 15},
		{Title: "Health", Width: 8},
		{Title: "Task", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	model := AppModel{
		botManager: manager,
		table:      t,
	}
	model.updateTable() // Первичное заполнение данными из конфига
	return model
}

// updateTable запрашивает свежие данные у менеджера и обновляет строки
func (m *AppModel) updateTable() {
	bots := m.botManager.ListBots()
	var rows []table.Row

	for i, b := range bots {
		id := fmt.Sprintf("%d", i+1)
		rows = append(rows, table.Row{id, b.Name, b.Status, b.Health, b.Task})
	}
	m.table.SetRows(rows)
}

func (m AppModel) Init() tea.Cmd {
	return doTick() // При старте TUI запускаем таймер обновления
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tickMsg:
		m.updateTable()
		return m, doTick() // Зацикливаем таймер

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
			
		case "s": // Запустить бота
			row := m.table.SelectedRow()
			if len(row) > 1 {
				botName := row[1]
				// Запускаем в горутине, чтобы блокировка TCP не заморозила интерфейс
				go func() { _ = m.botManager.StartBot(botName) }()
			}
			
		case "x": // Остановить бота
			row := m.table.SelectedRow()
			if len(row) > 1 {
				botName := row[1]
				go func() { _ = m.botManager.StopBot(botName) }()
			}
		}
		
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m AppModel) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render("🤖 Derealization Manager v1.0.0")

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: Выбор | [S]tart | [X]Stop | [Q]uit")

	return baseStyle.Render(
		fmt.Sprintf("%s\n\n%s\n\n%s", header, m.table.View(), footer),
	)
}