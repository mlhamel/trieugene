package store

import (
	"time"
)

type Outflow struct {
	ScrapedAt time.Time `json:"scraped_at"`
	TakenAt   time.Time `json:"taken_at"`
	Outflow   float64   `json:"outflow"`
}
