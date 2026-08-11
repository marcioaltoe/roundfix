# R06 — outside acceptance evidence

This row does not credit Spec 0080's PRD, TechSpec, Tasks, adopted backlog
copies, or earlier QA report as proof.

## Independent measurement

`docs/specs/_archived/0071-verification-cost/baseline/2026-08-03-after.md`
is the original measurement report. It records, on the same 12-core macOS
machine:

- unchanged `make verify`: 4.9 seconds after the change, versus 142.2 seconds
  before;
- complete `go test ./... -count=1`: 88.9 seconds with a warm build cache and
  94.1 seconds with a cold build cache;
- typical changed-package costs of about 5, 48.5, and 54.5 seconds.

The report explains the cache boundary, lists reproduction commands, and
records three consecutive fresh full runs at 89.3, 87.9, and 88.9 seconds.
Git history places its archived source at commit `8cc814c4` on 2026-08-03,
before Spec 0080's 2026-08-06 adoption commit `142e3232`.

## Independent gate history

`docs/specs/_archived/0079-one-door-for-fleet-knowledge/qa/` contains five
same-day reports, not a Spec 0080 rehearsal. They preserve multiple full gate
rounds and distinct late findings. The exact 92/29/30-minute Run durations
cited by Spec 0080 are not present in those Markdown reports and were not
reconstructed here.

`docs/findings/2026-08-06-rollup-qa-gates-and-verification-evidence.md` is a
nineteen-member repository Finding rollup. Its members independently record
QA cycle cost, detector placement, report selection, fail-fast Verification,
and unreachable Pull Request journeys. Its consolidated learning places
detectors near the defect and preserves independent static findings together.

The independent measurement establishes the PRD's local-versus-complete cost
premise directly. The archived reports and rollup independently establish the
repeated terminal-gate churn premise, with the exact unpublished Run-duration
numbers disclosed rather than inferred.
