# Changelog

All notable changes to Roundfix are documented in this file.

## [0.0.1] - 2026-07-25

### Added

- Added the public Context-Driven Baseline workflow for greenfield setup,
  adoption, updates, and Readoption. Humans use one confirmation-gated
  `roundfix baseline` flow; automations use portable `baseline plan` and
  `baseline apply` documents bound to repository preimages and a Plan Digest.
- Added built-in and repository-owned Baseline Profiles, Profile drafts,
  guided Profile adaptation, Repository Capability alignment, and
  `--profile-file` support.
- Added byte-exhaustive instruction segmentation and supervised semantic
  classification so preserved rules can move to their owning agent guides
  without losing source accounting.
- Added generated Instruction Hierarchy, ADR lifecycle, Findings format,
  repository-extension, domain-language, autonomous-work, Secondbrain, and
  local Spec workflow guidance.
- Added typed Project Decisions and Project Constraints for UUID v7
  identifiers, repository-owned HTTP contracts, Better Auth route exceptions,
  and explicit authorization before changing lint, formatting, test-runner, or
  verification configuration.
- Added Agent Selection Profiles with Preferred Selections, ordered Fallback
  Chains, exact ACP capability proof, and the `profiles show`,
  `profiles configure`, and `profiles validate` commands.
- Added read-only Release Plan output for normal semantic-version planning and
  the confirmation-gated `release plan --reset-to v0.0.1` inventory.

### Changed

- Restarted the Roundfix-owned CLI, npm packages, setup generation, Release
  Plan schema, changelog, and distributed skills at version `0.0.1`.
- Kept Build Commit and Build Time as the source-state identity for local
  binaries without changing their semantic version.
- Made the Go CLI the authority for Context-Driven Baseline planning,
  application, recovery, Profile management, skill restoration, and canonical
  asset synchronization; `setup-context-driven` now provides recipes over that
  public API.
- Routed Task, QA, and review work through owned Agent Sessions selected by
  category profiles, with selection attempts and fallback events visible in
  persisted Runs, text output, Attach, and the Live Run View.
- Preferred Claude Opus 5 with `xhigh` reasoning for design, UI, UX, TUI, and
  frontend work, with Codex as the configured fallback.
- Accepted `acpx` `0.12.0` and newer compatible versions without downgrading a
  newer installation.
- Suppressed byte-identical repeated tool summaries in non-interactive console
  output while preserving every Run Event.
- Allowed repository-owned HTTP Contract exceptions to use every standard HTTP
  method instead of restricting all exceptions to `GET` and `POST`.

### Fixed

- Made cooperative Agent Session cancellation and forced-close ordering
  deterministic without changing production grace periods.
- Rejected partial Agent Selection overrides and failed closed when an exact
  runtime, model, and reasoning tuple cannot be proven before Run creation.
- Connected semantic classification to the public interactive Baseline flow,
  including repository-document, repository-rules, managed-entry, and rejected
  dispositions.
- Preserved cumulative Profile removals across adaptation retries and rejected
  drafts that retain Project Decisions after removing their render modules.
- Reused complete persisted Better Auth exceptions, including customized
  scopes and method sets, during Baseline upgrades.
- Detected machine-specific absolute paths in Markdown links and `file://`
  references during source-corpus validation.
- Removed the fixed 256-source ceiling from HTTP route discovery while keeping
  bounded candidate output.
- Kept confirmed Project Decisions and repository rules visible through
  generated guidance, rollback, stale-plan refusal, and Profile changes.

Earlier release sections are intentionally omitted from the restarted
changelog. Git history remains the source for prior implementation history.

[0.0.1]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.1
