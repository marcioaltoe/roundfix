---
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
prd: _prd.md
created: 2026-10-22
---

# The Roundfix skill's follow-up, status and unattended modes, and up to three questions per call — Technical Spec

## Executive Summary

Two surfaces change. The Baseline core module rewords
`clause.core.ask-user-answerable-decisions` from one question per call to at
most three independent questions per call, keeping its identity and
`stop-and-ask` enforcement; the module record command, `make
baseline-digests` and this repository's Managed Refresh regenerate the derived
files, and a new test file in `internal/baseline` holds the text, the
rendering and the `retained` disposition. The Roundfix skill gains a `##
Modes` section, OpenCode slash metadata and a topic reference
`references/modes.md`, and the mode is the first word after the skill's
name, `follow-up`, `status` or `unattended`, with or without a leading `--`;
a new test file in `skills` holds the mode detection,
the absence of argument placeholders, the follow-up step order, the
unattended answer rules, and the link between the reference's question limit
and the clause. No Go production code changes.

The trade-off accepted: the skill modes are prose an Agent follows, so the
tests prove the text's structure and its agreement with the Baseline, not an
Agent's behavior in a live harness; that criterion is declared under
Unreachable Acceptance.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted; the clause and the skill keep their identities. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  recipient is added; tests read repository files and temporary repositories.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0320 for the
  clause, and Build Orders 2 and 3 follow ADR-0320 for the modes and the answer
  placement and ADR-0330 for message triage and capture. Both are named with their
  titles in `_prd.md` → Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("same as cycle AO"), carrying the cycle AJ grant
  "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/modes.md`,
  `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at v0.109.0. Line numbers move when
other Specs merge; identifiers do not.

1. `internal/baseline/assets/modules/core.json`, rule
   `rule.core.decision-authority` version 5, holds the only clause of the
   rule, `clause.core.ask-user-answerable-decisions`, enforcement
   `stop-and-ask`, whose guidance says "Ask exactly one question per call or
   message and keep only one question pending, even when the tool accepts an
   array of questions.", "present the same single question with numbered or
   lettered options in chat" and "Never bundle questions". The module is
   version 21.
2. The clause renders as the `stop-and-ask` bullet of
   `docs/agents/agent-instructions.md` in this repository and in the Standard
   TypeScript golden under
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/`.
   No user guide or Roundfix skill reference repeats it (`grep -rn -i
   "exactly one question\|one question per\|question per call"` outside
   `docs/history`: only the module, the two guides and the owned interview
   skills `brainstorming`, `write-idea`, `write-prd`, `write-techspec` and
   `business-analyst`, which ask one question per message for their own
   interviews).
3. The Source Baseline
   `baseline.standard-typescript-monorepo-0.0.1` records the clause in its
   manifest with `stop-and-ask`, and its accounting maps the legacy clause to
   it as `replaced`. `classifySourceClauseTransition`
   (`internal/baseline/plan.go`) records a source clause `retained` when the
   selected Baseline has a clause with the same identity and enforcement.
   `clause_characterization_test.go` pins the clause's enforcement only.
4. Prototype, in a scratch worktree of v0.109.0 that never touched this
   clone: rewording the clause and raising the rule to version 6, then
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   wrote module version 22 into `core.json` and
   `internal/baseline/module-versions.json`; `make baseline-digests` rewrote
   the Standard TypeScript golden `docs/agents/agent-instructions.md`, the
   Standard TypeScript profile, `catalog.diagnostics.golden.json`,
   `catalog.digest`, `catalog.normalized.json` and the four plan
   characterization goldens (`advisory-only-divergences`, `clean-adoption`,
   `idempotent-replan-after-verified-apply`,
   `same-baseline-changed-profile-and-catalog-digests`);
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   updated `docs/agents/agent-instructions.md` and
   `docs/agents/setup-context.json`, and a second run reported no file
   change. `go test ./internal/baseline`, the `internal/cli` Baseline, Setup
   and Doctor tests and the documentation contract suite passed.
