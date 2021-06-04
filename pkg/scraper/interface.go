package scraper

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

type Parser interface {
	Run(ctx context.Context, consumer Consumer) error
}

type Scraper interface {
	Run(context.Context) ([][]string, error)
}

type Consumer func(context.Context, Result) error

type Result struct {
	ScrapedAt time.Time
	TakenAt   time.Time
	Outflow   interface{}
}

func (r *Result) Sha1() string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(fmt.Sprintf("%d:%f", r.TakenAt.Unix(), r.Outflow)))
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (r *Result) MD5() string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, fmt.Sprintf("%d:%f", r.TakenAt.Unix(), r.Outflow))
	if err != nil {
		return ""
	}
	return string(hasher.Sum(nil))
}
