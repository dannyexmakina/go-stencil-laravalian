package command

import (
	"github.com/dannyexmakina/go-stencil-laravalian/app/constants"
	"github.com/spf13/cobra"
)

// RootCmd は `php artisan` 相当
var RootCmd = &cobra.Command{
	Use:   "console",
	Short: constants.AppName + " console",
}
