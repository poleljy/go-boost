package util

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	uid := uuid.NewV4()
	return strings.Replace(uid.String(), "-", "", -1)
}

func UUIDShort() string {
	uuid := UUID()
	return uuid[:10]
}
