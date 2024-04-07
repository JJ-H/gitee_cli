package issue

import (
	"gitee_cli/config"
	"gitee_cli/internal/api/enterprises"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var enterprise enterprises.Enterprise

var IssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		entPath, err := cmd.Flags().GetString("ent")
		if entPath == "" {
			if entPath = config.Conf.DefaultEntPath; entPath == "" {
				color.Red("请指定企业 path")
				return
			}
		}
		enterprise, err = enterprises.Find(entPath)
		if err != nil {
			color.Red("企业未找到！")
			return
		}
	},
}

func init() {
	IssueCmd.AddCommand(CreateCmd)
	IssueCmd.AddCommand(ListCmd)
	IssueCmd.AddCommand(ViewCmd)
	IssueCmd.PersistentFlags().StringP("ent", "e", "", "specify the selector_tui path")
}
