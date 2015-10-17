package logger

import (
	"log"
	"util"
)

var logger *util.Logger

func SetLogger(customLogger *util.Logger) {
	logger = customLogger
}

func Info(content ...interface{}) {
	if logger == nil {
		log.Print(content...)
	} else {
		logger.Info(content...)
	}
}

func Infof(format string, content ...interface{}) {
	if logger == nil {
		log.Printf(format, content...)
	} else {
		logger.Infof(format, content...)
	}
}

func Warning(content ...interface{}) {
	if logger == nil {
		log.Print(content...)
	} else {
		logger.Warning(content...)
	}
}

func Warningf(format string, content ...interface{}) {
	if logger == nil {
		log.Printf(format, content...)
	} else {
		logger.Warningf(format, content...)
	}
}

func Notice(content ...interface{}) {
	if logger == nil {
		log.Print(content...)
	} else {
		logger.Notice(content...)
	}
}

func Noticef(format string, content ...interface{}) {
	if logger == nil {
		log.Printf(format, content...)
	} else {
		logger.Noticef(format, content...)
	}
}

func Debug(content ...interface{}) {
	if logger == nil {
		log.Print(content...)
	} else {
		logger.Debug(content...)
	}
}

func Debugf(format string, content ...interface{}) {
	if logger == nil {
		log.Printf(format, content...)
	} else {
		logger.Debugf(format, content...)
	}
}
