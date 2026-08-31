---
status: completed
type: test
---

# Task: The CLI fixtures build a qa seed the way Roundfix now renders it

Corrective, from QA finding F-001, and consequent rather than independent: the
derived Verification made these fixtures stale. `implementTaskContent` assigns a
per-seed fixture command to every seed including `qa` ones, so every assembled
CLI journey now stops at the mechanical stage with
`SC-QA-VERIFICATION-AUTHORED`, and `rtk make verify` exits 2.

## Work

- A `qa` seed with no explicit command renders the command Roundfix derives for
  its Spec, exactly as a real graph now must. Reference the derivation rather
  than transcribing its text: a fixture that copies the string stops testing the
  day the derivation legitimately changes.
- A seed that names its own command keeps it, so a fixture can still exercise
  the refusal deliberately.
- Leave the non-`qa` seeds alone. Their fixture command is what makes those
  journeys meaningful.
- The whole `internal/cli` suite passes on a fresh cache, which is where the
  failure appeared.

## References

- `_prd.md` → Goal 3, Core Feature 3
- QA finding F-001 in `qa/qa-report-2026-08-31.md`
- The consequent-fix ordering in `docs/agents/agent-instructions.md`: this lands
  after the change that made it necessary, never folded into it

## Verification
- `grep -q "DerivedQAVerification" internal/cli/implement_test.go || exit 1; go test -count=1 ./internal/cli`

## Result

Implementation:

- `implementTaskContent` now uses `spec.DerivedQAVerification(slug)` when a
  `qa` seed has no explicit Verification command. The fixture follows the
  production derivation rather than copying its command text.
- A non-empty seed Verification still takes precedence for `qa` seeds, so
  refusal fixtures can keep an authored command deliberately.
- Non-`qa` seeds still default to the scripted Agent-marker command.
- `TestImplementTaskContentChoosesVerificationByTaskType` locks those three
  cases to one rendered command each.

Focused-check evidence:

- Before the implementation change,
  `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test ./internal/cli -run '^TestImplementTaskContentChoosesVerificationByTaskType$' -count=1`
  failed only `qa_seed_derives_its_default`: the fixture rendered the
  Agent-marker command instead of the value returned by
  `spec.DerivedQAVerification`.
- After the final code edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test ./internal/cli -run '^(TestImplementTaskContentChoosesVerificationByTaskType|TestRunImplementUsesConfiguredExternalSpecRootEndToEnd|TestRunImplementSettledQAGateReportsWithoutRun|TestRunImplementQAVerdictMatrix|TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict)$' -count=1`
  passed. This covers the three fixture-selection cases and four assembled QA
  journeys without running the Task's declared Verification.
- `rtk gofmt -d internal/cli/implement_test.go` produced no diff.
- `rtk git diff --check` passed.
- Changed-path inspection found only `internal/cli/implement_test.go` and this
  assigned Task file; the Daemon-owned `status: in_progress` edit was
  preserved.

The referenced `qa/qa-report-2026-08-31.md` is absent from this worktree. The
Task's embedded F-001 reproduction and the focused pre-fix failure supplied the
available defect evidence. The full fresh-cache `internal/cli` suite was not
run because it is the Task's declared Verification and remains Daemon-owned.

### Verification feedback — attempt 1

- The Daemon diagnostic was inspected at its supplied artifact path. The suite
  reached the CLI tests and timed out in
  `TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks`
  while waiting for two Agent starts after observing only `task_01`.
- That cancellation test creates two non-`qa` seeds and gives both an explicit
  Verification command. It does not execute the `qa` default-selection branch
  changed by this Task, and the diagnostic did not report a fixture derivation
  or `SC-QA-VERIFICATION-AUTHORED` failure.
- `rtk env GOCACHE=/private/tmp/roundfix-task08-feedback-gocache go test ./internal/cli -run '^TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks$' -count=1`
  passed in isolation.
- `rtk env GOCACHE=/private/tmp/roundfix-task08-feedback-gocache go test ./internal/cli -run '^(TestImplementTaskContentChoosesVerificationByTaskType|TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks)$' -count=1`
  passed, so the fixture-selection regression and the reported cancellation
  test do not reproduce an interaction.
- No code changed during this feedback turn. Altering the unrelated
  cancellation test or weakening its synchronization assertions would exceed
  this Task's fixture slice. The declared Verification remains for the
  Daemon's single rerun.
