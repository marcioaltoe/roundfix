---
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
prd: _prd.md
created: 2026-11-03
---

# The follow-up names what waits on someone else and where each item went — Technical Spec

## Executive Summary

The Follow-up Command reads every regular file under `docs/_inbox/`, folder
markers included, and orders every open Backlog Entry and unresolved Finding
in the four classes of ADR-0319, with no way for an item to say that it waits
on someone else. The design skips dot-named entries in the one inbox walk,
shared with `roundfix deliver plan`; reads two optional frontmatter fields,
`waiting_on` and `since`, into a Waiting Declaration; puts a valid one in a
fifth class, `waiting`, that the list holds back unless `--all` is given and
that never reaches the judge; reports a malformed one on its item; and
teaches the Roundfix skill's follow-up mode to pass `--all`, write the
declaration and name each listed item's destination. ADR-0353 records the
class and its default.

The trade-off accepted: the declaration is validated only by the command that
reads it, so a malformed one is found when the follow-up runs, not by
`roundfix spec check`.

## Project Constraints

- Identifier strategy: not applicable — no identifier or slug is minted; the
  new frontmatter and JSON keys are optional and additive. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the judge transport, keys and
  recipients are unchanged; fewer items reach it. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0319 keeps the order in code,
  ADR-0319: "The order is computed in code, never asked."; the `waiting`
  class is decided by the declaration, outside its four classes. ADR-0338
  makes the follow-up mode claim the handoff and leave no handed-in message
  pending, ADR-0338: "No message the mode handed in stays pending when the turn
  ends"; the report's destinations extend that accounting. ADR-0330
  keeps the command away from ai-memory, ADR-0330: "Roundfix's own code never
  talks to ai-memory". ADR-0353 is new. ADR Alternatives below relates each.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AW answers of 2026-11-03 quoted in the authorization record; the Governed
  Paths are the Roundfix skill files the Tasks change; no Makefile, workflow,
  module, lint or Baseline change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0343-the-follow-up-names-what-waits-and-where-each-item-went/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-03 at `4e61bf6d`; `_prd.md` → Validation records the runs.

1. `deliveryPlanIntent` (`internal/cli/deliver_plan.go`) globs
   `docs/backlog/*.md` and `docs/findings/*.md`, then walks `docs/_inbox/`
   with `filepath.WalkDir` and appends an `inbox` intent for every regular
   file, hidden or not, descending into every folder. Its two callers are
   `roundfix deliver plan` and `followUpInventory`.
2. `followUpInventory` (`internal/cli/follow_up_inventory.go`) turns open
   backlog intents into `backlog` items, `pending` and `partial` findings into
   `finding` items and every inbox intent into a `note` item, reading each
   through `readFollowUpMarkdown`, whose `followUpMarkdown` struct names
   `type`, `type-hint`, `status`, `created`, `created_at` and `origin` and
   ignores every other key. A read error makes the three sources
   `unavailable`.
3. `runFollowUpCommand` (`internal/cli/follow_up.go`) sends every item to
   `judge.PrioritizeFollowUp` and orders them with `followUpOrder` into
   priorities 1 to 4; `renderFollowUpText` prints `Judge: answered <n> of
   <len(items)>`, `Judge: not asked: no item` as a fixed string, and each
   source as `<name>: <state> <count>[: <detail>]`.
4. `parseFollowUpCommand` knows `--status`, `--messages` and `--format`, and
   refuses `--messages` with `--status`. `followUpUsage` and the root usage in
   `internal/cli/cli.go` print the synopsis.
5. `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes`
   (`internal/cli/live_documents_docscontract_test.go`) requires the guide and
   both copies of the follow-up reference to list every JSON field of the
   document, source, judge and item structs.
6. The Roundfix skill's `modes` reference runs `roundfix follow-up --format
   json` in step 1 and reports in step 8 the handoff's items and each consumed
   message, but no destination for a listed Finding or Backlog Entry.

## System Architecture

No package, command or file outside `internal/cli` is added. Changes by
component:

- `internal/cli/deliver_plan.go`: the inbox walk skips dot-named entries.
- `internal/cli/follow_up_inventory.go`: the Waiting Declaration reader.
- `internal/cli/follow_up.go`: `--all`, the waiting class, the judge input,
  the JSON keys and the text lines; `internal/cli/cli.go`: the root synopsis.
- The follow-up and deliver guides and references, the `modes` reference,
  `CONTEXT.md` and ADR-0353.

## Implementation Design

### Interfaces

```go
// internal/cli/follow_up.go
type followUpSource struct {
	// ...existing fields unchanged...
	Waiting int `json:"waiting,omitempty"` // items of this source with a valid Waiting Declaration
}
type followUpItem struct {
	// ...existing fields unchanged, then:
	WaitingOn    string `json:"waiting_on,omitempty"`
	Since        string `json:"since,omitempty"`
	WaitingError string `json:"waiting_error,omitempty"`
}
type followUpRequest struct {
	messages, format string
	status, all      bool
}
// followUpDocument gains the unexported all bool that renderFollowUpText reads.

