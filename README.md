# deisrel

[![Build Status](https://travis-ci.org/deis/deisrel.svg?branch=master)](https://travis-ci.org/deis/deisrel)
[![codebeat badge](https://codebeat.co/badges/46e06b60-7e4c-4daf-875b-c7c07ee56035)](https://codebeat.co/projects/github-com-deis-deisrel)
[linux-amd64 Download](https://storage.googleapis.com/deisrel/deisrel-latest-linux-amd64)
[darwin-amd64 Download](https://storage.googleapis.com/deisrel/deisrel-latest-darwin-amd64)

Deis (pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS) that adds a
developer-friendly layer to any [Kubernetes](http://kubernetes.io) cluster, making it easy to
deploy and manage applications on your own servers.

For more information about the Deis Workflow, please visit the main project page at
<https://github.com/deis/workflow>.

We welcome your input! If you have feedback, please [submit an issue][issues]. If you'd like to participate in development, please read the "Development" section below and [submit a pull request][prs].

# About

`deisrel` is a utility tool for automating Deis product releases. The idea is that it provides a
way to automate the release process for Deis Workflow without human intervention, eventually being
able to automate the release process through a Continuous Integration server like <https://ci.deis.io/>.

# Installing deisrel

You can install `deisrel` by installing the latest release from [bintray](https://bintray.com/):

[linux-amd64](https://storage.googleapis.com/deisrel/deisrel-latest-linux-amd64)
[darwin-amd64](https://storage.googleapis.com/deisrel/deisrel-latest-darwin-amd64)

Alternatively, you can compile this project from source using Go 1.6+:

	$ git clone https://github.com/deis/deisrel
	$ cd deisrel
	$ make bootstrap build
	$ ./deisrel

Once done, you can then move the client binary anywhere on your PATH:

	$ mv deisrel /usr/local/bin/

# Usage

In order to use `deisrel`, you must first add a [GitHub access token](https://github.com/settings/tokens) to your environment:

	$ export GITHUB_ACCESS_TOKEN="myaccesstoken"

Then use `deisrel help` to explore the commands more in-depth.

For example, to generate an aggregated changelog for [Deis Workflow][workflow]:

	$ deisrel changelog global v2.0.0-beta3 v2.0.0-beta4

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
