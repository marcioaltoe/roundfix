# Changelog

All notable changes to Roundfix are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/), versions follow
[Semantic Versioning](https://semver.org/), and commit subjects use
[Conventional Commits](https://www.conventionalcommits.org/) (validated by
Cocogitto). The tag drives the release: pushing `vX.Y.Z` runs the release
workflow, which publishes the npm packages and creates the GitHub Release from
this file's matching section.

## [0.2.0] - 2026-07-07

Run discovery, a redesigned cockpit, first-class knowledge-workspace support,
and outcome notifications — driven by two days of dogfooding Roundfix on
itself and on external repositories.

### Added

- **Run Browser** — `roundfix runs` and no-argument `roundfix attach` open a
  navigable TUI over the repository's Runs: Active Runs by default, `a`
  toggles the full history, `Enter` opens the read-only Live Run View, and
  closing it returns to a refreshed browser. Cancel is side-effect free.
- **`runs list`** — deterministic Run discovery for agents and scripts: eight
  stable columns (run id, state, kind, target, agent, UTC start, duration,
  local branch), an Active-only default bounded to the 20 newest,
  `--state <active|terminal|all>`, `--limit N` (`0` unbounded), `--all` for
  every repository, and exactly one trailing stderr note naming hidden Runs
  and the widening flag.
- **External Spec Root** — `specs.root` points Spec artifacts at any
  directory, including a knowledge-workspace repository nested behind a
  `docs/specs` symlink. Resolved once against the user's checkout and carried
  into Run and Task Worktrees; external task files and QA Reports are never
  staged into code-repository commits (settled on disk, committed by their
  owning repository), and staging unconditionally skips paths that cross a
  symbolic link instead of failing the Run.
- **Run Outcome Notifications** — operational Runs (`resolve`, `watch`,
  `implement`) fire one best-effort notification at their terminal outcome:
  native desktop by default (macOS `osascript`, Linux `notify-send`), or
  `notify.command` with a `ROUNDFIX_RUN_ID`/`ROUNDFIX_OUTCOME`/
  `ROUNDFIX_KIND`/`ROUNDFIX_TARGET` environment contract and a 30s bound.
  Detached Runs notify from the detached process. Failures warn and journal;
  they never change a Run's report, outcome, or exit code.
- **Review artifacts in history** — after a clean integration, `resolve` and
  `watch` commit the Run's review artifacts in one separate Daemon-owned docs
  commit (`docs: review round NNN for pr <n>`) that rides the Final Push;
  `fetch` still never commits, and external artifact roots are never staged
  (ADR-0036, superseding ADR-0029's never-commit clause).
- **Skill bundle** — the binary ships all 14 Roundfix-owned skills (the
  operational `roundfix` skill plus the 13 authorial workflow skills);
  `skills list` separates the bundled set from the recommended external
  skills, and `skills install` writes the full bundle.
- **Docs** — operational usage guide (`docs/usage.md`), CONTEXT-driven
  development write-up with attribution
  (`docs/context-driven-development.md`), release runbook
  (`docs/release-runbook.md`), and an npm-first README Install section.

### Changed

- **Cockpit visual fidelity** — the Live Run View now renders to the approved
  design: one style-token source (cyan section labels and active borders,
  green done, amber running or waiting, red locked or failed, muted gray
  timestamps and paths), carded Work Queue with an accent-bordered selection
  and Batch elapsed stamps, a timeline with automatic ▼/▶ Batch collapse, an
  aligned timestamp gutter, one bounded summary line per event (raw payloads
  live only in the Detail Modal), a `Live · detail hidden/open` indicator, a
  styled Detail Modal, and explanatory per-Run-kind empty states.
  `ROUNDFIX_COLOR=never`/`NO_COLOR` keeps every distinction through text
  markers. Every fidelity behavior is byte-pinned by tests.
- Final Push for review Runs now runs from the user checkout after
  integration, so the review-artifacts docs commit rides it.

### Fixed

- **Symlinked spec trees no longer kill Runs** — the per-Task commit died on
  `git add` when the task file's path crossed a symbolic link (field reports:
  conexus, tax-poc), after verification had already passed. The commit
  boundary now degrades: unstageable paths are dropped with a journaled
  warning, and a Task whose only changes are external settles completed
  without a commit.
- **Cleanup never demotes a Clean Run** — Roundfix-owned worktree removal uses
  `git worktree remove --force` (bootstrap debris no longer blocks it), and a
  cleanup failure after successful integration warns and journals instead of
  converting the outcome to Failed.
- **Claude runtime under a Claude orchestrator** — the acpx child environment
  strips Claude Code's nested-session guard (`CLAUDECODE`), so a
  Claude-driven session can launch claude-runtime Runs.
- `attach` with an unknown value (such as a picker number) now explains that
  picker numbers are not stable Run ids and points at the Run Browser and
  `runs list`.

### Removed

- The unreleased `--active` flag on `runs list`, superseded by `--state`
  before any release carried it.

### Notes

- Repositories that used interim git shims to survive the symlink commit bug
  must remove them after upgrading — they mask the real commit boundary. Then
  set `specs.root` to the knowledge-workspace path in `.roundfixrc.yml`.

## [0.1.0] - 2026-07-06

First public release. Roundfix drives two kinds of work to a clean, verified,
integrated outcome with local coding agents.

### Added

- **Review resolution** — `fetch`, `resolve`, and `watch` pull unresolved
  CodeRabbit review threads into local Round and Review Issue artifacts, drive a
  selected ACP Runtime (`codex`, `claude`, `opencode`) over bounded Batches,
  verify, commit per Batch, resolve source threads, and push only when nothing
  unresolved remains (`watch --until-clean`).
- **Spec execution** — `implement` runs a Spec's Task Graph in isolated Run and
  Task Worktrees with a ready-set wave scheduler (`worktree.concurrency`), one
  commit per verified Task, an optional `--qa` gate, and porcelain-only
  integration onto the user's branch.
- **Detached Runs** — `--detach` on `resolve`, `watch`, and `implement`
  re-executes Roundfix as a session leader that survives its caller; followed
  with `attach`, controlled with `stop`.
- **`doctor`** — read-only machine-readiness diagnosis (Node.js, pinned acpx,
  the configured Agent probe, and codex runtime hygiene) that mutates nothing.
- **Support commands** — `setup`, `init`, `upgrade`, `settle`, `archive`, `gc`,
  `stop`, and `skills` (`check`/`install`).
- **Worktree bootstrap** — `worktree.bootstrap` prepares each new worktree
  (install dependencies, migrate and seed databases, warm caches) after
  `worktree.copy` and before Agent work, so stateful monorepos verify cleanly.
- **Run store retention** — `store.journal_retention` plus `roundfix gc` bound
  the Run Event Journal and reclaim run-artifact storage; Active Runs, `runs`
  rows, and locks are never pruned.
- **Review artifact placement** — resolved into the repository's spec tree by
  default: an explicit Artifact Directory wins, else a spec-associated PR writes
  under `docs/specs/<slug>/reviews/`, else `docs/specs/_reviews/pr-<n>/`.
- **npm distribution** — `npx roundfix`, `bunx roundfix`, and global installs
  from one channel via a launcher package plus per-platform binary packages,
  built and published by a tag-triggered release workflow that also uploads
  GitHub Release assets for `roundfix upgrade`.

### Notes

- Requires Node.js 22.13+, acpx 0.12.0, and the GitHub CLI authenticated for the
  target repository.
- codex runtime hygiene gates on the `com.apple.quarantine` attribute (the real
  XProtect trigger) and code-signature validity — not `spctl --assess`, which
  rejects any signed CLI that is not a notarized app.

[0.2.0]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.2.0
[0.1.0]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.1.0
