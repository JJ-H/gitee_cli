package pull_request

import (
	"gitee_cli/internal/api/pull_request"
	"github.com/spf13/cobra"
)

var ReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		iid := args[0]
		pull_request.Review(iid)
	},
}
