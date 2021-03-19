package store

import (
	"context"
	"io/ioutil"
	"os"
	"path"
)

type Local struct {
	prefix string
}

func NewLocal(prefix string) Store {
	return &Local{prefix: prefix}
}

func (l *Local) Setup(ctx context.Context) error {
	err := os.MkdirAll(l.prefix, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s *Local) Persist(ctx context.Context, filename string, data string) error {
	dataAsByte := []byte(data)

	filenameWithPrefix := path.Join(s.prefix, filename)

	dirName := path.Dir(filenameWithPrefix)

	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filenameWithPrefix, dataAsByte, 0644)
	if err != nil {
		return err
	}

	return nil
}
