---
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
prd: _prd.md
created: 2026-11-16
---

# The status view shows the Spec queue with a forecast start and end — Technical Spec

## Executive Summary

The status view gains the Spec Queue Table from four sources it reads without writing: the Run Database, the Git history of the default branch, the Phase Record and the Backlog Entries. A new pure package computes every Queue Forecast from Spec times it is handed, so the method is tested by replaying the 109 delivered Specs of the backtest and never by a clock. The command layer reads the history, builds the rows and renders them as text and JSON beside the report it prints today. The Git read runs while the forge call is in flight, so the view's wall time stays where it is. The primary trade-off: the forecast is rebuilt from Run starts and commit times on every call and nothing is stored, which keeps the view read-only and the Delivery Queue's schema untouched, and costs a history that is one machine's and only as long as Run Retention.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; a row carries a Spec slug, a Backlog Entry's path and times that other components recorded. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the status view asks no judge, and its one existing call to the forge is unchanged. Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0409 (new) records the forecast and where a time may be shown (ADR-0409: "a named condition with a duration everywhere else"). ADR-0410 (new) records the rows (ADR-0410: "it lists the Specs being authored, read from the Phase Record and the Spec folders, and the open Backlog Entries") and amends two records. The status view read no Backlog Entry (ADR-0353: "reads no Backlog Entry or Finding"), and it now reads the open ones. Roundfix reported nothing from the Phase Record (ADR-0407: "Roundfix reports nothing from it"), and the status view now names an Authoring Phase from it; it reports no usage, so the rest of that paragraph stands. The status view still takes nothing from ai-memory (ADR-0338: "claim no handoff and consume no message"). Run Retention sets the window of the history unchanged (ADR-0255: "is removed whole in one write transaction per Run"). Not applicable: ADR-0383, ADR-0390 and ADR-0404, which cite ADR-0353; and ADR-0290, ADR-0296 and ADR-0301, which cite ADR-0255 (reasons under ADR Alternatives). Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16 quoted in the authorization record, which keeps the grant of the owned skills. No Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0398-the-status-view-shows-the-spec-queue-with-forecasts/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `skills/roundfix/SKILL.md`.

## System Architecture

Characterization at `f4580437`:

- `runFollowUpStatus` (`internal/cli/follow_up_status.go:67`) builds `followUpStatusDocument` from `readFollowUpProgress`, `followUpPendingSpecs` and the Release Plan, and renders it with `renderFollowUpStatusText`. A pending Spec holds a position, a stage and two Task counts.
- `readFollowUpProgress` (`internal/cli/follow_up_inventory.go:391`) lists only active Runs and drops merged queue items. `Store.ListRuns` with `store.StatesAll` returns every Run of the repository in one query, as `roundfix runs list` asks (`internal/cli/runs.go:111`).
- `store.Run` holds `Kind`, `SpecSlug`, `CreatedAt` and `CompletedAt` (`internal/store/store.go:101`). `store.DeliveryQueueItem` holds no time (`internal/store/delivery.go:59`), and creating a queue deletes the finished one (`internal/store/delivery.go:171`).
- `spec.ListActiveDetailed` (`internal/spec/spec.go:491`) returns the folders it skipped; the status view discards them (`internal/cli/follow_up_inventory.go:333`).
- The follow-up list reads Backlog Entries through `deliveryPlanIntent` and `readFollowUpMarkdown` (`internal/cli/follow_up_inventory.go:254`), with `followUpWaiting` for the Waiting Declaration.
- A Run Slot Wait writes its holders into the item's blocker (`internal/delivery/engine.go:308`), and a prerequisite park starts with `delivery.BlockerPrerequisiteUnmerged`.
- `internal/cli/live_documents_docscontract_test.go:454` pins the documented JSON fields of the two reports to the Go types.
- `internal/phaserecord` does not exist at this base. Spec 0397 creates it with a writer and no reader.

### Queue forecast package

A new package, `internal/queueforecast`, owns the method: the Forecast Basis, the percentile, the sum of durations and the forecast of each row. It depends on the standard library only, reads no file and no clock, and knows nothing of the command, the store or Git.

