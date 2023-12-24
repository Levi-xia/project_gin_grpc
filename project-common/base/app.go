package base

import (
	"os"
	"strings"
)

func GetCurApp() string {
	workDir, _ := os.Getwd()
	arr := strings.Split(workDir, "/")
	app := arr[len(arr)-1]
	return app
}