// internal/cli/follow_up_inventory.go
type followUpMarkdown struct {
	// ...existing fields unchanged...
	WaitingOn yaml.Node `yaml:"waiting_on"`
	Since     yaml.Node `yaml:"since"`
}
// followUpWaiting reads a Waiting Declaration from decoded frontmatter.
func followUpWaiting(fields followUpMarkdown) (waitingOn, since, waitingError string)
```

### Data Models

- Backlog Entry and Finding frontmatter: optional `waiting_on` (one line of
  text) and `since` (`YYYY-MM-DD`). Read only for an open Backlog Entry and a
  `pending` or `partial` Finding.
- Follow-up JSON (`roundfix/follow-up/v1`, unchanged name): source key
  `waiting` and item keys `waiting_on`, `since`, `waiting_error`, each
  omitted when zero or empty. A Waiting Item has `priority` 5,
  `priority_name` `waiting`, `basis` `waiting`, and `answers`, `blocks` and
  `reason` null.

### API Contracts

1. API Contract: hidden inbox entries. In `deliveryPlanIntent`'s walk, an
   entry other than the walk root whose name starts with `.` is skipped: a
   folder with `fs.SkipDir`, a file by returning nil. Backlog and finding globs
   are unchanged. `roundfix deliver plan` and the `repository-inbox` source
   both stop listing such entries.
2. API Contract: the Waiting Declaration. `followUpWaiting` treats an absent
   key or a YAML null as absent and returns, in this order:
   1. `waiting_on` present but not a scalar, blank after trimming, or holding
      a line break → `waitingError` `waiting_on must be one line of non-empty text`;
   2. `waiting_on` valid and `since` absent → `since is missing`;
   3. `waiting_on` valid and `since` not a scalar that `time.Parse` accepts
      with layout `2006-01-02` → `since must be a YYYY-MM-DD date`;
   4. `waiting_on` absent and `since` present → `since needs waiting_on`;
   5. both valid → the trimmed `waiting_on` and the `since` text, no error;
   6. both absent → three empty strings.
   `followUpInventory` calls it for `backlog` and `finding` items only and
   copies the result into the item; a non-scalar value never fails the read.
   It sets `Waiting` on the `backlog` and `findings` sources to the number of
   their items with a non-empty `WaitingOn`; `Count` keeps counting every
   item read.
3. API Contract: held back by default. `runFollowUpCommand` splits the
   inventory into ranked items (empty `WaitingOn`) and Waiting Items. Only
   ranked items are sent to `judge.PrioritizeFollowUp` and ordered by
   `followUpOrder`. Without `--all` the Waiting Items are dropped from
   `Items`. With `--all` they are appended after the ranked items, sorted by
   `Since` and then `Ref` in byte order, with the fields of Data Models, and
   `Rank` runs from 1 over the whole list. `state` is `nothing_to_do` when
   `Items` is empty and nothing is in progress, as before.
4. API Contract: the judge state. With no ranked item the judge is
   `not_asked` with reason `no item` when `Items` is empty and `only waiting
   items` otherwise. Answered and partial states compare answers with the
   number of ranked items.
5. API Contract: `--all`. `parseFollowUpCommand` registers `--all` (bool).
   After the existing `--messages` check, `--status` with `--all` returns
   `--all applies only to the follow-up list; drop it with --status`, printed
   by the existing failure path with exit 2. `followUpUsage` gains
   `[--all]` after `[--status]` in its synopsis and the option line
   `  --all       Also list items waiting on someone else, held back by default`
   after the `--status` line; the root usage line in `cli.go` gains the same
   `[--all]`.
6. API Contract: the text lines. `renderFollowUpText` appends to an item line
   ` — waiting on <waiting_on> since <since>` for a Waiting Item and
   ` — invalid waiting declaration: <waiting_error>` for an item with an
   error, after the existing `from` and `blocks` suffixes. A source with
   `Waiting > 0` prints `: <n> waiting` after its count, followed by
   `, held back without --all` when `--all` was not given. The judge line
   prints `Judge: not asked: <reason>`, and its answered line counts ranked
   items only. Every other line keeps its bytes.
7. API Contract: the status view. `runFollowUpStatus` is unchanged; it never
   reads a Backlog Entry or Finding and prints no waiting count.
8. API Contract: the follow-up mode. The `modes` reference says: after
   `follow-up`, the word `all` or `--all` makes step 1 run the command with
   `--all`; step 2's plan leaves an item waiting only by writing its Waiting
   Declaration (`waiting_on` and `since`; a Finding also updates
   `updated_at`); and step 8's report names, for every listed `finding` and
   `backlog` item, the Spec it became, the Spec it was grouped into, or
   `waiting on <waiting_on> since <since>`, plus the held-back count the
   sources give, and a listed item without a destination is a gap in the
   report, as a handed-in message left pending is.
9. API Contract: the written rules. The follow-up guide and both copies of
   the follow-up reference state the hidden-entry rule, the declaration, the
   validation reasons, `--all`, the waiting class, the judge exclusion and the
   four new JSON keys in their field sentences; the deliver guide and both
   copies of the deliver reference state the hidden-entry rule for inbox rows;
   `CONTEXT.md` adds **Waiting Declaration** and **Waiting Item** and revises
   Follow-up Item, Follow-up Priority and Follow-up Command; ADR-0353 is
   created from the ADR-0353 section below.

### Surface Transcripts

The fixture of transcripts 2 and 3 is a committed Git repository holding
`docs/_inbox/.gitkeep`; an open Backlog Entry
`docs/backlog/2026-10-21-export-widget.md` (`type: feat`, `created:
2026-10-21`, title `Export widget`); and three pending Findings:
`docs/findings/2026-10-20-erp-totals.md` (`created_at: 2026-10-20`,
`waiting_on: ERP vendor support answer`, `since: 2026-10-25`, title `ERP
totals differ`), `docs/findings/2026-10-22-broken-import.md`
(`created_at: 2026-10-22`, title `Broken import`) and
`docs/findings/2026-10-23-late-feed.md` (`created_at: 2026-10-23`,
`waiting_on: supplier sample data`, no `since`, title `Late feed`). No
judge key is set. `<namespace>` is the origin's repository name or the
directory name, `<root>` the repository root.

1. Surface Transcript: the help names `--all`.

   ```transcript
   $ roundfix follow-up --help
   stdout:
   Usage:
     roundfix follow-up [--status] [--all] [--messages <path>] [--format <text|json>]

   List this repository's pending work in Follow-up Priority order, or show its status.
   It is read-only: it never writes a file, commits, pushes or prompts.
   Options:
     --status    Show work in progress, pending Specs, the next release and blockers
     --all       Also list items waiting on someone else, held back by default
     --messages  JSON file of ai-memory messages to list with the pending work
     --format    text (default) or json

   Exits 0 for a report, including one with nothing to do; 2 for invalid usage,
   configuration or messages file; 1 for an unreadable Run Database.
   stderr:
   exit: 0
   ```

2. Surface Transcript: the list holds the Waiting Item back and reports the malformed declaration.

   ```transcript
   $ roundfix follow-up
   stdout:
   Follow-up for <namespace> in <root>
   Items: 3
     1. health finding 2026-10-22 docs/findings/2026-10-22-broken-import.md: Broken import
     2. health finding 2026-10-23 docs/findings/2026-10-23-late-feed.md: Late feed — invalid waiting declaration: since is missing
     3. feature feat 2026-10-21 docs/backlog/2026-10-21-export-widget.md: Export widget
   In progress: 0
   Judge: unavailable: <reason>; ordered by type
   Sources:
     ai-memory-messages: missing: no --messages file
     backlog: checked 1
     findings: checked 3: 1 waiting, held back without --all
     repository-inbox: checked 0
     specs: checked 0
     runs: checked 0
     delivery-queue: checked 0
     pull-requests: <pull request source>
     git: checked 0
   stderr:
   exit: 0
   ```

3. Surface Transcript: `--all` lists the Waiting Item last.

   ```transcript
   $ roundfix follow-up --all
   stdout:
   Follow-up for <namespace> in <root>
   Items: 4
     1. health finding 2026-10-22 docs/findings/2026-10-22-broken-import.md: Broken import
     2. health finding 2026-10-23 docs/findings/2026-10-23-late-feed.md: Late feed — invalid waiting declaration: since is missing
     3. feature feat 2026-10-21 docs/backlog/2026-10-21-export-widget.md: Export widget
     4. waiting finding 2026-10-20 docs/findings/2026-10-20-erp-totals.md: ERP totals differ — waiting on ERP vendor support answer since 2026-10-25
   In progress: 0
   Judge: unavailable: <reason>; ordered by type
   Sources:
     ai-memory-messages: missing: no --messages file
     backlog: checked 1
     findings: checked 3: 1 waiting
     repository-inbox: checked 0
     specs: checked 0
     runs: checked 0
     delivery-queue: checked 0
     pull-requests: <pull request source>
     git: checked 0
   stderr:
   exit: 0
   ```

4. Surface Transcript: the status view refuses `--all`.

   ```transcript
   $ roundfix follow-up --status --all
   stdout:
   stderr:
   roundfix: follow-up failed: --all applies only to the follow-up list; drop it with --status
   Run 'roundfix follow-up --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. A repository whose Backlog Entries and Findings declare neither field gets the same follow-up text and JSON bytes as before.
