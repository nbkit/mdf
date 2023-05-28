package utils

import (
	"encoding/json"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/log"
	"mime/multipart"
)

type ReqContext struct {
	Domain   string                  `json:"domain" form:"domain"`
	Widget   string                  `json:"widget"  form:"widget"`
	ID       string                  `json:"id" form:"id"`
	Type     string                  `json:"type" form:"type"`
	Code     string                  `json:"code" form:"code"`
	Name     string                  `json:"name" form:"name"`
	Tag      string                  `json:"tag" form:"tag"`
	IDS      []string                `json:"ids" form:"ids"`
	Page     int                     `json:"page" form:"page"`
	PageSize int                     `json:"page_size" form:"page_size"`
	Action   string                  `json:"action" form:"action"` // 动作编码
	Rule     string                  `json:"rule" form:"rule"`     //规则编码
	Q        string                  `json:"q" form:"q"`           //模糊查询条件
	Entity   string                  `json:"entity" form:"entity"`
	Data     interface{}             `json:"data" form:"data"`     //数据
	Extend   map[string]interface{}  `json:"extend" form:"extend"` //数据
	Param    map[string]interface{}  `json:"param" form:"param"`   //数据
	Files    []*multipart.FileHeader `json:"-" form:"files"`
	canceled bool                    `json:"-" form:"-"`
}

func NewReqContext() *ReqContext {
	return &ReqContext{}
}
func (s *ReqContext) Bind(c *gin.Context) *ReqContext {
	if err := c.ShouldBind(s); err != nil {
		log.ErrorD(err)
	}
	if s.ID == "" {
		s.ID = c.Param("id")
	}
	if s.ID == "" {
		s.ID = c.Query("id")
	}
	return s
}

func (s *ReqContext) DataBind(obj interface{}) error {
	if b, err := s.dataMarshalJSON(); err != nil {
		return err
	} else {
		return json.Unmarshal(b, obj)
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (s ReqContext) dataMarshalJSON() ([]byte, error) {
	if s.Data == nil {
		return []byte(""), nil
	}
	bytes, err := json.Marshal(s.Data)
	return bytes, err
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
		Entity: s.Entity, Data: s.Data,
	}
}
