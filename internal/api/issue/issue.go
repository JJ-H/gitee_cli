package issue

import (
	"fmt"
	"gitee_cli/internal/api/issue_state"
	"gitee_cli/internal/api/issue_type"
	"gitee_cli/internal/api/user"
	"gitee_cli/utils/http_utils"
)

const endpoint = "https://api.gitee.com/enterprises/%d/issues"

type Issue struct {
	Id          int                    `json:"id"`
	Ident       string                 `json:"ident"`
	Title       string                 `json:"title"`
	Url         string                 `json:"issue_url"`
	Description string                 `json:"description"`
	IssueState  issue_state.IssueState `json:"issue_state"`
	Assignee    user.Member            `json:"assignee"`
	IssueType   issue_type.IssueType   `json:"issue_type"`
}

type listRes struct {
	Data       []Issue `json:"data"`
	TotalCount int     `json:"total_count"`
}

func Find(enterpriseId int, params map[string]string) ([]Issue, error) {
	url := fmt.Sprintf(endpoint, enterpriseId)
	g := http_utils.NewGiteeClient("GET", url, params, nil)
	g.SetCookieAuth()
	r, err := http_utils.DoAndDecode[listRes](g, "获取任务列表失败")
	if err != nil {
		return nil, err
	}
	return r.Data, nil
}

func Create(enterpriseId int, payload map[string]interface{}) (Issue, error) {
	url := fmt.Sprintf(endpoint, enterpriseId)
	g := http_utils.NewGiteeClient("POST", url, nil, payload)
	g.SetCookieAuth()
	return http_utils.DoAndDecode[Issue](g, "创建工作项失败")
}

func Update(enterpriseId int, issueId int, payload map[string]interface{}) (Issue, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issues/%d", enterpriseId, issueId)
	g := http_utils.NewGiteeClient("PUT", url, nil, payload)
	g.SetCookieAuth()
	return http_utils.DoAndDecode[Issue](g, "更新工作项失败")
}

func Detail(enterpriseId int, ident string) (Issue, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issues/%s?qt=ident", enterpriseId, ident)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	g.SetCookieAuth()
	return http_utils.DoAndDecode[Issue](g, "获取工作项详情失败")
}

func FillOptions(issues []Issue, optionMap map[string]int, options []string) (map[string]int, []string) {
	if len(issues) == 0 {
		return optionMap, options
	}
	for _, issue := range issues {
		key := fmt.Sprintf("[%s] %s", issue.Ident, issue.Title)
		optionMap[key] = issue.Id
		options = append(options, key)
	}
	return optionMap, options
}
