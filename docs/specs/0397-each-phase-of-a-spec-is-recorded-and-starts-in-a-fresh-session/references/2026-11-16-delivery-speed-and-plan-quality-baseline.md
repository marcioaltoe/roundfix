# Baseline: how fast Specs are delivered and how well their plans hold, before Spec 0397 (recorded by the author, 2026-11-16, 07:05–07:25 BRT)

The maintainer added two conditions to the success criterion on 2026-11-16: no or almost no loss of speed, and better artifacts for implementations. This file records what both looked like before the Spec, from three sources read without changing them: the Run Database of this repository on the Mac mini, the Archive Records under `docs/history/specs/`, and the Git history of the default branch at `2caa1d17`. The commands ran in the repository's checkout with a binary built from that commit (Roundfix 0.145.1).

## 1. Runs: `roundfix runs list --state terminal --limit 0 --format plain`

The Run Database holds 177 terminal Implement Runs for 112 Specs, created from 2026-10-16 17:52 UTC to 2026-11-15 11:49 UTC: 110 ended Clean, 66 Unresolved and 1 Failed.

| Specs whose first Run started | Specs | First Run ended Clean | Runs per Spec | Duration of a Clean Run: median (first to third quartile) |
| --- | --- | --- | --- | --- |
| In the 30 days the database holds | 112 | 67 (60%) | 1.58 (177 Runs) | 47 min (24 to 66), 110 Runs |
| On or after 2026-11-01 | 47 | 32 (68%) | 1.38 (65 Runs) | 61 min (38 to 85), 47 Runs |
| On or after 2026-11-09 | 19 | 15 (79%) | 1.26 (24 Runs) | 56 min (37 to 71), 19 Runs |

Of the 112 Specs, 68 (61%) needed one Run and 14 needed three or more; 45 (40%) had at least one Run that did not end Clean.

## 2. Failed Verification and corrective Tasks: `roundfix runs causes --since 2026-10-16 --until 2026-11-16 --format json`

The report covers the same 177 Runs and 112 archived Specs, each with its source revision available. It lists 3 items, all failed Verification attempts: two of class `repository-convention` on one Task of Spec 0386 and one of class `implementation-defect` on one Task of Spec 0393. It lists no corrective Task, which the command defines as a Task Graph node numbered after the QA Task.

Per-Task verdicts exist for 19 Specs and 43 Tasks, the Specs whose first Run started on or after 2026-11-09 (0376 to 0396), because this repository keeps the Run Event Journal for 168 hours: 43 verdicts passed and 3 failed (6.5% of 46), with 2 feedback rounds and 1 Task that needed more than one Run.

## 3. Archive Records and the Git history

The 58 Archive Records created from 2026-11-01 to 2026-11-14 are Specs 0337 to 0396: 36 of type `fix`, 19 of type `feat` and 3 without a type.

- Time from the commit that adds the Spec's PRD to the default branch to the commit that adds its Archive Record: median 6.0 hours, first quartile 3.7, third quartile 10.9, 90th percentile 19.2, longest 24.4. For the 30 Specs delivered from 2026-11-08 to 2026-11-15 the median is 7.4 hours (4.2 to 12.4).
- Disposition: 56 are `partial` and 2 are `pass`. With 97% of the records in one value, the disposition does not tell a good plan from a poor one.
- `unproven` entries: 23 in all, in 22 of the 58 records (38%).
- 41 Specs were archived on the day their record says they were created and 17 on the next day.

## 4. The band used for a median of 20 Specs

A median of 20 values moves by chance. Its standard error is taken as `1.2533 × (interquartile range ÷ 1.349) ÷ √20`, the normal approximation. For the plan-to-merge time the interquartile range is 7.2 hours and the standard error 1.5 hours. For the duration of a Clean Run the interquartile range is 47 minutes and the standard error 10 minutes. A share measured on 20 Specs moves in steps of 5 points, and a share near 70% to 80% has a standard deviation of about 9 to 10 points.

## Limits

- The Run Database is one machine's. From 2026-11-09 to 2026-11-11 this repository's Runs ran on another machine, so 47 of the 58 Specs created since 2026-11-01 have Runs here; the 11 without are Specs 0365 to 0375.
- A Run's duration is not divided by its Tasks, so a week of larger Specs reads as a slower week.
- A first Run that does not end Clean is not always the plan's fault: a runtime failure or a busy machine ends a Run Unresolved too.
- The first-Run rate was already rising inside the window, from 60% to 79%, before this Spec exists.
- Findings of the strict check at plan time are corrected before a plan is committed and nothing records them, so they have no baseline.
