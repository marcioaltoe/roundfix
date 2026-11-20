---
status: accepted
knowledge: true
created_at: 2026-11-17T18:45:00Z
updated_at: 2026-11-17T18:45:00Z
deprecated_at: null
superseded_by: null
---

# The Corrections Report prints counts and sets two models side by side only above twenty Specs each

The Corrections Report shows every rate as a count over a count beside the
Delivery Type and the size in Tasks, prints no percentage, and sets two authoring
models side by side only for the same Authoring Phase and Delivery Type and only
when each has at least 20 delivered Specs. Below that it prints
`insufficient sample` with the counts it has, and it never names a better
model. The report exists to inform the choice of an authoring model, and a
small or mixed sample read as a verdict would misinform it.

## Context

The study of 2026-11-17 found 4 delivered Specs authored on Fable and 14 on
Opus. Counted per Spec, the Fable Specs had 1.50 corrections each against
1.14; counted per Task, 0.40 against 0.48: the order flips with the
denominator. The four Fable Specs were all features with an Idea, and 11 of
the 14 Opus Specs were fixes, so the model and the Delivery Type moved together.

The maintainer chose on 2026-11-17 "Completo e honesto (Recommended)", a view
that shows the number of Specs, their type and their size beside every rate
and "amostra insuficiente" below a minimum; "Comparar o que vier", with no
alternation of models; and "20 Specs (Recommended)" as the minimum.

Published practice points the same way. The CONSORT statement asks authors to
"State the results in absolute numbers when feasible (eg, 10/20, not 50%)"
(Moher, Schulz and Altman, 2001). The data presentation standards of the US
National Center for Health Statistics suppress a proportion whose denominator
is under 30 and print the reason (NCHS, 2017). Telling 10% from 50% with 80%
power needs 23 per group, and 10% from 30% needs 69 (Kim, 2016).

## Decision

**Counts over counts.** A rate is printed as its numerator over its
denominator, for Specs and for Tasks. The report prints no percentage, and
its JSON carries the two integers.

**Groups.** Plan defects are counted for each Authoring Phase, Delivery Type,
model and reasoning effort, with the Specs and Tasks that group authored, and
a model's total across its efforts is shown beside them. A plan defect counts
for the phase that wrote the artifact it names.

**The comparison.** For one Authoring Phase and one Delivery Type, the report
sets side by side the groups, and the model totals, that each hold at least
20 delivered Specs. With fewer than two of them it prints
`insufficient sample`, each group's count and the 20 needed.

**What it states.** The report states its window, the Roundfix Home it read,
the Specs whose model is `unknown`, the Corrections that are `undeclared` and
the Specs it left out of a count, with the reason.

**No verdict.** The report ranks nothing: its rows follow the Spec number and
its groups a fixed order of phase, type and name. It prints no word that
names a better or a worse model.

**Beside, never counted.** Dismissed review findings and failed Verification
attempts by Cause Class are shown beside the Corrections of each Spec and
enter no count of Corrections.

## Considered Options

- **A percentage with a confidence interval.** Rejected: the interval of a
  group of this size is wider than any difference worth acting on, and the
  two counts say the same without inviting the reading.
- **A minimum of 30**, as the health statistics standard sets. Not taken: the
  maintainer chose 20 as a floor for showing two counts side by side, not as
  a test, and each count is always printed.
- **No minimum.** Rejected by the maintainer.
- **A statistical test and its result.** Rejected: the Specs are not assigned
  to a model at random, so a test would answer a question the data cannot
  ask.
- **Roundfix alternating the authoring model between Specs of one type.** Not
  taken: the maintainer chose to compare the Specs as they come. It is the
  design that can answer the question, and it stays a later step.
- **Hiding a group below the minimum.** Rejected: the counts of every group
  are shown, and only the comparison is withheld.

## Consequences

- The comparison may read `insufficient sample` for weeks, or for as long as
  each phase and Delivery Type is authored by one model.
- Twenty is a floor, not evidence: at 20 Specs a share moves in steps of five
  points.
- A reader who wants a percentage computes it from the two integers.
- The report reaches as far back as this machine's records, so a wider
  comparison needs facts that travel with the Spec.
