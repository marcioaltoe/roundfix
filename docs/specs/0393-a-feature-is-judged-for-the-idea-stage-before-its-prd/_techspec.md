---
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
prd: _prd.md
created: 2026-11-14
---

# A feature is judged for the Idea stage before its PRD — Technical Spec

## Executive Summary

The judge gains one judgment, `idea-stage`: three Noul questions asked in one
request about one item that is not yet a Spec. `roundfix spec judge` gains the
stage `idea`, which takes `--item` in place of a Spec slug, resolves the item
through the Follow-up Command's own inventory, builds the state the follow-up
judge already sends, and prints a verdict with each probability and the lines
of the record. The three route clauses of the Baseline's Spec routing rule keep
their ids and carry the rule to every adopter; `brainstorming`, `write-idea`,
`write-prd` and the Roundfix skill's follow-up mode run and record the
judgment; the glossary adds **Idea Stage Judgment**. ADR-0403 and ADR-0404
record the decisions.

The trade-off accepted: the command can send only an item the follow-up
inventory holds, so a request that exists only in a session is judged by the
agent alone, and the thresholds stay at the documentation's starting points
with a wide undecided band. Both follow from the maintainer's data boundary
and from a measurement too small to fix a number; the record every feature PRD
keeps is what a later Spec needs to narrow the band.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the item is
  named by the reference the Follow-up Command prints. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — one request per judgment to the
  judge's existing recipients, with the existing keys read from the
  environment and never printed, under the existing ceiling. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0403 (new) records the rule: "The
  verdict starts the stage; it refuses nothing." ADR-0404 (new) records what
  the judge reads: "Nothing but an inventoried item can become a request".
  ADR-0200 applies unchanged: "A raised judgment is a finding the authoring
  model answers". ADR-0201 applies unchanged: "Each request appends one line
  to a Judge Log". ADR-0319 applies unchanged: its request state, "one item's
  kind, type, origin, title and an excerpt of its text", is reused as it is.
  ADR-0400 applies: one context in three implementation Tasks that each earn
  their verification cycle. ADR-0209 does not apply: this judgment reads no
  Finding and does not use the grouping reader. ADR-0231 and ADR-0239 apply
  unchanged: the ceiling and the key order are theirs. ADR-0283
  applies unchanged: inside a Run the judge sends nothing. ADR-0353 applies
  unchanged: a Waiting Item is never sent to the judge. ADR-0208, ADR-0235,
  ADR-0273, ADR-0383, ADR-0390 and ADR-0396 cite one of those and do not
  apply, for the reasons `_prd.md` → Project Constraints gives: this Spec
  changes no grouping clause, starts no Agent Session, adds no key variable or
  Run child, and changes neither the message rules nor the follow-up order.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — the maintainer's answer of 2026-11-14
  "Mesmas do BE (Recommended)", quoted in `_authorization.md`; the Governed
  Paths are the spec-workflow module, the Standard TypeScript profile and its
  formatter golden, and this repository's managed Spec routing guide and setup
  manifest, which task_02 changes through their commands, and the skill files
  task_01 and task_03 change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0393-a-feature-is-judged-for-the-idea-stage-before-its-prd/_authorization.md`;
  bounded files: `.agents/skills/brainstorming/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-idea/references/idea-template.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/brainstorming/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`.

## Current behavior

Read on 2026-11-14 at `251ed81e`; `_prd.md` → Validation records the runs.

1. `internal/judge/questions.json` holds six judgments, one `pinned_model`
   (`jev-1.13`) with its `accepted_model_pattern`, the two transports and the
   language gate; `Load` in `internal/judge/questions.go` reads them into
   `Questions`.
2. `PrioritizeFollowUp` in `internal/judge/follow_up.go` is the only judgment
   that reads an item that is not a Spec. For each item it cuts the text to
   `item_max_chars` (4,000) runes, scans the title, origin, excerpt, type and
   kind with `secretscan.ScanContent`, and sends the state
   `{"project": …, "item": {"kind", "type", "origin", "title", "text"}}`,
   whatever the item's language. It skips every item inside a Run, without a
   key, at the ceiling and with an unreadable Judge Log, retries HTTP 429 and
   529 up to three attempts, checks the ceiling before each attempt and
   appends a Judge Log row per attempt. `client.ask` in
   `internal/judge/client.go` sends a map of questions only for the kind
   `follow-up-priority`, and `evaluate` reads a map of Nouls only for that
   kind; an answer from a model outside the pin is `skipped` with `answered by
   <model>, thresholds belong to <pinned model>`.
3. `parseSpecJudgeCommand` and `runSpecJudgeCommand` in
   `internal/cli/spec_judge.go` accept one Spec slug, `--stage
   <prd|techspec|tasks>` and `--format <text|json>`. Any other stage fails
   with `unsupported --stage "<value>"; use prd, techspec or tasks`, and every
   failure prints `roundfix: spec judge failed: <error>` and `Run 'roundfix
   spec judge --help' for usage.` on standard error and exits 2. The command
   hands the judge only the variables `Questions.KeyVariables` names, and none
   when the Run Marker is set. `judge.Run` calls `PlanSpec`, which fails
   without `_prd.md`.
