# Final QA journey matrix

Every row starts pending and requires evidence from a fresh disposable Fluxus
clone created during the final `qa-gate`. Evidence from an earlier Task or QA
run does not satisfy either row.

| Journey | Actor | Entry point | Expected observable | Independent confirmation | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| Fluxus greenfield | Fluxus maintainer | Public Baseline Command against a fresh disposable Fluxus greenfield clone | Explicit project decisions produce the reviewed Plan, apply, formatter and repository Verification recommendation | Fresh Plan has zero managed delta; audit and empty reapply remain verified | Capture command transcript and disposable clone identity during final `qa-gate` | pending |
| Fluxus update | Fluxus maintainer | Public Baseline Command against a separate fresh disposable Fluxus update clone | Keep-defaults reuses the persisted Better Auth HTTP reason and reaches a Plan without manual repair | Formatter and repository Verification preserve managed bytes; audit and empty reapply report zero managed delta | Capture command transcript and disposable clone identity during final `qa-gate` | pending |
