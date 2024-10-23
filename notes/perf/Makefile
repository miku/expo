SHELL := /bin/bash

perf.data: add_vec
	perf record -g -o perf.data ./add_vec

add_vec: add_vec.cpp
	# -g -- Produce debugging information in the operating system's native
	#  format (stabs, COFF, XCOFF, or DWARF).
	g++ -O2 -g -o add_vec add_vec.cpp

.PHONY: clean
clean:
	rm -f add_vec
	rm -f perf.data
	rm -f perf.data.old
