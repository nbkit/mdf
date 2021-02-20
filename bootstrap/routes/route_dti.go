package routes

import (
	"github.com/nbkit/mdf/bootstrap/dti"
	"github.com/nbkit/mdf/gin"
)

func routeDti(engine *gin.Engine) {
	engine.POST("/dti/{group:string}/{name:path}", func(ctx *gin.Context) {
		hand := &dti.DtiHandProc{Group: ctx.Param("group"), Ctx: ctx, Path: ctx.Param("name")}
		if hand.Group != "" {
			hand.Do()
		}
	})
	engine.POST("/dti/{name:path}", func(ctx *gin.Context) {
		hand := &dti.DtiHandProc{Group: "dti", Ctx: ctx, Path: ctx.Param("name")}
		if hand.Group != "" {
			hand.Do()
		}
	})
}
