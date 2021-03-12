package main

import (
	"github.com/nbkit/mdf/bootstrap"
)

func main() {
	server := bootstrap.NewServer()
	server.Start()
}
