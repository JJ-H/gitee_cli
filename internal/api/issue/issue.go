package issue

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee_cli/utils/http_utils"
)

const Endpoint = "https://api.gitee.com/enterprises/%d/issues"

type Issue struct {
	Id          int    `json:"id"`
	Ident       string `json:"ident"`
	Title       string `json:"title"`
	Url         string `json:"issue_url"`
	Description string `json:"description"`
}

func Find(enterpriseId int, params map[string]string) ([]Issue, error) {
	url := fmt.Sprintf(Endpoint, enterpriseId)
	giteeClient := http_utils.NewGiteeClient("GET", url, params, nil)
	giteeClient.SetCookieAuth()

	_, err := giteeClient.Do()
	if err != nil || giteeClient.IsFail() {
		return []Issue{}, err
	}

	data, _ := giteeClient.GetRespBody()
	type res struct {
		Data       []Issue `json:"data"`
		TotalCount int     `json:"total_count"`
	}

	var _data res

	json.Unmarshal(data, &_data)

	return _data.Data, nil
}

func Create(enterpriseId int, payload map[string]interface{}) (Issue, error) {
	url := fmt.Sprintf(Endpoint, enterpriseId)
	giteeClient := http_utils.NewGiteeClient("POST", url, nil, payload)
	giteeClient.SetCookieAuth()

	giteeClient.Do()

	if giteeClient.IsFail() {
		return Issue{}, errors.New("创建工作项失败！")
	}

	issue := Issue{}

	data, _ := giteeClient.GetRespBody()
	json.Unmarshal(data, &issue)

	return issue, nil
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

func Detail(enterpriseId int, ident string) (Issue, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/issues/%s?qt=ident", enterpriseId, ident)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	giteeClient.Do()
	if giteeClient.IsFail() {
		return Issue{}, errors.New("获取工作想失败！")
	}

	data, _ := giteeClient.GetRespBody()
	issue := Issue{}
	json.Unmarshal(data, &issue)
	return issue, nil
}
