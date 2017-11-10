package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type SimpleLogger struct {
	//file log
	fileLevel  int
	fileLogger *log.Logger
	baseFile   *os.File
	outputPath string
	maxLine    int
	logLine    int
	//shell log
	shellLevel int
	isColored  bool
	buffer     chan *Msg
}

const (
	COLOR_DEBUG_LEVEL_DESC = "[\x1b[32mdebug\x1b[0m] "
	COLOR_INFO_LEVEL_DESC  = "[\x1b[36minfo \x1b[0m] "
	COLOR_WARN_LEVEL_DESC  = "[\x1b[33mwarn \x1b[0m] "
	COLOR_ERROR_LEVEL_DESC = "[\x1b[31merror\x1b[0m] "
	COLOR_FATAL_LEVEL_DESC = "[\x1b[31mfatal\x1b[0m] "
)

var color_format = []string{
	COLOR_DEBUG_LEVEL_DESC,
	COLOR_INFO_LEVEL_DESC,
	COLOR_WARN_LEVEL_DESC,
	COLOR_ERROR_LEVEL_DESC,
	COLOR_FATAL_LEVEL_DESC,
}

var std_format = []string{
	DEBUG_LEVEL_DESC,
	INFO_LEVEL_DESC,
	WARN_LEVEL_DESC,
	ERROR_LEVEL_DESC,
	FATAL_LEVEL_DESC,
}

func (self *SimpleLogger) SetColored(colored bool) {
	self.isColored = colored
}

func (self *SimpleLogger) doPrintf(level int, levelDesc, format string, a ...interface{}) {
	nformat := levelDesc + format
	if level >= self.fileLevel {
		if self.logLine > self.maxLine {
			if self.baseFile != nil {
				self.baseFile.Sync()
				self.baseFile.Close()
			}
			self.baseFile, _ = self.createLogFile(self.outputPath)
			self.fileLogger = log.New(self.baseFile, "", log.Ldate|log.Lmicroseconds)
			self.logLine = 0
		}
		if self.fileLogger != nil {
			self.fileLogger.Printf(nformat, a...)
			self.logLine++
		}
	}
	if level >= self.shellLevel {
		sel_fmt := color_format
		if !self.isColored {
			sel_fmt = std_format
		}
		nformat := sel_fmt[level] + format
		log.Printf(nformat, a...)
	}
}

func (self *SimpleLogger) DoPrintf(level int, levelDesc, format string, a ...interface{}) {
	if self.buffer == nil {
		self.doPrintf(level, levelDesc, format, a...)
		return
	}
	self.buffer <- &Msg{level, levelDesc, format, a}
}

func (self *SimpleLogger) setFileOutDir(path string) {
	self.outputPath = path
}

func (self *SimpleLogger) setFileLogLevel(level int) {
	self.fileLevel = level
}
func (self *SimpleLogger) setShellLogLevel(level int) {
	self.shellLevel = level
}

func (self *SimpleLogger) createLogFile(dir string) (*os.File, error) {
	now := time.Now()
	filename := fmt.Sprintf("_%d-%02d-%02d_%02d-%02d-%02d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
	if dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}
	file, err := os.Create(path.Join(dir, filename))
	if err != nil {
		log.Printf("create file failed %v", err)
		return nil, err
	}
	return file, nil
}

func (self *SimpleLogger) run() {
	go func() {
		for {
			m, ok := <-self.buffer
			if ok {
				self.doPrintf(m.level, m.levelDesc, m.fmt, m.param...)
			}
		}
	}()
}

func CreateLogger(path string, fileLevel, shellLevel, maxLine, bufSize int) *SimpleLogger {
	logger := &SimpleLogger{}
	logger.setFileOutDir(path)
	logger.setFileLogLevel(fileLevel)
	logger.setShellLogLevel(shellLevel)
	logger.maxLine = maxLine
	logger.logLine = maxLine + 1
	if bufSize > 10 {
		log.Printf("log start with async mode.")
		logger.buffer = make(chan *Msg, bufSize)
		logger.run()
	}
	return logger
}

var DefaultLogger = CreateLogger("test", FATAL_LEVEL, DEBUG_LEVEL, 10000, 1000)
