package util

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	FilenameSuffixInDay    = "20060102"
	FilenameSuffixInSecond = "20060102150405"
	StandardLogPrefix      = "2006/01/02 15:04:05"
)

const (
	debugLevel = iota
	infoLevel
	noticeLevel
	warningLevel
	errorLevel
)

const (
	nocolor = 0
	red     = 30 + iota
	green
	yellow
	blue
	purple
	cyan
)

var (
	logPrefix = map[int]string{
		debugLevel:   "DEBUG",
		infoLevel:    "INFO",
		noticeLevel:  "NOTICE",
		warningLevel: "WARNING",
		errorLevel:   "ERROR",
	}
	logColor = map[int]int{
		debugLevel:   cyan,
		infoLevel:    nocolor,
		noticeLevel:  green,
		warningLevel: yellow,
		errorLevel:   red,
	}
)

type Logger struct {
	sync.Mutex
	filenamePrefix        string
	filenameSuffixFormat  string
	currentFilenameSuffix string
	fileWriter            io.WriteCloser
}

func NewLogger(filenamePrefix string, filenameSuffixFormat string) *Logger {
	fileWriter := getFileWriter(filenamePrefix, filenameSuffixFormat)
	return &Logger{
		filenamePrefix:        filenamePrefix,
		filenameSuffixFormat:  filenameSuffixFormat,
		currentFilenameSuffix: time.Now().Format(filenameSuffixFormat),
		fileWriter:            fileWriter,
	}
}

func (l *Logger) updateInnerLogger(now time.Time) {
	filenameSuffix := now.Format(l.filenameSuffixFormat)
	if filenameSuffix != l.currentFilenameSuffix {
		l.fileWriter.Close()
		l.fileWriter = getFileWriter(l.filenamePrefix, l.filenameSuffixFormat)
		l.currentFilenameSuffix = filenameSuffix
	}
}

func (l *Logger) write(level int, format string, content ...interface{}) {
	now := time.Now()

	var s string
	if format == "" {
		s = renderColor(fmt.Sprintf("%s\t[%s]\t%s\n", now.Format(StandardLogPrefix), logPrefix[level], fmt.Sprint(content...)),
			logColor[level])
	} else {
		s = renderColor(fmt.Sprintf("%s\t[%s]\t%s\n", now.Format(StandardLogPrefix), logPrefix[level], fmt.Sprintf(format, content...)),
			logColor[level])
	}

	l.Lock()
	defer l.Unlock()
	l.updateInnerLogger(now)
	l.fileWriter.Write([]byte(s))
}

func (l *Logger) Info(content ...interface{}) {
	l.write(infoLevel, "", content...)
}

func (l *Logger) Infof(format string, content ...interface{}) {
	l.write(infoLevel, format, content...)
}

func (l *Logger) Warning(content ...interface{}) {
	l.write(warningLevel, "", content...)
}

func (l *Logger) Warningf(format string, content ...interface{}) {
	l.write(warningLevel, format, content...)
}

func (l *Logger) Notice(content ...interface{}) {
	l.write(noticeLevel, "", content...)
}

func (l *Logger) Noticef(format string, content ...interface{}) {
	l.write(noticeLevel, format, content...)
}

func (l *Logger) Debug(content ...interface{}) {
	l.write(debugLevel, "", content...)
}

func (l *Logger) Debugf(format string, content ...interface{}) {
	l.write(debugLevel, format, content...)
}

func (l *Logger) Error(content ...interface{}) {
	l.write(errorLevel, "", content...)
}

func (l *Logger) Errorf(format string, content ...interface{}) {
	l.write(errorLevel, format, content...)
}

func renderColor(s string, color int) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
}

func getFileWriter(filenamePrefix string, filenameSuffixFormat string) io.WriteCloser {
	filenameSuffix := time.Now().Format(filenameSuffixFormat)
	fileWriter, err := os.OpenFile(fmt.Sprintf("%s.%s", filenamePrefix, filenameSuffix), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return fileWriter
}
