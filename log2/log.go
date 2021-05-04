package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

type Level int

const MaxChanSize = 100
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type Log struct {
	Filename      string
	Path          string
	MaxSize       int
	MinLevel      Level
	LogChan       chan *LogInfo
	FileObj       *os.File
	HighLvFileObj *os.File
}
type LogInfo struct {
	Filename string
	Funcname string
	Msg      string
	Nowtime  string
	LineNum  int
	Lv       string
}

func NewLog(basePath, fileName string, maxSize int, minLevel Level) Log {
	fileObj, highLvFileObj := initFile(basePath, fileName)
	logChan := make(chan *LogInfo, MaxChanSize)
	logObj := Log{
		Filename:      fileName,
		Path:          basePath,
		MaxSize:       maxSize,
		MinLevel:      minLevel,
		LogChan:       logChan,
		FileObj:       fileObj,
		HighLvFileObj: highLvFileObj,
	}
	go logObj.writeLog()
	return logObj
}

func initFile(basePath, name string) (*os.File, *os.File) {
	fullPath := path.Join(basePath, name+".log")
	fileObj, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic("获取日志文件句柄失败")
	}
	fullPath = path.Join(basePath, name+".err.log")
	highLvFileObj, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic("获取日志文件句柄失败")
	}
	return fileObj, highLvFileObj
}

func (l *Log) writeLog() {
	for loginfo := range l.LogChan {
		//记录低等级日志
		fileInfo, err := l.FileObj.Stat()
		if err != nil {
			fmt.Printf("获取文件信息失败,err : %v", err)
		}
		rd := fmt.Sprintf("[%s][%s][file:%s func:%s line:%d]:%s", loginfo.Lv, loginfo.Nowtime, loginfo.Filename, loginfo.Funcname, loginfo.LineNum, loginfo.Msg)
		if fileInfo.Size() > int64(l.MaxSize) {
			//文件太大，需要切割
			oldName := path.Join(l.Path, l.Filename+".log")
			timestamp := time.Now().Unix()
			newName := fmt.Sprintf("%s.bak%d", oldName, timestamp)
			l.FileObj.Close()
			os.Rename(oldName, newName)
			l.FileObj, err = os.OpenFile(oldName, os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				panic("获取日志文件句柄失败")
			}
		}
		//写入日志
		fmt.Fprintln(l.FileObj, rd)
		//记录高等级日志
		if loginfo.Lv == "ERROR" || loginfo.Lv == "FATAL" {
			fileInfo, err = l.HighLvFileObj.Stat()
			if err != nil {
				fmt.Printf("获取文件信息失败,err : %v", err)
			}
			if fileInfo.Size() > int64(l.MaxSize) {
				//文件太大，需要切割
				oldName := path.Join(l.Path, l.Filename+".err.log")
				timestamp := time.Now().Unix()
				newName := fmt.Sprintf("%s.bak%d", oldName, timestamp)
				l.FileObj.Close()
				os.Rename(oldName, newName)
				l.FileObj, err = os.OpenFile(oldName, os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					panic("获取日志文件句柄失败")
				}
			}
			//写入日志
			fmt.Fprintln(l.HighLvFileObj, rd)
		}
	}
}

func (l Log) record(msg, lv string) {
	//记录错误信息到通道中
	current := time.Now().Format("2006/01/02 15:04:05")
	pc, file, line, ok := runtime.Caller(2)
	fmt.Println(runtime.FuncForPC(pc).Name(), file, line)
	if !ok {
		fmt.Println("获取错误信息失败")
	}
	// Filename string
	// Funcname string
	// LineNum  int
	filename := file
	funcname := runtime.FuncForPC(pc).Name()
	loginfo := &LogInfo{
		Filename: path.Base(filename),
		Funcname: funcname,
		LineNum:  line,
		Lv:       lv,
		Nowtime:  current,
		Msg:      msg,
	}
	select {
	case l.LogChan <- loginfo:
	default:
		//在通道已满时这个线程可能会阻塞。使用select方法，通道已满时，执行default将新的错误信息丢弃
	}
}

func (l Log) Debug(msg string) {
	//判断是否比最低等级高
	if l.MinLevel <= DEBUG {
		//交给一个goroutine，让它生成错误信息，并记录在通道中
		go l.record(msg, "DEBUG")
	}
}

func (l Log) Info(msg string) {
	//判断是否比最低等级高
	if l.MinLevel <= INFO {
		//交给一个goroutine，让它生成错误信息，并记录在通道中
		l.record(msg, "INFO")
	}
}

func (l Log) Warning(msg string) {
	//判断是否比最低等级高
	if l.MinLevel <= WARNING {
		//交给一个goroutine，让它生成错误信息，并记录在通道中
		l.record(msg, "WARNING")
	}
}

func (l Log) Error(msg string) {
	//判断是否比最低等级高
	if l.MinLevel <= ERROR {
		//交给一个goroutine，让它生成错误信息，并记录在通道中
		l.record(msg, "ERROR")
	}
}

func (l Log) Fatal(msg string) {
	//判断是否比最低等级高
	if l.MinLevel <= FATAL {
		//交给一个goroutine，让它生成错误信息，并记录在通道中
		l.record(msg, "FATAL")
	}
}
