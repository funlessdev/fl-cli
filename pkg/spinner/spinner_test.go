package spinner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theckman/yacspin"
)

func TestCreateSpinner(t *testing.T) {
	cfg := yacspin.Config{}
	spinner, err := CreateSpinner(cfg)
	assert.NoError(t, err)
	assert.Equal(t, spinner.Status(), yacspin.SpinnerStopped)
}
