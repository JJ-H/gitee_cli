package pull_request

import (
	"fmt"
	"gitee_cli/internal/api/pull_request"
	"gitee_cli/utils"
	"gitee_cli/utils/git_utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strings"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request",
	Long:  "Create a pull request",
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		base, _ := cmd.Flags().GetString("base")
		head, _ := cmd.Flags().GetString("head")
		skipBody, _ := cmd.Flags().GetBool("skip-body")
		draft, _ := cmd.Flags().GetBool("draft")
		assignee, _ := cmd.Flags().GetString("assignees")
		tester, _ := cmd.Flags().GetString("testers")
		prune, _ := cmd.Flags().GetBool("prune")
		if !git_utils.IsGitDir() {
			color.Red("请在仓库目录下执行该命令！")
			return
		}

		baseRepo, err := git_utils.ParseCurrentRepo()
		if err != nil {
			color.Red("获取当前仓库异常！")
			return
		}

		if title == "" {
			title = utils.ReadFromInput("请输入标题", title)
			title = strings.TrimSpace(title)
		}

		if head == "" {
			branch, _ := git_utils.GetCurrentBranch()
			if branch == "" {
				branch = utils.ReadFromInput("请输入起始分支", branch)
				head = strings.TrimSpace(branch)
				if head == "" {
					color.Red("无效的起始分支！")
					return
				}
			}
			head = branch
		}

		if base == "" {
			base = utils.ReadFromInput("请输入目标分支", base)
			base = strings.TrimSpace(base)
			if base == "" {
				color.Red("无效的目标分支！")
			}
		}

		if body == "" && !skipBody {
			body = utils.ReadFromEditor(utils.InitialEditor("填写 Pull Request 内容", ""), body)
		}

		pullRequest, err := pull_request.CreatePr(baseRepo, base, head, title, body, assignee, tester, draft, prune)

		if err != nil {
			color.Red(err.Error())
			var input string
			input = utils.ReadFromInput(utils.Yellow("是否重试？(y/n/q)"), input)
			if input == "y" || input == "yes" {
				pullRequest, err = pull_request.CreatePr(baseRepo, base, head, title, body, assignee, tester, draft, prune)
				if err != nil {
					color.Red(err.Error())
					return
				}
			} else {
				return
			}
		}
		fmt.Printf("创建 PR「%s」 成功，访问地址：%s\n", utils.Yellow(pullRequest.Title), utils.Cyan(pullRequest.HtmlUrl))
	},
}

func init() {
	CreateCmd.Flags().StringP("title", "t", "", "Title of the pull request")
	CreateCmd.Flags().StringP("body", "b", "", "Body of the pull request")
	CreateCmd.Flags().StringP("base", "B", "", "The branch into which you want your code merged")
	CreateCmd.Flags().StringP("head", "H", "", "The branch that contains commits for your pull request (default [current branch])")
	CreateCmd.Flags().BoolP("skip-body", "", false, "Skip adding a body to the pull request")
	CreateCmd.Flags().BoolP("draft", "", false, "Create a draft pull request")
	CreateCmd.Flags().StringP("assignees", "a", "", "Assign the pull request to users to code review, multi user split by , user1,user2")
	CreateCmd.Flags().StringP("testers", "", "", "Assign the pull request to users to test, multi user split by , user1,user2")
	CreateCmd.Flags().BoolP("prune", "", true, "Prune source branch after pr merged")
}
