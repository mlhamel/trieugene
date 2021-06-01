package scraper

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/scraper"
)

type DummyScraper struct {
	cfg *config.Config
}

func NewDummyScraper(cfg *config.Config) scraper.Scraper {
	return &HttpScraper{
		cfg: cfg,
	}
}

func (d *DummyScraper) Run(ctx context.Context) ([][]string, error) {
	return [][]string{}, nil
}
