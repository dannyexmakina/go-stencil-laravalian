# Eloquent Model → GORM の構造体

## Laravel

```php
class User extends Model
{
    use SoftDeletes;
    const STATUS_ACTIVE = 'active';
    protected $fillable = ['name', 'email', 'status'];
}
```

## このステンシル

[app/model/user.go](../app/model/user.go)

```go
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

type User struct {
	gorm.Model
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"uniqueIndex"`
	Status string `json:"status"`
}
```

## 対応表

| Laravel (Eloquent) | GORM |
| --- | --- |
| `id`, `created_at`, `updated_at` | `gorm.Model` を埋め込む |
| `use SoftDeletes` | `gorm.Model` の `DeletedAt` で自動的に有効 |
| `$fillable` | 不要（構造体のフィールドがそのまま対象） |
| `$hidden` | `json:"-"` タグ |
| `$casts` | フィールドの型そのもの（`time.Time`, `bool` など） |
| `User::all()` | `db.Find(&users)` |
| `User::find($id)` | `db.First(&user, id)` |
| `User::where('status', 'active')->get()` | `db.Where("status = ?", "active").Find(&users)` |
| `User::create([...])` | `db.Create(&user)` |
| `$user->save()` | `db.Save(&user)` |
| `$user->delete()` | `db.Delete(&user)` |
| `hasMany` / `belongsTo` | フィールド + `gorm:"foreignKey:..."`、取得は `db.Preload("Posts")` |
| `ModelNotFoundException` | `errors.Is(err, gorm.ErrRecordNotFound)` |

## 注意点

- **Active Recordではない**: `$user->save()` のようにモデル自身がDBを操作するのではなく、`db` に構造体を渡します。モデルはただのデータです。
- **結果は引数のポインタに書き込まれる**: `db.Find(&users)` の `&` を忘れないこと。
- **ゼロ値に注意**: `db.Where(&User{Status: ""})` のように構造体で条件を渡すと、空文字や `0` は条件から無視されます。条件は文字列 + プレースホルダで書くのが安全です。
- **アクセサ・ミューテタはない**: 必要ならメソッドとして書きます（`func (u *User) IsActive() bool`）。
- モデル固有の定数は、このステンシルでは同じファイルに置いています（`UserStatusActive` のように型名を接頭辞にする）。

## Goの慣習では

ORMを使わず `database/sql` や、SQLから型安全なコードを生成する **sqlc** を使うプロジェクトも多くあります。GORMは「Eloquentに近い書き心地」を優先した選択です。
