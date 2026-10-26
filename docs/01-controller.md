# Controller → handler

## Laravel

```php
class UserController extends Controller
{
    public function __construct(private UserService $users) {}

    public function index()
    {
        return response()->json($this->users->list());
    }
}
```

## このステンシル

[app/controller/user_controller.go](../app/controller/user_controller.go)

```go
type UserController struct {
	users *service.UserService
}

func NewUserController(users *service.UserService) *UserController {
	return &UserController{users: users}
}

func (ctl *UserController) Index(c echo.Context) error {
	users, err := ctl.users.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}
```

## 対応表

| Laravel | Go (Echo) |
| --- | --- |
| `class XxxController` | `type XxxController struct` |
| コンストラクタインジェクション | `NewXxxController(...)` 関数で受け取る |
| アクションメソッド `index()` | `func (ctl *XxxController) Index(c echo.Context) error` |
| `$request` | `c echo.Context`（`c.Request()` で `*http.Request`） |
| `$request->route('id')` | `c.Param("id")` |
| `$request->query('page')` | `c.QueryParam("page")` |
| `$request->all()` + FormRequest | `c.Bind(&req)` + validator（[06-validation.md](06-validation.md)） |
| `return response()->json($x, 201)` | `return c.JSON(http.StatusCreated, x)` |
| `abort(404)` | `return echo.NewHTTPError(http.StatusNotFound)` |

## 注意点

- アクションは**必ず `error` を返す**。レスポンスを書いたら `return c.JSON(...)` のようにその戻り値を返します。`error` を返すとEchoのエラーハンドラが500などに変換します。
- コントローラーの構造体は**起動時に1つだけ作られ、全リクエストで共有**されます。フィールドにリクエスト固有の値を保存しないでください。
- 戻り値の構造体はそのままJSONになります。JSONのキー名は構造体タグ `json:"name"` で指定します（Laravelの `$hidden` は `json:"-"`）。

## Goの慣習では

Goコミュニティでは「controller」ではなく **handler** と呼び、`internal/handler/` に置くのが一般的です。また、構造体を作らず `func ListUsers(svc *service.UserService) echo.HandlerFunc` のように関数で返す書き方もよく使われます。
