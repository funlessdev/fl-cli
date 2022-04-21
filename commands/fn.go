package commands

import (
	"fmt"

	"github.com/funlessdev/funless-cli/client"
)

type fn struct {
	Name string `arg:"" help:"name of the function to invoke"`
}

func (f *fn) Run(invoker client.Invoker) error {
	res, err := invoker.Invoke(f.Name)
	if err != nil {
		return err
	}

	fmt.Println(res.Status)
	return nil
}
