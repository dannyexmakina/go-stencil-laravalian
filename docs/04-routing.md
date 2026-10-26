# routes/web.php → ルーティング

## Laravel

```php
Route::get('/users', [UserController::class, 'index']);
Route::post('/users', [UserController::class, 'store']);
Route::middleware('auth')->prefix('admin')->group(function () { ... });
```

## このステンシル

[routes/web.go](../routes/web.go)

```go
func Register(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.RequestLogger)

	users := controller.NewUserController(service.NewUserService(db))
	e.GET("/users", users.Index)
	e.POST("/users", users.Store)
}
```

## 対応表

| Laravel | Echo |
| --- | --- |
| `Route::get('/users', ...)` | `e.GET("/users", handler)` |
| `/users/{id}` | `/users/:id`（取得は `c.Param("id")`） |
| `Route::prefix('admin')->group(...)` | `g := e.Group("/admin")` |
| `->middleware('auth')` | `e.Group("/admin", authMiddleware)` または `e.GET(path, h, mw)` |
| `Route::resource(...)` | なし。1本ずつ書く |
| `->name('users.index')` / `route()` | `e.GET(...).Name = "users.index"` / `e.Reverse("users.index")` |
| ルートモデルバインディング | なし。`c.Param("id")` から自分で取得する |

## 注意点

- ルート定義は**ただの関数呼び出し**です。`routes/` はファイルを置けば読み込まれるわけではなく、`cmd/server/main.go` から `routes.Register(e, db)` を呼んでいます。
- コントローラーの組み立て（依存の注入）もここで行っています。規模が大きくなったら組み立ては `main.go` 側に寄せるのがおすすめです（[08-dependency-injection.md](08-dependency-injection.md)）。
