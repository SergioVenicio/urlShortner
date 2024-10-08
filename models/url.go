package models

import (
	"time"
)

type Hit struct {
	Time       time.Time   `json:"time"`
	RequestURI string      `json:"request_uri"`
	RemoteAddr string      `json:"remote_addr"`
	Host       string      `json:"host"`
	Header     interface{} `json:"header"`
	Method     string      `json:"method"`
}

type URLMetadata struct {
	URL  string `json:"url"`
	Hits []*Hit `json:"hits"`
}

func (m *URLMetadata) AddHit(h *Hit) {
	m.Hits = append(m.Hits, h)
}

type URL struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}
