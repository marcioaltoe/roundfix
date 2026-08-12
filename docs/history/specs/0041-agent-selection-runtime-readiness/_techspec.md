---
spec: 0041-agent-selection-runtime-readiness
prd: _prd.md
created: 2026-07-17
---

# Agent Selection Runtime Readiness — Technical Spec

## Executive Summary

Replace Roundfix's name-only adapter discovery and fixed
`set reasoning_effort` assumption with one shared capability-driven Agent
Selection prover. It first verifies the effective adapter command, official
package lineage, and supported version. A disposable ACP Session then exposes
the adapter's model and configuration controls. The prover converts the
canonical requested tuple into an adapter assignment plan, applies it, reads
the effective advertised state, and returns a structured proof or a classified
failure. Setup, profile configuration/validation, Doctor, and operational
preflight consume that same result.

The key boundary is intent versus transport. Roundfix continues to store and
report `codex / gpt-5.6-sol / high`; it does not store an ACP adapter's variant
ID as the user-facing model. The assignment plan may use an independent
reasoning option or an advertised model variant, but proof succeeds only when
that transport representation maps unambiguously back to the requested tuple.
An empty effort remains an explicit model-managed selection and is never used
as recovery for rejected `high`.

## System Architecture

- **`internal/agent`** owns advertised ACP capabilities, adapter
  identity/version proof, canonical-to-adapter assignment planning,
  disposable-session proof, effective-state comparison, classified failures,
  and cleanup.
  `ACPXRunner.Probe` delegates exact selection proof to this boundary instead
  of trusting any PATH executable named `codex-acp` and always issuing
  `set reasoning_effort`.
- **`internal/config`** keeps canonical Agent Selection Profiles and official
  identifiers. Built-ins and `DefaultConfigYAML` move the general Codex
  fallback from Terra/max to GPT-5.5/xhigh; recommendations remain advisory.
- **`internal/cli`** exposes one readiness coordinator used by setup,
  `profiles configure`, `profiles validate`, Doctor, and operational
  preflight. CLI code owns output, exit codes, no-mutation ordering, and the
  all-or-none override grammar.
- **`internal/daemon`** continues to consume canonical proven selections and
  the configured Fallback Chain. It does not understand ACP transport variant
  IDs and retains ADR-0050's notification-before-activation boundary.
- **Docs and Roundfix Skill** replace bare `--agent` examples with profile-led
  commands and document complete one-Run overrides. Canonical and embedded
  skill copies remain synchronized.
- **Spec 0036** appends Repository Skill Set readiness after the profile-aware
  Doctor surface and reuses its dependency seams without duplicating model
  proof.

## Discovery Evidence and Corrected Assumptions

The 2026-07-17 dogfood environment first supplied these observations:

```text
Codex CLI:                 0.144.5
ACPX:                      0.12.0
@zed-industries/codex-acp: 0.16.0
```

- The adapter session advertised `gpt-5.6-sol` but not Terra or Luna.
- Its advertised config options contained model and mode controls but no
  independent reasoning control.
- `session/set_config_option reasoning_effort=high` failed with ACP `-32602`.
- The contemporaneous Codex model metadata listed `high` as supported by Sol
  and listed Sol, Terra, and Luna as official identifiers.
- ACPX `0.12.0` documents its built-in Codex command as
  `npx -y @agentclientprotocol/codex-acp`, but the generated ACPX User Config
  contained `agents.codex.command = "codex-acp"`. That override resolved the
  legacy global package and masked ACPX's current built-in.

A second disposable proof used the official adapter directly:

```text
@agentclientprotocol/codex-acp: 1.1.4
advertised models:              gpt-5.6-sol, gpt-5.6-terra,
                                gpt-5.6-luna, gpt-5.5, ...
advertised reasoning:           low, medium, high, xhigh, max, ultra
proved selection:               gpt-5.6-sol / high
```

The proof used ACP session configuration only, sent no Agent prompt, and
closed the disposable Session. It confirms that Sol/high is a valid default
through the current official adapter and that Terra/Luna are valid advertised
identifiers. It does not make Terra or Luna generated defaults.

