# Agent instructions

This repository is Roundfix: a local-first Go CLI and daemon that picks up work
items — PR review issues and Spec Task Graphs — resolves them through the user's
selected ACP runtime, and pushes only when nothing unresolved remains. Stdlib
`flag` dispatch and a Bubble Tea v2 TUI.

Project map: `cmd/roundfix/` is the thin CLI entry point; behavior lives in
`internal/...` — `internal/cli/` owns parsing, output, and exit behavior;
`internal/app/` holds app metadata.

## Repository rules

- Keep the project KISS: prefer the smallest behavior that satisfies the
  documented product contract.
- **NEVER** copy names, branding, package names, comments, examples, or
  generated artifacts from reference projects into this repository.
- **HARD RULE — roundfix skill sync**: before opening any PR, confirm the
  roundfix skill still matches the shipped CLI behavior; a PR that changes CLI
  behavior ships the skill update too.
- **HARD RULE — skill ownership**: repo-owned authorial workflow skills may be
  adapted locally; every other skill is upstream-managed and **MUST NOT** be
  modified here.
- **HARD RULE — sanctioned digest regeneration**: after an expressly authorized
  Roundfix-owned Skill or Baseline module edit, run `make baseline-digests`.
  Every derived pin rewritten by that command is deterministic fallout of the
  authorized source edit and needs no separate express authorization. A
  hand-edited pin value remains an unauthorized mutation.
- The release plan the baseline requires is `roundfix release plan`, and the
  maintainer decisions it defers to live in `docs/user-guide/release-runbook.md`.

## Anti-patterns (immediate rejection)

1. Introducing Cobra, testify, or any dependency the stdlib covers.
