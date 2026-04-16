package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger = logrus.New()
var GlobalFields = map[string]interface{}{}

// this is the entry point for the logger as it initializes the logger
func InitLogger(logType string) {

	// Set logger format
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
		ForceQuote:      true,
	})

	setLoggerLevel(viper.GetString("LOG_LEVEL"))
	GlobalFields = map[string]interface{}{
		"log_type": logType,
	}
}

func setLoggerLevel(l string) {
	l = strings.ToUpper(l)
	switch l {
	case "PANIC":
		Logger.SetLevel(logrus.PanicLevel)
	case "FATAL":
		Logger.SetLevel(logrus.FatalLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "TRACE":
		Logger.SetLevel(logrus.TraceLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

func Info(msg string) {
	Logger.Info(msg)
}

func Debug(msg string) {
	Logger.Debug(msg)
}

func Error(err error, msg string) {
	Logger.WithError(errors.WithStack(err)).Error(msg)
}

func ErrorWithoutSentry(err error, msg string) {
	Logger.WithError(errors.WithStack(err)).Error(msg)
}

func Warn(msg string) {
	Logger.Warn(msg)
}

func Fatal(err error, msg string) {
	Logger.WithError(errors.WithStack(err)).Fatal(msg)
}

func Panic(err error, msg string) {
	Logger.WithError(errors.WithStack(err)).Panic(msg)
}

func CreateLogMsg(message ...string) string {
	otrMessage := strings.Join(message, " | ")
	return fmt.Sprintf("%s | %s", strings.ToUpper(viper.GetString("APP_ENV")), otrMessage)
}