5. `.agents/skills/roundfix/SKILL.md` (version 0.1.127, 15.9 KB of a
   20,000-byte budget held by `TestRoundfixSkillEntryFileStaysWithinItsBudget`)
   has no modes and no `$` character; its reference index has 18 rows and
   `TestRoundfixSkillIndexNamesEveryReference` requires exactly 18 files.
   `skillCoverageReferenceIndex` skips index rows whose commands column starts
   with `—`, so a topic reference needs no Skill Coverage Map entry. The
   frontmatter parser in `skills/skills.go` reads only `name` and `version`;
   other metadata keys are free. `TestEveryOwnedSkillVersionIsRecorded`
   refuses changed content under an unchanged version, so the version is
   raised by hand before the record command.
6. Spec 0306, authored in parallel and delivered first, owns the source both
   Specs answer and delivers the read-only follow-up command with its command
   guide, its Roundfix reference `references/follow-up.md` and its index row,
   which raise the reference count to 19; this Spec's Tasks do not touch them.
7. Amendment characterization, 2026-10-24, at v0.118.0 (`beb6b66b`), after
   Spec 0306 merged: `roundfix --help` lists `roundfix follow-up [--status]
   [--messages <path>] [--format <text|json>]` and no top-level `status` or
   `unattended` command (`roundfix status` and `roundfix unattended` both
   print `unknown command`); `deliver status` exists as a subcommand. The
   Roundfix skill is at version 0.1.135 with 19 references, its entry file is
   16,395 bytes and holds no `$` character, and no Roundfix skill file
   mentions a first word or a bare mode word. So the bare word `follow-up`
   after the skill's name is the one word that also names a CLI command, and
   the skill text has to say which one it selects.

## System Architecture

- `internal/baseline/assets/modules/core.json`: the reworded clause, rule
  version 6, module version by the record command.
- Derived Baseline files and this repository's guides: regenerated only by
  the sanctioned commands.
- `internal/baseline/question_batch_clause_test.go` (new).
- `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`
  (new), their mirrors under `skills/roundfix/`, the version record.
- `skills/roundfix_modes_test.go` (new) and `skills/roundfix_layout_test.go`.
- `CONTEXT.md`: **Skill Mode**.

No package, import edge or schema changes.

## Implementation Design

### Interfaces

```go
// internal/baseline/question_batch_clause_test.go (test-only)
const askUserClauseID = "clause.core.ask-user-answerable-decisions"
const askUserClauseLimit = "Ask at most three questions per call or message, never more than the tool accepts, and keep only that call pending."
const askUserClauseIndependence = "Put questions in one call only when they are independent of each other; a question whose options depend on another answer waits for that answer."
const askUserClauseRetired = "Ask exactly one question per call or message"

// skills/roundfix_modes_test.go (test-only)
func roundfixSection(text, heading string) string // body under a "## " heading up to the next "## " heading
```

### Data Models

None changes shape. The core module's rule `rule.core.decision-authority`
becomes version 6 and the module takes the next version the record command
writes. The Roundfix skill's frontmatter `metadata` gains the string key
`opencode/slash` with the value `"true"`.

### API Contracts

1. API Contract: the clause text. The guidance of
   `clause.core.ask-user-answerable-decisions` is exactly: "Never guess a
   decision the user can answer cheaply. Use the structured user-interaction
   tool actually exposed and permitted in the current session:
   `AskUserQuestion` in Claude Code, `request_user_input` in Codex, or
   `question` in OpenCode; Codex runtimes that expose
   `request_user_input_async` may use that equivalent. Interpret skill
   references to `AskUserQuestion` as this interaction contract across
   runtimes. Respect the tool's advertised schema and mode restrictions; do not
   invent tool calls or switch session modes merely to ask a question. Ask at
   most three questions per call or message, never more than the tool accepts,
   and keep only that call pending. Put questions in one call only when they
   are independent of each other; a question whose options depend on another
   answer waits for that answer. Wait for the user's response before asking the
   next call or doing work that depends on an answer; an asynchronous request
   may leave independent work running. Make each question answerable without
   rereading the session: state what was found and why it forces a choice,
   provide 2-3 concise, mutually exclusive options with the consequence of each
   choice, and put the recommended option first with `(Recommended)` in its
   label. Preserve custom text input; do not add an `Other` option when the
   interface already supplies one. If no permitted structured tool is
   available, present the same questions numbered in chat, each with numbered
   or lettered options, and accept a custom answer to each. Never merge two
   decisions into one question, replace the options with an open-prose
   question, or treat a recommendation, preselection, dismissal, timeout, or
   missing answer as the user's choice." Its id, enforcement and rule keep
   their values; it gains no `replaces`.
