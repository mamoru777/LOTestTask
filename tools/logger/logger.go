package logger

import (
	"fmt"
	"io"
	"os"
)

type Config struct {
	LogLevel LogLevel
}

type consoleLog struct {
	message  any
	layer    string
	logLevel LogLevel
}

type Logger struct {
	logCh    chan consoleLog
	logLevel LogLevel
	w        io.Writer
}

func InitLogger(cfg Config, w io.Writer) *Logger {
	logger := &Logger{
		logCh:    make(chan consoleLog),
		w:        w,
		logLevel: cfg.LogLevel,
	}

	go func() {
		for log := range logger.logCh {
			logger.handle(log)
		}
	}()

	return logger
}

func (l *Logger) Error(log any, layer string) {
	l.logCh <- consoleLog{message: log, logLevel: LevelError, layer: layer}
}

func (l *Logger) Warning(log any, layer string) {
	l.logCh <- consoleLog{message: log, logLevel: LevelWarning, layer: layer}
}

func (l *Logger) Info(log any, layer string) {
	l.logCh <- consoleLog{message: log, logLevel: LevelInfo, layer: layer}
}

func (l *Logger) Fatal(log any, layer string) {
	l.logCh <- consoleLog{message: log, logLevel: LevelFatal, layer: layer}
}

func (l *Logger) Debug(log any, layer string) {
	l.logCh <- consoleLog{message: log, logLevel: LevelInfo, layer: layer}
}

func (l *Logger) handle(log consoleLog) {
	if l.logLevel.GreaterThan(log.logLevel) {
		return
	}

	var msg string

	switch v := log.message.(type) {
	case error:
		msg = v.Error()
	default:
		msg = fmt.Sprintf("%v", v)
	}

	buf := make([]byte, 0)

	var delimer byte = ' '

	buf = append(buf, getColor(log.logLevel)...)
	buf = append(buf, fmt.Sprintf("[%s]", log.logLevel.ToUpper())...)
	buf = append(buf, colorReset...)
	buf = append(buf, delimer)

	buf = append(buf, fmt.Sprintf("[%s]", log.layer)...)
	buf = append(buf, delimer)

	buf = append(buf, msg...)
	buf = append(buf, '\n')

	_, err := l.w.Write(buf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not write log: %s\n", err)
	}
}

func getColor(level LogLevel) string {
	switch level {
	case LevelFatal:
		return colorMagenta
	case LevelError:
		return colorRed
	case LevelWarning:
		return colorYellow
	case LevelInfo:
		return colorBlue
	case LevelDebug:
		return colorCyan
	default:
		return colorWhite
	}
}

func (l *Logger) Close() {
	close(l.logCh)
}
