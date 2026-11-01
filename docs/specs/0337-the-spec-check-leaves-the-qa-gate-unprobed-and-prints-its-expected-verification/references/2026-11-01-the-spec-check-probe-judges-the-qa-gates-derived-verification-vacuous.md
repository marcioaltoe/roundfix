---
type: fix
status: promoted
created: 2026-11-01
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
reason: null
---

# The Spec Check probe judges the QA gate's derived Verification vacuous

## Symptom

`roundfix spec check --strict --run-verification` exits `1` on a Spec whose
pending QA gate depends on an honest corrective Task, when a failed QA Report
from an earlier Run is kept under the Spec's `qa/` folder. The probe reports
the corrective Task `honest` and the QA gate's derived Verification `vacuous`:
the derived command finds the newest `qa-report-*.md` and only requires a
non-empty `verdict:`, so the kept `verdict: fail` report satisfies it before
any work.

## Where

- `internal/cli/spec_check.go`, `probeSpecVerifications` (about lines
  447 to 480 at `340c9004`): the pre-work probe iterates every Task of
  `graph.Tasks`, the QA gate included, and skips only completed Tasks.
- `internal/spec/task.go`, `taskVerification` (about line 136): a QA Task's
  Verification must equal `DerivedQAVerification` (`slices.Equal`), so the
  author cannot strengthen it; a differing command is refused with
  `SC-QA-VERIFICATION-AUTHORED` (`internal/speccheck/coherence.go`, about
  line 431).
- The Daemon already separates the gate: `taskPlanWithoutQAGate`
  (`internal/daemon/task_engine.go`) removes the Task of type `qa` before the
  scheduler, so the pre-work probe never runs on it; the gate's acceptance is
  `spec.QAReportEligibility` and the audited head (`runQAGate`, about line
  3724 in the reporter's reading), not the shell command.

## Expected

The pre-work effect probe of `spec check --run-verification` is not applied
to the QA gate's derived command: the check reports it not probed with a
reason, as it does for completed Tasks. The renderer, the eligibility rule
and the audited head stay unchanged, and no check is weakened. A regression
covers a kept failed QA report beside an honest corrective Task.

## Evidence

Reported by the Fluxus project (`default/gesttione-solutions-fluxus`)
through ai-memory messages `01a1987b` and `01a1987c` on 2026-11-01, with
Roundfix 0.127.0. Their Spec 0140 parked Run Unresolved on a QA `fail`. The
operator ran `roundfix reconcile --carry-forward`, kept the failed QA Report
and its evidence byte-identical, added a corrective `task_05` and made the QA
gate need it. The strict check with `--run-verification` then exited `1`:
`task_05` honest, the QA gate's derived Verification `vacuous`. Their
read-only diagnosis names the four locations above and proposes the fix in
Expected. Meanwhile Fluxus used an authorized QA override. The operator
reproduced the shape on `340c9004` with that commit's `bin/roundfix` in a
scratch clone: a pending `task_01` with `test -f corrective-output` reports
`honest`, and the pending QA gate beside a committed
`qa/qa-report-2026-10-31.md` with `verdict: fail` reports `vacuous`, exit 1.
