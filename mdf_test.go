package mdf

import (
	"github.com/nbkit/mdf/utils"
	"testing"
)

func TestDate_Date(t *testing.T) {
	d := utils.TimeNow().Format(utils.TimeFormatStr("yyyy-MM-dd HH:mm:ss.SSSSSSSS"))
	t.Log(d)

	server := NewServer(Option{
		EnabledFeature: false,
	})

	server.Start()
}
