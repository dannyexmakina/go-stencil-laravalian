package main

import (
	"log"

	"github.com/dannyexmakina/go-stencil-laravalian/app/job"
	"github.com/dannyexmakina/go-stencil-laravalian/config"
	"github.com/hibiken/asynq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: cfg.RedisAddr}, asynq.Config{Concurrency: 10})

	mux := asynq.NewServeMux()
	mux.HandleFunc(job.TypeSendWelcomeEmail, job.HandleSendWelcomeEmail)

	log.Fatal(srv.Run(mux))
}
