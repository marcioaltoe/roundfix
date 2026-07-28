---
task: task_08
spec: 0052-claude-adapter-standardization
status: pending
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
