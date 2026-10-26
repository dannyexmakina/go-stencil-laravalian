# Middleware

## Laravel

```php
class RequestLogger
{
    public function handle(Request $request, Closure $next)
    {
        $start = microtime(true);
        $response = $next($request);
        Log::info(...);
        return $response;
    }
}
```

## このステンシル

[app/middleware/request_logger.go](../app/middleware/request_logger.go)

```go
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		log.Printf("%s %s (%s)", c.Request().Method, c.Path(), time.Since(start))
		return err
	}
}
```

## 対応表

| Laravel | Echo |
| --- | --- |
| `handle($request, $next)` | `func(next echo.HandlerFunc) echo.HandlerFunc` |
| `$next($request)` | `next(c)` |
| 処理を止める（`abort(401)`） | `next` を呼ばずに `return echo.NewHTTPError(http.StatusUnauthorized)` |
| グローバルミドルウェア（Kernel） | `e.Use(mw)` |
| ルートミドルウェア | `e.Group("/x", mw)` / `e.GET(path, h, mw)` |
| `$request->attributes->set()` / `auth()->user()` | `c.Set("user", u)` / `c.Get("user").(*model.User)` |
| パラメータ付き（`throttle:60,1`） | 引数を受け取ってミドルウェアを返す関数を作る |

```go
// パラメータ付きミドルウェア
func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// role を使ったチェック
			return next(c)
		}
	}
}
```

## 既製品

CORS、リカバリー（panic→500）、レートリミットなどは `github.com/labstack/echo/v4/middleware` に揃っています。自作する前に確認してください（このステンシルの `app/middleware` とパッケージ名が同じなので、importでは別名を付けます）。
