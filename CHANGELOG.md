# Changelog

All notable changes to Roundfix are documented in this file.

## [0.7.0] - 2026-08-26

A Run whose verified work meets a refusing commit hook no longer dies losing
it, and a gate that refuses at its own precondition no longer poisons every
later run of the same Spec.

### Added

- **A hook refusal is a classified, recoverable state.** When the repository's
  commit hook refuses after the authoritative Verification passed, the Daemon
  classifies it `hook_refused` — by the hook's own output when it names itself,
  and by a structural signal from the resolved hooks directory when it prints
  only its finding, which is the husky/lint-staged shape that used to die as a
  bare git error. The Task stays `completed`, the work stays staged, and the
  Run names the recovery command. The invariant that a commit hook must never
  be stricter than the authoritative Verification ships in the Baseline and
  renders into every adopting repository. Measured origin: three Runs of one
  Spec died this way in a fleet repository, each with correct, verified work
  recovered by hand.

- **Settle recovers completed-but-uncommitted work.** `settle` accepts a Task
  that is `completed` while a surface still holds its work uncommitted — the
  exact state a hook refusal leaves. It resolves the surface in priority order,
  re-runs the Task's Verification verbatim in place with no Agent turn, stages
  the whole surface including deletions and renames, commits with the standard
  Task message, and integrates. A failing re-run stops with the surface
  untouched.

- **A refused gate writes its refusal once.** A QA gate stopped at a
  precondition writes a report its own contract accepts: one terminal
  `blocked | precondition` row carrying the refusing check and its reason,
  instead of an empty Results table that blocked every later run until someone
  deleted evidence. Wired end to end — the Daemon's mechanical stage and the
  qa-gate skill — and proven by running the stage over the refusal writer's own
  bytes. Measured origin: two Specs deadlocked this way on 2026-08-12/14, and a
  third reproduced it live during this release's own QA.

- **The mechanical stage reads the newest report only, and only its Results
  table.** A superseded refusal cannot block the run that supersedes it, and a
  comparison table a gate writes in its prose as evidence is no longer parsed
  as result rows — three prose cells each became a separate blocker in the
  measured case.

- **Archival lifecycles are Baseline rules.** The active documentation
  directories hold live work only: findings that are done or deferred archive
  to `docs/history/findings/` through Rollup absorption with an `absorbed_by`
  license; backlog entries in any status other than `open` or `promoted`
  archive with no license; handoffs archive only on the user's explicit
  confirmation, and all together. A Rollup that licenses archived members must
  itself stay active — moving one unresolved 261 licenses in the measured
  sweep.

### Changed

- **The ADR set is consolidated.** Twelve new records (ADR-0138 through
  ADR-0149) each absorb two or three same-subject decisions; the 25 superseded
  members carry `superseded_by` and live under `docs/history/adr/`, and every
  external citation follows its consolidation. 133 active records became 120.

- **One gap-promotion rule.** The Spec Consistency Check's gap promotion is
  exported and shared, so the Daemon's gate precondition and
  `roundfix spec check --strict` cannot drift apart — the refusal names that
  command, so they must agree.

## [0.6.0] - 2026-08-14

Retired documentation gets one home, and a repository still on the old layout
reaches the new one by running Roundfix rather than by moving files by hand.

### Added

- **One history root.** `docs/history/` holds every retired family — specs,
  findings, decision records, declined intent, and Review Artifacts. Before
  this, retired material sat inside the active directories it was retired from,
  one archive folder per tree, so history and live work shared a parent and
  excluding history meant knowing each tree separately. The resolver stays the
  single place that knows where a family lives, so no consumer carries a path of
  its own. See ADR-0120.

- **The Baseline migrates a repository it finds on an older layout.** Running it
  detects retired material wherever a previous layout put it and moves it into
  the new root. Relocations enter the Plan as an ordered ledger of source,
  destination, and content identity, carrying no bytes — a plan's size tracks
  the change it makes rather than the history it moves. They apply inside the
  existing transaction, verify each destination against its recorded identity,
  and roll back completely. See ADR-0121.

  Proven against a repository this work did not build: a disposable clone of an
  adopted repository with 413 archived files and 12.4 MiB migrated with its
  bytes intact and its second run reporting nothing left to do.

