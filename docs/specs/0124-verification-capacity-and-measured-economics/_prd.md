---
spec: 0124-verification-capacity-and-measured-economics
status: active
created: 2026-09-08
surfaces: [backend, infra, docs]
---

# Verification cost and failures are explained by measurements

Historical gate failures move between concurrent process tests while isolated reruns pass. The legacy Force Stop failure lost its assertion log, and the gate cost distribution has not been measured after later improvements. Current local/CI cache tiers also need a clear evidence boundary. The maintainer needs a measured capacity policy and failure records that support a cause rather than a guessed timeout change.

This Spec is **in authoring**. Its protected scope and open decisions are
pending approval. The accompanying authorization is a proposal with no grant;
there is no TechSpec, Task Graph, or authority to start implementation.

## Project Constraints

- Identifier strategy: not applicable — retain existing Run, Task, and diagnostic identities; measurements do not introduce a new persisted identity scheme. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — verification scheduling and local/CI test execution introduce no authentication or HTTP policy. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0056 separates per-Run Task and Verification Capacity and expressly does not coordinate the entire machine. ADR-0117 requires each defect to be checked by the stage that can produce it. Preserve these boundaries unless a measured, explicitly approved design revises them. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0119 is an archived historical decision cited by an adopted measurement. Preserve that historical evidence; it is not a new operative authorization or a reason to weaken current refusal behavior.
  ADR-0147 supersedes the historical refusal decision: preserve an adapter-origin refusal and use the catalog as the net, not as authority to claim a better message or successful selection.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `Makefile`, `.github/workflows/ci-verify.yml`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Measure current gate time by phase and execution context before choosing capacity changes.
- Preserve each failing diagnostic so reruns cannot erase the only causal evidence.
- Keep incremental reuse and complete fresh verification distinguishable.
- Make existing copied-lock analyzer coverage part of Verification after its defects are repaired.

## Core Features

1. Record comparable fresh baselines for local execution, loaded Agent/Daemon execution, and CI, identifying test/package parallelism, cache mode, and host conditions without copying secret environment values.
2. Capture full failing output and process ownership evidence before rerunning the legacy Force Stop and worktree-load cases. Unknown root cause remains an explicit outcome if the experiments cannot reproduce it.
3. Select a bounded capacity policy from those measurements; do not adopt a core-count formula or raise deadlines solely because the original failure disappears.
4. Preserve the direct unfiltered Go gate and existing two-tier semantics. Decide whether a named cold-cache convenience command adds value beyond current fresh test targets.
5. After the runtime-state prerequisite passes, add ordinary go-vet coverage to the approved Verification composition with a negative control that demonstrates it can refuse a copied-lock regression.

6. Characterize proposed controls for fixture deadlines, external-record lookups inside gates, fragile literals, test-mass infrastructure and toolchain/formatter drift. Record a measured adopt/defer decision for each proposal rather than treating every historical suggestion as an approved detector. No gate, cache or concurrency policy changes before its named cause and negative controls are established.

## Non-Goals / Out of Scope

- No arbitrary timeout extension, retry-until-green policy, diagnostic suppression, or weakened assertions.
- No machine-wide scheduler, new observability service, or new dependency in this proposal.
- No assertion that the historical saturation or ERANGE hypotheses are already proven.

## Acceptance evidence

Outside-evidence row: Replay the existing failure signature if it is observable and retain a typed unknown otherwise. Compare cold and cached runs under recorded conditions; neither a cached pass nor an unobserved failed assertion can establish a performance repair.

The later Task Graph and QA must prove each Core Feature with observed
positive and negative cases. Source inspection and proposed checks do not
constitute implementation or terminal QA evidence.

## Open Questions

- Approve the bounded characterization scope and the proposed Makefile/CI paths.
- Set an acceptable measurement window and probe-resource budget before load experiments.
- Approve any resulting parallelism or timing-policy change after its cause and before/after evidence exist; this proposal does not choose a number.

Until answered, all proposed limits and protected mutations remain unapproved.

## Source ownership

The maintainer selected this intent for implementation. Ordinary sources now
have one primary owner and one copy under that owner's `references/` directory.
Active Rollups remain as shared archive-license roots; their dated addenda map
every remaining family to its consuming Spec. Adoption is not execution approval.

- [2026-08-08-go-clean-testcache-clears-a-cache-the-gate-does-not-use.md](references/2026-08-08-go-clean-testcache-clears-a-cache-the-gate-does-not-use.md)
- [2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md](references/2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md)
- [2026-08-15-the-gate-runs-a-saturated-suite-inside-a-dense-run.md](references/2026-08-15-the-gate-runs-a-saturated-suite-inside-a-dense-run.md)
- [2026-08-11-a-git-worktree-that-fails-only-under-load.md](references/2026-08-11-a-git-worktree-that-fails-only-under-load.md)
- [2026-09-08-force-stop-legacy-owner-has-an-unexplained-ci-failure.md](references/2026-09-08-force-stop-legacy-owner-has-an-unexplained-ci-failure.md)

## Research basis

Secondbrain `inbox/roundfix/_triaged/2026-09-02-flake-em-force-stop-legacy-owner-so-em-ci.md` explains the lost diagnostic; `wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md` keeps a non-observation separate from a pass. Exa read [Go command documentation](https://pkg.go.dev/cmd/go) and returned the official [Go test command contract](https://pkg.go.dev/cmd/go/internal/test@go1.26.5): test-result reuse differs from execution, -count=1 disables that reuse, and -parallel governs one test binary rather than all package processes. This prevents treating one parallelism flag as a bound on total host load.

The local Secondbrain index was read and its query workflow used before
authoring. The sources above affected the stated requirements and limits;
they do not approve this Spec or substitute for live capability/behavior proof.

## Authoring checkpoint

Spec 0123 removes current copied-lock diagnostics before continuous analyzer gating. Spec 0119 owns experimental execution authority. The performance targets are pending measurements, not fabricated success thresholds.

Record the maintainer's bounded decision in [_authorization.md](_authorization.md),
then commit that approval record separately before its consuming tooling
changes. Only after that checkpoint may the TechSpec settle the design and a
Task Graph authorize execution. PRD-stage checks will report their actual
scope and any pending authorization findings; a green partial check is not
implementation readiness.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
