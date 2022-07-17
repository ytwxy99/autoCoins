package utils

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// GetNowTimeStamp get now timestamp
func GetNowTimeStamp() int64 {
	return time.Now().Unix()
}

// GetNowTime get now time.time
func GetNowTime() time.Time {
	nowTimstamp := time.Now().Unix()
	return time.Unix(nowTimstamp, 0)
}

// GetNowData get now data
func GetNowData() string {
	nowTimstamp := time.Now().Unix()
	tm := time.Unix(nowTimstamp, 0)

	return tm.Format("2006-01-02 15:04:05")
}

func GetData(timstamp int64) string {
	tm := time.Unix(timstamp, 0)

	return tm.Format("2006-01-02 15:04:05")
}

// GetOldTimeStamp get old timestap
func GetOldTimeStamp(years int, months int, days int) int64 {
	return time.Now().AddDate(years, months, days).Unix()
}

// GetOldData get old data
func GetOldData(years int, months int, days int) string {
	return time.Now().AddDate(years, months, days).String()
}

// StringToFloat32 string convert to float32
func StringToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		logrus.Error("string convert to float32 error: %v", s)
		return 0
	}

	return float32(f)
}

// StringToFloat64 string convert to float64
func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logrus.Error("string convert to float32 error: %v", s)
		return 0
	}

	return float64(f)
}

// Float32ToString float32 convert to string
func Float32ToString(f float32) string {
	return strconv.FormatFloat(float64(f), 'g', 10, 64)
}

func CurrentHourTimestamp() int64 {
	return time.Now().Unix() - int64(time.Now().Second()) - int64(60*time.Now().Minute())
}

func NextHourTimestamp() int64 {
	addDuration, _ := time.ParseDuration("1h")

	return time.Unix(CurrentHourTimestamp(), 0).Add(addDuration).Unix()
}

func CurrentHalfHourTimestamp() int64 {
	addDuration, _ := time.ParseDuration("30m")
	return time.Unix(CurrentHourTimestamp(), 0).Add(addDuration).Unix()
}

func IsTradeTime() bool {
	currenTimeStamp := time.Now().Unix()
	before1mDuration, _ := time.ParseDuration("-1m")
	beforeHalfTime := time.Unix(CurrentHalfHourTimestamp(), 0).Add(before1mDuration).Unix()

	if currenTimeStamp >= beforeHalfTime && currenTimeStamp <= CurrentHalfHourTimestamp() {
		return true
	}

	beforeNextHourTime := time.Unix(NextHourTimestamp(), 0).Add(before1mDuration).Unix()

	if currenTimeStamp >= beforeNextHourTime && currenTimeStamp <= NextHourTimestamp() {
		return true
	}

	return false
}
