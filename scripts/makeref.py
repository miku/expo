#!/usr/bin/env python

import glob

for path in sorted(glob.glob("static/1brc*svg")):
    name = path.replace("static/", "").replace(".svg", "")
    print(f"### {name}")
    print()
    print(f"![]({path})")
    print()
