# go-stencil-laravalian

Laravel経験者がGoでWebアプリを書き始めるためのステンシル。

## このステンシルについて

これはLaravel経験者向けに、意図的に"Goらしくない"ディレクトリ構成を採用したステンシルです。`controller/`, `service/`, `model/` という命名、`app/` 直下の階層構造は、Goコミュニティの慣習（`handler/`, `internal/` を基点とする構成など）とは異なります。

この構成は「Laravelでの思考習慣を保ったままGoを書き始められること」を優先しています。Goに慣れてきたら、本来の慣習に寄せていくことを推奨します。

## 移行ガイド

Laravelの各機能がGoで何に相当するかは [docs/](docs/README.md) に機能ごとにまとめています。

## ディレクトリ構造

```
.
├── cmd/
│   ├── server/         # Webサーバー (Echo)
│   ├── console/        # Artisanコマンド相当 (Cobra)
│   ├── scheduler/      # cron常駐 (robfig/cron)
│   └── worker/         # ジョブワーカー (asynq)
├── app/
│   ├── controller/     # Controllers
│   ├── service/        # Services
│   ├── model/          # Models（モデル固有の定数も同ファイル）
│   ├── middleware/     # Middleware
│   ├── command/        # Cobraコマンドの実装
│   ├── schedule/       # cronスケジュール定義
│   ├── job/            # asynqジョブの実装
│   └── constants/      # アプリ横断の定数
├── config/             # 設定 (envconfig)
├── routes/             # ルーティング定義
├── migrations/         # マイグレーション (Goose、SQLはバイナリに埋め込み)
├── docs/               # Laravel → Go 移行ガイド
└── build/              # ビルド成果物の出力先（git管理外）
```

## 想定ライブラリ

- Echo（ルーティング）
- GORM（ORM）
- Goose（マイグレーション）
- Cobra（CLI）
- robfig/cron（スケジューラー）
- asynq（ジョブキュー）
- godotenv（.env読み込み）
- envconfig（設定）
- go-playground/validator（バリデーション）
- air（ホットリロード）

## 動かし方

```sh
cp .env.example .env
go run ./cmd/console migrate      # テーブル作成（goose up）
go run ./cmd/server               # http://localhost:8080/ , GET/POST /users
go run ./cmd/console greet Taylor # Cobraコマンド
go run ./cmd/scheduler            # 毎分ログを出すcron
go run ./cmd/worker               # asynqワーカー（要Redis）
```

サンプルのDBはcgo不要のSQLite（glebarez/sqlite）です（MySQL / PostgreSQL への切り替えは [docs/14-database.md](docs/14-database.md)）。テーブルは `migrations/` のSQLをGooseで適用して作成します。

## Makefile

| コマンド | 内容 |
| --- | --- |
| `make build` | `cmd/` 以下の全バイナリ（server / console / scheduler / worker）を `build/` に出力 |
| `make dev` | air でサーバーをホットリロード起動（出力先は `build/server`）。air 未導入なら `go run github.com/air-verse/air@latest` で代用 |
| `make migrate` | マイグレーションを実行（`ARGS=down` でロールバック、`ARGS=status` で状態確認） |
| `make clean` | `build/` の中身を削除（`.gitkeep` は残す） |

## ステータス

ごく簡単なサンプル実装のみ。

## ライセンス

[MIT](LICENSE)
