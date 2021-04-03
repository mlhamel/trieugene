package scraper

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/mlhamel/trieugene/pkg/config"
)

const url = "https://www.cehq.gouv.qc.ca/suivihydro/fichier_donnees.asp?NoStation=040204"

type HttpScraper struct {
	cfg    *config.Config
	client *http.Client
}

func NewHttpScraper(cfg *config.Config) *HttpScraper {
	return &HttpScraper{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (h *HttpScraper) Run(ctx context.Context) ([][]string, error) {
	response, err := h.client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error while downloading: %w", err)
	}

	h.cfg.Logger().Info().
		Int("statusCode", response.StatusCode).
		Int64("ContentLength", response.ContentLength).
		Msgf("Response from %s", url)

	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	return reader.ReadAll()
}
