package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/checker/check"
	"example.com/checker/store"
	"example.com/checker/target"
)

func main() {
	log.SetFlags(0)

	targets, err := loadTargets("targets.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := store.Open("checker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, t := range targets {
		r := check.Run(t)
		if _, err := db.Save(r); err != nil {
			log.Fatal(err)
		}
	}

	recent, err := db.ListRecent(10)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range recent {
		if r.Error != "" {
			fmt.Printf("%s  %-10s NG  %s\n", r.CheckedAt.Format(time.TimeOnly), r.Name, r.Error)
			continue
		}
		fmt.Printf("%s  %-10s %d  %dms\n", r.CheckedAt.Format(time.TimeOnly), r.Name, r.StatusCode, r.LatencyMs)
	}
}

// Unmarshal は構造体タグどおりに値を流し込むだけで New の検証を通らないため、
// 読み込み後に 1 件ずつ New に通して不正な設定を弾く。
func loadTargets(path string) ([]target.Target, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var raw []target.Target
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	targets := make([]target.Target, 0, len(raw))
	for i, r := range raw {
		t, err := target.New(r.Name, r.URL)
		if err != nil {
			return nil, fmt.Errorf("%s[%d]: %w", path, i, err)
		}
		targets = append(targets, t)
	}
	return targets, nil
}
