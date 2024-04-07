package enterprise

import (
	"github.com/spf13/cobra"
)

var EntCmd = &cobra.Command{
	Use:     "enterprise",
	Aliases: []string{"ent"},
	Short:   "Manage enterprises",
}
