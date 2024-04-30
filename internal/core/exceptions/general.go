package exceptions

import (
	"fmt"
)

var ErrNoRecords error

func NotFoundErr(model string) error {
	ErrNoRecords = fmt.Errorf("%s not found", model)

	return ErrNoRecords
}
