---
spec: 0192-owned-skills-that-describe-the-product-as-it-is
prd: _prd.md
created: 2026-09-30
---

# Owned skills that describe the product as it is — Technical Spec

## Executive Summary

Twelve owned skills and five guide files get text corrections, and two new
repository tests guard the class that drifted most: a CLI command nobody
documented, and a guide link that resolves to nothing. No Go production code
changes. The tests live in `internal/docscontract`, which `make verify-docs`
already runs at the Pull Request boundary, so no Makefile or CI change is
needed. The main trade-off is a literal-name contract: the tests prove that
`roundfix <command path>` appears in the skill and the guide, not that the
description is right. That is cheap and cannot be satisfied by a paraphrase,
and it leaves accuracy to review and the QA gate.

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Markdown files and repository
  tests only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0057, ADR-0116, ADR-0156, ADR-0167
  and ADR-0179 hold and the skills are aligned with them; the gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104 and ADR-0155; ADR-0093,
  ADR-0117, ADR-0166 and ADR-0176 hold; ADR-0097, ADR-0160, ADR-0168 and ADR-0170 do not
  apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for every owned skill, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`,
  `.agents/skills/implement-spec/SKILL.md`, `skills/implement-spec/SKILL.md`,
  `.agents/skills/council/SKILL.md`, `skills/council/SKILL.md`,
  `.agents/skills/council/references/archetypes.md`,
  `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md`,
  `.agents/skills/write-idea/references/opportunity-scan.md`,
  `.agents/skills/business-analyst/SKILL.md`,
  `skills/business-analyst/SKILL.md`, `.agents/skills/brainstorming/SKILL.md`,
  `skills/brainstorming/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

Four groups of files change, and they share nothing:

| Group | Files | Task |
| --- | --- | --- |
| Roundfix skill and command contract | the Roundfix skill, its mirror, a new test file | task_01 |
| User guide | five guide files, the README, a new test file | task_02 |
| Authoring skills | `write-prd`, `write-techspec`, `write-tasks`, their templates and mirrors | task_03 |
| Gate, lifecycle and discovery skills | `qa-gate`, `archive-spec`, `setup-context-driven`, `implement-spec`, `council`, `write-idea`, `business-analyst`, `brainstorming`, two reference files and mirrors | task_04 |

Every skill edit is made in `.agents/skills/`, the canonical tree. The
sanctioned `make skills-sync` copies each owned skill folder to `skills/`, and
`make baseline-digests` follows. A scratch clone measured that these edits
leave every derived artifact unchanged, so parallel Tasks share no generated
file.

## Implementation Design

### Interfaces

Two test files join `internal/docscontract`, under its existing `docscontract`
build tag. They reach the CLI through `cli.Run`, the seam
`publicdocs_test.go` already uses.

```go
// internal/docscontract/command_documentation_test.go (task_01)

// commandPaths returns each command path a root help usage line names, in
// first-seen order and without duplicates.
func commandPaths(help string) []string

// undocumentedCommands returns the paths whose literal "roundfix <path>" is
// absent from text.
func undocumentedCommands(paths []string, text string) []string
```

```go
// internal/docscontract/user_guide_contract_test.go (task_02)

// brokenLinks returns "file:line: destination" for each relative Markdown
// link in content that does not resolve from the file's directory.
func brokenLinks(repoRoot, file, content string) []string
```

### The command path rule

`commandPaths` reads only lines that start with two spaces and `roundfix `.
It takes the tokens after `roundfix` until the first token that starts with
`[`, `--`, `(` or `<`.

- No token before the stop means a flag-only line such as `roundfix --help`;
  it yields nothing.
- When the stopping token has the form `<a|b|c>`, each alternative is a
  subcommand: `roundfix window <set|show|clear>` yields `window set`,
  `window show` and `window clear`.
- Any other stopping token ends the path: `roundfix archive <slug>` yields
  `archive`.

On v0.21.0 this yields fifty-one paths. The contract test refuses fewer than
forty and requires `deliver retry`, `window clear` and
`baseline capabilities check` to be among them, so a changed help format
cannot pass with an empty list.

### The Roundfix skill

Edits to `.agents/skills/roundfix/SKILL.md`, all outside `### QA settlement`:

