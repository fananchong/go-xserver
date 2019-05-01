package utils

import "time"

// GetMillisecondTimestamp : 当前的毫秒时间戳
func GetMillisecondTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
