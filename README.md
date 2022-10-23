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

# fl

This is the repository for the CLI of the Funless (FL) platform, a new generation research-driven serverless platform.

The CLI can be used both to deploy the platform and to create, delete and invoke functions on it.

## Using the CLI

The CLI currently exposes two sets of commands, `fn` and `admin`, for function manipulation and deployment respectively.
### fn

The `fn` commands are used to create, delete and invoke functions.

#### `fn create`

The `fn create` command is used to create functions and store them in the platform's permanent storage.

It takes a mandatory argument (the function's name) and requires either a `--source-file` or a `--source-dir` additional parameter, for the function's source.

In case a source directory is passed, a suitable container from [fl-runtimes](https://github.com/funlessdev/fl-runtimes) is pulled and used to build the source.



#### `fn delete`

The `fn delete` command is used to remove functions from the platform's permanent storage.



#### `fn invoke`

The `fn invoke` command is used to run functions on the platform. Both keyword arguments and json arguments can be passed using the `-a` or the `-j` flag.



### admin

The `admin` commands are used to deploy and delete an instance of the platform; currently the only deployment option available is a development instance.

All `admin` commands can be used with the `a` shortcut, so `fl admin dev` is the same as `fl a dev`.

#### `admin dev`

The `admin dev` command is used to spin up a development version of the platform, with 1 Core and 1 Worker, on the local machine.

Custom images for Core and Worker can be used with the `--core` and `--worker` flags.


#### `admin reset`

The `admin reset` command is used to delete all containers of the development installation of the platform; pulled images will not be deleted.


## Contributing

Anyone is welcome to contribute to this project or any other Funless project. 

You can contribute by testing the projects, opening tickets, writing documentation, sharing new ideas for future works and, of course,
by contributing code. 

You can pick an issue or create a new one, comment on it that you will take priority and then fork the repo so you're free to work on it.
Once you feel ready open a Pull Request to send your code to us.

## License

This project is under the Apache 2.0 license.