# The queue after v0.7.0 (2026-08-27)

Session paused for a machine restart. `main` is at `aba7a020`, working tree
clean, no Active Run, no feature branch open.

## Read this first: the installed binary cannot open the Run Database

Spec 0117 raised the Run Database to **schema 13**. The database is machine-wide
(`~/.roundfix/roundfix.db`) and shared by the seven repositories on this machine
that carry a `.roundfixrc.yml`. The released `roundfix` 0.7.0 supports schema 12
and now refuses every command with:

```text
Run Database has schema version 13, but this binary supports schema version 12
```

The guard fails closed with a clear message and corrupts nothing, but every
repository on this machine is blocked until one of these:

1. **Cut v0.8.0** from `main` — the proper fix, and `main` already carries 0117.
2. `cp bin/roundfix ~/go/bin/roundfix` — unblocks immediately with an
   unreleased binary. A bridge, not a resolution.

This happened because an unreleased schema migration ran against a shared
resource. Worth remembering as a rule: a schema migration in a machine-wide
database is a one-way door for every consumer on that machine, not just the
repository that ran it.

## Shipped and released this session

**v0.7.0** — published 2026-08-26, npm launcher and five platform packages at
0.7.0, GitHub Release with five assets.

- **Spec 0098** — a hook refusal after a passing Verification is classified
  `hook_refused`, including hooks that print only their finding (the
  husky/lint-staged shape, which used to die as a bare git error). The Task
  stays `completed` with work staged; `settle` accepts completed-but-uncommitted
  Tasks, re-runs Verification in place, stages deletions, and integrates. The
  hook-strictness invariant renders into every adopting repository.
- **Spec 0113** — a gate refused at a precondition writes one terminal
  `blocked | precondition` row instead of an empty Results table; the mechanical
  stage reads the newest report only and parses rows solely from the table under
  `## Results`.
- **Archival lifecycles in the Baseline** — findings (through Rollup absorption
  with an `absorbed_by` license), backlog (`open`/`promoted` only), handoffs
  (only on explicit confirmation, all together).
- **ADRs consolidated** 133 → 120: twelve new records (0138–0149) each absorbing
  two or three same-subject decisions, 25 superseded into `docs/history/adr/`.

## Merged after v0.7.0, not yet released

