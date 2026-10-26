package config

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// OpenDB はDB接続を開く。サンプルでは cgo 不要の SQLite を使う
func OpenDB(c *Config) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(c.DBPath), &gorm.Config{})
}
