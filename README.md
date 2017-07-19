# Dairycart  [![Build Status](https://travis-ci.org/dairycart/dairycart.svg?branch=master)](https://travis-ci.org/dairycart/dairycart) [![codecov](https://codecov.io/gh/dairycart/dairycart/branch/master/graph/badge.svg)](https://codecov.io/gh/dairycart/dairycart)

Dairycart is an open-source eCommerce platform written in Go.

## Status

Dairycart is currently pre-Alpha, and isn't suitable for production use as of yet. There is no stable, defined API, and until version 1.0 is released and tagged, any route should be considered volatile and subject to change with no notice.

## Dependencies

To run Dairycart locally, all you need is Docker, and a bash terminal. :simple-smile:

## Running

To run the Dairycart API server, simply execute `debug.sh` after cloning this repository locally.

## Testing

There are two test suites, currently, a set of unit tests for the API server, and a suite of integration tests, which attempt to test the API server with real-world style requests. You can run the unit tests by executing `run_unit_tests.sh`, and you can run the integraiton tests similarly by executing `run_integration_tests.sh`

## Trello

The Trello board for this project can be found [here](https://trello.com/b/z3lgKd59/dairycart)

## Reporting Issues

Thusfar, Dairycart's development has been a solo effort, and this may still be the case for some time. That being the case, there are most definitely issues with the
code that warrant reporting. Any issues can be reported using the [Github Issue Tracker](https://github.com/dairycart/dairycart/issues/new) for this repository. Please ensure any issues filed adhere to the [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting involved

General instructions on _how_ to contribute can be found in [CONTRIBUTING](CONTRIBUTING.md).

## Open source licensing info

This project is licensed under the [MIT License](https://en.wikipedia.org/wiki/MIT_License)

A big thanks to the [CFPB](https://github.com/cfpb/open-source-project-template) for the template this file is based on.