---
type: refactor # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The canonical method lives in rendered guides instead of skills

## Opportunity

The CONTEXT-driven method is canonical: it should read the same in every
repository that adopts it. Today it is delivered by rendering prose into
`docs/agents/`, which each repository then personalises in the same files.

Measured on 2026-08-09 in this repository:

| | |
| --- | ---: |
| `docs/agents/` | 13 files, 1,450 lines |
| Managed regions rendered from the catalog | 13 |
| Repository rules interleaved with them | 41 |
| Catalog producing them | 16 modules, 103 KB of JSON |

So 103 KB of catalog JSON exists to render 1,450 lines of prose into files that
also hold 41 local rules, and most of a `baseline update` plan's output is the
retention evidence proving it did not eat them. A `baseline update` run the same
day reported 9 file changes, 4 digest mismatches, 41 retention entries and 2
nested-carrier warnings — for prose.

The skill channel already carries canonical content without any of that: 14
skills ship with the CLI, `make skills-sync` keeps copies identical, and
`skills/baseline_skill_contract_test.go` fails when a required clause is edited
away. On 2026-08-09 those contracts caught five real regressions, including a
gate edit that removed the Project Constraint audit.

## Value

Separating canonical from local removes the retention machinery entirely for
these files: the baseline stops writing them, so it has nothing to preserve.
Updates become a version bump rather than a diff-and-merge over managed regions,
and the digest-mismatch class disappears — the same run found a clause,
"Task Type-selected Agent Session", present in the repository and in no catalog
module, which applying the plan would have removed.

## Shape

Non-binding, and the shape matters more than the move.

The obvious split is canonical to a skill and personalisation to `AGENTS.md`.
The risk it carries is that **a guide is always-on and a skill is on demand**:
a clause in `docs/agents/` reached from `CLAUDE.md` governs every task, while a
clause in a skill governs only a dispatched one. Today mandatory dispatch is
itself prose in `skill-dispatch.md` and nothing verifies it, so the move would
trade a weak guarantee for a weaker one unless dispatch becomes checkable.

The division worth settling is therefore by enforcement rather than by
canonicity:

- a clause a mechanical check already decides is documentation, and can live
  wherever it is cheapest to read;
- a clause that governs every task and has no check stays always-on and minimal;
- a clause that governs one kind of work belongs to that skill;
- repository personalisation belongs in `AGENTS.md`.

That ordering matters because 2026-08-09 proved prose alone does not hold: a
finding filed on 2026-08-06 predicted Spec 0090's F-001 exactly, and nothing
stopped it until the rule became a checker detector.

Also worth settling: whether each part of the method gets its own skill, or one
skill carries the method with a `references/` directory. The second keeps one
dispatch to enforce; the first keeps each part loadable on its own.

Related: `docs/specs/0085-what-an-agent-reads-before-it-decides` already owns
what an Agent reads, so this either extends it or supersedes part of it.
