package util

var id int64 = 0

func GenID() int64 {
	return id + 1
}
