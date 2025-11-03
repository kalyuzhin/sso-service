package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kalyuzhin/sso-service/internal/app"
	"github.com/kalyuzhin/sso-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()

	a, err := app.New(ctx, cfg.Database.GetDSN(), cfg.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := a.GRPCServer.Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	a.GRPCServer.Stop()
}
