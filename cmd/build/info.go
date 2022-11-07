package build

import "time"

var (
	version  = time.Now().Format(time.RFC3339)
	commit   = ""
	date     = ""
	dateTime time.Time
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func Date() time.Time {
	if !dateTime.IsZero() {
		return dateTime
	}
	dateTime, _ = time.Parse(time.RFC3339, date)
	if dateTime.IsZero() {
		dateTime = time.Now()
	}
	return dateTime
}
