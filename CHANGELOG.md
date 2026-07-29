# Changelog

All notable changes to Roundfix are documented in this file.

## [0.0.2] - 2026-07-29

### Added

- Added Adapter Readiness for the Claude ACP Runtime: official
  `@agentclientprotocol/claude-agent-acp` at the pinned minimum, recognized
  legacy lineages that fail with the official install action, and lineage
  proof that never accepts a matching executable name. The Doctor Command
  reports adapter evidence for every runtime the required profiles reference,
  and the Setup Command migrates a stale Claude override the way it already
  migrated Codex.
- Added `make baseline-digests`, one sanctioned command that regenerates every
  derived Baseline digest from its canonical source across both chains — the
  Skill-digest snapshots and the module chain covering the Source Baseline
  manifest, formatter goldens, and catalog fixtures. It names each regenerated
  artifact, reports an explicit no-change result, and describes failures with a
  stable error code, the failing stage, retryability, and a next step.
- Added a repository-green precondition: a Task whose Verification is the
  repository-wide gate now proves the repository is green before any Agent
  Session is created, and settles with a reason distinct from a post-Agent
  Verification failure.
- Added typed Review Source Evidence, the Review Skipped outcome, bounded
  transient retry, unknown Review Issue counts, current-head approval
  evidence, artifact-only evidence inheritance, notification receipts, and the
  Supervisor monitor command in the Detached startup report.
- Added the Reconcile Command for terminal Run Worktrees and Run Branches,
  with proof-based classification and `--apply` as its only mutation switch.
- Added compare-and-set terminal completion, owner-identity proof for Force
  Stop, and registered Agent Session cleanup.
- Added Repository Skill Set readiness to the Doctor Command as a blocking
  check with per-ownership remediation.
- Added the autonomous Spec delivery contract to the Roundfix Skill: the loop
  from branch to merged Pull Request, the rule that QA is requested once when
  the Task Graph closes, and the conditions that stop the loop.

### Changed

- Advertised Agent Model identifiers are opaque when the adapter advertises an
  independent reasoning control, so a bracketed identifier is selectable as
  printed and a context-window annotation is never accepted as a reasoning
  effort. The built-in `frontend` profile moves to the proven
  `claude / opus / xhigh`, and the pinned minimum Codex adapter version rises.
- Task decomposition is decided autonomously from a PRD and TechSpec that
  passed their own gates, escalating only when authority, architecture, or
  blast radius is at stake.
- The repository gate behaves identically inside and outside a sandbox: the Go
  build cache defaults to a repository-local ignored path when unset, and an
  exported value still wins.
- Preflight refusals name the fallback boundary, so an operator learns that
  Fallback Chains activate only after Run creation instead of guessing why a
  configured fallback did not apply.
- Refreshed the upstream-managed skill set.

### Fixed

- Fixed the Daemon staging path so an executable file can no longer reach a
  Task commit, and stopped a bare `go build` from leaving an untracked binary.
- Fixed same-day QA report selection, which ordered names as raw strings and
  therefore returned the day's first report instead of its newest, so a Spec
  that passed on a rerun kept reporting the earlier failure.
- Fixed Settle so it never rewrites a terminal outcome.

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

[0.0.2]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.2
[0.0.1]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.1
