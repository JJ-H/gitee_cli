package issue_state

import (
	"fmt"
	"gitee_cli/utils/http_utils"
)

type IssueState struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
}

func ListWithIssue(entId int, issueId int) ([]IssueState, error) {
	type res struct {
		Data       []IssueState `json:"data"`
		TotalCount int          `json:"total_count"`
	}
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issues/%d/issue_states", entId, issueId)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	g.SetCookieAuth()
	r, err := http_utils.DoAndDecode[res](g, "获取任务状态列表失败")
	if err != nil {
		return nil, err
	}
	return r.Data, nil
}

func FillOptions(issueStates []IssueState, optionMap map[string]int, options []string) (map[string]int, []string) {
	if len(issueStates) == 0 {
		return optionMap, options
	}
	for _, issueState := range issueStates {
		optionMap[issueState.Title] = issueState.Id
		options = append(options, issueState.Title)
	}
	return optionMap, options
}
