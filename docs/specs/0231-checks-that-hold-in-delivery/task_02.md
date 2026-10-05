---
task: task_02
spec: 0231-checks-that-hold-in-delivery
status: pending
type: infra
complexity: low
---

# Task 02: CI's pull request job tests the head merged with the current default branch and records it

## Overview

The finding of 2026-10-05 ([a Pull Request check ran on a stale merge](references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md))
is that a re-run of CI's pull request job keeps the event's merge commit and
a reopen can still receive GitHub's old test merge, so no attempt saw the
default branch fix. This Task makes the job check out the head commit, merge
the default branch tip it fetched at job start, verify against that tip and
record it as the `tested-base` notice the Delivery Queue reads (ADR-0236).

## Requirements

1. MUST answer the finding of 2026-10-05 named in the Overview by changing
   `.github/workflows/ci-verify.yml` and no other file; this is the Governed
   Path `_authorization.md` bounds, and the change is its own commit.
2. MUST set the checkout step's `ref` to the pull request's head commit
   (`github.event.pull_request.head.sha`) on `pull_request` events and to the
   default (empty) on pushes, keeping `fetch-depth: 0` and
   `persist-credentials: false` (TechSpec Invariant 10).
3. MUST add, before `Verify changed`, a step named exactly
   `Merge the current base branch`, run only on `pull_request` events, that
   reads the base branch name from `github.base_ref` through an environment
   variable `BASE_REF` (never interpolated into the script), resolves
   `refs/remotes/origin/$BASE_REF` to its full commit id, merges it with
   `--no-ff --no-edit` under a CI-only committer identity passed with `-c`,
   writes `VERIFY_BASE=<sha>` to `$GITHUB_ENV`, and prints exactly one
   `::notice title=tested-base::<sha>` line (API Contract 1). A conflict MUST
   fail the step.
4. MUST remove the `VERIFY_BASE` taken from the event's base commit from
   `Verify changed`, so it verifies against the tip the merge step exported.
5. MUST leave the `push` path unchanged: the `Verify` step with its
   `VERIFY_TEST_TARGET`, `SUITE_BUDGET_SECONDS` and `GOFLAGS`, and
   `Verify docs` on both events; and MUST NOT change any other workflow.

## Subtasks

- [ ] Check out the head commit on pull requests.
- [ ] Add the merge step and its annotation.
- [ ] Point `Verify changed` at the exported tip.

## Acceptance Criteria

- [ ] In a disposable clone whose default branch moved after the head
      branched, the step's script leaves a merge whose parents are the head
      and the moved tip, exports that tip and prints one `tested-base` notice
      for it.
- [ ] The push path and `Verify docs` are unchanged.

## Context

- instruction: `docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`
- interface: `.github/workflows/ci-verify.yml`

## Verification

- `export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1; d="$(mktemp -d)" && script="$(ruby -ryaml -e 'steps = YAML.load_file(ARGV[0]).fetch("jobs").fetch("verify").fetch("steps"); step = steps.find { |s| s["name"] == "Merge the current base branch" } or abort("no merge step"); puts step.fetch("run")' .github/workflows/ci-verify.yml)" && git init -q -b main "$d/origin" && git -C "$d/origin" -c user.name=t -c user.email=t@example.com commit -q --allow-empty -m base && git -C "$d/origin" checkout -q -b feature && git -C "$d/origin" -c user.name=t -c user.email=t@example.com commit -q --allow-empty -m head && head="$(git -C "$d/origin" rev-parse feature)" && git -C "$d/origin" checkout -q main && git -C "$d/origin" -c user.name=t -c user.email=t@example.com commit -q --allow-empty -m moved && tip="$(git -C "$d/origin" rev-parse main)" && git clone -q "$d/origin" "$d/ci" && git -C "$d/ci" checkout -q --detach "$head" && out="$(cd "$d/ci" && BASE_REF=main GITHUB_ENV="$d/env" bash -e -c "$script")" && test "$(git -C "$d/ci" rev-parse HEAD^1)" = "$head" && test "$(git -C "$d/ci" rev-parse HEAD^2)" = "$tip" && grep -qx "VERIFY_BASE=$tip" "$d/env" && test "$(printf '%s\n' "$out" | grep -cx "::notice title=tested-base::$tip")" = 1` — expected: exit 0; before this Task the workflow has no merge step, so the script lookup aborts and the command fails.
- `f=.github/workflows/ci-verify.yml && grep -qF 'github.event.pull_request.head.sha' "$f" && ! grep -qF 'github.event.pull_request.base.sha' "$f" && ruby -ryaml -e 'steps = YAML.load_file(ARGV[0]).fetch("jobs").fetch("verify").fetch("steps"); names = steps.map { |s| s["name"] }; checkout = steps.find { |s| s["uses"].to_s.start_with?("actions/checkout@") } or abort("no checkout"); m = names.index("Merge the current base branch"); v = names.index("Verify changed"); full = names.index("Verify"); ok = m && v && full && m < v && checkout.fetch("with")["ref"].to_s.include?("github.event.pull_request.head.sha") && checkout.fetch("with")["fetch-depth"] == 0 && steps[m]["if"].to_s.include?("pull_request") && !steps[v].fetch("env", {}).key?("VERIFY_BASE") && steps[full]["if"].to_s.include?("push") && steps[full].fetch("env").fetch("VERIFY_TEST_TARGET") == "test-budget"; exit(ok ? 0 : 1)' "$f"` — expected: exit 0; before this Task the checkout has no head ref and `Verify changed` reads the event's base commit, so the command fails.

## References

- `_prd.md` → Goals; Core Feature 2; Success Metric 2; Acceptance evidence
- `_techspec.md` → Interfaces; Invariant 10; API Contract 1; Integration Points; Build Order 2
- ADR-0236; ADR-0179
