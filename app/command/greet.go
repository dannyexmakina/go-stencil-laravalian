package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// `go run ./cmd/console greet Taylor`
var greetCmd = &cobra.Command{
	Use:   "greet [name]",
	Short: "挨拶を表示する",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hello, %s!\n", args[0])
	},
}

func init() {
	RootCmd.AddCommand(greetCmd)
}
