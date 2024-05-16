package tui

import (
	"fmt"
	"gitee_cli/config"
	"gitee_cli/internal/api/enterprises"
	"gitee_cli/internal/api/issue"
	"gitee_cli/internal/api/issue_state"
	"gitee_cli/internal/api/pull_request"
	"gitee_cli/utils/git_utils"
	"gitee_cli/utils/tui/selector_tui"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"os"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

const (
	Issue = iota
	IssueType
	PullRequest
	Enterprise
	SSHKey
)

// Table TODO 初始化改成 options 模式
type Table struct {
	table        table.Model
	SelectedKey  string
	ViewMode     bool
	ResourceType int
	Enterprise   enterprises.Enterprise
}

func NewTableModel(enterprise enterprises.Enterprise, resourceType int, columns []table.Column, rows []table.Row) Table {

	height := 5
	if resourceType == PullRequest {
		height = 10
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	if resourceType == PullRequest {
		s.Selected = s.Selected.Background(lipgloss.Color("#8B4789"))
		// 避免与 diff 快捷键冲突
		t.KeyMap.HalfPageDown = key.NewBinding(
			key.WithDisabled(),
		)
	}
	t.SetStyles(s)
	return Table{table: t, Enterprise: enterprise, ResourceType: resourceType}
}

func NewTable(enterprise enterprises.Enterprise, resourceType int, columns []table.Column, rows []table.Row) *tea.Program {
	if resourceType == SSHKey || resourceType == Enterprise {
		return tea.NewProgram(NewTableModel(enterprise, resourceType, columns, rows), tea.WithAltScreen())
	}
	return tea.NewProgram(NewTableModel(enterprise, resourceType, columns, rows))
}

func (t Table) Init() tea.Cmd { return nil }

func (t Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if t.table.Focused() {
				t.table.Blur()
			} else {
				t.table.Focus()
			}
		case "c":
			if t.ResourceType == Issue || t.ResourceType == Enterprise {
				clipboard.WriteAll(t.table.SelectedRow()[0])
			} else if t.ResourceType == PullRequest {
				clipboard.WriteAll(t.table.SelectedRow()[1])
			}
		case "q", "ctrl+c":
			t.SelectedKey = ""
			return t, tea.Quit
		case "s":
			if t.ResourceType == Issue {
				targetIssue, err := issue.Detail(t.Enterprise.Id, t.table.SelectedRow()[0])
				if err != nil {
					color.Red("获取任务详情失败！")
					return t, tea.Quit
				}
				if issueStates, err := issue_state.ListWithIssue(t.Enterprise.Id, targetIssue.Id); err == nil {
					var model tea.Model
					options := make([]string, 0)
					optionMap := make(map[string]int, 0)
					optionMap, options = issue_state.FillOptions(issueStates, optionMap, options)
					promote := "请选择要变更的状态"
					selector := selector_tui.NewMapSelector(optionMap, options, promote, true)
					if model, err = selector.Run(); err != nil {
						color.Red("任务状态选择器加载失败！")
						os.Exit(1)
					}
					mapSelector, _ := model.(selector_tui.MapSelector)

					issueStateId, err := mapSelector.SelectedValue()
					if err != nil {
						color.Red(err.Error())
						os.Exit(1)
					}
					if targetIssue, err = issue.Update(t.Enterprise.Id, targetIssue.Id, map[string]interface{}{
						"issue_state_id": issueStateId,
					}); err != nil {
						color.Red(err.Error())
						os.Exit(1)
					}
					return t, tea.ClearScrollArea
				} else {
					color.Red("获取任务状态列表失败！")
					os.Exit(1)
				}
			}
		case "v":
			if t.ResourceType == Issue {
				if _issue, err := issue.Detail(t.Enterprise.Id, t.table.SelectedRow()[0]); err == nil {
					NewPager(_issue.Title, _issue.Description, Markdown).Run()
				} else {
					color.Red("获取任务详情失败！")
					return t, tea.Quit
				}
			} else if t.ResourceType == PullRequest {
				path, _ := git_utils.ParseCurrentRepo()
				if path == "" {
					path = config.Conf.DefaultPathWithNamespace
				}
				t.SelectedKey = t.table.SelectedRow()[1]
				if pullRequerst, err := pull_request.Detail(t.SelectedKey, path); err == nil {
					NewPager(pullRequerst.Title, pullRequerst.Body, Markdown).Run()
				} else {
					color.Red("获取pr详情失败！")
					return t, tea.Quit
				}
			}
		case "d":
			if t.ResourceType == PullRequest {
				path, _ := git_utils.ParseCurrentRepo()
				t.SelectedKey = t.table.SelectedRow()[1]
				if path == "" {
					path = config.Conf.DefaultPathWithNamespace
				}
				if diff, err := pull_request.FetchPatchContent(t.SelectedKey, path); err == nil {
					NewPager(t.table.SelectedRow()[0], diff, Diff).Run()
				} else {
					color.Red(err.Error())
					return t, tea.Quit
				}
			}
		case "enter":
			t.SelectedKey = t.table.SelectedRow()[0]
			if t.ResourceType == Issue {
				url := fmt.Sprintf("https://e.gitee.com/%s/dashboard?issue=%s", t.Enterprise.Path, t.SelectedKey)
				browser.OpenURL(url)
			} else if t.ResourceType == PullRequest {
				path, _ := git_utils.ParseCurrentRepo()
				if path == "" {
					path = config.Conf.DefaultPathWithNamespace
				}
				t.SelectedKey = t.table.SelectedRow()[1]
				url := fmt.Sprintf("https://gitee.com/%s/pulls/%s", path, t.SelectedKey)
				if os.Getenv("CONVERT_ENT_URL") != "" {
					url = fmt.Sprintf("https://e.gitee.com/%s/repos/%s/pulls/%s", strings.Split(path, "/")[0], path, t.SelectedKey)
				}
				browser.OpenURL(url)
			} else if t.ResourceType == Enterprise {
				path := t.table.SelectedRow()[2]
				url := fmt.Sprintf("https://e.gitee.com/%s", path)
				browser.OpenURL(url)
			} else {
				return t, tea.Quit
			}
		}
	}
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t Table) View() string {
	return baseStyle.Render(t.table.View()) + "\n"
}
