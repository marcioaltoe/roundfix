---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
members:
  - 2026-08-07-changing-the-http-contract-discards-its-exceptions.md
  - 2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md
  - 2026-08-07-the-setup-refresh-interviews-a-repository-that-already-answered.md
  - 2026-08-07-two-http-contract-defaults-and-only-one-is-read.md
  - 2026-07-23-setup-context-driven-adoption-process-improvements.md
  - 2026-07-24-greenfield-agent-guidance-acceptance-target.md
  - 2026-07-26-baseline-profile-refresh-retention-gap.md
  - 2026-07-26-vortex-baseline-capability-remediation.md
  - 2026-07-27-bare-go-build-writes-an-untracked-root-binary.md
  - 2026-07-27-derived-skill-digest-pins-have-no-regeneration-path.md
  - 2026-07-27-sandboxed-agents-cannot-reach-the-default-go-cache.md
  - 2026-07-28-tooling-tasks-need-a-green-repo-and-an-undocumented-commit-choreography.md
  - 2026-07-29-doctor-requires-roundfix-own-development-skills.md
  - 2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md
  - 2026-08-01-characterization-corpus-is-outside-the-regeneration-command.md
  - 2026-08-05-what-this-repository-should-change-after-a-full-queue-night.md
---

# Baseline and derived tooling — one owner must reach every consequence (2026-08-06)

The Baseline and tooling findings describe one contract boundary: a declared
source change must preserve semantic guidance, identify every derived artifact,
and provide one sanctioned path that can regenerate the complete consequence
set. File hashes, scattered pins, and repository-specific recovery knowledge
do not provide that guarantee.

## Consolidated learning

- Baseline adoption must account for retained Normative Clauses and
  capabilities, not only file identity or a successful apply.
- A regeneration command must own the whole derivation chain, including the
  pins and characterization corpus that would otherwise block its own run.
- Tooling Tasks need explicit authority, a green precondition, deterministic
  local caches, and commit choreography that distinguishes prerequisite work
  from consequences of the authorized change.
- Readiness diagnostics must derive the Repository Skill Set and name the
  failed probe and a next action that can actually reach green.

## Live edge

Several member defects shipped through Specs 0054, 0057, 0061, 0062, and 0067.
The rollup remains `pending` because the accumulated evidence still asks for
one mechanically owned derivation path and semantic retention proof across the
whole Baseline lifecycle.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.

## Addendum — 2026-09-08 — Current triage

Existing deliveries address semantic retention (0057), repository-derived
skill requirements (0061), declared Baseline regeneration (0062 and 0067),
manifest-driven refresh (0082), and unrecorded managed-region refresh (0084).

Three member defects remain visible in current source. The `http-contract`
Change branch in `internal/cli/baseline_human.go` returns only `mode`, dropping
recorded exceptions and source provenance. The TypeScript Profile still names
an inert `Post-only` HTTP default beside the catalog's live default. Greenfield
classification returns without decisions while
`internal/baseline/preservation.go` still requires them for stale managed
sources. `skills/` also has no declared regeneration ownership, now captured in
`docs/backlog/2026-09-08-skill-regeneration-declares-its-owned-outputs.md`.

Route the residuals to provisional P2, Baseline decisions and complete
regeneration. Keep this Rollup active to license its 16 archived members;
historical absorption did not establish that every member defect was fixed.

## Addendum — 2026-09-08 — Complete implementation routing

The maintainer selected the residual work for the queue. Its primary owner is
[0121-baseline-decisions-and-complete-regeneration](../specs/0121-baseline-decisions-and-complete-regeneration/_prd.md).
- [0119-spec-contained-authorization](../specs/0119-spec-contained-authorization/_prd.md): operative authority.
- [0124-verification-capacity-and-measured-economics](../specs/0124-verification-capacity-and-measured-economics/_prd.md): measurement-driven tooling changes.

`done` records complete routing to implementation Specs, not a passing repair or
QA verdict. This Rollup remains in the active findings directory because its
archived members still name this basename as their absorption license. Those
licenses and original observations are preserved; retirement waits for the
durable replacement contract in 0120. Shipped mechanisms remain regression
obligations rather than duplicate implementation Tasks.
