package scraper

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mlhamel/trieugene/pkg/config"
)

const url = "https://www.cehq.gouv.qc.ca/suivihydro/fichier_donnees.asp?NoStation=040204"

type Scraper struct {
	cfg      *config.Config
	client   *http.Client
	consumer func(context.Context, Result) error
}

type Result struct {
	ScrapedAt time.Time
	TakenAt   time.Time
	Outflow   float64
}

func (r *Result) Sha1() string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(fmt.Sprintf("%d:%f", r.TakenAt, r.Outflow)))
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (r *Result) MD5() string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, fmt.Sprintf("%d:%f", r.TakenAt, r.Outflow))
	if err != nil {
		return ""
	}
	return string(hasher.Sum(nil))
}

func NewScraper(cfg *config.Config, consumer func(context.Context, Result) error) *Scraper {
	return &Scraper{cfg: cfg, client: &http.Client{}, consumer: consumer}
}

func (s *Scraper) Run(ctx context.Context) error {
	response, err := s.client.Get(url)

	if err != nil {
		return fmt.Errorf("Error while downloading: %w", err)
	}

	s.cfg.Logger().Info().
		Int("statusCode", response.StatusCode).
		Int64("ContentLength", response.ContentLength).
		Msgf("Response from %s", url)

	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()

	var record Result

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
			s.cfg.Logger().Error().Err(err).Msg("Error parsing date")
			continue
		}

		rawFlow := strings.Replace(each[2], ",", ".", 1)
		rawFlow = strings.Replace(rawFlow, "*", "", 1)

		outflow, err := strconv.ParseFloat(rawFlow, 64)

		if err != nil {
			s.cfg.Logger().Error().Err(err).Msg("Error parsing flow")
			continue
		}

		record.ScrapedAt = time.Now().UTC()
		record.TakenAt = takenAt.UTC()
		record.Outflow = outflow

		s.cfg.Logger().
			Info().
			Time("ScrapedAt", record.ScrapedAt).
			Time("TakenAt", record.TakenAt).
			Float64("Outflow", record.Outflow).
			Msg("Record parsed")

		if err = s.consumer(ctx, record); err != nil {
			return err
		}
	}

	return nil
}
