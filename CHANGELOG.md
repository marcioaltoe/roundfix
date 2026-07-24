# Changelog

All notable changes to Roundfix are documented in this file.

## [0.0.1] - 2026-07-23

### Added

- Local-first Review Runs and Spec Task Graph execution through configured ACP
  Runtimes.
- CONTEXT-driven repository setup with managed profiles, Baseline Readoption,
  Repository Capability evidence, and Roundfix-owned workflow skills.
- Read-only Release Plan support for normal semantic-version planning and the
  confirmation-gated `--reset-to v0.0.1` inventory.

### Changed

- Restarted the Roundfix-owned CLI, npm packages, setup generation, Release
  Plan schema, changelog, and distributed skills at version `0.0.1`.
- Kept Build Commit and Build Time as the source-state identity for local
  binaries without changing their semantic version.

Earlier release sections are intentionally omitted from the restarted
changelog. Git history remains the source for prior implementation history.

[0.0.1]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.1
