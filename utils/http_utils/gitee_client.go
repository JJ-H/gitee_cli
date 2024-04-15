package http_utils

import (
	"bytes"
	"encoding/json"
	"gitee_cli/config"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type GiteeClient struct {
	Url        string
	Method     string
	Payload    interface{}
	Headers    map[string]string
	Response   *http.Response
	CookieAuth bool
	Query      map[string]string
}

type ErrMsgV5 struct {
	Message string `json:"message"`
}

func NewGiteeClient(method, urlString string, query map[string]string, payload interface{}) *GiteeClient {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	if query != nil {
		queryParams := parsedUrl.Query()
		for k, v := range query {
			queryParams.Set(k, v)
		}
		parsedUrl.RawQuery = queryParams.Encode()
	}
	return &GiteeClient{
		Method:  method,
		Url:     parsedUrl.String(),
		Payload: payload,
		Query:   query,
	}
}

func (g *GiteeClient) SetHeaders(headers map[string]string) {
	g.Headers = headers
}

func (g *GiteeClient) Do() error {
	// 多次调用首先置空
	g.Response = nil
	_payload, _ := json.Marshal(g.Payload)
	req, _ := http.NewRequest(g.Method, g.Url, bytes.NewReader(_payload))
	req.Header.Set("Content-Type", "application/json")
	cookie := config.Conf.CookiesJar
	accessToken := config.Conf.AccessToken
	if accessToken != "" && !g.CookieAuth {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else if cookie != "" {
		req.Header.Set("Cookie", cookie)
	} else {
		color.Red("授权错误！")
		os.Exit(1)
	}
	for key, value := range g.Headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}

	var resp *http.Response
	var err error

	if resp, err = client.Do(req); err != nil {
		return err
	}
	g.Response = resp
	return nil
}

func (g *GiteeClient) IsSuccess() bool {
	if g.Response == nil {
		return false
	}

	successMap := map[int]struct{}{
		http.StatusOK:        struct{}{},
		http.StatusCreated:   struct{}{},
		http.StatusNoContent: struct{}{},
	}

	if _, ok := successMap[g.Response.StatusCode]; ok {
		return true
	}
	return false
}

func (g *GiteeClient) IsFail() bool {
	return !g.IsSuccess()
}

func (g *GiteeClient) GetRespBody() ([]byte, error) {
	return ioutil.ReadAll(g.Response.Body)
}

func (g *GiteeClient) SetCookieAuth() {
	g.CookieAuth = true
}
