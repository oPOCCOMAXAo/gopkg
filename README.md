# gopkg

Library of packages for quick use.

## example

- [config](example/config) - example app config with some services;
- [app](example/app) - example app run with config initiation, dependencies building and wiring and lifecycle management.

## pkg/services

Each package with services contains:

- Config - struct with `env` tags for package `github.com/caarlos0/env/v11`;
- Package - struct with all exported services, grouped for easier dependency injection;
- MakePackage - constructor for Package, entrypoint.

Each dependency has public constructor and can be used without entire package build.

### General services

- [logger](pkg/services/logger) - entrypoint for slog.Logger.
  Please adopt slog.Logger and construct it from single entrypoint.

- [lifecycle](pkg/services/lifecycle) - app lifecycle controller, implement interface for servable/shutdownable services and register them here.
  Example of use: [ginserver](pkg/services/ginserver).

Read app example: [Run()](example/app/run.go).

## pkg/utils

Each package contains utility functions/structs with could be used without constructor and lifecycle.

IMPORTANT: utils mustn't use pkg/services.
