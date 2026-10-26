# DB接続（SQLite → MySQL / PostgreSQL）

## 現状

サンプルは**SQLite専用**です。手軽に動かせるよう、cgo不要のドライバ `glebarez/sqlite` を使っています。

SQLiteに依存しているのは次の3か所だけです。モデル・サービス・コントローラーはGORM経由でDBを扱っているので、変更は不要です。

| 箇所 | SQLite依存の内容 |
| --- | --- |
| [config/database.go](../config/database.go) | GORMドライバ `glebarez/sqlite` |
| [config/config.go](../config/config.go) | 設定がファイルパス `DB_PATH` のみ |
| [app/command/migrate.go](../app/command/migrate.go) と [migrations/](../migrations/) | `goose.SetDialect("sqlite3")` と SQLite用のSQL |

## Laravel との対応

| Laravel | このステンシル |
| --- | --- |
| `config/database.php` の `connections` | [config/database.go](../config/database.go) の `OpenDB` |
| `DB_CONNECTION=mysql` | `DB_DRIVER=mysql`（以下の手順で追加する） |
| `DB_HOST` / `DB_PORT` / `DB_DATABASE` / `DB_USERNAME` / `DB_PASSWORD` | 個別に持つか、まとめて `DB_DSN` にする |
| PDOドライバ（`pdo_mysql` 拡張） | Goのドライバパッケージを `go get` する |

## 切り替え手順

### 1. ドライバを追加する

```sh
# MySQL
go get gorm.io/driver/mysql
# PostgreSQL
go get gorm.io/driver/postgres
```

### 2. 設定を変える

[config/config.go](../config/config.go) の `DBPath` を、ドライバ名と接続文字列（DSN）に置き換えます。

```go
type Config struct {
	AppPort   string `envconfig:"APP_PORT" default:"8080"`
	DBDriver  string `envconfig:"DB_DRIVER" default:"mysql"`
	DBDSN     string `envconfig:"DB_DSN" required:"true"`
	RedisAddr string `envconfig:"REDIS_ADDR" default:"127.0.0.1:6379"`
}
```

`.env` / `.env.example` の例：

```sh
# MySQL
DB_DRIVER=mysql
DB_DSN=app:secret@tcp(127.0.0.1:3306)/app?charset=utf8mb4&parseTime=True&loc=Local

# PostgreSQL
DB_DRIVER=postgres
DB_DSN=host=127.0.0.1 port=5432 user=app password=secret dbname=app sslmode=disable TimeZone=Asia/Tokyo
```

> MySQLでは `parseTime=True` を必ず付けてください。付けないと `DATETIME` を `time.Time` に読み込めずエラーになります。

Laravelのように `DB_HOST` などを個別に持ちたい場合は、フィールドを分けて `OpenDB` の中で `fmt.Sprintf` でDSNを組み立てます。

### 3. 接続処理を変える

[config/database.go](../config/database.go)：

```go
package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(c *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch c.DBDriver {
	case "mysql":
		dialector = mysql.Open(c.DBDSN)
	case "postgres":
		dialector = postgres.Open(c.DBDSN)
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER: %s", c.DBDriver)
	}
	return gorm.Open(dialector, &gorm.Config{})
}
```

使わないドライバは import しなくて構いません（1つに決めるなら `switch` も不要です）。

### 4. goose の dialect を変える

[app/command/migrate.go](../app/command/migrate.go) の `"sqlite3"` を設定値に置き換えます。goose の dialect 名は `mysql` / `postgres` で、上の `DB_DRIVER` と同じ名前にしておけばそのまま渡せます。

```go
if err := goose.SetDialect(cfg.DBDriver); err != nil {
	return err
}
```

### 5. マイグレーションのSQLを書き直す

[migrations/20261026000000_create_users_table.sql](../migrations/20261026000000_create_users_table.sql) はSQLiteの書き方なので、使うDBに合わせます。`gorm.Model` を埋め込んだモデルに合うよう、`id` / `created_at` / `updated_at` / `deleted_at` は残します。

MySQL：

```sql
-- +goose Up
CREATE TABLE users (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    status     VARCHAR(32)  NOT NULL,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    UNIQUE KEY idx_users_email (email),
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE users;
```

PostgreSQL：

```sql
-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    status     VARCHAR(32)  NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- +goose Down
DROP TABLE users;
```

### 6. 不要になったものを片付ける

```sh
go mod tidy   # 使わなくなった glebarez/sqlite が go.mod から消える
```

`.gitignore` の `*.sqlite` も不要になります。

### 7. 確認する

```sh
go run ./cmd/console migrate status
go run ./cmd/console migrate
go run ./cmd/server
```

## 注意点

- **DBは事前に作っておく**: Laravelと同じく、`CREATE DATABASE` はマイグレーションでは行いません。
- **SQLiteとMySQL/PostgreSQLを環境で切り替えたい場合**（ローカルはSQLite、本番はMySQLなど）、Goのコードは `switch` で両対応できますが、**マイグレーションのSQLはDBごとに必要**です。`migrations/sqlite/`・`migrations/mysql/` のようにディレクトリを分け、`goose.RunContext` に渡すディレクトリを `DB_DRIVER` で切り替えます。ただしローカルと本番でDBの挙動が変わるので、Docker などで本番と同じDBをローカルでも使うほうがおすすめです。
- **接続プール**: Laravel（PHP-FPM）はリクエストごとに接続しますが、Goは常駐プロセスが接続を使い回します。必要に応じて上限を設定してください。

  ```go
  sqlDB, _ := db.DB()
  sqlDB.SetMaxOpenConns(25)
  sqlDB.SetMaxIdleConns(25)
  sqlDB.SetConnMaxLifetime(5 * time.Minute)
  ```

- **PostgreSQLの検索条件**: `?` プレースホルダはGORMが変換するので、そのまま `db.Where("status = ?", ...)` で書けます。
