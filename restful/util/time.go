package util

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetTimeStd(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetDateStd(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetUnix() string {
	now := time.Now().Unix()
	return strconv.FormatInt(now, 10)
}

func GetUnixNano() string {
	now := time.Now().UnixNano()
	return strconv.FormatInt(now, 10)
}

func TimeNow() time.Time {
	t := time.Now()
	_, offset := t.Zone()
	if offset != 28800 {
		t = t.Add(time.Hour * 8)
	}
	return t
}

func TimeNowStd() string {
	return TimeNow().Format("2006-01-02 15:04:05")
}

func GetTimeLocal(t time.Time) time.Time {
	_, offset := t.Zone()
	if offset != 28800 {
		t = t.Add(time.Hour * 8)
	}
	return t
}

func GetTimeLocalStd(t time.Time) string {
	_, offset := t.Zone()
	if offset != 28800 {
		t = t.Add(time.Hour * 8)
	}
	return t.Format("2006-01-02 15:04:05")
}

// sql支持的时间类型
// SqlTime format json time field by myself
type SqlTime struct {
	time.Time
}

func NewSqlTime(t time.Time) SqlTime {
	return SqlTime{
		Time: t,
	}
}

// MarshalJSON on SqlTime format Time field with %Y-%m-%d %H:%M:%S
func (t SqlTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	nullTime := fmt.Sprintf("\"%s\"", `0001-01-01 00:00:00`)

	if formatted == nullTime {
		formatted = "\"\""
	}
	return []byte(formatted), nil
}

func (t *SqlTime) UnmarshalJSON(v []byte) error {
	timeStr := strings.Replace(string(v), "\"", "", -1)
	if len(timeStr) == 0 || timeStr == "" {
		timeStr = "0001-01-01 00:00:00"
	}

	value, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err != nil {
		return err
	}
	*t = SqlTime{Time: value}
	return nil
}

// Value insert timestamp into mysql need this function.
func (t SqlTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *SqlTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = SqlTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
