cmds = ttt
pkg = ttt
g = ttt.go

GCFLAGS=-I ../games/bin
LDFLAGS=-L ../games/bin

include $(GOROOT)/src/Make.$(GOARCH)
include ../Make.pkg

test: all; ./bin/ttt -a localhost:1234 localhost:1235

player1: test

player2:;./bin/ttt localhost:1235 localhost:1234
