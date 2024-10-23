#!/bin/bash

# Run and record profiles, convert them to flamegraphs, etc.

rm -f time.log
for f in $(find . -maxdepth 1 -name "1brc-*"| sort -r); do
    echo "running $f..."
    NAME=$(basename $f)
    echo $f >> time.log
    { time $f ; } 2>> time.log
    echo >> time.log
    $f -cpuprofile cpu.pprof && bash pprof_to_svg.sh && cp cpu.svg static/$NAME.svg
    rm -f cpu.svg
done