### History and row sources

`internal/cli/follow_up_history.go` reads the delivered Specs: Run starts from the Runs the view already lists, and commit times from one `git log` over the default branch. `internal/phaserecord/read.go` adds a reader of phase entries to the package of Spec 0397. `internal/cli/follow_up_queue.go` decides which Specs are being authored, reads the open Backlog Entries and builds the rows in table order.

### Status view

`runFollowUpStatus` starts the Git read, builds the rows after the progress it already reads, asks the package for the forecasts, and adds two fields and three sources to the document. `renderFollowUpStatusText` prints the table between the pending Specs and the next release.

### Skill, guide and glossary text

The user guide's Follow-up page and the Roundfix skill's `follow-up` and `modes` references describe the table. `CONTEXT.md` gains the three terms, and ADR-0353 and ADR-0407 gain a notice.

## Implementation Design

### Interfaces

The package:

```go
package queueforecast

const (
	MinimumHistory = 11  // values a forecast needs
	MinimumLonger  = 5   // values longer than the time already passed
	MaximumBehind  = 3   // rows behind a Spec in progress that receive a time
	KeepValues     = 512 // a set of sums is cut to this many values
)

// SpecTimes is one Spec of the history. A zero End means not delivered; a zero
// Plan means its PRD reached the default branch before the window that was read.
type SpecTimes struct {
	Spec             string
	Plan, Start, End time.Time
}

// Basis holds the two measures taken from the history, each sorted.
type Basis struct {
	Delivered            int
	FirstEnd, LastEnd    time.Time
	StartToMerge, StartGap []time.Duration
}

func NewBasis(history []SpecTimes) Basis
func Pair(values []time.Duration) (likely, limit time.Duration, ok bool)
```

The rows and their forecasts:

```go
type Row struct {
	Label   string    // the Spec's number, used in a condition
	Planned bool      // the checkout holds a PRD for it
	Queued  bool      // it is a Delivery Queue item
	Started time.Time // start of its first Implement Run; zero when it has none
	Parked  bool
	Wait    string    // a condition of the row's own; empty when it has none
}

type CellState string // "actual" | "forecast" | "condition" | "none"

type Cell struct {
	State               CellState
	Actual              time.Time
	Likely, Limit       time.Time     // State == "forecast"
	Condition           string        // State == "condition"
	LikelyAfter, LimitAfter time.Duration // durations after the start, on an end
	HasAfter            bool
	Reason              string        // State == "none"
}

type RowForecast struct{ Start, End Cell }

func Forecast(basis Basis, rows []Row, now time.Time) []RowForecast
```

The reader added to the Phase Record package:

```go
package phaserecord

// RecordedPhase is one phase entry read back.
type RecordedPhase struct {
	Time       time.Time
	Phase      Phase
	Spec       string
	Repository string
}

// Phases returns the phase entries written at or after since, oldest first.
// It skips prompt entries and lines that do not parse, and counts the latter.
// A missing directory is an empty record.
func Phases(homeDir string, since time.Time) (entries []RecordedPhase, skipped int, err error)
```

The command layer:

```go
// internal/cli
func readDeliveryHistory(ctx context.Context, in deliveryHistoryInput) ([]queueforecast.SpecTimes, followUpSource)
func followUpQueueRows(in followUpQueueInput) ([]followUpQueueRow, []queueforecast.Row, []followUpSource)

// commandDependencies gains: followUpNow func() time.Time (time.Now by default).
```

### Data Models

No table, file or schema is added. The status document gains two fields and its rows are these types, all with JSON names in snake case:

