package bytes

import (
	"math"
	"testing"
)

func TestByte(t *testing.T) {
	tests := []struct {
		arg         int64
		wantDecimal string
		wantBinary  string
	}{
		{arg: 0, wantDecimal: "0 B", wantBinary: "0 B"},
		{arg: 27, wantDecimal: "27 B", wantBinary: "27 B"},
		{arg: 1000, wantDecimal: "1.0 kB", wantBinary: "1000 B"},
		{arg: 1024, wantDecimal: "1.0 kB", wantBinary: "1.0 KiB"},
		{arg: 31589105, wantDecimal: "31.6 MB", wantBinary: "30.1 MiB"},
		{arg: 1855425871872, wantDecimal: "1.9 TB", wantBinary: "1.7 TiB"},
		{arg: math.MaxInt64, wantDecimal: "9.2 EB", wantBinary: "8.0 EiB"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := ByteDecimal(tt.arg); got != tt.wantDecimal {
				t.Errorf("ByteDecimal() = %v, want %v", got, tt.wantDecimal)
			}

			if got := ByteBinary(tt.arg); got != tt.wantBinary {
				t.Errorf("ByteBinary() = %v, want %v", got, tt.wantBinary)
			}
		})
	}
}
