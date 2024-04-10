package cmd

import (
	"fmt"
	"gitee_cli/cmd/auth"
	"gitee_cli/cmd/enterprise"
	"gitee_cli/cmd/issue"
	"gitee_cli/cmd/pull_request"
	sshkey "gitee_cli/cmd/ssh-key"
	"gitee_cli/cmd/user"
	"gitee_cli/utils"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:     "gitee",
	Short:   "Gitee In terminal",
	Long:    "Gitee CLI is a tool which interact with gitee server seamlessly via terminal",
	Version: "0.0.3",
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(pull_request.Pr)
	RootCmd.AddCommand(issue.IssueCmd)
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(enterprise.EntCmd)
	RootCmd.AddCommand(sshkey.SshKeyCommand)
	RootCmd.AddCommand(user.UserCmd)
	RootCmd.SetVersionTemplate(fmt.Sprintf("Gitee CLI Version %s\n", utils.Cyan(RootCmd.Version)))
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