- **Review Artifacts become a retired family.** A review whose Spec is known
  already travelled with that Spec; the orphan case — a pull request with no
  owning Spec — now has a terminal home. Which orphan is finished is decided by
  local Git reachability, never by the hosting provider, because the migration
  must run offline. A head the repository cannot decide stays live: keeping a
  finished artifact costs a directory, while retiring a live one breaks the
  Round still running. See ADR-0123.

### Fixed

- **A collision no longer passes silently.** A relocation whose destination is
  occupied is refused by name, its siblings still move, the run does not report
  success, and the next run stops calling a repository with a pending migration
  current.

- **A retired decision is not a pending one.** The layout guide treats
  `proposed`, `rejected`, `deprecated` and `superseded` alike as inactive, which
  is right for whether a decision is in force and wrong as a retirement rule.
  Only the last three retire a record; a proposal stays where an author can find
  it. See ADR-0122.

- **Review never reads history.** The exclusion is path-anchored, which also
  closes a hole the previous prefix match never covered: a review directory
  inside an active Spec was not excluded at all.

## [0.4.0] - 2026-08-06

This release closes a gate that failed open, and then makes the loop that
depends on it able to run unattended.

### Fixed

- **A green check is no longer read as a review.** When CodeRabbit declines to
  review — a rate limit, a path-filter skip — it publishes a check whose
  conclusion is **success**, deliberately, so the block does not prevent a
  merge. Roundfix read that check, fell through to the settled path, and
  recorded the head as `verified`. Pull Request #107 merged 125 files and
  +5,621 lines with no review at all because of it.

  `verified` now requires positive evidence that a review ran. A refusal
  resolves `skipped` with the source's own reason, read from the authoritative
  comment rather than the conclusion, and recognised by class so a shape nobody
  has seen yet cannot approve a merge. Anything unrecognised resolves
  `pending`. The default is inverted rather than the refusal list extended: a
  list fixes today's shape and leaves tomorrow's open.

- **The Run Database waits 30s for a held lock**, up from 5s. The database is
  machine-wide, so concurrent Runs, a following event stream, and a `runs list`
  all contend on it. At 5s an Agent Batch died twice writing its completion
  event — an infrastructure failure with no work defect that failed the Task
  and cost a recovery Run each time.

- **A Review Run validates its target before acting.** `watch --pr N` named the
  Pull Request and resolved the branch from the checkout, and nothing checked
  that the two agreed. Preflight now compares them and refuses with both
  branches and both revisions named, creating no Run. A checkout that moves
  mid-Run reaches a terminal interruption distinct from `Failed`, with its
  Review Issues left unsettled — because a re-run is right for an environmental
  stop and wrong for a genuine failure.

### Added

- **Roundfix asks for the review.** With automatic review disabled, nothing
  started one and `watch --until-clean` waited for evidence nobody requested.
  After a Round's Final Push, `watch` and `resolve` now publish one review
  request for the pushed head, idempotent per head SHA through a marker in the
  comment body. `fetch` publishes none. Preflight refuses a configuration pair
  that would strand a Run or buy two reviews per Round.

- **`roundfix gc compact`** previews bytes before, reclaimable, and projected
  after, then compacts behind three guards: an Active Run, another writer, and
  insufficient temporary capacity each refuse by name. Compaction is explicit
  and never a side effect of a retention sweep. Retention deleted rows and
  nothing returned the pages; the database had reached 1.4 GB with its journal
  already pruned.

- **`roundfix gc sanitize`** discovers every Artifact Root from durable Run
  metadata rather than the current repository alone, classifies each one, and
  removes only directories proven eligible or absent. Review Artifacts and
  every ambiguous path are preserved with the reason.

- **`roundfix storage report`** is read-only, takes no flags, and groups
  measured bytes and row counts by repository, state, table, and Artifact Root.

- **Four Spec Consistency Check rules** refuse at authoring time what used to
  be advice: a Verification that cannot fail when no work was done, mutually
  unsatisfiable requirements, a rehearsal Task that declares no cases, and a
  loop order that diverges between the places stating it.

