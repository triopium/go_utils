# Go Utils
Go Utils are packages which helps with developement of golang programs and libraries

## Go Configure

Go configure is library which helps:

- define cli flags, cli subcommands, environment variables, configurataion files.
- establish priority between various sources of variables
- populate arbitrary struct with effective variables

## Go Tester

Go tester helps creating temporary directory and populate it with
source test files.  When running in manual mode it waits for user
input after the tests ends or error is encountered before the
temporary directory is cleaned up. Useful when you need to check
test artifacts created during test.

## Go Helper

Various helper functions which helps as shorthand for common tasks.

### Go workspace usage

-add repo dir to workspace file
go.work
-------
go 1.23.2

use (
	./fulltext
	./openmedia
	./omserve
	~/runs/golang/go2023/go_utils
)

go work sum
