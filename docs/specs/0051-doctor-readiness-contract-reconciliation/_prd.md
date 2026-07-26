---
spec: 0051-doctor-readiness-contract-reconciliation
status: active
created: 2026-07-26
surfaces: [cli, backend, docs]
---

# Doctor Readiness Contract Reconciliation

A follow-up review found that the Doctor Command's Repository Skill Set check
can outlive cancellation, can inherit the wrong working-directory behavior
from repository-root handling, and can print remediation that is ambiguous
when both ownership groups fail. The shared external-skill hash also leaves
collation-equal paths dependent on input order, while the newly added Go module
metadata is not tidy. This correction makes those contracts explicit and
deterministic without rewriting the branch or revisiting the implementation
history already accepted by the maintainer.

## Project Constraints

- Identifier strategy: not applicable — this correction creates no
  project-owned Internal Identifier; existing skill names and
  repository-relative paths remain the only identities. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all behavior remains local,
  offline, and filesystem-only, with no authentication, authorization,
  transport, or HTTP change. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 require Agent
  Selection Profile proof to remain independent and exact; ADR-0066 and
  ADR-0072 keep shared Baseline hash consumption in the Go CLI; ADR-0077
  requires this readable Project Constraint snapshot. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization was given
  on 2026-07-26 to update `go.mod` and `go.sum` through Go tooling for the
  `golang.org/x/text` collation dependency, and to keep
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` aligned
  with shipped CLI behavior. These four paths are the complete protected
  tooling boundary for this Spec; no other protected tooling mutation is
  authorized. Source: `docs/agents/agent-instructions.md`.

## Goals

- Make every blocking Repository Skill Set filesystem operation honor the
  Doctor Command's cancellation context.
- Give collation-equal but byte-distinct paths one deterministic total order.
- Preserve the pre-existing Agent Selection Profile working-directory
  behavior while keeping Repository Skill Set inspection anchored only to the
  resolved Git root.
- Use canonical Repository Skill Set terminology and fail-closed,
  ownership-specific remediation.
- Restore Go module metadata to the state produced by Go tooling.

## User Stories

1. As an operator cancelling Doctor, I want Repository Skill Set inspection to
   stop cooperatively, so a large repository does not keep Doctor alive after
   cancellation.
2. As a maintainer hashing Unicode paths, I want equivalent collation keys to
   have a deterministic tie-breaker, so the digest does not depend on input
   order.
3. As a developer running Doctor outside Git, I want profile proof to retain
   its usable process working directory while Repository Skill Set readiness
   reports that a Git repository is required.
4. As an Agent reading a mixed-ownership failure, I want one fail-closed shell
   command chain, so later remediation does not run after an earlier action
   fails.
5. As a maintainer reviewing dependency metadata, I want `go mod tidy` to
   produce no diff, so direct and indirect requirements reflect real imports.

## Core Features

1. **Context-aware inspection.** Public Repository Skill Set entry points take
   `context.Context` first, propagate it through repository walks and reads,
   and return cancellation or deadline errors without losing error-chain
   inspection.
2. **Total external hash order.** External-skill hashing keeps American
   English collation as the primary ordering contract and uses normalized path
   bytes to break equality between distinct paths.
3. **Separated working-directory authority.** Agent Selection Profile proof
   uses the resolved Git root when present and the process working directory
   otherwise. Repository Skill Set inspection uses only a non-empty resolved
   Git root and never substitutes the process working directory.
4. **Deterministic failures.** A symlinked external lock retains external
   ownership, missing-root detail uses the canonical Repository Skill Set
   term, and multiple remediation commands form one `&&`-joined fail-closed
   chain.
5. **Tool-produced module metadata.** Go tooling moves the text dependency to
   the direct requirement set and removes stale checksums without introducing
   another dependency or changing versions.
6. **Synchronized guidance.** User documentation and the canonical and shipped
   Roundfix Skill pair describe the exact Doctor output and remediation
   contract.

## User Experience

Doctor keeps its current ordered, one-line-per-check text interface. Outside a
Git repository, profile proof still receives a real working directory and the
skills line reports that Repository Skill Set readiness requires a Git
repository. When both owned and external repair actions are needed, `next:`
contains one `&&`-joined command chain in owned-then-external order.

Cancellation remains a command-level failure signal. Doctor does not install,
update, write, or execute any remediation automatically.

## Non-Goals / Out of Scope

- Rebuilding, rebasing, force-pushing, dropping, or otherwise rewriting any
  existing branch commit.
- Editing any archived Spec, including Specs 0036 and 0050.
- Reverting or editing the accepted upstream-managed skill update,
  `skills-lock.json`, or `skills/recommended.txt`.
- Changing Baseline assets, generated Baseline parity artifacts, or the
  repository-owned digest preservation behavior already implemented.
- Editing `.coderabbit.yaml` or `.roundfixrc.yml`.
- Moving Repository Skill Set readiness into `HealthChecker`; the existing
  injected Doctor seam remains the ownership boundary.
- Making owned and external skill-root symlinks share one classification.
  Both remain blocking, but preserve their ownership-specific contracts.
- Adding a Doctor output schema, new flag, network call, or automatic repair.

## Success Metrics

- A pre-cancelled context and a context cancelled during skill traversal both
  stop inspection and remain detectable with `errors.Is`.
- Every permutation of a corpus containing collation-equal Unicode paths
  produces the same external-skill digest.
- Outside Git, profile proof receives the process working directory, the
  repository checker is not invoked, and Doctor prints canonical Repository
  Skill Set guidance.
- A symlinked lock produces only external remediation; a mixed failure prints
  one owned-then-external chain joined by `&&`.
- `go mod tidy -diff` is empty after implementation.
- Focused tests, affected-package race tests, public Doctor execution, and the
  repository verification gate pass without mutating checked repository
  state.

## Decisions

- Preserve all accepted branch history and correct only the newly confirmed
  contract gaps.
- Keep American English collation for compatibility, but define normalized
  byte ordering as the deterministic tie-breaker that the external CLI leaves
  unspecified.
- Propagate the command context through the existing Doctor dependency seam
  rather than adding a second cancellation mechanism.
- Restore the process working-directory fallback only for Agent Selection
  Profile proof, never for Repository Skill Set inspection.
- Keep ownership-specific symlink behavior and the existing `HealthChecker`
  boundary unchanged.
- Add no ADR because this Spec narrows existing cancellation, hashing, and CLI
  contracts instead of creating a new architectural boundary.

## Open Questions

None.
