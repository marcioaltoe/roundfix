---
type: fix
status: open
created: 2026-11-19
spec: null
reason: null
---

# The Baseline semantic analysis proves its selections under a fixed 30-second deadline

## Symptom

On a loaded machine, the exact selection proof of the Baseline semantic analysis can time out after 30 seconds even when the repository raised `defaults.profile_proof_setup_timeout`. The analysis is then unavailable and preserves every original Source Baseline Entry as one lossless range, so the classification falls back to manual work; nothing is lost, but the analysis the operator asked for does not run. Low impact.

## Where

`Analyzer.proveSelection` in `internal/baselineacp/analyzer.go:330` builds `agent.ProbeRequest{Runtime: selection, WorkDir: workDir}` without `SetupTimeout`, so `ProveExactSelection` uses the built-in `acpxPreflightSetupTimeout` (`internal/agent/selection_assignment.go:118-121`, 30 s). The analyzer is built by `baselineacp.NewDefaultAnalyzer()` in `internal/cli/baseline_human.go:260` before any configuration is loaded, and the `roundfix baseline` command reads no Roundfix configuration on that path. ADR-0344 says "Baseline semantic analysis keeps its stop".

## Expected

The analysis proves under the configured `defaults.profile_proof_setup_timeout` when the target repository's configuration can be read, and keeps today's built-in deadline, without failing the command, when it cannot. Spec 0411 adds `config.ProfileProofSetupTimeout(homeDir, workDir)`, a reader of that one key that ignores every other key, which this path can reuse.

## Evidence

Measured by the author of Spec 0411 on 2026-11-18 at `a6a4a15b` (Spec 0411 `references/2026-11-18-characterization.md`, section 1, row 9). No field report yet.

## Decisions

- 2026-11-19, about 00:05 BRT (AskUserQuestion during the planning of Spec 0411): the maintainer chose to record it — "(b) prazo da análise do baseline".
