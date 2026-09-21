---
status: approved
granted: 2026-09-21
action: apply one declared-acceptance eligibility policy across archive and the derived QA Verification, and keep the shipped skill true to it
consuming: 0152-one-declared-acceptance-policy
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0152

Two grants are consumed here.

The skill files ride the standing authorization of 2026-09-18, valid across the
remaining slices of this queue for one purpose: keeping the shipped skill true
to the CLI behavior the slice itself delivers.

`internal/spec/archive.go` was granted on 2026-09-21, bounded to calling a
single shared eligibility decision instead of carrying its own copy of the rule.

`internal/spec/archive_test.go` was granted on 2026-09-21, after the QA gate
found the first record had missed it. Task 02 pins archive's accepted and
refused reports with tests, and that is where archive's tests live. The record
had named the source file and not its test, which is the same class of miss that
produced the widening on Spec 0151 — a path predicted rather than measured
against the governed set.

## Why each governed path is unavoidable

Archive is one of the two places that decide whether the newest QA Report is
acceptable, and it holds the policy the repository documented. Leaving it
untouched would mean adding a second authoritative implementation elsewhere and
holding the two in sync with a test — the arrangement that produced the defect
this Spec repairs.

The derived Verification command is public contract surface the skill describes,
and this Spec changes what it accepts.

`internal/docscontract/testdata/corpus-golden.json` is covered by the standing
authorization of 2026-09-21, bounded to recording the corpus counts this tree
actually measures.

The golden did not go stale because of this Spec. Measured both ways: with Spec
0152 present the sweep counts 9 `SC-CONTRACT-UNDECLARED` and 9
`SC-METRIC-UNDECLARED`, and with it removed it counts the same 9 — this Spec
declares both, so it contributes nothing. The golden said 10 because retiring
Spec 0128 dropped the active corpus from ten Specs to nine, and that commit went
straight to `main` without a pull request, so the docs gate never ran on it.
`main` is red on `make verify-docs` today; this pull request is what surfaced it.

The same counts are pinned in two places. `internal/spec/archive_layout_characterization_test.go`
carries them hardcoded alongside the golden's `update` prose, and compares all
three, so both files move together or neither passes. It is covered by the same
standing authorization, bounded to recording the measured counts.

## What is not governed

Computed as the intersection of this Spec's changed paths with the literal set
in `internal/speccheck/governed.go`, rather than by probing paths guessed in
advance. Exactly two governed files are touched, both listed above. Everything
else the Spec changes — `internal/spec/qa.go`, `internal/spec/qa_test.go`,
`internal/spec/task.go`, `internal/spec/task_test.go`,
`internal/spec/errors.go`, `internal/cli/cli.go`, `internal/cli/qa_report.go`
and `internal/cli/qa_report_test.go` — is ordinary source that no authorization
has bounded.

## Approved bounded mutation

In `internal/spec/archive.go`: replace the inline verdict judgement with a call
to the single shared eligibility decision. Archive's own behavior does not
change — the same reports are accepted and the same reports are refused, for the
same reasons.

In `internal/spec/archive_test.go`: pin that behavior, so the refactor is proved
rather than asserted.

In the Roundfix skill and its mirror, state:

- that one eligibility policy decides whether the newest QA Report is
  acceptable, and that settlement and archive both apply it;
- that a `partial` verdict whose blocked rows are declared unreachable settles
  the terminal `qa` Task;
- that `fail`, an undeclared partial, a missing or unparseable report, and a
  `pass` carrying blocked rows each still refuse.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`. If any derived pin changes as a result, rewrite it only with
`make baseline-digests`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No action, operation or path beyond those above.
- No change to which QA Report is newest, or to what makes a blocked row
  unreachable.
- No change to the skill's version, or to any behavior it documents beyond this
  Spec's own.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