4. `followUpInventory` in `internal/cli/follow_up_inventory.go` returns the
   Follow-up Items: the entries of a `--messages` file (kind `message`,
   reference `message:<id>`, type the entry's hint or `message`), open Backlog
   Entries (kind `backlog`, type from the front matter), unresolved Findings
   (kind and type `finding`), notes under `docs/_inbox/` (kind and type
   `note`) and active Specs (kind `spec`), each with its title, text, origin
   and, for a Backlog Entry or a Finding with a valid Waiting Declaration, its
   `WaitingOn` and `Since`. `followUpNamespace` gives the project's name. No
   non-test Go file reads `_idea.md`.
5. `rule.spec.routing` in `internal/baseline/assets/modules/spec-workflow.json`
   (module version 24) holds `clause.spec.routing-01-large-initiative`,
   "Large or fuzzy product initiative: run `write-idea` → `write-prd` →
   `write-techspec` → `write-tasks`."; `clause.spec.routing-02-standard-feature`,
   "Standard feature that changes product behavior: run `write-prd` →
   `write-techspec` → `write-tasks`; skip the TechSpec only when the feature
   has no architectural surface."; and `clause.spec.routing-05-task-graph`,
   which begins "Use `brainstorming` before creative or feature work, start
   with the smaller sufficient route when two routes fit, and execute
   implementation from the Task Graph." and ends with the Surface Transcript
   sentence. They render into this repository's `docs/agents/spec-routing.md`,
   the Standard TypeScript formatter golden of that guide, the profile pin,
   the catalog snapshots and the plan goldens.
   `internal/baseline/authoring_evidence_clauses_test.go` pins the Surface
   Transcript sentence once in the module, the golden and the repository
   guide; `internal/baseline/clause_characterization_test.go` pins the three
   ids and their force. The Source Baseline
   `baseline.standard-typescript-monorepo-0.0.1` holds the three ids with
   that baseline's wording and is frozen history.
6. `.agents/skills/brainstorming/SKILL.md` routes "Product-level idea" to
   `write-idea` and "Feature ready to specify" to `write-prd`.
   `.agents/skills/write-prd/SKILL.md` reads an existing `_idea.md` as "the
   exploration" and never sends a request to `write-idea`.
   `.agents/skills/write-idea/SKILL.md` says when the stage earns its cost and
   nothing about who starts it. `.agents/skills/roundfix/references/modes.md`
   plans in step 2, gathers decisions in step 3 and authors in step 5 with no
   word about the Idea stage. `.agents/skills/roundfix/references/spec.md`
   documents the judge under `### Advisory judge`, and
   `docs/user-guide/commands/spec.md` under `### spec judge`; both state the
   unknown-stage message.
7. `CONTEXT.md` defines **Advisory Judgment** as a judgment "about a Spec
   artifact" or "about a source relationship"; no entry names a judgment of
   the Idea stage.

## System Architecture

- `internal/judge`: the `idea-stage` judgment in `questions.json`, its loader
  fields, and `JudgeIdeaStage` in a new `idea_stage.go`; `client.ask` and
  `evaluate` treat `idea-stage` as the second kind that asks several Nouls in
  one request (task_01).
- `internal/cli`: the stage `idea` of `roundfix spec judge` in a new
  `spec_judge_idea.go`, called from `runSpecJudgeCommand`; it resolves the
  item with `followUpInventory` and renders text and JSON (task_01).
- `internal/baseline/assets/modules/spec-workflow.json`: the three reworded
  route clauses, their derived files and this repository's managed refresh
  (task_02).
- The owned skills, their templates and tests, and `CONTEXT.md` (task_03); the
  `spec` command guide and the Roundfix skill's `spec` reference travel with
  the command (task_01).

No new package, no new dependency and no change to the Run Database, the
Follow-up Command, the archive or any other stage of the judge.

## Implementation Design

### Interfaces

In `internal/judge/questions.go`, beside `FollowUpJudgment`, with `Questions`
gaining the field `IdeaStage IdeaStageJudgment` loaded from the key
`idea-stage`:

```go
type IdeaStageCriterion struct {
	ID       string   `json:"id"`
	Summary  string   `json:"summary"`
	Question Question `json:"question"`
}

type IdeaStageJudgment struct {
	QuestionID     string               `json:"question_id"`
	Criteria       []IdeaStageCriterion `json:"criteria"`
	YesAtLeast     float64              `json:"yes_at_least"`
	NoAtMost       float64              `json:"no_at_most"`
	ItemMaxChars   int                  `json:"item_max_chars"`
	NotJudgedTypes []string             `json:"not_judged_types"`
}
```

In the new `internal/judge/idea_stage.go`:

```go
type IdeaStageRequest struct {
	RunID, RepoRoot, Project string
	Item                     FollowUpItem
	Keys                     map[string]string
	HomeDir                  string
	Transport                http.RoundTripper
	Now                      func() time.Time
}

type IdeaStageAnswer struct {
	ID          string
	Probability float64
	Answer      string // "yes", "no" or "unsure"
}

type IdeaStageReport struct {
	State                                  string // "answered", "judge_unavailable" or "not_applicable"
	Reason                                 *string
	Verdict                                string // "indicated", "not indicated", "undecided" or ""
	Answers                                []IdeaStageAnswer
	Model, Transport, KeyVariable          *string
	Calls, InputTokens                     int
	CostUSD, MonthCostUSD, MonthCeilingUSD float64
}

func JudgeIdeaStage(ctx context.Context, q Questions, req IdeaStageRequest) (IdeaStageReport, error)
func IdeaStageVerdict(answers []string) string
```

In `internal/cli`: `specJudgeRequest` gains `item` and `messages`; a new
unexported `ideaStageItem(items []followUpItem, reference string)
(followUpItem, error)` is the only place a reference becomes an item; and
`renderIdeaStageText(doc ideaStageDocument) string` renders the report.

### Data Models

The `idea-stage` entry of `internal/judge/questions.json`, inserted after
`follow-up-priority`, word for word (these are the questions as measured):

   ````json
   "idea-stage": {
     "question_id": "idea_stage",
     "criteria": [
       {
         "id": "shape_open",
         "summary": "the right shape of the solution is still open",
         "question": {
           "type": "noul",
           "instructions": "`item` is one pending request addressed to the software project `project`. Is the right shape of the solution for `item` still open?",
           "criteria": {
             "true": "At least two materially different designs could satisfy the item, and the item does not settle which one to build.",
             "false": "The item names one concrete design, or only one reasonable design exists."
           }
         }
       },
       {
         "id": "unverified_assumption",
         "summary": "it rests on an assumption about its users, the business or the market that it gives no evidence for",
         "question": {
           "type": "noul",
           "instructions": "`item` is one pending request addressed to the software project `project`. Does `item` rest on an assumption about its users, the business or the market that it gives no evidence for?",
           "criteria": {
             "true": "A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong.",
             "false": "The need is shown by evidence in the item itself, such as a measured defect, a reported failure or a decision the product owner already made."
           }
         }
       },
       {
         "id": "market_question",
         "summary": "deciding to build it depends on a market or competitive question",
         "question": {
           "type": "noul",
           "instructions": "`item` is one pending request addressed to the software project `project`. Does deciding whether to build `item` depend on a market or competitive question?",
           "criteria": {
             "true": "The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it.",
             "false": "The decision depends only on the product's own behavior and its current users."
           }
         }
       }
     ],
     "yes_at_least": 0.7,
     "no_at_most": 0.3,
     "item_max_chars": 4000,
     "not_judged_types": [
       "fix",
       "perf",
       "refactor",
       "finding"
     ]
   }
   ````

The JSON report of the command, schema `roundfix/idea-stage-judgment/v1`, keys
in this order. `reason`, `verdict`, `judged_by`, `model`, `transport` and
`key_variable` are `null` when they do not apply, and a criterion's
`probability` and `answer` are `null` unless the judge answered:

```json
{"schema":"roundfix/idea-stage-judgment/v1","item":{"ref":"docs/backlog/2026-11-02-teams-can-see-each-others-run-results.md","kind":"backlog","type":"feat","title":"Teams can see each other's Run results"},"state":"answered","reason":null,"verdict":"indicated","judged_by":"judge","criteria":[{"id":"shape_open","summary":"the right shape of the solution is still open","yes":"At least two materially different designs could satisfy the item, and the item does not settle which one to build.","no":"The item names one concrete design, or only one reasonable design exists.","probability":0.98,"answer":"yes"}],"thresholds":{"yes_at_least":0.7,"no_at_most":0.3},"model":"jev-1.13.0","transport":"typesafe","key_variable":"ROUNDFIX_TYPESAFE_API_KEY","calls":1,"input_tokens":512,"cost_usd":0.000021504,"month_cost_usd":0.012,"month_ceiling_usd":5}
```

(The example shows one of the three criteria.) The Judge Log gains no field:
each attempt appends one `roundfix/judge-log/v1` row with `judgment`
`idea-stage`, `artifact` the item's reference, `target` its kind,
`question_id` `idea_stage`, the three probabilities under `probabilities`, and
`outcome` `answered` or `skipped`.

The record in an Idea and a PRD is the section `## Idea-stage judgment` with
these top-level lines, in this order (API Contract 7 holds the templates):

```text
- item: <the judged reference, or `the request in the session`>
- judged-by: judge (<model>) | agent (<why the judge did not answer>) | nobody (the user asked for the Idea stage)
- shape_open: <yes|no|unsure>[ (<probability>)]
- unverified_assumption: <yes|no|unsure>[ (<probability>)]
- market_question: <yes|no|unsure>[ (<probability>)]
- verdict: indicated | not indicated | undecided | not judged
- stage: ran | dismissed — <the user's reason> | not run
- shape: changed | split | shelved | unchanged | not applicable[ — <one line>]
```

### API Contracts

The four-backtick blocks below are indented by three spaces that are not part
of their text.

1. **The judgment.** `JudgeIdeaStage` judges `req.Item` and never returns an
   error for a judgment; its checks run in this order, and the first that
   holds decides the report:
   1. The item's `Type` is one of `not_judged_types`: `State`
      `not_applicable`, `Reason` `the item is typed <type>; only a feature is
      judged for the Idea stage`. Nothing is read or sent.
   2. `req.RunID` is not empty: `State` `judge_unavailable`, `Reason` `inside
      Run <run-id>, named by ROUNDFIX_RUN_ID`. No transport is selected, no
      key or Judge Log is read and nothing is sent.
   3. No key variable is set: `Reason` `<first variable> is not set (nor
      <second>, nor <third>)`, the form `judge.Run` already prints.
   4. The Judge Log cannot be read: `judge log unreadable: <error>`. The
      month's cost is at or above the ceiling: `monthly ceiling reached
      (US$<cost> of US$<ceiling>)`, as the other judgments format it.
   5. The title, origin, excerpt, type and kind, joined by newlines as
      `PrioritizeFollowUp` joins them, match a secret pattern: `not sent:
      secret pattern <kind>`. Nothing is sent.
   6. Otherwise one request carries the state
      `{"project": <Project>, "item": {"kind", "type", "origin", "title",
      "text"}}`, the text cut to `item_max_chars` runes, and the three
      questions keyed by their ids. Attempts, the retry of HTTP 429 and 529,
      the ceiling check before each attempt, key redaction and the Judge Log
      row per attempt follow `PrioritizeFollowUp`. An answer counts only when
      every criterion returned a Noul between 0 and 1 and the reported model
      matches the pin; otherwise `State` is `judge_unavailable` with the
      reason the other judgments use (`answered by <model>, thresholds belong
      to <pinned model>`, `unreadable answer`, `service unavailable (…)`,
      `key refused (HTTP <n>)`, `request refused (HTTP <n>)`, `judge log not
      writable: <error>`).
   An answered report has `State` `answered`, one `IdeaStageAnswer` per
   criterion in the order of `criteria`, each `yes` at `yes_at_least` or
   more, `no` at `no_at_most` or less and `unsure` otherwise, and `Verdict`
   from `IdeaStageVerdict`: `indicated` when any answer is `yes`, `not
   indicated` when every answer is `no`, `undecided` otherwise.
2. **The command.** `roundfix spec judge --stage idea --item <reference>
   [--messages <path>] [--format <text|json>]`. With `--stage idea` the
   command takes no Spec slug, calls `followUpInventory` with the messages
   path resolved as the Follow-up Command resolves it, and finds the item
   whose reference equals `--item` byte for byte through `ideaStageItem`. A
   Waiting Item is reported `not_applicable` with the reason `the item waits
   on <waiting_on> since <since>` before the judge is called. The project is
   `followUpNamespace`; the keys, the ceiling, the transport and the clock
   are passed as `runSpecJudgeCommand` passes them today. Exit `0` whenever a
   report is printed. Exit `2`, with `roundfix: spec judge failed: <error>`
   and `Run 'roundfix spec judge --help' for usage.` on standard error and
   nothing on standard output, for these errors, word for word:
   - `--stage idea requires --item, the reference of the item to judge`
   - `--stage idea takes no Spec slug; name the item with --item`
   - `--item applies only to --stage idea`
   - `--messages applies only to --stage idea`
   - `no pending item "<reference>"; name a reference that 'roundfix follow-up' lists`
   - `"<reference>" is an active Spec; the Idea stage is judged before a PRD exists`
   - `unsupported --stage "<value>"; use idea, prd, techspec or tasks`
   - an unreadable or invalid messages file, with the Follow-up Command's own
     error text for it.
3. **Text output.** Probabilities print with two decimals. When the judge
   answered:

   ```text
   Idea stage for <reference>: <verdict>
     shape_open: <answer> (<probability>) — the right shape of the solution is still open
     unverified_assumption: <answer> (<probability>) — it rests on an assumption about its users, the business or the market that it gives no evidence for
     market_question: <answer> (<probability>) — deciding to build it depends on a market or competitive question
   Rule: indicated when any criterion is yes, not indicated when every criterion is no, undecided otherwise; yes at 0.70 or more, no at 0.30 or less.
   Record:
     - item: <reference>
     - judged-by: judge (<model>)
     - shape_open: <answer> (<probability>)
     - unverified_assumption: <answer> (<probability>)
     - market_question: <answer> (<probability>)
     - verdict: <verdict>
   Judge: <calls> call(s), <input tokens> input tokens, US$<cost>; month US$<month cost> of US$<ceiling>; model <model> via <transport> on <key variable>
   ```

   The `Judge:` line formats its numbers as the summary of the other stages
   does. When the judge did not answer (`judge_unavailable`), the output is
   Surface Transcript 4, followed by the same `Judge:` line only when at least
   one request was sent. When the item is `not_applicable`, the output is the
   one line of Surface Transcript 3. Each criterion's summary and its `yes`
   and `no` lines are read from `questions.json`, never from a second copy.
4. **Help.** `specJudgeUsage` keeps every existing line and gains: the
   synopsis line `  roundfix spec judge --stage idea --item <reference>
   [--messages <path>] [--format <text|json>]` right after the existing one;
   a paragraph after the first one that reads "With --stage idea it judges,
   before a PRD exists, whether one item that is not yet a Spec needs the
   Idea stage: an open Backlog Entry, a note in docs/_inbox/ or, with
   --messages, a message entry, named by the reference that 'roundfix
   follow-up' prints. It never gates."; the option lines `--item` ("With
   --stage idea, the reference of the item to judge") and `--messages` ("With
   --stage idea, the JSON file of ai-memory messages that holds the item");
   `idea` in the `--stage` option line; and "or unknown item" at the end of
   the exit code 2 line. `specUsage` and the top-level `usage` gain the same
   synopsis line after their existing `spec judge` line.
5. **Baseline clauses.** In `rule.spec.routing`, with the rule and module
   versions raised by the record command, three clauses keep their ids,
   `enforcement: mandatory` and the absence of `replaces`; every other clause
   is unchanged, and the `guidance` texts become, word for word:
   - `clause.spec.routing-01-large-initiative`:

   ````text
   Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`. A feature takes the same route when its Idea stage is indicated. Before writing the PRD of a feature that has no `_idea.md` and that the user did not send to `write-idea`, judge three criteria about the request: whether the right shape of the solution is still open, whether the request rests on an assumption about its users, the business or the market that it gives no evidence for, and whether deciding to build it depends on a market or competitive question. When the request is recorded as an item that is not yet a Spec, which is an open Backlog Entry, a note in `docs/_inbox/` or a message handed in as a file, run `roundfix spec judge --stage idea --item <reference>` and take the judge's answers; when the command reports that the judge did not answer, or the request has no such record, answer the same three criteria yourself, each as yes, no or unsure, and say that the judge did not answer. The Idea stage is indicated when any criterion is yes, not indicated when every criterion is no, and undecided otherwise. Indicated: start `write-idea` without asking, show the verdict, each criterion's answer with its probability and who judged, and let the user dismiss the stage at its first question. Undecided, or in doubt: the user decides, seeing the same answers. Not indicated: continue to `write-prd`. Record who judged, each answer, the verdict, whether the stage ran or was dismissed and why, and what the Idea did to the shape of the request, in the Idea when one is written and always in the PRD. The judgment informs and never refuses: no command fails for it, and a dismissal needs only its recorded reason. When the stage is indicated or undecided and no user is present to answer, write no PRD for that feature: keep it as a Backlog Entry that records the judgment and waits, with `waiting_on` and `since` in its frontmatter, and never answer the stage's questions or the user's decision yourself. A fix or a refactor is never judged for the Idea stage.
   ````

   - `clause.spec.routing-02-standard-feature`:

   ````text
   Standard feature that changes product behavior: run `write-prd` → `write-techspec` → `write-tasks` once its Idea stage is judged not indicated or the user dismissed it; skip the TechSpec only when the feature has no architectural surface.
   ````

   - `clause.spec.routing-05-task-graph`: its first sentence becomes the
     sentence below, and the Surface Transcript sentence after it stays byte
     for byte:

   ````text
   Use `brainstorming` before creative or feature work, start with the smaller sufficient route when two routes fit, except that the Idea-stage judgment, never the smaller route, decides whether a feature starts at `write-idea`, and execute implementation from the Task Graph.
   ````

6. **Skill texts.** In the canonical files, word for word; `make skills-sync`
   copies each to its mirror, and no other byte of these files changes.
   - `.agents/skills/brainstorming/SKILL.md`: the bullet "- **Feature ready
     to specify** → `write-prd` (synthesizes this conversation; it will not
     re-interview)." gains, before its final period, ", once the Idea-stage
     judgment below does not send it to `write-idea`"; and this section is
     inserted right before `## Key principles`:

   ````markdown
   ## Idea-stage judgment

   Before a feature goes to `write-prd`, judge whether it needs the Idea stage, unless the user asked for `write-idea` or the feature already has an `_idea.md`. The criteria are three: the right shape of the solution is still open; the request rests on an assumption about its users, the business or the market that it gives no evidence for; deciding to build it depends on a market or competitive question.

   When the repository runs Roundfix and the request is recorded as an open Backlog Entry, a note in `docs/_inbox/` or a message handed in as a file, run `roundfix spec judge --stage idea --item <reference>` and take the judge's answers. Otherwise, or when the command says the judge did not answer, answer the three criteria yourself, each as yes, no or unsure, and say that the judge did not answer.

   - Any criterion yes: the Idea stage is indicated. Start `write-idea` and show the answers; the user may dismiss the stage at its first question.
   - Every criterion no: go to `write-prd`.
   - Anything else, or any doubt: ask the user whether the Idea stage runs, showing the answers.

   Never present your own answers as the judge's, and never answer the Idea stage's questions for the user. Hand the answers to the skill that runs next, which records them. A fix or a refactor is never judged.
   ````

   - `.agents/skills/write-idea/SKILL.md`: this section is inserted right
     before `## Ground rules`:

   ````markdown
   ## When the Idea-stage judgment starts this stage

   The repository's Spec routing starts this stage without the user asking when the Idea-stage judgment of a feature reads indicated, and when the user chose it for an undecided one. Then open the first question by showing the judgment: the verdict, each criterion's answer, with its probability when the judge answered, and who judged, the judge with its model or the agent because the judge did not answer. Give that first question one last option that dismisses the stage, and say that a dismissal is recorded with its reason. When the user dismisses, stop: write no `_idea.md`, and hand the judgment and the user's reason to `write-prd`, which records them.

   Never answer this stage's questions for the user. When no user is present, as in a delegated subagent, do not run the stage: the feature stays a Backlog Entry that records the judgment and waits, with `waiting_on` and `since` in its frontmatter.

   Write the judgment into the `## Idea-stage judgment` section of `_idea.md`, in the fielded form the template shows. When the draft is approved, set its `shape` line to what this stage did to the request as it arrived, `changed`, `split`, `shelved` or `unchanged`, followed by one line that says what.
   ````

   - `.agents/skills/write-prd/SKILL.md`: in `## Inputs`, "Clarify only what
     the idea left open." is followed, in the same paragraph, by "Copy its
     `## Idea-stage judgment` section into the PRD unchanged."; and this
     section is inserted right before `## Ground rules`:

   ````markdown
   ## Idea-stage judgment

   A feature PRD is written only after its Idea stage was judged. When `_idea.md` exists, the judgment is its `## Idea-stage judgment` section. When the conversation, or the brief that delegated this skill, already holds the judgment of this request, with the user's dismissal or choice, record that one. Otherwise judge now, before any other step, as the repository's Spec routing says. The criteria are three: the right shape of the solution is still open; the request rests on an assumption about its users, the business or the market that it gives no evidence for; deciding to build it depends on a market or competitive question. When the request is recorded as an open Backlog Entry, a note in `docs/_inbox/` or a message handed in as a file, run `roundfix spec judge --stage idea --item <reference>`; when the command says the judge did not answer, or the request has no such record, answer the three criteria yourself, each as yes, no or unsure, and say that the judge did not answer.

   - Indicated, any criterion yes: stop this skill and start `write-idea`, showing the verdict and the answers. The user may dismiss the stage at its first question; then this skill resumes.
   - Undecided, no criterion yes and not every one no, or any doubt: ask the user whether the Idea stage runs, showing the same answers.
   - Not indicated, every criterion no: continue.

   The judgment informs and never refuses, and a fix or a refactor is never judged. Never present your own answers as the judge's. When the stage is indicated or undecided and no user is present, as in a delegated subagent, write no PRD: keep the feature as a Backlog Entry that records the judgment and waits, with `waiting_on` and `since` in its frontmatter.

   Record the judgment in the PRD's `## Idea-stage judgment` section, in the fielded form the template shows: the item, who judged, each criterion's answer, the verdict, `stage` (`ran`, `dismissed` with the user's reason, or `not run`) and `shape`, what the Idea did to the request (`changed`, `split` or `unchanged`, or `not applicable` when the stage did not run). Every feature PRD carries this section.
   ````

   - `.agents/skills/roundfix/references/modes.md`, under `## --follow-up`
     and `## --unattended`, each sentence appended to the end of the paragraph
     or step named:
     - step 2: "For each planned feature that has no Idea, run the Idea-stage
       judgment as `docs/agents/spec-routing.md` says, with `roundfix spec
       judge --stage idea --item <reference>` and, for a message entry, the
       same messages file, and show its verdict and answers in the plan: an
       indicated feature starts at `write-idea`, an undecided one waits for
       the user's choice, and a not indicated one goes to its PRD."
     - step 3: "The questions of an indicated Idea stage, and the choice for
       an undecided one, are asked here, while the user is present and before
       any Spec is authored; the user may dismiss an indicated stage at its
       first question."
     - step 5: "A feature whose Idea stage is indicated or undecided is
       authored only after the user answered; without an answer it stays a
       Backlog Entry that records the judgment and waits, its Waiting
       Declaration naming the Idea stage."
     - `## --unattended`, a new paragraph after the one that begins "An item
       that does not become a Spec now": "The Idea stage is the user's: this
       mode asks its questions and never answers them, and it never decides
       an undecided Idea-stage judgment. A feature whose judgment is
       indicated or undecided and gets no answer becomes or stays a Backlog
       Entry that records the judgment, with a Waiting Declaration whose
       `waiting_on` names the user's answers to the Idea stage, and the
       report names it."
   - `.agents/skills/roundfix/references/spec.md`, under `### Advisory
     judge`: the bash block gains the line `roundfix spec judge --stage idea
     --item <reference> [--messages <path>] [--format <text|json>]` after its
     last line; the unknown-stage sentence names `idea, prd, techspec or
     tasks`; and this section is inserted right before `## Spec close audit`:

   ````markdown
   ### Idea-stage judgment

   ```bash
   roundfix spec judge --stage idea --item docs/backlog/<file>.md
   roundfix spec judge --stage idea --item message:<id> --messages <path> --format json
   ```

   Before the PRD of a feature is written, `--stage idea` judges whether the
   feature needs the Idea stage. It takes no Spec slug: `--item` names one item
   that is not yet a Spec, by the reference `roundfix follow-up` prints for it,
   an open Backlog Entry, a note in `docs/_inbox/` or, with the same
   `--messages` file, a message entry. The judge receives the project's name and
   that one item's kind, type, origin, title and the first 4,000 characters of
   its text, after the secret scan, in whatever language it is written; nothing
   else is sent.

   The judge answers three criteria, each as the probability of yes:
   `shape_open` (the right shape of the solution is still open),
   `unverified_assumption` (it rests on an assumption about its users, the
   business or the market that it gives no evidence for) and `market_question`
   (deciding to build it depends on a market or competitive question). An answer
   is yes at 0.70 or more, no at 0.30 or less and unsure between. The verdict is
   `indicated` when any criterion is yes, `not indicated` when every criterion is
   no and `undecided` otherwise. The output always shows each answer with its
   probability, and ends with a `Record:` block to copy into the
   `## Idea-stage judgment` section of the Idea or the PRD.

   Follow the repository's Spec routing: indicated starts `write-idea` and shows
   the answers, and the user may dismiss the stage at its first question;
   undecided is the user's choice, seeing the same answers; not indicated goes to
   `write-prd`. The judgment never gates and changes no other command.

   The command exits `0` whenever it ran. When the judge did not answer (inside a
   Run, without a key, at the monthly ceiling, with an unreadable Judge Log, on a
   service failure, for an item that matches a secret pattern, or when another
   model answered), the first line reads `not judged: <reason>` and the output
   lists the three criteria with their yes and no meanings: answer them yourself,
   say that the judge did not answer, and never present your answers as the
   judge's. An item typed `fix`, `perf`, `refactor` or `finding` and a Waiting
   Item are not judged and nothing is sent; an active Spec, an unknown reference
   and a usage error exit `2`. The keys, the Judge Log, the monthly ceiling and
   the Run Marker are those of the Advisory Judge above.

   `--format json` prints one document with the schema
   `roundfix/idea-stage-judgment/v1`: the `item` (its `ref`, `kind`, `type` and
   `title`), the `state` (`answered`, `judge_unavailable` or `not_applicable`),
   the `reason` when the item was not judged, the `verdict`, `judged_by`, the
   `criteria` (each with its `id`, `summary`, `yes` and `no` meanings,
   `probability` and `answer`), the `thresholds`, and the judge's `model`,
   `transport`, `key_variable`, `calls`, `input_tokens`, `cost_usd`,
   `month_cost_usd` and `month_ceiling_usd`.
   ````

   - `docs/user-guide/commands/spec.md`: the synopsis block of `### spec
     judge` gains the same line; its unknown-stage sentence names `idea, prd,
     techspec or tasks`; and a section `### spec judge --stage idea`, with
     the body of the section above from its first bash block to its last
     paragraph, is inserted right after the paragraph that holds the
     unknown-stage sentence.
7. **The record templates.** In
   `.agents/skills/write-idea/references/idea-template.md`, inserted right
   before `## Decisions` inside the template's fenced block:

   ````markdown
   ## Idea-stage judgment

   <!-- The judgment that started this stage, in this fielded form; write-prd copies the section into the PRD.
        judged-by: `judge (<model>)`, `agent (<why the judge did not answer>)` or `nobody (the user asked for the Idea stage)`.
        Each criterion: yes, no or unsure, with the probability in parentheses when the judge answered; omit the three lines when nobody judged.
        verdict: indicated | not indicated | undecided | not judged.
        stage: ran. shape: changed | split | shelved | unchanged, then a dash and one line that says what. -->

   - item: <reference of the judged item, or `the request in the session`>
   - judged-by: judge (<model>)
   - shape_open: yes (0.85)
   - unverified_assumption: unsure (0.41)
   - market_question: no (0.09)
   - verdict: indicated
   - stage: ran
   - shape: changed — <what the Idea changed>
   ````

   In `.agents/skills/write-prd/references/prd-template.md`, inserted right
   before `## Glossary` inside the template's fenced block:

   ````markdown
   ## Idea-stage judgment

   <!-- Every feature PRD carries this record, in this fielded form; the minimal PRD of a fix or a refactor omits the section.
        Copy it from `_idea.md` when an Idea was written.
        judged-by: `judge (<model>)`, `agent (<why the judge did not answer>)` or `nobody (the user asked for the Idea stage)`.
        Each criterion: yes, no or unsure, with the probability in parentheses when the judge answered; omit the three lines when nobody judged.
        verdict: indicated | not indicated | undecided | not judged.
        stage: ran | dismissed — <the user's reason> | not run.
        shape: changed | split | unchanged, then a dash and one line that says what the Idea changed; not applicable when the stage did not run. -->

   - item: <reference of the judged item, or `the request in the session`>
   - judged-by: judge (<model>)
   - shape_open: no (0.12)
   - unverified_assumption: no (0.17)
   - market_question: no (0.07)
   - verdict: not indicated
   - stage: not run
   - shape: not applicable
   ````

8. **Glossary.** In `CONTEXT.md`, word for word. In **Advisory Judgment**,
   "or about a source relationship, such as whether a Finding or Backlog
   Entry shares a Spec's context, that the authoring model must answer"
   becomes "about a source relationship, such as whether a Finding or Backlog
   Entry shares a Spec's context, or about an item that is not yet a Spec,
   such as whether a feature needs the Idea stage, that the authoring model
   must answer", with the comma after "attributed claim" kept, and its
   closing parentheses gain a third link, to ADR-0403, written as the
   glossary's Markdown link to that file under `docs/adr/`. A new entry goes
   right after that entry's `_Avoid_` line, separated by one blank line, each
   ADR in its parentheses written as the same kind of link:

   ````text
   **Idea Stage Judgment**:
   The Advisory Judgment, made before a feature's PRD is written, of whether the feature needs the Idea stage, the study `write-idea` runs. It answers three yes or no criteria about the request: the right shape of the solution is still open, the request rests on an assumption it gives no evidence for, and deciding to build it depends on a market question. The judge answers them for an open Backlog Entry, a repository inbox note or a handed-in message through `roundfix spec judge --stage idea`, and the agent answers them, saying so, when the judge does not. The verdict is indicated when any criterion is yes, not indicated when every criterion is no and undecided otherwise: an indicated stage starts without being asked and the user may dismiss it with a recorded reason, an undecided one is the user's choice, and no command refuses anything. The PRD of every feature Spec records who judged, the answers, whether the stage ran and what it did to the shape of the request ([ADR-0403](docs/adr/0403-a-feature-is-judged-for-the-idea-stage-before-its-prd-and-the-judgment-never-refuses.md), [ADR-0404](docs/adr/0404-the-idea-stage-judge-reads-one-item-that-is-not-yet-a-spec.md)).
   _Avoid_: Idea gate, discovery check, routing verdict
   ````

   The texts, with those links, passed `TestGlossaryPhrasesAreDefinedTerms`
   and the other glossary contracts on the prototype.

### Surface Transcripts

The fixture is a committed Git repository with an `origin` remote, holding
`README.md` and two open Backlog Entries: `docs/backlog/2026-11-02-teams-can-see-each-others-run-results.md`,
typed `feat`, titled "Teams can see each other's Run results", and
`docs/backlog/2026-11-03-status-prints-the-header-twice.md`, typed `fix`.
`<reason>` is the reason the judge did not answer in the environment that runs
the command: without a judge key it names the key variables, and inside a Run
it names the Run Marker. These transcripts prescribe a surface that does not
exist yet; the Task's named tests assert them, and the reasons and the error
suffix are copied from the current binary.

1. Surface Transcript: the stage needs an item.

   ```transcript
   $ roundfix spec judge --stage idea
   stdout:
   stderr:
   roundfix: spec judge failed: --stage idea requires --item, the reference of the item to judge
   Run 'roundfix spec judge --help' for usage.
   exit: 2
   ```

2. Surface Transcript: a reference the inventory does not hold is refused.

   ```transcript
   $ roundfix spec judge --stage idea --item docs/backlog/2026-11-09-missing.md
   stdout:
   stderr:
   roundfix: spec judge failed: no pending item "docs/backlog/2026-11-09-missing.md"; name a reference that 'roundfix follow-up' lists
   Run 'roundfix spec judge --help' for usage.
   exit: 2
   ```

3. Surface Transcript: a fix is not judged.

   ```transcript
   $ roundfix spec judge --stage idea --item docs/backlog/2026-11-03-status-prints-the-header-twice.md
   stdout:
   Idea stage for docs/backlog/2026-11-03-status-prints-the-header-twice.md: not judged: the item is typed fix; only a feature is judged for the Idea stage
   stderr:
   exit: 0
   ```

4. Surface Transcript: when the judge does not answer, the criteria are the
   agent's to answer.

   ```transcript
   $ roundfix spec judge --stage idea --item docs/backlog/2026-11-02-teams-can-see-each-others-run-results.md
   stdout:
   Idea stage for docs/backlog/2026-11-02-teams-can-see-each-others-run-results.md: not judged: <reason>
   Answer the three criteria yourself, each as yes, no or unsure, and say that the judge did not answer:
     shape_open: the right shape of the solution is still open
       yes: At least two materially different designs could satisfy the item, and the item does not settle which one to build.
       no: The item names one concrete design, or only one reasonable design exists.
     unverified_assumption: it rests on an assumption about its users, the business or the market that it gives no evidence for
       yes: A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong.
       no: The need is shown by evidence in the item itself, such as a measured defect, a reported failure or a decision the product owner already made.
     market_question: deciding to build it depends on a market or competitive question
       yes: The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it.
       no: The decision depends only on the product's own behavior and its current users.
   Rule: indicated when any criterion is yes, not indicated when every criterion is no, undecided otherwise.
   Record:
     - item: docs/backlog/2026-11-02-teams-can-see-each-others-run-results.md
     - judged-by: agent (<reason>)
   stderr:
   exit: 0
   ```

### Invariants

1. A request's state is built only from one `followUpItem` that
   `ideaStageItem` found in the inventory by its reference; no flag, file or
   standard input can add text to it.
2. Nothing is sent for an item whose type is in `not_judged_types`, for a
   Waiting Item, for an item that matches a secret pattern, inside a Run,
   without a key, at the ceiling or with an unreadable Judge Log.
3. An answer from a model outside the pin never becomes a verdict.
4. The command exits `0` for every report and changes no other command's
   output or exit code; no command reads a verdict or a record.
5. The output never shows a verdict without its three probabilities, and
   never prints `judged-by: judge` unless the judge answered.
6. The three route clause ids, their force and the absence of `replaces` do
   not change; the first sentence of the first clause, the route and the
   TechSpec exception of the second, and the Surface Transcript sentence of
   the third stay word for word.
7. The questions, the summaries, the thresholds and `not_judged_types` live
   only in `questions.json`; the output and the JSON read them from there.

### Declared breaks

Measured with a prototype of the Baseline, skill and glossary changes on a
scratch worktree at `251ed81e`. The Go change was not prototyped.

- Behavior, task_01: `roundfix spec judge` accepts `--stage idea`, `--item`
  and `--messages`; the unknown-stage message names `idea`; the help of `spec
  judge`, `spec` and the top-level help gain lines. The Behavior Surfaces
  `command spec judge` and the `spec` command guide change, so the Behavior
  Surface Record is re-recorded. The Roundfix skill version rises.
- Behavior, task_02: every adopter's next managed refresh rewrites the three
  route clauses in `docs/agents/spec-routing.md`; this repository's guide and
  setup manifest change through the same refresh. The derived files that
  changed on the prototype are the module version record, the Standard
  TypeScript profile and its formatter golden of the Spec routing guide,
  `internal/baseline/testdata/catalog.digest`, `catalog.normalized.json` and
  `catalog.diagnostics.golden.json`, and the four goldens under
  `internal/baseline/testdata/plan-characterization/`.
- Behavior, task_03: the four skills carry the rule and the record, and their
  versions rise; the skill edits did not move the Behavior Surface Record on
  the prototype.
- Tests changed: in `internal/cli/spec_judge_test.go`,
  `TestSpecJudgeRefusesAnUnknownStageOrFormat` pins the unknown-stage message
  and moves to the new one (task_01). Every other existing test of
  `./internal/baseline` and `./skills`, the `Baseline` tests of
  `./internal/cli` and the glossary contracts passed unchanged on the
  prototype. No recorded test is removed or renamed, so
  `docs/references/coverage-record.json` is not re-recorded.
  `internal/cli/cli_test.go` is a Governed Path outside the grant and needs no
  edit: its help assertions are per line and every pinned line stays.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Core Feature 2, Success Metric 1 →
  API Contracts 1 to 3, the `idea-stage` entry, `JudgeIdeaStage`,
  `IdeaStageVerdict`, `renderIdeaStageText` (task_01).
- Goal 6, Core Feature 3, Core Feature 4, Success Metric 2, Success Metric 3,
  Success Metric 4 → API Contracts 1 and 2, Invariants 1 to 5,
  `ideaStageItem`, Surface Transcripts 1 to 4 (task_01).
- Goal 2, Goal 3, User Story 2, User Story 3, Core Feature 5, Core Feature 6,
  Success Metric 5 → API Contract 5, Invariant 6, the managed refresh
  (task_02).
- Goal 4, User Story 4, Core Feature 7, Success Metric 6 → API Contracts 3, 6
  and 7 and the record of Data Models (task_01, task_03).
- Goal 5, User Story 5 → API Contract 1, step 6 (the state is sent whatever
  its language), and the measurement reference.
- User Story 6, Core Feature 8, Core Feature 9, Success Metric 6 → API
  Contracts 5 and 6 (task_02, task_03).
- Core Feature 10, Success Metric 7 → API Contracts 4, 6 and 8, the glossary
  and documentation contracts (task_01, task_03).
- Success Metric 8 → the measurement reference, and the Judge Log row of API
  Contract 1.

Deviations from the PRD: None.

## Absence Evidence

1. The judgment and the stage — command:
   `grep -rn 'idea-stage\|IdeaStage\|"idea"' internal/judge internal/cli`;
   result: no match.
2. The rule and the record section — command:
   `grep -rn 'Idea stage is indicated\|Idea-stage judgment' internal/baseline/assets/modules docs/agents .agents/skills skills docs/user-guide CONTEXT.md`;
   result: no match.
3. The term — command: `grep -n 'Idea Stage Judgment' CONTEXT.md`; result: no
   match.
4. The new tests — command:
   `grep -rn 'TestIdeaStage\|TestSpecJudgeIdeaStage\|TestTheRouteClausesJudgeAFeature\|TestTheAuthoringSkillsJudgeAFeature\|TestTheFollowUpModeRunsTheIdeaStage\|TestTheSpecReferenceDocumentsTheIdeaStage' internal skills`;
   result: no match.
5. The report schema — command:
   `grep -rn 'roundfix/idea-stage-judgment' internal docs/user-guide .agents/skills`;
   result: no match.
6. The decision records — command: `grep -rln 'ADR-0403\|ADR-0404' docs/adr CONTEXT.md`;
   result: only the two new ADR files.

## ADR Alternatives

- ADR-0403: the decision this Spec records; it rejects a command that refuses
  a PRD without an Idea, asking the user every time, the rule alone without a
  judge, a structural check of the record and a record in the Archive Record
  or the front matter, and the design takes none of them.
- ADR-0404: the decision this Spec records; it rejects sending a request
  typed in the session, a language gate, thresholds tuned on the measured
  set, the question about worth, a second pinned model and judging Findings
  and fixes, and the design takes none of them.
- ADR-0200: applies unchanged; it adopted an advisory command over a gate and
  one pinned model, and the design keeps both: the stage exits zero whenever
  it ran and counts answers only from the pin.
- ADR-0201: applies unchanged; it kept the boundary by construction instead
  of by convention, and the design does the same through the inventory
  lookup.
- ADR-0209: not applicable; it decides which Findings and Backlog Entries the
  grouping question reads through its own reader, and this judgment reads its
  item through the follow-up inventory and sends no Finding.
- ADR-0400: applies; it rejected keeping a count, a size check in the strict
  check and a renamed term, and the design takes none of them: the Build Order
  is one context in three implementation Tasks that each earn their
  verification cycle.
- ADR-0319: applies unchanged; it rejected one ranking question, ordering by
  type alone, the English gate and letting the command call ai-memory, and
  the design takes none of them: the questions are separate Nouls, items are
  sent in any language and messages arrive only as a file.
- ADR-0231: applies unchanged; the ceiling is read as it decides.
- ADR-0239: applies unchanged; the judge's own key is read first.
- ADR-0283: applies unchanged; the Run Marker keeps this judgment silent too.
- ADR-0353: applies unchanged; it rejected judging Waiting Items, and the
  design does not send one.
- ADR-0208: not applicable; it decides how sources are grouped into Specs,
  and this Spec changes no grouping clause, so the design takes none of its
  rejected alternatives.
- ADR-0235: not applicable; it decides which models implementation work may
  run on, and this Spec starts no Agent Session and selects no model for one.
- ADR-0273: not applicable; it decides what the children of a Run may hold
  and what a Run pushes, and this Spec adds no key variable, no Run child and
  no evidence path.
- ADR-0383: not applicable; it decides when a Cross-project Message expects a
  reply and the `blocked_work` exception of the follow-up order, and this
  Spec changes neither.
- ADR-0390: not applicable; it decides the precedence of one blocking message
  for each sender and the answer to direct questions, and this Spec changes
  neither.
- ADR-0396: not applicable; it decides the completion notice a message may
  ask for, and this Spec sends no message.

## Integration Points

The judge's recipients, OpenRouter and TypeSafe, are reached only by the
released command with a key in the environment; no test, Verification or QA
row reaches them. ai-memory is reached only by a session, which hands messages
to the command as a file. The Follow-up Command's inventory is called as a
function and its output is unchanged.

## Testing Approach

Every Go test uses temporary repositories and homes, the fake judge transport
the existing judge tests use (`environment.dependencies.judgeTransport` in the
CLI, `Transport` in the package) and a fixed clock; none reaches the network,
ai-memory or the real `~/.roundfix`, and none reads a file under the Spec
root. Fake key values are built at run time. Each new top-level test calls
`t.Parallel()` first.

- `internal/judge/idea_stage_test.go` (new):
  - `TestIdeaStageAsksThreeCriteriaInOneRequest`: one request whose body
    holds the transport's request model, the three question ids with the
    measured instructions and criteria, and a state with exactly the keys
    `project` and `item`, the item with exactly `kind`, `type`, `origin`,
    `title` and `text`.
  - `TestIdeaStageVerdictUsesTheThresholds`: 0.70 is `yes`, 0.69 `unsure`,
    0.30 `no`, 0.31 `unsure`; one `yes` with two `no` is `indicated`; three
    `no` is `not indicated`; `unsure, no, no` is `undecided`; three `yes` is
    `indicated`.
  - `TestIdeaStageSendsNothingItMustNot`: zero requests, the State and the
    Reason of API Contract 1 for an item typed `fix`, `perf`, `refactor` and
    `finding`, a non-empty Run ID, no key, a month at the ceiling, an
    unreadable Judge Log and a secret-shaped text built at run time.
  - `TestIdeaStageFailsOpen`: a model outside the pin, an unreadable answer,
    HTTP 500 and HTTP 401 each give `judge_unavailable` with its reason and
    no verdict, and a Judge Log row.
  - `TestIdeaStageCutsTextOnARuneBoundary` and
    `TestIdeaStageLogsItsJudgment` (the Judge Log row of Data Models).
  - `TestTheIdeaStageQuestionsAreTheMeasuredOnes`: the three instructions,
    their criteria, the summaries, both thresholds and `not_judged_types`
    equal literals in the test, so a rewording fails until it is measured.
- `internal/cli/spec_judge_idea_test.go` (new), each Surface Transcript in
  one named test that asserts every line of both streams and the exit code:
  `TestSpecJudgeIdeaStageRequiresAnItem` (1),
  `TestSpecJudgeIdeaStageRefusesAnUnknownReference` (2),
  `TestSpecJudgeIdeaStageDoesNotJudgeAFix` (3) and
  `TestSpecJudgeIdeaStageListsTheCriteriaWhenTheJudgeDoesNotAnswer` (4,
  without a key, with the reason of the fixture). Also:
  - `TestSpecJudgeIdeaStagePrintsTheVerdictAndTheRecord`: with a fake judge
    answering 0.98, 0.84 and 0.17, the whole text of API Contract 3 and the
    whole JSON document; with 0.12, 0.17 and 0.07 `not indicated`; with 0.42,
    0.20 and 0.05 `undecided`.
  - `TestIdeaStageItemComesOnlyFromTheInventory`: `ideaStageItem` returns the
    item of an exact reference, refuses a path that exists on disk and is not
    in the inventory, a reference that differs by a `./` prefix, and an
    active Spec, each with the error of API Contract 2.
  - `TestSpecJudgeIdeaStageReadsAMessageEntry` (a `message:<id>` reference
    with `--messages`, and the usage errors of `--item` and `--messages`
    without `--stage idea` and of a slug with it) and
    `TestSpecJudgeIdeaStageDoesNotJudgeAWaitingItem`.
  - The existing tests of `spec_judge_test.go` pass with one edit, the
    unknown-stage message.
- `internal/baseline/idea_stage_route_clauses_test.go` (new), built like
  `authoring_evidence_clauses_test.go`:
  `TestTheRouteClausesJudgeAFeatureForTheIdeaStage` asserts, once each in the
  module clause and once each in the Standard TypeScript golden and the
  repository guide, whitespace-normalized, every sentence of API Contract 5
  as its own phrase, the kept phrases of Invariant 6 included, and that the
  three clauses are `mandatory` and carry no `replaces`.
- `skills/idea_stage_judgment_test.go` (new):
  `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage` (both copies of
  `brainstorming`, `write-idea`, `write-prd` and the two templates: the three
  criteria, the three verdict rules, "Never present your own answers as the
  judge's", the dismissal option, the waiting rule and every key of the
  record), `TestTheFollowUpModeRunsTheIdeaStageJudgment` (the four sentences
  of the `modes` reference, both copies) and
  `TestTheSpecReferenceDocumentsTheIdeaStage` (both copies of the `spec`
  reference: the synopsis, what is sent, the three ids, both thresholds, the
  not judged cases and the schema name).

Sabotages for the Task author to confirm on a `git archive` copy, each
observed by the named test through the function it mutates:

1. Skip the secret scan → `TestIdeaStageSendsNothingItMustNot` (the
   secret-shaped item is sent).
2. Empty `not_judged_types` → the same test (the `fix` item is sent) and
   `TestSpecJudgeIdeaStageDoesNotJudgeAFix`.
3. Compare with `>` in place of `>=` at `yes_at_least` →
   `TestIdeaStageVerdictUsesTheThresholds` (0.70 reads `unsure`).
4. Read `indicated` only when every answer is `yes` → the same test (one
   `yes` with two `no`).
5. Ignore the Run ID → `TestIdeaStageSendsNothingItMustNot`.
6. Accept an answer from a model outside the pin → `TestIdeaStageFailsOpen`.
7. Let `ideaStageItem` read a path that is not in the inventory →
   `TestIdeaStageItemComesOnlyFromTheInventory`.
8. Print `judged-by: judge` when the judge did not answer →
   `TestSpecJudgeIdeaStageListsTheCriteriaWhenTheJudgeDoesNotAnswer`.
9. Remove "The judgment informs and never refuses" from the first clause →
   `TestTheRouteClausesJudgeAFeatureForTheIdeaStage`.
10. Remove "Never present your own answers as the judge's" from `write-prd` →
    `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage`.

## Build Order

1. task_01: the judgment and the command. `internal/judge/questions.json`,
   `questions.go`, `client.go`, the new `idea_stage.go` and
   `idea_stage_test.go`; `internal/cli/spec_judge.go`, `spec_check.go` and
   `cli.go` (the three help texts), the new `spec_judge_idea.go` and
   `spec_judge_idea_test.go`, and the message pin in `spec_judge_test.go`;
   `docs/user-guide/commands/spec.md`; both copies of the Roundfix skill's
   `spec` reference with `TestTheSpecReferenceDocumentsTheIdeaStage` in the
   new `skills/idea_stage_judgment_test.go`; the Roundfix skill's version
   fields and `skills/testdata/owned-skill-versions.json`; and
   `docs/references/behavior-surfaces.json`, re-recorded as its last edit.
2. task_02: the three route clauses, `idea_stage_route_clauses_test.go`, the
   module version record, the derived Baseline files of Declared breaks and
   this repository's managed refresh (depends on: 1, because the first clause
   names the command, and to keep one writer of shared records at a time).
3. task_03: the `brainstorming`, `write-idea` and `write-prd` skills with
   their two templates and mirrors, both copies of the Roundfix skill's
   `modes` reference, the two remaining tests of
   `skills/idea_stage_judgment_test.go`, the four skills' version fields and
   the version record, and `CONTEXT.md` (depends on: 1, 2; the skills restate
   the rule task_02 fixes and both task_01 and task_03 raise the Roundfix
   skill version).
4. task_04: the QA gate (depends on: 3).

Three implementation Tasks and the QA gate, one context, each Task with its
own verification cycle: the command with its guide, the Baseline with its
derived files, and the skills with the glossary. They run in series: task_01
and task_03 share the Roundfix skill's version fields, the version record and
`skills/idea_stage_judgment_test.go`.

No other Spec is required. Specs 0391 and 0392 deliver first; the shared files
are named in `_prd.md` → Prerequisites, and each Task regenerates derived
files and records from the tree it runs on. When the skills are edited
together with the module, `make skills-sync` runs before `make
baseline-digests`, which refuses a canonical skill that differs from its
mirror.

## Risks & Considerations

- The Go surface was designed and not prototyped: the transcripts, the JSON
  document and the help lines are prescriptions, and only the reasons and the
  error suffix were copied from the current binary. The Task author checks
  each help and message pin on the default branch when the Task runs.
- `followUpInventory` reads every source to find one item. It is the same
  read the Follow-up Command does and it keeps the boundary in one place; a
  repository with a very large backlog pays that read once per judgment.
- An item is found by its exact reference. A session that passes a path with
  another spelling gets the unknown-item error, which names the Follow-up
  Command as the source of references.
- A message entry is judged from the same file the follow-up mode already
  writes; outside that mode a session must write the file to have a message
  judged, or answer the criteria itself.
- A verdict near a cut point can change between two runs; the record keeps
  the probabilities, and a re-run is not asked for.
- The clause text is long. It renders once, in the Spec routing guide of every
  adopter.
- Adopters on an older Roundfix keep the old clauses until their next
  Baseline update, and a newer guide with an older binary names a stage that
  binary refuses with a usage error; the rule then falls to the agent's own
  answers, which is the path without a judge.

## Glossary

- adds: **Idea Stage Judgment**
- changes: **Advisory Judgment**

## Decisions

Autonomous route: no person present; decisions taken from `_prd.md` → Decisions, the operator's cycle BF brief for Spec 0393, the measurement and the characterization above.

- The judgment is a stage of `roundfix spec judge` with the item named by
  `--item`: authors already run that command for advisory judgments, and a
  flag keeps the positional argument a Spec slug.
- The item is resolved through `followUpInventory` and nothing else, so the
  data boundary holds by construction. See ADR-0404.
- The type guard and the Waiting Item guard run before any check of the
  judge, so a fix reads the same inside and outside a Run.
- The criteria, their summaries and their yes and no meanings are printed
  from `questions.json` when the judge does not answer, so the agent answers
  the judge's own criteria.
- The record is a fielded section of the PRD and of the Idea; no Go code
  reads it in this Spec. See ADR-0403.
- The measurement reference is committed under `docs/references/`, because
  the thresholds and the later calibration cite it after this Spec's folder
  is archived.
- The rule lives in the first route clause, which keeps its id, so no Source
  Baseline row changes.
