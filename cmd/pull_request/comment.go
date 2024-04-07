package pull_request

import (
	"fmt"
	"gitee_cli/internal/api/pull_request"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var CommentCmd = &cobra.Command{
	Use:     "comment",
	Short:   "Comment pull request",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"note"},
	Run: func(cmd *cobra.Command, args []string) {
		comment := args[0]
		iid, _ := cmd.Flags().GetInt("iid")

		if iid == 0 {
			color.Red("请给定有效的 pull request 序号！")
			os.Exit(1)
		}
		if err := pull_request.Note(iid, comment); err != nil {
			fmt.Println(err.Error())
			return
		}
		color.Green("评论成功！")
	},
}

func init() {
	CommentCmd.Flags().IntP("iid", "i", 0, "Pull request number")
}