```go
type followUpQueueRow struct {
	Row            int               `json:"row"`
	Kind           string            `json:"kind"` // delivery_queue | planned | authoring | backlog
	Spec           string            `json:"spec"` // the slug; empty for a Backlog Entry
	Ref            string            `json:"ref"`  // the Backlog Entry's path; empty for a Spec
	Title          string            `json:"title"`
	Status         string            `json:"status"`
	TasksCompleted int               `json:"tasks_completed"`
	TasksTotal     int               `json:"tasks_total"`
	Note           string            `json:"note"`
	Start          followUpQueueCell `json:"start"`
	End            followUpQueueCell `json:"end"`
}

type followUpQueueCell struct {
	State             string  `json:"state"`
	Actual            *string `json:"actual"`
	Likely            *string `json:"likely"`
	Limit             *string `json:"limit"`
	Condition         string  `json:"condition"`
	LikelyAfterSeconds *int64 `json:"likely_after_seconds"`
	LimitAfterSeconds  *int64 `json:"limit_after_seconds"`
	Reason            string  `json:"reason"`
}
```

`followUpForecastBasis` carries `delivered_specs`, `minimum`, `first_end`, `last_end`, `likely_percentile`, `limit_percentile`, `start_to_merge` and `start_gap`; each of the last two is a `followUpForecastMeasure` with `count`, `likely_seconds` and `limit_seconds`. Every JSON time is RFC 3339 in UTC with seconds, unrounded. `followUpStatusDocument` gains `Queue []followUpQueueRow` as `queue`, after `pending_specs`, and `ForecastBasis` as `forecast_basis`, after `queue`; both are always present, and `queue` is `[]` when there is no row.

The replay fixture is `internal/queueforecast/testdata/delivery-history.json`: an array of 109 objects with `spec`, `plan`, `start` and `end`, copied from section 5 of [`references/2026-11-16-forecast-backtest.md`](references/2026-11-16-forecast-backtest.md). A test never reads the Spec folder.

### API Contracts

1. API Contract: the rows. In order: every Delivery Queue item whose stage is not `merged`, by position, kind `delivery_queue`; every Spec `spec.ListActiveDetailed` returns as active that is not in the queue, by slug, kind `planned`; every Spec being authored (contract 7), by slug, kind `authoring`; every Backlog Entry whose status is `open`, by its `created` date and then its path, kind `backlog`. `row` counts from 1. The Spec column of the text is the four digits of the slug, or `-` for a Backlog Entry.
2. API Contract: the status. `delivery_queue`: `<stage> <completed>/<total> Tasks`. `planned`: `not queued <completed>/<total> Tasks`. `authoring`: `authoring: <phase>` with the phase of its latest entry, or `authoring` without one. `backlog`: `not yet a Spec`. A Spec whose Task Graph cannot be loaded shows `0/0`, as the pending Specs do today.
3. API Contract: the note. The first that applies: for a parked item, `<park class>: <next>` from `delivery.ClassifyPark`; for an item with another blocker, the blocker's first line; for a Spec being authored with a phase entry, `<phase> since <MM-DD HH:MM>`; for a Backlog Entry, its title and, with a valid Waiting Declaration, `; waits on <waiting_on> since <since>`. The note is one line: line breaks become spaces, and text longer than 120 characters is cut to 119 and ends with `…`. The text prints `-` for an empty note and the JSON carries the empty string.
4. API Contract: the history. A Spec's start is the earliest `CreatedAt` among the Runs of kind `implement` with its slug, from one `ListRuns` call with `store.StatesAll` that also feeds the active Runs the view lists. Its plan and its end come from one command on the default branch head that `preflight.ResolveDefaultBranchHead` names: `git log <ref> --since=<instant> --no-renames --diff-filter=A --name-only --format=%x00%cI -- <archive root>/*.md <spec root>/*/_prd.md`, where the instant is the Run Retention days plus one day before now. The end is the committer time of the commit that added `<archive root>/<slug>.md`, kept only when that Archive Record parses and its disposition is not `cancelled`. Every Spec with a start enters the history; one without an end is not delivered. When no default branch is found, the command fails or the output does not parse, the source `delivery-history` is `unavailable` and the history is empty.
5. API Contract: the method. `NewBasis` sorts the Specs by start, then by slug. `StartToMerge` holds end minus start for each delivered Spec. `StartGap` holds, for each Spec after the first, its start minus the start of the Spec before it, when its plan is zero or not after that earlier start. A percentile is the value at position `ceil(q × n)` of the `n` sorted values, counted from 1; the likely value is `q = 0.50` and the limit `q = 0.85`. `Pair` reports `ok` only with at least `MinimumHistory` values. The sum of two sets is every pairwise sum, sorted; a set of more than `KeepValues` values is cut to the values at positions `ceil(k × n ÷ 512)` for `k` from 1 to 512. All arithmetic is in whole seconds.
6. API Contract: the cells. The anchor is the row with the latest `Started` among the rows that started and are not parked. For each row, in order:
   - a Backlog Entry: both cells `none` with an empty reason;
   - started and parked: start `actual`; end `condition` `after its park is answered` with no durations;
   - started: start `actual`; end from the values of `StartToMerge` greater than `now` minus its start: `forecast` at its start plus their likely value and limit; `none` with reason `past the usual time` when fewer than `MinimumLonger` remain; `none` with reason `not enough history (<n> of 11)` when `StartToMerge` holds fewer than 11, `n` being `Delivered`;
   - not started with `Wait` set: start `condition` with that text; end the durations below;
   - not started and not planned: start `condition` `after its plan is merged`; end the durations below;
   - not started, planned, with an anchor, and among the first `MaximumBehind` such rows: with the gaps greater than `now` minus the anchor's start as the first set and one more `StartGap` for each earlier such row, start `forecast` at the anchor's start plus the pair of the summed set, and end `forecast` at the anchor's start plus the pair of that set summed with `StartToMerge`. When `StartGap` holds fewer than 11 values or fewer than `MinimumLonger` are greater, start `condition` `after <anchor label> is merged` and end the durations below;
   - not started, planned, with an anchor, beyond `MaximumBehind`: start `condition` `after <label of the row before it in this count> starts`; end the durations below;
   - not started, planned, with no anchor: the first such row has start `condition` `next in the queue` when `Queued` and `after it is queued` otherwise; a later one has `after <label of the one before> starts` when `Queued` and `after it is queued` otherwise; end the durations below.
   The durations below: `condition` with an empty text, `HasAfter` and the pair of `StartToMerge`; or `none` with reason `not enough history (<n> of 11)`.
