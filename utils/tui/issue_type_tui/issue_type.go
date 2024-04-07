package issue_type_tui

import (
	"errors"
	"gitee_cli/internal/api/enterprises"
	"gitee_cli/internal/api/issue_type"
	"gitee_cli/utils/tui"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func NewIssueTypeSelector(category string, issueTypes []issue_type.IssueType) *tea.Program {

	rows := make([]table.Row, 0)
	columns := []table.Column{{Title: category, Width: 20}}
	for _, issueType := range issueTypes {
		rows = append(rows, []string{issueType.Title})
	}

	t := tui.NewTableModel(enterprises.Enterprise{}, tui.IssueType, columns, rows)
	return tea.NewProgram(t)
}

func SelectedValue(issueTypes []issue_type.IssueType, selectedKey string) (int, error) {
	if selectedKey == "" {
		os.Exit(0)
	}
	for _, issueType := range issueTypes {
		if issueType.Title == selectedKey {
			return issueType.Id, nil
		}
	}
	return 0, errors.New("无效的选项！")
}
