package issue

import (
	"fmt"
	"gitee_cli/internal/api/issue"
	"gitee_cli/internal/api/user"
	"gitee_cli/utils/tui"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List issues",
	Aliases: []string{"search"},
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("assignee")
		entPath, _ := cmd.Flags().GetString("ent")
		limit, _ := cmd.Flags().GetInt("limit")

		var search string
		payload := make(map[string]string)
		if len(args) != 0 {
			search = args[0]
		}
		if username == "" && search == "" {
			payload["only_related_me"] = "1"
		}

		if username != "" {
			if assignee, err := user.FindUser(username); err == nil {
				payload["assignee_id"] = strconv.Itoa(assignee.Id)
			}
		}
		payload["search"] = search
		payload["page"] = "1"
		payload["per_page"] = strconv.Itoa(limit)
		_issues, _ := issue.Find(enterprise.Id, payload)

		columns := []table.Column{
			{Title: "Ident", Width: 8},
			{Title: "Title", Width: 60},
		}

		rows := make([]table.Row, 0)

		for _, issue := range _issues {
			rows = append(rows, table.Row{issue.Ident, issue.Title})
		}

		issueTable := tui.NewTable(enterprise, tui.Issue, columns, rows)

		var model tea.Model
		var err error
		if model, err = issueTable.Run(); err != nil {
			color.Red("任务渲染失败！")
			os.Exit(1)
		}

		_table, _ := model.(tui.Table)
		ident := _table.SelectedKey
		if ident == "" {
			os.Exit(0)
		}
		url := fmt.Sprintf("https://e.gitee.com/%s/dashboard?issue=%s", entPath, ident)
		browser.OpenURL(url)
	},
}

func init() {
	ListCmd.Flags().StringP("assignee", "A", "", "filter issue assignee")
	ListCmd.Flags().IntP("limit", "l", 10, "limit the number of issues returned")
}
