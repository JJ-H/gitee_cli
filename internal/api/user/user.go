package user

import (
	"encoding/json"
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
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.Do()

	if giteeClient.IsFail() {
		return User{}, errors.New("查询用户失败！")
	}

	users := make([]User, 0)

	data, _ := giteeClient.GetRespBody()
	if err := json.Unmarshal(data, &users); err != nil {
		return User{}, errors.New("查询用户失败，解析响应失败！")
	}

	if len(users) == 0 {
		return User{}, nil
	}

	return users[0], nil
}

func FindMember(keyword string, enterpriseId int) (Member, error) {
	url := fmt.Sprintf("https://api.gitee.com/enterprises/%d/members?search=%s", enterpriseId, keyword)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	giteeClient.Do()

	data, _ := giteeClient.GetRespBody()
	if giteeClient.IsFail() {
		return Member{}, errors.New("查询企业成员失败！")
	}

	members := struct {
		Data       []Member `json:"data"`
		TotalCount int      `json:"total_count"`
	}{}

	if err := json.Unmarshal(data, &members); err != nil {
		return Member{}, errors.New("查询用户失败，解析响应失败！")
	}

	if len(members.Data) == 0 {
		return Member{}, nil
	}

	return members.Data[0], nil

}

func BasicUser() (User, error) {
	url := "https://api.gitee.com/enterprises/users"
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	err := giteeClient.Do()
	if err != nil || giteeClient.IsFail() {
		return User{}, errors.New("查询用户失败！")
	}
	user := User{}

	data, _ := giteeClient.GetRespBody()

	if err := json.Unmarshal(data, &user); err != nil {
		return User{}, errors.New("查询用户失败，解析响应失败！")
	}

	return user, nil
}