Therefore ADR-0040's blanket statement that GPT-5.6 is model-managed is not a
safe compatibility rule, and executable presence is not adapter readiness.
Production code must not read the private Codex model cache or ACPX Session
files used during bounded diagnosis. Only verified adapter identity plus live
advertised ACP capabilities and the resulting effective session state may
prove an operational selection.

## Data Model

Add Agent-owned capability and proof values. Names are illustrative; the
implementation may refine them without changing the public JSON contracts:

```go
type SelectionCapabilities struct {
    Models          []ModelCapability
    ReasoningOption *SelectCapability
}

type ModelCapability struct {
    AdapterValue    string
    CanonicalModel  string
    ReasoningEffort string
}

type SelectionAssignment struct {
    Requested       config.AgentSelection
    AdapterModel    string
    ReasoningKey    string
    ReasoningValue  string
    Encoding        string // independent, model_variant, model_managed
}

type SelectionProof struct {
    Selection       config.AgentSelection
    Assignment      SelectionAssignment
    AdapterCommand  string
    AdapterVersion  string
    Status          string
}
```

`SelectionCapabilities` is disposable-session evidence, not persisted product
configuration. `SelectionProof` may be projected into CLI JSON and Run Events,
but runtime-private payloads and unbounded option descriptions are omitted.
Canonical profile identity, deduplication, persistence, and fallback equality
continue to use the exact runtime/model/reasoning tuple.

## Capability Acquisition

Roundfix must obtain capability evidence from the ACP Session used for proof.
ACP `session/new` and `session/set_config_option` expose authoritative
`configOptions`; the latter returns the complete current configuration because
changing one option may change dependent options. Extend the runner seam to
consume that public machine-readable protocol result, including the full state
returned after model and reasoning updates. If ACPX cannot project the required
documented response without private persistence, pin a version that can or add
the missing public projection upstream. Roundfix must not read
`~/.acpx/sessions` directly in production.

The snapshot parser accepts only the documented fields needed for proof:

- current model value;
- advertised model option values;
- advertised select configuration IDs, current values, and option values;
- adapter/runtime identity and command version when available.

Missing, malformed, contradictory, or duplicate capability mappings fail
closed. Human names and descriptions are diagnostic only. The adapter value is
authoritative for assignment; a mapping to a canonical tuple is accepted only
when it is explicit in the advertised value/metadata or follows an adapter
format covered by focused compatibility tests. Roundfix does not infer
reasoning from marketing names such as Sol, Terra, or Luna.

## Assignment Planning

Given one canonical Agent Selection, planning follows this order:

1. Require a non-empty runtime and model. Preserve an empty effort as the
   explicitly requested model-managed state.
2. Find advertised model representations for the canonical model.
3. For a non-empty effort, prefer an advertised independent reasoning option
   when it contains the exact value.
4. Otherwise accept an advertised model variant only when its metadata maps
   unambiguously to the canonical model and exact effort.
5. If neither representation exists, return a typed unsupported-capability
   failure. Never retry with an empty effort or another model.
6. For an empty effort, select only the base/model-managed representation and
   assign no reasoning option.

The plan is deterministic. Multiple equivalent representations use a stable
preference documented above; ambiguous or conflicting representations fail
instead of guessing. Runtime-specific encoding knowledge stays inside
`internal/agent`, not CLI or config packages.

## Exact Disposable-Session Proof

The prover retains ADR-0039's unique disposable Agent Session and cleanup on
every path:

1. verify the pinned ACPX command and resolve the adapter command;
2. ensure the disposable session without relying on runtime-owned defaults as
   proof;
3. capture advertised capabilities;
4. build the assignment plan;
5. apply its model and optional reasoning operations in the plan's order;
6. capture the resulting effective state;
7. compare effective model and reasoning with the canonical requested tuple;
8. close the session under an independent bounded cleanup context;
9. return proof only after successful comparison and cleanup.

A command exit of zero without observable matching effective state is not
proof. A rejected application, absent control, stale advertised state, or
cleanup failure returns a structured error. The live Agent Session uses the
same assignment planner and operations. A capability-shape change between
preflight and live setup fails the selection before the first prompt; the
configured fallback path may then follow ADR-0050 after notification.

