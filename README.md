# simple-order-api
This project demonstrates unit testing practices for a simple order api that performs crud operations. It shows effective ways and gives examples to write a unit test with these concepts
- testint.T package
- Mocking
- Mocking with Mockery
- Table-Driven Tests
- Test Suite Concept

## How you run project and unit tests
All you have to do is clone the project with git command and open it in a suitable ide and run **go mod vendor** command.  Here i recommend for ide [Goland](https://www.jetbrains.com/go/) so you can run tests with its facilities. If you choose to run tests without an available ide, you can run **go test** command within tests file in spesific directories.
Makefile is also in project and for creating mock types, all you do is running "make mock" command.

## Dependency
It depends go 1.18 version 

## What's contained in this project
- Swagger - Use swag init command for creating swag presentation. [Swaggo](https://github.com/swaggo/swag) is used for that
- [Gin-gonic](https://github.com/gin-gonic/gin) is used for http web framework.
- Testing - Tests are implemented with [Testing Package](https://pkg.go.dev/testing) support.
- [Stretch Testify Package](https://github.com/stretchr/testify) is used for mocking, assertions and suite.
- [Mockery](https://github.com/vektra/mockery) is used for mock types automatically without writing any code.
