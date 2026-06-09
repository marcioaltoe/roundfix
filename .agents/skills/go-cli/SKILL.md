---
name: go-cli
description: Use when designing, building, modifying, or reviewing Go CLI behavior in this repository, including command names, flags, interactive prompts, stdout and stderr contracts, exit codes, help text, version output, and command dispatch.
---

# Go CLI

Keep the command surface boring, predictable, and testable.

## Workflow

1. Read `docs/product-brief.md` for command contract before changing behavior.
2. Keep `main` thin. Route through a testable function that accepts args,
   stdout, and stderr.
3. Parse commands before doing IO-heavy work.
4. Validate as much as possible during preflight and fail before remote fetches,
   database writes, or ACP startup.
5. Test stdout, stderr, and exit code separately.

## Command rules

- Keep command names aligned with the product contract: `fetch`, `resolve`,
  and `watch`.
- Use clear exit codes:
  - `0` for success.
  - `1` for runtime failure.
  - `2` for usage, validation, or not-implemented command paths.
- Write user-facing normal output to stdout.
- Write errors, validation failures, and next-action guidance to stderr.
- Do not call `os.Exit` outside `main`.
- Do not print from deep packages. Return data or errors to the CLI layer.
- Prefer standard-library parsing while the command surface is small.
- Do not add a CLI framework until it removes real complexity.

## Interactive prompts

- Suggest detected values when safe, but let the user override them.
- Remember that `fetch`, `resolve`, and `watch` collect different optional
  parameters.
- Treat ACP runtime selection as explicit. If the selected runtime fails, show
  that error and stop.
- Do not fallback to a different runtime automatically.
- Avoid making the user wait before preflight validation has passed.

## Help text

- Keep help text concise and operational.
- Mention only commands and options that exist.
- Prefer exact next commands in failure messages.
- Do not expose internal implementation details in help text.

## Verification

Use these checks before finishing CLI work:

```bash
rtk go test ./...
rtk go run ./cmd/roundfix --help
rtk go run ./cmd/roundfix --version
```
