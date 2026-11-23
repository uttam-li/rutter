# Global agent rules

## Step 0

Before any other tool, unless it is already loaded, call `skill({ name: "unslop" })`.

When the task matches, load the related skill:

- UI/UX: call `skill({ name: "modern-web-guidance" })` and `skill({ name: "ui-unslop" })`.
- TypeScript: call `skill({ name: "typescript-best-practices" })`.
- Browser: call `skill({ name: "browser" })`; use `playwright-cli` only as a last resort.

## Version control

Do not push to origin without approval. Create branches and commits locally only.

## Wayfinder

Put wayfinder maps and tickets in `.wayfinder/` at the project root unless told otherwise. If a map already exists, create a new one alongside it instead of extending it.

## Package managers

Use Bun by default. If Bun does not work, use pnpm; consider npm only if neither works.
