---
status: accepted
knowledge: true
created_at: 2026-11-16T14:05:00Z
updated_at: 2026-11-16T14:05:00Z
deprecated_at: null
superseded_by: null
---

# A queue forecast is a likely time and a limit from the repository's own deliveries, or a named condition

The status view forecasts when a Spec starts and when its delivery merges as
two times, the median and the 85th percentile of what this repository's
delivered Specs took. It gives clock times only where a replay over history
held the limit for at least 80% of the Specs, and a named condition with a
duration everywhere else. The forecast exists so that the maintainer plans his
time from the status, and the condition exists so that no time is invented.

## Context

On 2026-11-15 the maintainer asked for a table of the scheduled Spec queue in
the status view, with a forecast start and end for the Specs not yet started.
On 2026-11-16 he decided its shape: the end is the delivery merge, a forecast
is a likely time and a limit, a wait Roundfix cannot date is a condition with
a duration, and the method is accepted only if the actual end falls inside the
limit for at least 80% of about 100 Specs.

The replay ran on 109 Specs delivered in this repository from 2026-10-17 to
2026-11-15, each forecast made only from the Specs merged before it. A Spec
took a median of 1.94 hours from its first Implement Run to its delivery merge
and 6.28 hours from its plan reaching the default branch. Roundfix records
when each Run began and ended and no time of a Delivery Queue item, and a new
queue replaces a finished one with its items.

## Decision

**The history.** A delivered Spec has a start, the start of its first
Implement Run in the Run Database, and an end, the committer time of the
commit that added its Archive Record to the default branch. The forecast reads
the Specs the Run Database still holds. It takes two durations from them: the
time from start to merge, and the start gap, from the start of one Spec to the
start of the next when that next Spec's plan was already on the default
branch.

**The pair.** The likely time is the median and the limit is the 85th
percentile, by nearest rank. A forecast needs 11 values. Once time has passed,
it uses only the values greater than that time and needs 5 of them.

**Where a time is shown.** A Spec in progress shows its actual start and an
end forecast from its start. A planned Spec one, two or three rows behind a
Spec in progress shows a start and an end: one start gap for each row, then
the time from start to merge, added as pairwise sums.

**Where a condition is shown.** Every other Spec that has not started shows
the wait it depends on as its start and its end as the pair of durations after
its start: a Spec with nothing in progress ahead of it, a Spec that waits on a
Run slot, a prerequisite or a park, and a Spec whose plan is not on the
default branch. An open Backlog Entry shows no forecast.

**The basis is stated.** Every report that shows a forecast says how many
delivered Specs it was computed from, over which dates, and that the limit is
the 85th percentile.

## Considered Options

- **Plan to merge as the basis of a Spec that has not started**, the basis the
  maintainer named for a Spec being authored. Refused by his own acceptance
  rule: forecast when the plan reached the default branch, the limit held for
  74 of 96 Specs (77.1%), and for 55 of 75 when another Spec was ahead. Shown
  that number on 2026-11-16, he accepted the time from start to merge instead:
  "Aceitar a troca (Recommended)".
- **A time for a planned Spec with nothing in progress ahead.** Refused: from
  the moment it became ready, the limit on its start held for 59 of 78. That
  wait ends when somebody queues the Spec.
- **Adding the medians and adding the limits down the queue.** Refused: the
  limit held and the likely time came earlier with each row, 25 of 63, 14 of
  39 and 6 of 19. Pairwise sums kept it between 42.9% and 52.6%.
- **A basis by the number of Tasks.** Not taken: it held for 81 of 98 against
  80 with one basis for every Spec, with a wider limit.
- **One time in each cell.** Refused by the maintainer, who chose the pair.
- **An earliest time where a wait cannot be dated.** Refused by the
  maintainer, who chose the condition.
- **Recording when a Delivery Queue item starts and reaches each stage.** Not
  taken now: the forecast holds on what is already recorded.

## Consequences

- With a queue that holds one Spec at a time, most rows that have not started
  show a condition. The table then tells the order, the status and the notes,
  and a duration for each row.
- A repository with fewer than 11 delivered Specs in its Run Database shows no
  time, and says how many it has.
- The history is one machine's and as long as Run Retention keeps it. A change
  in how the repository delivers moves the forecast only as new Specs replace
  old ones.
- The limit is missed about one time in seven by design.
- A Spec being authored shows a condition and the durations from start to
  merge, not the plan-to-merge figures, which the maintainer accepted.
