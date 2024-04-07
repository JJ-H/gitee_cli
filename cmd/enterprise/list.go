package enterprise

import (
	enterprises2 "gitee_cli/internal/api/enterprises"
	"gitee_cli/utils/tui"
	"github.com/charmbracelet/bubbles/table"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all enterprises joined by me",
	Run: func(cmd *cobra.Command, args []string) {
		enterprises, err := enterprises2.List()
		if err != nil {
			color.Red("获取企业列表失败！")
			return
		}

		columns := []table.Column{
			{Title: "ID", Width: 8},
			{Title: "Name", Width: 28},
			{Title: "Path", Width: 28},
		}
		rows := make([]table.Row, 0)
		for _, ent := range enterprises {
			rows = append(rows, table.Row{strconv.Itoa(ent.Id), ent.Name, ent.Path})
		}

		entTable := tui.NewTable(enterprises2.Enterprise{}, tui.Enterprise, columns, rows)

		if _, err := entTable.Run(); err != nil {
			color.Red("企业渲染失败！")
			os.Exit(1)
		}
	},
}

func init() {
	EntCmd.AddCommand(ListCmd)
}
