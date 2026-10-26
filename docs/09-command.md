# Artisan コマンド → Cobra

## Laravel

```php
class Greet extends Command
{
    protected $signature = 'greet {name}';
    protected $description = '挨拶を表示する';

    public function handle()
    {
        $this->info("Hello, {$this->argument('name')}!");
    }
}
```

```sh
php artisan greet Taylor
```

## このステンシル

[app/command/greet.go](../app/command/greet.go)

```go
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
```

```sh
go run ./cmd/console greet Taylor
# ビルド後
./build/console greet Taylor
```

## 対応表

| Laravel | Cobra |
| --- | --- |
| `php artisan` | `cmd/console`（[app/command/root.go](../app/command/root.go) の `RootCmd`） |
| `$signature = 'greet {name}'` | `Use: "greet [name]"` + `Args: cobra.ExactArgs(1)` |
| `{--force}` オプション | `cmd.Flags().Bool("force", false, "...")` |
| `$this->argument('name')` | `args[0]` |
| `$this->option('force')` | `cmd.Flags().GetBool("force")` |
| `handle()` | `Run`（エラーを返すなら `RunE`） |
| `$description` | `Short` / `Long` |
| `php artisan list` | `console --help` |
| 登録（Kernel / 自動検出） | `init()` で `RootCmd.AddCommand(...)` |

## 注意点

- 新しいコマンドは `app/command/` にファイルを追加し、`init()` で `RootCmd` に登録するだけで使えます（同じパッケージの `init()` はすべて自動で実行されます）。
- DBが必要なコマンドは、`RunE` の中で `config.Load()` → `config.OpenDB()` してから処理します。
- 失敗を終了コードに反映させたい場合は `Run` ではなく `RunE` を使って `error` を返します。
