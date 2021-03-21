package main

import (
	"github.com/nbkit/mdf/bootstrap"
)

func main() {
	server := bootstrap.NewServer(bootstrap.Option{})
	server.Cache().Upgrade().Cache()
	server.Start()
}
