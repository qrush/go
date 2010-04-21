include $(GOROOT)/src/Make.$(GOARCH)

TARG=wagon
GOFILES=wagon.go

include $(GOROOT)/src/Make.cmd

doc:	;godoc -path=.. .

htmldoc: ;godoc -html=true -path=.. . > doc.html
test: ;./wagon
