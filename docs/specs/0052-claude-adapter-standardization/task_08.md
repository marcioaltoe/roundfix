---
task: task_08
spec: 0052-claude-adapter-standardization
status: completed
type: docs
complexity: low
---

# Task 08: Align the protected Roundfix Skill pair and derived digest pins

## Overview

Update the canonical Roundfix Skill so its adapter, Doctor, Setup, profile,
and selection guidance matches the shipped CLI — the Claude adapter contract
beside the Codex one, the raised Codex pin, the multi-runtime `adapter:`
line, the new frontend default, and the opaque-identifier rule — and
propagate the deterministic Skill-digest fallout to the derived Baseline
pins. This is the Spec's only protected-tooling mutation and is expressly
authorized in the PRD and TechSpec Tooling authority entries.

## Requirements

1. MUST update the canonical Skill and its embedded copy so the documented
   adapter contract covers Claude exactly as it covers Codex: official
   package, `0.63.0` pin, official install action, legacy-lineage failure
   behavior, and the raised Codex pin `1.1.5`.
2. MUST update the Skill's Doctor and Setup output examples, built-in
   profile description (`frontend` preferred `claude / opus / xhigh`), and
   selection guidance including the opaque-identifier rule.
3. MUST keep the Skill pair byte-identical to each other.
4. MUST propagate the changed Skill `contentDigest` through the derived
   pins: the setup asset's entry and top-level digest, the normalized
   catalog snapshot and its digest, and the parity-corpus fixture and
   manifest rows.
5. MUST change only the expressly authorized files: `.agents/skills/roundfix/SKILL.md`,
   `skills/roundfix/SKILL.md`, `internal/baseline/assets/setups/typescript-bun.json`,
   `internal/baseline/testdata/catalog.digest`,
   `internal/baseline/testdata/catalog.normalized.json`,
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
   `internal/baseline/testdata/parity-corpus/v1/manifest.json`, plus this
   Task file.

## Subtasks

- [ ] Rewrite the Skill's adapter, Doctor, Setup, and Agent selection
      sections to the shipped contract.
- [ ] Synchronize the embedded Skill copy byte-identically.
- [ ] Recompute and propagate the five derived digest pins from the
      canonical sources.

## Acceptance Criteria

- [ ] The Skill names `@agentclientprotocol/claude-agent-acp` at `0.63.0`
      and the official install action, and no longer states the Codex pin as
      `1.1.4` or the frontend preferred model as `claude-opus-5`.
- [ ] `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` are
      byte-identical.
- [ ] The embedded Baseline catalog validates: setup digests, normalized
      snapshot, and parity-corpus rows agree with the edited Skill.
- [ ] The change touches only the authorized files plus this Task file.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `internal/baseline/assets/setups/typescript-bun.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`

## Verification

- `make skills-sync-check` — expected: the Skill pair is synchronized.
- `go test -count=1 ./internal/baseline/` — expected: pass; the embedded catalog, setup digests, and parity corpus validate against the edited Skill.
- `grep -n 'agentclientprotocol/claude-agent-acp' skills/roundfix/SKILL.md` — expected: at least one match.
- `grep -n 'zed-industries/claude-code-acp' skills/roundfix/SKILL.md ; test $? -eq 1` — expected: no matches (exit 1).

## References

`_prd.md` → Core Feature 9, Project Constraints: Tooling authority;
`_techspec.md` → Build Order 7, Project Constraints: Tooling authority;
ADR-0079.

## Result

Implementation:

- The canonical and shipped Roundfix Skills now document official Codex
  `1.1.5` and Claude `0.63.0` adapter lineage, pinned `npx` commands,
  deterministic `npm install -g` actions, confirmed legacy-lineage migration,
  and fail-closed behavior for legacy, unknown, and below-pin adapters.
- Doctor and Setup guidance now covers every distinct runtime referenced by
  the effective required profiles and shows the shipped, runtime-sorted
  Claude-plus-Codex `adapter:` evidence line.
- Agent selection guidance now uses the built-in
  `claude / opus / xhigh` frontend Preferred Selection. It keeps dated Claude
  catalog labels advisory and explains the ADR-0079 opaque-identifier rule:
  `opus[1m]` and `opus` remain selectable with independent reasoning, while
  `[1m]` is not accepted as a reasoning effort.
- The Skill copies remain byte-identical. Their canonical
  `contentDigest` is
  `4ec6dbf928c6c726e08f60bbf980b3be85cc0554179c3959b3c5695c065395db`.
- The derived Baseline pins now carry setup digest
  `b3268bda435a66b80e93b4009bde92ee2844976cb9f328be75d1be2fe7973e63`,
  catalog digest
  `sha256:3062d035a1d6822e33aa9e1938f7fbb34278f5c2dcf7f6c44177a5fe7f843732`,
  and parity fixture identity
  `e7a97493633bbccca88bf9b9635adc2adf43f8cdf0df0cd3f14600847d2d2bae`.

Focused checks:

- Before editing,
  `rtk grep -n -E '1\.1\.4|claude-opus-5|claude-agent-acp|adapter: ok|Adapter Readiness|opaque|advertised' .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md docs/user-guide/commands.md docs/user-guide/configuration.md docs/user-guide/usage.md internal/agent/acpx_runner.go internal/cli/doctor.go internal/cli/setup.go internal/config/profiles.go internal/agent/selection_assignment.go internal/agent/selection_capabilities.go`
  found the stale `1.1.4`, single-runtime Doctor example, and
  `claude-opus-5` frontend-selection guidance in both Skill copies.
- A read-only inline `rtk node -e` digest-chain assertion passed after the
  final edit. It proved Skill byte identity and required/stale text,
  recomputed the Skill and setup digests, matched the setup asset's byte
  identity in `catalog.normalized.json`, recomputed `catalog.digest`, and
  matched the parity fixture's byte count and SHA-256 in the corpus manifest.
- `rtk go test -count=1 ./skills -run 'TestBaselineSkillContract|TestAuthorialSkillSync'`
  passed: `21 passed in 1 packages`.
- `rtk go test -count=1 ./internal/baseline -run 'TestCatalogDigest|TestBaselineCompatibilityCorpus'`
  passed: `2 passed in 1 packages`.
- `rtk git -c core.fsmonitor=false diff --check` exited `0`.
- `rtk git -c core.fsmonitor=false diff --name-only` listed exactly the seven
  expressly authorized protected files and this Task file.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  that gate.

Acceptance evidence:

- Adapter contract: both Skill copies name
  `@agentclientprotocol/claude-agent-acp` at `0.63.0`, its official install
  action, and Codex `1.1.5`; the focused text assertion found no `1.1.4`,
  fully scoped deprecated Claude Code package, or stale frontend Preferred
  Selection.
- Skill synchronization: byte comparison in the digest-chain assertion and
  the focused Skills tests both passed.
- Baseline consistency: the recomputed Skill, setup, normalized catalog,
  catalog digest, parity fixture, and parity manifest identities all agree;
  the targeted catalog digest and compatibility-corpus tests passed.
- Bounded scope: postflight Git evidence contains only the exact PRD/TechSpec
  allowlist plus `task_08.md`; the Daemon-owned frontmatter status remains
  `in_progress`.

Follow-up:

- `docs/agents/autonomous-work.md` still describes the older frontend model.
  It is outside this Task's exact protected-tooling allowlist and remains a
  separate documentation follow-up.
