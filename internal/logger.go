package llm

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-hclog"
)

type Logger struct {
	level      hclog.Level
	stderrFile *os.File
	stdoutFile *os.File
	stderr     hclog.Logger
	stdout     hclog.Logger
}

func NewLogger(name string, level hclog.Level) (*Logger, error) {
	var l *Logger
	dir, err := GetLogDirPath()
	if err != nil {
		return l, err
	}

	d := time.Now().Format("2006-01-02")

	stdoutFilename := filepath.Join(dir, fmt.Sprintf("./%s-stdout.log", d))
	stderrFilename := filepath.Join(dir, fmt.Sprintf("./%s-stderr.log", d))

	stdoutFile, err := tea.LogToFile(stdoutFilename, "")
	if err != nil {
		return l, err
	}
	stderrFile, err := tea.LogToFile(stderrFilename, "")
	if err != nil {
		return l, err
	}

	stdout := hclog.New(&hclog.LoggerOptions{
		Name:   name,
		Output: stdoutFile,
		Level:  level,
	})

	stderr := hclog.New(&hclog.LoggerOptions{
		Name:   name,
		Output: stderrFile,
		Level:  level,
	})

	return &Logger{
		level,
		stderrFile,
		stdoutFile,
		stderr,
		stdout,
	}, nil

}

func (l Logger) Log(level hclog.Level, msg string, args ...interface{}) {
	l.stdout.Log(level, msg, args)
}

func (l Logger) Trace(msg string, args ...interface{}) {
	l.stdout.Trace(msg, args...)
}

func (l Logger) Debug(msg string, args ...interface{}) {
	l.stdout.Debug(msg, args...)
}

func (l Logger) Info(msg string, args ...interface{}) {
	l.stdout.Info(msg, args...)
}

func (l Logger) Warn(msg string, args ...interface{}) {
	l.stdout.Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...interface{}) {
	l.stderr.Error(msg, args...)
}

func (l Logger) IsTrace() bool {
	return l.stdout.IsTrace()
}

func (l Logger) IsDebug() bool {
	return l.stdout.IsTrace()
}

func (l Logger) IsInfo() bool {
	return l.stdout.IsTrace()
}

func (l Logger) IsWarn() bool {
	return l.stdout.IsTrace()
}

func (l Logger) IsError() bool {
	return l.stdout.IsTrace()
}

func (l Logger) ImpliedArgs() []interface{} {
	return l.stdout.ImpliedArgs()
}

func (l Logger) With(args ...interface{}) hclog.Logger {
	return Logger{
		level:      l.level,
		stderrFile: l.stderrFile,
		stdoutFile: l.stdoutFile,
		stderr:     l.stderr.With(args...),
		stdout:     l.stdout.With(args...),
	}
}

func (l Logger) Name() string {
	return l.stdout.Name()
}

func (l Logger) Named(name string) hclog.Logger {
	return Logger{
		level:      l.level,
		stderrFile: l.stderrFile,
		stdoutFile: l.stdoutFile,
		stderr:     l.stderr.Named(name),
		stdout:     l.stdout.Named(name),
	}
}

func (l Logger) ResetNamed(name string) hclog.Logger {
	return Logger{
		level:      l.level,
		stderrFile: l.stderrFile,
		stdoutFile: l.stdoutFile,
		stderr:     l.stderr.ResetNamed(name),
		stdout:     l.stdout.ResetNamed(name),
	}
}

func (l Logger) SetLevel(level hclog.Level) {
	l.stdout.SetLevel(level)
	l.stderr.SetLevel(level)
}

func (l Logger) GetLevel() hclog.Level {
	return l.stdout.GetLevel()
}

func (l Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return l.stdout.StandardLogger(opts)
}

func (l Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return l.stdout.StandardWriter(opts)
}

func (l Logger) Close() {
	l.stderrFile.Close()
	l.stdoutFile.Close()
}

func MustGetLoggerFromContext(c context.Context, name string) Logger {
	l, ok := c.Value(ctxLogger{}).(Logger)
	if !ok {
		SimpleLogFatal("unable to fetch logger from context")
	}

	if name != "" {
		l = l.Named("modelselect").(Logger)
		if !ok {
			SimpleLogFatal("cannot cast hclog to llm.Logger")
		}
	}

	return l
}

func SetLoggerInContext(c context.Context, l Logger) context.Context {
	return context.WithValue(c, ctxLogger{}, l)
}

// SimpleLogFatal uses the go std log package to log and exit the application -- perfect for when
// instantiation of the llm.Logger itself, fails.
func SimpleLogFatal(msg string, args ...interface{}) {
	l := append([]interface{}{msg}, args)

	log.Fatal(l...)
}

// Get the absolute path to the folder where the log files live
func GetLogDirPath() (string, error) {
	var cfgPath string
	home, err := os.UserHomeDir()
	if err != nil {
		return cfgPath, fmt.Errorf("unable to get user home directory location: %w", err)
	}
	cfgPath = filepath.Join(home, ".llm")
	return cfgPath, nil
}

type ctxLogger struct{}
