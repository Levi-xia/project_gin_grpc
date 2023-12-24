package base

import (
	"os"
	"strings"
	"time"
	"com.levi/project-common/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LOG *zap.Logger

type LogConfig struct {
	DebugFileName string `json:"debugFileName"`
	InfoFileName  string `json:"infoFileName"`
	WarnFileName  string `json:"warnFileName"`
	ErrorFileName string `json:"errorFileName"`
	MaxSize       int    `json:"maxsize"`
	MaxAge        int    `json:"max_age"`
	MaxBackups    int    `json:"max_backups"`
}

func InitLog() (err error) {
	// 初始化日志
	curApp := GetCurApp()
	cfg := &LogConfig{
		DebugFileName: strings.ReplaceAll(config.GlobalConf.Zap.DebugFileName, "{app}", curApp),
		InfoFileName:  strings.ReplaceAll(config.GlobalConf.Zap.InfoFileName, "{app}", curApp),
		WarnFileName:  strings.ReplaceAll(config.GlobalConf.Zap.WarnFileName, "{app}", curApp),
		ErrorFileName: strings.ReplaceAll(config.GlobalConf.Zap.ErrorFileName, "{app}", curApp),
		MaxSize:       config.GlobalConf.Zap.MaxSize,
		MaxAge:        config.GlobalConf.Zap.MaxAge,
		MaxBackups:    config.GlobalConf.Zap.MaxBackups,
	}
	writeSyncerDebug := getLogWriter(cfg.DebugFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerInfo := getLogWriter(cfg.InfoFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerWarn := getLogWriter(cfg.WarnFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerError := getLogWriter(cfg.ErrorFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)

	encoder := getEncoder()
	// 文件输出
	debugCore := zapcore.NewCore(encoder, writeSyncerDebug, zapcore.DebugLevel)
	infoCore := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel)
	warnCore := zapcore.NewCore(encoder, writeSyncerWarn, zapcore.WarnLevel)
	errorCore := zapcore.NewCore(encoder, writeSyncerError, zapcore.ErrorLevel)
	// 标准输出
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	core := zapcore.NewTee(debugCore, infoCore, warnCore, errorCore, std)
	LOG = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(LOG)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
