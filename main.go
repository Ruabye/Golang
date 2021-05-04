package main

import (
	log "mycode/log2"
	"time"
)

func main() {
	lg := log.NewLog("", "log", 100, log.INFO)
	lg.Debug("这是一条debug信息")
	lg.Info("这是一条info信息")
	lg.Warning("这是一条warning信息")
	lg.Error("这是一条error信息")
	lg.Fatal("这是一条fatal信息")
	time.Sleep(3 * time.Second)
}
