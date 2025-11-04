package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	port  = ":9090"
	name  = "db.sqlite"
	limit = 10
)

type UpdReq struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	AddScr int64  `json:"addScr"`
}

type Entry struct {
	Rank int
	Id   string
	Name string
	Scr  int64
}

type UpdResp struct {
	Updated bool    `json:"updated"`
	CurTop  []Entry `json:"curTop"`
	NewTop  []Entry `json:"newTop"`
	Changed bool    `json:"changed"`
}

type Server struct {
	db *sql.DB
}

func initDB(db *sql.DB) error {
	const sql = `
	CREATE TABLE IF NOT EXISTS scores (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		scr INTEGER
	);`

	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы: %v", err)
	}
	return nil
}

func (s *Server) getTop(lmt int) ([]Entry, error) {
	query := fmt.Sprintf(`
		SELECT id, name, scr
		FROM scores
		ORDER BY scr DESC
		LIMIT %d
	`, lmt)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе топа: %v", err)
	}
	defer rows.Close()

	var entry []Entry
	count := 1
	for rows.Next() {
		var id, name string
		var scr int64
		if err := rows.Scan(&id, &name, &scr); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки топа: %v", err)
		}
		entry = append(entry, Entry{
			Rank: count,
			Id:   id,
			Name: name,
			Scr:  scr,
		})
		count++
	}
	return entry, nil
}

func (s *Server) updScr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "неверный формат JSON"}`, http.StatusBadRequest)
		return
	}

	curTop, err := s.getTop(limit)
	if err != nil {
		http.Error(w, `{"error": "ошибка чтения топа"}`, http.StatusInternalServerError)
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, `{"error": "ошибка транзакции"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var oldScr int64
	row := tx.QueryRow("SELECT scr FROM scores WHERE id = ?", req.Id)
	if err := row.Scan(&oldScr); err != nil && err != sql.ErrNoRows {
		return
	}

	newScr := oldScr + req.AddScr

	res, err := tx.Exec(`
		INSERT INTO scores (id, name, scr)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name,
			scr=excluded.scr
	`, req.Id, req.Name, newScr)

	var updated bool
	if err != nil {
		http.Error(w, `{"error": "ошибка обновления"}`, http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	updated = rows > 0

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error": "ошибка коммита"}`, http.StatusInternalServerError)
		return
	}

	newTop, err := s.getTop(limit)
	if err != nil {
		http.Error(w, `{"error": "ошибка чтения нового топа"}`, http.StatusInternalServerError)
		return
	}

	changed := false
	if len(curTop) != len(newTop) {
		changed = true
	} else {
		if len(curTop) > 0 {
			for i := 0; i < len(curTop); i++ {
				if curTop[i].Id != newTop[i].Id || curTop[i].Scr != newTop[i].Scr {
					changed = true
					break
				}
			}
		} else if len(newTop) > 0 {
			changed = true
		}
	}

	resp := UpdResp{
		Updated: updated,
		CurTop:  curTop,
		NewTop:  newTop,
		Changed: changed,
	}

	json.NewEncoder(w).Encode(resp)
}

func (s *Server) getLdr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lmtStr := r.URL.Query().Get("limit")
	lmt := limit
	if lmtStr != "" {
		if val, err := strconv.Atoi(lmtStr); err == nil {
			lmt = val
		}
	}

	entry, err := s.getTop(lmt)
	if err != nil {
		http.Error(w, `{"error": "ошибка получения топа"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string][]Entry{"entries": entry})
}

func main() {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("не удалось открыть базу данных: %v", err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatalf("не удалось инициализировать базу данных: %v", err)
	}

	s := &Server{db: db}

	http.HandleFunc("/score/update", s.updScr)
	http.HandleFunc("/leaderboard", s.getLdr)

	log.Printf("сервер запущен на порту %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}