2. No Waiting Item is ever sent to the judge, with or without --all.
3. Every item a source read is either listed or counted as held back on that source's line.
4. A malformed declaration never hides its item and never makes a source unavailable.
5. No dot-named entry under docs/_inbox/ is listed by the follow-up or by deliver plan.
6. The status view's bytes do not change.
```

### Declared breaks

- `TestFollowUpListHelpMatchesTranscript` (`internal/cli/follow_up_test.go`)
  and `TestFollowUpHelpTranscript` (`internal/cli/follow_up_status_test.go`)
  pin the help and the root synopsis; both move to transcript 1, names kept.
- `roundfix deliver plan` and the `repository-inbox` source stop listing
  dot-named inbox entries; no existing test pins one.

## ADR-0353

task_02 creates
`docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md`
with the lifecycle overlay (`status: accepted`, `created_at` and `updated_at`
`2026-11-03T18:00:00Z`, `deprecated_at: null`, `superseded_by: null`), the
title "A waiting item is held back from the follow-up list unless --all is
asked", and these parts:

- Context: ADR-0319 orders every Follow-up Item in four classes, ADR-0319:
  "The order is computed in code, never asked.", and ADR-0338 makes the
  follow-up mode claim the handoff and leave no handed-in message pending,
  ADR-0338: "No message the mode handed in stays pending when the turn ends";
  neither records an item that waits on something outside the session. On 2026-11-02
  Fluxus listed three Findings and a Backlog Entry that each waited on
  someone else and only free text said why; on 2026-11-03 the maintainer
  accepted the fields and asked for them to be left out by default and listed
  by a flag such as `--all`, quoted verbatim from `_prd.md` → Decisions.
- Decision: the Waiting Declaration, the Waiting Item in a fifth class
  `waiting` outside the four, held back by default and counted per source,
  listed by `--all` after the ranked items oldest `since` first, never sent to
  the judge, `--all` refused by the status view; validated by the command
  that reads it, a malformed declaration listing its item in its ordinary
  class with the reason; an item without either field keeps its bytes.
- Considered Options, each rejected or deferred with its reason: listing
  Waiting Items by default; dropping them without a count; judging them under
  `--all`; a Spec Check finding; the Baseline contracts in the same change
  (deferred to the open Backlog Entry
  `docs/backlog/2026-11-03-the-backlog-and-findings-contracts-name-the-waiting-declaration.md`);
  a status count.
- Consequences: the flag, the optional JSON keys, and the follow-up mode
  writing the declaration and naming each destination.

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → the inbox walk (API Contract 1),
  Invariant 5.
- Goal 2, Goal 3, Core Feature 2, Core Feature 3, Core Feature 4,
  Success Metric 2, Success Metric 3 → the reader, the split, the judge input
  and `--all` (API Contracts 2 to 6), transcripts 1 to 3, Invariants 2 and 3.
- Goal 4, Core Feature 5, Success Metric 4 → the validation order (API
  Contract 2), transcript 2, Invariant 4.
- Goal 6, Success Metric 5 → Invariant 1 and the characterization test.
- Success Metric 6 → API Contracts 5 and 7, transcript 4, Invariant 6.
- Success Metric 7 → transcripts 1 to 4.
- Goal 5, Core Feature 6 → the follow-up mode (API Contract 8).
- Core Feature 7, Success Metric 8 → the written rules (API Contract 9) and
  ADR-0353.

Deviations from the PRD: None.

## Absence Evidence

1. A hidden-entry skip in the inbox walk — command: `grep -n 'SkipDir\|HasPrefix(entry.Name' internal/cli/deliver_plan.go`; result: no match; the walk records every regular file.
2. A Waiting Declaration reader — command: `grep -rn 'waiting_on\|WaitingOn' internal/cli docs/user-guide .agents/skills/roundfix`; result: no match (the only `WaitingOn` in the module is the QA report's blocked-row field in `internal/speccheck`, unrelated).
3. A list flag `--all` on the follow-up — command: `grep -n '"all"' internal/cli/follow_up.go`; result: no match; the flags are `status`, `messages` and `format`.
4. A destination rule in the follow-up mode — command: `grep -n 'destination' .agents/skills/roundfix/references/modes.md`; result: no match.

## ADR Alternatives

- ADR-0319: followed; its rejected alternatives (one ranking question, type
  order as the default, the English gate, the command calling ai-memory) stay
  rejected: ranked items keep the four answers and the type fallback, and the
  `waiting` class sits outside the order instead of changing it.
- ADR-0338: extended; its rejected alternatives (consuming only in step 4,
  leaving an unhandled message pending, listing the handoff without accepting
  it, a new handoff field) stay rejected: the destination rule adds an account
  for listed items and changes no message or handoff rule.
- ADR-0330: followed; its rejected alternative of the command reading
  ai-memory stays rejected: the command reads only repository files and the
  messages file.
- ADR-0353: written by this Spec.

## Integration Points

None external. The command reads repository files; the judge is reached only
through the existing `internal/judge` path, with fewer items.

## Testing Approach

All tests are hermetic: temporary Git repositories, `t.TempDir()` homes, the
fake `gh` and judge transports of `follow_up_test.go`, no network.

- `internal/cli/deliver_plan_test.go`:
  `TestDeliverPlanIgnoresHiddenInboxEntries` — `docs/_inbox/.gitkeep`,
  `docs/_inbox/.obsidian/app.json` and `docs/_inbox/note.md` give exactly one
  inbox row, the note's.
- `internal/cli/follow_up_test.go`: `TestFollowUpIgnoresHiddenInboxEntries` —
  the same tree gives `repository-inbox` count 1 and one `note` item; the
  help test moves to transcript 1.
- `internal/cli/follow_up_status_test.go`: `TestFollowUpHelpTranscript`
  moves to transcript 1 and the new root synopsis.
- `internal/cli/follow_up_waiting_test.go` (new):
  - `TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration`: a fixture with a
    message, an open Backlog Entry, a pending Finding and a note, no
    declaration; text and JSON equal constants recorded from the starting
    tree before the change, with only the repository root substituted.
  - `TestFollowUpHoldsBackWaitingItems`: with a judge key and the fake judge,
    the Waiting Item is absent from `items`, the findings source has
    `count` 2 and `waiting` 1, the fake judge saw only the other titles, and
    an ordinary item's JSON carries none of the four new keys.
  - `TestFollowUpAllListsWaitingItemsLast`: with `--all`, two Waiting Items
    with `since` `2026-10-26` and `2026-10-25` come last in that order of
    `since` (`2026-10-25` first), with priority 5, `waiting` name and basis,
    their fields, ranks after the ranked items, and no judge request for
    either; with only Waiting Items the judge is `not_asked` with
    `only waiting items`.
  - `TestFollowUpReportsAnInvalidWaitingDeclaration`: a table of blank
    `waiting_on`, a list `waiting_on`, a two-line `waiting_on`, no `since`,
    `since: 2026-02-30`, `since: [x]` and `since` alone, each giving its
    `waiting_error` from API Contract 2 on an ordinary item, with every
    source `checked`.
  - `TestFollowUpHoldsBackWaitingItemsTranscript`,
    `TestFollowUpAllListsWaitingItemsTranscript` and
    `TestFollowUpStatusRefusesAllTranscript`: transcripts 2, 3 and 4 in
    full, `<root>` and `<namespace>` replaced by the fixture's values, the
    judge reason and the pull-request line by the fake's.
- `-tags docscontract`: `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes`
  over `./internal/cli` reads the new keys in the three field sentences.
- Existing suites: `./skills`, `-tags docscontract ./internal/docscontract`.
- Sabotages, each run on a temporary copy:
  - the dot check removed from the walk: `TestDeliverPlanIgnoresHiddenInboxEntries`
    and `TestFollowUpIgnoresHiddenInboxEntries` fail on `.gitkeep`;
  - only files skipped, folders walked: both fail on `.obsidian/app.json`;
  - Waiting Items sent to the judge: `TestFollowUpHoldsBackWaitingItems`
    fails on the fake judge's titles;
  - `--all` ignored: `TestFollowUpAllListsWaitingItemsLast` fails;
  - the `since` sort reversed: `TestFollowUpAllListsWaitingItemsLast` fails
    on the `2026-10-25` and `2026-10-26` pair;
  - `since` accepted when missing: `TestFollowUpReportsAnInvalidWaitingDeclaration`
    fails on the no-`since` case;
  - `omitempty` dropped from `waiting`: `TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration`
    fails on the JSON bytes.

## Build Order

1. Hidden inbox entries: the walk, its two tests, the follow-up and deliver
   guides and references, `CONTEXT.md` Follow-up Item, the Roundfix skill
   version and the Behavior Surface Record last (task_01).
2. The waiting class and `--all`: the reader, the split, the judge input, the
   flag, the text lines, the tests and transcripts, the follow-up guide and
   references, `CONTEXT.md`, ADR-0353, the skill version and the Behavior
   Surface Record last (task_02) (depends on: 1).
3. The follow-up mode: the `modes` reference, the `SKILL.md` mode summary,
   the skill version (task_03) (depends on: 2).
4. QA gate (task_04) (depends on: 3).

The Tasks run in series: task_01 and task_02 both edit the follow-up guide
and reference, `follow_up_test.go`, `CONTEXT.md`, the skill version record and
the Behavior Surface Record; task_03 teaches the flag and the declaration
task_02 builds and edits the skill version record.

Files Specs 0342 and 0344, authored in parallel, may also change, each kept to
a narrow edit here so they rebase cleanly: `skills/testdata/owned-skill-versions.json`
and the Roundfix `SKILL.md` pair (version fields and the Modes summary line),
`.agents/skills/roundfix/references/deliver.md` and its mirror,
`docs/user-guide/commands/deliver.md` (one sentence in the `deliver plan`
paragraph), `docs/references/behavior-surfaces.json` and `CONTEXT.md`. This
Spec does not edit the `implement`, `review`, `archive` or `spec` references,
the `write-*` or `archive-spec` skills, `.roundfixrc` schema or docs, or any
Baseline asset. Its functions are `deliveryPlanIntent`, `followUpInventory`,
`readFollowUpMarkdown`, `runFollowUpCommand`, `parseFollowUpCommand`,
`renderFollowUpText` and `followUpOrder`; a later Spec that touches one
rebases on it.

## Risks & Considerations

- A wait never expires: an item whose `since` is old stays held back until a
  person removes the declaration. The follow-up mode's report shows each
  held-back count, and `--all` shows the dates.
- A malformed declaration is found only when the follow-up runs; the Backlog
  Entry this plan adds proposes the Baseline contract text that would teach
  the fields where every agent reads the contracts.
- A repository that keeps real notes in dot-named files under `docs/_inbox/`
  stops seeing them; the docs layout defines raw notes only, and every
  report so far named `.gitkeep`.

## Glossary

- adds: **Waiting Declaration** — the `waiting_on` and `since` pair of a Backlog Entry or a Finding
- adds: **Waiting Item** — a Follow-up Item whose Waiting Declaration is valid
- changes: **Follow-up Item** — a hidden file or folder in `docs/_inbox/` is not one
- changes: **Follow-up Priority** — the fifth class `waiting`, outside the four
- changes: **Follow-up Command** — Waiting Items are held back unless `--all`

## Decisions

- The declaration is read with `yaml.Node` fields so that a non-scalar value
  is reported on its item instead of failing the three repository sources.
- A Waiting Item is placed after the ranked items by `since`, outside
  `followUpOrder`, so the four classes and their sort keep their code.
- The JSON keys are optional, so Invariant 1 holds without a schema version
  change.
- The hidden-entry rule lives in the one inbox walk, so `deliver plan`
  changes with the follow-up.
- Secondbrain: `wiki/index.md` read; the two `qmd query` runs named in
  `_prd.md` → Validation found Roundfix mirrors and the maintainer's GTD
  pages, none deciding a field. Exa: the two Getting Things Done sources
  quoted in `_prd.md` → Acceptance evidence.
- Three implementation Tasks in series and a QA gate.
