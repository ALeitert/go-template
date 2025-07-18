# Go Template

A basic template for a Go program with the following features:
- A makefile with linting, tidy check, unit test, and building of a Docker container.
- A GitHub workflow with basic checks for PRs.
- A framework to run multiple services in their own goroutine.
  That is useful when, for example, running an HTTP service which provides metrics to Prometheus in addition to the program's main logic.
- A configuration which is loaded from a YAML file.
  It supports default values in case in YAML file is incomplete.
