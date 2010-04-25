cmds = rps
pkg = rps
g = rps.go

GCFLAGS=-I ../games/bin
LDFLAGS=-L ../games/bin

include $(GOROOT)/src/Make.$(GOARCH)
include ../Make.pkg

test: all; ./test/run
