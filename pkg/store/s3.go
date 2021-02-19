package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mlhamel/trieugene/pkg/config"
)

type S3 struct {
	cfg    *config.Config
	client *s3.S3
}

const shortDuration = 100 * time.Millisecond

func NewS3(cfg *config.Config) Store {
	conf := aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.S3AccessKey(), cfg.S3SecretKey(), ""),
		Endpoint:         aws.String(cfg.S3URL()),
		Region:           aws.String(cfg.S3Region()),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	client := s3.New(session.New(&conf))

	return &S3{cfg: cfg, client: client}
}

func (s *S3) Persist(ctx context.Context, timestamp time.Time, name string, id string, data interface{}) error {
	ctx = context.Background()
	s.cfg.Logger().Debug().Msgf("Persisting %s/%s: %s", name, id, data)
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket()),
		Key:    aws.String(buildKey(name, timestamp.Unix(), id)),
		Body:   strings.NewReader(fmt.Sprintf("%v", data)),
	})
	if err != nil {
		s.cfg.Logger().Error().Err(err).Msgf("Persisting %s/%s: Failed", name, id)
		return err
	}

	s.cfg.Logger().Debug().Msgf("Persisting %s/%s: Succeed", name, id)

	return nil
}
