package logger

import (
	"fmt"
	"os"
	"time"
)

type (
	logLevel uint8
	Logger   struct {
		level logLevel
	}
)

func New(level string) *Logger {
	logg := Logger{
		level: errorLevel,
	}

	if levelIndex, ok := titleToLevel[level]; ok {
		logg.level = levelIndex
	}

	return &logg
}

func (l Logger) Fatal(msg string) {
	l.printMessage(msg, fatalLevel)
	os.Exit(1)
}

func (l Logger) Error(msg string) {
	l.printMessage(msg, errorLevel)
}

func (l Logger) Warning(msg string) {
	l.printMessage(msg, warningLevel)
}

func (l Logger) Info(msg string) {
	l.printMessage(msg, infoLevel)
}

func (l Logger) Debug(msg string) {
	l.printMessage(msg, debugLevel)
}

func (l Logger) Fatalf(format string, values ...any) {
	l.Fatal(fmt.Sprintf(format, values...))
}

func (l Logger) Errorf(format string, values ...any) {
	l.Error(fmt.Sprintf(format, values...))
}

func (l Logger) Warningf(format string, values ...any) {
	l.Warning(fmt.Sprintf(format, values...))
}

func (l Logger) Infof(format string, values ...any) {
	l.Info(fmt.Sprintf(format, values...))
}

func (l Logger) Debugf(format string, values ...any) {
	l.Debug(fmt.Sprintf(format, values...))
}

// Напечатать отформатированное сообщение в консоль.
func (l Logger) printMessage(msg string, level logLevel) {
	if level > l.level {
		return
	}

	logTime := time.Now().Format(time.DateTime + ".000")
	fmt.Printf("[%s] [%s] %s \n", logTime, level, msg)
}
