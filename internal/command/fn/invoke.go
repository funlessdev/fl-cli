package fn

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Invoke struct {
	Name      string            `arg:"" name:"name" help:"name of the function to invoke"`
	Namespace string            `name:"namespace" short:"n" default:"_" help:"namespace of the function to invoke"`
	Args      map[string]string `name:"args" short:"a" help:"arguments of the function to invoke" xor:"args"`
	JsonArgs  string            `name:"json" short:"j" help:"json encoded arguments of the function to invoke; overrides args" xor:"args"`
}

func (f *Invoke) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger) error {
	args := make(map[string]interface{}, len(f.Args))
	if f.Args != nil {
		for k, v := range f.Args {
			args[k] = v
		}
	} else if f.JsonArgs != "" {
		err := json.Unmarshal([]byte(f.JsonArgs), &args)
		if err != nil {
			return err
		}
	}
	res, err := fnHandler.Invoke(ctx, f.Name, f.Namespace, args)
	if err != nil {
		return extractError(err)
	}

	if res.Result != nil {
		decodedRes, err := json.Marshal(res.Result)
		if err != nil {
			return err
		}
		logger.Info(string(decodedRes))
	} else {
		return errors.New("received nil result")
	}

	return nil
}