- **`docs/backlog/`** gives typed intent a home — work that is not yet a Spec
  and not a dated finding, with a frontmatter contract and a stated boundary
  against findings.

### Changed

- **Skill compatibility is a declared minimum, not a content digest.** Editing
  an owned skill no longer changes what the binary claims about it, and no
  regeneration step follows. A version at or above the minimum satisfies; below
  blocks naming the skill, the minimum, the version found, and the upgrade
  path; unversioned or unresolvable stays its own state and an unreachable
  source is never reported as a missing skill. Third-party skills are never
  held to a version Roundfix invented for them. This resolves the known issue
  recorded in 0.3.1.

- **The loop states one order, once.** Implement the graph including its
  authored gate, archive, open the Pull Request, watch until Clean, merge. The
  case that motivated reordering — a Spec whose acceptance observes its own
  Pull Request — is absorbed by environment-blocked rows carrying equivalent
  evidence.

## [0.3.1] - 2026-08-03

Version 0.3.1 supersedes the 0.0.3 tag, which never published: its release
failed on the npm dist-tag and no coordinate reached the registry. The
platform packages already carried 0.3.0 from an earlier versioning era, so
0.3.1 is the first version above every published coordinate and the first
consistent release across the launcher and all five platform packages.
Everything the unpublished 0.0.3 section described ships here for the first
time, together with the work that followed it.

### Known issues

- Skill content is pinned by digest into the binary, so editing a shipped
  skill changes what Roundfix claims about itself. Compatibility becomes a
  per-skill declared minimum version in a later release; until then, keep
  skill edits and binary releases together.

### Breaking

- **Removed the `--qa` parameter from the Implement Command.** The QA gate is
  now authored into a Spec's Task Graph as a terminal `qa` Task, declared at
  decomposition with `qa: task_NN` in the manifest — or declined once with
  `qa: declined` and a reason. What runs is what the graph says. Passing
  `--qa` is an unknown-flag error whose message names the contract. Task
  Graphs authored before this contract keep working unchanged, and appending
  work after a gate has reported now invalidates that result at the next
  load, naming the inserted Tasks instead of quietly starting a second cycle.
  Interactive Input no longer asks about QA.

### Performance

- Cut the implementation loop's verification cost by 97%: an unchanged tree
  now completes the local gate in about 5 seconds against 142. The change was
  structural — code under test takes its environment and working directory as
  arguments instead of mutating the process, which both unblocked parallel
  execution (1 → 479 declarations in the CLI package) and made Go's test
  result cache trustworthy.
- Cut the full fresh suite by roughly a third, and removed work rather than
  coverage: the CLI package compiled the whole project seven times per run
  from empty caches, the coverage harness invoked the toolchain once per
  package instead of once, and eleven packages declared no parallelism at
  all. No test was deleted, skipped, or weakened.
- Reduced git subprocesses in production paths by roughly 13%: object reads
  that spawned one process per file now stream through a single batch
  process, and repository resolution asks its several questions in one
  invocation. Every Run pays these against the user's own repository.

### Added

- Added a Pull Request verification workflow. Until now no workflow ran the
  test suite before a merge — only the release workflow, on a tag push. CI
  forces a fresh run so its verdict never rests on a cache entry from another
  commit, while the local gate lets the cache decide what to re-run.
- Added npm Trusted Publishing: the release workflow authenticates publication
  through GitHub Actions OIDC instead of a long-lived repository token, on a
  Node 24 runtime with an asserted npm floor. A bounded per-coordinate token
  fallback keeps the release set whole while trusted-publisher configuration is
  still unproven, records every coordinate that needed it, and closes when one
  release publishes all six packages without it.
- Added a publication preflight that evaluates the launcher and every platform
  package as one release set after Verification and before cross-compilation.
  It refuses an already-used version, detects a post-unpublish cooldown, names
  the exact blocked coordinate, and reports a registry it could not read as
  undetermined rather than as ineligible — so no package publishes unless every
  coordinate is eligible.
- Added a publish-free rehearsal of that preflight, so the release set can be
  checked against the live registry without cutting a tag.
- Added a read-only capability re-check that resolves no decisions and writes
  nothing, so remediating a blocking Baseline divergence and re-checking is a
  loop rather than a full re-plan.
