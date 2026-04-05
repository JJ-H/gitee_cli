package git_utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	BRANCH_PREFIX = "refs/heads/"
	HTTP_PREFIX   = "https://gitee.com/"
	SSH_PREFIX    = "git@gitee.com:"
	GIT_SUFFIX    = ".git"
)

func CurrentDir() string {
	wd, _ := os.Getwd()
	return wd
}

func IsGitDir() bool {
	wd := CurrentDir()
	if _, err := os.Stat(fmt.Sprintf("%s/.git", wd)); err != nil {
		return false
	}
	return true
}

func GetCurrentBranch() (string, error) {
	catFile := exec.Command("cat", ".git/HEAD")
	extractBranch := exec.Command("awk", "{print $2}")

	var output bytes.Buffer
	catFile.Stdout = &output
	extractBranch.Stdin = &output
	if err := catFile.Run(); err != nil {
		return "", fmt.Errorf("读取 .git/HEAD 失败: %w", err)
	}
	res, err := extractBranch.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "", errors.New("获取当前分支异常")
	}
	return strings.TrimSpace(strings.TrimPrefix(string(res), BRANCH_PREFIX)), nil
}

func ParseCurrentRepo() (string, error) {
	if !IsGitDir() {
		return "", errors.New("请在仓库目录下执行该命令！")
	}
	gitRemote := exec.Command("git", "remote")
	gitRemote.Dir = CurrentDir()
	remoteOutput, err := gitRemote.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取 git remote 失败: %w", err)
	}
	remoteName := strings.Split(strings.TrimSpace(string(remoteOutput)), "\n")[0]
	if remoteName == "" {
		return "", errors.New("未找到 git remote，请先添加 remote")
	}
	getUrl := exec.Command("git", "remote", "get-url", remoteName)
	getUrl.Dir = CurrentDir()
	urlOutput, err := getUrl.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取 remote URL 失败: %w", err)
	}
	gitUrl := strings.TrimSpace(string(urlOutput))
	gitUrl = strings.TrimPrefix(gitUrl, HTTP_PREFIX)
	gitUrl = strings.TrimPrefix(gitUrl, SSH_PREFIX)
	pathWithNamespace := strings.TrimSuffix(gitUrl, GIT_SUFFIX)
	return pathWithNamespace, nil
}
