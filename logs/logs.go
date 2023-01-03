package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logs *zap.Logger

func init() {
	// var err error
	config := zap.NewProductionEncoderConfig()

	//  SugaredLogger
	// config := zap.NewProductionConfig()
	// config.EncoderConfig.TimeKey = "timestamp"
	// config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// config.EncoderConfig.StacktraceKey = ""

	// log, err = config.Build(zap.AddCallerSkip(1))
	// if err != nil {
	// 	panic(err)
	// }

	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	fileLogName := "logs/logfile.log"
	logFile, _ := os.OpenFile(fileLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	logs = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

}

func Info(message string, fields ...zap.Field) {
	logs.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logs.Info(message, fields...)
}

func Error(message interface{}, fields ...zap.Field) {
	// message ที่ส่งมาเป็น type error หรือไม่
	// msg, ok := message.(error)

	// ถอด type จาก interface ด้วย if/switch

	switch msg := message.(type) {
	case error:
		logs.Error(msg.Error())
	case string:
		logs.Error(msg, fields...)
	}
}
