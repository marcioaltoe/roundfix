# Changelog

All notable changes to Roundfix are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/), versions follow
[Semantic Versioning](https://semver.org/), and commit subjects use
[Conventional Commits](https://www.conventionalcommits.org/) (validated by
Cocogitto). The tag drives the release: pushing `vX.Y.Z` runs the release
workflow, which publishes the npm packages and creates the GitHub Release from
this file's matching section.

## [0.3.0] - 2026-07-15

Roundfix-owned model selection, review-loop integrity, settlement robustness,
and launch/recovery fixes — a week of dogfooding Roundfix on itself (specs
0023–0029), where three of the shipped fixes closed failure modes reproduced
live during the same cycle.

### Added

- **Agent Model ownership** — Roundfix selects the Agent Model and Default
  Reasoning Effort for every Agent Session explicitly; runtime-owned local
  configuration never participates. Per-runtime Project/User Config keys
  (`runtimes.<runtime>.model`, `runtimes.<runtime>.reasoning_effort`), ordered
  Codex and Claude Model Catalogs in Interactive Input, a `--reasoning-effort`
  override next to `--model`, and an availability preflight that fails before
  a Run exists, naming the rejected values and the recovery paths.
- **Model fallback guardrail** — when the selected model fails its probe, the
  preflight offers a probe-discovered Fallback Selection gated on explicit
  confirmation (or a copy-paste re-run command in non-interactive mode); no
  silent substitution, ever.
- **Optional reasoning effort** — an empty reasoning value is valid and means
  the Agent Model manages reasoning; Roundfix assigns the option only when
  configured non-empty.
- **Context-efficient Spec Runs** — the Spec Context Bundle gives each Task a
  bounded start context (the assigned Task, declared `## Context` paths, and
  prior changed files) instead of cold-start repository exploration.
- **Branch Integrity Preflight** — `fetch`, `resolve`, and `watch` refuse to
  start while unintegrated `roundfix/run-*` work or another Active Run is
  bound to the PR Head Branch: fast-forwardable pending work is integrated
  automatically and reported, anything else blocks naming the branch, ahead
  count, and the exact command. `--skip-branch-integrity` bypasses both
  guardrails only after publishing a PR audit comment — a failed publish
  fails the command.
- **Clean Unverified outcome** — after the Final Push, watch polls for the
  Review Source check through a grace window (`watch.check_grace_period`,
  default 5m); a check that never appears ends the Run `CleanUnverified` with
  its own exit code `3` instead of a silently-noted Clean.
- **Outcome Comments** — Review Issue outcomes propagate to GitHub per issue
  at Batch settlement: invalid and duplicated issues get an explanatory
  comment before their thread resolves, failed issues get the failure reason
  and stay open, and still-unresolved issues receive a run-end comment. All
  comments carry an idempotency marker, so retries never duplicate.
- **Terminal reasons everywhere** — Review Issue artifacts persist a
  `terminal_reason`, failed and skipped Tasks print an indented `reason:` line
  in the implement report, and the review report separates this Run's counts
  from the pull request's cumulative counts.
- **Orphaned-lock reclamation** — Runs record their owner process id; a lock
  whose owner is provably dead is reclaimed automatically (Run completed
  Failed with the reason journaled, one stderr warning) instead of blocking
  every relaunch until a manual force stop.
- **Task status synonym normalization** — `done` and hyphen/space variants of
  the canonical statuses normalize instead of voiding a finished Batch; the
  task file is rewritten to canonical form.
- **Settlement transparency** — settle prints one `commit <path>` line per
  committed path and warns when other failed Tasks share the worktree; the
  Daemon warns (Run Event + stderr) when a Task commit contains no change
  outside the Spec Root.
- **Runtime diagnostics** — preflight names a missing ACP adapter binary with
  its install command instead of a raw spawn failure, and `roundfix doctor`
  gains `adapter:` and `model:` check lines — the model line reports the
  effective Agent Model probe and, on failure, the runtime's currently
  advertised models.
- **CONTEXT-driven setup refresh** — the setup skill is now
  `setup-context-driven`: it scaffolds the full `docs/` layout (inbox, ADRs,
  agent guides, design, findings, handoffs, references, specs, user guide)
  and seeds `docs/agents/docs-layout.md` with each folder's job and a dated
  findings template.

### Changed

- **Review Runs execute in the user's checkout** — `fetch`, `resolve`, and
  `watch` no longer create a Run Worktree or Run Branch: review fixes are a
  delta over the published HEAD, batch commits land directly on the PR Head
  Branch, and Integration Pending no longer exists as a review outcome.
  Preflight requires a clean tracked working tree (untracked files allowed);
  after a failed batch, everything dirty is Agent work by construction.
  Worktree isolation remains the contract for spec Runs.
- **Actionable model rejections** — a Batch that dies because the Agent
  Session rejects the Agent Model settles with
  `Agent Model "<model>" not advertised by runtime "<runtime>"; advertised: <list>`
  in artifacts, journal, and report, replacing the opaque
  `agent/protocol error`.
- **Settle surface selection** — settle picks the first surface (Task
  Worktree, Run Worktree, current repository) where the target Task is
  actually `failed`, always prints `Settle surface: <path>`, and a refusal
  enumerates every candidate with the status found there.

### Fixed

- **Silent `--detach` death** — the Detached Run handshake is now two-phase
  (liveness within 10s, Run creation within a 5-minute ceiling), so a slow
  but healthy Preflight Validation no longer gets killed mid-flight with an
  empty stderr relay; every failure branch prints an explicit diagnostic
  naming the phase, exit code or signal, and the child's output or its
  absence.
- **Stale settle surfaces** — settle no longer resolves a kept worktree where
  the target Task never ran (refusing with a misleading `pending` status)
  while the authoritative checkout holds the Task `failed`.
- **Blocked relaunches after a killed Run** — covered by orphaned-lock
  reclamation above; a dead owning process no longer requires
  `roundfix stop --force` before a new Run can start.

### Deprecated

- **`defaults.model`** — ignored with exactly one stderr warning pointing at
  the per-runtime replacement (`runtimes.<runtime>.model`), following the
  removed-config-keys contract.

## [0.2.0] - 2026-07-07

Run discovery, a redesigned cockpit, first-class knowledge-workspace support,
and outcome notifications — driven by two days of dogfooding Roundfix on
itself and on external repositories.

### Added

- **Run Browser** — `roundfix runs` and no-argument `roundfix attach` open a
  machine-wide navigable TUI: every repository's Active Runs by default (with
  a repository column), `a` toggles the full history, `Enter` opens the
  read-only Live Run View, and closing it returns to a refreshed browser.
  Cancel is side-effect free, and no git repository is required.
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
