# Service Container → 手動での組み立て

## Laravel

型宣言しておけばサービスコンテナが自動で解決します。

```php
public function __construct(private UserService $users) {}
```

## このステンシル

コンテナはありません。`New...` 関数で依存を渡しながら、**起動時に一度だけ手で組み立てます**。

```go
// cmd/server/main.go → routes/web.go
cfg, _ := config.Load()
db, _ := config.OpenDB(cfg)

userService := service.NewUserService(db)
userController := controller.NewUserController(userService)
```

## 対応表

| Laravel | Go |
| --- | --- |
| 自動解決（オートワイヤリング） | `NewXxx(依存...)` を自分で呼ぶ |
| `$this->app->singleton(...)` | 起動時に1回作って使い回す（Goではこれが普通） |
| `$this->app->bind(Interface::class, Impl::class)` | インターフェース型の引数に実装を渡す |
| ServiceProvider | `main.go`（またはそこから呼ぶ組み立て関数） |
| `app(UserService::class)`（どこからでも取得） | しない。引数で受け取る |

## 注意点

- 面倒に見えますが、依存関係がコードに全部書かれているので「何がどこから来ているか」を追いやすいのが利点です。
- 組み立てが長くなってきたら、`wire`（コード生成）や `fx` などのDIツールもありますが、まずは手書きで十分です。
- 循環依存（AがBを、BがAを必要とする）はパッケージのimportレベルでコンパイルエラーになります。設計を見直すサインです。
