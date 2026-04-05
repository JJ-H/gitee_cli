package pull_request

import (
	"gitee_cli/internal/api/pull_request"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var CloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := pull_request.Close(args[0]); err != nil {
			color.Red(err.Error())
			return
		}
		color.Green("关闭 PR 成功🏅")
	},
}
