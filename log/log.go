/*
	Log is a simple logging system I'm using to help figure out when & where things are going wrong

	It allows you to specify a Logging Level, and if an attempt to write to the log at a lower log level than specified is attempted, then nothing happens instead
*/
package log

import "os"
import "io"
import "fmt"

/*
	Logging levels are, in order: Debug, Info, Warning, Error, Fatal, Unknown
	each has a corresponding function, which acts as a Print to the default output
	and each also has a Log function, which returns the default output if it is valid, and a fake output if invalid
*/
const (
	Log_Level_Debug = iota
	Log_Level_Info
	Log_Level_Warning
	Log_Level_Error
	Log_Level_Fatal
	Log_Level_Unknown
)

var loggerLevel = Log_Level_Debug

type dummyLogger struct {
}

func (dummyLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var dummy dummyLogger

/*
	SetLogLevel will set the log level.  Any further attempts to write at a lower level will fail
*/
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
