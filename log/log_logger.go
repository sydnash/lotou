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

func (self *SimpleLoger) DoPrintf(level int, levelDesc, format string, a ...interface{}) {
	format = levelDesc + format
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
			self.fileLogger.Printf(format, a...)
			self.logLine++
		}
	}
	if level >= self.shellLevel {
		log.Printf(format, a...)
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
