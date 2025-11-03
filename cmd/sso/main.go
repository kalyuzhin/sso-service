package main

import (
	"fmt"
	"github.com/kalyuzhin/sso-service/internal/app"
	"log"

	"github.com/kalyuzhin/sso-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	a := app.New(cfg.GRPC.Port)

	err = a.GRPCServer.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(cfg)
}
