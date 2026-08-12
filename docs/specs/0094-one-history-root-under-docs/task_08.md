---
task: task_08
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: low
---

# Task 08: Prove review reaches no history tree

## Overview

The review configuration excludes retired trees by matching a leading underscore
at any depth. The history root carries no underscore, so that match no longer
reaches it and an explicit path-anchored exclusion is required. The same anchoring
closes a hole the prefix match never covered: a review directory inside an active
Spec is not matched today and is therefore reviewed. And the rule-source patterns
admit no negation at all, so non-reachability there is asserted by a comment; this
slice replaces the comment with a check.

## Requirements

1. MUST exclude every path under the history root from review, by a path-anchored
   pattern rather than by a leading-underscore match.
2. MUST exclude a review directory inside an active Spec from review.
3. MUST fail when any configured rule-source pattern matches a path under the
   history root or under a Spec-owned review directory.
4. MUST fail when the exclusion covering the history root is removed or narrowed
   so it no longer matches.
5. MUST read the patterns from the repository's own review configuration rather
   than from a copy, so the check follows an edit to that file.
6. MUST correct the configuration comment that names the old root, so the prose
   and the resolver agree.
7. MUST NOT change any review configuration key other than the path filters, the
   rule-source patterns proven to reach a retired tree, and the comment.

## Subtasks

- [ ] Add the path-anchored history exclusion and the Spec-owned review exclusion.
- [ ] Add the check that rule-source patterns reach no retired tree.
- [ ] Add the check that the history exclusion still covers the root.
- [ ] Correct the configuration comment naming the root.

## Acceptance Criteria

- [ ] A path under the history root is excluded from review, proven by resolving
      the shipped patterns against one.
- [ ] A review directory inside an active Spec is excluded, which the pre-Task
      configuration does not achieve.
- [ ] A rule-source pattern that reaches a retired tree fails the check, proven by
      exercising the check against such a pattern.
- [ ] Removing the history exclusion fails the check, proven by exercising it
      against a configuration without that pattern.
- [ ] The check reads the repository's own review configuration.
- [ ] The configuration comment names the history root the resolver answers.

## Rehearsal Cases

- Case: a path under the history root resolved against the shipped filters;
  Observation: it is excluded.
- Case: a path under an active Spec's review directory; Observation: it is
  excluded, where the pre-Task configuration included it.
- Case: a rule-source pattern that matches a path under the history root;
  Observation: the check fails and names the offending pattern.
- Case: the shipped rule-source patterns; Observation: the check passes, so the
  current configuration is proven rather than assumed.
- Case: a configuration whose history exclusion is absent; Observation: the check
  fails.

## Bounded scope

This Task may create or modify only:

- `.coderabbit.yaml`
- `internal/docscontract/publicdocs_test.go`
- `docs/specs/0094-one-history-root-under-docs/task_08.md`

Express maintainer authorization:
`docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`,
extended on 2026-08-12 to cover the reachability check in this package's existing
canonical suite for contracts that read the published review configuration. An
earlier draft of this Task named no test path, and the Agent correctly refused
twice rather than widen the boundary itself.

## Verification

- `go test -count=1 -tags docscontract ./internal/docscontract -run 'ReviewHistory|ReviewArchiv' -v > /tmp/0094-task-08.log 2>&1; s=$?; grep -q '^--- PASS: .*Review' /tmp/0094-task-08.log || { cat /tmp/0094-task-08.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing check; fails when the check does not exist or selects no cases.
- `! grep -qi 'no tests to run' /tmp/0094-task-08.log` — expected: exits 0, refusing a vacuous run.
- `grep -q '"!docs/history/\*\*"' .coderabbit.yaml` — expected: exits 0, proving the path-anchored history exclusion is present rather than only asserted by a test.
- `! grep -n 'archive root: _archived/' .coderabbit.yaml` — expected: exits 0, proving the stale comment no longer names the old root.

## Context

- interface: `.coderabbit.yaml`

## References

`_techspec.md` → Build Order 8; Integration Points: the review tool; Testing
Approach: the rule-source reachability check. `_prd.md` → Core Feature 9; Goal 5;
User Story 5.

## Result

Implemented the two path-anchored review exclusions and replaced the unsafe
repository-wide rule-source globs with root-anchored patterns. The canonical
docs-contract suite now reads `.coderabbit.yaml`, resolves its path filters,
checks whether any rule-source glob can reach either protected tree, and derives
the expected history-root comment from `internal/spec.ArchiveDir`.

Focused checks:

- Before the configuration edit,
  `rtk go test -tags docscontract ./internal/docscontract -run '^TestReviewHistoryConfiguration$'`
  failed because `!docs/history/**` was absent.
- After the edits,
  `rtk go test -tags docscontract ./internal/docscontract -run '^TestReviewHistoryConfiguration'`
  passed all seven selected cases. These include injected `**/AGENTS.md` and
  `docs/specs/**/reviews/**/*.md` rule sources, plus configurations with each
  required exclusion removed.
- `rtk go test -tags docscontract ./internal/docscontract` ran 62 tests: 60
  passed and the two existing Spec-wide corpus checks failed on the planned,
  out-of-scope `HistoryMove` and `historyMoves` glossary findings and their
  golden count. No Task 08 review-history case failed.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  those commands and settlement.

Acceptance evidence:

- The shipped `!docs/history/**` pattern resolves against
  `docs/history/specs/example/_prd.md` and excludes it.
- The shipped `!docs/specs/**/reviews/**` pattern resolves against
  `docs/specs/example/reviews/round-01/issue.md`; removing that pattern proves
  the pre-Task filters do not exclude the path.
- Each injected unsafe rule-source pattern is rejected, and the error names the
  offending pattern.
- Removing either required exclusion makes the protected path reachable and
  makes validation fail with the missing pattern and path named.
- `readReviewConfiguration` loads the repository-root `.coderabbit.yaml`; no
  copied pattern list is used as the configuration under test.
- The configuration comment names `docs/history/`, derived in the check from
  the resolver's `docs/history/specs` answer.
