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
story and acceptance criterion. Both texts change to say the same things:

- **The coverage is declared.** A `qa` Task Requirement that starts with
  `MUST verify` or `MUST run` is a coverage source. A Requirement that covers
  every Acceptance Criterion of named Tasks contributes one source per
  criterion. A Task with no such Requirement gets a bounded default.
- **Some sources are always covered:**
  - the outside-evidence row;
  - the Pull Request row;
  - the repository Verification;
  - each PRD Unreachable Acceptance declaration;
  - the frontend sweep.
- **Coverage is named, complete and closed.**
  - Every row names in its provenance the sources it covers.
  - Every source appears in some row.
  - No row covers anything else.
  - How rows group sources is left open.
- **Blocking is scoped.** Once the matrix exists, a finding blocks only the rows
  that depend on it.

Verdict rules, report keys, the mechanical stage, and the stage's withholding of
the Agent Session are unchanged.

This design accepts three trade-offs.

1. **Only declared sources are covered.** A promise left out of the `qa` Task
   goes unchecked. The maintainer accepted that risk and assigned the
   strengthening of authoring to Spec 0129.
2. **Grouping is not fixed.** Two executors cover the same sources but may split
   them into different rows. A rule for row shape was tried in three review
   rounds, and each round found a new place where it needed interpretation.
3. **The changed-path audit is still repeated.** The mechanical stage receives
   only the current Run's Task commits and does not name them in its report, so
   the gate cannot know what to skip. That gap goes to follow-up.

## Project Constraints

- Identifier strategy: applicable — the QA Report frontmatter keys, the three typed blocked causes, the Mechanical Refusal Codes and matrix row identifiers keep their spelling and meaning; this Spec introduces no identifier, refusal code or report key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request or HTTP surface is created or read; the change is guidance for the gate executor and the text of its prompt. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's matrix, verdicts and evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and every row still ends in one of the typed blocked causes or a result.
  ADR-0088 applies: the QA gate is authored into the graph, which is where the declared coverage is read.
  ADR-0091 applies: the QA gate is a Task node of its own type, so the coverage is declared by the `qa` Task.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, so the stage still withholds the Agent Session on a blocking fact, no row re-checks a fact the stage decides, and no stage makes a verdict more permissive.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, so row inputs are fixed when the row is planned.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row stays a source every matrix covers.
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

The gate Agent reads the prompt first and the skill second, so the two texts must
agree. The prompt stays short and points at the skill for detail, as it does
now. The Daemon builds the prompt only when the mechanical stage did not withhold
the Agent. That is why scoped blocking applies only to a planned matrix and never
overrides the stage.

## Implementation Design

### Where each Core Feature lands in the skill

**Section 1, scope.**

- Replace the closing sentence that makes every promise and explicit exclusion a
  planned row. The scope is complete when coverage is complete and closed.
- State that declaring a matrix does not make the run partial (Core Feature 1).
- In the tooling-audit bullet, delete only the sentence that stops flow QA after
  any audit problem, and point instead to finding-dependent blocking (Core
  Feature 5). The bullet keeps its behavior: the gate still audits Task commits
  by command, and every clause the contract test pins stays verbatim.
- The existing Unreachable Acceptance paragraph stays.

**Section 2, rows.** Replace the "Add a row for" list with:

- the declaration rule and its criterion expansion (Core Feature 1);
- the non-waivable sources (Core Feature 2);
- the bounded default (Core Feature 3);
- the coverage rule (Core Feature 4): every row names in its provenance the
  sources it covers, every source appears in at least one row, no row names a
  source outside the set, rows may group sources one observation settles, no
  row's result is computed from other rows, and no row re-checks a Mechanical
  Refusal Code.

Behavior probes for high-risk journeys stay in the default path only. The
existing Results table gains no column; provenance is written where each row
already describes what it covers.

**Row input declaration.** Inputs are written when the row is planned as
`pending`, bounded to what the row reads. A row whose inputs grew after it ran is
not carriable (Core Feature 6). The carry-forward mechanics, kinds and
carried-row notation stay.

**Section 3, static gate.**

- A timeout or intermittent failure of the repository Verification is recorded
  as a failure, never an environment block. Contention, a non-reproducible run
  and the same failure on the unchanged delivery target are not proved
  environmental causes (Core Feature 7). This tightens, never loosens, the
  existing environment-caused rule for that one command.
- An analyzer that the `qa` Task names beyond the repository Verification runs
  over the changed packages. A diagnostic identical on the delivery target is
  recorded as observed, naming the Spec that owns it. A repository Verification
  failure stays blocking (Core Feature 8).

**Sections 5 and 6, verdict.** No rule changes. The verdict clause about an
authored QA Task that defines a partial run keeps its words; Section 1 now says
what that phrase does not cover.

### Anchor clauses

Task Verification and the gate replay look for these exact anchors, so the skill
must carry them verbatim in both copies:

```text
Each numbered `qa` Task Requirement that starts with `MUST verify` or `MUST run` is a declared verification Requirement
declaring a matrix does not make the run partial
Every coverage source appears in at least one row's provenance
no row's provenance names a source outside the coverage sources
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

The contract is a Go raw string literal, so it carries no backticks:

```text
QA contract:
- Run the qa-gate process for this Spec. Each numbered qa Task Requirement that
  starts with MUST verify or MUST run is a declared verification Requirement, and
  the declared Requirements with the non-waivable sources are the coverage
  sources; otherwise the coverage sources are the PRD user stories and Goals,
  each non-QA Task's Acceptance Criteria and the declared intentional breaks.
