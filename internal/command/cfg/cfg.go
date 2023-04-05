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

package cfg

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/funlessdev/fl-cli/pkg"
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

	var configBasePath string

	if config.Path == "" {
		configBasePath = pkg.ConfigFileName
	} else {
		configBasePath = path.Base(config.Path)
	}

	configText, _, err := homedir.ReadFromConfigDir(configBasePath)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	configString := string(configText[:])
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
	case "api_host":
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
