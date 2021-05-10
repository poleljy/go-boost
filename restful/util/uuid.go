package util

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}
