# Migration → Goose

## Laravel

```php
// database/migrations/2026_10_26_000000_create_users_table.php
public function up(): void
{
    Schema::create('users', function (Blueprint $table) {
        $table->id();
        $table->string('name');
        $table->string('email')->unique();
        $table->string('status');
        $table->timestamps();
        $table->softDeletes();
    });
}

public function down(): void
{
    Schema::dropIfExists('users');
}
```

```sh
php artisan migrate
```

## このステンシル

[migrations/20261026000000_create_users_table.sql](../migrations/20261026000000_create_users_table.sql)

```sql
-- +goose Up
CREATE TABLE users (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL,
    status     TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- +goose Down
DROP TABLE users;
```

実行は console コマンド経由です（[app/command/migrate.go](../app/command/migrate.go)）。

```sh
go run ./cmd/console migrate          # = up
go run ./cmd/console migrate down     # 1つ戻す
go run ./cmd/console migrate status
# または
make migrate ARGS=status
```

SQLファイルは [migrations/embed.go](../migrations/embed.go) の `//go:embed` でバイナリに埋め込んでいるので、本番では `./build/console migrate` だけで実行できます（SQLファイルを別途配置する必要はありません）。

## 対応表

| Laravel | Goose |
| --- | --- |
| PHPのスキーマビルダー（`$table->string('name')`） | **生のSQL**で書く（Go関数でも書ける） |
| `up()` / `down()` | 1ファイル内に `-- +goose Up` / `-- +goose Down` |
| `php artisan make:migration create_users_table` | `go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations create create_users_table sql` |
| `php artisan migrate` | `go run ./cmd/console migrate` |
| `php artisan migrate:rollback` | `go run ./cmd/console migrate down` |
| `php artisan migrate:status` | `go run ./cmd/console migrate status` |
| `migrations` テーブル | `goose_db_version` テーブル |
| Seeder | 標準ではない。コマンド（[09-command.md](09-command.md)）として作るのが手軽 |

## 注意点

- **モデルとテーブル定義は自動で同期されません**。GORMの構造体（[03-model.md](03-model.md)）にフィールドを足したら、対応するカラムを追加するマイグレーションも自分で書きます。`gorm.Model` を埋め込んでいる場合は `id`, `created_at`, `updated_at`, `deleted_at` が必要です。
- GORMの `AutoMigrate` は手軽ですが、カラムの削除や型変更、ロールバックができないため、このステンシルでは使っていません。
- SQLはDBごとに方言があります。サンプルはSQLite向けです。MySQL/PostgreSQLに切り替える手順は [14-database.md](14-database.md) を参照してください。
- 1つのマイグレーションの中で複数の文を書けますが、文の区切りはセミコロンです。関数定義などセミコロンを含む文は `-- +goose StatementBegin` / `-- +goose StatementEnd` で囲みます。
