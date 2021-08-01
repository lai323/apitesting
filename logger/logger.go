package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitLogger() {
	level := convStrLevel(viper.GetString("log.level"))
	stdout := viper.GetBool("log.stdout")
	format := viper.GetString("log.format")
	logdir := viper.GetString("log.logdir")
	logfile := viper.GetString("log.logfile")
	if logdir == "" {
		logdir = "./log"
	}
	if logfile == "" {
		logfile = "log"
	}
	logrus.SetLevel(level)

	fileSrc, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	var src io.Writer
	if stdout {
		stdoutSrc := os.Stdout
		src = io.MultiWriter(fileSrc, stdoutSrc)
	} else {
		src = io.MultiWriter(fileSrc)
	}
	logrus.SetOutput(src)

	logrus.SetFormatter(getformatter(format))
	logrus.AddHook(getLogHook(format, logdir, logfile))
	return
}

func newloggerErr(format string, a ...interface{}) error {
	return fmt.Errorf("NewLogger error: %s", fmt.Sprintf(format, a...))
}

func convStrLevel(l string) logrus.Level {
	var level = logrus.DebugLevel
	switch l {
	case "panic":
		level = logrus.PanicLevel
	case "fatal":
		level = logrus.FatalLevel
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	}
	return level
}

func NewLogger(loglevel, format string, stdout bool, logdir, logfile string) (*logrus.Logger, error) {
	var logger = logrus.New()
	var level = logrus.DebugLevel
	switch loglevel {
	case "panic":
		level = logrus.PanicLevel
	case "fatal":
		level = logrus.FatalLevel
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	default:
		return logger, newloggerErr(`NewLogger error log level not allow: %s allow: "panic" "fatal" "error" "warn" "info" "debug" "trace", current:`, loglevel)
	}
	logger.SetLevel(level)

	fileSrc, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return logger, newloggerErr(err.Error())
	}
	var src io.Writer
	if stdout {
		stdoutSrc := os.Stdout
		src = io.MultiWriter(fileSrc, stdoutSrc)
	} else {
		src = io.MultiWriter(fileSrc)
	}
	logger.SetOutput(src)

	logger.SetFormatter(getformatter(format))
	logger.AddHook(getLogHook(format, logdir, logfile))
	return logger, nil
}

func getformatter(f string) logrus.Formatter {
	switch f {
	case "text":
		return &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		}
	case "json":
		return &logrus.JSONFormatter{}
	case "text_disable_fulltimestamp":
		fallthrough
	default:
		return &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: false,
		}
	}
}

func getLogHook(format, logdir, logfile string) *lfshook.LfsHook {
	logWriter, _ := getLogWriter(logdir, logfile)
	writeMap := lfshook.WriterMap{
		logrus.TraceLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	return lfshook.NewHook(writeMap, getformatter(format))
}

func getLogWriter(logdir, logfile string) (*rotatelogs.RotateLogs, error) {
	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(logdir, os.ModePerm)
			if err != nil {
				return nil, newloggerErr("create log dir (%s) error %s", logdir, err.Error())
			}
		}
		return nil, newloggerErr(err.Error())
	}

	filepath := path.Join(logdir, logfile)
	logWriter, err := rotatelogs.New(
		filepath+".%Y%m%d.log",
		rotatelogs.WithLinkName(filepath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, newloggerErr(err.Error())
	}
	return logWriter, nil
}
