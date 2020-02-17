package config

import (
	"bufio"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"os"
	"strings"
)
import log "github.com/sirupsen/logrus"
import "github.com/rifflock/lfshook"
import "time"

const suffix string = "log"
const split string = "/"

// 创建目录
func createFolder(logPath string) {
	lastSplit := strings.LastIndex(logPath, split)
	if lastSplit != -1 {
		folder := logPath[0:lastSplit]
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"folder": folder,
			}).Error("create log folder failure")
		}
	}
}

// 创建日志扩展
func newLfsHook(logPath string) log.Hook {
	// 创建目录，防止日志出现在目录里
	createFolder(logPath)

	// 创建writer
	writer, err := rotatelogs.New(
		fmt.Sprintf("%v%v.%v", logPath, "%Y%m%d%H", suffix),

		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(fmt.Sprintf("%v.%v", logPath, suffix)),

		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour*24),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		// rotatelogs.WithMaxAge(time.Hour*24),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
		return nil
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{TimestampFormat: SystemConfig.LoggerConf.TimestampFormat})
	return lfsHook
}

// 初始化日志
func InitLogger() {
	// 设置日志格式
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: SystemConfig.LoggerConf.TimestampFormat})

	// 创建扩展，日志切割
	hook := newLfsHook(SystemConfig.LoggerConf.LogPath)
	if nil != hook {
		log.AddHook(hook)
	}

	log.Info("初始化日志服务")
	log.SetLevel(parseLevel(SystemConfig.LoggerConf.LogLevel))

	// 关闭标准输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Error("Open Src File err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)
}

func parseLevel(level string) log.Level {
	if level == "" {
		return log.InfoLevel
	}

	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.Info("日志级别错误，默认使用debug模式")
		return log.InfoLevel
	}
	return logLevel
}
