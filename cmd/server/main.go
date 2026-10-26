package main

import (
	"log"

	"github.com/dannyexmakina/go-stencil-laravalian/config"
	"github.com/dannyexmakina/go-stencil-laravalian/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// 事前に `go run ./cmd/console migrate` でテーブルを作成しておく
	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	routes.Register(e, db)
	log.Fatal(e.Start(":" + cfg.AppPort))
}
