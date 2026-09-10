---
status: approved
granted: 2026-09-09
action: establish Spec-contained authorization and committed-provenance command execution authority
consuming: 0119-spec-contained-authorization
paths:
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/speccheck/constraints.go
  - internal/speccheck/constraints_characterization_test.go
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
  - skills/write-tasks/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/spec-routing.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - internal/speccheck/governed.go
  - internal/speccheck/governed_repocontract_test.go
  - internal/suiteguardcontract/regeneration.go
  - internal/suiteguardcontract/regeneration_test.go
  - internal/speccheck/mechanical_test.go
  - internal/speccheck/coherence.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0119

On 2026-09-09 the maintainer answered the explicit scope request with
"Aprovar os 21 caminhos", approving the complete proposal recorded here. The
earlier proposal state is superseded; the paths above are the exact bounded
set, and nothing outside them is authorized by this record.

## Confirmed session authority

On 2026-09-08 the maintainer requested Specs containing their authorizations
and the same requirement in canonical guidance. The confirmed unattended
delivery scope is implementation, the configured review policy, PR, and squash
merge after required checks pass. The configured policy permits codex, claude,
coderabbit or explicit none; none records intentional omission, and an enabled
provider's failure still blocks. Spec 0130 delivered the first grant in this
shape on 2026-09-09, so the record schema below is the maintainer's practice
made canonical rather than a new convention introduced here.

## Approved bounded mutation

1. Make the Spec-contained `_authorization.md` the canonical home for tooling
   authorization records in the Baseline modules and their rendered guides,
   and in the Roundfix-owned PRD, TechSpec and Task authoring instructions and
   templates in both the canonical `.agents/skills/` tree and the shipped
   `skills/` bundle.
2. Extend the constraint reader so grant validation follows the artifact's
   role rather than a date in its filename, and so a record cited from inside
   its own Spec resolves. Preserve every existing explicit legacy format.
3. Keep the historical Governed Path set monotonic while adding the owned
   shipped templates the current path predicate misses.

Preserve clause identities wherever an existing obligation can be extended
instead of replaced. Ordinary product-source implementation paths, including
the new `internal/spec` authorization reader, are not governed and belong to
the TechSpec's build order; discussing them here does not make them tooling
grants.

## Approved execution trust contract

On 2026-09-09 the maintainer selected committed provenance as the trust
contract for executing authored Verification, answering "Confiança por
procedência commitada".

- Authored commands may execute when the Spec artifacts that carry them are
  tracked in this repository and byte-identical to their committed bytes at
  the resolved revision.
- A Spec Root outside the repository's Git tree, an untracked or modified Spec
  artifact, or command text that differs from the committed bytes requires an
  execution approval recorded in that Spec's `_authorization.md`, naming the
  approved source revision, before any shell execution.
- Read-only `roundfix spec check` without `--run-verification` stays available
  for every source, trusted or not.
- This contract governs `spec check --run-verification`, Implement dispatch and
  Settle alike. It grants no network access, no credential access, no new
  sandbox, and no different execution privilege.

## Limits and commit order

- This record lands in `main` ancestry before the consuming squash delivery.
  A separate commit inside the consuming pull request would be flattened with
  the change and cannot preserve prior approval.
- Do not edit an upstream-managed skill or an archived grant.
- Do not infer approval from silence, preselection, a pending question, or a
  checker that does not yet validate the new record location.
- No paid API use, release, tag, deployment, destructive cleanup, or
  branch-policy exception is granted by this record.
- Corrective work stays inside the standing cap of two corrective Tasks per
  Spec; a third means the decomposition is wrong and the Run stops.
- Verification remains Daemon-owned. ADR-0014, ADR-0057, ADR-0096, ADR-0117,
  ADR-0130 and ADR-0149 remain operative and are not revised by this grant.

## Amendment — 2026-09-09 — widened scope and typed operations

The maintainer answered the second review's scope question with "Ampliar a 0119
para cobrir os três", widening this grant on the same day it was granted.

Two paths join the bounded set: `internal/suiteguardcontract/regeneration.go`
and its test. That file owns a second grant parser which reads only the legacy
and active Spec roots, so it cannot see an archived Spec's grant and cannot
preserve a multi-Spec consuming list. Leaving it outside the boundary would ship
two readers disagreeing about the same record.

The frontmatter now carries `operations`, the machine-readable permission list
Core Feature 5 requires. A free-text `action` cannot distinguish approval to
implement from approval to release, so an audit could not honor this record's
own no-release limit. The listed operations are exactly the maintainer's
confirmed through-merge delivery scope. `release`, `tag`, `deploy` and paid
consumption are absent, and absence is refusal rather than silence.

This amendment lands in its own commit before the work that consumes it, and it
widens nothing beyond the two named paths and the typed restatement of limits
this record already carried in prose.

## Amendment — 2026-09-10 — the ancestor audit's own test

`internal/speccheck/mechanical_test.go` joins the bounded set. The ancestor
grant audit lives in `internal/speccheck/mechanical.go`, which is ordinary
source, but its test file became governed when this Spec widened the Governed
Path set to every path an operative record has bounded — the archived Spec 0130
grant bounds exactly that file. The Spec's own widening made its own audit
Task's change out of bounds, and the terminal QA gate caught it as F-003.

This extension is limited to the purpose already granted: proving the ancestor
audit that this record's own action names. It adds no new action, no operation,
and no path beyond the one file.

`internal/speccheck/coherence.go` joins for the same reason at one remove. The
authorized edit to `internal/speccheck/constraints.go` changed the signature of
the constraint-row detector, and that file holds the call site the change made
stale. It is a consequent fix in the baseline's exact sense: necessary only
because the authorized change made something else stale, and one line long. The
grant covers the call site rather than leaving the authorized change unbuildable.

It lands in `main` on its own, before the implementation candidate that consumes
it. A commit on the consuming branch would be flattened by the squash into the
same target commit as the change it authorizes, which is the self-approval this
Spec exists to refuse. The rule applied to itself.

## Approval evidence — 2026-09-10

On 2026-09-10 the maintainer answered the explicit request for these two exact
paths with "Aprovar os dois caminhos", approving the amendment above and raising
the effective bounded set to twenty-five.

The decision was requested rather than inferred from the standing
purpose-bounded extension authority, because Core Feature 2 of this Spec
requires an amendment to carry a newly recorded maintainer decision before
dependent work resumes. Shipping the Spec that defines what an authorization
record must contain, with an amendment its own Core Feature 2 would reject,
would refute the contract being delivered. A standing delegation is authority to
ask narrowly, not a substitute for the recorded decision.

## Sanctioned regeneration

The repository-owned command resolves its generated outputs. This declaration
records the digest regeneration that follows the approved source edits above;
it adds no source paths. The `skills/` bundle regenerated by `make skills-sync`
is already enumerated path by path in the grant above.

```yaml
command: make baseline-digests
```
