package pull_request

import (
	"fmt"
	"gitee_cli/config"
	"gitee_cli/internal/api/enterprises"
	"gitee_cli/internal/api/pull_request"
	"gitee_cli/utils"
	"gitee_cli/utils/git_utils"
	"gitee_cli/utils/tui"
	"github.com/charmbracelet/bubbles/table"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests related",
	Run: func(cmd *cobra.Command, args []string) {
		keyword, _ := cmd.Flags().GetString("keyword")
		scope, _ := cmd.Flags().GetString("scope")
		commitSha, _ := cmd.Flags().GetString("commit")
		openInBrowser, _ := cmd.Flags().GetBool("open")
		convertEntUrl, _ := cmd.Flags().GetBool("convert")
		if commitSha != "" {
			pathWithNamespace, err := git_utils.ParseCurrentRepo()
			if err != nil {
				color.Red(err.Error())
				return
			}
			pr, err := pull_request.FindPullRequestByIid(commitSha, pathWithNamespace)
			if err != nil {
				color.Red(err.Error())
				return
			}
			if openInBrowser {
				browser.OpenURL(pr.HtmlUrl)
				return
			}
			fmt.Printf("该 commit 由 PR: 「%v」 合入，访问地址: %s\n", utils.Green(pr.Title), utils.Blue(pr.HtmlUrl))
			return
		}
		pullRequests := pull_request.List(scope)
		if keyword != "" {
			pullRequests = pull_request.FuzzySearch(pullRequests, keyword)
		}
		if len(pullRequests) == 0 {
			color.Cyan("无匹配的 Pull requests.")
			return
		}

		columns := []table.Column{
			{Title: "PR 标题", Width: 60},
			{Title: "IID", Width: 10},
			{Title: "创建人", Width: 12},
			{Title: fmt.Sprintf("%s审查状态", config.Conf.UserName), Width: 18},
			{Title: "冲突", Width: 10},
			{Title: "是否可合入", Width: 10},
		}
		rows := make([]table.Row, 0)
		for _, pr := range pullRequests {
			mergeCheckMsg := "No"
			if pr.CanMergeCheck {
				mergeCheckMsg = utils.Green("Yes")
			}
			conflictMsg := "No"
			if !pr.Mergeable {
				conflictMsg = utils.Magenta("Yes")
			}
			acceptMsg := utils.Magenta("未审查")
			if pr.User.Id == 0 {
				acceptMsg = "-"
			} else if pr.User.Accept {
				acceptMsg = utils.Green("通过")
			}
			rows = append(rows, table.Row{pr.Title, strconv.Itoa(pr.Number), pr.Creator.Name, acceptMsg, conflictMsg, mergeCheckMsg})
		}

		prTable := tui.NewTable(enterprises.Enterprise{}, tui.PullRequest, columns, rows)

		if convertEntUrl {
			os.Setenv("CONVERT_ENT_URL", "true")
		}

		defer func() {
			os.Unsetenv("CONVERT_ENT_URL")
		}()
		//var model tea.Model
		var err error
		if _, err = prTable.Run(); err != nil {
			color.Red("pr渲染失败！")
			os.Exit(1)
		}
		//utils.PrRender(pullRequests, convert)
	},
}

func init() {
	ListCmd.Flags().StringP("keyword", "k", "", "filter pr by keyword")
	//ListCmd.Flags().BoolP("reviewed", "r", false, "filter pr by review state")
	ListCmd.Flags().StringP("scope", "s", "", "filter pr by scope (owner)")
	ListCmd.Flags().StringP("commit", "c", "", "find pr by commit")
	ListCmd.Flags().BoolP("open", "o", false, "open in browser, only effective for searching pr via commit sha")
	ListCmd.Flags().BoolP("convert", "", false, "transfer url in enterprise")
}
