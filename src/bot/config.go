package bot

import (
	"os"
)

type Config struct {
	Token string
}

func ReadConfig() (*Config, error) {
	var cfg Config
	cfg.Token = os.Getenv("BOT_TOKEN")
	return &cfg, nil
}
