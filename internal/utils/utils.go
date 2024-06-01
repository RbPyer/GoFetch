package utils

import "fmt"


func UtsnameToStr(u [65]int8) string {
	result := make([]byte, 65)
	for i := 0; i < 65 && u[i] != 0; i++ {
		result = append(result, byte(u[i]))
	}
	return string(result)
}


func Int64ToTimeStr(num int64) string {
	var days, hours, minutes, seconds int64
	days = num / 86400
	hours = num % 86400 / 3600
	minutes = num % 3600 / 60
	seconds = num % 60

	return fmt.Sprintf("%d days %d hours %d minutes %d seconds", days, hours, minutes, seconds)
}