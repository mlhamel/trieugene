package store

import (
	"fmt"
	"time"
)

func buildKey(timestamp time.Time, key string) string {
	return fmt.Sprintf("%s/%d", key, timestamp.Unix())
}