- **Commands.** Sections for `roundfix window set`, `roundfix window show`,
  `roundfix window clear`, `roundfix qa-report accept`,
  `roundfix baseline capabilities check`, `roundfix baseline profile init` and
  `roundfix init`, each written from the command's own `--help` and its
  `docs/user-guide/commands.md` section. `qa-report accept` takes no `--help`;
  its usage is printed by `roundfix qa-report --help`.
- **Event categories.** The Supervisor Run Event Stream paragraph lists
  `task-status`, `batch`, `verification`, `outcome`, and `agent-selection`,
  as `roundfix events --help` does.
- **Delivery order.** The sentence that starts "Follow one order per Spec:"
  states the order the Delivery queue section already gives: implement the
  graph including its authored gate, run the configured pre-PR review, archive
  on the branch, pass the repository gate, open the Pull Request, verify
  current-head checks, and merge. The numbered steps of "Autonomous Spec
  delivery" follow it: a step titled **Review the candidate.** comes before
  **Archive and publish.**, and `roundfix watch … --until-clean` is described
  as the step for a repository whose Review Source gives Pull Request
  feedback, not as a step every Spec takes.
- **Reopen.** The step about corrective Tasks says that a Task added as a
  dependency of a QA gate that already settled `completed` needs
  `roundfix reopen --spec <slug>` first, and that a failed or pending gate
  needs no reopen.
- **Version.** Both version fields stay at `0.0.2`. Raising the Roundfix skill's version breaks three existing `./skills` tests that read the owned-skill floor from its version line (measured on 2026-09-30); the Spec that makes the floor follow the bundle raises it.

The adapter version in `agents/openai.yaml` is left to Spec 0189.

### The user guide

- **`--qa`.** `commands.md`, `usage.md` and `README.md` drop the flag from the
  `implement` synopsis and examples. They say the gate is the Spec's authored
  terminal `qa` Task and that `--qa` is refused as an unknown flag.
- **Archive destination.** `commands.md` and
  `context-driven-development.md` state `docs/history/specs/<slug>/` for the
  built-in Spec Root and `<specs.root>/_archived/<slug>/` otherwise.
- **Links.** Fifteen links point at `../specs/_archived/…`,
  `../findings/_archived/…`, `../specs/0036-…` and two retired ADRs. Each is
  retargeted to the document's current path under `docs/history/`. A link
  whose document no longer exists anywhere is removed and its sentence kept.
- **Skills line.** The three places that print
  `skills: ok (39 required: 14 Roundfix-owned, 25 external)` show the shape
  `skills: ok (<required> required: <owned> Roundfix-owned, <external> external)`
  and say the numbers come from the repository's Repository Skill Set.
- **Missing commands.** `commands.md` gains entries for `roundfix spec audit`
  and `roundfix storage report`, written from each command's `--help`.

### The authoring skills

- **Authorization record.** `write-prd`, `write-techspec` and `write-tasks`
  state that every Spec a Delivery Queue delivers carries an approved
  `_authorization.md`, with or without protected tooling. Its `operations`
  list grants `implement`, `commit`, `push`, `pull_request` and `merge`, and a
  Spec that changes no Governed Path records `paths: []`. The two templates
  say the same in their Tooling authority comment. They keep exactly four
  `<applicable | not applicable>` placeholders, four
  ``Source: `docs/agents/`` placeholders, and the phrases
  `express maintainer authorization`, `bounded files` and
  `no protected tooling mutation`, which a contract test counts.
- **Backend guide.** `write-prd` and `write-techspec` keep the literal
  `docs/agents/backend.md` and add that a repository without that guide cites
  the guide that owns the policy for its surfaces.
- **Unreachable Acceptance.** The PRD template gains an optional
  `## Unreachable Acceptance` section with the three fields the parser reads:
  `criterion`, `reason` and `satisfied-by`.
- **Coverage Map.** The TechSpec template's comment asks for one line per PRD
  goal, user story, Core Feature and Success Metric.
- **`write-tasks`.** It states the delivery order above, names
  `SC-CITATION-UNSUPPORTED` instead of denying the check, says
  `roundfix reopen --spec <slug>` is needed when a Task is added after a
  completed gate, and says code for another operating system is verified by
  building its non-test code for that system, never with `go vet`, which also
  compiles tests written for the host. Its task template says the Daemon
  writes `status` during a Run and `implement-task` writes it only when run
  standalone.