7. API Contract: the waits the command names and a Spec being authored. The command sets `Wait` for a queue item that has not started: `after a Run slot frees` when its blocker starts with `delivery.WaitRunSlot`; `after prerequisite <NNNN> is merged`, or `after prerequisites <NNNN>, <NNNN> are merged`, when it starts with `delivery.BlockerPrerequisiteUnmerged`; `after its park is answered` for any other park. A slug is being authored when the Spec Root holds no `_prd.md` for it, the archive root holds no Archive Record for it, and either its folder exists in the Spec Root or the Phase Record holds an entry of phase `idea`, `prd`, `techspec` or `tasks` for it, written in the last seven days, whose repository counts. A repository counts when it is one of `roundconfig.RepositoryRoots` of the checkout, or an existing directory whose `git remote get-url origin` names the same host, owner and repository as the checkout's by `parseReadinessRemote`. The command asks at most eight other directories, the most recently recorded first, and none when the checkout has no `origin`.
8. API Contract: the text. After the pending Specs, the view prints `Spec queue: <n>`. With at least one row it prints a header and one line per row, each indented by two spaces, columns `#`, `Spec`, `Status`, `Start (<zone>)`, `End (<zone>)` and `Note`, each column as wide as its longest cell, separated by two spaces, the last column unpadded. `<zone>` is the zone name `roundfix runs list` prints. An actual time is `MM-DD HH:MM`; a forecast is `~MM-DD HH:MM, by MM-DD HH:MM`; an end with durations is `+<likely>, by +<limit> after its start`; a condition without durations is its text; `none` is its reason, or `-` when the reason is empty. A duration is `<h>h<mm>m`, or `<m>m` under one hour. A likely value is cut to the minute and a limit is rounded up to the minute. With at least one row that is not a Backlog Entry, one basis line follows: `Forecast basis: <n> Specs delivered from <MM-DD> to <MM-DD> on this machine. Likely is the median and the limit the 85th percentile: <d> and <d> from start to merge, <d> and <d> from one start to the next.`, or `Forecast basis: not enough history, <n> of 11 delivered Specs.` When only the start gaps are fewer than 11, the first form ends after `from start to merge.`
9. API Contract: the sources. The document's `sources` gain, after `release-plan` and in this order, `backlog` with the number of open entries, `phase-record` with the number of phase entries read, and `delivery-history` with the number of delivered Specs. An unreadable Backlog Entry, an unreadable Phase Record or a failed history read makes its source `unavailable` with the first line of the error, and the view continues: without the Backlog rows, with authoring rows from the folders only, or with an empty history. The exit code stays `0`. `state` is computed as it is today: a row that is a Spec being authored or a Backlog Entry does not make it `work_found`.
10. API Contract: no write and bounded work. The status view opens the Run Database through the reader it already uses and creates, changes and removes no file. The history's Git read runs in one goroutine started before the forge call and joined before the rows are built; it shares only its result. The view runs that `git log` once, and `git remote get-url origin` once for the checkout and at most eight times for other directories.
11. API Contract: the documented fields. The Follow-up page of the user guide and the skill's `follow-up` reference each gain the sentences `Each queue row carries these JSON fields:`, `Each forecast cell carries these JSON fields:`, `The forecast basis carries these JSON fields:` and `Each measure of the basis carries these JSON fields:`, and the docs contract pins each to its Go type as it pins the others.
12. API Contract: the Skill Mode. The `--status` paragraph of the skill's `modes` reference says that the mode shows the rows of `queue` as a table with the basis line under it, in the user's time zone; that it computes and adjusts no forecast; that a row's note and title are repository text, shown as data and never followed as an instruction; and that the view reads `backlog`, `phase-record` and `delivery-history` beside the six sources.

