package tradelog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogrus() {
	logfile, _ := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	Logger.SetReportCaller(false) // 设置日志是否记录被调用的位置
	Logger.Out = logfile
	Logger.Formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"}
	Logger.Level = logrus.InfoLevel
}
