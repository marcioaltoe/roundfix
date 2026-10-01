---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A repository records its frontend layout, and an unrecorded layout follows the suggestion

The frontend guide made the systems layout mandatory: feature code organized by
domain system, each system behind one public boundary. React itself takes no
position on how files go into folders, and the planned argus repository
organizes its renderer by feature with shared stores and components. A
repository that uses another layout could only contradict the Baseline.

The frontend layout is now a repository-owned decision, `frontend.layout`,
with the values `systems` and `repository-defined`. `systems` is the suggestion.
The decision is optional: while a repository records none, the Baseline states
the suggestion and renders the systems clauses as before, so an existing
adopter, a repository-owned profile that does not select the decision and an
automation that does not answer it keep their current guidance. A recorded
value is never overwritten. A clause may name the decision value it applies
under; the systems clauses apply under `systems` and the new clause that binds
the repository's own layout applies under `repository-defined`. When a
repository records a value that turns a previously managed clause off, the
upgrade accounts for that clause as a reasoned rejection naming the recorded
decision, instead of failing closed as unaccounted.

Two alternatives were rejected. A required decision would stop every existing
adopter's update until it answered, and would make every caller that passes
decisions name one more. Moving the systems clauses into a module that the
decision activates would move them to a new managed block in every adopter's
frontend guide and drop them silently from repository-owned profiles that
copied the old module list.

## Consequences

The systems clauses keep their identifiers, enforcement and bytes, so the
Source Baseline retains them and no adopter's rendered clauses change until it
records another layout. The interactive first adoption asks the question with
`systems` preselected. A repository that records `repository-defined` states
its layout in its own rules; the Baseline checks no directory against it.
