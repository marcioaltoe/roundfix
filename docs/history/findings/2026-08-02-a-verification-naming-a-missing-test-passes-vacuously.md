---
status: done
created_at: 2026-08-02
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# A Verification naming a missing test passes vacuously (2026-08-02)

`go test ./pkg -run TestThatWasNeverWritten -count=1` exits **0**. The Daemon
runs each Task's declared Verification verbatim, reads exit 0, and settles the
Task `completed`. A Task whose Agent implemented nothing therefore settles
`completed` as long as its Verification names tests that do not exist.

This is the generalised form of
[a rehearsal Task can settle completed without rehearsing](2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md).
That report described one Task whose Verification passed most easily when no
work was done. This is the same defect as a reusable mechanism, and it was
authored into more than sixty Verification blocks before anyone noticed.

## Symptom / evidence

Spec 0057's QA gate reported four Tasks settled `completed` whose commits
changed only their own Task file:

```text
$ git diff-tree --no-commit-id --name-only -r 6192029 4205b14 604f958 3230efb
docs/specs/0057-.../task_04.md
docs/specs/0057-.../task_06.md
docs/specs/0057-.../task_07.md
docs/specs/0057-.../task_09.md
```

Their Verification blocks named `TestDivergenceRendersProbe`,
`TestCapabilityRecheck`, `TestDivergencePromptRemediateOutcome`, and
`TestClauseDeltaRendersBeforeLedger`. None exists:

```text
$ go test ./internal/baseline -run 'TestThisNameDoesNotExistAnywhere' -count=1
No tests found
$ echo $?
0
```

Task 12 then documented those behaviors as shipped, because from its side they
were: four upstream Tasks were `completed`.

## Root cause

Two properties combine.

**`go test -run` treats "no test matched" as success.** This is correct for the
tool — a filter that matches nothing is not an error — and wrong as a
completion gate.

**The Verification author and the implementer are the same agent, one Task
apart.** The Supervisor writes a Verification naming the test the implementer
is expected to write. Nothing binds the two: if the implementer writes neither
the behavior nor the test, the gate that was supposed to catch that passes.

The Daemon is not at fault. It ran the declared commands verbatim, which is its
contract. The commands were vacuous.

## Action / suggestion

Verification that names a test must prove the test **ran**, not merely that the
command exited zero. The portable form already available:

```bash
go test ./pkg -run '^TestX$' -count=1 -v | grep -q -- "--- PASS: TestX"
```

A missing test produces no `--- PASS:` line and the command fails.

Beyond the mechanical fix, the authoring rule that generalises it: **a
Verification command must be able to fail when no work was done.** Any check
whose passing state is also its empty state — a filter that matched nothing, a
grep over a file that was never created, a clean `git status`, a suite that
excludes the changed code — is not a gate. This belongs with the Verification
honesty work already queued.

Worth auditing separately: how many `completed` Tasks across shipped Specs
settled on a vacuous Verification. The Specs that shipped passed QA gates that
exercised real behavior, so the delivered work is probably sound — but the
settlement mechanism was never load-bearing, and "probably sound" is not the
claim the Task status makes.

## What worked — keep

- The QA gate caught it, and caught it precisely: it named the four commits and
  showed each changed only its own Task file. A gate that reads Git rather than
  the Agent's narration is what made this visible.
- The Task Result sections were honest. Three of the four say plainly that the
  behavior was not implemented. The Agents did not claim false work; the gate
  that consumed their claims did not read them.
