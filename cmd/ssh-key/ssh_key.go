package ssh_key

import (
	"github.com/spf13/cobra"
)

var SshKeyCommand = &cobra.Command{
	Use:   "ssh-key",
	Short: "Manage ssh-keys",
}
