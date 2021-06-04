package jobs

import (
	"context"
	"testing"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
)

func Test(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	manager := trieugene.NewDummyManager(cfg)
	expectedJob := trieugene.ExpectsJob{}
	dummyScraper := scraper.NewDummyScraper(cfg)
	parser := scraper.NewParser(cfg, dummyScraper)

	csvJob := NewCsvJob(&CsvJobKwargs{
		Cfg:      cfg,
		Manager:  manager,
		StoreJob: &expectedJob,
		Parser:   parser,
		Scraper:  dummyScraper,
	})

	csvJob.Run(ctx)
}