- **Versions.** `write-prd` and `write-techspec` move to `0.0.3`;
  `write-tasks` moves from `0.0.3` to `0.0.4`.

### The structured question form

The clause in `docs/agents/agent-instructions.md` asks for exactly one question
per call, two or three concise and mutually exclusive options with the
consequence of each, the recommended option first with `(Recommended)` in its
label, and no `Other` option where the interface already supplies a custom
answer. Five skills show another form. Each is aligned where it stands:

- **`write-prd`** (task_03). The example block under "Clarify" keeps its
  question and three options. The first option is labelled `(Recommended)`
  and each option states its consequence. The `D) Other — describe` line goes.
  The paragraph below says to ask through the structured question tool the
  session exposes, and that without one the same single question is shown in
  chat and a custom answer is accepted.
- **`write-idea`** (task_04). Its ground rule names the same form instead of
  `D) Other — describe` as the escape.
- **`brainstorming`** (task_04). Its "Multiple choice preferred" rule asks for
  two or three options with the recommended one first, and drops
  `D) Other — describe`.
- **`business-analyst`** (task_04). The Mode 2 block drops option D, labels
  option A `(Recommended)`, and its rule says two or three options per
  decision.
- **`council`** (task_04). Step 7 replaces the compound question "Which path
  are you taking, and what triggers would cause you to revisit this
  decision?" with two questions asked one at a time. The first asks which
  path, with the two or three paths the synthesis named as options and the
  recommended one first. After the answer, the second asks which trigger
  reopens the decision, with two or three triggers drawn from the synthesis
  risks. Both answers are recorded verbatim under `## Decision Captured`.

`write-techspec` refers to the PRD stage's protocol and needs no edit.

### The gate, lifecycle and discovery skills

- **`qa-gate`.**
  - Its static-gate step says a Daemon-assigned gate records the Daemon's
    repository Verification result and does not run it again, and that a
    standalone gate runs the repository's selected Verification.
  - Its closing template keeps the seeded `| # | Status | Provenance |` header
    in `## Results`.
  - Its closing rule permits Pull Request preparation on `pass` or a
    qualifying declared `partial`.
  - Version `0.0.4`.
- **`archive-spec`.** A normal archive also runs through
  `roundfix archive <slug>`, which verifies, stamps and moves. The manual
  stamp and `git mv` steps go, the skill says never to hand-edit archive front
  matter, and the commit subject is `docs: archive <slug>`, the one the
  Delivery Queue writes. Version `0.0.3`.
- **`setup-context-driven`.** It keeps the phrase a contract test requires and
  adds the four terminal Finding statuses. Its plan review list names
  `historyMoves` and the `baseline.history.citation` warnings. Version
  `0.0.3`.
- **`implement-spec`.** The argument hint drops `[--from task_NN]`. Version
  `0.1.1`.
- **`council`.** `references/archetypes.md` says each advisor's prompt is
  built from its definition in that file, as the skill body already does.
  Version `0.0.3`.
- **`write-idea`.** `references/opportunity-scan.md` names step 7. Version
  `0.0.3`.
- **`business-analyst` and `brainstorming`.** Only the question form above
  changes. Version `0.0.3` each.

### Data Models

None. No schema, record or payload changes.

### API Contracts

None. No command, flag, output or exit code changes; the Spec changes what is
written about them and adds repository tests.

## Coverage Map

- Goal 1 → The authoring skills; The gate, lifecycle and discovery skills; The
  structured question form.
- Goal 2 → The Roundfix skill; The user guide; The command path rule.
- Goal 3 → The user guide.
- Goal 4 → Interfaces; The command path rule.
- Core Feature 1 → The Roundfix skill.
- Core Feature 2 → The user guide.
- Core Feature 3 → The authoring skills.
- Core Feature 4 → The gate, lifecycle and discovery skills.
- Core Feature 5 → The structured question form.
- Core Feature 6 → Interfaces; The command path rule.
- Core Feature 7 → The Roundfix skill; The authoring skills; The gate,
  lifecycle and discovery skills.
- Success Metric 1 → Testing Approach 1, Testing Approach 2.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 3.

## Integration Points

