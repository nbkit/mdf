package main

import (
	"github.com/nbkit/mdf"
	"github.com/nbkit/mdf/log"
	"os"
)

func runApp() error {
	defer func() {
		if r := recover(); r != nil {
			log.ErrorD(r)
			os.Exit(0)
		}
	}()
	server := mdf.NewServer(mdf.Config{})

	server.Use(func(s *mdf.Server) {
		//s.GetEngine().Use(cors.AllCross())
		//rules.Register()
	})
	server.Start()
	return nil
}