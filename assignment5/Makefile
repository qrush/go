include $(GOROOT)/src/Make.$(GOARCH)

TARG=matrix

GOFILES=\
	matrix.go\

doc: ;godoc -path=.. .
htmldoc: ;godoc -html=true -path=.. . > doc.html

include $(GOROOT)/src/Make.pkg
