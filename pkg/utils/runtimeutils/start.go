package runtimeutils

import (
	"strconv"
	"time"
)

//nolint:gochecknoglobals
var (
	startTime = time.Now().Unix()

	startTimeString = strconv.FormatInt(startTime, 10)
)

func GetStartTime() int64 {
	return startTime
}

func GetStartTimeString() string {
	return startTimeString
}
