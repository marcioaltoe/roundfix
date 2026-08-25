---
granted: 2026-08-12
action: Mutate the bounded authoring skills, Baseline catalog assets, setup-owned guides, and the verification Makefile that the eight corrective Specs 0095, 0096, 0098, 0104, 0105, 0106, 0107 and 0108 require.
consuming: 0095-a-verification-that-ran-before-anyone-believed-it, 0096-a-failure-the-agent-can-read, 0098-a-hook-that-cannot-outrank-the-gate, 0104-a-gate-that-cannot-certify-its-own-cache, 0105-the-gates-own-economics, 0106-a-decision-that-reaches-every-artifact, 0107-the-authoring-rules-the-guides-do-not-carry, 0108-what-an-agent-loads-to-answer-one-question
paths:
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/implement-task/SKILL.md
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/decisions.json
  - docs/agents/autonomous-work.md
  - docs/agents/docs-layout.md
  - docs/agents/spec-routing.md
  - docs/agents/agent-instructions.md
  - docs/agents/secondbrain.md
  - docs/agents/skill-dispatch.md
  - docs/agents/specific-repository.md
  - docs/agents/setup-context.json
  - Makefile
---

# Tooling authorization — the authoring and Baseline corrections (2026-08-12)

On 2026-08-12 the maintainer reviewed the eight Specs in
`docs/design/2026-08-12-the-spec-set-the-evidence-asks-for.md` that cannot be
written without a tooling grant, was told that a per-Spec record is the tighter
shape and that an umbrella grant is looser than the bounded-file rule intends,
and chose the umbrella:

> Um grant guarda-chuva para as 8

This record is that grant. It is deliberately broader than the fourteen records
before it, and the per-Spec breakdown below exists so that the changed-file audit
each Task runs still has a boundary narrower than the union above.

## What this covers, per Spec

**0095 — a Verification that ran before anyone believed it.** The authoring
contract gains the rule that a Verification command passes only by exiting zero,
with the working forms recorded beside the vacuity rule it belongs next to.
Bounded to `.agents/skills/write-tasks/SKILL.md`.

**0096 — a failure the Agent can read.** The authoring contract gains the
sanctioned exits when the corrective-Task ceiling is reached. Bounded to
`.agents/skills/write-tasks/SKILL.md`.

**0098 — a hook that cannot outrank the gate.** The autonomous-work module gains
the invariant that a commit hook may never be stricter than the authoritative
Verification, and the recovery ladder for a Run that dies without one. Bounded to
`internal/baseline/assets/modules/autonomous-work.json` and the guide it renders,
`docs/agents/autonomous-work.md`.

**0104 — a gate that cannot certify its own cache.** A make target that clears
the cache the gate actually uses, and the guidance that currently points at a
command clearing a different one. Bounded to `Makefile` and
`docs/agents/specific-repository.md`.

**0105 — the gate's own economics.** The QA gate skill applies the
equivalent-evidence path to the Pull Request row by default and takes ownership of
the QA Task's Verification; the task-authoring skill stops the review-resolution
Agent from rewriting an authored Verification. Bounded to
`.agents/skills/qa-gate/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, and
`.agents/skills/implement-task/SKILL.md`.

**0106 — a decision that reaches every artifact.** Exclusion plans a deletion,
carrier discovery excludes the Baseline's own fixtures, the HTTP Contract keeps
its exceptions through a change, the inert profile default is settled, and the
Setup Manifest's recorded catalog digest is compared against the embedded
catalog. Bounded to `internal/baseline/assets/decisions.json`,
`internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
`internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
and `docs/agents/setup-context.json`.

**0107 — the authoring rules the guides do not carry.** The measured authoring
rules from three fleet repositories reach the modules that render the guides, and
the Secondbrain guide gains a production clause with a trigger. Bounded to
`internal/baseline/assets/modules/autonomous-work.json`,
`internal/baseline/assets/modules/context-workflow.json`,
`internal/baseline/assets/modules/spec-workflow.json`,
`internal/baseline/assets/modules/secondbrain.json`,
`internal/baseline/assets/modules/core.json`, and the guides they render:
`docs/agents/autonomous-work.md`, `docs/agents/docs-layout.md`,
`docs/agents/spec-routing.md`, `docs/agents/agent-instructions.md`,
`docs/agents/secondbrain.md`.

**0108 — what an Agent loads to answer one question.** The canonical method
divides by enforcement rather than by canonicity, and the operational skill splits
its nine subjects into references. Bounded to every
`.agents/skills/*/SKILL.md` listed above, `docs/agents/skill-dispatch.md`, and the
modules that declare dispatch.

## Bounded by purpose

Each Spec's Tasks may change only the paths its own section names, plus their own
Task files. The union in the frontmatter exists so a checker can enumerate the
grant; it is not a licence for one Spec to reach another's paths. A Task whose
changed-file audit finds a path outside its Spec's section fails, exactly as it
would under a per-Spec record.

This grant covers no other Baseline clause, no other module, no other guide, no
other skill, no CI workflow, no dependency manifest or lockfile, and no version
pin. Derived Baseline pins and skill mirrors rewritten by `make baseline-digests`
and `make skills-sync` are sanctioned fallout under ADR-0081, not separate
targets.

## Sanctioned regeneration

The clause above states this in prose; the block below states it in the form the
changed-path audit reads. The detector may treat these exact paths only as
outputs of the declared command. This declaration does not make any output a
freely editable bounded file; the QA gate still verifies that its bytes match the
canonical sources under `.agents/skills/`.

```yaml
command: make skills-sync
outputs:
  - skills/archive-spec/SKILL.md
  - skills/implement-task/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/roundfix/agents/openai.yaml
  - skills/write-idea/SKILL.md
  - skills/write-idea/references/idea-template.md
  - skills/write-idea/references/opportunity-scan.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-tasks/SKILL.md
  - skills/write-tasks/references/task-template.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
```

Only the eight skills this grant bounds have their mirrors declared. A mirror
whose canonical source is outside the frontmatter above is not covered, so a
Task cannot reach an unauthorized skill by regenerating the bundle.

```yaml
command: make baseline-digests
outputs:
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/autonomous-work.md
  - internal/baseline/assets/source-baselines/index.json
  - internal/baseline/preservation_test.go
  - internal/baseline/testdata/catalog.diagnostics.golden.json
  - internal/baseline/testdata/catalog.digest
  - internal/baseline/testdata/catalog.normalized.json
  - internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json
  - internal/baseline/testdata/plan-characterization/clean-adoption.golden.json
  - internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json
  - internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json
```

These outputs are regenerated whenever the Baseline modules or their rendered
guides change; Spec 0098's documentation of the hook-strictness invariant in
`internal/baseline/assets/modules/autonomous-work.json` triggers this
regeneration. The outputs remain read-only; the QA gate verifies that the
regenerated bytes match the canonical digest.

## Chronology

This record lands as its own commit, before any commit that consumes it. A
prerequisite fix repairing something already red may land before either. A
consequent fix made necessary by an authorized change lands after it, never folded
into it.

## What the maintainer was told

That a per-Spec record keeps each grant drawn around one purpose, which is what
the fourteen existing records do and what the QA gate's changed-file audit reads;
and that a single record spanning eight Specs widens every one of them to the
union unless the per-Spec sections above are treated as binding. They are treated
as binding here for that reason.
