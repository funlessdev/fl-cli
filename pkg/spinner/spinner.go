package spinner

import (
	"github.com/theckman/yacspin"
)

func CreateSpinner(cfg yacspin.Config) (*yacspin.Spinner, error) {
	s, err := yacspin.New(cfg)
	if err != nil {
		// exitf("failed to make spinner from struct: %v", err)
		return nil, err
	}

	return s, nil
}
