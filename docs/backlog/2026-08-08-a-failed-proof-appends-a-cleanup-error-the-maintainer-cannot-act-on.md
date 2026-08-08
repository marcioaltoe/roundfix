---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A failed proof appends a cleanup error the maintainer cannot act on

## Opportunity

When an Exact Agent Selection Proof fails, `roundfix profiles validate` and
`roundfix doctor` print the real cause correctly and then append a second,
unrelated instruction the maintainer has no way to satisfy.

Observed by Spec 0088's QA gate on 2026-08-08 (finding F01), reproducing with a
scratch repository whose `data.preferred` names
`runtime: opencode`, `model: opencode-go/not-a-real-model`,
`reasoning_effort: ""`. The classification is right — `model_not_advertised`,
naming the affected `data` category and the advertised catalog — and then the
output continues:

```
close disposable Agent Session "roundfix-preflight-…": … No named session …
for cwd … and agent opencode; recovery: rerun Agent Selection readiness after
the Session can be closed; next: restore Agent Session cleanup
```

There is no cleanup problem. The disposable Session was never opened, because
the proof failed before opening it; closing a Session that does not exist is the
expected outcome, not a fault. `ProveExactSelection` joins the setup error with
the cleanup error and both reach the surface, so the second one reads as a
second defect with its own recovery instruction.

## Value

The cost is misdirection at exactly the moment a maintainer is diagnosing. The
printed `next:` is the last line, which is the line a reader acts on, and it
points at Agent Session cleanup rather than at the model identifier that is
actually wrong. Spec 0088 measured a whole session lost to a misdiagnosis of a
readiness surface; this is the same failure mode in a smaller frame.

It also weakens a contract the repository states plainly: a failed check carries
one deterministic next action. Two next actions, one of them spurious, is not
that.

## Shape

Non-binding. A cleanup failure that only reports a Session absent when the setup
error already explains why no Session exists is not an independent fault and
should not be joined into the surfaced error — it belongs in a warning at most.
The narrower question worth settling with it: whether any cleanup error should
ever outrank the setup error that caused it in the surfaced `next:` action.
