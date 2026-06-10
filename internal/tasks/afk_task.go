// internal/tasks/afk_task.go
package tasks

type AFKTask struct {
	status TaskStatus
}

func NewAFKTask() *AFKTask {
	return &AFKTask{
		status: StatusRunning,
	}
}

func (t *AFKTask) Name() string {
	return "AFK"
}

func (t *AFKTask) Description() string {
	return "Бот стоит на месте"
}

func (t *AFKTask) Status() TaskStatus {
	return t.status
}

// Tick принимает BotContext
func (t *AFKTask) Tick(b BotContext) error {
	// AFK ничего не делает
	return nil
}