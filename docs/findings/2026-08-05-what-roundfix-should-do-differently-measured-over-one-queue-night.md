---
date: 2026-08-05
surface: internal/cli, internal/daemon, internal/spec, .agents/skills
status: open
---

# What Roundfix should do differently, measured over one queue night

Roundfix drove five Specs from PRD to archive in one session. The Agent work was
not the constraint: **33 implementation Tasks settled on Verification attempt 1,
and exactly one failed on the work.** Every other failure came from the
Supervisor's authoring or from the loop's own mechanics.

This records the product behaviours that cost the most, each with what it cost.
It is written from the Supervisor's seat and corroborated by eight field
findings written from fluxus and vortex sessions, which observe the same ceiling
from the consumer side.

## 1. Progress is not readable without writing a parser

`roundfix runs list` gives a Run's state. `roundfix attach` gives a human the
Live Run View. Between them there is no compact, non-interactive answer to
*which Task is running and which have settled.*

Every status check this session piped `roundfix events --filter task-status`
through a hand-written Python one-liner to reduce JSONL to eight lines. For a
Supervisor that is an agent rather than a person, that is the difference
between checking progress and building a tool to check progress.

Worse: `roundfix events --follow` emits nothing to a pipe until it exits. A
1 h 28 min Run reported **once, at the end**. Twice this session I described a
Task as still running when it had settled twenty minutes earlier, because the
stream I was reading had buffered.

**Suggested behaviour.** `roundfix runs show <id>` — one line per Task with
phase, status, and elapsed; exits immediately; `--format json`. And make
`--follow` line-buffered so a pipe sees events as they happen.

## 2. Nothing watches a Pull Request after it is opened

Spec 0070's review landed 17 minutes after its Pull Request opened and sat
**1 h 32 min** unaddressed, because the Supervisor was monitoring the Run and no
signal arrives when a review lands. `roundfix watch --until-clean` is exactly
the tool, and it was started only after the maintainer asked for status.

The loop documented in `docs/agents/autonomous-work.md` already says to open the
Pull Request and then watch. The gap is that they are two commands, so the
second is skippable — and it was skipped nine times before it mattered, because
every earlier review was skipped by path filters.

**Suggested behaviour.** Let the Pull Request step start the watch: a flag on
`watch` that opens the Pull Request first, or an `implement` mode that chains
into it. One command that cannot be half-done.

## 3. A gate failure has no path back

The QA report states the required repair in prose — often precisely. Nothing
consumes it. Every cycle this session went: read the report, decide whether the
defect is authoring or code, write a corrective Task by hand, amend the
manifest, reset the gate Task to `pending`, relaunch.

That last step is undocumented and was discovered by reading
`internal/spec/spec.go`: `StaleGateError` refuses to load a graph whose gate is
`failed` while its dependency closure grew, so the gate **must** be reset to
`pending` by hand or the next Run will not start.

**Suggested behaviour.** At minimum, make the refusal name the fix —
"reset `task_09` to `pending`" — instead of describing the invalid state. Better:
let `implement` reset an authored gate whose closure grew, since that is the
only valid response to the condition it detects.

## 4. Fail-fast Verification spends the only repair turn on the first of N defects

Corroborated independently from a fluxus session
([finding](2026-08-04-fail-fast-verification-spends-the-single-repair-turn-on-the-first-of-n-defects.md)).
A Task with two independent defects cannot settle: attempt 1 stops at the first,
the single Verification Feedback turn repairs exactly that, and attempt 2 fails
on the second — which was never visible.

**Suggested behaviour.** Run the declared Verification commands to completion
and return **all** failures in one Feedback turn. The commands are already
independent; only the early exit makes them serial.

## 5. An accepted Review Issue has no terminal state

Corroborated from a vortex session
([finding](2026-08-04-an-accepted-gap-has-no-terminal-state-so-the-loop-cannot-close.md)):
closing one Pull Request took **six maintainer interventions, four of which
carried no judgement at all**. The proximate cause is narrow — there is no
status meaning *"valid finding, deliberately not fixed, accepted."* Every
available terminal status either lies about the finding or blocks the loop
forever.

This session's own review rounds landed `invalid` on findings that were
genuinely wrong, which worked. The gap is the finding that is **right** and
deliberately not acted on.

**Suggested behaviour.** An `accepted` terminal status carrying a required
reason, distinct from `invalid` and from `resolved`.

## 6. The Task Graph is read from the commit, not the working tree

Every Spec costs **two** Pull Requests: one to make `_tasks.md` readable to a
process running on the same machine, and one for the work. Five Specs this
session, ten Pull Requests, five of them pure ceremony.

**Suggested behaviour.** For a local Run against the current checkout, read the
graph from the working tree — or offer `--graph-from-worktree` explicitly. The
commit requirement is right for a detached or remote Run and wrong for the
common case.

## 7. The tooling boundary splits one repair into two Tasks

An authorized tooling Task may mutate only its bounded files. Spec 0070's F-001
was one fix with two halves — split a test assertion, and add the gate step that
runs it — and became `task_10` and `task_11`. Correct under the rule, and it
consumed **both** corrective Tasks the contract allows on a single finding.

**Suggested behaviour.** Let the corrective-Task ceiling count *findings*, not
Tasks. The rule exists to catch a Spec whose decomposition is wrong, and one
finding split by a boundary is not that.

## 8. Presence-awareness skips the check that matters most

Spec 0064's `SC-TOOLING-UNAUTHORIZED` — the detector for the defect that put it
first in the queue — is skipped **entirely** when `_techspec.md` is absent, which
is the majority of a live queue. Spec 0059's PRD cited a 2026-07-28
authorization record that does not exist, and the sweep reported it clean.

ADR-0094 was implemented as *skip the detector if any input is missing* rather
than *check the artifacts that are present*.

**Suggested behaviour.** Per-artifact evaluation: run every check whose inputs
exist, and record a skip only for the specific artifact that is absent.

## 9. `--delete-branch` destroys the handle reconcile resolves by

Fixed by Spec 0068, recorded here as the shape to avoid repeating: a Run whose
target branch its own squash merge deleted degraded to `unknown` forever,
because the name was the identifier. Content is what integration means.

**Suggested behaviour, generalised.** Where Roundfix records a Git ref as an
identity, it should be able to fall back to content. Refs are deleted by normal
workflow; content is not.

## Evidence

- Session of 2026-08-04 into 2026-08-05: Specs 0064, 0076, 0070, 0068, 0066
  archived; 0067 stopped on a maintainer decision; 21 Pull Requests merged.
- Task settlement: 33 implementation Tasks, 1 failure attributable to the work.
- Timings: Spec 0070's review latency 21:23→22:55; Spec 0066's Run 1 h 57 min
  reported once at exit.
- Corroborating field findings from fluxus and vortex sessions, versioned in
  `docs/findings/` on 2026-08-04.
- `docs/handoffs/2026-08-05-the-night-the-queue-moved.md`.
