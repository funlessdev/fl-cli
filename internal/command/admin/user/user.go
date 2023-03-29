// Copyright 2023 Giuseppe De Palma, Matteo Trentin
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

package admin_user

import (
	"context"
	"fmt"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type User struct {
	Create CreateUser `cmd:"" name:"create" aliases:"c" help:"Create a new FunLess user"`
	List   ListUsers  `cmd:"" name:"list" aliases:"l" help:"List all FunLess users"`
}

type CreateUser struct {
	Name string `arg:"" name:"name" help:"The name of the new user"`
}

func (u *CreateUser) Run(ctx context.Context, userHandler client.UserHandler, logger log.FLogger) error {
	logger.StartSpinner("Creating user...")
	res, err := userHandler.Create(ctx, u.Name)
	_ = logger.StopSpinner(err)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("User %s created. Auth token:\n", res.Name))
	logger.Info(res.Token)
	return err
}

func (u *User) Help() string {
	return `
DESCRIPTION
	
	Create a new FunLess user. A user is a (name, token) pair, used to 
	authenticate to the FunLess API.
	To create new user specify an unique name. The token will be generated 
	automatically by the FunLess Platform.

EXAMPLES

	$ fl admin user create userA\n
`
}

type ListUsers struct {
}

func (u *ListUsers) Run(ctx context.Context, userHandler client.UserHandler, logger log.FLogger) error {
	logger.StartSpinner("Listing existing users...")
	res, err := userHandler.List(ctx)
	_ = logger.StopSpinner(err)
	if err != nil {
		return err
	}
	logger.Info("Users:")
	for _, user := range res.Names {
		logger.Info(fmt.Sprintf("- %s", user))
	}
	return err
}

func (u *ListUsers) Help() string {
	return `
DESCRIPTION

	List all existing FunLess user names. To create a new user use 
	the "create" command.
	To get the token of an existing user use the "token" command (TODO). 

EXAMPLES

	$ fl admin user list
`
}
