package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

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

func CheckCPUMon(path string) (string, error) {
	data, err := os.ReadFile(filepath.Join(path, "name"))
	if err != nil {
		return "", err
	}
	content := string(data)
	switch content {
	case "coretemp\n", "k10temp\n":
		return path, nil
	default:
		err = errors.New("no info about temperature file")
		return "", err
	}
}
