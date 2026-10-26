# go-stencil-laravalian

Laravel経験者がGoでWebアプリを書き始めるためのステンシル。

## このステンシルについて

これはLaravel経験者向けに、意図的に"Goらしくない"ディレクトリ構成を採用したステンシルです。`controller/`, `service/`, `model/` という命名、`app/` 直下の階層構造は、Goコミュニティの慣習（`handler/`, `internal/` を基点とする構成など）とは異なります。

この構成は「Laravelでの思考習慣を保ったままGoを書き始められること」を優先しています。Goに慣れてきたら、本来の慣習に寄せていくことを推奨します。

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
└── migrations/         # マイグレーション (Goose or golang-migrate)
```

## 想定ライブラリ

- Echo（ルーティング）
- GORM（ORM）
- Goose または golang-migrate（マイグレーション、選択制）
- Cobra（CLI）
- robfig/cron（スケジューラー）
- asynq（ジョブキュー）
- godotenv（.env読み込み）
- envconfig（設定）
- go-playground/validator（バリデーション）
- air（ホットリロード）

## ステータス

構造のみ。中身の実装はこれから。
