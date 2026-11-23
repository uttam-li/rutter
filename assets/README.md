# Brand assets

| File | Use |
| --- | --- |
| `logo.svg` | Primary mark, 512×512 dark tile |
| `logo-512.png` | Raster version for READMEs and previews |
| `favicon.svg` | Scalable favicon |
| `favicon.ico` | Legacy favicon with 16, 32, and 48 px images |
| `favicon-16.png`, `favicon-32.png` | PNG favicon fallbacks |
| `apple-touch-icon.png` | 180×180 iOS home-screen icon |

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
| Route | `#c4ffdf` → `#56dda0` |
| Node | `#5ea8ff` |

The tile keeps the mark legible on both light and dark pages. Regenerate the PNG and ICO files with `rsvg-convert` after editing the SVG source.
