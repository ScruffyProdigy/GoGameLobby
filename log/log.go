package log

import "os"
import "io"
import "fmt"

const (
	Log_Level_Debug = iota
	Log_Level_Info
	Log_Level_Warning
	Log_Level_Error
	Log_Level_Fatal
	Log_Level_Unknown
)

var loggerLevel = Log_Level_Debug

type DummyLogger struct {
}

func (DummyLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var dummy DummyLogger

func SetLogLevel(level int) {
	loggerLevel = level
}

func log(message string, level int) {
	if level >= loggerLevel {
		fmt.Print(message)
	}
}

func logger(level int) io.Writer {
	if level >= loggerLevel {
		return os.Stdout
	}
	return dummy
}

func Debug(message string) {
	log("\nDebug: "+message, Log_Level_Debug)
}

func DebugLog() io.Writer {
	return logger(Log_Level_Debug)
}

func Info(message string) {
	log("\nInfo: "+message, Log_Level_Info)
}

func InfoLog() io.Writer {
	return logger(Log_Level_Info)
}

func Warning(message string) {
	log("\nWarning: "+message, Log_Level_Warning)
}

func WarningLog() io.Writer {
	return logger(Log_Level_Warning)
}

func Error(message string) {
	log("\nError: "+message, Log_Level_Error)
}

func ErrorLog() io.Writer {
	return logger(Log_Level_Error)
}

func Fatal(message string) {
	log("\nFatal: "+message, Log_Level_Fatal)
}

func FatalLog() io.Writer {
	return logger(Log_Level_Fatal)
}

func Unknown(message string) {
	log("\nUnknown: "+message, Log_Level_Unknown)
}

func UnknownLog() io.Writer {
	return logger(Log_Level_Unknown)
}
