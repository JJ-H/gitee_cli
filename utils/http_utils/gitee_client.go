package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitee_cli/config"
	"github.com/fatih/color"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const defaultTimeout = 30 * time.Second

type GiteeClient struct {
	Url        string
	Method     string
	Payload    interface{}
	Headers    map[string]string
	Response   *http.Response
	CookieAuth bool
	Query      map[string]string
	initErr    error
}

type ErrMsgV5 struct {
	Message string `json:"message"`
}

func NewGiteeClient(method, urlString string, query map[string]string, payload interface{}) *GiteeClient {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return &GiteeClient{initErr: fmt.Errorf("invalid URL %q: %w", urlString, err)}
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
	if g.initErr != nil {
		return g.initErr
	}
	g.Response = nil
	_payload, err := json.Marshal(g.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(g.Method, g.Url, bytes.NewReader(_payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
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
	client := &http.Client{Timeout: defaultTimeout}

	resp, err := client.Do(req)
	if err != nil {
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
		http.StatusOK:        {},
		http.StatusCreated:   {},
		http.StatusNoContent: {},
	}

	_, ok := successMap[g.Response.StatusCode]
	return ok
}

func (g *GiteeClient) IsFail() bool {
	return !g.IsSuccess()
}

func (g *GiteeClient) GetRespBody() ([]byte, error) {
	defer g.Response.Body.Close()
	return io.ReadAll(g.Response.Body)
}

func (g *GiteeClient) SetCookieAuth() {
	g.CookieAuth = true
}
