package util

import (
	"fmt"
	"testing"
)

func TestNow(t *testing.T) {
	fmt.Println(TimeNow())
	fmt.Println(TimeNowStd())
}
