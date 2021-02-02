package jobs

import (
	"errors"
	"fmt"
)

var ErrInvalidMsg = errors.Unwrap(fmt.Errorf("Invalid args for message processing job"))
var ErrInvalidData = errors.Unwrap(fmt.Errorf("Invalid args for data processing job"))