2. API Contract: retention. A Standard TypeScript adopter's Managed Refresh
   with the new catalog records the clause `retained` with the target
   `clause.core.ask-user-answerable-decisions`; the Source Baseline corpus,
   manifest, accounting and retention transitions keep their bytes.
3. API Contract: the skill entry. `.agents/skills/roundfix/SKILL.md` gains, in
   its frontmatter, `opencode/slash: "true"` under `metadata` and the
   description clause "and run its `--follow-up`, `--status` and
   `--unattended` Skill Modes"; in its body, a `## Modes` section right before
   `## Archive Relink and History Citations` whose first paragraph is exactly
   "A request names a Skill Mode with the first word after the skill's name:
   `follow-up`, `status` or `unattended`, with or without a leading `--`, in
   every harness: `/roundfix follow-up` or `/roundfix --follow-up` in Claude
   Code, `$roundfix status` or `$roundfix --status` in Codex, and `@roundfix
   unattended`, `@roundfix --unattended` or `/roundfix unattended` in
   OpenCode. A Skill Mode is not the CLI command `roundfix follow-up`, which
   the follow-up mode runs as one of its steps. Read the mode from the request
   text; a request whose first word after the skill's name is anything else,
   such as `/roundfix watch 123`, uses the skill as the rest of this file
   describes.", one bullet per mode, and a line that links the
   modes reference and says to read it before running a mode; and the
   index row `| [modes](references/modes.md) | — (the `--follow-up`,
   `--status` and `--unattended` Skill Modes; each command they run lives in its
   own reference) | Running a Skill Mode. |` between the `implement` and
   `profiles` rows. The body after the frontmatter contains no `$` followed by
   a digit or by `ARGUMENTS`.
4. API Contract: the modes reference, first half. `references/modes.md` starts
   `# Skill Modes` and has `## How a mode is chosen`, `## Optional sources` and
   `## --status` sections with the content Build Order 2 lists; on
   whitespace-normalized text its `## How a mode is chosen` section contains
   "first word after the skill's name", "with or without a leading `--`" and
   `` `roundfix follow-up` ``.
5. API Contract: the modes reference, second half. It adds `## --follow-up`,
   an ordered list of eight steps, and `## --unattended`, which contains these
   three sentences verbatim: "Ask at most three questions per call, never more
   than the tool accepts, as the stop-and-ask clause of
   `docs/agents/agent-instructions.md` allows, and put only independent
   questions in one call."; "Answers and authorizations about an item that
   becomes a Spec go into that Spec's `_authorization.md`, and its decisions
   also go into the PRD's `## Decisions`."; "An item that does not become a
   Spec now becomes a Backlog Entry or a Finding that records the information
   and decisions and never an authorization."
6. API Contract: product neutrality. No file under the Roundfix skill names a
   personal assistant or an announcement channel: neither `Jarvis` nor `GT3`
   appears in any letter case.
7. API Contract: commands named. The modes reference names the follow-up
   command, its status view, its message-file option and schema and its
   configuration keys exactly as the tree's `roundfix --help`, the command's
   `--help` and its Roundfix reference state when the Task starts, and names
   no command or flag the tree lacks. Spec 0306 plans them as `roundfix
   follow-up [--status] [--messages <path>] [--format <text|json>]` and the
   schema `roundfix/follow-up-messages/v1`; Spec 0317, delivered before this
   Spec, removes the command's Secondbrain source and deprecates
   `secondbrain.root` and `secondbrain.inbox_namespace`, so the modes
   reference names neither key; the tree, not this plan, is authoritative.

### Surface Transcripts

None. No command, flag, exit code or output line changes.

### Invariants

