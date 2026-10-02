# Task acceptance measurement

The question was written anew for this Spec, so these results are comparable
with the 2026-09-30 measurement only in direction and scale. The harness
selected OpenRouter first and reported transport `openrouter`; every answered
response reported the pinned model `jev-1.13`.

## Run accounting

| measure | value |
| --- | ---: |
| calls | 237 |
| input tokens | 131,258 |
| cost | US$0.005512836 |
| answered | 237 |
| excluded | 7 |
| repaired | 31 |
| bootstrap draws discarded | 0 of 2,000 (0%) |

The seven exclusions were one missing Task file and six Tasks rejected by the
judge language gate. The record contains no repaired-answer substitutions;
`repaired` is the Run Database label derived from at least one failed
Verification verdict.

## Statistics

| statistic | value |
| --- | ---: |
| repaired base rate | 0.130802 |
| AUROC | 0.466881 |
| 95% cluster-bootstrap interval | [0.358034, 0.574594] |
| random baseline AUROC | 0.536172 |
| acceptance-criteria-length AUROC | 0.640698 |
| flagged at `noul < 0.5` | 38 |
| precision | 0.052632 |
| recall | 0.064516 |

The 2026-09-30 published figure was AUROC 0.62 [0.55, 0.70] over 300 Tasks
with a 23% repaired base rate. This is a new question and a new measurement,
so it does not reproduce that sample. Its signal is lower than the prior
figure, its interval is wholly below the rule’s 0.70 lower-bound threshold,
and its length baseline is stronger than the judgment score.

## Decision

The fixed `_techspec.md` → Decision rules, rule 2, first requires at least
150 answered Tasks, at least 20 repaired Tasks, an unblocked record, and no
more than the allowed discarded share. Those gates are met: 237 answered,
31 repaired, status `measured`, and 0% discarded. The adoption branch then
requires an interval lower bound of at least 0.70, at least 10 flagged Tasks,
and precision at least twice the base rate. Although 38 Tasks were flagged,
the lower bound is 0.358034 and precision 0.052632 is below twice the base
rate 0.261603, so the rule does not fire. No Backlog Entry is created.

Verdict: do not adopt
