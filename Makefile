SHELL := /bin/sh

.PHONY: help dev down logs ps db-up backend-dev backend-test frontend-install frontend-dev

help:
	@printf '%s\n' 'Available commands:'
	@printf '%s\n' '  make dev              Start local Docker stack'
	@printf '%s\n' '  make down             Stop local Docker stack'
	@printf '%s\n' '  make logs             Follow Docker logs'
	@printf '%s\n' '  make ps               Show Docker services'
	@printf '%s\n' '  make db-up            Start PostgreSQL only'
	@printf '%s\n' '  make backend-dev      Run backend locally'
	@printf '%s\n' '  make backend-test     Run Go tests'
	@printf '%s\n' '  make frontend-install Install frontend dependencies'
	@printf '%s\n' '  make frontend-dev     Run frontend dev server locally'

dev:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

db-up:
	docker compose up -d postgres

backend-dev:
	cd backend && go run ./cmd/api

backend-test:
	cd backend && go test ./...

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev
