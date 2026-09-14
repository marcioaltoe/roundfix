---
spec: 0138-a-matrix-the-spec-declares
prd: _prd.md
created: 2026-09-14
---

# A matrix the Spec declares — Technical Spec

## Executive Summary

The change touches two texts and no running code path: the qa-gate skill, which
the gate's Agent follows, and the QA contract the Daemon places in that Agent's
prompt. Today the contract tells the gate to validate every user story and
acceptance criterion. Both texts change to say three things:

- **The matrix is declared.** Each `qa` Task Requirement that starts with
  `MUST verify` or `MUST run` is one row. A Task without such a Requirement gets
  a bounded default.
- **Obligations stay.** The outside-evidence row, the Pull Request row, the
  repository Verification and the frontend sweep cannot be waived.
- **Blocking is scoped.** Once the matrix exists, a finding blocks only the rows
  that depend on it.

Verdict rules, report keys, the mechanical stage and its withholding of the Agent
Session all stay as they are.

The main trade-off is that the matrix now covers only what the Spec wrote down,
so a promise left out of the `qa` Task goes unchecked. The maintainer accepted
that risk and assigned the strengthening of authoring to Spec 0129.

A second trade-off is that the gate keeps repeating the changed-path audit. The
mechanical stage receives only the current Run's Task commits, and its report
does not name the commits it audited, so the gate cannot know what to skip. That
gap goes to follow-up rather than into this Spec.

## Project Constraints

