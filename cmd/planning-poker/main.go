package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/size12/planning-poker/internal/app"
	"github.com/size12/planning-poker/internal/config"
)

func main() {

	cfg := config.GetConfig()

	poker, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed create app: %v\n", err)
	}

	go poker.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	<-c
	poker.Shutdown()
}
