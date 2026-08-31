# kickoff is settled: implement-spec remains the loop (2026-08-31)

The 2026-08-26 triage listed **`kickoff` versus `implement-spec`** as a decision
owed before either could be packaged. This records the answer so a later reader
does not reopen it.

## The decision

`implement-spec` remains the Roundfix-owned autonomous loop. The proposed
`kickoff` skill is not adopted.

Confirmed by the maintainer on 2026-08-31. The skill was never committed to this
repository, so there is no removal commit to cite: it existed only in the
maintainer's working environment and was deleted there before the confirmation.
That is why this record carries the reasoning rather than pointing at a diff —
the reasoning is the only durable artifact the decision leaves behind.

## Why the question existed

The proposed skill declared `implement-spec` a rival loop that contradicts
`docs/agents/autonomous-work.md`, and published on its own authority. Both
cannot be the loop at once: had `kickoff` shipped alongside, the fleet would
have received a loop plus a note telling readers not to use it.

## What was kept

The half of `kickoff` that carried measured value already shipped. Its
session-cutoff script — whose own header documented the naive-comparison bug it
existed to avoid — became **Spec 0117, a Run Window the Preflight owns**,
released in v0.8.0. The discipline became a Run-creation precondition rather
than a script a supervisor must remember to call, which is the whole reason the
triage put that half first.

## What was declined, and why it was not portable anyway

The rest was not adoptable as written, independently of the rivalry:

- it cited `ADR-0041`, which does not exist in this repository — ADR numbers are
  per repository, so the citation could not resolve;
- it carried a consuming repository's domain vocabulary and spend ceiling;
- it declared a review policy that ADR-0118 places in a typed Baseline decision
  rather than in a skill;
- it was written in Portuguese, while every Roundfix-owned skill ships in
  English.

Adopting it would have meant rewriting all four before the rivalry question was
even reached.

## Consequence for the queue

**No `kickoff` implementation shipped, and none will.** Nothing from the skill
entered this repository as a skill, a command, or a rule. The one thing that did
ship — the Run Window of Spec 0117 — was authored here from the measured problem
its script documented, not adapted from its code.

The queue item therefore closes with no implementation owed. The remaining Specs
are `0105` and `0097`, plus the defects captured for Triage in the Secondbrain
inbox.
