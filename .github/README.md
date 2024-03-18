<h1 align="center">Go Password Manager</h1>

<p align="center">
  <a href="https://github.com/raulaguila/go-pass/releases" target="_blank" style="text-decoration: none;">
    <img src="https://img.shields.io/github/v/release/raulaguila/go-pass.svg?style=flat&labelColor=0D1117">
  </a>
  <img src="https://img.shields.io/github/repo-size/raulaguila/go-pass?style=flat&labelColor=0D1117">
  <img src="https://img.shields.io/github/stars/raulaguila/go-pass?style=flat&labelColor=0D1117">
  <a href="../LICENSE" target="_blank" style="text-decoration: none;">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat&labelColor=0D1117">
  </a>
  <a href="https://goreportcard.com/report/github.com/raulaguila/go-pass" target="_blank" style="text-decoration: none;">
    <img src="https://goreportcard.com/badge/github.com/raulaguila/go-pass?style=flat&labelColor=0D1117">
  </a>
  <a href="https://github.com/raulaguila/go-pass/actions?query=workflow%3Ago-test" target="_blank" style="text-decoration: none;">
    <img src="https://github.com/raulaguila/go-pass/actions/workflows/go-test.yml/badge.svg">
  </a>
  <a href="https://github.com/raulaguila/go-pass/actions?query=workflow%3Ago-build" target="_blank" style="text-decoration: none;">
    <img src="https://github.com/raulaguila/go-pass/actions/workflows/go-build.yml/badge.svg">
  </a>
</p>

## Requirements/dependencies

--

## Getting Started

- Help with make command

```sh
Usage:
        make [COMMAND]
        make help

Commands:

help                           Display help screen
init                           Create environment variables
compose-up                     Run docker compose up for create and start containers
compose-build                  Run docker compose up build for create and start containers
compose-down                   Run docker compose down for stopping and removing containers, networks
compose-remove                 Run docker compose down for stopping and removing containers, networks, volumes
compose-exec                   Run docker compose exec to access bash container
compose-log                    Run docker compose logs to show logger container
compose-top                    Run docker compose top to display the running containers processes
go-fmt                         Run go fmt
go-vet                         Run go vet
go-test                        Run go test
go-test-cover                  Run go test with coverage report
go-test-cover-html             Run go test with HTML coverage report
```

- Build and run project

```sh
make init
# Change db authentication parameters in `.env` file created on `/configs`
make compose-build
```

- Remove project

```sh
make compose-remove
```

## Features

--

## Code status

- Development

## Contributors

<a href="https://github.com/raulaguila" target="_blank">
  <img src="https://contrib.rocks/image?repo=raulaguila/go-template">
</a>

## License

Copyright © 2023 [raulaguila](https://github.com/raulaguila).
This project is [MIT](../LICENSE) licensed.
