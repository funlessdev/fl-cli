package log

import (
	"fmt"
	"time"

	"github.com/theckman/yacspin"
)

type loggerBuilder struct {
	debug   bool
	spinCfg yacspin.Config
	err     error
}

func NewLoggerBuilder() builder {
	return &loggerBuilder{
		spinCfg: yacspin.Config{
			SuffixAutoColon:   true,
			StopCharacter:     "✓",
			StopColors:        []string{"fgGreen"},
			StopFailCharacter: "✗",
			StopFailColors:    []string{"fgRed"},
		},
	}
}

func (l *loggerBuilder) WithDebug(b bool) builder {
	l.debug = b
	return l
}

func (l *loggerBuilder) SpinnerFrequency(freq time.Duration) builder {

	if freq <= time.Duration(0) {
		l.err = fmt.Errorf("spinner frequency must be greater than 0, %w", l.err)
		return l
	}

	if l.err != nil {
		return l
	}

	l.spinCfg.Frequency = freq
	return l
}

// // SpinnerCharSet sets the character set to use for the spinner animation [0 to 90].
func (l *loggerBuilder) SpinnerCharSet(charset int) builder {
	if charset < 0 || charset > 90 {
		l.err = fmt.Errorf("spinner character set must be between 0 and 90, %w", l.err)
		return l
	}
	if l.err != nil {
		return l
	}
	l.spinCfg.CharSet = yacspin.CharSets[charset]
	return l
}

// // Build returns a new logger instance.
func (l *loggerBuilder) Build() (FLogger, error) {

	if l.err != nil {
		return nil, l.err
	}

	s, err := yacspin.New(l.spinCfg)
	if err != nil {
		return nil, err
	}

	logger := &FLoggerImpl{
		debug:          l.debug,
		currentMessage: "",
		spinner:        s,
	}
	return logger, nil
}
