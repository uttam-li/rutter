#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

rsvg-convert -w 512 -h 512 logo.svg -o logo-512.png
rsvg-convert -w 180 -h 180 favicon.svg -o apple-touch-icon.png
rsvg-convert -w 32 -h 32 favicon.svg -o favicon-32.png
rsvg-convert -w 16 -h 16 favicon.svg -o favicon-16.png

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
rsvg-convert -w 48 -h 48 favicon.svg -o "$tmp/favicon-48.png"

python3 - "$tmp/favicon-48.png" <<'PY'
import struct
import sys
from pathlib import Path

sizes = [(16, Path("favicon-16.png")), (32, Path("favicon-32.png")), (48, Path(sys.argv[1]))]
images = [(size, path.read_bytes()) for size, path in sizes]

header = struct.pack("<HHH", 0, 1, len(images))
offset = 6 + 16 * len(images)
entries = b""
for size, data in images:
    entries += struct.pack("<BBBBHHII", size, size, 0, 0, 1, 32, len(data), offset)
    offset += len(data)

Path("favicon.ico").write_bytes(header + entries + b"".join(data for _, data in images))
PY
