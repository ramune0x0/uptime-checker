// 疎通確認の結果を SQLite に保存・読み出しする。
package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"example.com/checker/check"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping %s: %w", path, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS result (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			url         TEXT NOT NULL,
			status_code INTEGER NOT NULL,
			latency_ms  INTEGER NOT NULL,
			error       TEXT NOT NULL,
			checked_at  TEXT NOT NULL
		)`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Save(r check.Result) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO result (name, url, status_code, latency_ms, error, checked_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		r.Name, r.URL, r.StatusCode, r.LatencyMs, r.Error, r.CheckedAt.Format(time.RFC3339),
	)
	if err != nil {
		return 0, fmt.Errorf("save %s: %w", r.Name, err)
	}
	return res.LastInsertId()
}

func (s *Store) ListRecent(limit int) ([]check.Result, error) {
	rows, err := s.db.Query(
		`SELECT name, url, status_code, latency_ms, error, checked_at
		 FROM result ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent: %w", err)
	}
	defer rows.Close()

	var results []check.Result
	for rows.Next() {
		var r check.Result
		var checkedAt string
		if err := rows.Scan(&r.Name, &r.URL, &r.StatusCode, &r.LatencyMs, &r.Error, &checkedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		r.CheckedAt, _ = time.Parse(time.RFC3339, checkedAt)
		results = append(results, r)
	}
	return results, rows.Err()
}
