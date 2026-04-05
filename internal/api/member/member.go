package member

import (
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
	type res struct {
		Data       []Member `json:"data"`
		TotalCount int      `json:"total_count"`
	}
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/members", enterpriseId)
	g := http_utils.NewGiteeClient("GET", url, params, nil)
	g.SetCookieAuth()
	r, err := http_utils.DoAndDecode[res](g, "获取成员列表失败")
	if err != nil {
		return nil, err
	}
	return r.Data, nil
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
