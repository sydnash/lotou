package utils

import "time"

func Now() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}
