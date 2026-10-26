# Queue / Job → asynq

## Laravel

```php
class SendWelcomeEmail implements ShouldQueue
{
    public function __construct(public int $userId) {}

    public function handle() { ... }
}

SendWelcomeEmail::dispatch($user->id);
```

```sh
php artisan queue:work
```

## このステンシル

[app/job/send_welcome_email.go](../app/job/send_welcome_email.go)

```go
const TypeSendWelcomeEmail = "email:welcome"

type SendWelcomeEmailPayload struct {
	UserID uint `json:"user_id"`
}

func NewSendWelcomeEmailTask(userID uint) (*asynq.Task, error) { ... }

func HandleSendWelcomeEmail(ctx context.Context, t *asynq.Task) error { ... }
```

ワーカー側の登録は [cmd/worker/main.go](../cmd/worker/main.go)：

```go
mux := asynq.NewServeMux()
mux.HandleFunc(job.TypeSendWelcomeEmail, job.HandleSendWelcomeEmail)
```

dispatch する側（例）：

```go
client := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr})
defer client.Close()

task, _ := job.NewSendWelcomeEmailTask(user.ID)
_, err := client.Enqueue(task)
```

## 対応表

| Laravel | asynq |
| --- | --- |
| Job クラス | タスク名の定数 + ペイロード構造体 + ハンドラ関数 |
| コンストラクタの引数（シリアライズされる） | ペイロード構造体（JSONにして渡す） |
| `handle()` | `func(ctx context.Context, t *asynq.Task) error` |
| `Job::dispatch(...)` | `client.Enqueue(task)` |
| `->delay(now()->addMinutes(5))` | `client.Enqueue(task, asynq.ProcessIn(5*time.Minute))` |
| `->onQueue('emails')` | `asynq.Queue("emails")` オプション |
| `$tries = 3` | `asynq.MaxRetry(3)` |
| `php artisan queue:work` | `go run ./cmd/worker` |
| `failed_jobs` テーブル | Redis 上のアーカイブ（Web UI の asynqmon で確認可能） |
| `QUEUE_CONNECTION=redis` | asynq は Redis 専用 |

## 注意点

- ハンドラが `error` を返すとリトライされます。リトライしても無駄なエラーは `fmt.Errorf("...: %w", asynq.SkipRetry)` で返します。
- Laravelのように**モデルを丸ごと渡さず**、IDだけをペイロードに入れて、ハンドラ内でDBから取り直します（Laravelの `SerializesModels` が裏でやっていることと同じです）。
- ワーカーはRedisが起動していないと動きません。
