<!--
  ~ Copyright 2023 Giuseppe De Palma, Matteo Trentin
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
![tests](https://github.com/funlessdev/fl-cli/actions/workflows/test.yml/badge.svg)
[![release](https://badgen.net/github/release/funlessdev/fl-cli)](https://github.com/funlessdev/fl-cli/releases/)
![contributors](https://badgen.net/github/contributors/funlessdev/fl-cli) 

# The FunLess CLI

This is the repository for the FunLess platform CLI tool (**fl**), a new generation research-driven serverless platform.

The CLI is written in Go, using the [Kong](https://github.com/alecthomas/kong) framework. It can be used to deploy the platform, manage modules and create, delete, invoke functions.

### The Commands

The CLI offers 3 main commands:

- **admin**: currently used only to deploy and remove the platform on a local machine
- **fn**: used to create, delete and invoke functions
- **mod**: used to create, delete and get information on modules
- **template**: used to manage the template folder

Each command has a series of subcommands. Use `--help` to get more information on each command.

#### fl admin

The `admin` command has currently 1 subcommand: `deploy` which can be used chained with:

- `docker up`: to deploy the platform on a local machine using Docker containers
- `docker down`: to remove the platform from a local machine
- `kubernetes up`: to deploy the platform on a Kubernetes cluster (not yet updated)
- `kubernetes down`: to remove the platform from a Kubernetes cluster (not yet updated)

#### fl fn

The `fn` command is used for anything function related. It has currently 6 subcommands: 
  
- `invoke`: to invoke an existing function
- `build`: to build the wasm file from a function's source code
- `upload`: to upload a function's wasm file to the platform with a given name
- `create`: a combination of `build` and `upload`, the wasm file is removed after the upload
- `delete`: to delete a function from the platform
- `new`: to create new function's project files from a templates

#### fl mod

The `mod` command is used for anything module related. It has currently 5 subcommands:

- `create`: to create a new module
- `delete`: to delete a module
- `get`: to get information on a module (name and functions inside)
- `list`: to list all modules
- `update`: to update a module's name

#### fl template

The `template` command is used to pull or list templates. It has currently 2 subcommands:

- `pull`: to pull a template from a repository
- `list`: to list all templates in the current folder

## Installation

Right now only the linux version of the CLI is supported, although we build and release the CLI for windows and macos as well, therefore 
they are not guaranteed to work (if anyone wants to try on those platforms and give feedback/help we'd appreciate it).

The tool is in the [Release page](https://github.com/funlessdev/fl-cli/releases) of this repository, and can be downloaded from there.
In linux, just move the binary to a folder in your `$PATH` and you're good to go.

## Usage

We added the docker deployment in the tool to have an easy quick start. It uses `docker compose` under the hood so you
need to have a recent version of docker installed and running (you can find a more in depth quick start at [funless.dev](https://funless.dev/)). All you have to do is run:

```bash
fl admin deploy docker up
```

This will deploy the platform on your local machine, together with some helper services to handle the logs (Kibana with Elasticsearch). You can access Kibana interface at `localhost:5601`.

The core services, instead, are:

- the **core** container which exposes the JSON API to interact with the platform (at `localhost:4000`)
- the **worker** container which handles the execution of the functions (only 1 worker in the local deploy)
- the **postgres** container which is the database used by the platform
- the **prometheus** container which is the monitoring service used by the platform (at `localhost:9090`)

To remove the platform from your machine, just run

```bash
fl admin deploy docker down
```

### Using custom images

If you are working on the core or worker, there is an easy way to quick start the platform with your custom components.
You need to build the images for the core and/or the worker and then specify them in the `docker up` command:

```bash
fl admin deploy docker up --core <your-core-image> --worker <your-worker-image>
```

## Contributing

Anyone is welcome to contribute to this project or any other FunLess project. 

You can contribute by testing the projects, opening tickets, writing documentation, sharing new ideas for future works and, of course,
by contributing code. 

You can pick an issue or create a new one and fork the repo so you're free to work on it.
Once you feel ready open a Pull Request to send your code to us.

## License

This project is under the Apache 2.0 license.
