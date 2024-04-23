# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## 1.1.4 - 2024-04-23

### Changed

- Upgraded golang.org/x/sys v0.18.0 => v0.19.0
- Upgraded google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda => v0.0.0-20240415180920-8c6c420018be
- Upgraded google.golang.org/grpc v1.63.0 => v1.63.2
- Upgraded github.com/rivo/uniseg v0.4.4 => v0.4.7

### Security

- [CVE-2023-45288](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-45288) - Upgraded `golang.org/x/net` to v0.24.0

## 1.1.3 - 2024-04-04

### Changed

- upgraded github.com/google/uuid v1.5.0 => v1.6.0
- upgraded google.golang.org/grpc v1.60.1 => v1.63.0

### Security

- [CVE-2024-24786](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2024-24786) - Upgraded `google.golang.org/protobuf` to v1.33.0

## 1.1.2 - 2023-12-26

### Added

- Added a Dockerfile for the CLI

### Changed

- Makefile used `DIR` variable and builds artifacts into a `bin` directory
- Upgraded:
  - `github.com/fatih/color` to `v1.16.0`
  - `github.com/google/uuid` to `v1.5.0`
  - `google.golang.org/grpc` to `1.60.1`
- Updated GitHub Action workflows

## 1.1.1 - 2023-12-14

### Added

- Attempt to recover broken gRPC message stream if lost

## 1.1.0 - 2023-11-04

### Added

- Added a `parrot` command to the CLI

## 1.0.0 - 2023-10-29

- Initial release