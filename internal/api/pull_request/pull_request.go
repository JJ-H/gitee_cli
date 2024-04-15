package pull_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee_cli/config"
	"gitee_cli/utils/git_utils"
	"gitee_cli/utils/http_utils"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type PullRequest struct {
	Id            int        `json:"id"`
	Title         string     `json:"title"`
	HtmlUrl       string     `json:"html_url"`
	Mergeable     bool       `json:"mergeable"`
	CanMergeCheck bool       `json:"can_merge_check"`
	PatchUrl      string     `json:"patch_url"`
	Draft         bool       `json:"draft"`
	Creator       creator    `json:"user"`
	Assignees     []assignee `json:"assignees"`
	User          assignee   `json:"-"`
	Number        int        `json:"number"`
	Body          string     `json:"body"`
}

type assignee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Accept bool   `json:"accept"`
}

type creator struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func Note(iid int, note string) error {
	// https://gitee.com/api/v5/repos/{owner}/{repo}/pulls/{number}/comments
	pathWithNamespace, err := git_utils.ParseCurrentRepo()
	if err != nil {
		pathWithNamespace = config.Conf.DefaultPathWithNamespace
	}
	url := fmt.Sprintf(config.Conf.ApiPrefix+"/repos/%s/pulls/%d/comments", pathWithNamespace, iid)

	payload := map[string]string{"body": note}
	giteeClient := http_utils.NewGiteeClient("POST", url, nil, payload)

	giteeClient.Do()
	if giteeClient.IsFail() {
		return errors.New(fmt.Sprintf("评论 pr %d 失败！", iid))
	}
	return nil
}

func List(scope string) []PullRequest {
	pathWithNamespace, err := git_utils.ParseCurrentRepo()
	if err != nil {
		pathWithNamespace = config.Conf.DefaultPathWithNamespace
	}
	url := fmt.Sprintf(config.Conf.ApiPrefix+"/repos/%s/pulls?state=open&sort=created&direction=desc&page=1&per_page=100", pathWithNamespace)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)

	pullRequests := make([]PullRequest, 50)

	giteeClient.Do()

	res, _ := giteeClient.GetRespBody()
	err = json.Unmarshal(res, &pullRequests)
	if err != nil {
		fmt.Println("解析 Pull Request 异常！")
		return nil
	}
	pullRequests = filterPullRequest(pullRequests, scope)
	return pullRequests
}

func filterPullRequest(pullRequests []PullRequest, scope string) []PullRequest {
	if len(pullRequests) == 0 {
		return pullRequests
	}

	// 过滤指定 User 审查的 PR
	filteredPullRequests := make([]PullRequest, 0)
	userId := config.Conf.UserId

	if scope == "owner" {
		for _, pr := range pullRequests {
			if pr.Creator.Id == userId {
				filteredPullRequests = append(filteredPullRequests, pr)
			}
		}
	} else {
		for _, pr := range pullRequests {
			assignees := pr.Assignees
			for _, assignee := range assignees {
				if assignee.Id == userId {
					pr.User = assignee
					filteredPullRequests = append(filteredPullRequests, pr)
					break
				}
			}
		}
	}

	return filteredPullRequests
}

func (pr PullRequest) TransferUrlToEnt() string {
	url := pr.HtmlUrl
	data := strings.Split(url, "/")
	iid := data[len(data)-1]
	pathWithName, err := git_utils.ParseCurrentRepo()
	if err != nil {
		pathWithName = config.Conf.DefaultPathWithNamespace
	}
	return fmt.Sprintf("https://e.gitee.com/oschina/repos/%s/pulls/%s", pathWithName, iid)
}

func FuzzySearch(pullRequests []PullRequest, keyword string) []PullRequest {
	if len(pullRequests) == 0 {
		return pullRequests
	}

	// 模糊搜索
	fuzzyPullRequests := make([]PullRequest, 0)

	for _, pr := range pullRequests {
		if strings.Contains(pr.Title, keyword) {
			fuzzyPullRequests = append(fuzzyPullRequests, pr)
		}
	}
	return fuzzyPullRequests
}

func FindPullRequestByIid(commitSha, pathWithNamespace string) (PullRequest, error) {

	currentBranch, err := git_utils.GetCurrentBranch()
	if err != nil {
		return PullRequest{}, err
	}

	iid, err := findPrIidBySha(commitSha, git_utils.CurrentDir(), currentBranch)
	if err != nil {
		return PullRequest{}, err
	}
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%d", pathWithNamespace, iid)

	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.Do()
	if giteeClient.IsFail() {
		return PullRequest{}, errors.New("查找 pr 异常！")
	}
	res, _ := giteeClient.GetRespBody()

	pullRequest := PullRequest{}
	err = json.Unmarshal(res, &pullRequest)
	if err != nil {
		fmt.Println("解析 Pull Request 异常！")
		return PullRequest{}, errors.New("解析 Pull Request 异常！")
	}
	return pullRequest, nil
}

func findPrIidBySha(commitSha, dir, currentBranch string) (int, error) {
	command := fmt.Sprintf("git log --merges --ancestry-path --oneline %s..%s | grep 'pull request' | tail -n1 | awk '{print $2\";\"$3}'", commitSha, currentBranch)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = dir
	res, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		return 0, errors.New("获取 Pull Request ID 失败！")
	}
	result := strings.Split(string(res), ";")
	if len(result) != 2 {
		return 0, errors.New("未找到匹配的 Pull Request！")
	}
	iid, _ := strconv.Atoi(strings.TrimPrefix(result[0], "!"))
	title := result[1]
	if title == "" {
		return 0, errors.New("获取 Pull Request ID 失败！")
	}
	return iid, nil
}

