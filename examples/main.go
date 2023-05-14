package main

import (
	"fmt"
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"os"
	"reflect"
	"time"
)

func main() {
	log.Info().WithForceOutput().Msgf("应用开始进入===============:%v ，启动参数为：%v \n", utils.TimeNow(), os.Args)
	if utils.PathExists("storage/refs/zoneinfo.zip") {
		os.Setenv("ZONEINFO", utils.JoinCurrentPath("storage/refs/zoneinfo.zip"))
		time.LoadLocation("Asia/Shanghai")
	}
	log.WarnF("启动参数：%v", os.Args)
	log.WarnF("环境变量：%v", os.Environ())

	t := utils.ToTime("2022-01-02")
	log.ErrorD(t)

	t = utils.ToTime("2022-1-02 09:34")
	log.ErrorD(t)

	t = utils.ToTime("2022/1-02")
	log.ErrorD(t)
	//flow := utils.NewFlowContext()
	//flow.Request.Action = "query"
	//execCommonRule(newCommonDisable(), flow)

	//if err := runApp(); err != nil {
	//	os.Exit(1)
	//}
}
func execCommonRule(rule rule.IRule, flow *utils.FlowContext) {
	reflectValue := reflect.ValueOf(rule)
	methodName := flow.Request.Action

	//methodName = utils.FirstCaseToUpper(methodName, true)
	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		reflectValue = reflectValue.Addr()
	}
	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		switch method := methodValue.Interface().(type) {
		case func():
			method()
		case func(*utils.FlowContext):
			method(flow)
		default:
			flow.Error(fmt.Errorf("unsupported function %v", methodName))
		}
	}
}

type commonDisable struct {
}

func newCommonDisable() commonDisable {
	return commonDisable{}
}

func (s commonDisable) Register() rule.MDRule {
	return rule.MDRule{Action: "*", Widget: "aa", Sequence: 50}
}
func (s commonDisable) query(flow *utils.FlowContext) {
	if flow.Request.ID == "" {
		flow.Error("缺少 ID 参数！")
		return
	}
}
func (s commonDisable) Exec(flow *utils.FlowContext) {
	if flow.Request.ID == "" {
		flow.Error("缺少 ID 参数！")
		return
	}
}
