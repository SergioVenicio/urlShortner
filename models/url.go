package models

import (
	"time"
)

type Hit struct {
	Time       time.Time   `json:"time"`
	Remote     string      `json:"remote"`
	RequestURI string      `json:"request_uri"`
	RemoteAddr string      `json:"remote_addr"`
	Host       string      `json:"host"`
	Header     interface{} `json:"header"`
	Method     string      `json:"method"`
}

type URLMetadata struct {
	Hits []*Hit `json:"hits"`
}

func (m *URLMetadata) AddHit(h *Hit) {
	m.Hits = append(m.Hits, h)
}

type URL struct {
	ID       string      `json:"id"`
	Source   string      `json:"source"`
	Metadata URLMetadata `json:"metadata"`
}

func NewURL(id string, source string, hits int) *URL {
	return &URL{
		ID:       id,
		Source:   source,
		Metadata: URLMetadata{Hits: []*Hit{}},
	}
}
