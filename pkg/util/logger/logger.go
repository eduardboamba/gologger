package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

const (
	OutputConsole        = 1
	OutputFile           = 2
	OutputConsoleAndFile = 3
	OutputBuffer         = 4

	LevelFatal = 0
	LevelError = 1
	LevelInfo  = 2
	LevelDebug = 3
)

var (
	logError *log.Logger
	logInfo  *log.Logger
	logDebug *log.Logger
)

type LogConfiguration struct {
	flags  int
	output int
	file   *os.File
	level  int
	buff   *bytes.Buffer
}

var logConfig LogConfiguration

func init() {
	logConfig = getDefaultConfiguration()

	logError = log.New(getWriter(), "[ERROR]", logConfig.flags)
	logInfo = log.New(getWriter(), "[INFO]", logConfig.flags)
	logDebug = log.New(getWriter(), "[DEBUG]", logConfig.flags)
}

func getDefaultConfiguration() LogConfiguration {
	return LogConfiguration{
		flags:  log.Ldate | log.Ltime,
		output: OutputConsole,
		level:  LevelDebug,
	}
}

func getWriter() io.Writer {
	switch logConfig.output {
	case OutputConsole:
		return os.Stdout

	case OutputFile:
		return logConfig.file

	case OutputConsoleAndFile:
		return io.MultiWriter(logConfig.file, os.Stdout)

	case OutputBuffer:
		return logConfig.buff
	}

	/* No supported match, so writing to null device;
	The writing is successful, even though nothing is done
	*/
	return ioutil.Discard
}

func updateLoggerConfig() {
	// Update flags
	logError.SetFlags(logConfig.flags)
	logInfo.SetFlags(logConfig.flags)
	logDebug.SetFlags(logConfig.flags)

	// Update output
	logError.SetOutput(getWriter())
	logInfo.SetOutput(getWriter())
	logDebug.SetOutput(getWriter())
}

// SetFlags accepts an OR combination of flags defined in the standard log package
// setting the flags to 0 will remove any formatting
func SetFlags(flags int) {
	logConfig.flags = flags

	updateLoggerConfig()
}

func SetOutputConsole() {
	logConfig.output = OutputConsole

	updateLoggerConfig()
}

// SetOutputFile redirects all output to `fileName`
// if OpenFile fails, function returns the error
// the file is created if it doesn't exist, opened in write-only mode and all writes are appended
func SetOutputFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logConfig.file = file
	logConfig.output = OutputFile

	updateLoggerConfig()

	return nil
}

// SetOutputConsoleAndFile behaves similar to SetOutputFile
// same output is redirected also to the console
func SetOutputConsoleAndFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logConfig.file = file
	logConfig.output = OutputConsoleAndFile

	updateLoggerConfig()

	return nil
}

func SetOutputBuffer(bytes *bytes.Buffer) {
	logConfig.output = OutputBuffer
	logConfig.buff = bytes

	updateLoggerConfig()
}

func SetLogLevel(level int) error {
	if level != LevelFatal && level != LevelError && level != LevelInfo && level != LevelDebug {
		return errors.Errorf("Unsupported log level %d", level)
	}

	logConfig.level = level
	return nil
}

func Fatal(v ...interface{}) {
	pc, fn, line, _ := runtime.Caller(1)
	logError.Fatalf("at %s[%s:%d] %v", shortenName(runtime.FuncForPC(pc).Name(), "."), shortenName(fn, "/"), line, v)
}

func Error(v ...interface{}) {
	if logConfig.level >= LevelError {
		pc, fn, line, _ := runtime.Caller(1)
		prefix := fmt.Sprintf("at %s[%s:%d]", shortenName(runtime.FuncForPC(pc).Name(), "."), shortenName(fn, "/"), line)

		logError.Println(prefix, v)
	}
}

func Info(v ...interface{}) {
	if logConfig.level >= LevelInfo {
		logInfo.Println(v...)
	}
}

func Debug(v ...interface{}) {
	if logConfig.level >= LevelDebug {
		logDebug.Println(v...)
	}
}

func shortenName(name, token string) string {

	if index := strings.LastIndex(name, token); index != -1 {
		a := []rune(name)
		return string(a[index+1 : len(name)])
	}

	return name
}