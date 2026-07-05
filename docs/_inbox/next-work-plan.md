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

**Open, routed below:**

| Finding | What                                                             | Spec             |
| ------- | ---------------------------------------------------------------- | ---------------- |
| R1-17   | Review artifacts belong with the Spec, not a loose root          | 0010             |
| R3-1    | Force-stopped Runs keep empty worktrees/branches forever         | 0009             |
| R3-2    | CodeRabbit issue titles are raw table fragments (emoji/markdown) | 0010             |
| R3-3    | merge-readiness `missing` path needs a docs expectation note     | 0010             |
| R3-4    | Status-poll stderr line repeats every interval                   | 0010             |
| R1-12   | Prompt-contract drift (templating, work-plan item 5)             | deferred         |
| R1-16   | codex full-access sandbox preset unavailable via acpx            | upstream/observe |

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

## Spec 0010 — Review Artifacts, Run Logs, and Spec Archiving

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

## Spec 0011 — npm Distribution (tag → npx/bunx/global)

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

## Sequencing

0009 (parallel — the layer to implement now) → 0010 (storage/lifecycle) →
0011 (distribution). Each is one `roundfix implement` cycle with Codex, QA
pass, then archive. Distribution last so the released binary already carries
the parallel + storage work.
