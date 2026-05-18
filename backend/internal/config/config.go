package config

import "os"

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	FrontendURL string
}

func Load() Config {
	return Config{
		AppEnv:      env("APP_ENV", "development"),
		Port:        env("BACKEND_PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://team_task_tracker:team_task_tracker@localhost:15432/team_task_tracker?sslmode=disable"),
		FrontendURL: env("FRONTEND_URL", "http://localhost:5173"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
