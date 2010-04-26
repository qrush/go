cmds = ttt
pkg = ttt
g = ttt.go

GCFLAGS=-I ../games/bin
LDFLAGS=-L ../games/bin

include $(GOROOT)/src/Make.$(GOARCH)
include ../Make.pkg

test: all; ./test/run
