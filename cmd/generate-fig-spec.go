package cmd

import (
	"github.com/withfig/autocomplete-tools/integrations/cobra"
)

func init() {
	RootCmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())
}
