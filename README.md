# Game Results Server (Go + Python client)

[Русская версия README](docs/README-ru.md)

A Go server for saving game results and a Python client for sending updates and reading the leaderboard. All communication occurs through HTTP routes.

## Architecture and Stack

- Languages: Go (server), Python (client)
- Server: two routes
  - POST /score/update - updating a player's score
  - GET /leaderboard - retrieving the current leaderboard
- Data storage: SQLite (default db.sqlite file)
- Client: Python script src/append_user.py, using requests to interact with the server

## Project structure

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

## How to run

1) Installing dependencies

    - Go:

    ```bash
    go mod download
    ```

    - Python (client):

    ```bash
    pip install -r requirements.txt
    ```

2) Starting the server

    - By default, the server listens on port 9090 and uses the db.sqlite database file
    - Example of starting the server:

    ```bash
    go run main.go
    ```

    or

     ```bash
     go build ./... && ./your_executable
     ```

3) Launching the Python client

- Example of launching the client (uses the default server <http://localhost:9090>):

 ```bash
 python3 src/append_user.py
 ```

- If you need to change the server address, you can edit the SERVER_URL constant in src/append_user.py

---

## How the routes work

- POST /score/update
- Purpose: update the player's score
- Request body (JSON):
 {
 "id": "string",
 "name": "string",
 "addScr": int
 }
- Response: a JSON server response with the fields Updated, CurTop, NewTop, and Changed

- GET /leaderboard
- Purpose: get the current leaderboard
- Query parameter: limit (optional) - maximum number of entries
- Response: JSON of the form {"entries": [{ "Rank": int, "Id": "string", "Name": "string", "Scr": int }...]}

---

## Examples of use

- Example request to update the score (curl):

 ```bash
 curl -X POST <http://localhost:9090/score/update> \
 -H "Content-Type: application/json" \
 -d '{"id": "user123", "name": "Иван", "addScr": 10}'
 ```

- Example request to get the top (curl):

 ```bash
 curl <http://localhost:9090/leaderboard?limit=5>
 ```

- Example of using the Python client (src/append_user.py) with the default server:

 ```bash
 python3 src/append_user.py
 ```

 Note: The data structure and formats match the definitions in the server and client code.

---

## Data storage

- By default, data is stored in the local SQLite database db.sqlite.
- The scores table has the following fields: id (TEXT, PRIMARY KEY), name (TEXT), and scr (INTEGER).
- In production mode, storage can be adapted, but in the current version it is a local database.

---

## Configuration

- PORT and other parameters are built directly into the code:
- Port: port 9090 (constant port)
- Database: db.sqlite (file name is set by the name constant)
- Maximum number of entries in the top: limit = 10
- To change the settings, edit the corresponding constants in main.go and recompile the project.

---

## Logging and observability

- Server logs are displayed in stdout/stderr. If necessary, you can expand logging and add tracing.

---

## Testing

- Manual testing via curl and Python client.
- Unit and integration tests can be added later using standard Go and Python tools.

---

## Contribution

- Fork the repository, create a feature/your-feature branch
- Make changes and create a pull request
- Add tests and relevant documentation

---

## License

- The project is licensed under the [GNU GPU v3 License](docs/LICENSE-en)
