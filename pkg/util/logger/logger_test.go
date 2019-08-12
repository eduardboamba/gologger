package logger_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/eduardboamba/gologger/pkg/util/logger"
	"gotest.tools/assert"
)

const staticMsg = "test message"

const (
	Rdate      = `\d{4}\/\d{2}\/\d{2}`
	Rtime      = `\d{2}:\d{2}:\d{2}`
	Rline      = `\d+`
	Rshortfile = `[A-Za-z0-9_\-]+\.go:` + Rline
	Rfunc      = `at [A-Za-z0-9_\-]+`

	RError = `\[ERROR\]`
	RInfo  = `\[INFO\]`
	RDebug = `\[DEBUG\]`
)

var tests = []struct {
	name    string
	flags   int
	level   int
	pattern string
}{
	{
		name:    "log error",
		flags:   log.Ldate | log.Ltime,
		level:   logger.LevelError,
		pattern: RError + Rdate + " " + Rtime + " " + Rfunc + `\[` + Rshortfile + `\]` + " " + `\[` + staticMsg + `\]`,
	},
	{
		name:    "log info",
		flags:   log.Ldate | log.Ltime,
		level:   logger.LevelInfo,
		pattern: RInfo + Rdate + " " + Rtime + " " + staticMsg,
	},
	{
		name:    "log debug",
		flags:   log.Ldate | log.Ltime,
		level:   logger.LevelDebug,
		pattern: RDebug + Rdate + " " + Rtime + " " + staticMsg,
	},
}

func TestLogStatements(t *testing.T) {
	for _, test := range tests {
		t.Log("Test: " + test.name)

		// Log to buffer, to compare output
		buff := new(bytes.Buffer)
		logger.SetOutputBuffer(buff)

		log.SetFlags(test.flags)
		if err := logger.SetLogLevel(test.level); err != nil {
			t.Errorf("Setting log level failed with %v", err)
		}

		switch test.level {
		case logger.LevelError:
			logger.Error(staticMsg)

		case logger.LevelInfo:
			logger.Info(staticMsg, 4)

		case logger.LevelDebug:
			logger.Debug(staticMsg, 3)
		}

		line := buff.String()
		line = line[0 : len(line)-1]
		pattern := "^" + test.pattern
		matched, err4 := regexp.MatchString(pattern, line)
		if err4 != nil {
			t.Fatal("pattern did not compile:", err4)
		}
		if !matched {
			t.Errorf("log output should match %q is %q", pattern, line)
		}
	}
}

func TestLogErrors(t *testing.T) {
	var err = logger.SetLogLevel(100)
	assert.Error(t, err, "Unsupported log level 100")

	err = logger.SetLogLevel(logger.LevelError)
	assert.NilError(t, err)

	// Create a temporary file
	file, err := ioutil.TempFile("", "output-*.out")
	if err != nil {
		t.Fatal("Cannot create temporary file:", err)
	}
	defer os.Remove(file.Name())

	err = logger.SetOutputConsoleAndFile("")
	assert.ErrorType(t, err, &os.PathError{})

	err = logger.SetOutputConsoleAndFile(file.Name())
	assert.NilError(t, err)

	err = logger.SetOutputFile("")
	assert.ErrorType(t, err, &os.PathError{})

	err = logger.SetOutputFile(file.Name())
	assert.NilError(t, err)
}