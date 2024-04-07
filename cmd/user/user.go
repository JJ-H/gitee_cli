package user

import "github.com/spf13/cobra"

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "User related command",
}

func init() {
	UserCmd.AddCommand(SearchCmd)
}
