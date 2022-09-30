// Copyright 2022 Giuseppe De Palma, Matteo Trentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/theckman/yacspin"
)

func TestBuilder(t *testing.T) {

	t.Run("NewLoggerBuilder returns new builder", func(t *testing.T) {
		builder := NewLoggerBuilder()
		assert.NotNil(t, builder)
	})

	t.Run("NewLoggerBuilder returns builder with default values", func(t *testing.T) {
		builder := NewLoggerBuilder()
		assert.False(t, builder.(*loggerBuilder).debug)
		assert.True(t, builder.(*loggerBuilder).spinCfg.SuffixAutoColon)
		assert.Equal(t, "✓", builder.(*loggerBuilder).spinCfg.StopCharacter)
		assert.Equal(t, []string{"fgGreen"}, builder.(*loggerBuilder).spinCfg.StopColors)
		assert.Equal(t, "✗", builder.(*loggerBuilder).spinCfg.StopFailCharacter)
		assert.Equal(t, []string{"fgRed"}, builder.(*loggerBuilder).spinCfg.StopFailColors)
		assert.Equal(t, os.Stdout, builder.(*loggerBuilder).writer)
		assert.Equal(t, false, builder.(*loggerBuilder).disableAnimation)
	})

	t.Run("Build returns new logger", func(t *testing.T) {
		logger, err := NewLoggerBuilder().Build()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		assert.False(t, logger.(*FLoggerImpl).debug)
	})

	t.Run("WithDebug true activates the debug in logger", func(t *testing.T) {
		logger, err := NewLoggerBuilder().WithDebug(true).Build()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		assert.True(t, logger.(*FLoggerImpl).debug)
	})

	t.Run("WithWriter sets the writer in logger", func(t *testing.T) {
		var outbuf bytes.Buffer
		logger, err := NewLoggerBuilder().WithWriter(&outbuf).Build()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		assert.Equal(t, &outbuf, logger.(*FLoggerImpl).writer)
	})
	t.Run("Not using WithWriter keeps default os.Stdout", func(t *testing.T) {
		logger, err := NewLoggerBuilder().Build()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		assert.Equal(t, os.Stdout, logger.(*FLoggerImpl).writer)
	})

	t.Run("SpinnerFrequency generates error if input is less or equal to 0", func(t *testing.T) {
		logger, err := NewLoggerBuilder().SpinnerFrequency(0 * time.Millisecond).Build()
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("SpinnerFrequency appends err if err in builder already present", func(t *testing.T) {
		builder := NewLoggerBuilder()
		builder.(*loggerBuilder).err = errors.New("test err")
		_, err := builder.SpinnerFrequency(0 * time.Millisecond).Build()
		assert.Error(t, err)
		assert.Equal(t, "spinner frequency must be greater than 0, test err", err.Error())
	})

	t.Run("SpinnerCharSet fails if input is out of range [0,90]", func(t *testing.T) {
		logger, err := NewLoggerBuilder().SpinnerCharSet(91).Build()
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("DisableAnimation disables animation in logger", func(t *testing.T) {
		logger, err := NewLoggerBuilder().DisableAnimation().Build()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		assert.True(t, logger.(*FLoggerImpl).disableAnimation)
	})
}

func setupWriterLogger(disableAnimation bool) (FLogger, *bytes.Buffer) {
	var outbuf bytes.Buffer
	logger, _ := NewLoggerBuilder().WithWriter(&outbuf).Build()
	logger.(*FLoggerImpl).disableAnimation = disableAnimation
	return logger, &outbuf
}
func TestSpinner(t *testing.T) {

	t.Run("StartSpinner starts spinner and saves message", func(t *testing.T) {
		logger, _ := setupWriterLogger(false)

		_ = logger.StartSpinner("test")

		assert.Equal(t, "test", logger.(*FLoggerImpl).currentMessage)
		assert.Equal(t, yacspin.SpinnerRunning, logger.(*FLoggerImpl).spinner.Status())
	})

	t.Run("StartSpinner called when spinner is already running returns error", func(t *testing.T) {
		logger, _ := setupWriterLogger(false)

		_ = logger.StartSpinner("test")
		err := logger.StartSpinner("test")

		assert.Error(t, err)
	})

	t.Run("StopSpinner stops spinner with success and returns nil when given nil", func(t *testing.T) {
		logger, _ := setupWriterLogger(false)

		_ = logger.StartSpinner("test")
		err := logger.StopSpinner(nil)

		assert.NoError(t, err)
		assert.Equal(t, yacspin.SpinnerStopped, logger.(*FLoggerImpl).spinner.Status())
	})

	t.Run("StopSpinner stops spinner and returns error when given error", func(t *testing.T) {
		logger, _ := setupWriterLogger(false)

		_ = logger.StartSpinner("test")

		inputErr := errors.New("test err")
		err := logger.StopSpinner(inputErr)

		assert.EqualError(t, err, "test err")
		assert.Equal(t, yacspin.SpinnerStopped, logger.(*FLoggerImpl).spinner.Status())
	})

	t.Run("StopSpinner called when spinner is not running returns error", func(t *testing.T) {
		logger, _ := setupWriterLogger(false)

		err := logger.StopSpinner(nil)

		assert.Error(t, err)
	})

	t.Run("SpinnerMessage sets currentMessage", func(t *testing.T) {
		logger, _ := setupWriterLogger(false)

		logger.SpinnerMessage("test")

		assert.Equal(t, "test", logger.(*FLoggerImpl).currentMessage)
	})

	t.Run("StartSpinner with disableAnimation does not start the spinner", func(t *testing.T) {
		logger, _ := setupWriterLogger(true)
		_ = logger.StartSpinner("test")
		assert.Equal(t, yacspin.SpinnerStopped, logger.(*FLoggerImpl).spinner.Status())
	})

	t.Run("StartSpinner with disableAnimation does a simple print", func(t *testing.T) {
		logger, outbuf := setupWriterLogger(true)
		_ = logger.StartSpinner("test")
		assert.Equal(t, "test\n", outbuf.String())
	})

	t.Run("StopSpinner with disableAnimation prints done if no error occured and failed if err", func(t *testing.T) {
		logger, outbuf := setupWriterLogger(true)

		_ = logger.StopSpinner(nil)
		assert.Equal(t, "done\n", outbuf.String())
		outbuf.Reset()

		_ = logger.StopSpinner(errors.New("test err"))
		assert.Equal(t, "failed\n", outbuf.String())
	})

}

func ExampleFLoggerImpl_Info() {
	logger, _ := NewLoggerBuilder().Build()
	logger.Info("test info log")
	// Output: test info log
}

func ExampleFLoggerImpl_Infof() {
	logger, _ := NewLoggerBuilder().Build()
	logger.Infof("test info log %s", "with format")
	// Output: test info log with format
}

func ExampleFLoggerImpl_Debug() {
	logger, _ := NewLoggerBuilder().Build()
	logger.Debug("test debug log with debug disabled")
	// Output:
}

func ExampleFLoggerImpl_Debugf() {
	logger, _ := NewLoggerBuilder().Build()
	logger.Debugf("test debug log %s with debug disabled", "with format")
	// Output:
}

func ExampleFLoggerImpl_Debug_withDebugEnabled() {
	logger, _ := NewLoggerBuilder().WithDebug(true).Build()
	logger.Debug("test debug log with debug enabled")
	// Output: DEBUG: test debug log with debug enabled
}

func ExampleFLoggerImpl_Debugf_withDebugEnabled() {
	logger, _ := NewLoggerBuilder().WithDebug(true).Build()
	logger.Debugf("test debug log %s with debug enabled", "with format")
	// Output: DEBUG: test debug log with format with debug enabled
}
