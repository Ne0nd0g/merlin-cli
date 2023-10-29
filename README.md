[![Build](https://github.com/Ne0nd0g/merlin-cli/actions/workflows/go.yml/badge.svg)](https://github.com/Ne0nd0g/merlin-cli/actions/workflows/go.yml)
[![GoReportCard](https://goreportcard.com/badge/github.com/Ne0nd0g/merlin-cli)](https://goreportcard.com/report/github.com/Ne0nd0g/merlin-cli)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Release](https://img.shields.io/github/release/Ne0nd0g/merlin-cli.svg)](https://github.com/Ne0nd0g/merlin-cli/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/Ne0nd0g/merlin-cli/total.svg)](https://github.com/Ne0nd0g/merlin-cli/releases)
[![Qodana](https://github.com/Ne0nd0g/merlin-cli/actions/workflows/qodana.yml/badge.svg)](https://github.com/Ne0nd0g/merlin-cli/actions/workflows/qodana.yml)
[![CodeQL](https://github.com/Ne0nd0g/merlin-cli/actions/workflows/codeql.yml/badge.svg)](https://github.com/Ne0nd0g/merlin-cli/actions/workflows/codeql.yml)
[![Twitter Follow](https://img.shields.io/twitter/follow/merlin_c2.svg?style=social&label=Follow)](https://twitter.com/merlin_c2)

# Merlin CLI

Merlin is composed of the following components:

* Merlin Server - The program that receives and handles Agent traffic and operator CLI commands to control the server and Agents
* Merlin Agent - The post-exploitation command and control Agent that runs on a compromised host
* Merlin CLI - The command line interface that allows operators to interact with the Merlin Server and Agents

**This repository covers the Merlin Command Line Interface (CLI) program**

Releases from the <https://github.com/Ne0nd0g/merlin/releases> page contain the Merlin Server, Agent, and CLI programs.

Merlin documentation can be found at <https://merlin-c2.readthedocs.io/en/latest/>

> The Main Merlin C2 repository can be found here: <https://github.com/Ne0nd0g/merlin>

## Command Line Interface

The CLI uses Google RPC (gRPC) protocol buffers over TLS to communicate with the Merlin Server.
All API calls require a password to authenticate to the server.

> **WARNING:** The default password is `merlin` and should always be changed to prevent unauthorized access

> **NOTE:** By default, the Merlin Server will generate a self-signed TLS certificate that will not be trusted by the CLI if the `secure` flag is used

```text
    $ ./merlin-cli -h
    Usage of merlin-cli:
      -addr string
            The address of the Merlin server to connect to (default "127.0.0.1:50051")
      -password string
            the password to connect to the Merlin server (default "merlin")
      -secure
            Require server TLS certificate verification
      -tlsCA string
            TLS Certificate Authority file path
      -tlsCert string
            TLS certificate file path
      -tlsKey string
            TLS private key file path
      -version
            Print the version number and exit
```

* `addr` - this flag specifies the address of the Merlin Server to connect to. The connection uses gRPC over TLS.
* `password` - this flag sets the password needed to authenticate all gRPC requests
* `secure` - this flag enables TLS certificate verification. When this flag is set, the CLI will verify the Server's TLS certificate
* `tlsCA` - this flag specifies a custom CA certificate file to validate and trust the Server's certificate
* `tlsCert` - this flag specifies the certificate file the Merlin CLI will use for mutual TLS authentication with the Merlin Server
* `tlsKey` - this flag specifies the private key file for the `tlsCert`
* `version` - this flag prints the version number of the Merlin Server and exits
