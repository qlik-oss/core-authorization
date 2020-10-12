# Qlik Core Authorization Examples

[![CircleCI](https://circleci.com/gh/qlik-oss/core-authorization.svg?style=svg)](https://circleci.com/gh/qlik-oss/core-authorization)

*As of 1 July 2020, Qlik Core is no longer available to new customers. No further maintenance will be done in this repository.*

This repo contains running examples showing how to configure Qlik Associative Engine for ABAC rule evaluation, and
examples on how some basic rules can be written.

The [Qlik Core](https://core.qlik.com/) documentation provides more in-depth technical documentation on how rules are
written and how the Qlik Associative Engine is configured for ABAC rule evalutaion. It is recommended to use that
documentation along with these examples.

The examples are running as Go test cases, where each test case verifies that a certain user is properly granted or
denied access, based on the rules provided to the engine instances.

## Engines and Users

Two engines are used with different types of users and rule sets. Both these engines have ABAC enabled and different
rule files provided to them. Check the [docker-compose.yml](./docker-compose.yml) for details on which command line
switches that are needed, and how rule files are provided to each engine.

### The Empty Engine

The **empty-engine** instance starts completely empty, with no documents or data available. On this engine, two kinds of
users are assumed:

1. _Admin_ users - These users have full access to everything on the engine.
1. _View_ users - These users have read-only access to documents and objects on the engine.

The file [empty_engine_rules.txt](./rules/empty_engine_rules.txt) is provided to enforce the access rules on this
engine.

### The Reload Engine

The **reload-engine** instance starts with no loaded data. A CSV file containing information on some famous movies is
provided. On this engine, three kinds of users are assumed:

1. _Create_ users - These users can create, update, and read documents but thay _cannot_ reload data.
1. _Reload_ users - These users can reload data into documents but they _cannot_ create documents.
1. _View_ users - These users have read-only access to documents and objects on the engine.

The file [reload_engine_rules.txt](./rules/reload_engine_rules.txt) is provided to enforce the access rules on this
engine.

## Access and Tests

To perform operations that shall either be granted or denied towards, functions in the file
[access.go](./access/access.go) are used. These functions use the [enigma-go](https://github.com/qlik-oss/enigma-go)
library to communicate with the engine instances.

To verify the expected behavior, the file [access_test.go](./access/access_test.go) contains Go test cases that call
the access functions with different claims based on which user that accesses the engine. Each test case verifies that
operations are properly granted or denied.

## Running

Running the tests is containerized. [Docker](https://www.docker.com/) and
[docker-compose](https://docs.docker.com/compose/) are required to run the tests. The
[End-User License Agreement](https://core.qlik.com/eula/) must be read and accepted.

Assuming that the current working directory is the repo root, set the `ACCEPT_EULA` environment variable as
appropriate and run the tests with:

```sh
ACCEPT_EULA=<yes/no> docker-compose run test-runner
```

The output should show all tests executed successfully. The command leaves two engine instances running. To remove
them, run:

```sh
docker-compose down
```
