package main

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"super-markdown-editor-web/agent"
	"super-markdown-editor-web/model"
	"time"
)

func init() {
	if _, err := os.Stat("./log"); err != nil {
		os.MkdirAll("./log", 0766)
	}
	baseLogPath := path.Join("./log",
		"agent.log")
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}
	formatter := &logrus.TextFormatter{
		// 不需要彩色日志
		DisableColors: false,
		// 定义时间戳格式
		TimestampFormat: "2006-01-02 15:04:05",
	}
	lfHook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}, formatter)
	logrus.AddHook(lfHook)
}

func main() {
	logrus.SetReportCaller(true)
	model.GetLocalEditorConf("./config.json")
	agent.RunAgent(true, model.GetLocalEditorConf(""))
}
