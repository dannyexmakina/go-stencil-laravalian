# config/*.php, .env → config

## Laravel

```php
// config/app.php
'port' => env('APP_PORT', 8080),

// 使う側
config('app.port');
```

## このステンシル

[config/config.go](../config/config.go)

```go
type Config struct {
	AppPort   string `envconfig:"APP_PORT" default:"8080"`
	DBPath    string `envconfig:"DB_PATH" default:"database.sqlite"`
	RedisAddr string `envconfig:"REDIS_ADDR" default:"127.0.0.1:6379"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
```

## 対応表

| Laravel | このステンシル |
| --- | --- |
| `.env` / `.env.example` | 同じ（godotenv が読み込む） |
| `env('KEY', 'default')` | タグ `envconfig:"KEY" default:"..."` |
| `config('app.port')` | `cfg.AppPort`（`Load()` で得た構造体を渡す） |
| `php artisan config:cache` | 不要（起動時に1回読むだけ） |

## 注意点

- `config()` ヘルパーのようなグローバルアクセスはありません。`main.go` で `config.Load()` し、必要なところへ渡します。
- 型付きなので、`int` や `bool` のフィールドに不正な値が入っていれば起動時にエラーになります。
- 必須にしたい値は `required:"true"` タグを付けると、未設定で起動が失敗します。
- 本番では `.env` を置かず、環境変数を直接設定する運用が一般的です（`godotenv.Load()` は `.env` が無くても失敗しないようにしています）。
