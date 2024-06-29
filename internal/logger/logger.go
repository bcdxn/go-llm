package logger

import (
	"fmt"
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

func New(dir string, level hclog.Level) (*Logger, error) {
	var l *Logger

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
		Name:   "llm",
		Output: stdoutFile,
		Level:  level,
	})

	stderr := hclog.New(&hclog.LoggerOptions{
		Name:   "llm",
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

func (l Logger) Named(name string) *Logger {
	return &Logger{
		level:      l.level,
		stderrFile: l.stderrFile,
		stdoutFile: l.stdoutFile,
		stderr:     l.stderr.Named(name),
		stdout:     l.stdout.Named(name),
	}
}

func (l Logger) GetLevel() hclog.Level {
	return l.stdout.GetLevel()
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

func (l Logger) Close() {
	l.stderrFile.Close()
	l.stdoutFile.Close()
}

func SimpleLogFatal(msg string, args ...interface{}) {
	l := append([]interface{}{msg}, args)

	log.Fatal(l...)
}

type CtxLogger struct{}
