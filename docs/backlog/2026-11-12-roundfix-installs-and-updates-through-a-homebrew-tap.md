---
type: feat
status: open
created: 2026-11-12
spec: null
reason: null
---

# Roundfix installs and updates through a Homebrew tap

## Opportunity

Recorded at the maintainer's request, relayed by jarvis in ai-memory message `01a1d27a-e25a-7180-86e3-69dc7cfa3ff3` of 2026-11-12 (22:35 BRT); the idea had been on the maintainer's Todoist since 2026-10-20. A Roundfix-owned Homebrew tap with a formula would let a macOS or Linux user install with `brew install` and update with `brew upgrade`, as OpenCode 2 does with its own tap.

## Value

The maintainer's machines update their whole toolchain through Homebrew and `dev-upgrade`; a formula would bring Roundfix into that routine instead of the separate npm install.

## Shape

A tap repository with a formula that downloads the release asset built by the existing release workflow for the platform, updated by the release workflow on each tag. Non-binding; no Spec yet.
