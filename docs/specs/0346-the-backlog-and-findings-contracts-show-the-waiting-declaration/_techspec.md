---
spec: 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
prd: _prd.md
created: 2026-11-04
---

# The Backlog and Findings contracts show the Waiting Declaration — Technical Spec

## Executive Summary

Two clauses of the `context-workflow` Baseline module gain two commented
frontmatter lines and one sentence each; nothing in Go production code
changes. The module records its next version, `make baseline-digests`
regenerates the derived catalog, golden and plan files, and a Managed
Refresh rewrites this repository's docs layout guide and Setup Manifest. The
user guide's published Findings template follows the generated contract,
because a repository contract test holds the two equal. The trade-off this
design accepts is that the keys are YAML comments rather than placeholder
keys: a copy made unchanged declares no wait and nothing has to validate a
placeholder, at the cost that an agent must uncomment two lines on purpose.
One new Baseline test proves the bytes, the YAML effect of the comments and
the Source Baseline retention of both clauses.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; both clauses keep their identities
  `clause.context.backlog-01-operational-contract` and
  `clause.context.findings-01-frontmatter`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the Baseline renders local guides. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0353: "Change the Baseline
  contracts in the same change: deferred to the open Backlog Entry"; this
  Spec is that change, and it keeps the ADR's rule that "The command that
  reads a declaration validates it." No new ADR. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AX answers of 2026-11-04 quoted in the authorization record; no
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0346-the-backlog-and-findings-contracts-show-the-waiting-declaration/_authorization.md`;
  bounded files: `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`.

## System Architecture

The change lives entirely in existing seams.

- `internal/baseline/assets/modules/context-workflow.json`, rule
  `rule.context.docs-layout` (version 18 at `6a0ed904`): the guidance of the
  two clauses. The clauses keep `id`, `enforcement: mandatory` and no
  `replaces`; the rule's `version` rises by one and the module's top-level
  `version` is written by the Module Version Record command (29 at
  `6a0ed904`; the next version is whatever the record holds plus one when
  the Task runs).
- The derived Baseline files that `make baseline-digests` owns: the
  Standard TypeScript formatter golden `docs/agents/docs-layout.md`, the
  golden digest in the Standard TypeScript profile, the three catalog files
  under `internal/baseline/testdata/` and the four plan characterization
  goldens.
- This repository's managed copies, written by the Managed Refresh
  (`roundfix baseline update --repo . --no-skills --no-history --yes`):
  `docs/agents/docs-layout.md` and the catalog and guide digests of
  `docs/agents/setup-context.json`.
- `docs/user-guide/context-driven-development.md`: the block between
  `baseline-findings-template:start` and `:end`, held equal to the generated
  Findings fence by `TestGuidanceCompositionDocumentation`
  (`internal/docscontract/publicdocs_test.go`), plus one explanatory
  paragraph after the lifecycle paragraph that follows it.
- `CONTEXT.md`: the **Waiting Declaration** entry.
- The readers stay as they are: `internal/cli/follow_up_inventory.go`
  decodes frontmatter with `gopkg.in/yaml.v3`, so a commented key is
  absent; `classifySourceClauseTransition` (`internal/baseline/plan.go`)
  classifies a previous clause whose identity and enforcement remain as
  `retained`.

## Implementation Design

### Interfaces

The new test file sketches its seams over the embedded catalog:

```go
// internal/baseline/waiting_declaration_clause_test.go
const waitingOnContractLine = "# waiting_on: <who or what it waits on, in one line> # optional, with since"
const sinceContractLine = "# since: YYYY-MM-DD # optional, the date the wait began"

var waitingDeclarationClauses = []struct{ id, sentence string }{ /* API Contracts 1-2 */ }

