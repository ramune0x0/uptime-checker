package target

import (
	"errors"
	"strings"
)

type Target struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func hasScheme(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func New(name, url string) (Target, error) {
	if name == "" {
		return Target{}, errors.New("empty name")
	}
	if url == "" {
		return Target{}, errors.New("empty url")
	}
	if !hasScheme(url) {
		return Target{}, errors.New("url must start with http:// or https://")
	}
	return Target{Name: name, URL: url}, nil
}
