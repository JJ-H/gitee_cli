package issue

import (
	"fmt"
	"gitee_cli/config"
	"gitee_cli/internal/api/issue"
	"gitee_cli/internal/api/issue_type"
	"gitee_cli/internal/api/member"
	"gitee_cli/utils"
	"gitee_cli/utils/tui"
	"gitee_cli/utils/tui/issue_type_tui"
	"gitee_cli/utils/tui/selector_tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strings"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a issue",
	Run: func(cmd *cobra.Command, args []string) {
		isBug, _ := cmd.Flags().GetBool("bug")
		isRequirement, _ := cmd.Flags().GetBool("feature")
		skipBody, _ := cmd.Flags().GetBool("skip-body")
		entPath, _ := cmd.Flags().GetString("ent")
		parentKeyWord, _ := cmd.Flags().GetString("parent")
		assigneeKeyWord, _ := cmd.Flags().GetString("assignee")
		candidateAssignees := make([]member.Member, 0)
		candidateTasks := make([]issue.Issue, 0)
		assigneeId := 0
		parentTaskId := 0

		if entPath == "" {
			if entPath = config.Conf.DefaultEntPath; entPath == "" {
				color.Red("请指定企业 path")
				return
			}
		}

		if parentKeyWord != "" {
			candidateTasks, _ = issue.Find(enterprise.Id, map[string]string{
				"search": parentKeyWord,
			})
		}

		if assigneeKeyWord != "" {
			candidateAssignees, _ = member.Find(enterprise.Id, map[string]string{
				"search": assigneeKeyWord,
			})
		}

		optionMap := make(map[string]int, 0)
		options := make([]string, 0)

		var category = issue_type.TASK
		var categoryText = "请选择任务类型"
		if isBug {
			category = issue_type.BUG
			categoryText = "请选择缺陷类型"
		} else if isRequirement {
			category = issue_type.REQUIREMENT
			categoryText = "请选择需求类型"
		}

		issueTypes, err := issue_type.List(category, entPath)
		if err != nil {
			color.Red("获取任务类型失败！")
			return
		}

		// 填充选项
		promote := "请选择要创建的工作项类型"
		issueTypeSelector := issue_type_tui.NewIssueTypeSelector(categoryText, issueTypes)
		var model tea.Model
		if model, err = issueTypeSelector.Run(); err != nil {
			color.Red("任务类型选择器加载失败！")
			return
		}

		_issueTypeSelector, _ := model.(tui.Table)
		issueTypeId, err := issue_type_tui.SelectedValue(issueTypes, _issueTypeSelector.SelectedKey)

		if err != nil {
			color.Red(err.Error())
			return
		}

		// 传统选择器
		var mapSelector selector_tui.MapSelector
		var selector *tea.Program
		if len(candidateTasks) != 0 {
			options = make([]string, 0)
			optionMap = make(map[string]int, 0)
			optionMap, options = issue.FillOptions(candidateTasks, optionMap, options)
			promote = "请选择要关联的父任务"
			selector := selector_tui.NewMapSelector(optionMap, options, promote, false)
			if model, err = selector.Run(); err != nil {
				color.Red("父任务选择器加载失败！")
				return
			}
			mapSelector, _ = model.(selector_tui.MapSelector)

			parentTaskId, err = mapSelector.SelectedValue()
			if err != nil {
				color.Red(err.Error())
				return
			}
		}

		if len(candidateAssignees) != 0 {
			options = make([]string, 0)
			optionMap = make(map[string]int, 0)
			optionMap, options = member.FillOptions(candidateAssignees, optionMap, options)
			promote = "请选择指派的负责人"
			selector = selector_tui.NewMapSelector(optionMap, options, promote, false)
			if model, err = selector.Run(); err != nil {
				color.Red("负责人选择器加载失败！")
				return
			}
			mapSelector, _ = model.(selector_tui.MapSelector)

			assigneeId, err = mapSelector.SelectedValue()
			if err != nil {
				color.Red(err.Error())
				return
			}
		}

		var title string
		title = utils.ReadFromInput("填写 Issue 标题", title)
		title = strings.TrimSpace(title)

		if title == "" {
			color.Red("请输入任务标题！")
			return
		}

		var description string

		// 获取模版
		if template, err := issue_type.FetchTemplate(issueTypeId, enterprise.Id); err == nil {
			description = template
		}

		if !skipBody {
			description = utils.ReadFromEditor(utils.InitialEditor("填写 Issue 描述", description), description)
		}

		payload := map[string]interface{}{
			"description_type": "md",
			"issue_type_id":    issueTypeId,
			"title":            title,
			"description":      description,
		}
		if parentTaskId != 0 {
			payload["parent_id"] = parentTaskId
		}

		if assigneeId != 0 {
			payload["assignee_id"] = assigneeId
		}

		issue, err := issue.Create(enterprise.Id, payload)
		if err != nil {
			color.Red("创建工作项失败！")
			return
		}
		fmt.Printf("创建工作项 「%s」成功，访问地址：%s\n", utils.Cyan(issue.Title), utils.Blue(issue.Url))
	},
}

func init() {
	CreateCmd.Flags().BoolP("task", "", true, "create a task")
	CreateCmd.Flags().BoolP("bug", "", false, "create a bug")
	CreateCmd.Flags().BoolP("feature", "", false, "create a feature")
	CreateCmd.Flags().BoolP("skip-body", "", false, "skip edit issue description")
	CreateCmd.Flags().StringP("parent", "p", "", "specify the parent task by search")
	CreateCmd.Flags().StringP("assignee", "A", "", "specify the assignee by search")
}
