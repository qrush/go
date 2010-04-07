cmds = mk
pkg = mk
g = mk.go

GCFLAGS=-I ../dag/bin
LDFLAGS=-L ../dag/bin

include $(GOROOT)/src/Make.$(GOARCH)
include ../Make.pkg

test: all; ./test/run
