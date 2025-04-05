package logger

import (
	"fmt"
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

var auditLogger *log.Logger

func init() {
	logFile := &lumberjack.Logger{
		Filename:   "logs/audit.log",
		MaxSize:    10, // Max size in MB before rotating
		MaxBackups: 3,  // Max number of old log files to keep
		MaxAge:     30,
		Compress:   true,
	}

	auditLogger = log.New(logFile, "AUDIT: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogAudit(userID *uint, action string, details string) {
	id := "unauthenticated"
	if userID != nil {
		id = fmt.Sprintf("UserID=%d", *userID)
	}
	msg := fmt.Sprintf("[%s] %s - %s", id, action, details)
	auditLogger.Println(msg)
}
