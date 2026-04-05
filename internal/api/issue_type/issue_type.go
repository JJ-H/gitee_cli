package issue_type

import (
	"errors"
	"fmt"
	"gitee_cli/internal/api/enterprises"
	"gitee_cli/utils/http_utils"
)

const (
	TASK = iota
	BUG
	REQUIREMENT
)

func typeCategory(t int) string {
	return map[int]string{
		TASK:        "task",
		BUG:         "bug",
		REQUIREMENT: "requirement",
	}[t]
}

type IssueType struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Template string `json:"template"`
}

func List(issueType int, entPath string) ([]IssueType, error) {
	ent, err := enterprises.Find(entPath)
	if err != nil {
		return nil, err
	}
	category := typeCategory(issueType)
	if category == "" {
		return nil, errors.New("无效的任务类型")
	}
	type res struct {
		Data       []IssueType `json:"data"`
		TotalCount int         `json:"total_count"`
	}
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issue_types/enterprise_issue_types?category=%s&page=1&per_page=100&state=1", ent.Id, category)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	g.SetCookieAuth()
	r, err := http_utils.DoAndDecode[res](g, "获取任务类型失败")
	if err != nil {
		return nil, err
	}
	return r.Data, nil
}

func FillOptions(issueTypes []IssueType, optionMap map[string]int, options []string) (map[string]int, []string) {
	if len(issueTypes) == 0 {
		return optionMap, options
	}
	for _, issueType := range issueTypes {
		optionMap[issueType.Title] = issueType.Id
		options = append(options, issueType.Title)
	}
	return optionMap, options
}

func FetchTemplate(issueTypeId, entId int) (string, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issue_types/%d", entId, issueTypeId)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	g.SetCookieAuth()
	issueType, err := http_utils.DoAndDecode[IssueType](g, "获取模板失败")
	if err != nil {
		return "", err
	}
	return issueType.Template, nil
}