## Shared Readiness Coordinator

Move tuple deduplication, category references, sequential proof, and error
formatting behind one internal service. It accepts resolved profiles and
returns ordered proofs keyed by canonical tuple plus all references:

```go
type ProfileReadiness struct {
    Proofs []ProfileProof
    Err    error
}

func ProveProfiles(
    context.Context,
    []config.ResolvedProfile,
    string,
) ProfileReadiness
```

Consumers differ only in transaction timing and presentation:

- `profiles validate` resolves requested categories and prints the complete
  proof set.
- operational preflight resolves the categories required by the Run and must
  finish before Run/worktree/artifact mutation.
- `profiles configure` validates the candidate document, proves every distinct
  candidate tuple, then asks for write confirmation; no proof or declined
  confirmation leaves the target bytes unchanged.
- setup constructs the proposed generated profile document in memory, proves
  it before either User Config or Project Config write, then follows its
  existing confirmation policy.
- Doctor resolves all required categories, invokes the same service, renders
  one aggregate readiness line, and continues its independent checks.

Proofs are not cached across command invocations. Within one invocation, an
exact tuple shared by categories or repeated fallback positions is proved once
and all references are reported.

## Generated Defaults and Official Identifiers

Update `builtinProfiles` and `DefaultConfigYAML`:

```yaml
general: &general
  preferred: { runtime: codex, model: gpt-5.6-sol, reasoning_effort: high }
  fallbacks:
    - { runtime: codex, model: gpt-5.5, reasoning_effort: xhigh }
```

`backend`, `qa`, and `review` use the same canonical tuple values without YAML
anchors in emitted config. Frontend retains `claude-fable-5 / medium` as
Preferred and its configured Codex fallback, subject to exact proof.

The Model Catalog continues to accept and present the official identifiers
`gpt-5.6-sol`, `gpt-5.6-terra`, and `gpt-5.6-luna`. A valid identifier is not a
claim of availability in every installed environment. The official adapter
`1.1.4` advertises all three, but recommendations remain advisory and exact
runtime proof remains the only readiness signal. Generated operational
profiles do not use Terra or Luna because the accepted product policy keeps
Sol/high with GPT-5.5/xhigh fallback; changing that policy requires a separate
decision, not a compatibility workaround.

## Adapter Provisioning and Identity

The supported Codex adapter is the official
`@agentclientprotocol/codex-acp` package at the version pinned by Roundfix. The
initial compatibility floor and tested pin for this Spec is `1.1.4`. Keep the
pin in one Agent-owned constant used by Setup, Doctor, and Preflight
Validation; do not scatter version strings through CLI code or docs.

Setup may use ACPX's built-in command only when Roundfix can prove that it
resolves to the same supported package/version. When Setup writes an explicit
ACPX agent override, it writes a deterministic supported command rather than a
bare name discovered on PATH. Existing bare overrides are resolved and
diagnosed. A legacy `@zed-industries/codex-acp`, unknown package lineage,
version below the compatibility floor, or version that cannot be proven fails
closed with one install/update command. Roundfix requests confirmation before
migrating an existing User Config and never silently rewrites it.

## Setup Transaction Boundary

Setup must never create the configuration and validate it afterward. Its order
is:

1. load existing config and resolve the effective adapter command;
2. prove official adapter lineage and the pinned supported version, proposing
   an authorized migration when an existing override is stale;
3. build proposed User Config and Project Config bytes in memory;
4. resolve the effective profiles that those proposals would produce;
5. prove every distinct Preferred Selection and fallback exactly;
6. if proof fails, print `profile readiness: failed`, the exact tuple,
   capability classification, installed component evidence, and next action;
7. exit non-zero without creating, truncating, or changing either config;
8. only after adapter and profile proof succeed, apply the existing
   prompt/`--yes` write policy.

`--no-input` performs diagnosis and writes nothing. A proof failure is not an
offer to substitute model-managed reasoning. When the adapter does not expose
the requested control, the next action prioritizes updating the incompatible
adapter/toolchain and rerunning setup; profile configuration is mentioned only
as an explicit user-owned alternative.

