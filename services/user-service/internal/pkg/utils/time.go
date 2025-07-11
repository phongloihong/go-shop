package utils

import "time"

func TimeNow() int64 {
	return time.Now().Unix()
}
