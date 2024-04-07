package issue_type

import (
	"encoding/json"
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
		return nil, errors.New("无效的任务类型！")
	}
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issue_types/enterprise_issue_types?category=%s&page=1&per_page=1000&state=1", ent.Id, category)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	if _, err := giteeClient.Do(); err != nil {
		return nil, err
	}

	var res = struct {
		Data       []IssueType `json:"data"`
		TotalCount int         `json:"total_count"`
	}{}

	data, _ := giteeClient.GetRespBody()

	json.Unmarshal(data, &res)

	return res.Data, nil
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
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	if _, err := giteeClient.Do(); err != nil {
		return "", errors.New("获取模板失败！")
	}
	issueType := IssueType{}
	data, _ := giteeClient.GetRespBody()
	json.Unmarshal(data, &issueType)
	return issueType.Template, nil
}
