# Laravel → Go 移行ガイド

Laravelの各機能が、このステンシル（およびGo一般）では何に相当するかを機能ごとにまとめたものです。各ページは「Laravelではこう → このステンシルではこう → Goの慣習ではこう」の順に書いています。

| Laravel | このステンシル | Goの慣習的な呼び方 | ページ |
| --- | --- | --- | --- |
| Controller | `app/controller/` (Echo) | handler | [01-controller.md](01-controller.md) |
| Service / Facade | `app/service/` | 関数・構造体のメソッド | [02-service.md](02-service.md) |
| Eloquent Model | `app/model/` (GORM) | model / repository | [03-model.md](03-model.md) |
| routes/web.php | `routes/` | router | [04-routing.md](04-routing.md) |
| Middleware | `app/middleware/` | middleware | [05-middleware.md](05-middleware.md) |
| FormRequest | リクエスト構造体 + validator | request DTO | [06-validation.md](06-validation.md) |
| config/*.php, .env | `config/` (envconfig + godotenv) | config | [07-config.md](07-config.md) |
| Service Container (DI) | `main.go` での手動組み立て | 明示的なDI | [08-dependency-injection.md](08-dependency-injection.md) |
| Artisan コマンド | `app/command/` (Cobra) | CLI | [09-command.md](09-command.md) |
| Task Scheduling | `app/schedule/` (robfig/cron) | cron | [10-schedule.md](10-schedule.md) |
| Queue / Job | `app/job/` (asynq) | worker / task | [11-job.md](11-job.md) |
| Migration | `migrations/` (Goose) | migration | [12-migration.md](12-migration.md) |
| DB接続 (config/database.php) | `config/database.go` (GORM) | - | [14-database.md](14-database.md) |
| 例外処理 | `error` 戻り値 | error handling | [13-error-handling.md](13-error-handling.md) |

## 最初に押さえる考え方の違い

- **魔法が少ない**: Laravelのファサード、自動DI、マジックメソッドに相当するものはありません。必要なものは引数で渡し、`main.go` で組み立てます。
- **例外ではなく `error` を返す**: 失敗する可能性のある関数は最後の戻り値で `error` を返し、呼び出し側がその場で処理します。
- **プロセスが常駐する**: PHP-FPMのようにリクエストごとに状態がリセットされません。パッケージ変数（グローバル変数）に状態を持たせると全リクエストで共有されます。
- **型で守る**: 配列（連想配列）の代わりに構造体を定義します。`$request->input('name')` ではなく `req.Name` です。
- **パッケージ = ディレクトリ**: 1ディレクトリ1パッケージ。名前が大文字で始まる識別子だけが外部から見えます（`public` / `private` の代わり）。
