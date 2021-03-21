package utils

import (
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/log"
	"mime/multipart"
)

type ReqContext struct {
	Domain    string                  `json:"domain" form:"domain"`
	Widget    string                  `json:"widget"  form:"widget"`
	Params    map[string]interface{}  `json:"params"  form:"params"` //一般指 页面 URI 参数
	ID        string                  `json:"id" form:"id"`
	IDS       []string                `json:"ids" form:"ids"`
	Page      int                     `json:"page" form:"page"`
	PageSize  int                     `json:"page_size" form:"page_size"`
	Action    string                  `json:"action" form:"action"`       // 动作编码
	Rule      string                  `json:"rule" form:"rule"`           //规则编码
	Q         string                  `json:"q" form:"q"`                 //模糊查询条件
	Condition interface{}             `json:"condition" form:"condition"` //附件条件
	Entity    string                  `json:"entity" form:"entity"`
	Data      interface{}             `json:"data" form:"data"` //数据
	Tag       string                  `json:"tag" form:"tag"`
	Files     []*multipart.FileHeader `json:"-" form:"files"`
	canceled  bool                    `json:"-" form:"-"`
}

func NewReqContext() *ReqContext {
	return &ReqContext{}
}
func (s *ReqContext) Bind(c *gin.Context) *ReqContext {
	if err := c.Bind(&s); err != nil {
		log.ErrorD(err)
	}
	if form, err := c.MultipartForm(); err != nil {
		log.ErrorD(err)
	} else if form != nil && form.File != nil {
		s.Files = form.File["files"]
	}
	return s
}

func (s *ReqContext) SetCancel(isCanceled bool) *ReqContext {
	s.canceled = isCanceled
	return s
}
func (s *ReqContext) Canceled() bool {
	return s.canceled
}
func (s *ReqContext) Adjust(fn func(req *ReqContext)) *ReqContext {
	fn(s)
	return s
}
func (s *ReqContext) Copy() *ReqContext {
	return &ReqContext{
		Widget: s.Widget, Domain: s.Domain,
		ID: s.ID, IDS: s.IDS,
		Page: s.Page, PageSize: s.PageSize, Q: s.Q,
		Action: s.Action, Rule: s.Rule,
		Condition: s.Condition, Entity: s.Entity, Data: s.Data,
	}
}
