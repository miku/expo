SHELL := /bin/bash
MAKEFLAGS := --jobs=$(shell nproc)
TARGETS := gen1brcdata \
    1brc-scratch \
	1brc-000-baseline \
	1brc-001-baseline-read-limitmem \
	1brc-002-plain-scanner \
	1brc-005-baseline-scan \
	1brc-010-baseline-scan-tweak \
	1brc-020-fanout \
	1brc-030-fanout-scanner \
	1brc-040-file-partition \
	1brc-050-mmap \
	1brc-060-mmap-float \
	1brc-060-mmap-int \
	1brc-070-mmap-int-tweaks \
	1brc-075-mmap-int-extra \
	1brc-076-mmap-int-za-key \
	1brc-080-mmap-int-static-map \
	1brc-081-mmap-int-static-map-za-key \
	1brc-082-mmap-faster-int-static-map-za-key \
	1brc-401-baseline \
	1brc-402-avoid-double-hashing \
	1brc-403-avoid-parse-float \
	1brc-404-temp-int32 \
	1brc-405-avoid-cut \
	1brc-406-no-scanner \
	1brc-407-custom-hash-table \
	1brc-408-parallel-baseline \
	1brc-409-parallel-opt \
	1brc-410-fast-semi \

.PHONY: all
all: $(TARGETS)

%: cmd/%/main.go
	go build -o $@ $<

.PHONY: clean
clean:
	rm -f $(TARGETS)
	rm -f cpu.txt cpu.png cpu.svg cpu.pprof

measurements.txt: gen1brcdata
	./gen1brcdata > measurements.txt

cities.txt: measurements.txt
	cut -d ';' -f 1 measurements.txt | LC_ALL=C sort -S30% -u > cities.txt

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. ./...
