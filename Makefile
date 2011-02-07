include $(GOROOT)/src/Make.inc

TARG=github.com/jbarham/cdb.go

GOFILES=\
	cdb.go\
	dump.go\
	make.go\
	hash.go\

include $(GOROOT)/src/Make.pkg
