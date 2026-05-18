package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var migrationFilePattern = regexp.MustCompile(`^(\d+)_.+\.sql$`)

type AppliedMigration struct {
	Version int64
	Name    string
}

type migrationFile struct {
	Version int64
	Name    string
	Path    string
}

func Up(ctx context.Context, db *pgxpool.Pool, dir string) ([]AppliedMigration, error) {
	files, err := loadMigrationFiles(dir)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint PRIMARY KEY,
			name text NOT NULL,
			applied_at timestamptz NOT NULL DEFAULT now()
		)
	`); err != nil {
		return nil, fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(1739018301)`); err != nil {
		return nil, fmt.Errorf("acquire migration lock: %w", err)
	}

	applied := make([]AppliedMigration, 0, len(files))
	for _, file := range files {
		alreadyApplied, err := isApplied(ctx, tx, file.Version)
		if err != nil {
			return nil, err
		}
		if alreadyApplied {
			continue
		}

		contents, err := os.ReadFile(file.Path)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", file.Name, err)
		}

		upSQL, err := extractUpSQL(string(contents))
		if err != nil {
			return nil, fmt.Errorf("parse migration %s: %w", file.Name, err)
		}

		if _, err := tx.Exec(ctx, upSQL); err != nil {
			return nil, fmt.Errorf("apply migration %s: %w", file.Name, err)
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO schema_migrations (version, name)
			VALUES ($1, $2)
		`, file.Version, file.Name); err != nil {
			return nil, fmt.Errorf("record migration %s: %w", file.Name, err)
		}

		applied = append(applied, AppliedMigration{
			Version: file.Version,
			Name:    file.Name,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit migrations: %w", err)
	}

	return applied, nil
}

func loadMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	files := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		matches := migrationFilePattern.FindStringSubmatch(name)
		if len(matches) != 2 {
			continue
		}

		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse migration version %s: %w", name, err)
		}

		files = append(files, migrationFile{
			Version: version,
			Name:    name,
			Path:    filepath.Join(dir, name),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	return files, nil
}

func extractUpSQL(contents string) (string, error) {
	upMarker := "-- +goose Up"
	downMarker := "-- +goose Down"

	upIndex := strings.Index(contents, upMarker)
	if upIndex == -1 {
		return "", fmt.Errorf("missing %q marker", upMarker)
	}

	body := contents[upIndex+len(upMarker):]
	if downIndex := strings.Index(body, downMarker); downIndex != -1 {
		body = body[:downIndex]
	}

	body = strings.TrimSpace(body)
	if body == "" {
		return "", fmt.Errorf("empty up migration")
	}

	return body, nil
}

func isApplied(ctx context.Context, tx pgx.Tx, version int64) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM schema_migrations
			WHERE version = $1
		)
	`, version).Scan(&exists); err != nil {
		return false, fmt.Errorf("check migration version %d: %w", version, err)
	}

	return exists, nil
}
