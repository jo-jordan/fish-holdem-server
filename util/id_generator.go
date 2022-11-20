package util

import "sync/atomic"

var id int64 = 0

func GenID() int64 {
	return atomic.AddInt64(&id, 1)
}
