package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

// levels
const (
	DEBUG_LEVEL = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
	FATAL_LEVEL
)

const (
	DEBUG_LEVEL_DESC = "[debug] "
	INFO_LEVEL_DESC  = "[info ] "
	WARN_LEVEL_DESC  = "[warn ] "
	ERROR_LEVEL_DESC = "[error] "
	FATAL_LEVEL_DESC = "[fatal] "
)

type Logger struct {
	//file log
	fileLevel  int
	fileLogger *log.Logger
	baseFile   *os.File
	outputPath string
	maxLine    int
	logLine    int

	//shell log
	shellLevel int
	buffer     chan *Msg
}

type Msg struct {
	level     int
	levelDesc string
	fmt       string
	param     []interface{}
}

var glogger *Logger

func send(level int, desc, format string, param ...interface{}) {
	if glogger == nil {
		return
	}
	m := &Msg{level, desc, format, param}
	glogger.buffer <- m
}

func Debug(format string, param ...interface{}) {
	send(DEBUG_LEVEL, DEBUG_LEVEL_DESC, format, param...)
}

func Info(format string, param ...interface{}) {
	send(INFO_LEVEL, INFO_LEVEL_DESC, format, param...)
}

func Warn(format string, param ...interface{}) {
	send(WARN_LEVEL, WARN_LEVEL_DESC, format, param...)
}

func Error(format string, param ...interface{}) {
	send(ERROR_LEVEL, ERROR_LEVEL_DESC, format, param...)
}

func Fatal(format string, param ...interface{}) {
	send(FATAL_LEVEL, FATAL_LEVEL_DESC, format, param...)
}

func createLogFile(dir string) (*os.File, error) {
	now := time.Now()
	filename := fmt.Sprintf("_%d-%02d-%02d_%02d-%02d-%02d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
	if dir != "" {
		os.MkdirAll(dir, os.ModeDir)
	}
	file, err := os.Create(path.Join(dir, filename))
	if err != nil {
		log.Printf("create file failed %v", err)
		return nil, err
	}
	return file, nil
}

func (self *Logger) setFileOutDir(path string) {
	self.outputPath = path
}

func (self *Logger) setFileLogLevel(level int) {
	self.fileLevel = level
}
func (self *Logger) setShellLogLevel(level int) {
	self.shellLevel = level
}

func Init(path string, fileLevel, shellLevel, maxLine, bufSize int) {
	if glogger != nil {
		return
	}
	glogger = &Logger{}
	glogger.setFileOutDir(path)
	glogger.setFileLogLevel(fileLevel)
	glogger.setShellLogLevel(shellLevel)
	glogger.maxLine = maxLine
	glogger.logLine = maxLine + 1
	glogger.buffer = make(chan *Msg, bufSize)

	glogger.run()
}

func (self *Logger) doPrintf(level int, levelDesc, format string, a ...interface{}) {
	format = levelDesc + format
	if level >= self.fileLevel {
		if self.logLine > self.maxLine {
			if self.baseFile != nil {
				self.baseFile.Sync()
				self.baseFile.Close()
			}
			self.baseFile, _ = createLogFile(self.outputPath)
			self.fileLogger = log.New(self.baseFile, "", log.LstdFlags|log.Lmicroseconds)
			self.logLine = 0
		}
		if self.fileLogger != nil {
			self.fileLogger.Printf(format, a...)
			self.logLine++
		}
	}
	if level >= self.shellLevel {
		log.Printf(format, a...)
	}
}

func (self *Logger) run() {
	go func() {
		for {
			m, ok := <-self.buffer
			if ok {
				self.doPrintf(m.level, m.levelDesc, m.fmt, m.param...)
			}
		}
	}()
}
