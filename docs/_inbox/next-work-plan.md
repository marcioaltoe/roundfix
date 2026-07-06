# Next work plan — organizing the 2026-07-05 requests

Consolidates Marcio's post-merge requests and the open items from three
dogfood findings logs into three themed specs. Branch:
`ma/parallel-and-distribution`.

## Findings audit — what is still open

Everything from `dogfood-findings.md` (round 1), `-2.md`, and `-3.md`
cross-checked against shipped specs 0001–0008.

**Resolved (no action):** round-1 findings 1–11, 13–15, 18–27 (shipped in
0003/0004/0005/0007/0008 or same-day); round-2 findings 2–4; the merge-ready
glossary gap. Round-2 finding 1 (acpx 10 MiB buffer) is upstream-only,
documented, defended by ADR-0020 — nothing more to do locally.

**Open, routed below.** (The three `dogfood-findings*.md` logs were deleted
2026-07-06 after this table absorbed every still-open item; git history keeps
them. `0009-parallel-scheduling` shipped Clean + QA pass — R3-1 done.)

| Finding | What                                                             | Spec             |
| ------- | ---------------------------------------------------------------- | ---------------- |
| R1-17   | Review artifacts belong with the Spec, not a loose root          | 0010             |
| R3-2    | CodeRabbit issue titles are raw table fragments (emoji/markdown) | 0010             |
| R3-3    | merge-readiness `missing` path needs a docs expectation note     | 0010             |
| R3-4    | Status-poll stderr line repeats every interval                   | 0010             |
| R3-5    | Removing a config key hard-failed the user's own config          | 0010 (shipped focus) |
| R3-6    | External kills orphan acpx adapter processes (40 leaked in a day)| 0010 (shipped focus) |
| R1-12   | Prompt-contract drift (templating, work-plan item 5)             | deferred         |
| R1-16   | codex full-access sandbox preset unavailable via acpx            | upstream/observe |
| R3-1    | Force-stopped Runs keep empty worktrees/branches forever         | DONE (0009)      |

**Round-3 robustness findings (from the 0009 dogfood, 2026-07-06) — route
into a robustness spec (0011 or a dedicated 0012):**

- **R3-5 config hard-break**: removing `resolve.concurrent` made every Run
  fail Preflight on any config still carrying it (broke the dogfood machine's
  own config). A removed config key must degrade to a deprecation **warning**
  or auto-migrate, never a hard Preflight failure.
- **R3-6 adapter orphans**: externally-killed Runs leak detached `codex-acp`
  adapter processes (40 accumulated over one day); `stop --force` cancels the
  acpx session but never reaps the OS process tree. Candidates: `stop --force`
  and the preflight sweep reap the Run's adapter tree; a `roundfix doctor`
  lists/kills orphaned adapters.
- **R3-7 `--detach` mode (Marcio, 2026-07-06)**: Runs die with their parent
  process — three Runs were killed mid-flight when the invoking session
  reaped its background tasks; the interim fix is caller-side
  `nohup ... & disown`, which no product should depend on. Add `--detach` to
  the operational commands: roundfix re-executes itself as a detached
  session leader, prints the run id, exits 0; `roundfix attach <run-id>`
  (already shipped) becomes the follow surface and `roundfix stop` the
  control surface. This is the stepping stone to `roundfix serve`
  (work-plan item 3), which remains the real long-term answer.
- **R3-8 codex notarization / XProtect false-positive (Marcio, 2026-07-06)**:
  macOS Tahoe fired "Malware Blocked: 'codex' was not opened because it
  contains malware" — a known OpenAI issue (openai/codex#24246, root cause
  #23649): the codex binary is OpenAI-signed but not Apple-notarized, so when
  an agent loop spawns it frequently, XProtect YARA-scans the ~193 MB binary
  every exec and intermittently blocks it, especially with the Homebrew-Cask
  `com.apple.quarantine` flag. Environment fixed 2026-07-06 (curl install to
  ~/.local/bin, no quarantine; CODEX_PATH pinned in fish+zsh). Product fix
  candidates: `roundfix setup`/`doctor` detects a quarantined or un-notarized
  codex on PATH and prints the curl-reinstall fix; the agent layer sets
  CODEX_PATH to a verified-clean codex when spawning codex-acp. Note: the
  same research validated R3-6 — Zed hit the identical `codex-acp` PPID=1
  orphan leak (zed-industries/zed#56146) and their fix (kill only PPID=1
  orphans) matches 0010's session-close reaping.

