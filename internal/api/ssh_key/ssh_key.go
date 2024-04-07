package ssh_key

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee_cli/config"
	"gitee_cli/utils/http_utils"
	"io/ioutil"
	"net/http"
	"os"
)

type SSHKey struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Key   string `json:"key"`
}

func AddKey(filepath, title string) (SSHKey, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return SSHKey{}, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return SSHKey{}, errors.New("读取公钥失败")
	}

	url := "https://gitee.com/api/v5/user/keys"
	payload := map[string]string{
		"key":   string(data),
		"title": title,
	}
	giteeClient := http_utils.NewGiteeClient("POST", url, nil, payload)

	giteeClient.Do()

	res, _ := giteeClient.GetRespBody()

	if giteeClient.IsFail() {
		errResponse := http_utils.ErrMsgV5{}
		err := json.Unmarshal(res, &errResponse)
		if err != nil {
			return SSHKey{}, errors.New("添加公钥失败")
		}
		return SSHKey{}, errors.New(errResponse.Message)
	}

	sshKey := SSHKey{}
	err = json.Unmarshal(res, &sshKey)
	if err != nil {
		return sshKey, errors.New("解析响应失败")
	}
	return sshKey, nil
}

func ListKeys() ([]SSHKey, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/users/%s/keys?access_token=%s", config.Conf.UserName, config.Conf.AccessToken)

	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	sshKeys := make([]SSHKey, 0)

	res, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &sshKeys)
	if err != nil {
		return nil, err
	}

	return sshKeys, nil
}

func DeleteKey(sshKeyId string) error {
	url := fmt.Sprintf("https://gitee.com/api/v5/user/keys/%s", sshKeyId)
	giteeClient := http_utils.NewGiteeClient("DELETE", url, nil, nil)

	giteeClient.Do()

	data, _ := giteeClient.GetRespBody()

	if giteeClient.IsFail() {
		if giteeClient.Response.StatusCode == http.StatusNotFound {
			return errors.New("公钥不存在")
		}
		errMsg := http_utils.ErrMsgV5{}
		json.Unmarshal(data, &errMsg)
		return errors.New(errMsg.Message)
	}
	return nil
}
