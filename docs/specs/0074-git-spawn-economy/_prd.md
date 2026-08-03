---
spec: 0074-git-spawn-economy
status: active
created: 2026-08-03
surfaces: [backend, test]
---

# Git spawn economy

This is an engineering-framed minimal PRD: the Spec is a performance refactor
with no product behavior change, entered through the TechSpec per
`docs/agents/spec-routing.md`.

## Problem

Verification was measured, not estimated: the full fresh suite runs at 36%
CPU utilization with kernel time three times user time, and a PATH shim
counted **13,926 git invocations in one run**. Roughly six thousand of those
are issued by Roundfix itself — `rev-parse` 3,518, `ls-tree` 984, `cat-file`
971, `rev-list` 829, `status` 349, `for-each-ref` 220 — because production
code reads repository facts one subprocess per fact and, in two loops, one
subprocess **per file**. Every real Run pays the same pattern against the
user's repository.

Test-side waste has already been removed (Spec 0071 and its follow-ups):
fixtures share builds and templates, 668 tests were parallelized, and
deleting the twelve heaviest tests was measured to buy only seven to
fourteen seconds. The floor that remains is production-issued spawns.

The maintainer's directive: a complete fresh execution — no cache of any
kind — under **60 seconds**. It stands at 67–70s.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier
  is created; existing command names, error codes, and package paths keep
  their contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no transport, credential, or
  deployment surface is touched; the governing clause prohibiting invented
  policy is unaffected. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0089 (code under test takes its
  environment explicitly) is extended to `internal/agent`; ADR-0090 (created
  by this Spec) governs batched object reads. Digest-bearing assets are not
  edited, so ADR-0081 regeneration is not triggered. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling configuration,
  workflow, skill, or version pin changes; the work lives in production Go
  and its tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. Production reads repository facts in batches: subprocess count
   proportional to **operations**, never to file count.
2. `internal/agent` takes its process environment explicitly per ADR-0089,
   so the 121 tests that today mutate process env may declare parallelism.
3. The complete fresh suite (`go test ./... -count=1`) completes under 60
   seconds on the reference machine, measured with the same commands before
   and after.
4. Observable behavior is unchanged: same outputs, same errors, same digests.

## Core Features

1. The two per-file `cat-file blob` loops (skills restore, assets sync
   provenance) read all objects through one `git cat-file --batch` process.
2. Repository resolution combines its multi-fact queries into single
   invocations where git's own interface allows it (`rev-parse` accepts
   multiple queries per call).
3. The ACPX child environment is composed from an explicit base supplied at
   the boundary, with the process environment as the default — the exact
   shape Spec 0071 gave `internal/cli`.
4. A committed spawn census (the shim procedure and its counts) taken before
   any change, and a before-and-after published under the Spec.

## Non-Goals

- No test is deleted, skipped, or weakened — measurement already showed
  removal cannot reach the target.
- No CLI surface, output, or exit-code change.
- No caching of repository facts across mutation boundaries; batching only
  within scopes where the repository state is provably immutable.
- No shared git client extraction: the seven per-package runners stay where
  they are. Batching lands inside the loops that need it.
