#!/bin/bash

PPROF=${1:-cpu.pprof}
SVG=${2:-cpu.svg}

go tool pprof -raw -output=cpu.txt "$PPROF"
stackcollapse-go.pl cpu.txt | flamegraph.pl > "$SVG"
