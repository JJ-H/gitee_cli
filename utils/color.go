package utils

import "github.com/fatih/color"

var (
	green   func(a ...interface{}) string = color.New(color.FgGreen).SprintFunc()
	blue    func(a ...interface{}) string = color.New(color.FgBlue).SprintFunc()
	yellow  func(a ...interface{}) string = color.New(color.FgYellow).SprintFunc()
	cyan    func(a ...interface{}) string = color.New(color.FgCyan).SprintFunc()
	red     func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()
	magenta func(a ...interface{}) string = color.New(color.FgMagenta).SprintFunc()
)

func Green(a ...interface{}) string {
	return green(a...)
}

func Blue(a ...interface{}) string {
	return blue(a...)
}

func Yellow(a ...interface{}) string {
	return yellow(a...)
}

func Cyan(a ...interface{}) string {
	return cyan(a...)
}

func Red(a ...interface{}) string {
	return red(a...)
}

func Magenta(a ...interface{}) string {
	return magenta(a...)
}
