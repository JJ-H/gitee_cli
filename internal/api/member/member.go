package member

import (
	"encoding/json"
	"fmt"
	"gitee_cli/utils/http_utils"
)

type Member struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Remark   string `json:"remark"`
}

func Find(enterpriseId int, params map[string]string) ([]Member, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/members", enterpriseId)
	giteeClient := http_utils.NewGiteeClient("GET", url, params, nil)
	giteeClient.SetCookieAuth()

	_, err := giteeClient.Do()
	if err != nil || giteeClient.IsFail() {
		return []Member{}, err
	}

	data, _ := giteeClient.GetRespBody()
	type res struct {
		Data       []Member `json:"data"`
		TotalCount int      `json:"total_count"`
	}

	var _data res

	json.Unmarshal(data, &_data)

	return _data.Data, nil
}

func FillOptions(members []Member, optionMap map[string]int, options []string) (map[string]int, []string) {
	if len(members) == 0 {
		return optionMap, options
	}

	for _, member := range members {
		key := fmt.Sprintf("%s(%s)", member.Name, member.Remark)
		optionMap[key] = member.Id
		options = append(options, key)
	}
	return optionMap, options
}
