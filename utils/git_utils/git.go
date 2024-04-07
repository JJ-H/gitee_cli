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
	err := catFile.Run()
	res, err := extractBranch.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "", errors.New("获取当前分支异常")
	}
	return strings.TrimSpace(strings.TrimPrefix(string(res), BRANCH_PREFIX)), nil
}

func ParseCurrentRepo() (string, error) {
	var err error
	var pathWithNamespace string
	if !IsGitDir() {
		return "", errors.New("请在仓库目录下执行该命令！")
	}
	gitRemote := exec.Command("git", "remote")
	gitRemote.Dir = CurrentDir()
	output, err := gitRemote.CombinedOutput()
	getUrl := exec.Command("git", "remote", "get-url", strings.Split(string(output), "\n")[0])
	getUrl.Dir = CurrentDir()
	output, err = getUrl.CombinedOutput()
	gitUrl := strings.Trim(string(output), "\n")
	gitUrl = strings.TrimPrefix(gitUrl, HTTP_PREFIX)
	gitUrl = strings.TrimPrefix(gitUrl, SSH_PREFIX)
	pathWithNamespace = strings.TrimSuffix(gitUrl, GIT_SUFFIX)
	return pathWithNamespace, err
}
