---
type: fix
status: promoted
created: 2026-10-26
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
reason: null
---

# An outside-evidence row that names what the authorization forbids does not qualify

## Symptom

A QA gate that cannot look up an outside source because the Spec's own
`_authorization.md` forbids it records the row as blocked and names the
prohibition. When the row names the network in the exact words
`blocked (environment: network prohibited by Spec authorization)`, the partial
qualifies. When it names what the authorization actually forbids in other
words, the partial does not qualify: settlement and `roundfix qa-report
accept` refuse it with
`rows_blocked_environment is <n>, 1 outside the pre-PR Pull Request row and
the network-denied and network-prohibited outside-evidence rows; expected 0
outside them`, the Delivery Queue parks the item `qa-environment-partial`,
and archiving needs a maintainer's QA Archive Override. Nothing in the work
failed.

Two Specs of this repository hit it:

- Spec 0317, QA report `qa-report-2026-10-25.md` at `ef1a4bab`, row
  `| QA-08 | blocked (environment: upstream checkout and network prohibited by Spec authorization) | outside-evidence row |`.
  The report itself names the mismatch: "The official
  `spec.ReadQAReport` parses this report as partial with environment=3 and
  declared=1. However, `spec.QAReportEligibility` refuses it".
- Spec 0309, QA report `qa-report-2026-10-26.md` at `06a9b7ce`, row
  `| R07 | blocked (environment: another repository prohibited by Spec authorization) | outside-evidence row |`.
  Rows R06 and R08 of the same report, written with the exact network phrase,
  were accepted.

Spec 0258 (QA report `qa-report-2026-10-13.md` at `078869bf`, rows 12–15)
wrote the exact network phrase four times before any form qualified; ADR-0276
made that phrase qualify, and the 0258 rows stay the control case.

## Where

`internal/spec/qa.go`: `QANetworkProhibitedStatus` is compared with `==` in
`qaReportResultsCounts`, so only the one network sentence counts. The QA gate,
`roundfix settle`, Daemon settlement, `roundfix qa-report accept`,
`roundfix archive` and the Delivery Queue's park decision all read that
count through `QAReportEligibility`. The `qa-gate` skill tells the agent to
write the network phrase "exactly", and agents keep writing the honest
variant when the authorization forbids something other than the network.

## Expected

An outside-evidence row blocked because the Spec's own authorization forbids
its lookup qualifies when its status is
`blocked (environment: <subject> prohibited by Spec authorization)`, with a
non-empty, bounded subject that names what the authorization forbids and
holds no parenthesis, and its provenance names the outside-evidence row. The
network phrase is that form with the subject `network`. Every other
qualification rule keeps its bytes: provenance, the frontmatter counts,
skipped rows, finding rows, declared rows and the refusal messages.

## Evidence

The three QA reports above, read at their source revisions on 2026-10-26.
Both Archive Records carry the override:
`docs/history/specs/0317-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`
(`disposition: qa-override`, approval of 2026-10-25 ~12:40 BRT, reason "QA-08
… was recorded as 'upstream checkout and network prohibited by Spec
authorization', outside the accepted wording") and
`docs/history/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff.md`
(`disposition: qa-override`, approval of 2026-10-26 ~07:55 BRT, reason "R07
… was recorded as 'another repository prohibited by Spec authorization',
outside the accepted wording"). In both the operator verified the blocked
source by hand before the maintainer approved.
