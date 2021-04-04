package mvc

import (
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/utils"
)

type Service struct {
	Flow *utils.FlowContext
}

func (s *Service) Bind(c *gin.Context) {
	s.Flow = utils.NewFlowContext().Bind(c)
}
