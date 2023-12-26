# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

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