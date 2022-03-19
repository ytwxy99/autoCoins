package utils

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// get now timestamp
func GetNowTimeStamp() int64 {
	return time.Now().Unix()
}

// get now time.time
func GetNowTime() time.Time {
	nowTimstamp := time.Now().Unix()
	return time.Unix(nowTimstamp, 0)
}

// get now data
func GetNowData() string {
	nowTimstamp := time.Now().Unix()
	tm := time.Unix(nowTimstamp, 0)

	return tm.Format("2006-01-02 15:04:05")
}

func GetData(timstamp int64) string {
	tm := time.Unix(timstamp, 0)

	return tm.Format("2006-01-02 15:04:05")
}

// get old timestap
func GetOldTimeStamp(years int, months int, days int) int64 {
	return time.Now().AddDate(years, months, days).Unix()
}

// get old data
func GetOldData(years int, months int, days int) string {
	return time.Now().AddDate(years, months, days).String()
}

// string convert to float32
func StringToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		logrus.Error("string convert to float32 error: %v", s)
		return 0
	}

	return float32(f)
}

// float32 convert to string
func Float32ToString(f float32) string {
	return strconv.FormatFloat(float64(f), 'g', 10, 64)
}
