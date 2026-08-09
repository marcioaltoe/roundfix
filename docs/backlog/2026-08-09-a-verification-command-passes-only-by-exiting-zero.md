---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A Verification command passes only by exiting zero, and task authoring does not say so

## Opportunity

The Daemon runs each Task's Verification through `sh -c` and treats a non-zero
exit as failure — `internal/daemon/daemon.go:101`. Nothing in the task-authoring
contract states that, and its absence is easy to author straight past.

The `write-tasks` skill tells an author that a Verification command "must be able
to fail when no work was done" and that a command selecting no cases is vacuous.
Both are true and neither is the missing rule. The missing rule is that the
command must **succeed by exiting zero**, which quietly forbids a whole family of
natural-looking assertions.

## Value

Measured on 2026-08-09, authoring Spec 0089. Six of its eight Tasks carried at
least one command whose intended pass condition was an empty result or an absent
string:

```text
git diff --name-only -- internal | grep -v '_characterization_test\.go$'
grep -rn 'must be empty for runtime' internal/config
```

The first exits 1 when it matches nothing, which is exactly when the Task
succeeded. The second was authored with the note "expected: exits non-zero",
which the Daemon cannot honor at all. Every one of those commands would have
failed its Task for doing the right thing.

It cost a Run. Task 01's Verification failed on the first attempt, the Agent
spent its repair turn diagnosing an authored gate rather than its own work, and
the diagnostic artifact the Daemon captured was **zero bytes** — so the failure
feedback carried no cause. That empty-diagnostic behavior is already filed from
2026-08-08 under a different root cause; this is a second way to reach it.

The working forms are unremarkable once known, and all three exit zero on
success under `sh -c`:

```sh
test -z "$(cmd | grep -v PATTERN)"
! grep -rq 'PATTERN' path
! cmd | grep -q 'PATTERN'
```

## Shape

Non-binding. The cheapest fix is one sentence in the authoring contract next to
the vacuity rule it belongs beside, with the three forms above as the worked
answer — `write-tasks` is a repo-owned authorial skill, so that edit needs
express maintainer authorization with bounded files.

Worth settling in the same work: whether a Verification command that produces no
diagnostic output should be surfaced differently from one that fails loudly,
since an empty artifact is what turned a one-line authoring slip into a spent
Agent turn. And whether the authoring contract should carry any executable
check at all — a Spec's gates are today only exercised the first time a Daemon
runs them, which is the most expensive moment to discover they are wrong.
