package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	gormlogger "gorm.io/gorm/logger"
)

var (
	ljLogger *lumberjack.Logger
)

func Init(logLevel string, logFile string, maxAge int, maxSize int) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if logFile != "" {
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		} else {
			ljLogger = &lumberjack.Logger{
				Filename: logFile,
				MaxAge:   maxAge,
				MaxSize:  maxSize,
			}
			writers = append(writers, ljLogger)
		}
	}

	logrus.SetOutput(io.MultiWriter(writers...))
}

func Close() {
	if ljLogger != nil {
		ljLogger.Close()
	}
}

type gormLogrusWriter struct{}

func (w *gormLogrusWriter) Printf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func NewGormLogger(level gormlogger.LogLevel) gormlogger.Interface {
	return gormlogger.New(
		&gormLogrusWriter{},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
