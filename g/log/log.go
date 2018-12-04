// vim: tabstop=4 shiftwidth=4 softtabstop=4
//
// Copyright 2016 www.meituan.com
//
// File: log.go
// Author: xuriwuyun <xuriwuyun@gmail.com>
// Created: 11/24/2016 17:40:58
// Last_modified: 12/15/2017 12:02:10

package log

import (
	"metric-transfer/g/log/hooks"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var (
	logDir   = flag.String("log-path", "/tmp", "log file path")
	logName  = flag.String("log-name", "", "log file name")
	logLevel = flag.String("log-level", "debug", "set log level: debug,info,warn,error,fatal,panic")
	logger   *logrus.Logger
)

type Fields logrus.Fields
type Entry struct {
	*logrus.Entry
}

func Init() {
	var logfile string
	if *logName == "" {
		file := filepath.Base(os.Args[0])
		logName = &file
	}
	if strings.HasPrefix(*logName, "/") {
		logfile = fmt.Sprintf("%s.log", *logName)
	} else {
		logfile = fmt.Sprintf("%s/%s.log", *logDir, *logName)
	}
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(fmt.Sprintf("open log file failed: %v\n", err))
	}

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		panic(fmt.Sprintf("invalid log level %v(debug,info,warn,error,fatal,panic)\n", *logLevel))
	}

	logger = logrus.New()
	logger.Level = level
	logger.Out = file
	logger.Hooks.Add(new(hooks.CallerHook))
	if level >= logrus.DebugLevel {
		logger.Hooks.Add(new(hooks.StdioHook))
	}

	logger.Formatter = &TextFormatter{
		TimestampFormat: "060102 15:04:05",
		SpacePadding:    0,
	}
	Info("Log Initialized")
}

func GetLogger() *logrus.Logger {
	return logger
}

func LogLevel() string {
	return logger.Level.String()
}

func WithField(key string, value interface{}) *Entry {
	entry := logger.WithField(key, value)
	return &Entry{entry}
}

func WithFields(fields Fields) *Entry {
	entry := logger.WithFields(logrus.Fields(fields))
	return &Entry{entry}
}

func (e *Entry) Debug(args ...interface{}) {
	e.Entry.Debug(args...)
}

func (e *Entry) Print(args ...interface{}) {
	e.Entry.Print(args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Entry.Info(args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Entry.Warn(args...)
}

func (e *Entry) Warning(args ...interface{}) {
	e.Entry.Warning(args...)
}

func (e *Entry) Error(args ...interface{}) {
	e.Entry.Error(args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Entry.Fatal(args...)
}

func (e *Entry) Panic(args ...interface{}) {
	e.Entry.Panic(args...)
}

// Entry Printf family functions

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Entry.Debugf(format, args...)
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Entry.Infof(format, args...)
}

func (e *Entry) Printf(format string, args ...interface{}) {
	e.Entry.Printf(format, args...)
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.Warnf(format, args...)
}

func (e *Entry) Warningf(format string, args ...interface{}) {
	e.Entry.Warningf(format, args...)
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Entry.Errorf(format, args...)
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Entry.Fatalf(format, args...)
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Entry.Panicf(format, args...)
}

// Entry Println family functions

func (e *Entry) Debugln(args ...interface{}) {
	e.Entry.Debugln(args...)
}

func (e *Entry) Infoln(args ...interface{}) {
	e.Entry.Infoln(args...)
}

func (e *Entry) Println(args ...interface{}) {
	e.Entry.Println(args...)
}

func (e *Entry) Warnln(args ...interface{}) {
	e.Entry.Warnln(args...)
}

func (e *Entry) Warningln(args ...interface{}) {
	e.Entry.Warningln(args...)
}

func (e *Entry) Errorln(args ...interface{}) {
	e.Entry.Errorln(args...)
}

func (e *Entry) Fatalln(args ...interface{}) {
	e.Entry.Fatalln(args...)
}

func (e *Entry) Panicln(args ...interface{}) {
	e.Entry.Panicln(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

func Println(args ...interface{}) {
	logger.Println(args...)
}

func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

func Warningln(args ...interface{}) {
	logger.Warningln(args...)
}

func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}
