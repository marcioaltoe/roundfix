---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Roundfix ships through npm platform packages

Releases are published to npm as a launcher package (`roundfix`) with per-platform binary packages under `optionalDependencies`, built by a tag-triggered workflow that cross-compiles the Go binary for darwin/linux/windows on both architectures — so `npx roundfix`, `bunx roundfix`, and global installs work from one channel on every platform. The same workflow uploads the binaries as GitHub Release assets, keeping the Upgrade Command's existing channel fed. Homebrew was deferred: it excludes Windows and demands a maintained tap for no reach npm does not already provide, and Node is already a hard prerequisite of the agent layer, so the npm runtime adds nothing new.
