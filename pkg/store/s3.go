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
)

type S3Params struct {
	AccessKey        string
	SecretKey        string
	URL              string
	Bucket           string
	Region           string
	DisableSSL       bool
	S3ForcePathStyle bool
}

type S3 struct {
	params *S3Params
	client *s3.S3
}

const shortDuration = 100 * time.Millisecond

func NewS3(params *S3Params) Store {
	conf := aws.Config{
		Credentials:      credentials.NewStaticCredentials(params.AccessKey, params.SecretKey, ""),
		Endpoint:         aws.String(params.URL),
		Region:           aws.String(params.Region),
		DisableSSL:       aws.Bool(params.DisableSSL),
		S3ForcePathStyle: aws.Bool(params.S3ForcePathStyle),
	}

	client := s3.New(session.New(&conf))

	return &S3{client: client, params: params}
}

func (s *S3) Setup(ctx context.Context) error {
	input := &s3.CreateBucketInput{Bucket: aws.String(s.params.Bucket)}
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
	return nil
}

func (s *S3) Persist(ctx context.Context, data *Data) error {
	datetime := time.Unix(data.Timestamp, 0)
	key := fmt.Sprintf("%s/%s/%s.json", data.Name, datetime.Format("20060102"), datetime.Format("1504"))
	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	bodyStr := string(body)

	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.params.Bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(bodyStr),
	})
	if err != nil {
		return err
	}

	return nil
}
