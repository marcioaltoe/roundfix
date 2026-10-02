---
type: fix
status: open
created: 2026-10-02
spec: null
reason: null
---

# The review flags an authorized QA override archive

## Symptom

The pre-PR review raised the same finding twice on 2026-10-02, on Specs 0215 and 0218: the Spec "is archived with task_05 failed and the final QA report verdict partial". Both Specs were archived with `roundfix archive --qa-override` under the maintainer's standing authorization for environment-only partials. Each archived `_prd.md` front matter records `qa_override: true`, the approval source, the reason, the observed QA outcome and the QA Task status. Each time the finding parked the item as `corrective-spec-required`, and the operator dismissed it with that evidence.

## Where

The review prompt and its Delivery Conventions, in `internal/cli/review*.go` and the convention validator from Spec 0203.

## Expected

The review treats an archive that carries `qa_override: true` with an approval source and a reason as authorized. A Delivery Convention states this, and the validator dismisses such a finding, as it does other convention restatements. An archive without the override record is still flagged.
