package store

import (
	"fmt"
)

func buildKey(args ...interface{}) string {
	var key string
	for a := range args {
		if key != "" {
			key = fmt.Sprintf("%s/%v", key, args[a])
		} else {
			key = fmt.Sprintf("%v", args[a])
		}
	}
	return key
}
