PREFIX = /usr/local
BINDIR =  $(PREFIX)/bin
ETCDIR = /etc

CMDS := $(notdir $(wildcard cmd/*))

GO = go
VERSION = 0.0.1
NAME = sboutils-$(VERSION)

all : $(CMDS)

$(CMDS):
	@echo "Compiling $@..."
	go build -o bin/$@ ./cmd/$@

.PHONY: all $(COMMANDS)
