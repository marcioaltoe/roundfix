---
status: accepted
created_at: 2026-10-08T00:00:00Z
updated_at: 2026-10-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec is authored against the QA rerun classes, and a retirement writes reduced history

The Baseline audit of 2026-10-08 (`docs/references/2026-10-08-baseline-audit.md`)
measured Specs 0225–0248. Five QA reruns came from Surface Transcript lines the
implementing test never asserted, while the code was correct. Four of the last
eleven Specs needed a QA Archive Override, and each override rested on an
outside-evidence row the QA gate could not reach: an operator-only log, a
lookup the Spec's authorization forbade, or a host the sandbox denied. Stale
wording left in a second guide failed QA or CI four times. Tests that read the
host's keys or bound a long macOS socket path parked or delayed Runs, Specs
authored in parallel needed hand fixes after a sibling merged, and a corrective
Task needed a hand-added `needs` edge before `roundfix reopen` saw it. On
2026-10-08 the maintainer asked to "Atualizar e verificar o baseline de regras
para que tenhamos o máximo de performance e resultado com o roundfix e
congruente com as mudanças".

**Each rerun class becomes an authoring rule.** The Spec workflow requires an
outside-evidence source the QA gate can check without the network before the
Run starts, a test that asserts every line of a Surface Transcript and its exit
code, and `requires` for a dependency on an unmerged Spec. The autonomous loop
re-checks a Spec authored beside another after rebasing, adds a corrective
Task to the gate's `needs` before reopening it, sweeps the old wording of a
changed message or rule and characterizes behavior before changing it. A new
Go clause keeps tests hermetic. The authoring skills say the same, and the
repository-owned guide carries the rules that hold only for Roundfix: the
suite guard, the `### QA settlement` section, clause retention, skill and
module records, derived files and the measured Governed Path set.

**A retirement writes reduced history directly.** ADR-0248 reduces the existing
history through a sanitize, and "Retired Review Artifacts and handoffs have no
reader, so they are removed". A new retirement now does the same at once
(operator decision, 2026-10-08, after the maintainer's "Não quero nada no
histórico que não seja relevante para o secondbrain"). A retired Finding or
Backlog Entry enters its history family as a Reduced History Entry: its front
matter carries the disposition, and it keeps its title, its first paragraph
and a final line naming a commit that already holds its full text. A finished
Review Artifact or a confirmed handoff is deleted and never enters history.
The sanitize keeps its contract for what earlier retirements left.

**Reworded clauses keep their identity.** ADR-0257: "A clause whose wording
changes keeps its identity". Every rule above extends a clause in place,
except the Go clause, which is new because no Source Baseline holds the Go
module.

## Considered options

- New clause identities for the transcript, prerequisite and loop rules.
  Rejected: each needs a Source Baseline row the regenerator never creates.
- Leaving retirement unchanged and sanitizing after each one. Rejected: a full
  entry written into history after the History Full Tag is a unit that tag
  does not hold, so the update refuses it (ADR-0254), and the text would sit in
  history until someone reduced it by hand.
- Changing the History Relocation code that still moves a finished legacy
  Review Artifact into `docs/history/reviews/`. Deferred: the update's Pending
  History removes it in the same plan, and the change is code, not guidance.

## Consequences

An adopter's next `roundfix baseline update` rewrites the Spec workflow, the
autonomous-work, the Go and the docs-layout guides. Every reworded clause is
retained, and the Go clause is added. The audit estimates about 9–11 minutes
of wall time saved per Spec from whole-transcript assertion and most of the
remaining overrides from reachable evidence. Those numbers rest on the audit's
counts, not on a measurement this ADR makes.
