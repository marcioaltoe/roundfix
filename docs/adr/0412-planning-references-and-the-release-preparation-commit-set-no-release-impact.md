---
status: accepted
knowledge: true
amended_by: [ADR-0438, ADR-0440]
created_at: 2026-11-16T23:30:00Z
updated_at: 2026-11-16T23:30:00Z
deprecated_at: null
superseded_by: null
---

# Planning references and the release-preparation commit set no release impact

The Release Plan Command asks for a classification by hand for two kinds of
commit that never change what a release ships. Both now classify themselves:
`docs/references/` is planning evidence, and a commit that changes only the
release paths of a Spec the range delivers is a release-preparation commit.
A Release Ceiling still never answers a classification by hand.

## Context

The Oraculum reported on 2026-11-16 (ai-memory message `01a1e695`) that every
release of a Spec with a Release Ceiling ended in
`manual_classification_required`. The Spec's own planning commit adds a file
under `docs/references/`, as `docs/agents/spec-routing.md` requires for outside
evidence. The plan counts only `docs/specs/`, `docs/backlog/`, `docs/adr/`,
`docs/findings/` and `docs/handoffs/` as planning evidence (ADR-0349). The
release-preparation commit, which changes the changelog and the version files,
was also listed as a commit to classify by hand.

This repository met both. Over its 51 releases from `v0.125.0` to `v0.147.0`,
the plan before the preparation commit stopped for a classification twice: once
for a `docs/references/` file of a planning commit (`v0.143.4`), and once for
the root `.ai-memory.toml` (`v0.142.1`). Every plan that included the
preparation commit stopped. In 46 of those ranges, every path the preparation
commit changes is a `release_paths` entry of a Spec the range delivers.

## Decision

**`docs/references/` is planning evidence.** It joins the five planning
directories of ADR-0349 as a maintenance-only path. The documentation layout
gives it external pointers and durable project reference documents, and no
release ships them. The rule stays a rule on paths: a commit that also changes
any other released path still crosses the boundary.

**A release-preparation commit sets no impact.** The plan reads the
authorization records of the Specs the range delivers before it classifies,
whenever the range holds a commit that has no automatic impact and crosses the
boundary, or when the plan needs the approval question answered. A commit is a
release-preparation commit when its automatic impact is `none`, it changes at
least one path, and every path it changes is a `release_paths` entry of a
record that grants a ceiling. Such a commit is inside the maintenance boundary,
and the change evidence says why (`releasePreparation`). A commit with a `feat`,
`fix`, `perf` or breaking impact keeps its impact even when it touches only
release paths.

**The ceiling's authority is unchanged.** It answers `approval_required` under
the five conditions of ADR-0349 and never answers a breaking plan or a
classification by hand. The two rules above remove the reason a plan needed a
classification; they never answer one.

**Extended by [ADR-0438](0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md) and [ADR-0440](0440-a-lesson-learned-is-written-only-when-a-test-holds-and-leaves-when-a-rule-states-it.md) (2026-11-19).** `.secondbrain-export` is maintenance-only and `docs/lessons/` is planning evidence.

## Considered Options

- Letting the Release Ceiling answer when the classification by hand covers
  only `docs:` commits with no impact was rejected: it is the alternative
  ADR-0349 rejected by the maintainer's decision ("a classification by hand is
  a new decision"), and it would widen the authority a ceiling grants.
- Counting a `docs/references/` file as planning evidence only when a Spec of
  the same commit cites it was rejected: the classifier would read file
  contents at every commit, and a reference changed alone, such as a Coverage
  Review in `skill-coverage.json`, ships nothing either.
- Counting `docs/references/` only in a `docs:` commit was rejected: no other
  planning directory depends on the commit type, and a commit with any released
  path already crosses the boundary.
- Recognising the preparation commit by its subject, such as
  `chore: release vX.Y.Z`, was rejected: adopters write `docs:` or other
  subjects, and a subject is not a declaration of which files a release may
  change.
- Adding a list of derived files a release may regenerate was rejected:
  `release_paths` already names every path the release may change, so a
  lockfile or a generated catalog the release rewrites belongs there.

## Consequences

- A planning change that commits outside evidence under `docs/references/` no
  longer stops the plan, so a Release Ceiling can answer a range that delivers
  that Spec.
- A plan rerun after the release-preparation commit, or on the head of a
  release Pull Request, reaches the same decision as before that commit when
  the delivered records list every path it changes.
- A preparation commit that changes a path no delivered record lists, or a
  range that delivers no Spec with a ceiling, still asks for a classification.
- The plan now reads the delivered records in ranges where it did not before.
  A Git failure while reading them fails the plan, as it already did for the
  approval answer.
