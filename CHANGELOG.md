<!--
  ~ Copyright 2022 Giuseppe De Palma, Matteo Trentin
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~ http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

## v0.3.3 (2023-02-02)

### Fix

- **docker-deploy**: download filebeat and .env files for new docker compose

## v0.3.2 (2023-01-28)

### Fix

- code file pointer with updated deps

### Refactor

- **logs**: uniform Info and Infof removing \n from Info

## v0.3.1 (2023-01-24)

### Fix

- **command**: raise error in json decoder if unknown fields are passed in apierror
- **mod**: add missing extracterror calls
- **wasm**: set unique name for builder containers

## v0.3.0 (2023-01-16)

### Feat

- **docker**: check containerwait exit status code
- **fn**: add check for necessary files in build and create
- **app**: add input validator to fnservice and modservice
- **client**: add name validation in fn_service and mod_service
- **client**: add input_validator struct and interface
- **mod**: fix mod get output formatting
- **mod**: add mod subcommands
- add modhandler interface over sdk
- add base files for mod subcommand
- **homedir**: add homedir pkg with utilities to handle .fl
- **fn-new**: implement new cmd with tests
- **template**: add template pull cmd
- **fn**: setup 'new' cmd with list flag
- **deploy**: add alias for deploy command
- **admin**: add deploy subcommands for k8s/docker
- **kubernetes_rm**: add kubernetes remove command
- **kubernetes**: add full kubernetes deployment
- **kubernetes**: add namespace, svc-account, role and rolebinding creation
- **fl_k8s**: add kubernetes yaml parsing
- add base k8s deploy command
- **deploy**: add base kubernetes deployer/remover types
- **upload**: add upload cmd to create functions from wasm files
- **wasm-builder**: add rename code.wasm
- **fn-build**: add build cmd to just build wasm binaries

### Fix

- update mod, fn_service, mod_service to new sdk types
- **mod**: fix print in mod list and mod get
- **mod_service**: fix nil pointer deref errors
- **fn_service**: fix nil pointer deref error
- **fn**: remove fn list command (moved to mod get)
- **fn**: update fn commands to new fnhandler
- **fn_service**: update fnhandler to new sdk version
- license header and linter errors
- check returned errors
- **license**: add license header
- **constants**: fix broken docker urls
- **create**: add get wasm file in builder to get correct file in create
- license headers and language flag type
- **fn**: add mandatory .wasm extension in fn create with single file
- **fn**: return error if empty file is passed in fn create

### Refactor

- move supported languages names/images in single struct
- **mod**: fix output formatting for mod subcommands
- **mod_service_test**: remove redundant type declaration
- add components replace in docker-compose
- add prometheus config download
- homdir in constants and create dir function
- **docker**: update docker deploy commands with new interface
- **docker**: substitute docker deployer/remover with dockershell interface
- **homedir**: return path to file from read/write
- **homedir**: move dir string to constant
- **template**: small improvements
- **template**: move new --list to template list cmd
- **template**: move copy code to pkg
- **fn**: add command aliases
- **deploy**: move deploy subcommands to separate folders
- **deploy**: move docker and k8s commands under deploy subcommand
- **deploy**: fix typo in kubernetes_remover struct
- **deploy**: move fl_k8s/parse to deploy/kubernetes_parse
- **constants**: revert default prometheus img
- **kubernetes**: remove unused parameters from k8s command
- change WithClientSet to WithConfig in k8s deployer/remover
- **kubernetes_deployer**: add getyamlcontent() auxiliary func; update interface
- **kubernetes**: add setupdeployer() function
- **fn**: rename mockInvoker to mockFnHandler
- **list**: remove redundant error check
- **create**: update create to build and upload with tmp file
- **create**: setup builder
- simplify docker client
- **wasm**: update builder
- implement new docker remover
- **admin-dev**: use new deployer and update tests
- **docker**: simplify docker client interface
- **docker-deployer**: separate setup
- **docker**: add docker package with interfaces
- **constants**: change builder map names
- **fn_service**: remove unused language parameter in create
- **fn**: add missing license headers
- **fn**: move fn functions to separate files
- remove unused Namespace field

## v0.2.1 (2022-10-31)

### Fix

- **dev**: pulling prometheus log

## v0.2.0 (2022-10-30)

### Feat

- **dev**: add prometheus in admin dev deployment
- **fn**: add out-dir param in fn create
- **fn**: update "fn create" command with build and source-dir
- **build**: add build module for function creation
- update fn and fn_service to new sdk version
- **dev**: add secret_key_base env in dev deployment
- add custom worker and core image in dev deployment
- **docker_utils**: add existence check in image pull

### Fix

- change host ip to local
- use correct prom image and expose port
- **fn**: check cast
- openapi error handling
- **build**: fix broken out-path in wasm builder
- **fn**: add default value for namespace in create/delete/invoke
- **fn**: fix out-dir and source-file interaction
- add missing license header in test file
- **fn**: mark language as not required in create
- **dev**: forward correct port in fl-core
- fix host and port configuration
- **license**: change with correct license header

### Refactor

- move names in constants and remove prom in reset
- move interfaces and impl in one file
- change parameter name
- **wasm**: use docker utils in wasm builder
- **dev**: drop ow and update local deployer
- **docker_utils**: move docker utils in its own package
- **main**: separate parse and run
- **main**: move app to its own package
- **main**: separate cli run in app.go
- **deploy**: add name for ctx parameter in dockerdeployer
- **build**: change wasm build functions to struct methods
- add type in Source flag and change resources to fixtures
- **fn**: migrate fn and fn_service to new api sdk version
- move images to Setup and pass kong default vars
- **docker_utils**: fix linting errors

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
