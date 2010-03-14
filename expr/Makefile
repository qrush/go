include $(GOROOT)/src/Make.$(GOARCH)

TARG=bin/expr
GOFILES=\
        expr.go

include $(GOROOT)/src/Make.cmd

doc:	;godoc -path=.. ../expr
test:	;./expr_test