- **Spec 0117 — a Run Window the Preflight owns** (PR #168, squash `aba7a020`).
  6/6 Tasks, QA `pass`, `Clean` on the first Run. `roundfix window
  <set|show|clear>`, repository-scoped in the Run Database; `HH:MM` resolves to
  the next occurrence (07:00 asked at 23:00 means tomorrow); re-setting without
  `--force` never moves a standing cutoff; a passed cutoff refuses `implement`
  from the Preflight with no Run and no side effects. The window governs
  starting, never finishing.

## The queue, in order

The honest note first: since v0.7.0 exactly one Spec shipped — 0117 — and it was
inserted at the front of the queue mid-session, ahead of the item the triage had
measured as the largest waste. That item is still unauthored. Start there.

### 1. Unresolved-Run work reuse — next, and already mapped

An implement Run that ends `Unresolved` leaves its Task commits on the Run
Branch, and the next `implement` reads the checkout, sees `status: pending`, and
re-executes work that was already completed and verified. Measured this session:
about **18 redundant Task executions** across three Runs of Spec 0098, which
stopped only when the Run Branch was merged by hand between Runs.

The seam survey is done, and it makes the work much smaller than it first
looked — **`roundfix reconcile --carry-forward` already does this job**:
it cherry-picks completed Task commits into the checkout, stamps
`## Carry-forward provenance`, and fast-forwards through a temporary detached
worktree. It refuses only on outcome:

- `internal/cli/reconcile.go:172-173` — *"Run %q is not Stopped; carry-forward accepts one stopped Run"*
- `internal/cli/reconcile.go:466` — `if run.State != store.StateStopped { continue }`

So the likely shape is extending the accepted-state set to `Unresolved`, keeping
every existing refusal (no passing verdict, ≠1 settlement commit per Task,
declared inputs moved, dirty checkout, external Specs Root). Verify that reading
before committing to it.

Supporting seams, all confirmed:

- Outcome decided at `internal/cli/implement.go:372-380`; written by
  `CompleteRun` at `:412`.
- **Integration is gated on Clean alone** — `internal/cli/implement.go:383`,
  `if outcome == store.StateClean {`. Nothing integrates on any other outcome.
  This one line is the defect's seam.
- Routing reads the **checkout commit**, not the Run Branch:
  `loadCommittedSpecGraph` at `implement.go:163`, resolved at `:450` against
  `gitState.HEAD`.
- `status: completed` is written to the Task Worktree or Run Worktree, never the
  checkout — `internal/daemon/task_engine.go:1558-1563` — and reaches the
  checkout only through that Clean-gated integration.
- `settle` cannot help: it refuses `pending` and tells the caller to run
  implement (`internal/cli/settle.go:504-512`), so it cannot reach a Task that
  is already `completed` and committed on an abandoned Run Branch.

Evidence to adopt: `docs/findings/2026-08-12-five-unresolved-runs-to-deliver-one-spec.md`
(kept pending for exactly this).

### 2. Spec 0116 — a verdict that states its own scope

PRD authored, source adopted, tooling authorization recorded for four authoring
skills. Needs TechSpec, graph, implementation. `--run-verification` appears zero
times across the authoring skills, so authors run the form of the check that
skips the probe. Measured three times this session: the full unscoped sweep
caught findings a stage-scoped run had skipped, including in this session's own
PRDs.

**Do this when authoring it:** write the authorization record into
`docs/workflow/authorizations/`. The grant exists only in the PRD prose
("granted 2026-08-26 in session"), and the changed-path audit resolves grants
from that directory — without a record there, the audit will refuse the skill
edits.

### 3. kickoff versus implement-spec — a decision, not an implementation

The maintainer's kickoff skill declares `implement-spec` a rival loop that
contradicts `docs/agents/autonomous-work.md` and publishes on its own.
`implement-spec` is Roundfix-owned and shipped. If kickoff is right, the fleet
receives a defective loop plus a note telling readers not to use it. Settle
whether kickoff replaces it before packaging either.

### 4. kickoff as an owned skill, after that decision

Not portable as written: it cites `ADR-0041` (absent here — ADR numbers are per
repository), carries a consuming repository's domain and spend ceiling, declares
a review policy that ADR-0118 says belongs in a typed Baseline decision, and is
in Portuguese while every owned skill ships in English. The session-cutoff half
of it is already delivered as Spec 0117.

### 5. Per-Agent-Session consumption — the blocking question is answered

`docs/history/backlog/2026-08-08-record-usage-per-agent-session.md` (deferred at
the 2026-08-26 triage) named its own blocker: *"whether ACP adapters expose
usage at all needs checking before any event shape is fixed."*

**They do.** The ACP schema defines `PromptResponse.usage` with `totalTokens`,
`inputTokens`, `outputTokens`, `thoughtTokens`; `claude-agent-acp` also carries
`cacheReadInputTokens`. Roundfix discards it in one line —
`internal/agent/acpx_runner.go:1829-1838` decodes the prompt result into a
single-field struct (`stopReason`), so any sibling is dropped by
`encoding/json`. ADR-0008's raw capture does not save it either: that covers
only `session/update` notifications of five known kinds, and the prompt *result*
never becomes a Run Event.

One caveat that shapes the design: the schema marks `usage` **UNSTABLE** —
"not part of the spec yet, and may be removed or changed at any point", with
`x-deserialize-default-on-error: true`. So absence must be an observable state,
never zero — which is what the backlog entry already asked for.

Worth weighing early: every cost claim made this session, including that
`opus/xhigh` costs roughly twice `opus/high` for about one point of result, came
from the binary's built-in recommendation table, not from consumption measured
in this repository.

### 6. Spec 0097 — a wave that cannot collide

Two Tasks a graph declares independent that edit the same file die at
integration. Measured first-hand: 0113's `task_05` completed, passed its
Verification, then failed with `integration conflict: internal/speccheck/mechanical_test.go`
against its sibling `task_07`. Note that 0097's PRD is otherwise unmeasured, and
its authoring rule half rides with the deferred 0107.

### 7. Spec 0105 — the gate's own economics

201 failed Tasks across five repositories in its PRD. Absorbs former Spec 0104
(the test-cache target) as a Task. Deliberately last: optimising the cost of a
verdict is the wrong order while verdicts still discard work.

## Machine state left behind

- **Run Worktrees removed**, 635 MB reclaimed. The six Run Branches they held
  are **preserved** (`git branch --list 'roundfix/run-*'`) — the checkouts are
  gone, no commit is. They report as unmerged only because their work reached
  `main` through squash merges, so ancestry cannot prove what content already
  did.
- **Three Run Worktrees from 2026-08-12 refused release** and were preserved by
  `reconcile`, correctly: their target branch `ma/0094-one-history-root-under-docs`
  resolves ambiguously, so integration cannot be proven. They predate this
  session.
- **Profiles** now follow the binary's own recommendation table — `sol/high` for
  backend-class work, `opus/high` where design judgment leads, `luna/max` for
  bounded docs, chore, and review. The prior configuration ran `opus/xhigh` on
  every backend and frontend Task.
- **Branch `ma/a-gate-that-cannot-certify-its-own-cache`** still exists locally
  and on the remote; its work is in `main` through PR #167.

## Process notes worth keeping

- **Run the full `roundfix spec check <slug>`, never `--stage prd`.** The
  stage-scoped form skips detectors and reports `No findings.` anyway; it hid an
  unaccounted ADR in two PRDs this session, including 0116's own.
- **Test every authored Verification against the unchanged tree before running
  a Spec** — `roundfix spec check <slug> --run-verification` reports each one
  `honest` or vacuous. All six of 0117's were proven honest before its Run, and
  0117 was the only Spec this session to reach `Clean` on the first attempt.
- **Never author two Tasks that rewrite the same file as siblings.** Serialize
  them and say why in `_tasks.md`.
- **Cut the release in this order:** release plan → approval → changelog section
  and version bump in one commit → tag that commit → workflow → copy the
  changelog section into the GitHub Release notes. The v0.7.0 tag was cut before
  the changelog was written and now points at a commit that lacks it.
