---
spec: 0138-a-matrix-the-spec-declares
prd: _prd.md
created: 2026-09-14
---

# A matrix the Spec declares — Technical Spec

## Executive Summary

The change touches two texts and no running code path. The first is the qa-gate
skill, which the gate's Agent follows. The second is the QA contract the Daemon
places in that Agent's prompt, which today tells the gate to validate every user
story and acceptance criterion. Both change to say:

- the `qa` Task's Requirements declare the matrix;
- an undeclared matrix derives rows from a bounded default;
- a finding blocks only its dependents;
- the gate does not repeat an audit the mechanical stage already made.

Verdict rules, report keys and the mechanical stage stay as they are.

The main trade-off is that the matrix now covers only what the Spec wrote down.
A `qa` Task that forgets a promise will not have it checked. The maintainer
accepted that risk and sent the strengthening of authoring to Spec 0129. The
second trade-off is that this Spec leaves a known gap: the mechanical stage sees
only the current Run's Task commits. Closing it would change which revision
authorizes a commit, which is an authority decision this Spec does not make, so
the gate keeps auditing the commits the stage never received.

## Project Constraints

- Identifier strategy: applicable — the QA Report frontmatter keys, the three typed blocked causes, the Mechanical Refusal Codes and matrix row identifiers keep their spelling and meaning; this Spec introduces no identifier, refusal code or report key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request or HTTP surface is created or read; the change is guidance for the gate executor and the text of its prompt. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's matrix, verdicts and evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and every row still ends in one of the typed blocked causes or a result.
  ADR-0088 applies: the QA gate is authored into the graph, which is where the declared matrix is read.
  ADR-0091 applies: the QA gate is a Task node of its own type, so the matrix is declared by the `qa` Task.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, so the gate stops re-deriving what the mechanical stage computes, and no stage makes a verdict more permissive.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, so row inputs are fixed when the row is planned.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row stays in every matrix.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the commit-dependent audit stays a gate fact executed as a command.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, which this Spec writes into the qa-gate skill and the QA contract.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule or citation reader.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-14, recorded in `docs/specs/0138-a-matrix-the-spec-declares/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved skill edit moves. The QA contract in the gate prompt is ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| qa-gate skill, canonical copy | Tells the gate Agent how to plan, run and close the matrix | Sections 1, 2, Row input declaration and 3 take the rules of Core Features 1-9 |
| qa-gate skill, distributed mirror | Embedded in the binary and installed into repositories | Regenerated from the canonical copy, never hand-edited |
| QA contract in the gate prompt | The rule summary the Daemon hands the gate Agent | The first bullet takes the declared matrix and bounded default (Core Feature 10) |
| Mechanical stage | Computes machine facts before the Agent turn | None |
| Workflow contract test for the skill | Pins the tooling-audit clauses the gate must keep | None; the clauses it pins stay in the skill |

The gate Agent reads the prompt first and the skill second, so the two texts must
say the same thing. The prompt stays short and points at the skill for detail,
as it does now.

## Implementation Design

### Where each Core Feature lands in the skill

- **Section 1, scope.**
  - Replace the closing sentence that makes every promise and explicit exclusion
    a planned row. The scope is complete when every declared source has a row.
  - State that declaring a matrix does not make the run partial (Core
    Feature 1).
  - Change the tooling-audit bullet so the gate audits by command only the Task
    commits the mechanical stage did not receive (Core Feature 6). The clauses
    that name actual paths from the Daemon-owned Task commit, the committed-path
    command, missing, late or untraceable authorization, and paths outside the
    exact bounded list stay word for word.
  - Delete the sentence that stops flow QA after any audit problem. Point
    instead to finding-dependent blocking (Core Feature 5).
- **Section 2, rows.** Replace the "Add a row for" list with two paths:
  - the declared matrix from the `qa` Task's Requirements (Core Feature 1);
  - the bounded default: PRD user stories and Goals, each non-QA Task's
    Acceptance Criteria, and declared intentional breaks (Core Feature 3).

  Add the obligations no declaration waives (Core Feature 2). Keep behavior
  probes for high-risk journeys in the default path only. Add the rule that
  every row has its own observable, with no aggregate row and no row that
  re-checks a Mechanical Refusal Code (Core Feature 4).
- **Row input declaration.** Inputs are written when the row is planned as
  `pending`, bounded to what the row reads. A row whose inputs grew after it ran
  is not carriable (Core Feature 7). The carry-forward mechanics, kinds and
  carried-row notation stay.
- **Section 3, static gate.**
  - A timeout or intermittent failure of the repository Verification is
    code-caused unless the unchanged delivery target reproduces it under the same
    command (Core Feature 8).
  - An analyzer the `qa` Task names beyond the repository Verification runs over
    the changed packages. A diagnostic identical on the delivery target is
    recorded as observed, with its owning Spec named. A repository Verification
    failure stays blocking (Core Feature 9).
- **Sections 5 and 6, verdict.** No rule changes. The verdict clause for an
  authored QA Task that defines a partial run keeps its words; Section 1 now says
  what that phrase does not cover.

### Anchor clauses

Task Verification and the gate replay look for these exact anchors, so the
skill must carry them verbatim in both copies:

```text
the `qa` Task's Requirements are the complete matrix
declaring a matrix does not make the run partial
A finding blocks only the rows that depend on it
Task commits the mechanical stage did not receive
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

