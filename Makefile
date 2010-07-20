include $(GOROOT)/src/Make.$(GOARCH)

TARG=trie

GOFILES=\
	trie.go\
	hyphen_trie.go\

include $(GOROOT)/src/Make.pkg