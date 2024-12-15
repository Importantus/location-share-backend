package utils

import (
	"strconv"
	"time"
)

func ParseTime(timestamp string) (time.Time, error) {
	parsedTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime := time.Unix(parsedTimestamp, 0)

	return parsedTime, err
}
