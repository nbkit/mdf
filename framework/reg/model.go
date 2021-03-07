package reg

import (
	"strings"

	"github.com/nbkit/mdf/utils"
)

type RegObject struct {
	Code    string           `json:"code"`
	Name    string           `json:"name"`
	Address string           `json:"address"`
	Content string           `json:"content"`
	Time    *utils.Time      `json:"time"`
	Configs *utils.EnvConfig `json:"configs"`
}

func (s RegObject) Key() string {
	return strings.ToLower(s.Code)
}
