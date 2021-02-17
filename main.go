package main

import (
	"fmt"

	"github.com/vivian-tangle/vivian-node/config"
	"github.com/vivian-tangle/vivian-node/listener"
)

func main() {
	fmt.Println("Hello world!")
	c := config.Config{}
	c.LoadConfig()
	listener := listener.Listener{Config: &c}
	listener.Listen()
}
