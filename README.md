Job Worker [![CI Status](https://github.com/andrejtokarcik/jobworker/workflows/CI/badge.svg)](https://github.com/andrejtokarcik/jobworker/actions)
==========

A prototype gRPC service providing OS process execution and monitoring,
written in Go.

Build & CI
----------

The default `make` workflow uses the multi-stage Docker build feature in order
to set up an environment with tool versions specified in the `build.env` file,
ensuring fast, predictable and reproducible results.

The basic available `make` targets are:

  - `build` (default)
  - `test`
  - `gen` to generate Go files for Protobuf & mocks
  - `check-gen` to verify that the generated files, the Go module
    dependencies, etc., are in sync with the codebase

The actual build steps performed by Docker are defined in `Makefile.local`.
By invoking

    make -f Makefile.local

it is therefore possible to avoid the intermediate container layers and to
perform the actions associated with a build target directly on the local host
system.

Finally, one can also locally execute the GitHub CI Actions using `make ci`
provided the [`act` tool](https://github.com/nektos/act) is installed on the
system.  This is especially useful to do as a last general check before pushing
to the mainline repository.
