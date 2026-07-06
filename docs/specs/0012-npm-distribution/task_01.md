---
task: task_01
spec: 0012-npm-distribution
status: pending
type: infra
complexity: medium
---

# Task 01: Platform mapping table and per-platform binary package scaffolding

## Overview

Lay the packaging foundation: one authoritative mapping table that ties each
supported target's npm platform token, Go `GOOS`/`GOARCH`, npm `os`/`cpu`
fields, package name, and release-asset name together, plus the per-platform
`@roundfix/cli-<os>-<arch>` package scaffolding it drives. This is the single
source later tasks and the workflow read from, so the npm-vs-Go token split
lives in exactly one place.

## Requirements

1. MUST define one mapping table covering darwin-arm64, darwin-x64, linux-arm64,
   linux-x64, and win32-x64, each row carrying: npm package name, npm `os`/`cpu`
   values, Go `GOOS`/`GOARCH`, and the release-asset name (which uses the Go
   `GOOS`/`GOARCH` tokens, e.g. `windows`/`amd64`, not the npm `win32`).
2. MUST scaffold each `@roundfix/cli-<os>-<arch>` package with a `package.json`
   setting `os`/`cpu` and a defined slot for the prebuilt binary.
3. MUST keep the packaging tree separate from Go sources (e.g. under
   `dist/npm/`) and MUST NOT alter any CLI behavior.
4. MUST make the mapping consumable by both the launcher shim (npm tokens) and
   the workflow/asset test (Go tokens) without duplicating the token split.

## Subtasks

- [ ] Mapping table for the five targets (npm ↔ Go ↔ asset names)
- [ ] Per-platform `package.json` scaffolds with `os`/`cpu`
- [ ] Binary slot/layout each package expects
- [ ] A `npm pack` dry-run (or equivalent) confirms file inclusion per package

## Acceptance Criteria

- [ ] The mapping table lists all five targets with correct npm `os`/`cpu`, Go `GOOS`/`GOARCH`, and asset-name tokens (win32 package ↔ windows asset).
- [ ] Each per-platform `package.json` declares `os`/`cpu` matching its target.
- [ ] `npm pack` (dry run) on a per-platform package includes the binary slot and nothing extraneous.
- [ ] No Go source or CLI behavior changed.

## Verification

- `cd dist/npm/<a-platform-package> && npm pack --dry-run` — expected: lists the intended files only.
- `rtk go build -buildvcs=false ./...` — expected: builds unchanged (no Go changes in this task). The `-buildvcs=false` flag matches the repo Makefile and avoids VCS-stamping failures in sandboxed Task Worktrees.

## References

`_prd.md` → User Stories 1-2; Core Feature 2. `_techspec.md` → Platform tokens —
npm vs Go, Build Order 1. ADR-0031.
