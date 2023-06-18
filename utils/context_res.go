package utils

import (
	"fmt"
	"github.com/nbkit/mdf/gin"
	"net/http"
	"strings"
)

type ResContext struct {
	data   Map
	errors []error
}

func NewResContext() *ResContext {
	return &ResContext{}
}

func (s ResContext) Data() Map {
	if s.data == nil {
		s.data = Map{}
	}
	return s.data
}
func (s ResContext) Errors() []error {
	return s.errors
}
func (s ResContext) Has(name string) bool {
	if _, ok := s.data[name]; ok {
		return true
	}
	return false
}

func (s ResContext) Get(name string) interface{} {
	if v, ok := s.data[name]; ok {
		return v
	}
	return nil
}
func (s *ResContext) Set(name string, value interface{}) *ResContext {
	if s.data == nil {
		s.data = Map{}
	}
	s.data[name] = value
	return s
}
func (s *ResContext) SetData(value Map) *ResContext {
	s.data = value
	return s
}

func (s *ResContext) SetError(err interface{}) *ResContext {
	s.errors = make([]error, 0)
	if err != nil {
		s.errors = append(s.errors, ToError(err))
	}
	return s
}

// Error()
// Error(error)
// Error("code is error")
// Error("%s is error","name")
func (s *ResContext) Error(err ...interface{}) error {
	if len(err) > 1 {
		s.errors = append(s.errors, ToError(fmt.Errorf(err[0].(string), err[1:]...)))
	} else if len(err) > 0 {
		s.errors = append(s.errors, ToError(err[0]))
	}
	if len(s.errors) > 0 {
		return s.errors[len(s.errors)-1]
	}
	return nil
}

func (s *ResContext) Bind(c *gin.Context) {
	if len(s.errors) > 0 {
		if _, ok := s.data["code"]; !ok {
			s.Set("code", http.StatusBadRequest)
		}
		errText := make([]string, 0)
		for _, err := range s.errors {
			errText = append(errText, err.Error())
		}
		s.Set("msg", strings.Join(errText, ";"))
		s.Set("errors", errText)

		c.PureJSON(http.StatusBadRequest, s.data)
	} else {
		if _, ok := s.data["code"]; !ok {
			s.Set("code", 200)
		}
		c.PureJSON(http.StatusOK, s.data)
	}
}

func (s *ResContext) Adjust(fn func(res *ResContext)) *ResContext {
	fn(s)
	return s
}
