package target

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type fileConfig struct {
	Targets []fileTarget `json:"targets"`
}

type fileTarget struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval string `json:"interval"`
}

type Target struct {
	Name     string
	URL      string
	Interval time.Duration
}

func Load(path string) ([]Target, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	var targets []Target
	for i, ft := range cfg.Targets {
		d, err := time.ParseDuration(ft.Interval)
		if err != nil {
			return nil, fmt.Errorf("targets[%d]: interval: %w", i, err)
		}
		if ft.Name == "" {
			return nil, fmt.Errorf("targets[%d]: name is empty", i)
		}
		if ft.URL == "" {
			return nil, fmt.Errorf("targets[%d]: url is empty", i)
		}
		if d <= 0 {
			return nil, fmt.Errorf("targets[%d]: interval is invalid", i)
		}
		targets = append(targets, Target{Name: ft.Name, URL: ft.URL, Interval: d})
	}
	return targets, nil
}
