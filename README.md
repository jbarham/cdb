# cdb.go - By John Barham

cdb.go is a [Go](http://golang.org/) package to read and write cdb ("constant database") files.

See the original cdb specification and C implementation by D. J. Bernstein
at http://cr.yp.to/cdb.html.

## Installation

Assuming you have a working Go environment, installation is simply:

	go get github.com/jbarham/cdb.go

The package documentation can be viewed online at
http://gopkgdoc.appspot.com/pkg/github.com/jbarham/cdb.go
or on the command line by running `go doc github.com/jbarham/cdb.go`.

The included self-test program `cdb_test.go` illustrates usage of the package.

## Utilities

The cdb.go package includes ports of the programs `cdbdump` and `cdbmake` from
the [original implementation](http://cr.yp.to/cdb/cdbmake.html).