package enterprises

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee_cli/utils/http_utils"
)

type Enterprise struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func List() ([]Enterprise, error) {
	url := "https://api.gitee.com/enterprises/list"
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	err := giteeClient.Do()
	if err != nil || giteeClient.IsFail() {
		return nil, err
	}
	data, _ := giteeClient.GetRespBody()

	type res struct {
		Data       []Enterprise `json:"data"`
		TotalCount int          `json:"total_count"`
	}

	var _data res

	json.Unmarshal(data, &_data)

	return _data.Data, nil
}

func Find(path string) (Enterprise, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/enterprises/%s", path)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	if err := giteeClient.Do(); err != nil {
		return Enterprise{}, errors.New("查询企业失败！")
	}
	data, _ := giteeClient.GetRespBody()
	ent := Enterprise{}

	json.Unmarshal(data, &ent)

	return ent, nil
}
