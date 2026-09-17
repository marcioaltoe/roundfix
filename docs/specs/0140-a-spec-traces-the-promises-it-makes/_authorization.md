---
status: approved
granted: 2026-09-17
action: extend the PRD, TechSpec and Task authoring skills so a Spec declares its Success Metrics and API Contracts as traceable units and names them where Tasks are written
consuming: 0140-a-spec-traces-the-promises-it-makes
paths:
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
  - skills/write-tasks/SKILL.md
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints_characterization_test.go
  - internal/docscontract/testdata/corpus-golden.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0140

On 2026-09-17 the maintainer chose to carve this Spec out of Spec 0129 and to
make traced promises its first slice. Asked the same day to approve this record
with its exact paths, operations and regeneration, the maintainer answered
"Aprovar como proposto". That approved the bounded scope below. The approval
lands in the delivery target's ancestry before the commits that consume it.

## Why a governed path is unavoidable

The rule this Spec writes is an authoring rule: what a PRD and a TechSpec must
declare, and in what shape. That text lives in `write-prd` and `write-techspec`,
both Roundfix-owned Skills, so changing it needs an express grant. Every Run
executes from the canonical copy under `.agents/skills/`, and the binary embeds
the distributed mirror under `skills/`, so both copies of each file are bounded
here.

Each skill's template is bounded beside its skill because the template is the
artifact an author copies; a rule stated only in the skill body would be
contradicted by the template that ships with it.

The Task authoring skill is bounded for the same reason. Its derivation step
tells the author to map every PRD user story and Core Feature, and its Task
checklist requires only user stories in References, so a rule the checker
enforces would have no guidance where Tasks are actually written. The maintainer
widened the grant to those two files, and their mirrors, on 2026-09-17 after
pre-PR review raised the gap.

Three more paths are governed because an authorization has already bounded them,
which ADR-0130 keeps governed for good: the checker's stage table, the
constraint characterization whose expectation the new unit kind moves, and the
corpus golden the new codes are counted in. The first QA gate of this Spec
refused on exactly those, and the maintainer widened the grant to them on
2026-09-17. The rest of the checker is ordinary source that no authorization has
bounded, so it needs no grant.

The Task template's distributed mirror is not a governed path — the mirror
pattern covers `SKILL.md` and agent files only — so it was removed from this
record. Bounding an ungoverned path fails the repository's own bounded-path
contract, and narrowing grants no authority.

## Approved bounded mutation

Edit the canonical skills and their templates so that:

- the PRD's Success Metrics section is a numbered list, or the single entry
  `None.` with the reason none applies;
- the TechSpec's API Contracts section is a numbered list, or the single entry
  `None.` with the reason none applies;
- each numbered Success Metric appears in the TechSpec Coverage Map;
- each numbered Success Metric and API Contract is named in some Task's
  References;
- each skill's report step refuses to recommend the next pipeline step while its
  own stage reports a promise finding;
- the Task authoring skill maps every declared promise, requires each in some
  Task's References, and its Task template's References example names one;
- the stage table places each coined code at the stage that can establish it;
- the constraint characterization keeps its subject while the new unit kind
  moves its expectation;
- the corpus golden counts the two coined codes.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`. The mirror paths are already bounded above. If any derived
pin changes as a result, rewrite it only with `make baseline-digests`.

## Sanctioned regeneration

The repository-owned command resolves its own generated outputs. This
declaration records the regeneration proposed above and adds no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No action, operation or path beyond those above.
- No edit to the qa-gate skill, the QA gate matrix, the QA Report or any verdict
  rule.
- No edit to the Task Type contract, the QA decision rules, or any Verification
  rule in the Task authoring skill; the only change there is the promise
  coverage the rule adds.
- No new or renumbered consistency code that changes an existing code's meaning,
  and no change to the severity of an existing finding.
- No Baseline module, guide inside setup-context markers, linter, analyzer or
  Verification configuration edit.
- No change to any archived Spec, which stays byte-identical.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
