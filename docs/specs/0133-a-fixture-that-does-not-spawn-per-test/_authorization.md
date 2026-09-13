---
status: approved
granted: 2026-09-13
action: narrow the absent-record rule to governed mutation and remove the per-test fixture spawn cost
consuming: 0133-a-fixture-that-does-not-spawn-per-test
paths:
  - internal/baseline/assets/modules/core.json
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0133

On 2026-09-13 the maintainer answered the explicit request for these three exact
paths with "Aprovar os três caminhos", approving the bounded scope above.

## Why a canonical path is unavoidable

The canonical clause in `internal/baseline/assets/modules/core.json` reads
"a proposed, absent, contradictory, or withdrawn record grants nothing", and the
implementation obeyed it literally: `internal/cli/implement.go` requires the
`implement` operation from every Spec, and a record with an empty `paths` list
is invalid, so a Spec touching no governed path can neither write a valid record
nor be implemented. Spec 0131 ran Clean in exactly that condition before the
operation enforcement landed.

Narrowing the behavior without narrowing the clause would leave the canonical
rule stating the opposite of what the code does. The two rendered artifacts are
regenerated from the module by the public Baseline update; they are not
hand-edited.

## Approved bounded mutation

State that an absent, proposed, contradictory or withdrawn record grants no
**governed mutation**, and does not by itself refuse implementation, commit or
push. Authority over protected tooling and authority to do the work become two
statements again, which is what tooling authority always meant.

The rest of the repair is ungoverned source: `internal/authorization`,
`internal/cli` and `internal/daemon`. The fixture cost repair touches only
`internal/daemon/task_engine_test.go`, also ungoverned.

## Limits

- No action, operation or path beyond the three above.
- The record stays required, and refusal stays fail-closed, for every change to
  a Governed Path. This narrows which actions an absent record blocks, never
  what a governed change needs.
- No paid API use, release, tag, deployment, or branch-policy exception.
- Verification remains Daemon-owned; Task status remains Daemon-written.

## Limits and commit order

This record lands in `main` ancestry before the consuming squash delivery. A
commit on the consuming branch would be flattened by the squash into the same
target commit as the change it authorizes.

## Sanctioned regeneration

The repository-owned command resolves its generated outputs. This declaration
records the digest regeneration that follows the approved module edit; it adds
no source paths.

```yaml
command: make baseline-digests
```
