package unix

import (
	"time"
)

const format = time.UnixDate

func NowUTC() string {
	return time.Now().UTC().Format(format)
}

func ToTimeUTC(s string) (time.Time, error) {
	return time.Parse(format, s)
}

func ToStringUTC(t time.Time) string {
	return t.Format(format)
}
