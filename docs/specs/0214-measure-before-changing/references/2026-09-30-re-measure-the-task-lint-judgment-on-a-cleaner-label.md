---
type: feat
status: promoted
created: 2026-09-30
spec: 0214-measure-before-changing
reason: null
---

# Re-measure the Task-lint judgment once Runs record a cleaner label

## Opportunity

The 2026-09-30 Jev measurement asked whether a Task file's acceptance criteria are all observable, and whether that predicts a Verification repair. The signal ran in the expected direction but was weak: AUROC 0.62 over 300 Tasks, 0.67 for backend Tasks only. The label was noisy, because it was read from whether the Agent-written `## Result` mentions Verification feedback, which reflects implementation difficulty as much as Task quality.

## Value

A Task whose acceptance an observer cannot check costs a repair round after the Agent has already worked. A judgment that flags it during `write-tasks` would move that cost to authoring, where it is one edit. The measurement cannot say whether the judgment is worth adopting until the label is clean.

## Shape

Once the Run Database records Verification attempts per Task, re-score the same acceptance question against that count on the pinned model version, with a random-Task baseline and a confidence interval. Adopt it as a third advisory judgment of `roundfix spec judge` only if it separates repaired from unrepaired Tasks well enough to act on one Task at a time. The finding that measured it is `docs/history/findings/2026-09-30-jev-judgments-measured-against-the-spec-archive.md`.