- **Supervision lesson (agent-side, not product)**: a Monitor filter must
  cover every terminal signature (`Preflight failed`, `reached ...`), not just
  the happy path — a config-death read as "still running" for hours. And
  `nohup ... & disown` detaches a Run from the harness (surviving external task
  stops) but removes the completion notification, so a comprehensive monitor
  is the only signal.

**Round-4 robustness findings (from the 0011 dogfood, 2026-07-06):**

- **R4-1 preflight sweep hard-fails on failed-Task debris**: the preflight
  sweep and `stop --force` reap terminal-Run worktree debris with
  `git worktree remove` **without `--force`**. A failed Task keeps its Task
  Worktree by design, and that worktree holds the agent's uncommitted edits, so
  `git worktree remove` refuses and the sweep aborts Preflight with exit 2 —
  bricking every future `implement` Run for that repo until a manual
  `git worktree remove --force`. roundfix cannot self-heal it (both reap paths
  share the bug). Fix candidate: force-remove terminal-Run debris the sweep
  owns, or skip-and-warn instead of hard-failing Preflight. Severity: high
  (any single failed Task blocks the repo). Interim: manual force-remove +
  branch delete.
- **R4-2 orphaned adapters survive session-close reaping**: terminal/failed
  Runs still leak PPID=1 `codex-acp` processes whose acpx session is already
  gone ("No named session … for cwd"), so 0010's session-close reaping has no
  session to close and the OS processes leak (9 accumulated from one Run).
  Extends R3-6: reaping must also cover adapters whose session is already
  detached — reap by explicit PID or process-tree, not only by session close.
  Interim: reap by explicit PID when no Run is active.
- **R4-3 concurrent Tasks editing one file conflict at integration**: the
  serialized cherry-pick integration (ADR-0026) fails the second concurrent
  Wave Task when two independent (`needs: []`) Tasks edit the same file. 0011's
  Wave-1 Tasks all touch `internal/cli/cli.go`, so concurrency 2 produced a
  spurious `integration conflict: internal/cli/cli.go` that failed task_02 even
  though its code was correct. Two mitigations, not mutually exclusive: (a)
  decomposition discipline — one Task owns a hot shared file per Wave, or make
  such Tasks depend on each other so they serialize; (b) product — a 3-way
  merge instead of a plain cherry-pick, or auto-serialize Wave Tasks that touch
  overlapping paths. Interim for the dogfood: run cli-heavy specs at
  `worktree.concurrency: 1`.
- **R4-4 outcome under-reports failures as "skipped"**: the deterministic
  implement outcome line re-reads each Task's status from the Run Worktree,
  where a failed Task's `failed` status never lands (it stays in the kept Task
  Worktree, un-integrated), so genuine verification/integration failures print
  as `skipped` with `0 failed`. The failures are visible in the stream
  (`Task task_NN failed: …`) but the summary misleads. Fix candidate: report
  the settled Task outcome the Daemon observed, not only the Run-Worktree file
  status.
- **R4-5 Run Event Journal never pruned → unbounded DB**: `run_events` is
  append-only with no retention, so `~/.roundfix/roundfix.db` reached ~220 MB
  (almost all journal rows carrying every agent payload). Routed to spec
  **0014-run-store-retention** (ADR-0033): retention window + `roundfix gc`
  (`--dry-run`) + best-effort preflight-sweep prune, pruning journals and
  artifact dirs of terminal Runs only — never `runs` rows or active-run locks.
  Marcio's call 2026-07-06: retention/prune + opt-in logs (0011 already makes
  the redundant per-Batch agent log files opt-in); keep the durable Run state.

## Spec 0009 — Parallel Scheduling (the "cited layer")

The worktree isolation from 0008 was built to host this. Turn sequential
per-Run execution into a ready-set DAG scheduler running independent
Tasks/Batches concurrently, each in its own worktree, integrating back on
one Run Branch.

- Ready-set over `needs` (spec) / independent Batches (review) replacing the
  fixed-size sequential loop; a concurrency cap (config, default sensible).
- **Worktree-per-Task** within a Run: each concurrent unit gets a fresh
  worktree branched from the Run Branch; successful units integrate onto the
  Run Branch (fast-forward or sequenced non-conflicting apply), the Run
  Branch integrates to the user branch as today (ADR-0024).
- **Worktree location config** (new): hierarchy repo Project Config > User
  Config > builtin default. The default is
  `~/.roundfix/worktrees/<repo-slug>/<timestamp-unique-id>/` — the slug and
  timestamp segments are NEVER configurable (they guarantee uniqueness and
  prevent collisions); config only sets the parent directory.
