package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Movie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	ActorCount int    `json:"actor_count"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "postgres://user:password@localhost:5432/moviesdb?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/movies", getMoviesHandler)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	q := r.URL.Query()
	yearMinStr := q.Get("year_min")
	yearMaxStr := q.Get("year_max")
	limitStr := q.Get("limit")
	offsetStr := q.Get("offset")

	limit := 10
	offset := 0
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
	}

	query := `
		SELECT m.id, m.title, m.year, COUNT(a.id) AS actor_count
		FROM movies m
		LEFT JOIN actors a ON m.id = a.movie_id
		WHERE 1=1
	`
	params := []interface{}{}
	paramIndex := 1

	if yearMinStr != "" {
		query += fmt.Sprintf(" AND m.year >= $%d", paramIndex)
		params = append(params, yearMinStr)
		paramIndex++
	}
	if yearMaxStr != "" {
		query += fmt.Sprintf(" AND m.year <= $%d", paramIndex)
		params = append(params, yearMaxStr)
		paramIndex++
	}

	query += fmt.Sprintf(`
		GROUP BY m.id
		ORDER BY m.year DESC
		LIMIT $%d OFFSET $%d
	`, paramIndex, paramIndex+1)
	params = append(params, limit, offset)

	rows, err := db.Query(query, params...)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.ActorCount); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		movies = append(movies, m)
	}

	queryTime := time.Since(start).Milliseconds()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Query-Time", fmt.Sprintf("%dms", queryTime))
	json.NewEncoder(w).Encode(movies)
}
