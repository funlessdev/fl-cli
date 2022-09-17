<!--
  ~ Licensed to the Apache Software Foundation (ASF) under one
  ~ or more contributor license agreements.  See the NOTICE file
  ~ distributed with this work for additional information
  ~ regarding copyright ownership.  The ASF licenses this file
  ~ to you under the Apache License, Version 2.0 (the
  ~ "License"); you may not use this file except in compliance
  ~ with the License.  You may obtain a copy of the License at
  ~
  ~   http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing,
  ~ software distributed under the License is distributed on an
  ~ "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  ~ KIND, either express or implied.  See the License for the
  ~ specific language governing permissions and limitations
  ~ under the License.
-->

## v0.1.0 (2022-09-17)

### Feat

- **dev**: bind logs to host folder
- **log**: add DisableAnimation in log builder
- **admin-dev**: add internal network for worker+runtimes in dev deployment
- add 'a' alias for admin cmd
- **log**: add parametric writer for logger
- add invoke, create and delete fn subcommands
- **admin-deploy**: reset sub command to remove local deployment
- **admin-deploy**: give proper names to fl containers
- **admin-deploy**: deploy containers attached on fl_net network
- **admin**: start core and worker container with admin deploy
- add version flag (-v or --version)
- **main**: add version flag
- wire logger into cli tool and use preflights in admin deploy
- **log**: introduce log pkg that wraps spinner
- add docker pkg to run preflight tests
- add sample admin command with spinner
- create client and fn service in main and bind it to kong
- **fn**: use FnService.Invoke when cmd fn is used
- **FnService**: add FnService with simple Invoke
- **client**: add send method and remove interface
- **client**: add an initial client package
- setup kong library with sample cli main
- create go project

### Fix

- **deploy-local**: pass correct network to worker RUNTIME_NETWORK env var
- **reset**: update reset with new deployer
- **license**: add license header
- **log**: don't clear currentMessage in StopSpinner
- **license**: add license header
- task install file and error warning
- **fn**: extract error content from swagger errors
- **deploy_local**: remove protocol from docker_host value in deployment
- **admin-deploy**: mount docker socket in container
- **license**: add license header
- **log**: return correct err in StopSpinner
- **log**: add license header
- **log**: add err handling for spinner start and stop
- **license**: add license headers
- **admin**: fix logging and add image pulling
- **Taskfile.yml**: fix paths
- **spinner**: add missing license header

### Refactor

- merge remove networks to reflect the setup
- **reset**: add reset functions to deployer interface
- **dev**: deployer arch and tests output
- **logs**: add testing mode for log and missing tests
- **admin**: remove init token and join
- **admin**: break admin subcmds in files
- **command**: move admin and fn in their own packages
- move private cli code in internal folder
- **go.mod**: change module name from funless-cli to fl-cli
- **fn**: remove writer from fn functions, add logger as parameter
- add writer and context as parameters for fn
- **client**: remove unused functions and tests in client
- **deploy_local**: move docker host code in utils function
- re-organize deploy code in admin pkg
- **admin-deploy**: remove preflight checks and start refact admin cmd
- **constants**: use latest tag for docker images
- **log**: builder pattern with tests
- **logs**: change log message style with a spinner each step
- **docker-package**: move shell interface in its file
- **log.go**: add DEBUG prefix for debug messages
- change log writer to fmt
- **preflight**: use the logger to print preflight status
- project structure refactor
- **license**: add license header to new files
- **client**: refactor client.go and reduce clientAPI to just get request
