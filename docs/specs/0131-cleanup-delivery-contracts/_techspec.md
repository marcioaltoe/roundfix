---
spec: 0131-cleanup-delivery-contracts
prd: _prd.md
created: 2026-09-09
---

# Cleanup delivery contracts — Technical Spec

## Executive Summary

Use the current declaration as the test input and obtain complete history in CI.
The trade-off is a larger checkout in exchange for reproducible historical
assertions. Keep the repair inside existing test and workflow seams.

## Project Constraints

- Identifier strategy: applicable — preserve existing clause IDs, Spec slugs and repository-relative paths; no new domain identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no application interface or credential access; the existing read-only CI checkout retains its credentials policy. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0153 defines explicit pre-PR review policy, including none; the fixture must track the approved clause without restoring the superseded Watch requirement. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on 2026-09-09: "Considere autorizada as mudanças e ajustes necessários para finalizar o trabalho nessa branch e fazemos o pull request, squash merge e sync main." The bounded repair is recorded in [_authorization.md](_authorization.md); bounded files: `internal/docscontract/corpus_test.go`, `.github/workflows/ci-verify.yml`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

The existing documentation contract test copies the shipped Markdown clause,
repository guide and Baseline module into a temporary repository. The public
Spec checker compares those declarations. CI runs the existing Make targets.
No new runtime component or dependency is needed.

## Implementation Design

### Interfaces

Find the current delivery-order clause in the authoritative module using its
existing clause identity or an equally bounded declared marker. Derive the
mutable order instead of copying the current sentence as a new test constant.
Preserve all three independent mutations and the existing assertions. Refuse
missing, empty or ambiguous targets. Keep JSON valid and do not mutate the
production parser or invent a second detector.

Set the existing actions/checkout input to fetch complete history. Keep
`persist-credentials: false`, permissions, action versions and check commands.
The historical tests must still fail when required evidence is absent.

### Data Models

No data model or glossary change. The final Task records the unchanged Vocabulary
Contract and tests the actual public checker through its existing fixture.

## Coverage Map

- Goal 1 and Core Feature 1 → current-clause mutation → task_01.
- Goal 2 and Core Feature 2 → complete CI history → task_01.
- OE-1 → historical Git read and current-candidate CI → task_01 and delivery checks.

## Integration Points

GitHub actions/checkout supplies repository history; the Go test invokes the
existing Spec checker. No new external service or privileged write is introduced.

## Testing Approach

The Task's declared Verification runs the existing three divergent-source cases
and verifies the full-history checkout declaration. The Supervisor then runs
complete `make verify verify-docs`, reviews the final candidate independently,
and requires GitHub's checks on that exact head before the requested squash merge.
A separate Agent QA Task is declined: this is maintenance of the test fixture and
its CI input, with no changed product journey; a second agent gate would repeat
these same checks. Retain a factual verification report for archival and clearly
record this QA decision instead of claiming a separate QA session ran.

## Build Order

1. Repair the two bounded files in one Task and record its focused evidence.
2. Complete repository verification and independent review (depends on: 1), record
   the result, archive this maintenance Spec and finish the existing PR.

## Risks & Considerations

Full history adds fetch cost but leaves the existing time budget intact. A weak
mutation could stop testing divergence; retain negative assertions and fail-fast
handling of an invalid source. Do not repeat the formatter mismatch from Spec
0130: on this host use Go and gofmt 1.26.7 explicitly for repository checks.

## Vocabulary Contract

No glossary term is added, changed or retired. The final Task confirms this.

## Decisions and research evidence

Secondbrain's 2026-09-09 capture on historical authorization and regeneration was
read: retain strict evidence boundaries rather than restore deleted docs or
bypass a guard. Current ADR-0153 supplies the review policy, not the retired test
sentence. Exa read the primary [actions/checkout documentation](https://github.com/actions/checkout),
which states that its default fetch is one commit and `fetch-depth: 0` obtains
all history. That directly supports the CI input repair. Exa also read the
[Go JSON documentation](https://pkg.go.dev/encoding/json#Unmarshal); its bounded
excerpt establishes the existing JSON decoding surface, not the correctness of
the proposed fixture. No JSON package migration is proposed. Deriving a mutation
from the existing clause is a local design choice, validated by the negative
cases and independent review.

Secondbrain source: `inbox/secondbrain/2026-09-09-autorizacoes-historicas-e-propriedade-da-regeneracao.md`.
