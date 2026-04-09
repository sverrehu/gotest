package internal

import (
	"time"
)

type Release struct {
	Version    string    `json:"version"`
	ReleasedAt time.Time `json:"released_at"`
}
