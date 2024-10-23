SHELL := /bin/bash
TARGETS := readms

.PHONY: all
all: $(TARGETS)

%: %.go
	go build -o $@ $^

.PHONY: clean
clean:
	rm -f $(TARGETS)
