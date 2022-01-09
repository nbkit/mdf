package main

import (
	"github.com/nbkit/mdf"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
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
		s.GetEngine().POST("id", func(context *gin.Context) {
			utils.NewFlowContext().Bind(context)
		})
	})
	server.Start(server.WithOptionMDF())
	return nil
}
