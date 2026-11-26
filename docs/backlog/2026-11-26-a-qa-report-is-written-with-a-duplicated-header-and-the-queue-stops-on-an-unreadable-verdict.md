---
type: fix
status: open
created: 2026-11-26
spec: null
reason: null
---

# A QA report is written with a duplicated header and the queue stops on an unreadable verdict

## Symptom

In the delivery of Vortex Spec 0222, the QA of one Run wrote `qa-report-2026-11-26-01.md` with two sets of `status`, `build`, `surfaces` and `user_flow_executables` in the front matter: the first closed with 8 executables, the second in-progress with `[]`. The sections `## Findings`, `## Blocked and skipped` and `## Coverage` appear twice. The queue stopped as "QA verdict unreadable".

## Open question

Which step writes the header twice: the QA gate skill's report authoring, or the runtime that seeds or amends the report. Not measured here: whether the runtime validates the report shape before it reads the verdict.

## Source

Cross-project message `01a2168c-7b8e-7580-9d8e-8f08971b2f10` from the Vortex project, 2026-11-26. Nature: reported (observed in a Run audited with Roundfix 0.153.3, no response expected). Recorded by triage on 2026-11-26.
