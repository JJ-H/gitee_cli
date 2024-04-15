package pull_request

import (
	"github.com/spf13/cobra"
)

var Pr = &cobra.Command{
	Use:     "pr",
	Aliases: []string{"pull_request"},
	Short:   "Manage pull requests",
}

func init() {
	Pr.AddCommand(ListCmd)
	Pr.AddCommand(CreateCmd)
	Pr.AddCommand(CommentCmd)
	Pr.AddCommand(CloseCmd)
	Pr.AddCommand(ReviewCmd)
}
