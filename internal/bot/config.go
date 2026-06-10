// internal/bot/config.go
package bot

type BotConfig struct {
	Name           string `yaml:"name"`
	Username       string `yaml:"username"`
	Server         string `yaml:"server"`
	AuthMode       string `yaml:"auth_mode"`
	AutoReconnect  bool   `yaml:"auto_reconnect"`
	ReconnectDelay int    `yaml:"reconnect_delay"`
}