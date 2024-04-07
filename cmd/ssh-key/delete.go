package ssh_key

import (
	"gitee_cli/internal/api/ssh_key"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var DeleteSshKey = &cobra.Command{
	Use:   "delete",
	Short: "delete a specified ssh key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sshKeyId := args[0]
		if sshKeyId == "" {
			color.Red("请提供正确的 SSH Key ID")
			return
		}
		err := ssh_key.DeleteKey(sshKeyId)
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.Green("删除公钥成功")
	},
}

func init() {
	SshKeyCommand.AddCommand(DeleteSshKey)
}
