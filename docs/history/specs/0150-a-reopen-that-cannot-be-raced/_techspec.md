---
spec: 0150-a-reopen-that-cannot-be-raced
status: active
created: 2026-09-19
surfaces: [backend, cli, docs]
---

# A reopen that cannot be raced

## Executive Summary

Two repairs to the command Spec 0149 shipped. The plan is rechecked against the
tree immediately before the replacement, and the atomic replacement follows a
symlinked Task path to its target instead of replacing the link.

## Project Constraints

- Identifier strategy: applicable — no identity is minted or renamed. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads and writes only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0155, ADR-0104 and ADR-0156 hold;
  this Spec repairs an operation under them. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the new exit-2 refusal is public CLI behavior.
  Express maintainer authorization: granted 2026-09-18 as a standing grant,
  consumed in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## The recheck

`preflightReopen` returns a plan: the QA Task path, its id, the report label and
the stale dependency ids. Between that read and `ReopenGate`'s rename, nothing
holds the tree.

The repair re-derives the plan immediately before the write and compares it to
the one preflight produced. The comparison is on what the record will claim: the
QA Task id, and the stale dependency set. If the gate is no longer stale, or the
set differs, the command refuses.

Re-deriving through the same recovery load the preflight uses keeps one
predicate. A recheck implemented differently from the check is a third thing
that can disagree with the gate.

The refusal exits 2, like every other condition the command declines on, and
names what changed.

## The symlinked path

`replaceTaskFile` creates its temporary file in `filepath.Dir(path)` and renames
onto `path`. When `path` is a symlink, the rename replaces the link.

`SetStatus` used `os.WriteFile` before Spec 0149 and wrote through the link, so
this is a regression that Spec introduced while making the write atomic. The
repair resolves the final path before choosing the temporary file's directory
and the rename target, so the write lands on the file the link points at and the
link survives.

Resolution stays inside the confinement the command already enforces: Spec 0149
checks the Task path lexically against the Spec directory and then checks the
resolved path against the configured Spec Root. Both checks stay.

## API Contracts

1. A reopen whose gate stopped being stale between preflight and the write
   refuses with exit 2 and writes nothing.
2. A reopen through a symlinked Task path rewrites the target and leaves the
   link a link.
3. Every refusal and every success Spec 0149 shipped is unchanged.

## Coverage Map

- Goal 1 → The recheck; API Contract 1.
- Goal 2 → The symlinked path; API Contract 2.
- Core Feature 1 → The recheck.
- Core Feature 2 → The recheck (refusal); API Contract 1.
- Core Feature 3 → The symlinked path; API Contract 2.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → The recheck, The symlinked path.

## Integration Points

- **`spec.ReopenGate` and `replaceTaskFile`.** The replacement seam both repairs
  land on.
- **Spec 0149.** Its refusals, its read-only active-Run guard and its
  confinement checks stay exactly as shipped.

## Testing Approach

1. **The gate changes underneath.** A fixture whose stale dependency is
   completed between the preflight read and the write is refused with exit 2,
   and its QA Task file is byte-identical afterwards. Fails on the tree as it
   stands, where the command writes from the stale plan.
2. **The symlinked Task path.** A fixture whose Task path is a symlink to a file
   inside the Spec Root has the target rewritten to `pending` with the
   invalidation recorded, and the path is still a symlink afterwards.
3. **Regression.** The reopen tests Spec 0149 shipped pass unchanged.
4. **The skill is true.** The shipped skill and its mirror name the recheck
   refusal, and the mirror matches the canonical file.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The recheck and its refusal, with tests (depends on: none).
2. The symlink-following replacement, with tests (depends on: none).
3. The shipped skill (depends on: 1).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **A recheck that never fires.** A test that only asserts the happy path would
  pass with the recheck deleted. The refusal case is a Success Metric for that
  reason.
- **Resolution that escapes confinement.** Following a link to write is only
  safe because the resolved path is already checked against the configured Spec
  Root. That check is a precondition of this repair, not an independent one.
- **A narrowing window, not a closed one.** The recheck shrinks the race to the
  interval between the recheck and the rename. Closing it entirely needs a lock,
  which this Spec declines as disproportionate and says so in its Non-Goals.
