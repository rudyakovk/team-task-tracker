# Team Task Tracker

Локальный team task tracker для небольших команд.

Текущий этап: базовый scaffold проекта. Бизнес-логика будет добавляться постепенно по плану из [docs/mvp-plan.md](docs/mvp-plan.md).

## Stack

- Backend: Go
- Frontend: React + TypeScript
- Database: PostgreSQL
- Local infrastructure: Docker Compose

## Local Development

Пока проект находится на Phase 0, доступны базовые сервисы:

- frontend: `http://localhost:5173`
- backend health: `http://localhost:8080/healthz`
- PostgreSQL: `localhost:15432`

Команды:

```sh
make dev
make down
make logs
make backend-test
```

Для локального запуска frontend без Docker:

```sh
make frontend-install
make frontend-dev
```

## Environment

Шаблон переменных окружения находится в `.env.example`.

Для Docker Compose сейчас используются безопасные development defaults. Позже, когда появятся миграции и подключение backend к базе, env-конфигурация будет расширена.
