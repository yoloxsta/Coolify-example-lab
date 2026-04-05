package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Notification struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "appuser"),
		getEnv("DB_PASSWORD", "apppass"),
		getEnv("DB_NAME", "appdb"),
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		message TEXT NOT NULL,
		read BOOLEAN DEFAULT false
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/notifications", corsMiddleware(notificationsHandler))
	http.HandleFunc("/api/notifications/", corsMiddleware(notificationActionHandler))

	log.Println("Notification service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		rows, err := db.Query(`SELECT id, message, "read" FROM notifications ORDER BY id DESC`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		notifs := []Notification{}
		for rows.Next() {
			var n Notification
			rows.Scan(&n.ID, &n.Message, &n.Read)
			notifs = append(notifs, n)
		}
		json.NewEncoder(w).Encode(notifs)

	case "POST":
		var n Notification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			http.Error(w, "invalid json", 400)
			return
		}
		err := db.QueryRow(`INSERT INTO notifications (message) VALUES ($1) RETURNING id`, n.Message).Scan(&n.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		n.Read = false
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(n)

	default:
		http.Error(w, "method not allowed", 405)
	}
}

func notificationActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	id := parts[len(parts)-1]

	switch r.Method {
	case "PATCH":
		_, err := db.Exec(`UPDATE notifications SET "read" = true WHERE id = $1`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"marked": id})

	case "DELETE":
		result, err := db.Exec(`DELETE FROM notifications WHERE id = $1`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"deleted": id})

	default:
		http.Error(w, "method not allowed", 405)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
