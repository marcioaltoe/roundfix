---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
absorbed_by: 0123-runtime-readiness-and-model-capabilities
members:
  - 2026-08-07-a-sixty-four-value-bound-locks-out-the-opencode-runtime.md
  - 2026-08-10-a-fake-adapter-goes-silent-under-a-dense-start.md
  - 2026-08-14-preflight-starves-when-the-machine-is-busy.md
  - 2026-07-26-claude-adapter-configoptions-migration.md
  - 2026-07-27-claude-adapter-standardization.md
  - 2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md
  - 2026-08-04-a-doctor-next-action-that-does-not-reach-green.md
  - 2026-08-04-the-documented-sandbox-escape-does-not-exist-on-the-codex-adapter.md
  - 2026-08-05-agent-full-access-passes-config-validation-and-fails-every-task.md
---

# Agent selection and execution environments — preflight must prove the Task's reality (2026-08-06)

These findings converge on the same failure mode: configuration can validate
while the selected adapter, model controls, or sandbox policy fails when the
Agent begins real work. Readiness is useful only when it exercises the same
selection and access contract the Task Session will receive.

## Consolidated learning

- Adapter lineage, advertised model controls, and model-name parsing belong to
  one selection proof; package presence alone is not capability evidence.
- Profile updates must merge the named category without deleting unrelated
  selections.
- Access-policy names must map to supported adapter behavior. A documented
  escape hatch or `agent_full_access` value that the runtime cannot honor is a
  preflight defect, not a Task failure.
- Doctor next actions must repair the failing predicate and then re-check the
  same predicate.

## Live edge

Specs 0052 and 0056 absorbed the first adapter and profile defects. The rollup
remains `pending` around end-to-end execution-policy proof: accepted config,
runtime launch, and Task filesystem access must describe the same environment.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.

## Addendum — 2026-09-08 — Current triage

The nine member Findings remain preserved under this active license. Existing
deliveries cover the Claude adapter and opaque identifiers (0052), profile
merge semantics (0056), bounded capability projection (0088), runtime effort
application (0089), selection proof (0091), and compiled fixtures (0103).

Current `internal/agent/acpx_runner.go` still applies full-access mode when
opening the work Session; this review found no Doctor proof of that effective
access policy. `internal/cli/profiles_validate.go` still includes every
configured fallback in preflight proof. The historical load-starvation reports
have no new reproduction in this triage. Route these residuals to provisional
P4, runtime readiness and model capabilities, and P5, Verification capacity and
measured economics. Regeneration ownership is shared with provisional P2.

This addendum records implemented scope from current source and archived
delivery records, not a fresh live-adapter test or complete resolution.

## Addendum — 2026-09-08 — Complete implementation routing

The maintainer selected the residual work for the queue. Its primary owner is
[0123-runtime-readiness-and-model-capabilities](../../specs/0123-runtime-readiness-and-model-capabilities/_prd.md).
- [0124-verification-capacity-and-measured-economics](../../specs/0124-verification-capacity-and-measured-economics/_prd.md): resource contention and characterized fallbacks.
- [0121-baseline-decisions-and-complete-regeneration](../../specs/0121-baseline-decisions-and-complete-regeneration/_prd.md): generated readiness remedies.

`done` records complete routing to implementation Specs, not a passing repair or
QA verdict. This Rollup remains in the active findings directory because its
archived members still name this basename as their absorption license. Those
licenses and original observations are preserved; retirement waits for the
durable replacement contract in 0120. Shipped mechanisms remain regression
obligations rather than duplicate implementation Tasks.

## Addendum — 2026-09-08 — Routing document removed

The maintainer requested removal of `docs/workflow/` and its routing documents.
The earlier citation remains a dated historical observation; its original bytes
can be read at Git revision `6b8ea48725cbca13974eee0b400b3482202874f6`.
The current primary and secondary Spec owners remain those in the complete-triage
addendum above; removing the old plan does not reopen or erase the members.

## Addendum — 2026-09-08 — Archived after complete routing

The maintainer requires terminal Findings and Rollups to leave the active
family directory. Every member now points directly to an existing active or
archived Spec, and this Rollup has its own direct Spec absorber. The earlier
statements retaining this file as an active license root are superseded by
this completed routing migration. Original observations and prior pointers
remain recorded; archival does not claim implementation of pending Specs.
