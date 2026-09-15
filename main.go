package main

import (
	"fmt"

	"github.com/Solaire89/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}
	err = cfg.SetUser("ray")
	if err != nil {
		fmt.Print(err)
	}
	newCfg, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%+v", newCfg)
}
