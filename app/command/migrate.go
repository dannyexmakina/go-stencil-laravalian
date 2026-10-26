package command

import (
	"github.com/dannyexmakina/go-stencil-laravalian/config"
	"github.com/dannyexmakina/go-stencil-laravalian/migrations"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// `go run ./cmd/console migrate [up|down|status]`（php artisan migrate 相当）
var migrateCmd = &cobra.Command{
	Use:       "migrate [up|down|status]",
	Short:     "マイグレーションを実行する",
	Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"up", "down", "status"},
	RunE: func(cmd *cobra.Command, args []string) error {
		action := "up"
		if len(args) == 1 {
			action = args[0]
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		db, err := config.OpenDB(cfg)
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		defer sqlDB.Close()

		goose.SetBaseFS(migrations.FS)
		if err := goose.SetDialect("sqlite3"); err != nil {
			return err
		}
		return goose.RunContext(cmd.Context(), action, sqlDB, ".")
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
