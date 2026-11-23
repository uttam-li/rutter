# Brand assets

The mark is a compass star: a mint bearing, a blue heading, and a transparent pivot. It comes from the older meaning of rutter, a book of routes for a voyage.

| File | Use |
| --- | --- |
| `logo.svg` | Primary mark, 512×512 dark tile |
| `logo-512.png` | Raster version for READMEs and previews |
| `favicon.svg` | Scalable favicon, simplified and enlarged for small sizes |
| `favicon.ico` | Legacy favicon with 16, 32, and 48 px images |
| `favicon-16.png`, `favicon-32.png` | PNG favicon fallbacks |
| `apple-touch-icon.png` | 180×180 iOS home-screen icon |
| `build.sh` | Regenerates every PNG and the ICO from the two SVGs |

## HTML

```html
<link rel="icon" href="/favicon.ico" sizes="any">
<link rel="icon" href="/favicon.svg" type="image/svg+xml">
<link rel="apple-touch-icon" href="/apple-touch-icon.png">
```

## Colors

| Role | Hex |
| --- | --- |
| Tile top | `#121820` |
| Tile bottom | `#080b0f` |
| Bearing | `#8df3c0` |
| Heading | `#5ea8ff` |

The tile keeps the mark legible on both light and dark pages. `favicon.svg` uses a larger star and a solid tile than `logo.svg` so it survives 16 px.

## Rebuild

```bash
./build.sh
```

The script uses `rsvg-convert` for the PNGs and Python's standard library to assemble the ICO. The Python step has no third-party dependencies.
