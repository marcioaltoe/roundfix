---
spec: 0142-an-absent-ref-is-named-absent
prd: _prd.md
created: 2026-09-17
---

# An absent ref is named absent — Technical Spec

## Executive Summary

Two changes on one path. The branch resolver stops collapsing "no ref" and
"several refs" into one count, and reports each with its own message. The Run
Branch set classifier stops failing as a whole when the recorded target branch
does not exist: it preserves each candidate in the set with a reason naming the
absent branch.

The trade-off this design accepts is that an absent target branch still proves
nothing about delivery, so nothing new is released. The change buys a truthful
reason and the classification of every other candidate, not recovered debris.

## Project Constraints

- Identifier strategy: applicable — Run identifiers, Run Branch names and the recorded target branch keep their spelling; refs are read by fully qualified name and no identifier is coined. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the classifier reads the local Git repository and the Run Database only, with no credential store or network transport. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the reconciliation outcome vocabulary and its preservation rules are governed by accepted decisions this design preserves. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0057 applies: the Daemon owns Run and Task status, and this design changes no status.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the adopted finding.
  ADR-0127 applies: process residue is a readiness observation, so naming an absent ref creates no Run and settles nothing.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/worktree`, which is not a Governed Path, and its test files are ordinary too. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Branch resolver | `internal/worktree` ref resolution | Answer absent, ambiguous, or resolved, each with its own message. |
| Run Branch set classifier | `internal/worktree` reconciliation | Preserve an absent-target set truthfully instead of failing it. |

No new package, file or directory is proposed.

## Implementation Design

### Absent and ambiguous are separate answers

The resolver already reads a branch through its fully qualified name. Its guard
counts candidate refs across the namespaces a short name could reach and refuses
when the count is not one, which makes zero and two indistinguishable. The guard
now reports:

- **absent** when no candidate ref matches, using the existence check already
  defined beside the resolver;
- **ambiguous** when more than one matches, keeping today's message and its
  branch name.

Callers that only need a head keep receiving an error; what changes is which
sentence a reader gets, and that the two cases can be told apart by their error.

### An absent target branch preserves its set

The classifier resolves the recorded target branch before it walks the
candidates, and a failure there returns an error for the whole set. When the
target branch is absent, the set is instead classified with each candidate
preserved, carrying a reason that names the absent branch. A target branch that
exists takes exactly the path it takes today: the QA Report comparison, the
content proof, the release and supersession rules, and the existing preservation
reasons.

An ambiguous target branch keeps failing the set, because the classifier cannot
know which ref the Run meant.

### Interfaces

No exported signature changes. The resolver's error carries which case it hit,
so the classifier can preserve on absence and keep failing on ambiguity.

### Data Models

No entity, schema, stored record or event payload changes.

### API Contracts

1. A reconciliation candidate whose recorded target branch has no local ref is
   reported as preserved, with a reason naming that branch as absent.
2. A refusal that reads `short ref is ambiguous` is emitted only when the short
   name matches more than one ref.
3. Every other reconciliation outcome, reason and summary count keeps its
   current text and meaning.

## Coverage Map

- Goal 1 → Branch resolver.
- Goal 2 → Run Branch set classifier.
- Goal 3 → Run Branch set classifier (present-target path).
- User Story 1 → Branch resolver; API Contract 2.
- User Story 2 → Run Branch set classifier.
- User Story 3 → Run Branch set classifier; API Contract 1.
- Core Feature 1 → Branch resolver.
- Core Feature 2 → Run Branch set classifier.
- Core Feature 3 → Run Branch set classifier; API Contract 3.
- Success Metric 1 → Testing Approach 1 and 3.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- API Contracts 1-3 → Branch resolver, Run Branch set classifier.

## Integration Points

- **Reconcile command.** It renders the reasons this design changes; its output
  contract is otherwise untouched.
- **Carry-forward.** It shares the classifier, so an absent target branch stops
  refusing its set there too, with the same reason.
- **Spec 0125.** Repository identity across linked worktrees and the delivered
  content evidence rule stay there.

## Testing Approach

1. **Resolver, at its existing fixture seam.** A fixture repository with no such
   ref answers absent; one with a tag and a branch of the same short name answers
   ambiguous; a plain branch resolves. The absent case fails on the tree as it
   stands today, which reports ambiguity.
2. **Classifier, at its existing fixture seam.** A Run set whose recorded target
   branch was deleted preserves every candidate with the absent-branch reason,
   while a second set whose target branch exists reaches its usual outcome in the
   same run.
3. **Outside evidence.** The adopted finding records fourteen candidates refused
   on this machine for a branch Git does not list. The replay reproduces that
   shape in a fixture and shows the new reason; where the recorded observation
   cannot be read, the row records that reason and does not block.
4. **Repository gate.** The terminal QA Task records the Daemon's `make verify`
   result as a fact.

## Build Order

1. Separate absent from ambiguous in the resolver, with its unit tests (depends
   on: none).
2. Preserve an absent-target set in the classifier, with its tests (depends on:
   1).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

- **A reason other tooling reads.** Any check pinned to the ambiguity sentence
  for an absent ref sees the new one. The repository's own tests are the only
  known readers, and they move with this change.
- **Preservation is unchanged.** Nothing new is released, so a mistake here
  cannot delete a Run Branch; the worst case is a wrong sentence, which is what
  this Spec exists to fix.

## Decisions

- **Two answers, not one count.** Absence and ambiguity have different repairs,
  so they need different messages.
- **Preserve on absence, fail on ambiguity.** An absent branch is a fact the
  classifier can state; an ambiguous one leaves it unable to know which ref the
  Run meant.
- **Nothing new is released.** Proving delivery without the target branch is Spec
  0125's evidence rule, deliberately left out.

## Vocabulary Contract

No token is coined. "Absent" and "ambiguous" describe existing ref states in
prose and name no new field, code or command.
