package util

import "time"

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