## One-Run Override Grammar

For `resolve`, `watch`, and `implement`, parse invocation selection flags as
one atomic object:

- none present: resolve category profiles normally;
- all of `--agent`, `--model`, and `--reasoning-effort` present: replace only
  each relevant profile's Preferred Selection and retain its Fallback Chain;
- any non-empty proper subset: exit `2` before Agent Session or Run creation.

An explicitly empty `--reasoning-effort ""` counts as present and requests
model-managed reasoning. Empty runtime or model remains invalid. Help and error
copy explain both valid forms. Remove bare `--agent` from canonical usage and
manifest commands. Fetch remains Agent-free.

This is intentionally not runtime-only override behavior: combining a runtime
from the invocation with model/reasoning from another runtime's profile can
produce an invalid or misleading tuple.

## Profile-Aware Doctor Integration

Doctor replaces the legacy single configured-runtime probe as its Agent
Selection readiness authority. It resolves `general`, `backend`, `frontend`,
`qa`, and `review`, including optional-category inheritance when applicable,
then invokes the shared coordinator. Its deterministic result is additive to
the existing Node, ACPX, adapter, Codex, and future Repository Skill Set lines.

Success includes the number of distinct tuples and category references.
Failure identifies the first failed tuple in stable proof order, all category
references, the classification, and next action while Doctor continues all
unrelated checks. JSON is not added to Doctor in this Spec. `profiles validate
--json` remains the detailed machine-readable proof surface.

Spec 0036 must inject and append its independent `skills:` result without
recreating or bypassing this coordinator.

## Error Taxonomy and Diagnostics

Extend selection proof failures with stable classifications:

- `adapter_lineage_unknown`;
- `adapter_version_unsupported`;
- `model_not_advertised`;
- `reasoning_control_not_advertised`;
- `model_variant_not_advertised`;
- `selection_rejected`;
- `effective_selection_mismatch`;
- `capability_evidence_invalid`;
- `session_cleanup_failed`.

Text errors name the operation and next action. JSON adds optional additive
fields for classification, encoding, adapter command/version, and bounded
advertised values while retaining the `roundfix/profiles-validate/v1` schema
unless an existing required field must change. Secrets, environment values,
session paths, raw private metadata, and unbounded adapter output are never
included.

When a profile lacks a distinct fallback or duplicates its Preferred
Selection, configuration validation remains pre-proof and no Session is
created. The error explicitly states that one additional distinct authorized
and proven Agent Selection is required. Roundfix never inserts a ranked model
to satisfy this invariant.

## Coverage Map

- Goals 1–3 and Stories 1–4 → generated defaults, official identifiers,
  capability acquisition, assignment planning, and exact proof.
- Goals 4–6 and Stories 1, 3, 7 → shared readiness coordinator, setup
  transaction, profile-aware Doctor, and operational preflight.
- Goals 7–8 and Stories 5–6, 8 → all-or-none override grammar and
  distinct-fallback diagnostics.
- ADR-0050 compatibility → live-session assignment reuse and pre-prompt-only
  fallback activation after notification.

## Testing Approach

- Agent package table tests model independent reasoning controls, advertised
  model variants, explicit model-managed selection, unsupported effort,
  unavailable model, ambiguous variants, malformed capability evidence,
  effective-state mismatch, cancellation, and cleanup error joining.
- A legacy-adapter contract fixture represents the observed
  `@zed-industries/codex-acp 0.16.0` shape: Sol
  advertised, no reasoning control, Terra/Luna absent. Sol/high, Terra/max, and
  Luna/max must fail with distinct classifications; Sol/model-managed may pass
  only when requested explicitly.
- An official-adapter `1.1.4` fixture advertises Sol, Terra, Luna, GPT-5.5, and
  the independent reasoning values. It proves Sol/high and GPT-5.5/xhigh
  without changing Roundfix's Model Catalog.
