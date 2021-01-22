package common

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"gopkg.in/natefinch/lumberjack.v2"
)

func LogInit(logPath string, maxDayAge int, maxSize int) *zap.SugaredLogger {
	var logger *zap.Logger
	if logPath == "" {
		logger, _ = zap.NewDevelopment()

		defer logger.Sync() // flushes buffer, if any
	} else {
		// lumberjack.Logger is already safe for concurrent use, so we don't need to
		// lock it.
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath, //"/var/log/myapp/foo.log",
			MaxSize:    maxSize, // megabytes
			MaxBackups: 1,
			MaxAge:     maxDayAge, // days
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			//zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		)
		logger = zap.New(core)
	}

	return logger.Sugar()
}

