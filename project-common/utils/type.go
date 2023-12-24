package utils

import "strconv"

// uint转string
func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

// uint转int64
func UintToInt64(u uint) int64 {
	return int64(u)
}

// string转int64
func StringToInt64(s string) (int64, error){
	i, err :=  strconv.ParseInt(s, 10, 64)
	return i, err
}