```text
1. The clause keeps its id, its stop-and-ask enforcement and its rule; only its guidance, the rule version and the module version change.
2. A refreshed adopter records the clause retained; no clause becomes unaccounted.
3. Every derived file changes only through the sanctioned commands; none is hand-edited.
4. A request whose first word after the skill's name is not `follow-up`, `status` or `unattended`, with or without a leading `--`, uses the Roundfix skill exactly as before; the existing sections keep their bytes apart from the version fields and the index.
5. The modes reference's question limit is the one the Baseline clause carries.
6. No authorization is recorded outside a Spec's _authorization.md.
7. No product text names a personal assistant or announcement channel.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1 and 2, User Story 2, Invariants 1-3 → the clause text and retention (API Contracts 1 and 2)
- Goal 2, Core Features 2 and 6, Success Metric 3, User Story 1, Invariant 4 → the skill entry and the first half of the modes reference (API Contracts 3 and 4)
- Goal 3, Core Feature 3, User Story 3 → the `--status` section (API Contract 4)
- Goal 4, Core Feature 4, Success Metric 4, User Story 4 → the `--follow-up` section and optional sources (API Contracts 4, 5 and 7)
- Goal 5, Core Feature 5, Success Metrics 4 and 5, User Story 5, Invariants 5-7 → the `--unattended` section and product neutrality (API Contracts 5 and 6)

## Absence Evidence

1. The three-question clause — command: `grep -rn "at most three questions" internal docs .agents skills`; result: no match at v0.109.0; the clause says "Ask exactly one question per call or message".
2. `references/modes.md` and the modes section — command: `grep -rn -e "## Modes" -e "references/modes.md" -e "opencode/slash" .agents/skills/roundfix skills/roundfix`; result: no match at v0.109.0.
3. `question_batch_clause_test.go` and `roundfix_modes_test.go` — command: `grep -rln -e askUserClauseLimit -e roundfixSection internal skills`; result: no match at v0.109.0.
4. The bare mode word — command: `grep -rn -e "with or without a leading" -e "first word after" .agents/skills/roundfix skills/roundfix`; result: no match at v0.118.0.

## ADR Alternatives

- ADR-0320: no rejected alternative is taken; the clause neither keeps one question per call nor adopts four or each tool's own limit, and no authorization is written into a Backlog Entry or Finding.
- ADR-0330: none of its rejected alternatives is taken; the `--follow-up` mode triages only ai-memory messages, consumes each after recording it, sends another project's finding as an ai-memory message to that project, and says what is missing when the session exposes no ai-memory tools; it never writes the Secondbrain.

## Integration Points

The embedded Baseline catalog and its derived files; the repository's own
Managed Refresh; the owned skill mirrors and version record. Spec 0306's
follow-up command is named by the skill text only. No test, Verification
command or QA row calls a provider, reaches the network, reads the real
`~/.roundfix`, reads the Secondbrain or calls ai-memory.

## Testing Approach

Every new top-level test and subtest calls `t.Parallel()` first and passes no
real-time duration; both new files compile for Windows (no platform API).

- `internal/baseline/question_batch_clause_test.go` (new):
  - `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall`: in
    `embeddedBaselineClauses`, the clause has enforcement `stop-and-ask`, no
    `replaces`, the limit and independence sentences once each on
    whitespace-normalized text, and not the retired sentence.
  - `TestTheAskUserClauseRendersInTheAgentInstructionsGuide`: subtests for the
    Standard TypeScript golden and `../../docs/agents/agent-instructions.md`,
    each holding both sentences once and not the retired one.
  - `TestAnAdopterRetainsTheAskUserClause`: through
    `newClauseReplacementAdopter`, the plan is ready, the delta records the
    clause `retained`, and the retention evidence names it as its own target.
- `skills/roundfix_modes_test.go` (new), over `roundfixLayoutSources()`
  (canonical and embedded), one subtest each:
  - `TestTheRoundfixSkillReadsItsModeFromTheRequest` (Build Order 2): the
    frontmatter `metadata` `opencode/slash` is `"true"`; the body after the
    frontmatter has no match of `\$(ARGUMENTS|[0-9])`; on
    whitespace-normalized text the `## Modes` section contains
    `` `--follow-up` ``, `` `--status` ``, `` `--unattended` ``,
    `` `/roundfix follow-up` ``, `` `/roundfix --follow-up` ``,
    `` `$roundfix status` ``, `` `$roundfix --status` ``,
    `` `@roundfix unattended` ``, `` `@roundfix --unattended` ``,
    `` `roundfix follow-up` ``, `` `/roundfix watch 123` ``,
    "with or without a leading `--`" and `(references/modes.md)`; the
    reference has non-empty `## How a mode is chosen`, `## Optional sources`
    and `## --status` sections, and its `## How a mode is chosen` section
    holds the three phrases of API Contract 4 on whitespace-normalized text.
    Each missing string is named in the failure.
  - `TestTheFollowUpModeRunsItsStepsInOrder` (Build Order 3): the `##
    --follow-up` section holds one ordered list whose items start with `1. `
    to `8. ` in order at the start of a line, an item running to the next
    item or the section's end; the first item that contains `--unattended`
    is item 3, `memory_message_pop` item 4, `--strict` item 6, `roundfix
    deliver` item 7 and `memory_message_send` item 8; and the reference
    contains none of `resolved_to`, `_triaged`, `Inbox Entr` and
    `secondbrain.root`.
  - `TestTheUnattendedModeAsksWithinTheBaselineLimit` (Build Order 3): the
    core module source (`../internal/baseline/assets/modules/core.json`)
    contains `at most three questions per call`; the `## --unattended`
    section contains it and the three sentences of API Contract 5 on
    whitespace-normalized text.
  - `TestTheRoundfixSkillNamesNoPersonalAssistant` (Build Order 3): no file
    under `roundfix/` contains `jarvis` or `gt3` once lower-cased.
