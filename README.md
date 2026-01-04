# Rest API

## Архитектура
Приложение поднимается при помощи команды ```docker-compose up --build -d```(не забудьте создать .env, пример приведён в разделе [Конфигурация](#custom-id)). Она запускает по очереди два докер-контейнера(контейнер с postgresql версии 16 и контейнер с go-приложением, поднимающим http-сервер на 8080 порту(Go версии 1.25)). Также после поднятия БД автоматически запускаются миграции, которые инициализируют таблицы и заполняют их тестовыми данными

## Эндпоинты

Запросы можно отсылать на 8080 порт localhost из системы, в которой поднят docker-compose

**GET http://127.0.0.1:8080/health** - проверка состояния сервера
**GET http://127.0.0.1:8080/api/v1/wallet/{{WALLET_UUID}}** - узнать детали кошелька по его UUID
**POST http://127.0.0.1:8080/api/v1/wallet** - совершить транзакцию

**Пример необходимого тела для POST-запроса:**
{
    "walletId": "893be64f-664a-43c2-9840-964e5e0d594f",
    "operationType": "WITHDRAW",
    "amount": 2001
}

<a id="custom-id"></a>
## Конфигурация

Для корректной инициализации системы требуется создать .env файле в корне проекта. Пример содержимого:

```
APP_PORT=8080

# Database
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=mydatabase

DB_MAX_CONNS=20
DB_MIN_CONNS=5
DB_MAX_CONN_LIFETIME=1h
```