- Identifier strategy: applicable — the QA Report frontmatter keys, the three typed blocked causes, the Mechanical Refusal Codes and matrix row identifiers keep their spelling and meaning; this Spec introduces no identifier, refusal code or report key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request or HTTP surface is created or read; the change is guidance for the gate executor and the text of its prompt. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's matrix, verdicts and evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and every row still ends in one of the typed blocked causes or a result.
  ADR-0088 applies: the QA gate is authored into the graph, which is where the declared matrix is read.
  ADR-0091 applies: the QA gate is a Task node of its own type, so the matrix is declared by the `qa` Task.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, so the stage still withholds the Agent Session on a blocking fact, no row re-checks a fact the stage decides, and no stage makes a verdict more permissive.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, so row inputs are fixed when the row is planned.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row stays in every matrix.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the commit-dependent audit stays a gate fact executed as a command.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, which this Spec writes into the qa-gate skill and the QA contract.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule or citation reader.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-14, recorded in `docs/specs/0138-a-matrix-the-spec-declares/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved skill edit moves. The QA contract in the gate prompt is ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| qa-gate skill, canonical copy | Tells the gate Agent how to plan, run and close the matrix | Sections 1, 2, Row input declaration and 3 take Core Features 1-8 |
| qa-gate skill, distributed mirror | Embedded in the binary and installed into repositories | Regenerated from the canonical copy, never hand-edited |
| QA contract in the gate prompt | The rule summary the Daemon hands the gate Agent | Its matrix bullet takes Core Feature 9 |
| Mechanical stage | Computes machine facts and withholds the Agent Session on a blocking fact before a matrix exists | None |
| Workflow contract test for the skill | Pins the tooling-audit clauses the gate must keep | None; every clause it pins stays |

The gate Agent reads the prompt first and the skill second, so the two must say
the same thing. The prompt stays short and points at the skill for detail, as it
does now. It is built only when the mechanical stage did not withhold the Agent,
which is why scoped blocking applies to a planned matrix and never overrides the
stage.

## Implementation Design

### Where each Core Feature lands in the skill

**Section 1, scope.**

- Replace the closing sentence that turns every promise and explicit exclusion
  into a planned row. The scope is now complete when every declared or default
  source has its row.
- State that declaring a matrix does not make the run partial (Core Feature 1).
- In the tooling-audit bullet, delete only the sentence that stops flow QA after
  any audit problem, and point instead to finding-dependent blocking (Core
  Feature 5). The bullet keeps its behavior: the gate still audits Task commits
  by command, and every clause the contract test pins stays verbatim.

**Section 2, rows.** Replace the "Add a row for" list with the declaration rule
(Core Feature 1):

- Each numbered `qa` Task Requirement that starts with `MUST verify` or
  `MUST run` is exactly one row.
- A `qa` Task with at least one such Requirement has declared its complete
  matrix.
- Other Requirements constrain how the gate runs.

Add the bounded default for an undeclared matrix (Core Feature 3) and the
obligations no declaration waives (Core Feature 2). Behavior probes for
high-risk journeys stay in the default path only. Add the rule that every row
has its own observable, with no aggregate row and no row that re-checks a
Mechanical Refusal Code (Core Feature 4).

**Row input declaration.** Inputs are written when the row is planned as
`pending`, bounded to what the row reads. A row whose inputs grew after it ran is
not carriable (Core Feature 6). The carry-forward mechanics, kinds and
carried-row notation stay.

**Section 3, static gate.**

- A timeout or intermittent failure of the repository Verification is
  code-caused unless the unchanged delivery target reproduces it under the same
  command (Core Feature 7).
- An analyzer the `qa` Task names beyond the repository Verification runs over
  the changed packages. A diagnostic identical on the delivery target is recorded
  as observed with its owning Spec named. A repository Verification failure stays
  blocking (Core Feature 8).

**Sections 5 and 6, verdict.** No rule changes. The verdict clause about an
authored QA Task that defines a partial run keeps its words; Section 1 now says
what that clause does not cover.

### Anchor clauses

Task Verification and the gate replay look for these exact anchors, so the
skill must carry them verbatim in both copies:

```text
Each numbered `qa` Task Requirement that starts with `MUST verify` or `MUST run` is exactly one row
declaring a matrix does not make the run partial
Once the matrix exists, a finding blocks only the rows that depend on it
fixed when the row is planned as `pending`
is not a proved environmental cause
identical on the delivery target
```

These sentences leave the skill:

```text
Any problem blocks flow QA after the complete command audit has been reported.
The scope is complete when every promise and explicit exclusion in the spec maps to a planned row
```

### QA contract shape

The contract is a Go raw string, so it carries no backticks:

```text
QA contract:
- Run the qa-gate process for this Spec. Each numbered qa Task Requirement that
  starts with MUST verify or MUST run is exactly one matrix row, and such
  Requirements are the complete matrix; otherwise derive rows from the PRD user
  stories and Goals, each non-QA Task's Acceptance Criteria and the declared
  intentional breaks. Always keep the outside-evidence row, the Pull Request
  row, the repository Verification and, when the PRD declares frontend, the
  frontend sweep.
- Once the matrix exists, a finding blocks only the rows that depend on it;
  every other row still runs.
