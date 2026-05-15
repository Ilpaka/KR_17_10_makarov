# go-products-api

REST API на Go (Gin): каталог товаров с JWT-аутентификацией и Swagger UI. Учебный проект.

## Стек

- Go 1.25
- [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк
- [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) — JWT-токены
- [swaggo/swag](https://github.com/swaggo/swag) + [gin-swagger](https://github.com/swaggo/gin-swagger) — OpenAPI 2.0 документация
- [joho/godotenv](https://github.com/joho/godotenv) — загрузка `.env`

## Требования

- Go ≥ 1.25
- (опционально) `swag` для регенерации Swagger-документации:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

## Установка

```bash
git clone https://github.com/Ilpaka/go-products-api.git
cd go-products-api
go mod download
```

## Переменные окружения

Создайте файл `.env` в корне проекта:

```
JWT_SECRET=your_secret_here
```

| Переменная   | Обязательная | Описание                                    |
|--------------|--------------|---------------------------------------------|
| `JWT_SECRET` | да           | Секрет для подписи и проверки JWT-токенов   |

## Запуск

```bash
go run ./cmd/server
```

Сервер слушает `:8090`.

## Swagger UI

После запуска документация доступна по адресу:
<http://localhost:8090/swagger/index.html>

## Аутентификация

Все эндпоинты `/products...` и `/add_products` защищены JWT. Сначала получите токен на `POST /auth/token`, затем передавайте его в заголовке `Authorization: Bearer <token>`.

Пример с `curl`:

```bash
TOKEN=$(curl -s -X POST localhost:8090/auth/token | python3 -c 'import sys,json; print(json.load(sys.stdin)["token"])')

curl -H "Authorization: Bearer $TOKEN" localhost:8090/products
```

## Эндпоинты

| Метод  | Путь              | Auth   | Описание                                                    |
|--------|-------------------|--------|-------------------------------------------------------------|
| POST   | `/auth/token`     | —      | Выдача JWT-токена (TTL 24ч)                                 |
| GET    | `/products`       | Bearer | Список товаров. Фильтры: `min_price`, `max_price`, `in_stock` |
| GET    | `/products/:id`   | Bearer | Получить товар по id                                        |
| POST   | `/add_products`   | Bearer | Создать товар                                               |
| PUT    | `/products/:id`   | Bearer | Обновить товар по id                                        |
| DELETE | `/products/:id`   | Bearer | Удалить товар по id                                         |

### Модель Products

```json
{
  "id": 1,
  "name": "Apple",
  "price": 100,
  "in_stock": 5
}
```

## Структура проекта

```
.
├── cmd/server/             # точка входа
│   └── main.go
├── internal/
│   ├── config/             # загрузка .env → Config
│   ├── model/              # доменные типы (Products, ErrorResponse, TokenResponse)
│   ├── repository/         # in-memory хранилище
│   ├── service/            # бизнес-логика, фильтры
│   ├── middleware/         # JWT-мидлварь
│   └── handler/            # Gin-хендлеры и роутер
├── docs/                   # сгенерированный Swagger (docs.go, swagger.json/yaml)
├── postman/                # Postman-коллекция
├── _template/              # референсный сниппет news handler (исключён из сборки)
├── .env                    # JWT_SECRET (не коммитить в реальном проекте)
├── go.mod / go.sum
└── task.md                 # исходная постановка задачи
```

## Postman

В `postman/postman_collection.json` лежит готовая коллекция со всеми эндпоинтами.

## Регенерация Swagger

После изменения аннотаций (`@Summary`, `@Param`, `@Success` и т.д.) или моделей:

```bash
swag init -g cmd/server/main.go --md docs --output docs --parseInternal
```

Флаг `--parseInternal` необходим, чтобы swag увидел типы из `internal/model`.