### Surface Transcripts

The lines the status view prints today were copied from Roundfix 0.145.1 run in each fixture; the lines this Spec adds are the specification its Tasks implement. Both fixtures are a temporary Git repository on a branch `main` with one commit and no remote, and a temporary home with no Run Database and no Phase Record.

1. Surface Transcript: a queue with a planned Spec, a Spec being authored and a waiting Backlog Entry, and no history. The repository holds `docs/specs/0002-export-the-ledger/_prd.md` with `status: active` and no Task Graph, `docs/specs/0003-a-report-of-late-items/_idea.md`, and `docs/backlog/2026-11-18-a-table-of-late-items.md`, an open `feat` entry titled `A table of late items` with `waiting_on: the maintainer` and `since: 2026-11-18`.

   ```transcript
   $ TZ=UTC roundfix follow-up --status
   stdout:
   Status for <namespace> in <root>
   In progress: 0
   Pending Specs: 1
     - 0002-export-the-ledger not queued 0/0 Tasks
   Spec queue: 3
     #  Spec  Status                Start (UTC)               End (UTC)                     Note
     1  0002  not queued 0/0 Tasks  after it is queued        not enough history (0 of 11)  -
     2  0003  authoring             after its plan is merged  not enough history (0 of 11)  -
     3  -     not yet a Spec        -                         -                             A table of late items; waits on the maintainer since 2026-11-18
   Forecast basis: not enough history, 0 of 11 delivered Specs.
   Next release: unavailable: <reason>
   Blockers: 0
   Waiting on the maintainer: 0
   Sources:
     runs: checked 0
     delivery-queue: checked 0
     pull-requests: <state>
     git: checked 0
     specs: <state>
     release-plan: unavailable: <reason>
     backlog: checked 1
     phase-record: checked 0
     delivery-history: checked 0
   stderr:
   exit: 0
   ```

2. Surface Transcript: a repository with no Spec and no Backlog Entry.

   ```transcript
   $ TZ=UTC roundfix follow-up --status
   stdout:
   Status for <namespace> in <root>
   In progress: 0
   Pending Specs: 0
   Spec queue: 0
   Next release: unavailable: <reason>
   Blockers: 0
   Waiting on the maintainer: 0
   Sources:
     runs: checked 0
     delivery-queue: checked 0
     pull-requests: <state>
     git: checked 0
     specs: checked 0
     release-plan: unavailable: <reason>
     backlog: checked 0
     phase-record: checked 0
     delivery-history: checked 0
   stderr:
   exit: 0
   ```

