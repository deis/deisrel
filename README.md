# deisrel

[![Build Status](https://travis-ci.org/deis/deisrel.svg?branch=master)](https://travis-ci.org/deis/deisrel)
[![codebeat badge](https://codebeat.co/badges/46e06b60-7e4c-4daf-875b-c7c07ee56035)](https://codebeat.co/projects/github-com-deis-deisrel)

[Download for 64 Bit Linux](https://storage.googleapis.com/deisrel/deisrel-latest-linux-amd64)

[Download for 64 Bit Darwin](https://storage.googleapis.com/deisrel/deisrel-latest-darwin-amd64)

Deis (pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS) that adds a
developer-friendly layer to any [Kubernetes](http://kubernetes.io) cluster, making it easy to
deploy and manage applications on your own servers.

For more information about the Deis Workflow, please visit the main project page at
<https://github.com/deis/workflow>.

We welcome your input! If you have feedback, please [submit an issue][issues]. If you'd like to participate in development, please read the "Development" section below and [submit a pull request][prs].

# About

`deisrel` is a utility tool for assisting with Deis product releases. It makes is easy to check what components need to be released before making a workflow release.

# Installing deisrel

You can install the latest version of `deisrel` from the following links:

- [linux-amd64](https://storage.googleapis.com/deisrel/deisrel-latest-linux-amd64)
- [darwin-amd64](https://storage.googleapis.com/deisrel/deisrel-latest-darwin-amd64)

Alternatively, you can compile this project from source using Go 1.6+:

	$ git clone https://github.com/deis/deisrel
	$ cd deisrel
	$ make bootstrap build
	$ ./deisrel

Once done, you can then move the client binary anywhere on your PATH:

	$ mv deisrel /usr/local/bin/

# Usage

deisrel requires two files, the `generate_params.toml` file from [deis/charts](https://github.com/deis/charts), and a repository mapping file, allowing deisrel to map between a repository and the component name.

Since the deis repositories are changing rapidly, right now the repository mapping file should be manually created.

Here's an example file:

```json
{
  "builder": ["builder"],
  "controller": ["controller"],
  "dockerbuilder": ["dockerbuilder"],
  "fluentd": ["fluentd"],
  "monitor": ["influxdb", "grafana", "telegraf"],
  "logger": ["logger"],
  "minio": ["minio"],
  "nsq": ["nsqd"],
  "postgres": ["database"],
  "redis": ["loggerRedis"],
  "registry": ["registry"],
  "router": ["router"],
  "slugbuilder": ["slugbuilder"],
  "slugrunner": ["slugrunner"],
  "workflow-manager": ["workflowManager"]
}
```

The key is the repository name, and the values are the names of components using that repository, as found in the `generate_params.toml` file.

With these two files, you can use deisrel to generate a report:

```console
$ deisrel path/to/generate_params.toml path/to/repomapping.json
builder         v2.2.0 -> v2.2.0 (dirty)
	builder has unrelased changes. See https://github.com/deis/builder/compare/v2.2.0...master
controller      v2.2.0 -> v2.2.1 (dirty)
	controller has unrelased changes. See https://github.com/deis/controller/compare/v2.2.1...master
database        v2.2.0 -> v2.2.0 (clean)
dockerbuilder   v2.2.0 -> v2.2.0 (clean)
fluentd         v2.2.0 -> v2.2.0 (clean)
grafana         v2.2.0 -> v2.2.0 (clean)
influxdb        v2.2.0 -> v2.2.0 (clean)
logger          v2.2.0 -> v2.2.0 (clean)
loggerRedis     v2.2.0 -> v2.2.0 (clean)
minio           v2.2.0 -> v2.2.0 (clean)
nsqd            v2.2.0 -> v2.2.0 (clean)
registry        v2.2.0 -> v2.2.0 (clean)
router          v2.2.0 -> v2.2.0 (clean)
slugbuilder     v2.2.0 -> v2.2.0 (clean)
slugrunner      v2.2.0 -> v2.2.0 (clean)
telegraf        v2.2.0 -> v2.2.0 (clean)
workflowManager v2.2.0 -> v2.2.0 (clean)
```

It's also possible to output the report in json using the `-o json` argument for easier machine parsing.

Github has a very low ratelimit for unauthenticated api requests. A personal oauth token can be used to bypass this restriction. Create an oauth token with no permission and set it to `$GH_TOKEN` and deisrel will
use it when making requests.

# Development

The Deis project welcomes contributions from all developers. The high level process for development matches many other open source projects. See below for an outline.

* Fork this repository
* Make your changes
* [Submit a pull request][prs] (PR) to this repository with your changes, and unit tests whenever possible
	* If your PR fixes any [issues][issues], make sure you write `Fixes #1234` in your PR description (where `#1234` is the number of the issue you're closing)
* The Deis core contributors will review your code. After each of them sign off on your code, they'll label your PR with `LGTM1` and `LGTM2` (respectively). Once that happens, a contributor will merge it

# License

Copyright 2013, 2014, 2015, 2016 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.


[issues]: https://github.com/deis/deisrel/issues
[prs]: https://github.com/deis/deisrel/pulls
[workflow]: https://github.com/deis/workflow
