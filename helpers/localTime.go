package helpers

import "time"

func GetLocalTime() *time.Location {
	location, err := time.LoadLocation("Asia/Jakarta")
	ErrorHandler(err)

	return location
}

func GetLocalDateNow() time.Time {
	location := GetLocalTime()
	return time.Now().In(location).Truncate(24 * time.Hour)
}