- The non-waivable sources are the outside-evidence row, the Pull Request row,
  the repository Verification, each PRD Unreachable Acceptance declaration and,
  when the PRD declares frontend, the frontend sweep.
- Every row names in its provenance the sources it covers; every source appears
  in at least one row, and no row covers anything else.
- Once the matrix exists, a finding blocks only the rows that depend on it;
  every other row still runs.
- <report naming bullet, unchanged>
- <verdict bullet, unchanged>
- Never commit, push, or open a pull request.
```

### Data Models

None. QA Report frontmatter, the Results table columns, row statuses,
blocked-cause counts and evidence layout are unchanged.

### API Contracts

None. No command, flag, exit code or output changes.

## Coverage Map

- Goal 1 → qa-gate section 2 (declaration rule, bounded default, coverage rule);
  QA contract.
- Goal 2 → qa-gate section 1 (blocking sentence removed); QA contract.
- Goal 3 → qa-gate section 2 (no row re-checks a Mechanical Refusal Code).
- Goal 4 → gate replay of Spec 0119, 0134, 0135, 0136 and 0137 reports (Testing
  Approach 4).
- User Stories 1-2 → qa-gate sections 1 and 2; QA contract.
- User Story 3 → qa-gate section 1; QA contract.
- User Story 4 → qa-gate Row input declaration.
- User Story 5 → qa-gate section 3 (analyzer scope).
- Core Features 1-4 → qa-gate sections 1 and 2.
- Core Feature 5 → qa-gate section 1.
- Core Feature 6 → qa-gate Row input declaration.
- Core Features 7-8 → qa-gate section 3.
- Core Feature 9 → QA contract in the gate prompt.

## Integration Points

- **Skill distribution.** The canonical skill regenerates its mirror through the
  repository's skill-sync target. The mirror parity test and the owned-skill
  version check keep the two copies identical. The skill keeps its version, so
  no setup minimum moves.
- **Workflow contract test.** It pins the tooling-audit clauses. The change keeps
  all of them, so the test needs no edit.
- **Mechanical stage.** It withholds the Agent Session on a blocking fact before a
  matrix exists. Nothing here changes that.
- **Spec 0122.** Its proposed authority also names the qa-gate skill. It builds
  on this text once this Spec lands.

## Testing Approach

1. **Prompt seam.** The existing QA prompt tests assert on `BuildQAPrompt`
   output. A new test asserts that the contract carries:
   - the declaration rule;
   - the bounded default;
   - the non-waivable sources, including Unreachable Acceptance declarations and
     the frontend sweep;
   - the coverage rule;
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
   diagnostics that are identical on the delivery target, attributing them to
   Spec 0123. That row exercises Core Feature 8 on real diagnostics.
4. **Outside evidence.** The terminal QA Task replays reports this Spec did not
   write:
   - the six Spec 0119 reports;
   - Spec 0136's reports;
   - Spec 0134's analyzer evidence;
   - the QA Reports of Specs 0134, 0135 and 0137, against their `qa` Tasks' declared
     sources.

   The replay applies the rules as written and does not re-run those Specs'
   gates. The terminal QA Task declares its own coverage under the same rule.

## Build Order

1. Rewrite the qa-gate skill's declaration, coverage, blocking, input and
   static-gate rules in the canonical copy, then regenerate the mirror (depends
   on: none).
2. Replace the QA contract's matrix bullet and add its prompt test (depends on:
   none).
3. Terminal QA (depends on: 1, 2).

Steps 1 and 2 touch disjoint files and may run in the same wave.

## Risks & Considerations

- **Fleet blast radius.** The skill ships to every repository that installs it.
  - A fleet `qa` Task with `MUST verify` Requirements now covers only those
    sources and the non-waivable ones.
  - A generic `qa` Task with no such Requirement keeps the bounded default.
  - The frontend sweep and the Pull Request, outside-evidence and Unreachable
    Acceptance sources stay mandatory.
- **Provenance quality.** Coverage is only checkable if provenance names sources
  precisely. The skill asks for the `qa` Task's own identifiers ("Requirement
  3", "Task 01 criterion 2", "repository Verification"), the form Specs 0134,
  0135 and 0137 already used. No mechanical detector enforces this yet; that
  remains a Non-Goal.
- **Pinned clauses.** Editing the tooling-audit bullet can drop a clause the
  workflow contract test pins. Task 1 deletes one sentence there, and its
  Verification runs that test.
- **Prompt and skill drift.** The anchor clauses give both texts one wording,
  and Task 2's test pins the prompt side.
- **Known gap left open.** The mechanical stage takes Task commits only from the
  current Run's start head onward and does not name them in its report. The gate
  therefore keeps auditing every Task commit. A pending finding carries this.

## Decisions

- **Declaration rule.** A `qa` Task declares its coverage with Requirements that
  start with `MUST verify` or `MUST run`, the form Specs 0134 to 0137 used. See
  ADR-0155.
- **Coverage, not row shape.** The maintainer chose this on 2026-09-15, after
  three review rounds each found a new interpretation gap in one row per
  Requirement.
- **Tooling audit.** The tooling-audit bullet keeps its behavior. Removing the
  repeated audit waits until the stage names the commits it audited.
- **No Daemon change.** Neither the mechanical stage's commit range nor its
  withholding rule changes.
- **Regeneration.** The mirror changes only through skill sync. Derived pins
  change only through the sanctioned digest regeneration, if at all.
- **ADR status.** ADR-0155 records the declared-coverage decision. It was
  accepted with the maintainer's approval of this Spec's grant.

## Vocabulary Contract

No token is coined.
