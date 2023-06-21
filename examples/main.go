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

	//if err := runApp(); err != nil {
	//	os.Exit(1)
	//}
}
