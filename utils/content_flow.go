package utils

import (
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/log"
)

type FlowContext struct {
	Token    *TokenContext
	Request  *ReqContext
	Response *ResContext
	Context  *gin.Context
}

func NewFlowContext() *FlowContext {
	return &FlowContext{Request: NewReqContext(), Response: NewResContext()}
}

func (s *FlowContext) Copy() *FlowContext {
	f := NewFlowContext()
	f.Request = s.Request.Copy()
	f.Token = s.Token
	f.Context = s.Context
	return f
}

func (s *FlowContext) Bind(c *gin.Context) *FlowContext {
	s.Context = c
	//bind req
	if err := c.Bind(s.Request); err != nil {
		log.ErrorD(err)
	}
	if form, err := c.MultipartForm(); err != nil {
		log.ErrorD(err)
	} else if form != nil && form.File != nil {
		s.Request.Files = form.File["files"]
	}

	//bind token

	if v, ok := s.Context.Get("context"); ok {
		if vv, is := v.(*TokenContext); is {
			s.Token = vv
		}
	}
	return s
}

// Token
func (s *FlowContext) UserID() string {
	return s.Token.UserID()
}

func (s *FlowContext) EntID() string {
	return s.Token.EntID()
}

func (s *FlowContext) OrgID() string {
	return s.Token.OrgID()
}

// Request
func (s *FlowContext) Canceled() bool {
	return s.Request.Canceled()
}

// Response
func (s FlowContext) Has(name string) bool {
	return s.Response.Has(name)
}
func (s *FlowContext) Set(name string, value interface{}) *FlowContext {
	s.Response.Set(name, value)
	return s
}

func (s *FlowContext) Get(name string) interface{} {
	return s.Response.Get(name)
}
func (s *FlowContext) SetData(value Map) *FlowContext {
	s.Response.SetData(value)
	return s
}
func (s *FlowContext) Adjust(fn func(req *FlowContext)) *FlowContext {
	fn(s)
	return s
}

func (s *FlowContext) Output() {
	s.Response.Bind(s.Context)
}
func (s *FlowContext) Error(err ...interface{}) error {
	return s.Response.Error(err...)
}
