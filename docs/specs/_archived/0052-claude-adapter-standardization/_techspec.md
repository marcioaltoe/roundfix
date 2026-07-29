---
spec: 0052-claude-adapter-standardization
prd: _prd.md
created: 2026-07-28
---

# Official Claude adapter and opaque model identifiers — Technical Spec

## Executive Summary

This Spec gives the claude ACP Runtime the same adapter lineage contract Codex
already has, and makes advertised Agent Model identifiers opaque whenever an
adapter advertises an independent reasoning control. The primary trade-off is
in Claude lineage proof: `claude-agent-acp --version` prints a bare semver with
no package name (verified: official `0.63.0` prints `0.63.0`; deprecated
`claude-code-acp` prints nothing), so the two-field probe contract Codex uses
cannot prove Claude package identity. We accept a hybrid proof — package
identity from the effective command when it names the official package, or from
the executable's resolved installation path otherwise — instead of restricting
support to the pinned `npx` form only. That keeps `npm install -g` and Homebrew
installs provable while preserving the rule that a matching executable name
alone is never Adapter Readiness proof. All parsing changes stay inside the
existing capability/assignment seams; no new package or directory is created.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier
  is created; adapter packages, advertised Agent Model identifiers, and Agent
  Selection tuples keep their source-native identities. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all changes are local adapter
  probing, capability parsing, configuration defaults, CLI reporting, and
  documentation; no authentication provider, credential, or HTTP route
  changes. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0037/0039 keep selection
  Roundfix-owned and proven through disposable Agent Sessions; ADR-0049 keeps
  profiles atomic; ADR-0050 keeps Fallback Chain activation
  post-Run-creation; ADR-0055 makes advertised capabilities and adapter
  lineage the proof boundary and forbids silent override edits; ADR-0079
  makes advertised model identifiers opaque under an independent reasoning
  control and supersedes ADR-0040's Claude premise while preserving its
  explicit-empty semantics. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

Everything extends existing modules:

- `internal/agent` — the adapter lineage contract, adapter inspection
  (`CheckAdapter`), capability parsing (`ParseSessionConfigOptions`,
  `parseModelCapability`), and assignment planning
  (`PlanSelectionAssignment`, `unsupportedSelection`).
- `internal/config` — built-in Agent Selection Profiles and default Project
  Config rendering.
- `internal/cli` — the Doctor Command adapter check, the Setup Command's acpx
  override migration, and the Preflight refusal diagnostics in the profile
  proof error.
- `docs/user-guide`, the Roundfix Skill pair, and the authorized derived
  digest pins — documentation of the new contract.

No new packages, layers, or directories. The Codex lineage machinery
(`inspectCodexAdapter`, `AdapterLineageError`, `AdapterVersionError`)
generalizes into runtime-parameterized form rather than growing a parallel
Claude copy.

## Implementation Design

### Interfaces

Runtime lineage contracts (in `internal/agent`, beside the existing Codex
constants):

```go
type adapterLineageContract struct {
    RuntimeID      string   // "codex", "claude"
    Package        string   // official npm package
    PinnedVersion  string   // minimum accepted version
    LegacyPackages []string // recognized deprecated lineages
    VersionOnly    bool     // --version prints bare semver (claude)
}

// ClaudeAdapterPackage = "@agentclientprotocol/claude-agent-acp"
// PinnedClaudeAdapterVersion = "0.63.0"
// PinnedCodexAdapterVersion rises "1.1.4" -> "1.1.5"
// ClaudeAdapterCommand() / ClaudeAdapterInstallCommand() mirror Codex helpers.
```

`CheckAdapter` drops the codex-only early return: runtimes with a lineage
contract (`codex`, `claude`, including `-custom` forms) go through
`inspectAdapter(ctx, invocation, contract)`; `opencode` keeps the current
bare-evidence path. `AdapterLineageError` and `AdapterVersionError` gain the
required package, pin, and install command as data instead of Codex-hardcoded
prose; classifications `adapter_lineage_unknown` and
`adapter_version_unsupported` are unchanged.

`inspectAdapter` proof per contract:

- Two-field `<package> <version>` output (Codex today, Claude if upstream
  adds it): package must equal the official package; version compared against
  the pin.
- `VersionOnly` output (Claude `0.63.0`): the version comes from the probe;
  package identity comes from, in order, (1) the effective command naming the
  official package (the `npx -y <pkg>@<ver>` form), or (2) the executable's
  symlink-resolved path containing a `node_modules/<package>/` segment —
  official proves, a legacy package classifies as legacy lineage, anything
  else fails as unproven with the official install action and pinned override
  as next action.
- Empty or malformed output (deprecated `claude-code-acp` prints nothing):
  lineage error; path resolution still runs so a resolvable legacy install
  gets the migration message rather than a generic probe failure.

