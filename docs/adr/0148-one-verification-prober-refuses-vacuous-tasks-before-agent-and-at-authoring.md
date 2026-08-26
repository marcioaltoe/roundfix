---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# One Verification prober refuses vacuous Tasks before the Agent and at authoring

The Daemon runs a Task's `## Verification` commands once against the unchanged tree before creating the Agent Session and refuses the Task when any command already passes — a command that exits zero before the work exists proves nothing about the work (measured on Spec 0089, whose gate matched a string already in the file and settled a Task whose work was never done). The authoring-time check asks the identical question earlier, so the vacuous/failing/unknown classification lives in one extracted prober that both callers use: a second implementation would be free to disagree, and a checker that approves what the probe later refuses is the same defect one layer up. The Daemon keeps its Run bookkeeping around the shared loop; the authoring caller supplies a working directory and nothing else.

Consolidates ADR-0109 and ADR-0124 (2026-08-26); both are archived under docs/history/adr/.