- Declared break: `TestRoundfixSkillIndexNamesEveryReference` in
  `skills/roundfix_layout_test.go` expects one more reference file than the
  value the tree holds when Build Order 2 starts (18 at v0.109.0, 19 at
  v0.118.0 after Spec 0306, so 20); no other existing test changes.
- Prototype record: on 2026-10-22, in a scratch worktree, the three Baseline
  tests and the first skill test, with a minimal skill edit, passed after the
  change and failed before it. These sabotages each failed the named case:
  the clause's enforcement set to `mandatory` → catalog validation fails both
  `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall` and
  `TestAnAdopterRetainsTheAskUserClause`; `docs/agents/agent-instructions.md`
  restored to v0.109.0 → the repository subtest of
  `TestTheAskUserClauseRendersInTheAgentInstructionsGuide`; `$ARGUMENTS`
  appended to the skill body → `TestTheRoundfixSkillReadsItsModeFromTheRequest/canonical`;
  the `opencode/slash` line removed → the same subtest; `@roundfix
  --follow-up` removed from the embedded entry only → its `embedded` subtest;
  the `## --status` heading renamed in the canonical reference → its
  `canonical` subtest. With a minimal second half of the reference, the three
  remaining skill tests passed, and these sabotages failed them: the third and
  fourth `--follow-up` items swapped → `TestTheFollowUpModeRunsItsStepsInOrder/canonical`
  (`--unattended` first in item 4); "three" changed to "four" in the
  reference's limit → `TestTheUnattendedModeAsksWithinTheBaselineLimit/canonical`;
  the same change in `core.json` only → the same test before its subtests;
  "and never an authorization" removed → its `canonical` subtest; `GT3`
  added to the embedded reference → `TestTheRoundfixSkillNamesNoPersonalAssistant/embedded`.
- Amendment of 2026-10-24: the prototype predates the bare mode words, and
  the amendment was not prototyped. The `@roundfix --follow-up` sabotage
  becomes `@roundfix unattended` removed from the embedded Modes section only
  (→ its `embedded` subtest), and two sabotages are added for the QA and the
  implementing Agent to prove: the sentence naming `roundfix follow-up`
  removed from the canonical Modes section, and "with or without a leading
  `--`" removed from the canonical reference's `## How a mode is chosen`
  section, each → the `canonical` subtest of
  `TestTheRoundfixSkillReadsItsModeFromTheRequest`.
