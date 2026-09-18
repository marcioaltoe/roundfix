---
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
prd: _prd.md
created: 2026-09-17
---

# A repository says who reviews before the Pull Request — Technical Spec

## Executive Summary

One key, one resolution, one report. Configuration gains a pre-Pull-Request
reviewer beside the existing settings, resolved through the precedence every
other key already uses, and the diagnostic command states the result with its
source. Nothing invokes a provider.

The trade-off this design accepts is that the policy becomes readable before it
becomes enforceable: a Run does not yet consume it, and publication does not yet
require its evidence. That keeps this slice provable on its own and leaves the
reviewer session, where the real cost and the real failure modes live, to the
slice that can also record what a review produced.

## Project Constraints

- Identifier strategy: applicable — the provider names `codex`, `claude`, `coderabbit` and `none` are the vocabulary the agent guides already use, and this Spec reuses them without coining or renaming. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — resolving and reporting a policy reads local configuration only; no provider is installed, authenticated or invoked, and no network transport is opened. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the configuration precedence and the read-only support surfaces are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0002 applies: configuration is YAML in User Config and Project Config, and this Spec adds a key inside that contract rather than a new configuration mechanism.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the agent guide that states the obligation and the configuration this repository already ships.
- Tooling authority: applicable — the Doctor check is public CLI behavior, and the repository's hard rule ships the Roundfix skill update with it. Express maintainer authorization: granted 2026-09-18, recorded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved edit moves. The configuration schema, the check itself and the user guide are ordinary source that no authorization has bounded, and this Spec edits no `.roundfixrc.yml`, Baseline asset, or guide inside setup-context markers. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Configuration schema | `internal/config` overlays and typed Config | Carry the pre-Pull-Request reviewer and refuse an unsupported value. |
| Precedence resolution | `internal/config` layering | Resolve Project over User over the built-in default, reporting which layer answered. |
| Readiness report | `internal/cli` health checks | State the resolved provider and its source as one read-only check. |
| Public guidance | `docs/user-guide` | Document the key, its values and its precedence. |

No new package is proposed. The key joins the existing overlay shape, and the
report joins the existing check list.

## Implementation Design

### The key and its values

A top-level `pre_pr_review` section carries `provider`, whose value is one of
`codex`, `claude`, `coderabbit` or `none`. It is separate from `review_source`,
which keeps every key and behavior it has today. The overlay pointer is nil when
a layer says nothing, which is how the layering already distinguishes silence
from a written value.

### Resolution and its source

Layering already walks built-in defaults, then User Config, then Project Config.
The resolution returns both the effective provider and the layer that supplied
it, so a reader can tell a project decision from an inherited one. Silence at
every layer resolves to `codex`, the built-in default the guides assume.

### Refusal at load

Validation refuses an unsupported value with an error naming the key, the value
and the four supported values, in the shape the existing agent validation error
uses. The refusal happens while configuration loads, before any command reads
the policy.

### The report

The diagnostic command gains one check whose detail states the provider and its
source, and states explicit `none` as review disabled by configuration. The
check performs no probe and contacts nothing; it reports what configuration
says.

### Interfaces

```go
// Config gains one section; its zero value keeps today's behavior.
type PrePRReview struct {
    Provider string // codex | claude | coderabbit | none
    Source   string // project | user | default
}
```

### Data Models

No stored record, Run Database table or event payload changes. The policy is
configuration, read on demand.

## API Contracts

1. Configuration accepts `pre_pr_review.provider` with the values `codex`,
   `claude`, `coderabbit` and `none`; any other value fails the load with an
   error naming the key, the value and the supported set.
2. The diagnostic command reports one additional check whose detail names the
   resolved provider and its source, and names `none` as review disabled by
   configuration.
3. Every existing configuration key, default, validation message and diagnostic
   check keeps its current behavior.

## Coverage Map

- Goal 1 → Configuration schema.
- Goal 2 → Precedence resolution, Readiness report.
- Goal 3 → Configuration schema (refusal at load).
- Goal 4 → Readiness report (no probe).
- User Stories 1-2 → Precedence resolution.
- User Story 3 → Configuration schema; Readiness report.
- User Story 4 → Readiness report; API Contract 2.
- User Story 5 → Configuration schema; API Contract 1.
- Core Feature 1 → Configuration schema.
- Core Feature 2 → Precedence resolution.
- Core Feature 3 → Readiness report.
- Core Feature 4 → Configuration schema.
- Core Feature 5 → Readiness report; Testing Approach 3.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 1.
- API Contracts 1-3 → Configuration schema, Readiness report.

## Integration Points

- **Agent guides.** They state the obligation and the vocabulary this key reuses.
  They are delivered inside setup-context markers and are not edited here.
- **`review_source`.** Untouched. The two settings are documented side by side so
  a reader can tell the acts apart.
- **Spec 0126.** The reviewer sessions, the finding dispositions, the configured
  omission record and the publication evidence stay there.

## Testing Approach

1. **Schema and refusal, at the configuration unit seam.** Each supported value
   loads; an unsupported value fails with the named error. The refusal case
   fails on the tree as it stands today, where the key does not exist.
2. **Precedence, at the layering seam.** Project over User over default, with
   the reported source for each, including both-silent.
3. **Report, at the diagnostic command seam.** The check states provider and
   source, states `none` as disabled, and the command still contacts nothing —
   asserted by the absence of any provider invocation in the check's
   dependencies.
4. **Outside evidence.** The agent guide's clause and this repository's own
   configuration are read as they stand: the guide names the four values, and
   the configuration carries no such key, so the resolved policy is the built-in
   default. Where the guide cannot be read, the row records that and does not
   block.
5. **Repository gate.** The terminal QA Task records the Daemon's `make verify`
   result as a fact.

## Build Order

1. Configuration key, its overlay, its validation and its precedence, with unit
   tests (depends on: none).
2. Diagnostic report of the resolved policy and its source, with its tests
   (depends on: 1).
3. Public guidance for the key, its values and its precedence (depends on: 1, 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **Two review settings side by side.** `review_source` and `pre_pr_review` are
  easy to confuse. The guidance documents them together and names the act each
  one governs.
- **A readable policy nobody enforces yet.** Until Spec 0126's slice lands, the
  key records intent without changing delivery. That is stated in the guidance so
  a reader does not assume enforcement.

## Decisions

- **A separate key, not a mode of `review_source`.** Different act, different
  time, different failure modes.
- **Silence inherits.** An absent key never means `none`; a repository that wants
  no reviewer writes it.
- **The source travels with the value.** A resolved provider without its layer
  cannot answer "is this the project's decision or mine".

## Vocabulary Contract

No token is coined. The provider names and the words project, user and default
already describe existing configuration layers.
