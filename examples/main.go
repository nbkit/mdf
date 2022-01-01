package main

import (
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"os"
	"time"
)

func main() {
	if utils.PathExists("storage/refs/zoneinfo.zip") {
		os.Setenv("ZONEINFO", utils.JoinCurrentPath("storage/refs/zoneinfo.zip"))
		time.LoadLocation("Asia/Shanghai")
	}
	log.WarnF("启动参数：%v", os.Args)
	log.WarnF("环境变量：%v", os.Environ())
	log.Info().Msgf("应用开始进入===============:%v ，启动参数为：%v \n", utils.TimeNow(), os.Args)
	if err := runApp(); err != nil {
		os.Exit(1)
	}
}
