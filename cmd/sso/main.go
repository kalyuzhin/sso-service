package main

import (
	"github.com/kalyuzhin/sso-service/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kalyuzhin/sso-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	a := app.New(cfg.GRPC.Port)

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
