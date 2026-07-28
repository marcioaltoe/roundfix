---
task: task_01
spec: 0052-claude-adapter-standardization
status: pending
type: backend
complexity: high
---

# Task 01: Prove official Claude adapter lineage

## Overview

Give the claude ACP Runtime the same Adapter Readiness contract Codex has:
official package, pinned minimum version, recognized legacy lineages, and a
lineage inspection that fails a deprecated or unproven adapter with the
official install action. Verifiable on its own through the adapter-check unit
seams — no Doctor, Setup, or capability change is required yet.

## Requirements

1. MUST define the Claude lineage contract: official package
   `@agentclientprotocol/claude-agent-acp`, pinned minimum version `0.63.0`,
   recognized legacy lineages `@zed-industries/claude-code-acp` and
   `@zed-industries/claude-agent-acp`, and pinned command/install helpers
   mirroring the Codex helpers.
2. MUST make the default claude adapter command resolve to the official
   package form instead of the deprecated bare `claude-code-acp`.
3. MUST generalize adapter inspection so runtimes with a lineage contract
   (`codex`, `claude`, and their `-custom` forms) are inspected instead of
   early-returning bare evidence; `opencode` keeps its current bare-evidence
   behavior.
4. MUST prove Claude package identity per the TechSpec hybrid: a two-field
   `<package> <version>` probe when available; otherwise a bare-semver probe
   whose package identity comes from the effective command naming the
   official package or from the executable's symlink-resolved installation
   path. A matching executable name alone is never proof. Empty or malformed
   probe output fails as a lineage error, still classifying a
   path-resolvable legacy install as legacy.
5. MUST carry the required package, pin, and install command in the lineage
   and version error types as data, keeping classifications
   `adapter_lineage_unknown` and `adapter_version_unsupported` and their
   install-command next actions.
6. MUST replace both `@zed-industries` install hints with the official
   install command and raise the pinned minimum Codex adapter version to
   `1.1.5`.
7. SHOULD keep probe output bounded and stderr discarded exactly as the
   existing Codex probe does.

## Subtasks

- [ ] Add the Claude lineage constants and command/install helpers beside the
      Codex ones, and raise the Codex pin.
- [ ] Parameterize the lineage/version error types and their messages by
      contract.
- [ ] Implement the generalized `inspectAdapter` with the hybrid Claude proof
      and wire `CheckAdapter` through it.
- [ ] Replace the deprecated default claude command and both legacy install
      hints.
- [ ] Extend the adapter-check test tables: official Claude (version-only
      output), legacy lineage via resolved path, unproven bare executable,
      below-pin version, opencode unchanged; update fake-adapter fixtures
      that still name `claude-code-acp`.

## Acceptance Criteria

- [ ] A claude adapter resolving to either `@zed-industries` lineage fails
      with classification `adapter_lineage_unknown` and a next action naming
      `npm install -g @agentclientprotocol/claude-agent-acp@0.63.0`.
- [ ] An official Claude adapter probed at `0.63.0` or newer returns evidence
      with the official package and probed version; below `0.63.0` fails with
      `adapter_version_unsupported` and the install command.
- [ ] A bare executable named `claude-agent-acp` whose resolved path proves
      no package lineage fails as unproven with the official install action;
      one whose resolved path contains the official package's
      `node_modules` segment passes.
- [ ] Codex adapter behavior is unchanged except the pin: version `1.1.4`
      now fails with `adapter_version_unsupported`.
- [ ] `opencode` still returns bare command evidence with no probe.
- [ ] No source or test path under `internal/` names
      `@zed-industries/claude-code-acp` or defaults claude to
      `claude-code-acp`.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/agent/ -run 'TestCheckAdapter|TestACPXProbe|TestResolveAdapterCommand'` — expected: pass, including the new Claude lineage cases.
- `grep -rn 'zed-industries/claude-code-acp' internal/ ; test $? -eq 1` — expected: no matches (exit 1).
- `grep -n '"claude"' internal/agent/acpx_runner.go | grep -c 'claude-code-acp' | grep -x 0` — expected: `0`; the default claude command no longer names the deprecated binary.

## References

`_prd.md` → User Story 1, Core Features 1–3, 7; `_techspec.md` → Build Order 1,
Interfaces: adapterLineageContract; ADR-0055.
