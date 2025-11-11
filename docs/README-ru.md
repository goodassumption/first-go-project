# Game Results Server (Go + Python клиент)

Сервер на Go для сохранения результатов игры и Python-клиент для отправки обновлений и чтения таблицы лидеров. Вся коммуникация происходит через HTTP-маршруты.

## Архитектура и стек

- Языки: Go (сервер), Python (клиент)
- Сервер: два маршрута
  - POST /score/update — обновление очков игрока
  - GET /leaderboard — получение текущей таблицы лидеров
- Хранение данных: SQLite (файл db.sqlite по умолчанию)
- Клиент: Python-скрипт src/append_user.py, использующий requests для взаимодействия с сервером

## Структура проекта

```text
ProjectRoot
├── docs
│ ├── LICENSE-en.md
│ ├── LICENSE-ru.md
| ├── README-en.md
│ └── README-ru.md
├── src
| ├── append_user.py # Python client: sends data to the server and reads the leaderboard
| ├── main.go # Go server: handles requests, works with the database
| ├── go.mod # Go dependencies (module)
| ├── go.sum # Go dependency checksums
| └── requirements.txt # Python client dependencies
├── LICENCE.md
├── README.md
└── .gitignore # Git exceptions
```

## Как запустить

1) Установка зависимостей

   - Go:
  
    ```bash
    go mod download
    ```

   - Python (клиент):

    ```bash
    pip install -r requirements.txt
    ```

2) Запуск сервера

   - По умолчанию сервер слушает на порту 9090 и использует файл базы данных db.sqlite
   - Пример запуска:

    ```bash
     go run main.go
     ```

    или

     ```bash
     go build ./... && ./your_executable
     ```

3) Запуск Python-клиента

- Пример запуска клиента (использует сервер по умолчанию <http://localhost:9090>):

    ```bash
    python3 src/append_user.py
    ```

- Если нужно изменить адрес сервера, можно отредактировать константу SERVER_URL в src/append_user.py

---

## Как работают маршруты

- POST /score/update
  - Назначение: обновление очков игрока
  - Тело запроса (JSON):
    {
      "id": "string",
      "name": "string",
      "addScr": int
    }
  - Ответ: JSON-ответ сервера с полями Updated, CurTop, NewTop, Changed

- GET /leaderboard
  - Назначение: получение текущей таблицы лидеров
  - Параметр запроса: limit (необязательно) — максимальное количество записей
  - Ответ: JSON вида {"entries": [{ "Rank": int, "Id": "string", "Name": "string", "Scr": int }...]}

---

## Примеры использования

- Пример запроса для обновления очков (curl):
  
  ```bash
  curl -X POST <http://localhost:9090/score/update> \
    -H "Content-Type: application/json" \
    -d '{"id": "user123", "name": "Иван", "addScr": 10}'
    ```

- Пример запроса для получения топа (curl):

    ```bash
    curl <http://localhost:9090/leaderboard?limit=5>
    ```

- Пример использования Python-клиента (src/append_user.py) с сервером по умолчанию:

    ```bash
    python3 src/append_user.py
    ```

Обратите внимание: структура и форматы данных соответствуют определению в коде сервера и клиента.

---

## Хранение данных

- По умолчанию данные сохраняются в локальной SQLite-базе db.sqlite.
- Таблица scores имеет поля: id (TEXT, PRIMARY KEY), name (TEXT), scr (INTEGER).
- В продакшн-режиме можно адаптировать хранение, но в текущей версии это локальная база данных.

---

## Конфигурация

- PORT и другие параметры заложены напрямую в код:
  - Порт: порт 9090 (константа port)
  - База данных: db.sqlite (имя файла задаётся константой name)
  - Максимальное число записей в топе: limit = 10
- Чтобы поменять настройки, отредактируйте соответствующие константы в main.go и перекомпилируйте проект.

---

## Логирование и наблюдаемость

- Логи сервера выводятся в stdout/stderr. При необходимости можно расширить логирование и добавить трассировку.

---

## Тестирование

- Ручное тестирование через curl и Python-клиент.
- Юнит и интеграционные тесты можно добавить позже с использованием стандартных инструментов Go и Python.

---

## Вклад

- Форкните репозиторий, создайте ветку feature/your-feature
- Внесите изменения и создайте pull request
- Добавляйте тесты и соответствующую документацию

---

## Лицензия

- Проект лицензирован [MIT license](LICENSE-ru.md)