- **`make verify-docs`.** Its `docs-test` step runs
  `go test -tags docscontract ./internal/docscontract`, so both new test files
  run at the Pull Request boundary with no Makefile change.
- **Skill contract tests.** `skills/baseline_skill_contract_test.go` and
  `skills/settlement_guidance_repocontract_test.go` require phrases and a
  byte-identical section in the skills this Spec edits. They are read, never
  edited.
- **Skill mirrors.** `TestAuthorialSkillSync` and `make skills-sync-check`
  require each mirror to equal its canonical folder.
- **Spec 0194.** It later splits the Roundfix skill and the commands guide into
  one file per command. The two contract tests then read the assembled text or
  the fragments, which is that Spec's change to make.

## Testing Approach

1. **Command contract for the skill.** New
   `internal/docscontract/command_documentation_test.go`:
   - `TestCommandPathsReadSubcommandsAndAlternatives` feeds a literal help
     text and expects the paths the rule defines, including the three
     `window` alternatives and nothing for a flag-only line.
   - `TestEveryCommandIsNamedInTheRoundfixSkill` reads the real root help
     through `cli.Run`, enforces the floor and the three anchor paths, and
     expects `undocumentedCommands` to return nothing for the skill.
   - `TestAnUndocumentedCommandIsReported` removes one command's literal from
     a copy of the text and expects exactly that path back.
2. **User guide contract.** New
   `internal/docscontract/user_guide_contract_test.go`:
   - `TestEveryCommandIsNamedInTheUserGuide` applies the same check to the
     concatenated `docs/user-guide/*.md`.
   - `TestUserGuideLinksResolve` expects `brokenLinks` to return nothing for
     each guide file and the README. It skips fenced code, URL schemes and
     fragment-only links, and strips fragments before resolving.
   - `TestABrokenUserGuideLinkIsReported` feeds a document with one dead link
     and one live link and expects only the dead one.
   - `TestUserGuideNamesNoRefusedQAFlag` expects no `--qa` that is not
     `--qa-override`.
3. **Phrase checks.** Each Task's Verification asserts, on whitespace-folded
   text, that each stale sentence is gone and its replacement present, and
   that both version fields carry the new version.
4. **Existing contracts.** Each skill Task's Verification runs the contract
   tests that read the skills it edits, the mirror comparison, and
   `make skills-sync-check`.
5. **Through the built binary.** The QA gate builds the candidate, runs each of
   the seven newly documented command forms with `--help` or its usage, and
   compares the skill's statements with the output. It also runs
   `roundfix implement --spec <slug> --qa` and expects the unknown-flag
   refusal.

## Build Order

1. The Roundfix skill and the command contract test, task_01 (depends on:
   none).
2. The user guide, its link and command contract tests, task_02 (depends on:
   1). It calls `commandPaths` and `undocumentedCommands` from step 1's file.
3. The authoring skills, task_03 (depends on: none).
4. The gate, lifecycle and discovery skills, task_04 (depends on: none).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

Steps 1, 3 and 4 share no file and run in one Wave.

## Risks & Considerations

- **A contract test requires a phrase an edit would remove.** Each skill Task
  runs the contract tests that read its skills and rewrites around a required
  phrase. The governed tests are outside the grant, so a Task that cannot keep
  a phrase stops and says which one.
- **The shared settlement section.** No Task may change `### QA settlement`.
  `TestSettlementGuidanceIsOneTable` is in the Verification of task_01 and
  task_04, the two Tasks that edit a skill carrying it.
- **Literal matching.** A command documented only under another spelling
  fails the contract. That is intended: the literal is what a reader searches
  for.
- **Skill and guide disagree for one Wave.** Until Spec 0193 lands, the
  generated autonomous-work guide states archive before review. The PRD
  records this.
- **Fleet skills.** The owned skills ship to adopting repositories. Every new
  sentence describes Roundfix commands and the Spec contract, which are the
  same for an adopter.

## Decisions

- The contract reads the root help and matches literal command names; see the
  PRD's Decisions.
- The tests live under the `docscontract` build tag, because their input is
  repository Markdown, which `make verify` deliberately leaves out.
- The closing QA report template adopts the seeded three-column Results
  header, because that is the header the Daemon writes and the eligibility
  policy reads.