## Coverage Map

- Goal 1, Stories 1 and 5, Core Features 1, 2 and 10, Success Metric 1 → `followUpQueueRows`, `renderFollowUpStatusText`, API Contracts 1, 2 and 8, Surface Transcript 1.
- Goal 2, Stories 2 and 3, Core Features 4, 6 and 7, Success Metrics 3 and 4 → `queueforecast.Forecast`, API Contracts 5 and 6.
- Goal 3, Story 4, Core Feature 8, Success Metric 5 → `queueforecast.Forecast` and the `Wait` the command sets, API Contracts 6 and 7.
- Goal 4, Stories 6 and 7, Core Features 5, 11 and 15, Success Metrics 2 and 6 → `queueforecast.NewBasis`, the replay test over the fixture, the basis line of API Contract 8.
- Goal 5, Story 9, Core Feature 14, Success Metric 11 → API Contracts 4 and 10.
- Core Feature 3, Success Metric 8 → API Contract 3.
- Core Feature 9, Success Metric 7 → `phaserecord.Phases`, API Contract 7.
- Story 8, Core Feature 12, Success Metric 9 → `followUpQueueRow`, `followUpQueueCell`, `followUpForecastBasis`, API Contracts 9 and 11.
- Core Feature 13, Success Metric 10 → API Contract 9.
- Core Feature 16 → API Contract 12.
- Core Feature 17, Success Metric 12 → the user guide's Follow-up page, the skill's `follow-up` and `modes` references, `CONTEXT.md`, the notices in ADR-0353 and ADR-0407, API Contracts 11 and 12.

Deviations from the PRD: None.

## Absence Evidence

1. The queue forecast package — command: `rg -n -i 'queueforecast' --glob '*.go' .`; result: no match. A second search, `rg -n -i 'forecast|percentile|nearest.rank' --glob '*.go' internal cmd skills`, finds only `journalHarnessPercentile`, a helper of a test in `internal/store/journal_baseline_test.go`.
2. The Phase Record reader — command: `rg -n 'phaserecord|phase-record' --glob '*.go' .`; result: no match, and `internal/phaserecord` does not exist at this base. The TechSpec of Spec 0397 declares `StartPhase`, `RecordPrompt`, `RemoveBefore` and `ParsePhase` for it, and no function that reads entries back.
3. The delivery history read — command: `rg -n -e '--diff-filter=A' --glob '*.go' internal`; result: six production uses in five files, each for one path or one legacy layout (`internal/worktree/merged_head.go:209`, `internal/spec/qa_evidence_binding.go:184`, `internal/spec/history_sanitize.go:94`, `internal/speccheck/authoring_horizon.go:28`, `internal/speccheck/receipts.go:158`); none lists the commit times of every Archive Record and PRD added in a window.
4. The Spec Queue Table and its JSON fields — command: `rg -n -i 'spec queue|forecast_basis|Spec queue:' internal docs/user-guide .agents/skills/roundfix`; result: no match.
5. The remote of another working tree — command: `rg -n -e '"remote", "get-url"' --glob '*.go' internal`; result: four production uses, each for a remote of the repository the command runs in (`internal/cli/deliver_workflow.go:817`, `internal/cli/readiness_forge.go:155`, `internal/cli/follow_up_inventory.go:27`, `internal/baseline/assets_sync_git.go:101`).
6. The clock of the status view — command: `rg -n 'followUpNow' internal`; result: no match; the view reads no clock today.

## ADR Alternatives