```text
QA contract:
- Run the qa-gate process for this Spec. The qa Task's Requirements are the
  complete matrix when they name what the gate verifies; otherwise derive rows
  from the PRD user stories and Goals, each non-QA Task's Acceptance Criteria and
  the declared intentional breaks. Always keep the outside-evidence row, the Pull
  Request row and the repository Verification.
- A finding blocks only the rows that depend on it; every other row still runs.
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

- Goal 1 → qa-gate section 2 (declared matrix, bounded default); QA contract.
- Goal 2 → qa-gate section 1 (blocking sentence removed); QA contract.
- Goal 3 → qa-gate section 1 (tooling-audit bullet), section 2 (no re-check of
  a Mechanical Refusal Code).
- Goal 4 → gate replay of Spec 0119, 0134 and 0136 reports (Testing Approach 4).
- User Story 1 → qa-gate sections 1 and 2; QA contract.
- User Story 2 → qa-gate section 1 and section 3.
- User Story 3 → qa-gate section 1 (tooling-audit bullet).
- User Story 4 → qa-gate Row input declaration.
- User Story 5 → qa-gate section 3 (analyzer scope).
- Core Features 1-3 → qa-gate sections 1 and 2.
- Core Feature 4 → qa-gate section 2.
- Core Feature 5 → qa-gate section 1.
- Core Feature 6 → qa-gate section 1.
- Core Feature 7 → qa-gate Row input declaration.
- Core Features 8-9 → qa-gate section 3.
- Core Feature 10 → QA contract in the gate prompt.

## Integration Points

- **Skill distribution.** The canonical skill regenerates its mirror through the
  repository's skill-sync target. The mirror parity test and the owned-skill
  version check hold the two together. The skill keeps its version, so no setup
  minimum moves.
- **Workflow contract test.** It pins the tooling-audit clauses listed above.
  The change keeps them, so the test needs no edit.
- **Spec 0122.** Its proposed authority also names the qa-gate skill. It builds on
  this text once this Spec lands.

## Testing Approach

1. **Prompt seam.** The existing QA prompt tests assert on `BuildQAPrompt`
   output. A new test asserts:
   - the contract names the declared matrix, the bounded default and
     finding-dependent blocking;
   - the contract no longer tells the gate to validate every user story and
     acceptance criterion.

   It fails on the tree as it stands today. The existing tests reference the
   contract through its constant and keep passing.
2. **Skill text.** Task Verification checks that each anchor clause is present
   and each removed sentence is absent, in both copies. It also runs the mirror
   parity test and the workflow contract test. The anchor checks fail on the tree
   as it stands today.
3. **Repository gate.** The terminal QA Task records `make verify` clean. It runs
   `go vet` over the prompt's package and records the pre-existing lock-copy
   diagnostics identical on the delivery target, attributed to Spec 0123. That
   row exercises Core Feature 9 on real diagnostics.
4. **Outside evidence.** The terminal QA Task replays reports this Spec did not
   write:
   - the six Spec 0119 reports, to show each real finding still reaches a planned
     row or a mechanical finding and each noise row does not recur;
   - Spec 0136's report, to show its governed-rename finding came from a `qa`
     Task Requirement;
   - Spec 0134's analyzer evidence, to show diagnostics identical on the delivery
     target.

   The replay applies the rules as written. It does not re-run those Specs'
   gates.

## Build Order

1. Rewrite the qa-gate skill's matrix, blocking, audit, input and static-gate
   rules in the canonical copy and regenerate the mirror (depends on: none).
2. Replace the QA contract's matrix bullet and add its prompt test (depends on:
   none).
3. Terminal QA (depends on: 1, 2).

Steps 1 and 2 touch disjoint files and may run in the same wave.

## Risks & Considerations

- **Fleet blast radius.** The skill ships to every repository that installs it.
  A fleet `qa` Task whose Requirements name only some targets now gets only
  those. The Requirement that only tells the gate to run declares nothing, which
  keeps generic `qa` Tasks on the bounded default. The frontend sweep and the
  Pull Request and outside-evidence rows stay mandatory.
- **Pinned clauses.** Rewording the tooling-audit bullet can drop a clause the
  workflow contract test pins. Task 1's Verification runs that test.
- **Prompt and skill drift.** The anchor clauses give both texts one wording, and
  Task 2's test pins the prompt side.
- **Known gap left open.** The mechanical stage takes Task commits only from the
  current Run's start head onward and uses that head as its delivery target, so
  a multi-Run Spec leaves earlier commits to the gate. Recording this as a
  follow-up keeps it visible without changing authority here.

## Decisions

- The declared matrix is read from the `qa` Task's Requirements, not from a new
  section, so authoring skills stay unchanged until Spec 0129.
- The tooling-audit change only removes repetition. Which revision authorizes a
  Task commit from an earlier Run is not decided here.
- No Daemon code changes; the mechanical stage's commit range is follow-up work.
- The mirror changes only through skill sync. Derived pins change only through
  the sanctioned digest regeneration, if at all.
- The declared-matrix decision is recorded as ADR-0155, accepted with the
  maintainer's approval of this Spec's grant.

## Vocabulary Contract

No token is coined.
