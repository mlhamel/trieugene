package scraper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/scraper"
)

type parserImpl struct {
	cfg     *config.Config
	scraper scraper.Scraper
}

func NewParser(cfg *config.Config, scraper scraper.Scraper) scraper.Parser {
	return &parserImpl{
		cfg:     cfg,
		scraper: scraper,
	}
}

func (p *parserImpl) Run(ctx context.Context, consumer scraper.Consumer) error {
	data, err := p.scraper.Run(ctx)

	if err != nil {
		return err
	}

	var record scraper.Result

	if err != nil {
		return fmt.Errorf("Error while parsing: %w", err)
	}

	for _, each := range data {
		if strings.TrimSpace(each[0]) == "Date" {
			continue
		}

		if strings.TrimSpace(each[0]) == "" {
			continue
		}

		rawDate := fmt.Sprintf("%s %s", each[0], strings.TrimRight(each[1], "\\"))

		takenAt, err := time.Parse("2006-01-02 15:04", rawDate)

		if err != nil {
			p.cfg.Logger().Error().Err(err).Msg("Error parsing date")
			continue
		}

		rawFlow := strings.Replace(each[2], ",", ".", 1)
		rawFlow = strings.Replace(rawFlow, "*", "", 1)

		outflow, err := strconv.ParseFloat(rawFlow, 64)

		if err != nil {
			p.cfg.Logger().Error().Err(err).Msg("Error parsing flow")
			continue
		}

		record.ScrapedAt = time.Now().UTC()
		record.TakenAt = takenAt.UTC()
		record.Outflow = outflow

		p.cfg.Logger().
			Info().
			Time("ScrapedAt", record.ScrapedAt).
			Time("TakenAt", record.TakenAt).
			Interface("Outflow", record.Outflow).
			Msg("Record parsed")

		if err = consumer(ctx, record); err != nil {
			return err
		}
	}

	return nil
}
