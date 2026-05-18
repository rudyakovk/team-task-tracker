package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"team-task-tracker/backend/internal/config"
	"team-task-tracker/backend/internal/database"
)

type seedConfig struct {
	WorkspaceName string
	AdminEmail    string
	AdminUsername string
	AdminPassword string
	AdminName     string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.Load()
	seed := seedConfig{
		WorkspaceName: env("SEED_WORKSPACE_NAME", "Local Workspace"),
		AdminEmail:    env("SEED_ADMIN_EMAIL", "admin@example.com"),
		AdminUsername: env("SEED_ADMIN_USERNAME", "admin"),
		AdminPassword: env("SEED_ADMIN_PASSWORD", "admin12345"),
		AdminName:     env("SEED_ADMIN_NAME", "Admin"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		logger.Error("begin seed transaction failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var workspaceID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO workspaces (name)
		SELECT $1
		WHERE NOT EXISTS (
			SELECT 1 FROM workspaces WHERE name = $1
		)
		RETURNING id
	`, seed.WorkspaceName).Scan(&workspaceID); err != nil {
		if err := tx.QueryRow(ctx, `
			SELECT id FROM workspaces WHERE name = $1
		`, seed.WorkspaceName).Scan(&workspaceID); err != nil {
			logger.Error("ensure workspace failed", "error", err)
			os.Exit(1)
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seed.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("hash admin password failed", "error", err)
		os.Exit(1)
	}

	var adminID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash, display_name, is_active)
		VALUES ($1, $2, $3, $4, true)
		ON CONFLICT (email) DO UPDATE SET
			username = EXCLUDED.username,
			password_hash = EXCLUDED.password_hash,
			display_name = EXCLUDED.display_name,
			is_active = true
		RETURNING id
	`, seed.AdminEmail, seed.AdminUsername, string(passwordHash), seed.AdminName).Scan(&adminID); err != nil {
		logger.Error("ensure admin user failed", "error", err)
		os.Exit(1)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO workspace_members (workspace_id, user_id, role)
		VALUES ($1, $2, 'admin')
		ON CONFLICT (workspace_id, user_id) DO UPDATE SET role = 'admin'
	`, workspaceID, adminID); err != nil {
		logger.Error("ensure admin membership failed", "error", err)
		os.Exit(1)
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("commit seed transaction failed", "error", err)
		os.Exit(1)
	}

	logger.Info(
		"seed data ready",
		"workspace", seed.WorkspaceName,
		"admin_email", seed.AdminEmail,
		"admin_username", seed.AdminUsername,
	)
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
