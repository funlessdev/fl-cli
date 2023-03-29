package cfg

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/homedir"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Cfg struct {
	Set CfgSet `cmd:"" name:"set" aliases:"s" help:"set a property in the config file"`
	Get CfgGet `cmd:"" name:"get" aliases:"g" help:"get a property from the config file"`
}

type CfgSet struct {
	Key   string `arg:"" enum:"${config_keys}" help:"name of the parameter that is being set"`
	Value string `arg:"" help:"value of the parameter that is being set"`
}

type CfgGet struct {
	Key string `arg:"" enum:"${config_keys}" help:"name of the parameter that is being read"`
}

func (g *CfgSet) Run(ctx context.Context, logger log.FLogger, config client.Config) error {

	configBasePath := path.Base(config.Path)
	configText, _, err := homedir.ReadFromConfigDir(configBasePath)
	configString := string(configText[:])
	if err != nil {
		return nil
	}

	r, _ := regexp.Compile(fmt.Sprintf("%s=(.+)", g.Key))

	var outConfig string

	if r.MatchString(configString) {
		outConfig = r.ReplaceAllLiteralString(configString, fmt.Sprintf("%s=%s", g.Key, g.Value))
	} else {
		outConfig = fmt.Sprintf("%s\n%s=%s\n", strings.Trim(configString, "\n"), g.Key, g.Value)
	}

	_, err = homedir.WriteToConfigDir(configBasePath, []byte(outConfig), true)
	if err != nil {
		return err
	}

	logger.Infof("Key %s set to %s in config (path %s).\n", g.Key, g.Value, config.Path)

	return nil
}

func (g *CfgGet) Run(ctx context.Context, logger log.FLogger, config client.Config) error {
	var cfgValue string
	switch g.Key {
	case "host":
		cfgValue = config.Host
	case "api_token":
		cfgValue = config.APIToken
	case "admin_token":
		cfgValue = config.AdminToken
	case "secret_key_base":
		cfgValue = config.SecretKeyBase
	}

	logger.Infof("%s=%s\n", g.Key, cfgValue)

	return nil
}
