---
task: task_04
spec: 0078-roundfix-asks-for-the-review
status: completed
type: chore
complexity: low
---

# Task 04: Turn the flow on for this repository

## Overview

The Task that must not land early. `.coderabbit.yaml` is already modified in
the working tree with manual review, and it takes effect when it reaches the
default branch — before the requester works, that file is what strands every
unattended Run.

This slice commits it together with the Project Config that makes Roundfix ask,
so manual mode and the asking arrive in the same change.

## Requirements

1. MUST add this block to `.coderabbit.yaml` as the first entry under
   `reviews:`, exactly as the maintainer authored it on 2026-08-05:

   ```yaml
   reviews:
     auto_review:
       enabled: false
       description_keyword: "coderabbit:review"
   ```

   Every other key in that file stays byte-identical; this Task adds three
   lines and changes nothing else. The keyword MUST be written explicitly: the
   Review Source schema defaults `description_keyword` to the empty string, and
   `enabled: false` with an empty keyword reviews nothing at all, silently.
2. MUST set `review_source.request_review: true` in `.roundfixrc.yml`, so this
   repository asks for the review its Review Source no longer starts.
3. MUST lower `watch.max_rounds` from `3` to `2` in `.roundfixrc.yml`, capping a
   pull request at three reviews.
4. MUST leave the resulting pair coherent under task_03's predicate, proven by
   running the repository's own Preflight Validation rather than by reading it.
5. MUST NOT change Go source. This Task turns on behaviour that already shipped.

## Subtasks

- [ ] Commit `.coderabbit.yaml` as the maintainer left it.
- [ ] Add `review_source.request_review: true` and set `max_rounds: 2`.

## Acceptance Criteria

- [ ] `.coderabbit.yaml` carries manual review with the `coderabbit:review`
      keyword, written explicitly rather than left to the empty default.
- [ ] `.roundfixrc.yml` enables the request and caps rounds at `2`.
- [ ] The repository's own configuration pair passes task_03's Preflight
      predicate, proven by a command rather than by inspection.
- [ ] No Go file changed.

## Context

- interface: `.roundfixrc.yml`
- interface: `.coderabbit.yaml`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `grep -q 'enabled: false' .coderabbit.yaml && grep -q 'description_keyword: "coderabbit:review"' .coderabbit.yaml`
  — expected: exit 0; manual review is committed with a non-empty keyword.
- `grep -q 'request_review: true' .roundfixrc.yml && grep -q 'max_rounds: 2' .roundfixrc.yml`
  — expected: exit 0; the repository asks, capped at two Rounds.
- `go run -buildvcs=false ./cmd/roundfix doctor > /dev/null` — expected: exit 0;
  the repository's effective configuration loads and validates.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Decisions; Success Metric 3.
- `_techspec.md` → Build Order 4; Risks & Considerations.

## Result

Configured this repository for the manual Review Source flow. CodeRabbit now
disables automatic review and retains the explicit `coderabbit:review`
description keyword. Project Config now asks for each post-push review and
sets Max Rounds to `2`.

Focused evidence:

- Manual-review configuration: `rtk git diff -- .coderabbit.yaml` showed only
  the required three-line `reviews.auto_review` insertion, with
  `enabled: false` and `description_keyword: "coderabbit:review"` as the first
  entry under `reviews:`.
- Project Config: the focused Preflight probe loaded the repository's actual
  `.roundfixrc.yml` through `config.Load` and reported
  `request_review=true max_rounds=2`.
- Preflight coherence: the same probe passed the loaded request value and this
  repository root to `preflight.Run` for `watch`; it exited `0` and reported
  `preflight=ok`.
- Go-source scope: no Go source was edited. The probe source and binary lived
  under `/tmp/roundfix-task04.eaopGG`, outside the repository.

The failed Daemon diagnostic was inspected before repair. No command from
this Task's `## Verification` section was rerun; the Daemon owns the final
configured Verification sequence.
