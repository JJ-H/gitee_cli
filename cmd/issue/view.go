package issue

import (
	"gitee_cli/internal/api/issue"
	"gitee_cli/utils/tui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display issue detail",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ident := args[0]

		if _issue, err := issue.Detail(enterprise.Id, ident); err == nil {
			tui.NewPager(_issue.Title, _issue.Description, tui.Markdown).Run()
		} else {
			color.Red("获取任务详情失败！")
			return
		}
	},
}
