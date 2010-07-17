include $(GOROOT)/src/Make.$(GOARCH)

TARG=trie

GOFILES=\
    defs.go\
	trie.go\
	value_trie.go\

include $(GOROOT)/src/Make.pkg