- <report naming bullet, unchanged>
- <verdict bullet, unchanged>
- Never commit, push, or open a pull request.
```

### Data Models

None. QA Report frontmatter, row statuses, blocked-cause counts and evidence
layout are unchanged.

### API Contracts

None. No command, flag, exit code or output changes.

## Coverage Map

- Goal 1 → qa-gate section 2 (declaration rule, bounded default); QA contract.
- Goal 2 → qa-gate section 1 (blocking sentence removed); QA contract.
- Goal 3 → qa-gate section 2 (no row re-checks a Mechanical Refusal Code).
- Goal 4 → gate replay of Spec 0119, 0134 and 0136 reports (Testing Approach 4).
- User Story 1 → qa-gate sections 1 and 2; QA contract.
- User Story 2 → qa-gate section 1; QA contract.
- User Story 3 → qa-gate Row input declaration.
- User Story 4 → qa-gate section 3 (analyzer scope).
- Core Features 1-4 → qa-gate sections 1 and 2.
- Core Feature 5 → qa-gate section 1.
- Core Feature 6 → qa-gate Row input declaration.
- Core Features 7-8 → qa-gate section 3.
- Core Feature 9 → QA contract in the gate prompt.

## Integration Points

- **Skill distribution.** The canonical skill regenerates its mirror through the
  repository's skill-sync target. The mirror parity test and the owned-skill
  version check hold the two together. The skill keeps its version, so no setup
  minimum moves.
- **Workflow contract test.** It pins the tooling-audit clauses. The change
  keeps all of them, so the test needs no edit.
- **Mechanical stage.** It withholds the Agent Session when it finds a blocking
  fact before a matrix exists. Nothing here changes that; scoped blocking
  applies only to the rows of a matrix the gate has planned.
- **Spec 0122.** Its proposed authority also names the qa-gate skill. It builds
  on this text once this Spec lands.

## Testing Approach

1. **Prompt seam.** The existing QA prompt tests assert on `BuildQAPrompt`
   output. A new test asserts that the contract carries:
   - the declaration rule;
   - the bounded default;
   - the non-waivable obligations, including the frontend sweep;
   - scoped blocking.

   It also asserts that the contract no longer tells the gate to validate every
   user story and acceptance criterion. The test fails on the tree as it stands
   today. The existing tests reference the contract through its constant and
   keep passing.
2. **Skill text.** Task Verification checks, in both copies, that each anchor
   clause is present and each removed sentence is absent. It also runs the mirror
   parity test and the workflow contract tests. The anchor checks fail on the
   tree as it stands today.
3. **Repository gate.** The terminal QA Task records `make verify` clean. It then
   runs `go vet` over the prompt's package and records the pre-existing lock-copy
   diagnostics identical on the delivery target, attributed to Spec 0123. That
   row exercises Core Feature 8 on real diagnostics.
4. **Outside evidence.** The terminal QA Task replays reports this Spec did not
   write:
   - the six Spec 0119 reports;
   - Spec 0136's reports;
   - Spec 0134's analyzer evidence;
   - the `qa` Tasks of Specs 0134 to 0137, under the declaration rule.

   The replay applies the rules as written and does not re-run those Specs'
   gates. The terminal QA Task declares its own matrix under the same rule.

## Build Order

1. Rewrite the qa-gate skill's declaration, default, blocking, input and
   static-gate rules in the canonical copy, then regenerate the mirror (depends
   on: none).
2. Replace the QA contract's matrix bullet and add its prompt test (depends on:
   none).
3. Terminal QA (depends on: 1, 2).

Steps 1 and 2 touch disjoint files and may run in the same wave.

## Risks & Considerations

- **Fleet blast radius.** The skill ships to every repository that installs it.
  - A fleet `qa` Task with `MUST verify` Requirements now gets only those rows.
  - A generic `qa` Task with no such Requirement keeps the bounded default.
  - The frontend sweep and the Pull Request and outside-evidence rows stay
    mandatory.
- **Pinned clauses.** Editing the tooling-audit bullet can drop a clause the
  workflow contract test pins. Task 1 deletes one sentence there and its
  Verification runs that test.
- **Prompt and skill drift.** The anchor clauses give both texts one wording, and
  Task 2's test pins the prompt side.
- **Known gap left open.** The mechanical stage takes Task commits only from the
  current Run's start head onward, and its report does not name them. The gate
  therefore keeps auditing every Task commit. A pending finding carries this.

## Decisions

- **Declaration rule.** A `qa` Task declares its matrix by Requirements that
  start with `MUST verify` or `MUST run`, the form Specs 0134 to 0137 used. A
  semantic test of whether Requirements "name what the gate verifies" was
  rejected because two Agents can judge it differently. See ADR-0155.
- **Tooling audit.** The tooling-audit bullet keeps its behavior. Removing the
  repeated audit waits on the stage naming the commits it audited.
- **No Daemon change.** Neither the mechanical stage's commit range nor its
  withholding rule changes.
- **Regeneration.** The mirror changes only through skill sync. Derived pins
  change only through the sanctioned digest regeneration, if at all.
- **ADR status.** The declared-matrix decision is recorded as ADR-0155, accepted
  with the maintainer's approval of this Spec's grant.

## Vocabulary Contract

No token is coined.
