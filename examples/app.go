package main

import (
	"github.com/nbkit/mdf"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/middleware/cors"
	"os"
)

func runApp() error {
	defer func() {
		if r := recover(); r != nil {
			log.ErrorD(r)
			os.Exit(0)
		}
	}()
	server := mdf.DefaultServer()
	server.Use(func(s *mdf.Server) {
		s.GetEngine().Use(cors.AllCross())
	})
	server.Start(server.WithOptionMigrate(), server.WithOptionAuthToken(false), server.WithOptionRule())
	return nil
}
