package main

import (
	"gitee_cli/cmd"
	_ "gitee_cli/config"
	runewidth "github.com/mattn/go-runewidth"
	"os"
)

func init() {
	os.Setenv("RUNEWIDTH_EASTASIAN", "0")
	runewidth.DefaultCondition.EastAsianWidth = false
}

func main() {
	cmd.Execute()
}
