# Go Todo REST API

Безопасный полнофункциональный REST API, построенный на Go, Gin, PostgreSQL и JWT-аутентификации.
Проект демонстрирует реализацию чистой архитектуры с разделением на слои, dependency injection и unit-тестируемым кодом.

## Возможности

- Аутентификация пользователей — регистрация и вход с использованием JWT-токенов
- Безопасность паролей — хеширование паролей с помощью bcrypt
- Защищённые маршруты — защита через JWT middleware
- Персональные задачи — у каждого пользователя собственная коллекция todo
- CRUD операции — создание, чтение, обновление и удаление задач
- Чистая архитектура — разделение на domain, ports, service, repository, handlers
- Graceful shutdown — корректное завершение активных запросов при остановке
- Миграции базы данных — версионные изменения схемы БД
- Hot Reloading — интеграция с Air для разработки

## API Endpoints

Публичные маршруты
| Метод | Endpoint         | Описание                        |
| ----- | ---------------- | ------------------------------- |
| GET   | `/`              | Проверка работы API             |
| POST  | `/auth/register` | Регистрация нового пользователя |
| POST  | `/auth/login`    | Вход и получение токена         |

Защищённые маршруты (JWT)
| Метод  | Endpoint      | Описание                         |
| ------ | ------------- | -------------------------------- |
| POST   | `/todo`       | Создать новую задачу             |
| GET    | `/todo`       | Получить все задачи пользователя |
| GET    | `/todo/:id`   | Получить конкретную задачу       |
| PATCH  | `/todo/:id`   | Обновить задачу                  |
| DELETE | `/todo/:id`   | Удалить задачу                   |

## Структура проекта

```
Go-Gin-Postgres-Todo-REST-API/
├── cmd/
│   └── api/
│       └── main.go                    # Composition root — собирает все слои
│
├── pkg/
│   ├── jwtutil/
│   │   └── jwtutil.go                 # GenerateToken(), ParseToken()
│   └── hasher/
│       └── bcrypt.go                  # Hash(), Compare()
│
├── internal/
│   ├── config/
│   │   └── config.go                  # Конфигурация окружения
│   │
│   ├── database/
│   │   └── postgres.go                # Подключение к базе данных
│   │
│   ├── domain/                        # Сущности и доменные ошибки
│   │   ├── todo.go                    # struct Todo, CreateTodoInput, UpdateTodoInput
│   │   ├── user.go                    # struct User, CreateUserInput
│   │   └── errors.go                  # ErrNotFound, ErrAlreadyExists,ErrUnauthorized,ErrInvalidInput
│   │
│   ├── ports/                         # Интерфейсы — контракты между слоями
│   │   ├── todo_repository.go         # interface TodoRepository
│   │   ├── user_repository.go         # interface UserRepository
│   │   ├── todo_service.go            # interface TodoService
│   │   └── user_service.go            # interface UserService
│   │
│   ├── service/                       # Бизнес-логика
│   │   ├── todo_service.go            # Реализует TodoService
│   │   └── user_service.go            # Реализует UserService
│   │
│   ├── repository/                    # Работа с базой данных
│   │   ├── todo_repository.go         # Реализует TodoRepository
│   │   └── user_repository.go         # Реализует UserRepository
│   │
│   ├── handlers/                      # HTTP обработчики
│   │   ├── todo_handler.go            # Маршруты todo
│   │   └── user_handler.go            # Маршруты аутентификации
│   │
│   └── middleware/
│       └── auth_middleware.go         # JWT middleware
│
├── migrations/
│   ├── 000001_create_todos_api_table.up.sql
│   ├── 000001_create_todos_api_table.down.sql
│   ├── 000002_create_users_api_table.up.sql
│   ├── 000002_create_users_api_table.down.sql
│   ├── 000003_add_user_id_to_todos_table.up.sql
│   └── 000003_add_user_id_to_todos_table.down.sql
│
├── .air.toml
├── .env.example
├── go.mod
└── go.sum
```

## Архитектура

Проект построен по принципам чистой архитектуры. Зависимости направлены строго в одну сторону:

```
handlers → ports ← service → ports ← repository → PostgreSQL
               ↑                ↑
            domain            domain
```

Каждый слой знает только о следующем через интерфейс из `ports/` — это позволяет тестировать `service/` без реальной базы данных.

## Используемые технологии

```
Go           — язык программирования backend
Gin          — веб-фреймворк
PostgreSQL   — реляционная база данных
pgx/v5       — драйвер PostgreSQL
JWT          — аутентификация
bcrypt       — хеширование паролей
golang-migrate — миграции базы данных
Air          — hot reload при разработке
godotenv     — управление переменными окружения
```

