---
spec: 0135-a-rename-the-gate-sees-from-both-sides
status: archived
created: 2026-09-14
surfaces: [backend]
archived: "2026-09-14"
source_slug: 0135-a-rename-the-gate-sees-from-both-sides
---


# A rename the gate sees from both sides

Spec 0134 made governed-mutation classification compare both directions, so a
Governed Path present before the Agent turn and absent after it is a mutation.
The pre-Pull-Request review of that work then proved the repair is necessary but
not reachable for a rename, and corrected a premise 0134 started from.

The worktree snapshot is not a listing of tracked files. It is the set of dirty
paths read from porcelain status, so a plain deletion already appears in the
after-snapshot and was already refused. The case that escapes is a rename: the
reader appends a record's destination and then skips the field carrying its
source, so `git mv Makefile notes.txt` yields a changed set naming only the
ungoverned destination. No comparison of snapshot pairs can recover a path that
never enters either snapshot, and 0134's own rename test could not catch this
because it injected a before-snapshot the real reader never produces.

Separately, resolving a record's location resolves a revision first, and a
failure there returns before the record is read. That early return reports an
unreadable record against the Spec Root field, so an unknown delivery target is
described as a missing Spec Root and sends the reader looking for the wrong
thing. The distinct reason for the condition already exists and is already
returned by the direct reader.

Both were found by the pre-Pull-Request review on 2026-09-14, against the
candidate carrying Specs 0119, 0132, 0133 and 0134.

## Project Constraints

- Identifier strategy: applicable — preserve the existing refusal codes, reason codes and record field names; no new identifier is introduced, and the repair consists of using a reason code that already exists. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic and is not narrowed here; ADR-0057 keeps the Daemon the exclusive writer of Task status. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon` and `internal/spec`, neither of which is a Governed Path, and their test files are ordinary too. Source: `docs/agents/agent-instructions.md`.

## Goals

- A rename of a Governed Path to an ungoverned name requires the same operation
  authority as changing it, proven through the reader the Daemon actually uses.
- An unresolvable delivery revision is reported as an unavailable revision
  naming that revision, so a stale target is not mistaken for a missing Spec
  Root.

## Core Features

1. The changed-path reader emits both the destination and the source of a
   rename or copy record, so a rename presents its governed side to the
   classifier that 0134 already made symmetric.
2. A governed rename is refused with the changed set the real reader returns,
   not with an injected snapshot pair.
3. Resolving a record's location reports an unresolvable revision with the
   existing unavailable-revision reason, naming the revision; a genuinely
   unresolvable Spec Root keeps its own reason and field.
4. Every refusal that works today keeps working with its existing token, and
   every outcome that is unresolved today stays unresolved.

## Non-Goals / Out of Scope

- The classifier's comparison, which 0134 settled and this Spec does not touch.
- External and symlinked Spec Root identity, which remains promoted to its own
  Spec with its own threat enumeration.
- Any change to the Governed Path set, to what a grant permits, or to the
  operation vocabulary.
- Repository-wide analyzer coverage, which Spec 0123 and Spec 0124 own.

## Declared intentional breaks

1. A Task that renames a Governed Path to an ungoverned name starts refusing,
   where it settled silently before.
2. Resolving a record against an unknown delivery revision starts reporting the
   unavailable-revision reason, where it reported an unreadable record against
   the Spec Root field before. The outcome stays unresolved either way.

Everything else keeps behaving as it does: refusal tokens, the operation
vocabulary, the Governed Path set, staging of a Task's own changes, and Daemon
ownership of Task status.

## Regression locks

- The reader's behavior is proven against porcelain output produced by a real
  rename in a repository the Task creates, never against a hand-written
  snapshot pair, because an injected pair is what hid this defect.
- The refusal fires on a positive observation — an observed governed path on
  either side of the record — never on missing or unreadable evidence.
- Both diagnostic reasons stay distinct, and the existing real-repository
  staging assertion keeps passing with the wider changed set.