- Adapter identity tests cover the supported official package, legacy Zed
  package, unknown same-named PATH executable, unsupported version, ACPX
  built-in resolution, existing override migration, declined migration, and
  byte-identical no-mutation failures.
- Coordinator tests prove exact-tuple deduplication, stable category references,
  fallback order, deterministic first failure, and no proof caching across
  invocations.
- Setup tests snapshot absent and existing User/Project Config bytes and prove
  zero mutation when preferred or fallback proof fails, including `--yes` and
  `--no-input`.
- Config tests assert Sol/high plus GPT-5.5/xhigh generated defaults and that
  Terra/Luna remain accepted catalog/recommendation values rather than emitted
  operational fallbacks.
- CLI tests cover no override, every partial flag subset, the complete triple,
  explicit empty reasoning, profile selection by Task Type/QA/review, exact
  exit `2`, stdout/stderr ownership, and no Run/worktree/session side effects.
- Doctor and `profiles validate` share injected proof fixtures and must report
  identical tuple/category failure evidence.
- Real QA uses disposable repositories and sessions to exercise both the
  legacy incompatible adapter and the current official adapter. It records
  versions, advertised controls, commands, output, and cleanup evidence
  without reading private runtime caches as product authority.
- Run `go test -race ./...` because the shared coordinator is used by
  operational paths that may prepare multiple Task categories.
- Run `make verify` and the Roundfix Skill synchronization checks.

## Build Order

1. Add the official adapter pin, effective-command identity/version proof,
   migration diagnostics, and focused tests.
2. Introduce capability/proof types, the public machine-readable ACP response
   seam, parsers, and classified errors (depends on: 1).
3. Implement deterministic assignment planning and exact disposable/live
   session application with focused Agent tests (depends on: 2).
4. Centralize profile resolution, tuple deduplication, and readiness proof;
   migrate `profiles validate` and operational preflight (depends on: 3).
5. Enforce all-or-none one-Run overrides and remove bare `--agent` from CLI
   usage/tests (depends on: 4).
6. Change generated defaults and make setup prove the in-memory proposal before
   writing (depends on: 4).
7. Make `profiles configure` prove candidates before confirmation/write and
   improve distinct-fallback diagnostics (depends on: 4).
8. Replace Doctor's legacy runtime probe with profile-aware readiness and align
   Spec 0036's injection/output order (depends on: 4).
9. Update CONTEXT, configuration/command docs, canonical Roundfix Skill,
   embedded copy, ADR/finding traceability, and QA evidence (depends on: 5–8).

## Risks and Considerations

- The pinned ACPX version may not project every documented ACP configuration
  response through its CLI. The implementation must pin a sufficient public
  contract or coordinate an upstream addition; reading ACPX's private session
  files is not an acceptable shortcut.
- Adapter model-variant syntax can evolve. Only advertised, unambiguous
  mappings are accepted; no static bracket construction is sent blindly.
- Exact effective-state proof adds an ACP inspection round trip. Quality and
  deterministic correctness outrank the small preflight latency and token-free
  control-plane cost.
- Blocking Setup on adapter identity may require migrating a legacy explicit
  override even when ACPX itself is current. This is intentional: executable
  presence is not proof that the selected adapter implements the required
  capability contract.
- A previously accepted bare `--agent` command becomes a usage error. This is
  an intentional public CLI correction and must be called out in release notes.
- Spec 0036 and this Spec both change Doctor dependencies/output. Implement
  this Spec first, then rebase 0036 on the shared coordinator rather than
  creating parallel selection checks.

## Decisions

- Capability-driven exact proof is the sole operational readiness authority.
- The official `@agentclientprotocol/codex-acp` package and supported pinned
  version are part of runtime readiness; same-name PATH discovery is not.
- Sol/high remains the canonical desired default and GPT-5.5/xhigh becomes its
  generated fallback.
- Model-managed reasoning is accepted only when explicitly configured empty.
- Setup blocks before persistence when generated profiles cannot be proven.
- One-Run Agent Selection overrides are all-or-none.
- Terra and Luna are valid official identifiers and advisory candidates;
  operational readiness remains subject to exact proof in the effective
  environment.
- See [ADR-0055](../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md).
