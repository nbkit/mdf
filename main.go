package main

import (
	"github.com/nbkit/mdf/bootstrap"
	"github.com/nbkit/mdf/log"
)

func main() {
	log.ErrorS().Latency().Output()
	server := bootstrap.NewServer()
	server.Start()
}
