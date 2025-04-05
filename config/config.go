// Пакет config используется для формирования конфигурации приложения

package config

import "github.com/caarlos0/env/v9"

type Config struct {
	ServerPort      uint     `env:"PORT" envDefault:"7540"`
	WebDir          string   `env:"WEB_DIR" envDefault:"./web"`
	IndexFileName   string   `env:"INDEX_FILE_NAME" envDefault:"index.html"`
	AllowedExtFiles []string `env:"ALLOWED_EXT_FILES" envDefault:".html,.js,.css,.ico"`
	DBName          string   `env:"DBFILE" envDefault:"scheduler.db"`
	UserPassword    string   `env:"PASSWORD" envDefault:""`
	JWTSecret       string   `env:"JWT_SECRET" envDefault:"jwt_secret_string"`
}

// Load загружает данные из переменных окружения в структуру Config
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
