package auth

import (
	"fmt"
	"gitee_cli/internal/api/user"
	"gitee_cli/utils"
	"os"

	"gitee_cli/config"
	"github.com/spf13/cobra"
)

var (
	CookiesFile   string
	GlobalCookies string
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate Gitee CLI with gitee enterprise",
	Run: func(cmd *cobra.Command, args []string) {
		if CookiesFile != "" {
			// Read cookies from file
			cookies, err := os.ReadFile(CookiesFile)
			if err != nil {
				fmt.Println("Failed to read cookies-file:", err)
				return
			}
			GlobalCookies = string(cookies)
		}

		// Save cookies to file for future usage
		if GlobalCookies != "" {
			if err := config.Update(map[string]interface{}{
				"cookies_jar": GlobalCookies,
			}); err != nil {
				fmt.Println("Failed to save cookies:", err)
				return
			}
		}
		user, err := user.BasicUser()
		if err != nil {
			fmt.Println("Authorize error: ", err)
			return
		}
		fmt.Println(fmt.Sprintf("Hi「%s」! You've %s authenticated", utils.Cyan(user.Name), utils.Green("successfully")))
	},
}

func init() {
	AuthCmd.Flags().StringVarP(&CookiesFile, "cookies-file", "f", "", "path to a file containing cookies")
}
