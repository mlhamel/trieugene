package store

import (
	"context"
	"encoding/json"
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
	s.cfg.Logger().Debug().Str("bucket", s.cfg.S3Bucket()).Msg("Starting Creating bucket")
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
	s.cfg.Logger().Debug().Str("bucket", s.cfg.S3Bucket()).Msg("Succeed Creating bucket")
	return nil
}

func (s *S3) Persist(ctx context.Context, data *Data) error {
	datetime := time.Unix(data.Timestamp, 0)
	key := fmt.Sprintf("%s/%s/%s.json", data.Name, datetime.Format("20060102"), datetime.Format("1504"))
	body, err := json.Marshal(data)

	if err != nil {
		s.cfg.Logger().Error().Err(err).Str("key", key).Msg("Failed marshaling data for persistence")
		return err
	}

	bodyStr := string(body)

	s.cfg.Logger().Debug().Str("key", key).Str("body", bodyStr).Msg("Starting Persistence")
	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket()),
		Key:    aws.String(key),
		Body:   strings.NewReader(bodyStr),
	})
	if err != nil {
		s.cfg.Logger().Error().Err(err).Str("key", key).Msg("Failed persistence")
		return err
	}

	s.cfg.Logger().Debug().Str("key", key).Msg("Succeed Persistence")

	return nil
}
