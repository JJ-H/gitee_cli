package ssh_key

import (
	"fmt"
	"gitee_cli/internal/api/enterprises"
	"gitee_cli/internal/api/ssh_key"
	"gitee_cli/utils/tui"
	"github.com/charmbracelet/bubbles/table"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var ListSshKey = &cobra.Command{
	Use:   "list",
	Short: "List personal ssh pub keys",
	Long:  "List personal ssh pub keys",
	Run: func(cmd *cobra.Command, args []string) {
		sshKeys, err := ssh_key.ListKeys()
		if err != nil {
			color.Red("获取ssh公钥列表失败")
			return
		}

		if len(sshKeys) == 0 {
			color.Green("暂未添加 SSH 公钥")
			os.Exit(0)
		}

		columns := []table.Column{
			{Title: "ID", Width: 8},
			{Title: "Key Sha", Width: 38},
			{Title: "Preview URL", Width: 32},
		}

		rows := make([]table.Row, 0)
		for _, key := range sshKeys {
			rows = append(rows, table.Row{strconv.Itoa(key.Id), key.Key[:50], fmt.Sprintf("https://gitee.com/keys/%d", key.Id)})
		}
		if _, err := tui.NewTable(enterprises.Enterprise{}, tui.SSHKey, columns, rows).Run(); err != nil {
			color.Red("获取 SSH 公钥失败！")
			os.Exit(1)
		}
	},
}

func init() {
	SshKeyCommand.AddCommand(ListSshKey)
}
