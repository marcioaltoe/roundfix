---
date: 2026-08-02
supersedes: 2026-07-31-pause-with-0058-in-remediation.md
---

# Handoff — the queue after 0057

Read `docs/workflow/spec-queue.md` for the ordered list. This document explains
what each Spec needs before it can start, what it costs, and what was learned
running the four that shipped.

## Where the loop stands

Shipped and archived: 0052, 0053, 0055, 0056, 0058, 0060, 0061, 0062.
0058 archived under `qa_override` — its remaining row needs a real tagged
release against six live npm trusted-publisher bindings.

In flight: **0057**, all 14 Tasks completed, Pull Request open, awaiting review
then the QA gate.

## The order, and why

| # | Spec | Blocked on | Notes |
| --- | --- | --- | --- |
| 1 | `0071-verification-cost` | nothing | Do this first. Measured, not estimated. |
| 2 | `0057` | in flight | Finish before starting new work. |
| 3 | `0059-run-storage-compaction` | nothing | Authorized 2026-07-28. |
| 4 | `0064-spec-artifact-consistency-gate` | nothing | Authorized 2026-08-02. |
| 5 | `0063-qa-cycle-economics` | nothing | Authorized 2026-08-02. |
| 6 | `0065-loop-order-and-verification-honesty` | **one path** | See below. |
| 7 | `0070-declared-unreachable-acceptance` | nothing | |
| 8 | `0066-run-teardown` | nothing | Tooling not applicable. |
| 9 | `0068-spec-close-audit` | nothing | Tooling not applicable. |
| 10 | `0069-review-run-targets-its-pull-request` | nothing | Tooling not applicable. |
| 11 | `0067-derived-artifact-boundary` | nothing | Authorized 2026-08-02. |

**0071 is first because it was measured.** `go test ./internal/baseline
-count=1` costs 109s on a warm cache; the two-package form costs 142s. Every
one of 0057's fourteen Tasks carried one as its last Verification line — about
28 minutes per pass, paid again on every retry. And the package holds 176 test
functions against 28 `t.Parallel()` calls on a twelve-core machine: adding
`t.Parallel()` to one test's subtests took it from 29s to 17s. Until this
lands, every later Spec pays the same tax multiplied by its Task count.

**0064 precedes the larger 0063** because it is cheaper and every Spec authored
after it is checked by it — including 0063.

## The one thing needing a decision

Spec **0065** may need `internal/baseline/assets/modules/autonomous-work.json`
to deliver the loop-order correction as a Baseline product rather than as
repository-local guidance. The 2026-08-02 authorization names the `Makefile`
and the owned skills; a module asset is neither. A precedent exists — the
2026-07-30 grant covered exactly that file — but the boundary is confirmed per
Spec, never inherited. Confirm before decomposing 0065.

## Maintainer prerequisites, not blocking the queue

Before the next release, and required for 0058's remaining evidence:

1. Six per-package trusted publishers on npmjs.com. The release runbook carries
   the exact table. npm validates none of it until the first publish attempt.
2. Repository variable `NPM_TRUSTED_PUBLISHING_FALLBACK=1`, or ADR-0084's
   bounded fallback never engages and a misconfigured publisher fails the
   release instead of completing it.

## What running four Specs taught

**A Verification must be able to fail when no work was done.**
`go test -run TestThatWasNeverWritten` exits 0, and four of 0057's Tasks
settled `completed` on exactly that. The hardened form:

```bash
go test ./pkg -run '^TestX$' -count=1 -v | grep -q -- "--- PASS: TestX"
```

The rule generalises past Go: any check whose passing state is also its empty
state — a filter that matched nothing, a grep over a file never created, a
clean `git status` — is not a gate.

**Capture the characterization corpus where the risk lives.** 0057's corpus
recorded plan outcomes in the Baseline package; the regression it existed to
catch appeared in a public command journey in the CLI package, so it could not
see it. Scope the corpus to the surface the change can break, not to the
package the change is written in.

**A remediation Task that introduces vocabulary must carry every document that
declares it, in the same slice.** Three findings came from missing this.

**`roundfix watch` binds to the main checkout, not to the `--pr` argument.**
Run it only from that Pull Request's branch and never switch branches while it
is Active. Spec 0069 fixes this; until then it is operator discipline.

**Return the checkout to `main` after closing each Spec or Run.** It removes
the root cause of the watch defect and makes outstanding work visible: a branch
outside `main` means something is unfinished.

**`reconcile` never integrates.** It classifies an Unresolved Run's branch as
unintegrated and preserves it. Integration is the Supervisor's job:
`git merge roundfix/run-<id>`. A `Clean` Run integrates itself.

**Prove before deleting.** Every branch and worktree removed this session was
first checked with `git diff --name-status` against its target for files unique
to it. When a squash merge deletes the target branch, `reconcile` degrades to
`unknown` and preserves forever — correct, and why Spec 0068 resolves
integration by content.

## Per-Spec cost, measured

| Spec | Tasks | QA cycles | Run time |
| --- | --- | --- | --- |
| 0058 | 8 | 4 | ~1h40m |
| 0062 | 7 | 2 | ~1h52m |
| 0056 | 7 | 2 | ~2h11m |
| 0057 | 14 | 3 | ~5h06m |

Task count multiplies fixed per-Task cost. That is the number 0071 attacks.