Capability parsing (`internal/agent/selection_capabilities.go`): resolve the
independent reasoning option before parsing the model list (today the model
loop runs first), then:

```go
// parseModelCapability(value, opaque bool)
// opaque (reasoning option advertised):
//   "opus[1m]" -> {AdapterValue: "opus[1m]", CanonicalModel: "opus",
//                  ReasoningEffort: "", ModelManaged: true}
// non-opaque (no reasoning option): existing canonical[effort] variant parse.
```

The bracket is still stripped into `CanonicalModel` so the canonical prefix
stays selectable, but no `ReasoningEffort` is extracted and the entry counts
as a model-managed base, which is exactly what `PlanSelectionAssignment`
needs for the independent encoding. `modelsForCanonical` additionally matches
`model.AdapterValue`, making the advertised form (`opus[1m]`) selectable.
`unsupportedSelection` renders each advertised model as its canonical form,
appending the raw advertised value when the two differ
(`opus (advertised opus[1m])`), so any identifier copied from the diagnostic
proves cleanly. `PlanSelectionAssignment` itself is unchanged: `opus / 1m`
now falls through the independent check (the reasoning option does not
advertise `1m`), finds no variants, and is rejected with the advertised
efforts in the message.

### Data Models

- `ModelCapability` keeps its shape; only the parse rules change.
- Built-in `frontend` profile (`internal/config`): Preferred Selection model
  `claude-opus-5` → `opus`; Fallback Chain unchanged; the rendered default
  Project Config mirrors it.
- The Model Catalog is unchanged: `opus` is an adapter identifier proven by
  Exact Agent Selection Proof, not a catalog entry.

### API Contracts

- Doctor `adapter:` line: inspects every distinct runtime referenced by the
  effective required profiles (Preferred Selections and Fallback Chains,
  deduplicated, sorted by runtime ID) instead of only the `general`
  preferred. One line, per-runtime evidence:
  `adapter: ok (claude: command=…; package=…; version=… | codex: …)`.
  Any failing runtime fails the line; `next:` names that runtime's install
  action. Exit-code behavior unchanged.
- Setup: the stale-adapter migration gate generalizes from codex-only to
  per-runtime. A stale or bare Claude override is diagnosed, a migration to
  `ClaudeAdapterCommand()` is proposed and proved before the offer, and
  decline preserves every byte — identical choreography to the Codex
  migration. `officialClaudeACPXOverride()` sits beside
  `officialCodexACPXOverride()`; the local-override seeding that today emits
  a bare `claude-agent-acp` command emits the official pinned form.
- Install hints: the `adapterInstallCommands` entries for `claude-code-acp`
  and `claude-agent-acp` (both currently `@zed-industries`, one wrong-scope)
  are replaced by the official install command via the same short-circuit
  used for Codex.
- Preflight refusal: `profileProofError` appends one sentence —
  `fallback: Fallback Chains activate only after Run creation (ADR-0050);
  Preflight proves every configured tuple and substitutes none` — so a
  refusal caused by a broken Preferred Selection explains why a configured
  fallback did not rescue the Run. Existing message fields keep their exact
  text.

## Coverage Map

- Goal 1 / Story 1 → generalized `CheckAdapter` + `inspectAdapter`, Doctor
  multi-runtime adapter check.
- Goal 1 / Story 2 → Setup per-runtime migration + `officialClaudeACPXOverride`.
- Goal 2 / Stories 3–4 → opaque `parseModelCapability`, `modelsForCanonical`
  AdapterValue match, `unsupportedSelection` rendering.
- Goal 3 / Story 5 → built-in `frontend` profile change.
- Story 6 → `profileProofError` fallback sentence.
- Goal 4 → install-hint fixes plus the documentation/Skill task; ADR-0040's
  premise is superseded by ADR-0079 without editing the legacy ADR file.
- PRD Core Feature 7 → `PinnedCodexAdapterVersion` raise and Setup proposal
  string.

## Integration Points

- npm registry — only through printed install commands; no network calls.
- acpx (`~/.acpx/config.json`) — read for effective commands; written only
  through the existing Setup confirmation flow (ADR-0055).
- The adapter processes themselves — the existing bounded `--version` probe;
  no new adapter invocations.

## Testing Approach

All coverage attaches to existing seams; no new test infrastructure:

- `internal/agent/acpx_runner_test.go` — extend the `CheckAdapter` tables:
  official Claude (version-only output + command-named package), legacy
  lineage via resolved path, bare same-named executable failing as unproven,
  below-pin version. `TestCheckAdapterPreservesNonCodexResolutionWithout-
  VersionExecution` is rewritten for the new claude contract; opencode keeps
  the bare-evidence assertion. Fake-adapter fixtures under `newFakeACPXHarness`
  and `newMacroFakeACPX` (`implement_test.go`) rename their claude command to
  the official adapter.