func CreatePr(baseRepo, baseRef, headRef, title, body, assignees, testers string, draft bool, prune bool) (PullRequest, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls", baseRepo)
	payload := map[string]interface{}{
		"base":                baseRef,
		"head":                headRef,
		"title":               title,
		"body":                body,
		"assignees":           assignees,
		"testers":             testers,
		"draft":               draft,
		"prune_source_branch": prune,
		"assignees_number":    len(strings.Split(assignees, ",")),
		"testers_number":      len(strings.Split(testers, ",")),
	}

	giteeClient := http_utils.NewGiteeClient("POST", url, nil, payload)
	err := giteeClient.Do()

	if err != nil {
		return PullRequest{}, errors.New("GiteeCilent 异常！")
	}
	res, _ := giteeClient.GetRespBody()

	if giteeClient.IsFail() {
		errResponse := http_utils.ErrMsgV5{}
		err := json.Unmarshal(res, &errResponse)
		if err != nil {
			return PullRequest{}, errors.New("创建 pull request 失败！")
		}
		return PullRequest{}, errors.New(errResponse.Message)
	}

	pullRequest := PullRequest{}
	if err := json.Unmarshal(res, &pullRequest); err != nil {
		return pullRequest, errors.New("解析响应失败")
	}
	return pullRequest, nil
}

func CreateLightPr(baseRepo, baseRef, prTitle string) (PullRequest, error) {
	content := "test"
	message := "test"
	unixTime := time.Now().Format("20060102150405")
	path := "test_" + unixTime + ".txt"
	branch := "test_" + unixTime
	// 新建分支
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/branches", baseRepo)
	payload := map[string]string{"refs": baseRef, "branch_name": branch}
	giteeClient := http_utils.NewGiteeClient("POST", url, nil, payload)

	giteeClient.Do()
	if giteeClient.IsFail() {
		return PullRequest{}, errors.New("创建 PR 失败")
	}

	giteeClient.Payload = map[string]string{"message": message, "content": content, "branch": branch}
	giteeClient.Url = fmt.Sprintf("https://gitee.com/api/v5/repos/%s/contents/%s", baseRepo, path)

	giteeClient.Do()
	if giteeClient.IsFail() {
		return PullRequest{}, errors.New("创建 PR 失败")
	}

	// 创建 pr
	giteeClient.Url = fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls", baseRepo)
	giteeClient.Payload = map[string]string{
		"title": prTitle,
		"head":  branch,
		"base":  baseRef,
	}
	giteeClient.Do()
	if giteeClient.IsFail() {
		return PullRequest{}, errors.New("创建 PR 失败")
	}
	res, _ := giteeClient.GetRespBody()
	pullRequest := PullRequest{}
	err := json.Unmarshal(res, &pullRequest)
	if err != nil {
		fmt.Println("解析 Pull Request 异常！")
		return PullRequest{}, errors.New("解析 Pull Request 异常！")
	}
	return pullRequest, nil
}

func Detail(iid, repoPath string) (PullRequest, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%s", repoPath, iid)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	if err := giteeClient.Do(); err != nil || giteeClient.IsFail() {
		return PullRequest{}, errors.New("获取 Pull Request 详情失败")
	}

	_pullRequest := PullRequest{}

	data, _ := giteeClient.GetRespBody()

	json.Unmarshal(data, &_pullRequest)
	return _pullRequest, nil
}

func FetchPatchContent(iid, repoPath string) (string, error) {
	url := fmt.Sprintf("https://gitee.com/%s/pulls/%s.diff", repoPath, iid)
	giteeClient := http_utils.NewGiteeClient("GET", url, nil, nil)
	giteeClient.SetCookieAuth()
	if err := giteeClient.Do(); err != nil || giteeClient.IsFail() {
		return "", errors.New("获取 Pull Request 补丁内容失败")
	}
	data, _ := giteeClient.GetRespBody()

	return string(data), nil
}

func Close(iid string) {
	pathWithNamespace, err := git_utils.ParseCurrentRepo()
	if err != nil {
		pathWithNamespace = config.Conf.DefaultPathWithNamespace
	}
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%s", pathWithNamespace, iid)
	payload := map[string]string{"state": "closed"}
	giteeClient := http_utils.NewGiteeClient("PATCH", url, nil, payload)
	if giteeClient.Do(); giteeClient.IsFail() {
		color.Red("关闭 PR 失败！")
		os.Exit(1)
	}
	color.Green("关闭 PR 成功🏅")
}

func Review(iid string) {
	pathWithNamespace, err := git_utils.ParseCurrentRepo()
	if err != nil {
		pathWithNamespace = config.Conf.DefaultPathWithNamespace
	}
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%s/review", pathWithNamespace, iid)
	giteeClient := http_utils.NewGiteeClient("POST", url, nil, nil)
	giteeClient.Do()
	if giteeClient.IsFail() {
		color.Red("审查通过失败！")
		os.Exit(1)
	}
	color.Green("审查通过成功🏅")
}