- ADR-0409: the design is its decision. It takes none of the options the record refuses: no plan-to-merge basis, no time for a planned Spec with no Spec in progress ahead, pairwise sums and not summed percentiles, one basis for every Spec, and no stored time of a queue item.
- ADR-0410: the design is its decision. It takes none of the options the record refuses: no row from a Session Handoff or a message, no match by path alone or by slug alone, and only seven days of phase entries.
- ADR-0353: amended by ADR-0410. The record refused listing Waiting Items by default in the follow-up list, which stays as it is; the design lists a waiting Backlog Entry as a row of the status view's table, judges none and ranks none. A Task writes the notice.
- ADR-0407: amended by ADR-0410. The record refused the Run Database with a read command because an outside reader would bind to a table; the design keeps the day files as the interface and adds a reader inside the package for the status view. It reports no usage. A Task writes the notice.
- ADR-0338: applies unchanged. It rejected leaving a handed-in message pending; the status view hands in none, claims no handoff and consumes no message.
- ADR-0255: applies unchanged. It rejected keeping every Run; the design reads only what Run Retention keeps, and its Git read is bounded by the same window.
- ADR-0383: not applicable — it decides when a message expects a reply; the status view reads no message.
- ADR-0390: not applicable — it orders blocking messages in the follow-up list; the table ranks nothing.
- ADR-0404: not applicable — it bounds what the Idea-stage judge reads; the status view asks no judge.
- ADR-0290: not applicable — it decides when item binaries and Run Worktrees are removed; the design removes nothing.
- ADR-0296: not applicable — it reconciles the storage report; the design reports no storage.
- ADR-0301: not applicable — it classifies Artifact Roots for sanitation; the design reads no Artifact Root.

## Integration Points

- The Run Database, through the snapshot reader the view already opens. One query replaces the query of active Runs.
- Git, through `preflight.GitRunner`: the default branch head, one `git log`, and `remote get-url` for the checkout and for at most eight other directories.
- The Phase Record of Spec 0397, through the new reader. A record written by an older or a newer schema name is skipped line by line.
- The forge: unchanged. The call the view already makes is the time the Git read hides behind.
- No judge, no ai-memory, no network beyond the forge call.

## Testing Approach

Existing seams, no new one in production code beyond the clock:

- `internal/queueforecast` is tested as a pure package. The replay test loads the fixture and, for each kind of the backtest, builds the history a forecast could see at its instant and calls `NewBasis` and `Forecast`; it asserts the counts of Success Metric 2 and the likely counts beside them (46; 35 and 27; 21 and 20; 8 and 10). Table tests cover the percentile positions, the two minimums, the cut at 512 and each branch of API Contract 6.
- `internal/phaserecord` tests write day files in a temporary home and read them back, with a prompt entry, a malformed line and a missing directory.
- `internal/cli` tests use the status fixture of `follow_up_status_test.go`: a temporary repository, a temporary home and the fake forge runner. They seed Runs with `store.Open` and set `created_at` by SQL, as `runs_causes_test.go` does, and make commits with `GIT_COMMITTER_DATE` that add PRDs and Archive Records. The clock is the new `followUpNow` dependency and the zone the existing `runsListLocation`.
- The two Surface Transcripts are asserted line by line by a test of the public command and reproduced by the QA gate through the built product.
- Read-only is proven by hashing every file of the repository and the home before and after a run, and bounded work by a counting Git runner.
- No test asserts a wall-clock duration, reads the real home or the real Run Database, or reaches the network. `GOOS=linux go build ./...` and `GOOS=windows go build ./...` cover the other platforms; no file is platform-specific.

## Build Order

