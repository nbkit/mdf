package main

import (
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"os"
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

	//if err := runApp(); err != nil {
	//	os.Exit(1)
	//}
}
