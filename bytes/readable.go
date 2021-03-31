package bytes

import "fmt"

func ByteDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ByteBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

/*
Input     Decimal (SI)     Binary (IEC)
               0            "0 B"            "0 B"
              27           "27 B"           "27 B"
             999          "999 B"          "999 B"
            1000         "1.0 kB"         "1000 B"
            1023         "1.0 kB"         "1023 B"
            1024         "1.0 kB"        "1.0 KiB"
            1728         "1.7 kB"        "1.7 KiB"
   1855425871872         "1.9 TB"        "1.7 TiB"
   math.MaxInt64         "9.2 EB"        "8.0 EiB
*/