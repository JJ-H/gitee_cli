package enterprises

import (
	"fmt"
	"gitee_cli/utils/http_utils"
)

type Enterprise struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func List() ([]Enterprise, error) {
	type res struct {
		Data       []Enterprise `json:"data"`
		TotalCount int          `json:"total_count"`
	}
	g := http_utils.NewGiteeClient("GET", "https://api.gitee.com/enterprises/list", nil, nil)
	g.SetCookieAuth()
	r, err := http_utils.DoAndDecode[res](g, "获取企业列表失败")
	if err != nil {
		return nil, err
	}
	return r.Data, nil
}

func Find(path string) (Enterprise, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/enterprises/%s", path)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	return http_utils.DoAndDecode[Enterprise](g, "查询企业失败")
}
