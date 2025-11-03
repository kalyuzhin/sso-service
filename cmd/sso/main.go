package main

import (
	"fmt"

	"github.com/kalyuzhin/sso-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

}
