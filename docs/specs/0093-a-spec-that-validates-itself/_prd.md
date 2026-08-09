---
spec: 0093-a-spec-that-validates-itself
status: active
created: 2026-08-09
surfaces: [backend, cli]
---

# A Spec that validates itself

Spec 0090's QA gate returned `fail` on two defects. One of them was a PRD
claiming that ADR-0083 makes `make verify` the authoritative gate. ADR-0083 is
"Adopted sources move to their owning Spec" and says nothing about verification;
the rule exists, but in `docs/agents/specific-repository.md`. The Spec invented a
provenance, and the invention travelled into its TechSpec and the queue document.

Finding that cost an Agent turn. The gate spent 461 tool calls and compacted its
context twice — the only phase of that Run to exhaust context, against 161 tool
calls for the average implementation Task. Reading two files would have caught
it.

It gets worse on inspection. Of the gate's sixteen matrix rows, **eight audit
artifacts rather than product**: the authored graph, tooling authorization,
control declarations, rubric chronology, outside evidence, the Vocabulary
Contract, Non-Goals, and the report's own contract. Half of the most expensive
step in the loop reads documents. The finding that failed the gate came from one
of those eight.

And the diagnosis was already written. A finding filed on 2026-08-06 says:

> Citation checks prove that an obligation was named, not obeyed.

That predicted F-001 three days early. The same finding explains why it happened
anyway: *a finding records evidence; it does not prevent recurrence until its
defect class reaches a Task, test, or consistency check.* The knowledge existed.
Nothing executable carried it.

`roundfix spec check` already carries nineteen such rules and runs in **0.04
seconds** for one Spec. Running it during authoring on 2026-08-09 caught
`SC-ADR-UNLISTED` on Spec 0090 and `SC-ADR-RELATED` on Specs 0091 and 0092 —
three defects that would otherwise have reached a gate. The checker is not
missing. Its rules and its timing are.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity. Checks key off
  the existing Spec slug and artifact paths. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the checker reads repository files
  and reaches no network surface. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0096 establishes that the QA gate
  proves machine facts before it spends an Agent turn; this Spec is that
  principle carried to its conclusion, moving the machine facts out of the gate
  entirely. ADR-0091 keeps the authored QA gate a Task node of its own type, and
  this Spec narrows what that gate checks without removing it. ADR-0104 requires
  an acceptance row on evidence this Spec did not author. ADR-0083 appears above
  only as the worked example of a false attribution — this Spec neither obeys
  nor revises it, and the mention is what the existing listing check cannot tell
  apart from a citation. ADR-0081 draws authorization around the cause rather than its computable effects, which is why the copies `make skills-sync` rewrites follow the authorized skill edit. This
  Spec adds ADR-0116 and ADR-0117. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the checker itself lives in
  `internal/speccheck` and `internal/cli`, which are ordinary source, but wiring
  it into the authoring skills and removing the QA gate's governance rows are
  protected-tooling edits. Covered by the standing grant at
  `docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`,
  whose purpose bound this Spec matches exactly: the same rules run earlier and
  cheaper, and none is dropped. Bounded files:
  `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`, and their generated copies under `skills/`,
  which `make skills-sync` rewrites; ADR-0081 draws authorization around the cause rather than its computable effects, so they follow the authorized skill edit. The
  `write-tasks` wiring shipped ahead of this Spec under its own narrow grant.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Spec artifact that cites a decision is checked against what that decision
   says, not only against whether it was listed.
2. Every check that can be decided by reading repository files runs during
   authoring, in under a second, and blocks the stage that produced the defect.
3. The QA gate stops spending Agent turns on questions a file read answers.
4. A rule that a finding already established becomes executable rather than
   remaining prose.

## Core Features

- **Semantic citation checks.** An artifact stating that ADR-XXXX establishes Y
  is read against ADR-XXXX. A citation whose target does not carry the claim is
  a finding, named with both texts. This is the check that would have caught
  F-001 for the price of two file reads.
- **Stage-scoped validation.** The checker answers for one artifact at the
  moment it is authored, so the PRD stage sees PRD rules and the Task stage sees
  the whole assembled Spec. Each authoring stage ends by running it and refuses
  to report while a finding stands.
- **A gate that reads product, not paperwork.** Every rule decidable from files
  leaves the QA gate's matrix. What remains is what needs judgment: the Spec's
  goals exercised through the surfaces a user reaches, with captured evidence.

## Non-Goals / Out of Scope

- Judging whether a Spec's goals are the right goals. The checker decides
  internal consistency; a wrong-but-consistent Spec is a human problem.
- Replacing the QA gate. It keeps the acceptance no file read can substitute
  for, and ADR-0091 keeps it a Task node of its own type.
- Checks that need commits to exist. Auditing which paths a Task actually
  touched against its bounded list is a post-commit question and stays in the
  gate, as a command rather than a judgement.
- Rewriting existing Spec artifacts to satisfy new rules. Archived Specs stay
  byte-identical; active ones are corrected when their stage next runs.

## Decisions

- The checker owns the rules; the skills own the timing. A skill that says
  "validate" without an executable check is the shape that let a 2026-08-06
  finding fail to prevent a 2026-08-09 defect.
- A citation check reports both texts rather than guessing intent. The failure
  mode of a semantic check is a false accusation, and a maintainer settles that
  in seconds when both quotes are in front of them.
