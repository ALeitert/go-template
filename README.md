# Go Template

A basic template for a Go program with the following features:
- Makefile with linting, tidy check, unit test, code genration, and building of a Docker container.
- GitHub workflow with basic checks for PRs.
- Framework to run multiple services in their own goroutine.
  That is useful when, for example, running an HTTP service which provides metrics to Prometheus in addition to the program's main logic.
- Configuration which is loaded from a YAML file.
  It supports default values in case in YAML file is incomplete.
- Database layer using [sqlc](https://sqlc.dev/).
- Server for Prometheus metrics.
- REST API server utilising OpenAPI.

Schemas for generating code are in the `schemas/` directory.
Currently that includes SQL migrations and queries, as well as OpenAPI definitions.
