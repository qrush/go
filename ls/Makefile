include $(GOROOT)/src/Make.$(GOARCH)

TARG=bin/ls
GOFILES=main.go

include $(GOROOT)/src/Make.cmd

doc: ;godoc -path=.. ../ls
htmldoc: ;godoc -html=true -path=.. ../ls > doc.html
test: all; ./test/run
