package main

import (
	"log"

	"github.com/dannyexmakina/go-stencil-laravalian/app/schedule"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	if err := schedule.Register(c); err != nil {
		log.Fatal(err)
	}
	c.Run() // ブロックする
}
