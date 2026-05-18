SHELL := /bin/sh

.PHONY: help dev down logs ps backend-test frontend-install frontend-dev

help:
	@printf '%s\n' 'Available commands:'
	@printf '%s\n' '  make dev              Start local Docker stack'
	@printf '%s\n' '  make down             Stop local Docker stack'
	@printf '%s\n' '  make logs             Follow Docker logs'
	@printf '%s\n' '  make ps               Show Docker services'
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

backend-test:
	cd backend && go test ./...

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

