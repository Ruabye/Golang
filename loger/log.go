package loger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

type LOGTYPE [2]interface{}

var (
	DEBUG   LOGTYPE = LOGTYPE{uint8(0), "DEBUG"}
	INFO    LOGTYPE = LOGTYPE{uint8(1), "INFO"}
	WARNING LOGTYPE = LOGTYPE{uint8(2), "WARNING"}
	ERROR   LOGTYPE = LOGTYPE{uint8(3), "ERROR"}
	FATAL   LOGTYPE = LOGTYPE{uint8(4), "FATAL"}
)

var (
	STD io.Writer = os.Stdout
	FIL io.Writer = nil
)

type Loger struct {
	Level  uint8
	Output io.Writer
}
type FileWriter struct {
	file     io.Writer
	FileName string
	FilePath string
	MaxSize  int64
}

func NewFileWriter(name, filepath string, msize int64) *FileWriter {
	fileObj, err := os.OpenFile(path.Join(filepath, name), os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return &FileWriter{
		file:     fileObj,
		FileName: name,
		FilePath: filepath,
		MaxSize:  msize,
	}
}
func (f *FileWriter) Write(p []byte) (n int, err error) {
	//文件是否需要切割
	fileInfo, err := os.Stat(path.Join(f.FilePath, f.FileName))
	if err != nil {
		//文件不存在则创建
		if os.IsExist(err) {
			fmt.Println("获取文件信息失败", err)
			return 0, nil
		}
	}
	if fileInfo.Size() > f.MaxSize {
		f = NewFileWriter("re"+f.FileName, f.FilePath, f.MaxSize)
		n, err = f.Write(p)
		return
	} else {
		n, err = fmt.Fprint(f.file, string(p))
		return
	}

}
func NewLoger(level LOGTYPE, out io.Writer) Loger {
	return Loger{
		//接口需要使用类型断言
		Level:  level[0].(uint8),
		Output: out,
	}
}

func (l Loger) Debug(fmtStr string, args ...interface{}) {
	l.writeLog(DEBUG, fmtStr, args...)
}
func (l Loger) Info(fmtStr string, args ...interface{}) {
	l.writeLog(INFO, fmtStr, args...)
}
func (l Loger) Warning(fmtStr string, args ...interface{}) {
	l.writeLog(WARNING, fmtStr, args...)
}
func (l Loger) Error(fmtStr string, args ...interface{}) {
	l.writeLog(ERROR, fmtStr, args...)
}
func (l Loger) Fatal(fmtStr string, args ...interface{}) {
	l.writeLog(FATAL, fmtStr, args...)
}

func (l Loger) writeLog(level LOGTYPE, fmtStr string, args ...interface{}) {
	if level[0].(uint8) < l.Level {
		return
	}
	cur := time.Now().Format("2006-01-02 15:04")
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		fmt.Println("找不到错误信息")
	}
	position := fmt.Sprintf("functionName：%s, fileName：%s, line: %d", runtime.FuncForPC(pc).Name(), path.Base(file), line)
	infos := fmt.Sprintf(fmtStr, args...)
	msg := fmt.Sprintf("[%s] [%s] [%s]: %s", level[1].(string), cur, position, infos)
	fmt.Fprintln(l.Output, msg)
}
