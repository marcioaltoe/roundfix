# Conexus Phase 06–07 — Roundfix field report (kept for analysis)

Source: Fable session in ~/dev/conexus, 2026-07-07. Scope: 8 Runs, 21 Tasks,
3 QA gates, 3 archives, 1 settle, 2 transport recoveries; symlink fix
validated in production mid-cycle. Original also at the reporter's
scratchpad/roundfix-improvement-report.md. **Not scheduled — parked for
triage.**

## 1. Confirmed bugs (by priority)

1.1 **Commit-hook failure kills the whole Run (HIGH)** — Run …3c9b6e7a:
task_03 verified and settled, then the daemon's `git commit` failed (husky
ran oxlint over root `scripts/`, outside the turbo lint scope) and the whole
Run died, killing the remaining graph. Asks: commit failure = task failure
(Run stays alive for independent tasks); hook stderr in the failure reason;
docs note "Verification ⊇ hook checks".

1.2 **acpx: session lost between Batches is not recreated (HIGH)** — Run
…16659fd4: session vanished after a long Batch (suspect: acpx `ttl: 300` vs
multi-minute verification); the next 2 tasks failed in seconds with missing
session, zero work. Asks: recreate session; retry on infra-class error
(exit 4); keep-alive/TTL guidance.

1.3 **acpx: agent/protocol error kills a nearly-done Batch (MEDIUM)** — task
with 465/468 tests green killed mid `bun run test` (10MiB buffer class).
Asks: single retry + acpx stderr tail in the log. Positive: worktree +
settle recovery worked perfectly.

1.4 **claude route broken at session/new (MEDIUM)** — `acpx claude sessions
new` → Internal error (reproduced outside Roundfix; direct handshake with
adapter 0.16.2 OK; drift vs acpx 0.12.0). Asks: preflight probes the
SELECTED agent (not only default); error naming the adapter command +
stderr; version matrix; doctor-style check for the claude route (also: an
orphaned `claude-agent-acp` shim from an asdf migration). NOTE: partially
addressed upstream 2026-07-07 — CLAUDECODE nested-guard strip (6e5668c) +
local adapter install; the acpx↔adapter version-drift and selected-agent
probe asks remain.

1.5 **Preflight sweep blocks all Runs (MEDIUM)** — a failed-bootstrap
worktree (untracked debris) had a reapable branch but `worktree remove`
refused without `--force` → every new implement exited 2 until manual
cleanup. NOTE: fixed upstream in 0022 (forced removal) — verify it also
covers the sweep-blocks-preflight path.

1.6 **`defaults.auto_commit: false` does not govern implement (LOW)** —
daemon committed anyway; only observed coupling is with `watch.auto_push`
in preflight. Ask: document or extend.

## 2. Design/UX

- Per-task runtime routing for mixed-runtime specs (`runtime:` frontmatter
  or type→agent map) — today's alternative is timed graceful-stop +
  relaunch with another `--agent`.
- `roundfix ps` / `status <run-id>` (ideally `--json`) — orchestrating
  detached today = polling console.log + frontmatter. NOTE: partially
  addressed by 0020 `runs list`; a per-run status/json surface remains open.
- Detached console log is noisy — skill dumps bury daemon milestones;
  truncate tool payloads (Journal keeps everything). NOTE: cockpit timeline
  got bounded summaries in 0021; the console.log itself remains raw.
- Title lowercasing mutilates acronyms — `mCP`, `sDK`; spare all-caps
  tokens.
- Dev build receives upgrade nag for an older release.
- Symlink-fix residual: status flips / `## Result` stay uncommitted in the
  external checkout → dozens of manual reconciliation commits per cycle;
  an opt-in auto-commit in the resolved-root repository would close the
  loop.

## 3. Discovery/docs

- `qa-report-*.md` naming convention only discoverable via strings on the
  binary (the error should name the glob).
- Verification exit-0 contract is implicit (prose bullets break silently).
- `journal.jsonl` missing at the documented path on a Failed Run.
- Docs note for monorepos with a host-singleton docker stack (never
  `compose up` from a worktree — caused the first failed Run).

## 4. What worked very well

The symlink fix (all boundaries); settle end to end; the worktree model
(07-b in parallel with 07-a's QA without friction; IntegrationPending
predicted and resolved with a clean rebase); automatic integration on Clean
(07-c Run 2: the ideal flow, zero touch); archive with correct gates; and
the status flips in the shared .knowledge as a live monitoring surface.

## 5. QA gate value evidence

The runtime gate found 5 defects invisible to hermetic tests: wholesale
catalog 404, missing includeSandbox, SDK double-init under StrictMode
(named popup + effect without cleanup), substring matcher in the validation
instrument, and a README teaching `bun --cwd <pkg> run <file>` — the exact
pitfall the Roundfix skill already documents (worth promoting into the
assigned-Batch contract).
