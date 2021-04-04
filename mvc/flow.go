package mvc

import (
	"github.com/nbkit/mdf/gin"
)

type Flow interface {
	Bind(c *gin.Context)
}