- **Worktree debris cleanup** (R3-1): terminal Runs with no commits beyond
  base get their worktree/branch reaped by the preflight sweep and by
  `stop --force`.
- Ordering, failure isolation (one Task's failure never blocks independents),
  and crash recovery must all survive concurrency; ADRs 0010/0013/0014/0023/
  0024 semantics preserved per unit.

## Spec 0011 — Review Artifacts, Run Logs, and Spec Archiving

Three storage/lifecycle fixes so nothing lands loose in `~/.roundfix`.

- **Review artifact location** (R1-17), hierarchy repo > global > default:
  - Explicit location set in Project or User Config → use it.
  - PR associated to a Spec (via `Roundfix-Spec` trailer or `--spec`) →
    `docs/specs/<slug>/reviews/`.
  - No spec, no config → default `docs/specs/_reviews/pr-<n>/round-NNN/`
    (in-repo, versioned, never loose).
  - Design tension carried from R1-17: in-repo artifacts vs the clean-tree
    preflight and worktree commits — resolve by committing them like `qa/`
    or excluding them from Run commits. Supersedes ADR-0003.
- **Run logs off by default** (extends R1-8): production Runs write no
  per-Run agent logs; a config key (User/Project) enables them for dev. The
  loose `~/.roundfix/artifacts/runs/*` accumulation stops.
- **Spec archiving** (new): after a Spec's Tasks are all completed and QA
  passes, archive it — move `docs/specs/<slug>/` to
  `docs/specs/_archived/<slug>/` with archive metadata. DECIDED 2026-07-05:
  the correct English spelling `_archived` stands (skills already use it
  everywhere; nothing to adjust). Decide in the PRD: a `roundfix archive`
  command, auto-archive on QA pass, or lean on the skill from the loop.
- Small: CodeRabbit issue-title derivation strips markup/emoji (R3-2);
  status-poll stderr dedup (R3-4); merge-readiness docs note (R3-3).

## Spec 0012 — npm Distribution (tag → npx/bunx/global)

Model: `~/dev/onioncry` (Rust binary shipped via npm platform packages).
Adapts to Go, which cross-compiles without per-platform toolchains — simpler.

- Root launcher package `roundfix` (`bin/roundfix.js` shim, `type: commonjs`,
  `optionalDependencies` per platform) so `npx roundfix` / `bunx roundfix` /
  `npm i -g roundfix` all work.
- Per-platform binary packages `@roundfix/cli-<os>-<arch>` (darwin-arm64,
  darwin-x64, linux-arm64, linux-x64, win32-x64) with `os`/`cpu` fields.
- GitHub Actions release workflow on `v*` tags: version-check (tag == npm ==
  app version), `make verify`, cross-compile per platform (Go `GOOS`/`GOARCH`
  matrix — no runner-per-OS needed like Rust), publish platform packages then
  the launcher; feeds the existing `roundfix upgrade` release channel (0007).
- The launcher shim forwards args and exit codes verbatim (roundfix's stable
  exit-code contract must survive the Node wrapper).
- brew deferred and correctly scoped: Homebrew is macOS/Linux only (no
  Windows) and needs a maintained tap/formula; npm covers all three platforms
  from one channel, so it goes first.

## Sequencing (revised 2026-07-06, robustness prioritized)

0009 parallel (SHIPPED, Clean + QA pass) → **0010-run-robustness** (SHIPPED,
R3-5/6/7: config deprecation, session-close reaping, Detached Runs — all four
tasks completed) → **0011-storage-lifecycle** (R1-17/R3-2/R3-3/R3-4 + opt-in
logs + Archive Command) → **0012-npm-distribution** → **0013-codex-runtime-
hygiene** (R3-8: `roundfix doctor` + verified-clean codex spawn). Each is one
`roundfix implement` cycle with Codex, QA pass, then archive.

Spec **0014-run-store-retention** was added 2026-07-06 (R4-5): implement it after
0013 and before the review-surface validation. Queue: 0011 → 0012 → 0013 → 0014 →
review validation (open PR + `roundfix watch --until-clean`).

Specs 0011/0012/0013 were authored 2026-07-06 (PRD + TechSpec + Task Graph each;
ADRs 0029/0030/0031 already existed, ADR-0032 added for R3-8; CONTEXT.md gained
Doctor Command and Archive Command). R3-8 was verified **not** implemented by
0010 and routed to its own spec 0013 (the doctor a doctor "waits for real
demand" note in 0010's non-goals — R3-8 is that demand). Distribution (0012)
still lands before 0013 only if the release should carry doctor; otherwise 0013
can run before 0012 since they are independent.
