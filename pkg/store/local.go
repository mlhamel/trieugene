package store

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/mlhamel/trieugene/pkg/config"
)

type Local struct {
	cfg    *config.Config
	prefix string
}

func NewLocal(cfg *config.Config) Store {
	return &Local{cfg: cfg}
}

func (l *Local) Setup(ctx context.Context) error {
	err := os.MkdirAll(l.cfg.LocalPrefix(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (l *Local) Persist(ctx context.Context, filename string, data string) error {
	dataAsByte := []byte(data)

	filenameWithPrefix := path.Join(l.cfg.LocalPrefix(), filename)

	dirName := path.Dir(filenameWithPrefix)

	l.cfg.Logger().Debug().Str("dirname", dirName).Msg("Creating directory")
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	l.cfg.Logger().Debug().Str("filenameWithPrefix", dirName).Msg("Persisting file")
	err = ioutil.WriteFile(filenameWithPrefix, dataAsByte, 0644)
	if err != nil {
		return err
	}

	return nil
}
