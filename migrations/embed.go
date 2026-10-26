package migrations

import "embed"

// FS はマイグレーションSQLをバイナリに埋め込む
//
//go:embed *.sql
var FS embed.FS
