package cmd

import (
	"fmt"
	"gitee_cli/config"
	// "os"

	// "gitee_cli/config"
	"github.com/spf13/cobra"
)

var ConfigCmdUsage = "Manage Gitee CLI config, Usage: config key [value]"
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: ConfigCmdUsage,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No config key provided.")
			return
		}

		if len(args) == 1 {
			key := args[0]
			value, err := config.Read(key)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(value)
			return
		}

		if len(args) == 2 {
			key := args[0]
			value := args[1]
			if err := config.Update(map[string]interface{}{
				key: value,
			}); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(value)
			return
		}
		fmt.Println(ConfigCmdUsage)

	},
}
