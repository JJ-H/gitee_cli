package ssh_key

import (
	"fmt"
	"gitee_cli/internal/api/ssh_key"
	"gitee_cli/utils"
	tui "gitee_cli/utils/tui/ssh_key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	fp "path/filepath"
)

var AddSshKey = &cobra.Command{
	Use:   "add",
	Short: "Add a ssh pub key for personal",
	Run: func(cmd *cobra.Command, args []string) {
		filepath, _ := cmd.Flags().GetString("filepath")
		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			color.Red("请指定 ssh key 标题")
			return
		}
		if filepath == "" {
			homeDir, _ := os.UserHomeDir()
			sshDir := fp.Join(homeDir, ".ssh")
			files, _ := fp.Glob(fp.Join(sshDir, "*.pub"))
			if len(files) == 0 {
				color.Red("请先生成 ssh 密钥对！")
				return
			}
			fileSelector := tui.InitialUploadSSHKeyTui(files)
			var data tea.Model
			var err error
			if data, err = fileSelector.Run(); err != nil {
				color.Red("公钥选择器出错，请指定公钥地址以上传！")
				return
			}
			fileSelectRes, _ := data.(tui.UploadSSHKeyTui)
			if fileSelectRes.Cursor == -1 {
				return
			}

			filepath = fileSelectRes.FileList[fileSelectRes.Cursor]
		}

		sshKey, err := ssh_key.AddKey(filepath, title)
		if err != nil {
			color.Red(err.Error())
			return
		}
		fmt.Printf("添加 ssh key 「%s」 成功，访问地址：%s\n", utils.Yellow(sshKey.Title), utils.Cyan(fmt.Sprintf("https://gitee.com/keys/%d", sshKey.Id)))
	},
}

func init() {
	SshKeyCommand.AddCommand(AddSshKey)
	AddSshKey.Flags().StringP("filepath", "f", "", "ssh pub key filepath")
	AddSshKey.Flags().StringP("title", "t", "", "title for ssh key")
}
