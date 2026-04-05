package pull_request

import (
	"gitee_cli/internal/api/pull_request"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := pull_request.Review(args[0]); err != nil {
			color.Red(err.Error())
			return
		}
		color.Green("审查通过🏅")
	},
}
