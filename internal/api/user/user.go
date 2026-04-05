package user

import (
	"errors"
	"fmt"
	"gitee_cli/utils/http_utils"
)

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}

type Member struct {
	Id       int    `json:"id"`
	Remark   string `json:"remark"`
	UserName string `json:"username"`
}

func FindUser(username string) (User, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/search/users?q=%s", username)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	users, err := http_utils.DoAndDecode[[]User](g, "查询用户失败")
	if err != nil {
		return User{}, err
	}
	if len(users) == 0 {
		return User{}, nil
	}
	return users[0], nil
}

func FindMember(keyword string, enterpriseId int) (Member, error) {
	type res struct {
		Data       []Member `json:"data"`
		TotalCount int      `json:"total_count"`
	}
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/members?search=%s", enterpriseId, keyword)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	g.SetCookieAuth()
	r, err := http_utils.DoAndDecode[res](g, "查询企业成员失败")
	if err != nil {
		return Member{}, err
	}
	if len(r.Data) == 0 {
		return Member{}, errors.New("未找到成员")
	}
	return r.Data[0], nil
}

func BasicUser() (User, error) {
	g := http_utils.NewGiteeClient("GET", "https://api.gitee.com/enterprises/users", nil, nil)
	g.SetCookieAuth()
	return http_utils.DoAndDecode[User](g, "查询用户失败")
}
