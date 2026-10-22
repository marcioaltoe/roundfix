---
type: fix
status: open
created: 2026-10-22
spec: null
reason: null
---

# A QA precondition refusal report records no audited commit, so a copied-in old refusal passes for the current result

## Symptom

On 2026-10-21 oraculum's Spec 0087 ran its QA gate four times. Only the first report, `qa-report-2026-10-21-01.md`, refused the precondition (`SC-WAVE-COLLISION`, correctly: at `57a7c55` the corrective `task_05` still needed `task_01`). The next three reports (`-02` to `-04`, `build:` `8736693`, `7472ad34` and `75a9dd0`) passed the precondition and failed on row R13 for two documentation findings. Each later Run carried the earlier reports into its own branch, and the refusal report records no `build:` (audited commit), so the operator read the old refusal as the current result four times, reported a gate defect to Roundfix, and archived with `--qa-override`. Spec 0297's authoring proved the gate correct on every head.

## Where

The QA report a gate writes when its precondition refuses (the qa-gate skill's refusal path and whatever template or Daemon code writes or checks that report), compared with the full report, which records `build:`.

## Expected

A refusal report records the commit it audited, as a full report does, so a reader (and any Roundfix command that picks the current QA report) can tell a copied-in older report from the current Run's result. Reports already written keep their bytes.

## Evidence

Spec 0297 (`docs/history/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check.md`, characterization) read oraculum's history read-only on 2026-10-22. Left out of 0297 because changing the report format is a decision its briefing did not cover.
