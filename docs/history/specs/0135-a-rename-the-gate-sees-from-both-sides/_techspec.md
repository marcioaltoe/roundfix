---
spec: 0135-a-rename-the-gate-sees-from-both-sides
status: active
created: 2026-09-14
---

# TechSpec — A rename the gate sees from both sides

## Executive Summary

Two readers report less than they found. The changed-path reader discards a
rename record's source, so the governed side of a rename never reaches the
classifier 0134 made symmetric. The location resolver reports a failed revision
resolution as an unreadable record against the Spec Root field, so a stale
delivery target is described as a missing Spec Root.

Both repairs are local and additive: emit the field already present in the
record, and use the reason code that already exists. Neither changes what a
grant permits, and neither narrows the Governed Path set.

## Project Constraints

- Identifier strategy: applicable — preserve the existing refusal codes, reason codes and record field names; the revision repair consists of using a reason code that already exists. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic and is not narrowed here; ADR-0057 keeps the Daemon the exclusive writer of Task status. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon` and `internal/spec`, and their test files are ordinary too. Source: `docs/agents/agent-instructions.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| Changed-path reader | Turn porcelain status into the set of paths the turn touched | Emit a rename or copy record's source alongside its destination |
| Governed-mutation classifier | Decide whether the turn mutated a Governed Path | None; it already compares both directions |
| Record location resolver | Resolve the repository, revision and path of a record before reading it | Report an unresolvable revision with the reason that already exists |

## Implementation Design

### A rename names both of its paths

A porcelain rename record carries the destination in its first field and the
source in a second. The reader appends the destination and then advances past
the source, so the source is discarded before any classifier sees it. Emitting
both states the truth about a rename for every reader of the changed set: both
paths did change, one by disappearing.

This is the layer the repair belongs to. The alternative — teaching the
classifier to parse rename records — would give the changed set two sources of
truth about the same event, which is how this class of defect returns.

The wider changed set reaches staging as well. Staging both sides of a rename is
what Git requires to record the rename at all, so the existing real-repository
staging behavior is a regression lock rather than a risk.

### An unresolvable revision says so

Resolution resolves the revision before the record is read, and its failure
path returns early through a helper that always reports an unreadable record
against the Spec Root field. The reader already owns a distinct reason for an
unavailable revision and already returns it when called directly with a bad
revision, which is why the existing test passes while the composed path
misreports. Use that reason on this branch, naming the revision that could not
be resolved.

This is a diagnostic contract, not an authority change: the outcome is
unresolved before and after, so nothing becomes permitted.

## Coverage Map

- PRD Goal 1 → Changed-path reader, Governed-mutation classifier.
- PRD Goal 2 → Record location resolver.
- Core Feature 1 → Changed-path reader.
- Core Feature 2 → Changed-path reader, Governed-mutation classifier.
- Core Feature 3 → Record location resolver.
- Core Feature 4 → both, by leaving every refusal token and every unresolved
  outcome intact.

## Testing Approach

Focused tests at the two existing seams; no new seam is needed.

1. A rename performed in a repository the test creates yields both of its paths
   from the reader.
2. A governed rename to an ungoverned name is refused when the changed set
   comes from the real reader rather than from an injected snapshot pair.
3. An ordinary removal and the real-repository staging behavior are unchanged
   by the wider changed set.
4. An unresolvable delivery revision reports the unavailable-revision reason
   naming the revision, and an unresolvable Spec Root keeps its own reason.
5. Every refusal that works today still works with its existing token.

Observation 1 rests on evidence this Spec did not author: Git's own porcelain
v1 status format, which places a rename record's source in a second
NUL-separated field. The assertion reads what Git actually emits for a rename
the test performs, rather than a string this Spec composed. If Git cannot be
run, the row records blocked with that reason.

## Build Order

1. Carry a rename's source through the changed-path reader, proven against
   real porcelain output (depends on: none).
2. Report an unresolvable revision as an unavailable revision (depends on:
   none).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

- The wider changed set is read by staging as well as by classification. The
  existing real-repository staging assertion is the lock that proves staging
  still stages exactly the Agent's own changes.
- A rename whose source is ungoverned and whose destination is governed was
  already refused through the destination. Emitting both sides keeps that
  refusal and adds the mirror case.

## Decisions

- Repair the reader, not the classifier, so the changed set keeps one source of
  truth about what a turn touched.
- Prove the reader against Git's own output. An injected snapshot pair is what
  allowed 0134's rename test to pass while the pipeline stayed broken, so the
  premise of the new test is that the boundary produces the input.

## Vocabulary Contract

No token is coined. Both repairs reuse existing reason codes and refusal
tokens, and their glossary owners are unchanged.
