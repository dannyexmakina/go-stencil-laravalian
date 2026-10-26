# 例外処理 → error 戻り値

Laravel（PHP）から来ると一番戸惑う部分です。

## Laravel

```php
try {
    $user = User::findOrFail($id);
} catch (ModelNotFoundException $e) {
    abort(404);
}
```

## Go

```go
var user model.User
if err := db.First(&user, id).Error; err != nil {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return err
}
```

## 対応表

| Laravel / PHP | Go |
| --- | --- |
| `throw new Exception('...')` | `return errors.New("...")` |
| 例外を投げる関数 | 最後の戻り値が `error` の関数 |
| `try { } catch (X $e) { }` | `if err != nil { }` をその場で書く |
| `catch (ModelNotFoundException $e)` | `errors.Is(err, gorm.ErrRecordNotFound)` |
| カスタム例外クラス | `var ErrXxx = errors.New(...)` またはエラー型 |
| `$e->getPrevious()`（原因の連鎖） | `fmt.Errorf("ユーザー作成: %w", err)` でラップ |
| `instanceof` で例外の種類を判定 | `errors.Is` / `errors.As` |
| `app/Exceptions/Handler.php` | `e.HTTPErrorHandler` |
| `abort(404)` | `return echo.NewHTTPError(http.StatusNotFound)` |

## 書き方の基本

```go
// サービス層: 意味のあるエラーを定義し、文脈を付けて返す
var ErrEmailTaken = errors.New("email already taken")

func (s *UserService) Create(name, email string) (*model.User, error) {
	if s.emailExists(email) {
		return nil, ErrEmailTaken
	}
	u := &model.User{Name: name, Email: email}
	if err := s.db.Create(u).Error; err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

// コントローラー層: エラーの種類に応じてHTTPステータスに変換する
u, err := ctl.users.Create(req.Name, req.Email)
if errors.Is(err, service.ErrEmailTaken) {
	return echo.NewHTTPError(http.StatusConflict, err.Error())
}
if err != nil {
	return err // 500
}
```

## 注意点

- `err` を `_` で捨てない。`if err != nil` が多いのは冗長ではなく、Goではそれが正しい書き方です。
- `panic` は例外の代わりではありません。「プログラムのバグで続行不可能」な場合だけに使います（Echoのリカバリーミドルウェアで500に変換できますが、それに頼った設計にはしません）。
- サービス層ではHTTPのことを知らない `error` を返し、HTTPステータスへの変換はコントローラー層で行います。
