package main

import (
	"github.com/nbkit/mdf"
	"github.com/nbkit/mdf/db"
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
	var aaa interface{}
	db.Default().Table("dddd").Select("FDafdaf").Take(&aaa)
	server := mdf.DefaultServer()
	server.Use(func(s *mdf.Server) {
		//s.GetEngine().Use(cors.AllCross())
		//rules.Register()
	})
	server.Start(server.WithOptionMDF())
	return nil
}
