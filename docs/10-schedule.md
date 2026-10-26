# Task Scheduling → robfig/cron

## Laravel

```php
// app/Console/Kernel.php
$schedule->command('report:send')->dailyAt('09:00');
$schedule->call(fn () => ...)->everyMinute();
```

```cron
* * * * * php artisan schedule:run
```

## このステンシル

[app/schedule/schedule.go](../app/schedule/schedule.go)

```go
func Register(c *cron.Cron) error {
	_, err := c.AddFunc("* * * * *", func() {
		log.Println("every minute task")
	})
	return err
}
```

```sh
go run ./cmd/scheduler   # 常駐プロセスとして起動
```

## 対応表

| Laravel | robfig/cron |
| --- | --- |
| `->everyMinute()` | `"* * * * *"` または `"@every 1m"` |
| `->hourly()` | `"@hourly"` |
| `->dailyAt('09:00')` | `"0 9 * * *"` |
| `->weekly()` | `"@weekly"` |
| `->timezone('Asia/Tokyo')` | `cron.New(cron.WithLocation(loc))` または `"CRON_TZ=Asia/Tokyo 0 9 * * *"` |
| `->withoutOverlapping()` | `cron.New(cron.WithChain(cron.SkipIfStillRunning(logger)))` |
| `schedule:run` を毎分cronで実行 | `cmd/scheduler` を**常駐**させる（OSのcron設定は不要） |

## 注意点

- 仕組みが逆です。Laravelは「OSのcronが毎分 `schedule:run` を叩く」、こちらは「プロセスが常駐して自分で時刻を見る」です。systemdやコンテナで常時起動させてください。
- 複数台で起動すると**全台で実行されます**。`onOneServer()` 相当はないので、スケジューラーは1台だけで動かすか、ロックを自前で取ります。
- 重い処理はスケジューラーで直接実行せず、ジョブをキューに積んで worker に任せるのがおすすめです（[11-job.md](11-job.md)）。
