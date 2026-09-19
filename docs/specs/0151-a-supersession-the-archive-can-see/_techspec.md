---
spec: 0151-a-supersession-the-archive-can-see
status: active
created: 2026-09-19
surfaces: [backend, cli, docs]
---

# A supersession the archive can see

## Executive Summary

A command records, in a Spec whose content another Spec delivered, which Spec
delivered it and why. Archive reads that record and accepts it in place of the
completed Tasks and passing gate such a Spec can never have, while every other
precondition it enforces stays in force.

## Project Constraints

- Identifier strategy: applicable — both Specs are named by slug and keep their
  identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads and writes only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0155, ADR-0104 and ADR-0156 hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — a new command and a changed archive
  precondition are public CLI behavior. Express maintainer authorization:
  granted 2026-09-18 as a standing grant, consumed in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Where the record lives

The amendment goes in a file of its own, `_supersession.md`, beside the Spec's
other underscore-prefixed records. Two reasons.

It keeps the PRD untouched, which is what "preserved prior Results and evidence"
asks for: the document that stated the Spec's intent still states it, unedited,
and a reader comparing intent to outcome sees both.

It gives archive something to read that cannot be confused with prose. A
paragraph inside a PRD would have to be parsed out of narrative text; a file
with frontmatter is a record.

Its frontmatter carries the superseding slug, the date and the reason. Its body
carries whatever explanation the command was given.

## The command

```
roundfix supersede --spec <slug> --by <slug> --reason <text>
```

- **Exit 0.** `_supersession.md` is written. Nothing else in the Spec changes.
- **Exit 2.** Preflight Validation failed: the superseded Spec is unknown, the
  superseding slug is neither active nor archived, the two slugs are the same,
  or the Spec already carries a supersession. The reason names which condition
  failed and nothing is written.

It creates no Run, writes no Run Event Journal entry, commits nothing and pushes
nothing — the envelope `settle` and `reopen` already share.

## What archive changes

Archive's completed-Tasks-and-passing-gate evidence becomes one of two accepted
proofs:

- the existing proof, for a Spec with a Task Graph; or
- a recorded supersession, for a Spec without one.

A Spec that has a Task Graph is unaffected even if a supersession is present:
its own evidence rules still decide. The supersession path exists for the Spec
that was never decomposed, which is the only Spec that cannot satisfy the
existing rule.

Every other archive precondition — self-containment, destination, the rest —
applies unchanged to both paths.

## API Contracts

1. `supersede` accepts `--spec`, `--by` and `--reason`, refuses unknown flags,
   and refuses each named condition with exit 2.
2. Archive succeeds for a Spec with no Task Graph and a recorded supersession,
   and still refuses that Spec without one.
3. Archive's behavior for a Spec with a Task Graph is unchanged.

## Coverage Map

- Goal 1 → What archive changes; API Contract 2.
- Goal 2 → Where the record lives.
- Goal 3 → Where the record lives; Testing Approach 3.
- Core Feature 1 → Where the record lives; The command.
- Core Feature 2 → What archive changes; API Contract 2.
- Core Feature 3 → The command (refusals); API Contract 1.
- Core Feature 4 → What archive changes; API Contract 3.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 1.
- API Contracts 1-3 → The command, What archive changes.

## Integration Points

- **`roundfix archive`.** The one precondition that gains a second accepted
  proof.
- **Spec 0128.** The Spec this exists for: its content shipped in Spec 0147 and
  it has been stuck active ever since.
- **Spec 0129.** Premise falsification and consumer revalidation stay there.

## Testing Approach

1. **The refusals.** An unknown superseding slug, a Spec superseding itself, and
   a second supersession on an already-amended Spec are each refused, and the
   Spec directory is byte-identical afterwards. Fails on the tree as it stands,
   where the command does not exist.
2. **The accepted path.** A fixture with no Task Graph and a recorded
   supersession archives; the same fixture without one refuses with the
   missing-Task-Graph reason.
3. **Preservation.** After the amendment, every file the fixture already had is
   byte-identical, and only `_supersession.md` is new.
4. **The unchanged path.** A fixture with a Task Graph keeps archive's existing
   behavior, with and without a supersession present.
5. **The skill is true.** The shipped skill and its mirror describe the command
   and the changed archive precondition, and the mirror matches the canonical
   file.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The supersession record and the command that writes it, with its refusals
   and tests (depends on: none).
2. Archive's second accepted proof, with tests for both paths (depends on: 1).
3. The shipped skill and the user guide (depends on: 1, 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **An archive that stops checking.** If the supersession path skipped archive's
  other preconditions, the record would become a way around the gate rather than
  a different proof for it. Core Feature 4 and Testing Approach 4 are the
  controls.
- **A record nobody can trust.** A supersession naming a Spec that does not
  exist is worse than no record, because it reads as an answer. The existence
  check is a refusal, not a warning.
- **Scope creep into falsification.** Recording "this shipped elsewhere" and
  recording "this premise was wrong" are different claims with different
  evidence. Only the first is here.
