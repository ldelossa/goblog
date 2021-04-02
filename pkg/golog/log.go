package golog

import (
	"os"

	"github.com/fatih/color"
)

func Error(format string, args ...interface{}) {
	color.Red(format, args...)
}

func Fatal(format string, args ...interface{}) {
	color.Red(format, args...)
	os.Exit(1)
}

func Info(format string, args ...interface{}) {
	color.Blue(format, args...)
}

func Warning(format string, args ...interface{}) {
	color.Yellow(format, args...)
}
