package main

import (
	"mycode/loger"
	"time"
)

func main() {
	f := loger.NewFileWriter("runninglogs.log", "", 10)
	lg := loger.NewLoger(loger.DEBUG, f)
	lg.Debug("这是一条debug信息")
	time.Sleep(time.Second)
	lg.Info("这是一条info信息")
	time.Sleep(time.Second)
	lg.Warning("这是一条warning信息")
	time.Sleep(time.Second)
	lg.Error("这是一条error信息")
	time.Sleep(time.Second)
	lg.Fatal("这是一条fatal信息")
	time.Sleep(time.Second)
	lg.Debug("这是一条debug信息")
	time.Sleep(time.Second)
	lg.Info("这是一条info信息")
	time.Sleep(time.Second)
	lg.Warning("这是一条warning信息")
	time.Sleep(time.Second)
	lg.Error("这是一条error信息")
	time.Sleep(time.Second)
	lg.Fatal("这是一条fatal信息")
	time.Sleep(time.Second)
}
