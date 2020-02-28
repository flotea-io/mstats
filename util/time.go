package util

import "time"

func CurrentTime() string {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	return ts
}

func TimeFromString(timeString string) time.Time {
	parse, _ := time.Parse("2006-01-02 15:04:05", timeString)
	return parse
}

func CurrentTimeWithoutSeconds() string {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04")
	return ts
}
