---
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: active
created: 2026-08-14
surfaces: [backend, docs]
---

# A gate report that does not block its successor

The QA gate refuses before spending an Agent turn when a machine fact says it
should, which is right and cheap. What it then writes is a report its own
contract calls malformed, and every later run of the same Spec reads that report
and refuses on it — with a prescribed fix that is impossible, because the run
never built a matrix to materialize rows from. It happened twice in one Spec, and
the only exit was deleting evidence, which is the one move this repository's rules
single out. A second, independent refusal in the same family names a row's blocked
cause as untyped when it is typed correctly: what the parser wants is a literal
the diagnostic never mentions, and the mismatch produces a second symptom that
sends a reader hunting a counting bug that does not exist.

## Project Constraints

- Identifier strategy: applicable — QA Report, the verdict vocabulary and the
  three typed blocked-cause counts are glossary terms this Spec changes the
  reading and writing of. The closing node checks whether the work introduced or
  changed a term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential or
  request is created or read. The work is report writing, report parsing and the
  diagnostics they emit. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0096 makes the gate prove machine
  facts before spending an Agent turn, which is the behaviour that produces the
  refusal this Spec repairs rather than removes; ADR-0080 owns verdict semantics
  and the three typed blocked-cause counts, which this Spec reads and must not
  redefine; ADR-0097 lets a row carry forward only on declared, unmoved evidence,
  which bounds what a report may claim about a prior run. The decisions built on
  those are accounted: ADR-0091 makes the gate a Task node of its own type,
  ADR-0093 checks Spec consistency by citation rather than inference, ADR-0104
  makes a Spec accept on evidence it did not author, and ADR-0117 places a check
  with the stage that can produce its defect — this Spec changes none of the four,
  and the last is why its repairs belong in the stage that writes and reads the
  report rather than in a later gate round. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the gate's report writer and
  mechanical stage plus their tests, and the `qa-gate` contract's stated row form
  is read rather than edited. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A refused gate leaves nothing that blocks its own next run.
2. A refusal reports what it actually needs, in the words the parser accepts.
3. No repair for a gate defect requires deleting evidence.

## User Stories

1. As a maintainer whose gate refused at a precondition, I want the next run to
   judge my Spec rather than the previous run's report, so that fixing the
   precondition is enough.
2. As a maintainer reading a report-shape refusal, I want it to name the literal
   the parser requires, so that I do not rewrite a row that was already typed
   correctly.
3. As a maintainer, I want a blocked-row count that disagrees with the table to
   name the rows it failed to parse, so that I do not look for a counting defect
   that is a parsing symptom.

## Core Features

1. **A precondition refusal writes a report its contract accepts.** Either it
   writes none, or it writes one whose Results table carries a terminal row
   recording the refusal — which is more useful than an empty table and is what
   actually happened.
2. **The mechanical stage reads the current report.** A superseded report does
   not refuse the run that supersedes it.
3. **A shape refusal names the form it wants.** When a row carries the right
   type and the wrong phrasing, the diagnostic names the missing literal rather
   than repeating the three type prefixes the row already satisfies.
4. **One cause reports as one finding.** A count that disagrees with the table
   because rows failed to parse is reported as that parse failure, not as a
   separate counting defect.
5. **A row typed correctly is not refused for prose.** Whether the parser should
   require its declared phrasing at all, or accept the type and treat what follows
   as free text, is settled rather than left to the parser's current behaviour.

## Non-Goals / Out of Scope

- Removing or weakening the mechanical stage. It refuses before spending an Agent
  turn, which is the behaviour that made both defects cheap to find.
- Changing verdict semantics or the three typed blocked-cause counts.
- Changing what the gate executes once its preconditions pass.
- Retroactively repairing reports already written.

## Success Metrics

- A gate that refuses at a precondition, then has that precondition fixed, passes
  on the next run without any report being deleted.
- The two measured refusals — an empty Results table and a row missing its
  declared literal — each produce a diagnostic naming what to change, proven by
  replaying both.
- A parse failure reports once rather than as two findings.

## Decisions

- The refusal path is repaired rather than the refusal removed. A gate that stops
  early on a machine fact is doing what it was built for; what it leaves behind is
  the defect.

## Open Questions

- Whether a precondition refusal should write a report at all. Writing a valid one
  keeps the evidence trail complete; writing none removes the whole class of
  successor-blocking. The design settles which, and the second is simpler.
- Whether the declared row phrasing is a contract or a convention. If it is a
  contract, the parser is right and only the message needs work; if a convention,
  the parser should accept the type alone.
