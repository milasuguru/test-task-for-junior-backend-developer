# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go 1.23+
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

docker compose up --build

После запуска сервис будет доступен по адресу http://localhost:8080.

Если postgres уже запускался ранее со старой схемой, пересоздай volume:

docker compose down -v
docker compose up --build

## Локальный запуск (без Docker)

DATABASE_DSN="postgres://postgres:postgres@localhost:5432/taskservice?sslmode=disable" go run ./cmd/api

## Swagger

Swagger UI: http://localhost:8080/swagger/
OpenAPI JSON: http://localhost:8080/swagger/openapi.json

## API

Базовый префикс: /api/v1

### Задачи
- POST   /api/v1/tasks
- GET    /api/v1/tasks
- GET    /api/v1/tasks/{id}
- PUT    /api/v1/tasks/{id}
- DELETE /api/v1/tasks/{id}

### Расписания
- POST   /api/v1/schedules
- GET    /api/v1/schedules
- GET    /api/v1/schedules/{id}
- DELETE /api/v1/schedules/{id}

## Типы расписаний

### daily — каждый N-й день
{"title": "Обзвон пациентов", "type": "daily", "every_n_days": 2}

### monthly — определённое число месяца (1-30)
{"title": "Инвентаризация", "type": "monthly", "month_day": 15}

### specific — конкретные даты
{"title": "Плановый осмотр", "type": "specific", "specific_dates": ["2026-05-28", "2026-06-10"]}

### parity — чётные или нечётные дни
{"title": "Отчётность", "type": "parity", "parity": "even"}
{"title": "Обход пациентов", "type": "parity", "parity": "odd"}

## Принятые решения и допущения

1. При создании расписания задачи сразу генерируются на generate_days_ahead дней вперёд (по умолчанию 30).
2. Каждую ночь в полночь cron-job генерирует новые задачи по всем активным расписаниям.
3. Для типа monthly числа 29 и 30 пропускаются в месяцах где их нет (февраль и др.).
4. Для типа specific даты за пределами горизонта генерации сохраняются, но задачи по ним создадутся в следующем цикле cron-job.
5. Поле generate_days_ahead позволяет настроить горизонт генерации для каждого расписания отдельно.
6. Удаление расписания не удаляет уже созданные задачи.