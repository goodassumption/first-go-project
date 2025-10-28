package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

type Server struct {
	db *sql.DB
}

func initDB(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data TEXT
	);`
	db.Exec(sql)
}

// createItem обрабатывает POST-запросы для создания новой записи
func (s *Server) createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	json.NewDecoder(r.Body).Decode(&item) // Читаем JSON из тела запроса в структуру Item

	// Вставляем данные в базу, ID генерируется автоматически
	res, _ := s.db.Exec("INSERT INTO items (data) VALUES (?)", item.Data)
	id, _ := res.LastInsertId() // Получаем сгенерированный ID
	item.ID = int(id)

	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок ответа
	json.NewEncoder(w).Encode(item)                    // Отправляем обратно созданный объект с его новым ID
}

func (s *Server) getItems(w http.ResponseWriter, r *http.Request) {
	rows, _ := s.db.Query("SELECT id, data FROM items") // Выбираем все записи
	defer rows.Close()                                  // Закрываем строки после завершения функции

	var items []Item // Слайс для сбора всех записей
	for rows.Next() {
		var item Item
		rows.Scan(&item.ID, &item.Data) // Сканируем данные из строки базы в структуру
		items = append(items, item)     // Добавляем в список
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items) // Отправляем список всех записей в формате JSON
}

func main() {
	db, _ := sql.Open("sqlite3", "example.sqlite") // Открываем/создаем файл базы данных
	defer db.Close()                               // Отложенное закрытие базы данных при выходе из main

	initDB(db) // Инициализируем базу данных (создаем таблицу)

	s := &Server{db: db} // Создаем экземпляр сервера

	// Настройка маршрутизатора
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // Если метод POST - создаем
			s.createItem(w, r)
			return
		}
		if r.Method == "GET" { // Если метод GET - получаем
			s.getItems(w, r)
			return
		}
	})

	// Запуск HTTP-сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
