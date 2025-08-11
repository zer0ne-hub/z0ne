# z0ne

Attack Surface Mapper and Reconnaissance Framework

## About

A modular attack surface mapper designed for CTFs and pentesting. A simple orchestrator of
tools for target reconnaissance. Most tools are from the awesome [projectdiscovery](https://projectdiscovery.io/) project, making this a "glorified" golang wrapper around them.

This tool is good for getting a quick overview of a target and its attack surface.
Its features include port scanning, subdomain enumeration, DNS resolution, and more. It can
give you a quick headstart in your reconnaissance phase and generate a simple markdown
report of its findings.

## Installation

Intended to be a go package, part of a larger framework for
[zer0ne](https://github.com/zer0ne-hub/z0ne), this individual tool can be installed via:

```bash
go install github.com/zer0ne-hub/z0ne@latest
```
Make sure you have [Go](https://go.dev/) installed on your machine and its GOPATH correctly
set up for binaries.

## Usage

After installing the tool, you can use it by running `z0ne` or `z0ne --help` to get a quick
overview of its usage.

### Deps

For those interested in the development of this tool, you can check out the dependencies in the
`go.mod` file but mostly the following:

- [Cobra](https://cobra.dev/) is used for command line parsing and execution
- [project discovery](https://projectdiscovery.io/) tools (naabu, subfinder, dnsx, httpx, katana, nuclei, uncover) are used in the `recon` module and orchestrated by the `core` module
- [taskfile](https://taskfile.dev/) is used to automate some build and run tasks
