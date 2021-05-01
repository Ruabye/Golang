package main

import (
	"mycode/loger"
)

func main() {
	f := loger.NewFileWriter("log.log", "", 20)
	lg := loger.NewLoger(loger.ERROR, f)
	lg.Debug("这是一条debug信息")
	lg.Info("这是一条info信息")
	lg.Warning("这是一条warning信息")
	lg.Error("这是一条error信息")
	lg.Fatal("这是一条fatal信息")
}