- `internal/agent/selection_capabilities_test.go` — new fixture: bracketed
  identifiers plus an independent `effort` option ⇒ opaque parse, `opus` and
  `opus[1m]` both plan `independent`, `opus / 1m` rejected naming advertised
  efforts; the existing `modelVariantCapabilityFixture` (bracket, no
  reasoning option) keeps proving the variant encoding.
- `internal/cli/doctor_test.go` — multi-runtime adapter line: both runtimes
  reported, deterministic order, failing claude names the official install
  action; legacy-lineage assertions mirror the existing Codex case.
- `internal/cli/cli_test.go` — Setup: Claude migration accept/decline tables
  beside the Codex ones; fresh-machine offer writes the official pinned
  Claude override; unrelated-bytes preservation unchanged.
- `internal/cli/implement_test.go` — Preflight refusal asserts the ADR-0050
  fallback sentence.
- Documentation contract tests already compare the Skill pair byte-identity
  and shipped-help wording; they gate the docs task.

## Build Order

1. Claude lineage contract in `internal/agent`: constants, generalized
   lineage errors, `inspectAdapter`, `CheckAdapter` without the non-codex
   early return, official install hints, Codex pin raise to `1.1.5`.
2. Opaque advertised identifiers in `internal/agent`: reasoning-option-first
   parse order, opaque `parseModelCapability`, `modelsForCanonical`
   AdapterValue match, `unsupportedSelection` dual-form rendering.
3. Built-in `frontend` Preferred Selection → `claude / opus / xhigh` in
   `internal/config` (depends on: 2).
4. Doctor multi-runtime Adapter Readiness (depends on: 1, 3).
5. Setup per-runtime migration and official Claude override (depends on: 1).
6. Preflight refusal fallback sentence in `internal/cli` (no dependency;
   schedule after 2 to keep diagnostic tests stable).
7. Documentation, user guide, Roundfix Skill pair, and authorized derived
   digest pins (depends on: 1–6).

## Risks & Considerations

- **Version-only lineage proof.** If a future official release changes its
  `--version` shape, the probe accepts the two-field form first, so upstream
  adding a package name is forward-compatible. Path resolution can fail on
  layouts without a `node_modules` segment (for example a source checkout);
  those fail closed with the pinned-override next action, which is the
  documented remedy.
- **Dedup collision — observed, not hypothetical.** The QA gate proved on
  2026-07-28 that the official adapter **echoes the requested model back into
  the advertised list**: probing with `--model opus` returns
  `default, opus[1m], claude-fable-5[1m], sonnet, haiku, opus`. Under opaque
  parsing both `opus[1m]` and `opus` yield canonical `opus` with empty effort,
  so the canonical-keyed dedup raised `ambiguous_model_variant` and failed the
  built-in frontend tuple on every real Run. The resolution is recorded in
  Decisions below: under opaque parsing the advertised values are already
  unique, so entries sharing a canonical alias are an alias group rather than
  an ambiguity, and the fail-closed variant check stays exactly as-is for
  adapters that use the variant encoding.
- **Doctor line-shape change.** The `adapter:` detail becomes per-runtime;
  the deterministic-order test and the Skill/docs examples change in the same
  Spec, so no drift window.
- **Windows shims.** `.cmd` shims do not resolve through symlinks; bare
  Claude commands there fail as unproven with the pinned-override remedy.
  Codex behavior is unaffected.

## Decisions

- One runtime-parameterized lineage inspection replaces a parallel Claude
  copy of the Codex machinery; error types carry the contract as data.
- Claude package identity comes from the effective command or the resolved
  installation path because the adapter's version output carries no package
  name (verified 2026-07-28); a matching executable name alone remains
  insufficient.
- Advertised identifiers under an independent reasoning control parse to a
  model-managed base with a canonical alias, keeping both `opus` and
  `opus[1m]` selectable with explicit effort. See ADR-0079.
- Alias groups are not ambiguity. Under opaque parsing the dedup keys on
  `AdapterValue` — already unique by construction — and never raises the
  variant-ambiguity issue; the `(CanonicalModel, ReasoningEffort)` check keeps
  failing closed for adapters that genuinely use the variant encoding.
  `modelsForCanonical` orders an exact `AdapterValue` match ahead of a
  canonical-alias match, so a requested `opus` binds to the echoed `opus` when
  the adapter returns one and to `opus[1m]` when it does not — deterministic
  either way.
- The `-custom` runtime suffix keeps its existing semantics: `claude-custom`
  now inspects lineage exactly like `codex-custom` does.
