package log

import "time"

type (
	FLogger interface {
		SpinnerSuffix(string)
		SpinnerMessage(string)
		StartSpinner(string)
		StopSpinner(error) error

		Info(args ...interface{})
		Infof(format string, args ...interface{})

		Debug(args ...interface{})
		Debugf(format string, args ...interface{})
	}

	builder interface {
		WithDebug(bool) builder
		SpinnerFrequency(time.Duration) builder
		SpinnerCharSet(int) builder
		Build() (FLogger, error)
	}
)
