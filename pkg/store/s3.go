package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

func NewS3Production(cfg *config.Config) Store {
	conf := aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKey(), cfg.S3SecretKey(), ""),
		Endpoint:    aws.String(cfg.S3URL()),
		Region:      aws.String(cfg.S3Region()),
	}

	client := s3.New(session.New(&conf))

	return &S3{cfg: cfg, client: client}
}

func (s *S3) Setup(ctx context.Context) error {
	s.cfg.Logger().Debug().Msgf("Creating bucket %s", s.cfg.S3Bucket())
	input := &s3.CreateBucketInput{Bucket: aws.String(s.cfg.S3Bucket())}
	_, err := s.client.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				err = nil
			default:
				return aerr
			}
		} else {
			return err
		}
	}
	s.cfg.Logger().Debug().Msgf("Creating bucket %s: Succeed", s.cfg.S3Bucket())
	return nil
}

func (s *S3) Persist(ctx context.Context, timestamp int64, name string, id string, data interface{}) error {
	ctx = context.Background()
	datetime := time.Unix(timestamp, 0)
	key := fmt.Sprintf("%s/%s/%s.json", name, datetime.Format("200601021504"), id)

	s.cfg.Logger().Debug().Msgf("Persisting %s", key)
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket()),
		Key:    aws.String(key),
		Body:   strings.NewReader(fmt.Sprintf("%v", data)),
	})
	if err != nil {
		s.cfg.Logger().Error().Err(err).Msgf("Persisting %s: Failed", key)
		return err
	}

	s.cfg.Logger().Debug().Msgf("Persisting %s: Succeed", key)

	return nil
}
