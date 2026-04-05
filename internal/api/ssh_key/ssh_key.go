package ssh_key

import (
	"encoding/json"
	"fmt"
	"gitee_cli/config"
	"gitee_cli/utils/http_utils"
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
	data, err := os.ReadFile(filepath)
	if err != nil {
		return SSHKey{}, fmt.Errorf("读取公钥失败: %w", err)
	}
	payload := map[string]string{"key": string(data), "title": title}
	g := http_utils.NewGiteeClient("POST", "https://gitee.com/api/v5/user/keys", nil, payload)
	return http_utils.DoAndDecode[SSHKey](g, "添加公钥失败")
}

func ListKeys() ([]SSHKey, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/users/%s/keys?access_token=%s", config.Conf.UserName, config.Conf.AccessToken)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	return http_utils.DoAndDecode[[]SSHKey](g, "获取公钥列表失败")
}

func DeleteKey(sshKeyId string) error {
	url := fmt.Sprintf("https://gitee.com/api/v5/user/keys/%s", sshKeyId)
	g := http_utils.NewGiteeClient("DELETE", url, nil, nil)
	if err := g.Do(); err != nil {
		return fmt.Errorf("删除公钥失败: %w", err)
	}
	if g.IsFail() {
		if g.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("公钥不存在")
		}
		body, _ := g.GetRespBody()
		var apiErr http_utils.ErrMsgV5
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			return fmt.Errorf("删除公钥失败: %s", apiErr.Message)
		}
		return fmt.Errorf("删除公钥失败: HTTP %d", g.Response.StatusCode)
	}
	return nil
}
