# QA command evidence — build 54fbf020

All commands ran from the Run Worktree root on 2026-08-05. Commands were
unpiped when their exit status was used as gate evidence.

## Graph, authorization, and changed paths

- `rtk git rev-parse HEAD` → exit 0:
  `54fbf0206682e79b6025ae565a5776fc6482cfe5`.
- `_tasks.md` names `task_05` as the sole `qa:` node. Its only dependency is
  completed `task_06`, which is the only non-QA leaf. Tasks 01-04 and 06 are
  `completed`; task_05 remains Daemon-owned `pending`.
- `rtk git merge-base --is-ancestor e257a7e4 f894538a` → exit 0. The Spec and
  standing authorization commit predates the standalone protected Skill task.
- `rtk git diff-tree --no-commit-id --name-only -r cbbbafc1` → task_01 file and
  `internal/reviewsource/coderabbit/{coderabbit.go,coderabbit_test.go}`.
- `rtk git diff-tree --no-commit-id --name-only -r 8b20dc85` → task_02 file and
  the same bounded CodeRabbit implementation/test pair.
- `rtk git diff-tree --no-commit-id --name-only -r a61edff4` → task_03 file and
  bounded CLI, watch, and run-event implementation/tests.
- `rtk git diff-tree --no-commit-id --name-only -r f894538a` → task_04 file,
  canonical/mirrored Roundfix Skill, and only sanctioned files beneath
  `internal/baseline/assets/setups` and `internal/baseline/testdata`; no Go
  source.
- `rtk git diff-tree --no-commit-id --name-only -r 54fbf020` → exactly
  `task_06.md`, `internal/cli/cli.go`, and `internal/cli/cli_test.go`.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` → exit
  0.

## Skill and repository gates

- `rtk make skills-sync-check` → exit 0; 4 Skill tests passed.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` → exit 0; all 14
  required Roundfix-owned skills passed, including `qa-gate` and
  `evidence-gate`.
- `rtk make verify` → exit 0. Summary: 3,406 Go tests across 26 packages; 1
  isolated corpus-budget test; 4 Skill tests; full Skill check; build of
  `bin/roundfix`; `bin/roundfix spec check` reported no finding for Spec 0077.
  The Spec check recorded the pre-existing optional skips for no TechSpec
  Vocabulary Contract and no references index; neither is a finding or an
  acceptance promise of this Spec.
- `rtk git -c core.fsmonitor=false status --short` after verification listed
  only this untracked QA rerun report/evidence; the gate introduced no tracked
  mutation.

## Prior finding F-001 reproduction

- `rtk proxy ./bin/roundfix watch --help` → exit 0. The built public help says:
  “The only check-or-status route to a verified head is a recognised
  review-completed current-head CodeRabbit check or commit status”; an
  unrecognised successful signal resolves pending because “a green check is not
  evidence that a review ran”; a refusal resolves skipped and watch “will not
  merge that head or clear it for merge.” The removed permissive phrase
  `a successful CodeRabbit check or commit status` is absent.
- `rtk go test ./internal/cli -count=1 -run '^TestRunCommandHelp$/watch$' -v` →
  exit 0; 2 test results passed in the CLI package.

## Recorded Review Source and assembled command flows

- The focused CodeRabbit classifier command ran eight named test groups and
  exited 0. Its visible subtests passed for:
  - recorded PR #107 rate-limit payload → `skipped`;
  - recorded completed review → `verified`;
  - unknown successful check name and output → `pending`;
  - recognised completion with no unresolved thread → `verified`;
  - recognised completion with an unresolved thread → `reviewed`;
  - current approval verified/reviewed hierarchy;
  - rate-limit and path-filter refusal classes, including title-case variants;
  - missing authoritative comment staying pending;
  - stale completed/refused evidence not settling the current head;
  - bounded refusal reason and GitHub rate-limit comment JSON mapping.
- `rtk go test ./internal/reviewsource/coderabbit -count=1` → exit 0; 77
  CodeRabbit package tests passed, including the recorded corpus.
- The focused watch command ran four named tests and exited 0:
  `TestRunReviewSkippedStopsBeforeFetch`,
  `TestRunUnrecognisedPendingEvidenceStopsBeforeFetchWithDiagnostic`,
  `TestRunReviewSkippedDuringMergeReadyPreservesTerminalEvidence`, and the
  adjacent verified-evidence canary
  `TestRunReviewEvidenceSharedByPreFetchAndMergeReady` all passed.
- The focused public CLI runner ran four named flows and exited 0. The refusal
  and unrecognised-green cases published their reasons, returned non-Clean,
  and asserted no Final Push; the fetched-zero and inherited verified-evidence
  canaries still reached their expected positive behavior.
- The focused Run Event Stream projection ran three named tests and exited 0.
  Refusal projected `evidence_state: skipped` plus the refusal reason; unknown
  projected `evidence_state: pending` plus the unrecognised diagnostic; the
  broader outcome-context projection round-tripped successfully.

## Policy and Non-Goal sweep

- `rtk git diff --name-status 0bfa6c4e..54fbf020 -- internal/reviewsource/coderabbit internal/watch internal/runevent internal/cli`
  → exit 0; the assembled product delta is limited to nine existing files in
  the classifier, watch, Run Event Stream, and CLI surfaces named by the Spec.
  No provider configuration, fallback provider, QA-verdict, or new accepted-gap
  surface was added.
- `rtk git diff -G 'retry|re-request|retrigger|backoff|capacity' --unified=0 0bfa6c4e..54fbf020 -- internal/reviewsource/coderabbit internal/watch internal/runevent internal/cli`
  → exit 0 with empty output; no automatic retry, re-request, retrigger,
  backoff, or capacity policy was introduced in the product delta.
- The canonical Roundfix Skill states that a refusal is not a transient
  failure and that Roundfix does not automatically retrigger or retry a refused
  head; it assigns that policy to follow-on work.

## Closed-report inspection

- Frontmatter inspection → `status: closed`, `verdict: pass`,
  `rows_blocked_environment: 1`, `rows_blocked_finding: 0`, and
  `rows_blocked_declared: 0`.
- Results-table parser → exit 0 with
  `rows=19 pass=18 environment=1 fail=0 pending=0 skipped=0`.
- Planned-matrix pending-row parser → exit 0; no terminal matrix row remains
  pending.
- Evidence-path existence and trailing-whitespace checks → exit 0.
- Final Git status contains only this collision-safe QA report and its evidence
  directory; no Task status, product source, commit, branch, push, or Pull
  Request was changed by QA.
