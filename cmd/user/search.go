package user

import (
	"fmt"
	"gitee_cli/config"
	"gitee_cli/internal/api/enterprises"
	"gitee_cli/internal/api/user"
	"gitee_cli/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a user info, Usage: gitee user search {username}",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			color.Red("请给定用户名！")
			return
		}
		username := args[0]
		isEntPath, _ := cmd.Flags().GetBool("ent")
		if isEntPath {
			entPath := config.Conf.DefaultEntPath

			if entPath == "" {
				color.Red("请使用 gitee config default_ent_path xxx 指定默认 path!")
				return
			}

			enterprise, err := enterprises.Find(entPath)

			if err != nil {
				color.Red("企业未找到！")
				return
			}

			var member user.Member

			if member, err = user.FindMember(username, enterprise.Id); err == nil {
				fmt.Printf("成员 ID：%s\n成员名称：%s\n用户名：%s\n", utils.Cyan(member.Id), utils.Cyan(member.Remark), utils.Blue(member.UserName))
			} else {
				color.Red("未查找到对应用户！")
				return
			}
		} else {
			user, err := user.FindUser(username)

			if err != nil {
				color.Red("查询用户失败！%v", err)
				return
			}

			if user.Id == 0 {
				color.Red("未查找到对应用户！")
				return
			}
			fmt.Printf("用户 ID：%s\n用户名称：%s\n用户主页：%s\n", utils.Cyan(user.Id), utils.Cyan(user.Name), utils.Blue(user.HtmlUrl))
		}

	},
}

func init() {
	SearchCmd.Flags().BoolP("ent", "e", false, "search member from current enterprise")
}
