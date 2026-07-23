# Activity / Analytics Service

Дополнительный сервис для LMS Microservices. Он не участвует в основной бизнес-логике: основные сервисы только публикуют события в Kafka, а аналитика читает их асинхронно.

## Что делает сервис

- Читает Kafka topics: `users.events`, `courses.events`, `homework.events`.
- Сохраняет все события в `activity_events`.
- Обновляет `course_stats`.
- Обновляет `user_activity_stats`.
- Публикует собственные события в `analytics.events`:
  - `activity.saved`
  - `course.stats.updated`
  - `user.stats.updated`

## API

Сервис ожидает headers от Gateway:

```http
X-User-ID: 1
X-User-Role: admin
```

### Health check

```http
GET /health
```

### История активности

```http
GET /activity/events?limit=50&offset=0
```

Доступ: `admin`.

### Активность пользователя

```http
GET /activity/users/:id?limit=50&offset=0
```

Доступ: `admin`.

### Статистика курса

```http
GET /analytics/courses/:id
```

Доступ: `teacher`, `admin`.

### Общий dashboard

```http
GET /analytics/dashboard
```

Доступ: `admin`.

## Таблицы

- `activity_events`
- `course_stats`
- `user_activity_stats`

Миграции выполняются через `AutoMigrate` при старте сервиса.

## Пример docker-compose блока

```yaml
activity-analytics-service:
  build:
    context: ./activity-analytics-service
  ports:
    - "8084:8084"
  environment:
    HTTP_PORT: "8084"
    DATABASE_DSN: "host=postgres user=postgres password=postgres dbname=analytics_db port=5432 sslmode=disable"
    KAFKA_BROKERS: "kafka:9092"
    KAFKA_GROUP_ID: "activity-analytics-service"
    KAFKA_TOPICS: "users.events,courses.events,homework.events"
    ANALYTICS_TOPIC: "analytics.events"
  depends_on:
    - postgres
    - kafka
```

Если в проекте используются отдельные базы, нужно создать `analytics_db` в PostgreSQL. Если используются схемы, поменяй `DATABASE_DSN` и настройки миграции под схему команды.

## Проверка вручную

После запуска можно отправить тестовое событие в `courses.events`:

```json
{
  "event": "student.enrolled",
  "course_id": 1,
  "student_id": 15,
  "created_at": "2026-06-09T12:20:00Z"
}
```

После обработки:

- `/activity/events` покажет сохраненное событие.
- `/analytics/courses/1` покажет `students_count: 1`.
- `/activity/users/15` покажет `courses_enrolled_count: 1`.