- Amendment of 2026-10-24 (cycle AP batch 6): `TestTheFollowUpModeRunsItsStepsInOrder`
  now locates `memory_message_pop` in item 4 and `memory_message_send` in
  item 8 and refuses `resolved_to`, `_triaged`, `Inbox Entr` and
  `secondbrain.root` anywhere in the reference. Not prototyped; two
  sabotages are added for the QA and the implementing Agent to prove, each
  on the canonical reference: `memory_message_send` removed from item 8, and
  `resolved_to` appended to item 4, each → the `canonical` subtest.

## Build Order

1. The question clause: the reworded clause, rule version 6, the module
   record, `make baseline-digests`, the Managed Refresh and
   `question_batch_clause_test.go` (no dependency).
2. The skill entry and the first half of the modes reference: the `## Modes`
   section, the metadata and description, the index row, `## How a mode is
   chosen` (how each harness passes the request text; the mode is the first
   word after the skill's name when it is `follow-up`, `status` or
   `unattended`, with or without a leading `--`; it is the Skill Mode, not the
   CLI command `roundfix follow-up`, which the follow-up mode runs; any other
   first word, or none, means ordinary use, as in `/roundfix watch 123`;
   when the words after the mode word name another mode, one question names
   them; a mode stops and names `roundfix upgrade` when the tree's Roundfix lacks the
   follow-up command), `## Optional sources` (the Secondbrain through
   `docs/agents/secondbrain.md`, read-only, as the knowledge base authoring
   consults; ai-memory when its MCP tools are exposed, its messages
   listed without consuming and handed to the follow-up command as a file,
   Jev through the command's own report, a release announcement only when the
   user's own instructions define one; a missing source is named in the report
   and the mode continues), `## --status` (the follow-up command's status
   view; in progress, pending Specs in delivery order, next release, blockers,
   what waits on the user; writes, asks, consumes and arms nothing), the
   layout count, `TestTheRoundfixSkillReadsItsModeFromTheRequest`, the version,
   the mirrors and **Skill Mode** (no dependency on Build Order 1: they share
   no file).
3. The second half of the modes reference: `## --follow-up` with its eight
   steps and `## --unattended` with its rules and answer placement, the three
   remaining skill tests, the version and the mirrors (depends on 1, whose
   clause text its test reads, and on 2, whose files it extends).
4. QA gate (depends on: 1, 2, 3).

No Spec is listed in `requires`. Spec 0306 delivers first; Build Orders 2 and 3
read the follow-up command's name and options from the tree and never invent
them. Specs 0301 to 0306 may raise the Roundfix skill's version and add a
reference; Build Orders 2 and 3 raise both over what the tree holds.

## Risks & Considerations

- An Agent may join dependent questions in one call; the clause names the
  independence condition, and the dependent case waits.
- A user who types `/roundfix follow-up` to run only the CLI command gets
  the whole follow-up mode, whose first step runs that command; the modes
  paragraph and the reference say so, and the read-only `status` mode, or the
  command typed in a shell, gives the bare listing.
- The modes are prose; a harness that changes how it passes the request text
  would break detection silently. The unproven step records the first live
  use, and the skill reads the whole request text rather than one variable.
- `opencode/slash` exposes `/roundfix` in every OpenCode adopter's TUI; it
  only invokes the skill, which already loads on demand.
- The owned interview skills keep one question per message; a reader may see
  that as a contradiction. ADR-0320 records why it is not.

## Glossary

- adds: **Skill Mode**

## Decisions

- The clause is reworded in place and keeps its identity, so the refresh
  records it `retained`; see ADR-0320.
- One topic reference holds the three modes; the follow-up command keeps its
  own reference from Spec 0306.
- The tests hold structure and agreement with the Baseline, and the live
  harness behavior is declared unreachable.
- Amendment of 2026-10-24, the maintainer's answer "Aceitar as duas formas
  (Recomendado)": the first word after the skill's name selects the mode
  with or without a leading `--`, in all three harnesses; only that word
  counts, so a mode word later in an ordinary request selects nothing; the
  bare `follow-up` is the Skill Mode, not the CLI command; the version raise
  stays one patch per skill Task, with no extra release cost.