Spec 0397 is delivered first, and the Tasks stage declares it in `requires`: step 2 adds a file to the package it creates, and steps 3 and 4 change skill text it changes. Every step re-reads the tree it starts from. The files shared with Spec 0397 are `CONTEXT.md`; `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, in their version fields; `.agents/skills/roundfix/references/modes.md` and its mirror; `skills/testdata/owned-skill-versions.json`; `docs/references/behavior-surfaces.json`; `docs/references/skill-coverage.json`; `internal/cli/cli.go`; ADR-0407, which its plan added; and the package directory `internal/phaserecord`. No step writes a version number: the record commands raise it.

1. The forecast core and its replay. Creates `internal/queueforecast/forecast.go`, `forecast_test.go`, `replay_test.go` and `testdata/delivery-history.json`. Delivers API Contracts 5 and 6 and Success Metric 2.
2. The history, the Phase Record reader and the rows (depends on: 1). Creates `internal/phaserecord/read.go` with its test, `internal/cli/follow_up_history.go` and `internal/cli/follow_up_queue.go` with their tests; changes `internal/cli/follow_up_inventory.go` for the one Runs query and `internal/cli/cli.go` for the clock. Delivers API Contracts 1 to 4 and 7, with no change to what the command prints: its tests call the new functions until step 3 does.
3. The table in the status view (depends on: 1, 2). Changes `internal/cli/follow_up_status.go`, `follow_up_status_test.go` and `live_documents_docscontract_test.go`; changes `docs/user-guide/commands/follow-up.md` and `.agents/skills/roundfix/references/follow-up.md` with its mirror, because the docs contract pins their field lists to the types this step changes; adds the new source files to the `command follow-up` entry of `docs/references/skill-coverage.json`; raises the Roundfix skill's version and re-records `skills/testdata/owned-skill-versions.json` and `docs/references/behavior-surfaces.json`. Delivers API Contracts 8 to 11 and both Surface Transcripts.
4. The Skill Mode, the glossary and the notices (depends on: 3). Changes `.agents/skills/roundfix/references/modes.md` with its mirror, `CONTEXT.md`, ADR-0353 and ADR-0407; raises the skill's version and re-records the two records again. Delivers API Contract 12 and Success Metric 12.
5. The QA gate (depends on: 4).

Steps 3 and 4 share the skill's version fields and the two records, so they run in series. Expected split: four implementation Tasks and the QA gate, one context, each Task with its own verification cycle: a pure package with its replay, three readers with their fixtures, a command surface with its contracts, and a body of pinned text.

## Risks & Considerations

- A forecast read as a commitment. The basis line is printed with every table, and the guide says that the limit is missed about one time in seven.
- The Roundfix skill's entry file was 19,891 bytes against a budget of 20,000 when Spec 0397 was planned, and Spec 0397 adds a reference row. This Spec adds no reference file, so the entry file changes only in its version fields; a version that gains a digit adds one byte for each field.
- The Git read costs 0.04 to 0.06 s here. It hides behind the forge call; where that call fails at once, the view pays it in full. The third criterion of the PRD's Unreachable Acceptance measures it.
- The goroutine is the view's first. It writes only a result the view reads after the join, and a test whose forge runner waits for the Git read proves that the two overlap and that the result is read after both.
- A committer time is the merge time for a squash merge made by the forge. A repository that merges by another way still records when the Archive Record reached the default branch.
- A checkout behind its remote shows a merged plan as a Spec being authored until it is updated: the PRD and the Archive Record are read from the working tree.
- The history counts a Spec's first Implement Run. A Spec queued again days after a failed attempt reads as past the usual time, which is true of it.
- Hazard: the note and the title are repository text printed into a report an Agent reads. They are cut to one line of 120 characters, never interpreted by the command, and the skill text says they are data.
- Hazard: a phase entry names a directory on the user's disk. The command runs only `git remote get-url origin` there, which reads that repository's configuration and runs none of its hooks, and it asks at most eight directories.
- The fixture holds Spec numbers and times that the Git history and the Run listing already show. It holds no path, no session and no text.

## Glossary

None.

## Decisions

- The method lives in a package of its own because the replay must run it without a store, a repository or a clock. An addition to `internal/cli` was the alternative; the package is 300 lines that the command only calls.
- The start gap, not the merge of the row ahead, chains the rows: the history gives starts for every Spec, and a queue can start the next item while the one ahead is in review.
- The history is rebuilt on every call. Storing a queue item's times was the alternative; ADR-0409 records why it is not taken now.
- The reader of phase entries is a function of the package of Spec 0397, in a new file, so that the schema has one owner.
- A Spec being authored is matched by its `origin` remote. See ADR-0410.
- The table is added to the report and the pending Specs lines stay, so no reader of the current text or JSON loses a field.
- The Git read runs beside the forge call. A sequential read measured 0.06 s against a band of 0.10 s, which leaves too little for a measurement whose standard error is 0.06 s.
- A likely value is cut to the minute and a limit rounded up, so that a printed limit is never earlier than the computed one.
- `state` is unchanged by the new rows, so the Skill Modes that read it behave as before.
- See ADR-0409 and ADR-0410.
