package log

import (
	"log"
	"os"
	"time"

	"github.com/theckman/yacspin"
)

type FLogger interface {
	SpinnerSuffix(suffix string)
	SpinnerMessage(msg string)
	StartSpinner(msg string)
	StopSpinner(success bool)

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type BaseFLogger struct {
	debug   bool
	logger  *log.Logger
	spinner *yacspin.Spinner
}

func NewBaseLogger(debug bool) (*BaseFLogger, error) {
	cfg := yacspin.Config{
		Frequency:         150 * time.Millisecond,
		Colors:            []string{"fgYellow"},
		CharSet:           yacspin.CharSets[59],
		SuffixAutoColon:   true,
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopMessage:       "done",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		return nil, err
	}
	return &BaseFLogger{debug: debug, logger: log.New(os.Stdout, "", 5), spinner: s}, nil
}

func (l *BaseFLogger) SpinnerSuffix(suffix string) {
	l.spinner.Suffix(suffix)
}

func (l *BaseFLogger) SpinnerMessage(msg string) {
	l.spinner.Message(msg)
}

func (l *BaseFLogger) StartSpinner(msg string) {
	l.spinner.Message(msg)
	l.spinner.Start()
}

func (l *BaseFLogger) StopSpinner(success bool) {
	if success {
		l.spinner.Stop()
	} else {
		l.spinner.StopFail()
	}
}

func (l *BaseFLogger) Info(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *BaseFLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *BaseFLogger) Debug(args ...interface{}) {
	if l.debug {
		l.logger.Println(args...)
	}
}

func (l *BaseFLogger) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.logger.Printf(format, args...)
	}
}
