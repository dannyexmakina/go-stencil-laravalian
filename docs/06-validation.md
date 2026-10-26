# FormRequest → リクエスト構造体 + validator

## Laravel

```php
class StoreUserRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'name'  => 'required|max:255',
            'email' => 'required|email',
        ];
    }
}
```

## このステンシル

[app/controller/user_controller.go](../app/controller/user_controller.go)

```go
type StoreUserRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email"`
}

func (ctl *UserController) Store(c echo.Context) error {
	var req StoreUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	// ...
}
```

## ルールの対応

| Laravel | go-playground/validator |
| --- | --- |
| `required` | `required` |
| `max:255`（文字列） | `max=255` |
| `email` | `email` |
| `in:a,b` | `oneof=a b` |
| `integer|min:1` | フィールドを `int` にして `min=1` |
| `nullable` | ポインタ型（`*string`）+ `omitempty` |
| `confirmed` | `eqfield=PasswordConfirmation` |
| `unique:users` | なし。サービス層でDBを確認する |

## 注意点

- Laravelのように**自動では実行されません**。`c.Bind` → `validate.Struct` をアクション内で明示的に呼びます。何度も書くのが面倒なら `e.Validator` に登録して `c.Validate(&req)` を使う方法もあります。
- エラーメッセージは英語の開発者向けです。ユーザー向けの日本語メッセージが必要なら `validator.ValidationErrors` を走査して自分で組み立てます。
- `c.Bind` はJSON・フォーム・クエリ・パスパラメータをまとめて読みます。受け付けたくない値が入らないよう、リクエスト専用の構造体を作り、モデルの構造体を直接 `Bind` しないでください（mass assignment対策）。
