package pull_request

import (
	"errors"
	"fmt"
	"gitee_cli/config"
	"gitee_cli/utils/git_utils"
	"gitee_cli/utils/http_utils"
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

func RepoPath() string {
	path, err := git_utils.ParseCurrentRepo()
	if err != nil {
		return config.Conf.DefaultPathWithNamespace
	}
	return path
}

func Note(iid int, note string) error {
	url := fmt.Sprintf(config.Conf.ApiPrefix+"/repos/%s/pulls/%d/comments", RepoPath(), iid)
	g := http_utils.NewGiteeClient("POST", url, nil, map[string]string{"body": note})
	if err := g.Do(); err != nil || g.IsFail() {
		return fmt.Errorf("评论 PR #%d 失败", iid)
	}
	return nil
}

func List(scope string, limit int) []PullRequest {
	url := fmt.Sprintf(config.Conf.ApiPrefix+"/repos/%s/pulls?state=open&sort=created&direction=desc&page=1&per_page=%d", RepoPath(), limit)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	pullRequests, err := http_utils.DoAndDecode[[]PullRequest](g, "获取 PR 列表失败")
	if err != nil {
		return nil
	}
	return filterPullRequest(pullRequests, scope)
}

func filterPullRequest(pullRequests []PullRequest, scope string) []PullRequest {
	if len(pullRequests) == 0 {
		return pullRequests
	}
	userId := config.Conf.UserId
	filtered := make([]PullRequest, 0)
	if scope == "owner" {
		for _, pr := range pullRequests {
			if pr.Creator.Id == userId {
				filtered = append(filtered, pr)
			}
		}
	} else {
		for _, pr := range pullRequests {
			for _, a := range pr.Assignees {
				if a.Id == userId {
					pr.User = a
					filtered = append(filtered, pr)
					break
				}
			}
		}
	}
	return filtered
}

func (pr PullRequest) TransferUrlToEnt() string {
	data := strings.Split(pr.HtmlUrl, "/")
	iid := data[len(data)-1]
	path, err := git_utils.ParseCurrentRepo()
	if err != nil {
		path = config.Conf.DefaultPathWithNamespace
	}
	return fmt.Sprintf("https://e.gitee.com/oschina/repos/%s/pulls/%s", path, iid)
}

func FuzzySearch(pullRequests []PullRequest, keyword string) []PullRequest {
	if len(pullRequests) == 0 {
		return pullRequests
	}
	result := make([]PullRequest, 0)
	for _, pr := range pullRequests {
		if strings.Contains(pr.Title, keyword) {
			result = append(result, pr)
		}
	}
	return result
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
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	return http_utils.DoAndDecode[PullRequest](g, "查找 PR 失败")
}

func findPrIidBySha(commitSha, dir, currentBranch string) (int, error) {
	command := fmt.Sprintf("git log --merges --ancestry-path --oneline %s..%s | grep 'pull request' | tail -n1 | awk '{print $2\";\"$3}'", commitSha, currentBranch)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = dir
	res, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("获取 PR ID 失败: %w", err)
	}
	result := strings.Split(string(res), ";")
	if len(result) != 2 || strings.TrimSpace(result[1]) == "" {
		return 0, errors.New("未找到匹配的 Pull Request")
	}
	iid, _ := strconv.Atoi(strings.TrimPrefix(result[0], "!"))
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
	g := http_utils.NewGiteeClient("POST", url, nil, payload)
	return http_utils.DoAndDecode[PullRequest](g, "创建 PR 失败")
}

func CreateLightPr(baseRepo, baseRef, prTitle string) (PullRequest, error) {
	unixTime := time.Now().Format("20060102150405")
	branch := "test_" + unixTime
	path := "test_" + unixTime + ".txt"

	g := http_utils.NewGiteeClient("POST", fmt.Sprintf("https://gitee.com/api/v5/repos/%s/branches", baseRepo), nil,
		map[string]string{"refs": baseRef, "branch_name": branch})
	if err := g.Do(); err != nil || g.IsFail() {
		return PullRequest{}, errors.New("创建临时分支失败")
	}

	g.Url = fmt.Sprintf("https://gitee.com/api/v5/repos/%s/contents/%s", baseRepo, path)
	g.Payload = map[string]string{"message": "test", "content": "test", "branch": branch}
	if err := g.Do(); err != nil || g.IsFail() {
		return PullRequest{}, errors.New("创建临时提交失败")
	}

	g.Url = fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls", baseRepo)
	g.Payload = map[string]string{"title": prTitle, "head": branch, "base": baseRef}
	return http_utils.DoAndDecode[PullRequest](g, "创建 PR 失败")
}

func Detail(iid, repoPath string) (PullRequest, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%s", repoPath, iid)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	return http_utils.DoAndDecode[PullRequest](g, "获取 PR 详情失败")
}

func FetchPatchContent(iid, repoPath string) (string, error) {
	url := fmt.Sprintf("https://gitee.com/%s/pulls/%s.diff", repoPath, iid)
	g := http_utils.NewGiteeClient("GET", url, nil, nil)
	g.SetCookieAuth()
	if err := g.Do(); err != nil || g.IsFail() {
		return "", errors.New("获取 PR diff 内容失败")
	}
	data, err := g.GetRespBody()
	if err != nil {
		return "", fmt.Errorf("读取 PR diff 失败: %w", err)
	}
	return string(data), nil
}

func Close(iid string) error {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%s", RepoPath(), iid)
	g := http_utils.NewGiteeClient("PATCH", url, nil, map[string]string{"state": "closed"})
	if err := g.Do(); err != nil || g.IsFail() {
		return fmt.Errorf("关闭 PR #%s 失败", iid)
	}
	return nil
}

func Review(iid string) error {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls/%s/review", RepoPath(), iid)
	g := http_utils.NewGiteeClient("POST", url, nil, nil)
	if err := g.Do(); err != nil || g.IsFail() {
		return fmt.Errorf("审查 PR #%s 失败", iid)
	}
	return nil
}
