package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type SimpleLoger struct {
	//file log
	fileLevel  int
	fileLogger *log.Logger
	baseFile   *os.File
	outputPath string
	maxLine    int
	logLine    int
	//shell log
	shellLevel int
}

const (
	COLOR_DEBUG_LEVEL_DESC = "[\x1b[32mdebug\x1b[0m] "
	COLOR_INFO_LEVEL_DESC  = "[\x1b[36minfo \x1b[0m] "
	COLOR_WARN_LEVEL_DESC  = "[\x1b[33mwarn \x1b[0m] "
	COLOR_ERROR_LEVEL_DESC = "[\x1b[31merror\x1b[0m] "
	COLOR_FATAL_LEVEL_DESC = "[\x1b[31mfatal\x1b[0m] "
)

var color_format = []string{COLOR_DEBUG_LEVEL_DESC, COLOR_INFO_LEVEL_DESC, COLOR_WARN_LEVEL_DESC, COLOR_ERROR_LEVEL_DESC, COLOR_FATAL_LEVEL_DESC}

func (self *SimpleLoger) DoPrintf(level int, levelDesc, format string, a ...interface{}) {
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
		nformat := color_format[level] + format
		log.Printf(nformat, a...)
	}
}

func (self *SimpleLoger) setFileOutDir(path string) {
	self.outputPath = path
}

func (self *SimpleLoger) setFileLogLevel(level int) {
	self.fileLevel = level
}
func (self *SimpleLoger) setShellLogLevel(level int) {
	self.shellLevel = level
}

func (self *SimpleLoger) createLogFile(dir string) (*os.File, error) {
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

func CreateLogger(path string, fileLevel, shellLevel, maxLine, bufSize int) *SimpleLoger {
	logger := &SimpleLoger{}
	logger.setFileOutDir(path)
	logger.setFileLogLevel(fileLevel)
	logger.setShellLogLevel(shellLevel)
	logger.maxLine = maxLine
	logger.logLine = maxLine + 1
	return logger
}

var DefaultLogger = CreateLogger("test", FATAL_LEVEL, DEBUG_LEVEL, 10000, 1000)