- Added a fourth outcome to the Baseline divergence prompt — exit without
  writing, print per-divergence remediation, name the re-check command —
  journaled distinctly from a decline, so pausing for repository work is no
  longer recorded as refusing the Baseline.
- Added a five-axis result matrix to Baseline runs: approved postimages,
  semantic retention, profile alignment, repository Verification, and
  idempotence, each reported as verified or not run. Completion language now
  requires verified retention and a passing idempotence check.
- Added an explicit removal declaration to `profiles configure`, so deleting a
  configured category is always intentional.
- Added the autonomous-loop discipline to the Context-Driven Baseline, so an
  adopting repository receives it from `roundfix baseline` rather than from a
  hand-copied block.

### Changed

- `profiles configure` now merges by Agent Work Category instead of replacing
  the whole `profiles` map. Configuring one category leaves every other
  configured category byte-identical, including its Fallback Chain, comments,
  key order, and indentation. A per-category summary — added, replaced,
  removed — renders before any write, so a removal can no longer hide inside
  reformatting churn.
- Executable capability discovery now resolves a bounded symlink chain and
  judges the target without ever executing it. Tools installed through Homebrew
  or Docker Desktop stopped reporting as missing, and a cycle, a broken link,
  and a non-executable target each produce their own diagnostic instead of a
  bare absence.
- Every unsatisfied Baseline divergence now renders the probe it evaluated —
  inspected paths with their states, or the inspected PATH candidate — in text
  as well as machine output, so remediation no longer requires reading catalog
  assets or source. Divergences group by requirement strength, and every
  advisory states that it does not block readiness or apply.
- A Profile or catalog digest change under an unchanged Baseline identifier is
  now an upgrade rather than a refresh: every previous managed Normative Clause
  receives an explicit disposition, or planning exits action-required. A ready
  update plan can no longer carry an empty retention ledger.
- Baseline carrier classification stops warning about the managed artifacts an
  apply just wrote. Only unmanaged nested carriers warn, and a carrier that
  cannot be classified keeps its warning.
- The QA gate is now reachable for a Spec whose acceptance observes its own
  Pull Request, and its verdict semantics distinguish rows blocked by
  environment from rows blocked by a finding.
- A Spec now owns the sources it adopts: `write-prd` moves adopted inbox notes
  and findings into the Spec's own `references/`, and archiving verifies the
  Spec is self-contained.
- Force Stop reads owner identity from the kernel instead of forking `ps`, so
  it survives a host that cannot fork, and an unreadable identity is reported
  distinctly from a proven mismatch.
- The release workflow's GitHub Actions are pinned to the current majors, which
  run on the Node 24 runtime GitHub now requires.

### Fixed

- Fixed `make baseline-digests` refusing the very edit it exists to regenerate.
  A module edit changed a generated guide, which invalidated the derived pin,
  which refused the catalog load, which prevented the refresh — so the
  command's own remediation was the command itself. Regeneration now defers the
  pins it rewrites and re-validates strictly afterwards; every other load stays
  strict.
- Fixed a declined `profiles configure` in a non-interactive context exiting
  zero, so automation can no longer read a refusal as a successful write.
- Fixed a fragment deleting configured profiles by omission. Removal is now a
  declared operation and can never be a side effect of an incomplete file.
- Fixed removal of a profile category disturbing unrelated blank lines and
  comments around its neighbours.
- Fixed the required Repository Skill Set being a fixed list, so a repository
  is no longer asked for skills its stack has no reason to hold.
- Fixed the autonomous loop requesting the QA gate before the surfaces its
  acceptance observes exist.
- Fixed the npm dist-tag being applied implicitly, which failed publication
  when a platform package already carried a higher version than the one being
  released.
- Fixed a Verification command starting after its Run's context was cancelled:
  capacity could be granted in the same instant the cancellation landed, so a
  Run that was tearing down could still spawn one more child process.
- Fixed the detached-child code path mutating process-global state while its
  test ran in parallel.

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

[0.7.0]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.7.0
[0.6.0]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.6.0
[0.4.0]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.4.0
[0.3.1]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.3.1
[0.0.3]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.3
[0.0.2]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.2
[0.0.1]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.0.1
