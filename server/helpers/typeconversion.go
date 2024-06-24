package helper

import "strconv"

func Int64ToStirng(i int64) string {
	return strconv.Itoa(int(i))
}

func IntToStirng(i int) string {
	return strconv.Itoa(i)
}

func StirngToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
