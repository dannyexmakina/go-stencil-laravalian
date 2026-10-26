# Service → 関数・構造体のメソッド

## Laravel

```php
class UserService
{
    public function create(string $name, string $email): User
    {
        return User::create(['name' => $name, 'email' => $email]);
    }
}
```

## このステンシル

[app/service/user_service.go](../app/service/user_service.go)

```go
type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(name, email string) (*model.User, error) {
	u := &model.User{Name: name, Email: email, Status: model.UserStatusActive}
	err := s.db.Create(u).Error
	return u, err
}
```

## 対応表

| Laravel | Go |
| --- | --- |
| サービスクラス | 依存（DBなど）をフィールドに持つ構造体 |
| `public function` | 大文字始まりのメソッド `Create` |
| `private function` | 小文字始まりのメソッド `create`（パッケージ内でのみ見える） |
| 戻り値 `User` / 例外 | `(*model.User, error)` の2値 |
| static メソッド / ヘルパー関数 | パッケージ関数 `func FormatName(...)` |
| `DB::transaction(fn() => ...)` | `s.db.Transaction(func(tx *gorm.DB) error { ... })` |

## 構造体にするか、関数にするか

- **依存（DB、外部APIクライアントなど）が必要** → 構造体 + メソッド。依存はフィールドに持ちます。
- **依存がなく入力だけで結果が決まる** → ただのパッケージ関数で十分です。Laravelでstaticヘルパーやトレイトにしていたものはこちら。

```go
// 依存がないならクラスにしない
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
```

## 注意点

- ファサード（`Cache::get()` など）のように「どこからでも呼べる」ものは作らず、必要なものは引数かフィールドで受け取ります。テストで差し替えやすくなります。
- テストでDBを差し替えたい場合は、使う側でインターフェースを定義します（Goでは「インターフェースは使う側が定義する」のが慣習です）。

```go
type UserLister interface {
	List() ([]model.User, error)
}
```

## Goの慣習では

`service` という名前自体は珍しくありませんが、`internal/user/service.go` のように**機能（ドメイン）単位でパッケージを切る**構成が好まれます。レイヤー単位（controller/service/model）で切るのはLaravel寄りの構成です。