func TestTheOperationalContractsShowTheWaitingDeclarationCommented(t *testing.T)
func TestAnAdopterKeepsTheWaitingContractClauses(t *testing.T)
func firstFrontmatterBlock(t *testing.T, guidance string) string
```

### Data Models

No data model changes. The Backlog and Findings frontmatter keys
`waiting_on` and `since` already exist (ADR-0353); the templates only show
them, commented.

### API Contracts

1. API Contract: the Backlog Operational Contract. In the guidance of
   `clause.context.backlog-01-operational-contract`, the copyable
   frontmatter gains, right after the `reason: null # required for terminal
   closure without a consuming Spec` line and before the closing `---`, the
   two lines

   ```yaml
   # waiting_on: <who or what it waits on, in one line> # optional, with since
   # since: YYYY-MM-DD # optional, the date the wait began
   ```

   and, after that fenced block and before `For` followed by the `feat`
   type heading, one paragraph:
   "An `open` entry that waits on someone or something outside the
   repository may uncomment `waiting_on`, naming who or what in one line,
   and `since`, the `YYYY-MM-DD` date the wait began." Every other byte of
   the clause stays.
2. API Contract: the Findings Operational Contract. In the guidance of
   `clause.context.findings-01-frontmatter`, the copyable frontmatter gains
   the same two lines right after `updated_at: YYYY-MM-DD` and before the
   closing `---`, and after the closing fence of the template one paragraph:
   "A `pending` or `partial` Finding that waits on someone or something
   outside the repository may uncomment `waiting_on`, naming who or what in
   one line, and `since`, the `YYYY-MM-DD` date the wait began." Every other
   byte of the clause stays.
3. API Contract: the update preview. On a repository whose Setup Manifest
   records the previous catalog, `roundfix baseline update --no-skills
   --no-history --format text` reports `File changes: 2` with
   `update docs/agents/docs-layout.md` and
   `update docs/agents/setup-context.json`, and `Retention evidence: 0`
   (a Managed Refresh). A plan that accounts a Source Baseline transition
   records both clauses with disposition `retained`, targets equal to their
   own identities and enforcement `mandatory`. No new flag, output line or
   exit code.

### Surface Transcripts

None. No command surface changes; API Contract 3 names lines an existing
command already prints for any guide refresh.

## Coverage Map

- Goal 1, Story 1, Core Feature 1, Success Metric 1 → API Contracts (item 1),
  Testing Approach (`TestTheOperationalContractsShowTheWaitingDeclarationCommented`).
- Goal 1, Story 2, Core Feature 2, Success Metric 1 → API Contracts (item 2),
  Testing Approach (the same test).
- Goal 2, Core Feature 3, Success Metric 2 → Interfaces
  (`firstFrontmatterBlock`), Testing Approach (the YAML subtests).
- Goal 3, Success Metric 4 → System Architecture (clause identity and
  `classifySourceClauseTransition`), Testing Approach
  (`TestAnAdopterKeepsTheWaitingContractClauses`).
- Goal 4, Story 3, Core Feature 4, Success Metric 3 → System Architecture
  (derived files, Managed Refresh, user guide template), Build Order (step 1).
- Success Metric 5 → API Contracts (item 3), Build Order (step 1, second
  Managed Refresh).
- Core Feature 5 → System Architecture (`CONTEXT.md`), Build Order (step 1).

Deviations from the PRD: None.

## Absence Evidence

1. `waiting_declaration_clause_test.go` (the clause test) — command: `grep -rn 'waiting_on' internal/baseline internal/speccheck docs/agents`; result: no match at `6a0ed904`; the only readers and tests of the keys are `internal/cli/follow_up.go`, `internal/cli/follow_up_inventory.go` and `internal/cli/follow_up_waiting_test.go`, which read entries, not the Baseline.

## ADR Alternatives

- ADR-0353: the ADR rejected adding a Spec Check finding for a malformed
  declaration and deferred the Baseline change to the adopted Backlog Entry;
  this design takes only the deferred Baseline change and still adds no
  check or validation, so the rejected alternative stays rejected.

## Integration Points

None outside the repository. Adopters receive the change through their own
`roundfix baseline update` after the release.

## Testing Approach

- `internal/baseline` (new file `waiting_declaration_clause_test.go`):
  reads the embedded catalog as `release_clause_test.go` does, asserts
  API Contracts 1-2 in the module, the Standard TypeScript golden and this
  repository's guide, parses each copyable frontmatter with `yaml.v3`
  (no waiting key as copied; both keys once uncommented), and runs
  `classifySourceClauseTransition` from
  `baseline.standard-typescript-monorepo-0.0.1` as
  `cross_project_continuity_test.go` does.
- Existing seams that must keep passing unchanged:
  `TestTheRetirementClausesCarryTheirText`,
  `TestLifecycleClausesCarryTheScopedWording`,
  `TestFindingsOperationalContract`, `TestNoTwoBaselineClausesShareText`,
  `TestBaselineClauseForceIsCharacterized`, the catalog, characterization,
  formatter and module-version tests, and the docscontract package
  (`TestGuidanceCompositionDocumentation`, `TestTheSkillCoverageMapIsCurrent`,
  `TestTheSetupManifestIsCurrent`).
- Sabotage proven in a scratch copy on 2026-11-04: uncommenting
  `waiting_on` in the Backlog contract fails the YAML subtest of
  `clause.context.backlog-01-operational-contract`; dropping the subject of
  the Findings sentence fails its sentence assertion; leaving this
  repository's guide at its previous bytes fails the `repository guide`
  subtest on both lines and both sentences.

## Build Order

1. Edit the two clauses and the rule version; add the clause test; record
   the module version; run `make baseline-digests`; run the Managed Refresh
   twice; update the user guide's Findings template and paragraph, the
   `CONTEXT.md` entry and the Skill Coverage Map review; re-record the
   Behavior Surface Record last.
2. Final QA gate (depends on: 1).

## Risks & Considerations

- Rebase order: this Spec delivers last in cycle AX. Another cycle AX Spec
  that changes a Baseline module moves the catalog digests, the plan goldens
  and the module version record. The Task regenerates them by command on its
  own base, so a rebase conflict in those files is resolved by rerunning the
  commands, never by hand.
- The formatter golden is regenerated by `make baseline-digests`; an
  unexpected rewrap of the sentence would show as a Verification failure of
  the whitespace-normalized checks, not as silent drift.
- The new frontmatter lines start with `#`. Inside the Markdown fence they
  are YAML comments, not headings; the Markdown formatter leaves fenced
  content unchanged (proved by the prototype golden).

## Glossary

- changes: **Waiting Declaration** — the entry names the Baseline's Backlog and Findings Operational Contracts that show it commented

## Decisions

- Commented lines over placeholder keys or `null` values: a copied
  placeholder would be a malformed declaration, and `null` would read as a
  value to fill in.
- The sentence sits after the Backlog frontmatter block, before the
  per-type bodies, and after the whole Findings template, so the user
  guide's template block stays a pure copy of the fence.
- One implementation Task: the module, its derived files, the managed
  guides and the user guide template must change in one settlement for the
  repository gate to stay green.
- No new ADR; ADR-0353 already decided the keys and their validation.
