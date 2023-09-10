package main

import (
	"fmt"
	"log"

	"github.com/size12/planning-poker/internal/entity"
)

func main() {
	player, err := entity.NewPlayer()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(player)
}
