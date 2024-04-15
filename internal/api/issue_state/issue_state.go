package issue_state

import (
	"encoding/json"
	"fmt"
	"gitee_cli/utils/http_utils"
)

type IssueState struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
}

func ListWithIssue(entId int, issueId int) ([]IssueState, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issues/%d/issue_states", entId, issueId)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	if err := giteeClient.Do(); err != nil {
		return nil, err
	}

	var res = struct {
		Data       []IssueState `json:"data"`
		TotalCount int          `json:"total_count"`
	}{}

	data, _ := giteeClient.GetRespBody()

	json.Unmarshal(data, &res)

	return res.Data, nil
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
