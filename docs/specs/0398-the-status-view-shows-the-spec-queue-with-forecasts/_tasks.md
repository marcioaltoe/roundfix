---
schema: spec-tasks/v1
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
requires: [0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session]
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The status view shows the Spec queue with a forecast start and end

| id      | title                                                                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Queue Forecast is computed from the delivered Specs and replayed over the recorded history      | backend | high       | —                |
| task_02 | The delivery history, the phase entries and the queue rows are read for the status view             | backend | high       | task_01          |
| task_03 | The status view shows the Spec Queue Table in its text and JSON reports                             | backend | high       | task_01, task_02 |
| task_04 | The status Skill Mode relays the table, the glossary defines its terms and two records are amended  | docs    | medium     | task_03          |
| task_05 | Run the final QA gate                                                                              | qa      | high       | task_04          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series, in the order of `_techspec.md` → Build Order. Tasks
02 and 03 share the Follow-up guide, the Skill Coverage Map, the Behavior
Surface Record and the two reader files of the command; tasks 03 and 04
share the version fields of the Roundfix skill, the owned-skill version
record and the Behavior Surface Record. Each implementation Task runs the
suites that consume its changed contracts in its own Verification, so no
separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, the Spec queue
the status view shows, in four implementation Tasks and the QA gate, and
each Task earns its verification cycle: a pure package with its replay,
three readers with their fixtures, a command surface with its contracts, and
a body of pinned text with two amended records.

## Required Spec and shared files

Spec 0397 was delivered first, and the manifest requires it. Task 02 adds a
file to the package `internal/phaserecord` that Spec 0397 created and reads
the entries it writes; task 04 writes glossary text that names the Phase
Record, which Spec 0397 defines. Every Task re-reads the tree it starts
from and keeps its edits inside its own sentences and fields. The files this
Spec shares with the Tasks of Spec 0397:

- `CONTEXT.md` (0397 tasks 01, 03 and 04; here task 04);
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, in
  their version fields only (0397 tasks 01 to 04; here tasks 03 and 04);
- `skills/testdata/owned-skill-versions.json` (0397 tasks 01 to 04; here
  tasks 03 and 04);
- `.agents/skills/roundfix/references/modes.md` and its mirror (0397 task
  04, outside its `--status` section; here task 04, inside it);
- `docs/references/behavior-surfaces.json` (0397 tasks 01 to 03; here tasks
  02 to 04) and `docs/references/skill-coverage.json` (0397 tasks 01 and 02;
  here tasks 02 and 03);
- `internal/cli/cli.go` (0397 tasks 01 to 03; here task 02, one field of the
  command dependencies);
- the package directory `internal/phaserecord` (0397 tasks 01 and 02 create
  `record.go` and `prompt.go`; here task 02 creates `read.go`);
- ADR-0407, which the plan of Spec 0397 added and none of its Tasks changes;
  here task 04 adds the notice.

No Task writes a version number: the record commands raise it on the tree
they run on.

## What the Tasks stage settled beyond the TechSpec

Each was measured or derived on the planning branch.

- The replay was repeated in whole seconds through the shape of
  `_techspec.md` → Interfaces (`NewBasis` over the history a forecast could
  see, then `Forecast` over rows) and reproduces every count of Success
  Metric 2. Changing the percentile position, the minimum of 11 or the
  pairwise sum changes at least one count; removing the cut at 512 changes
  none, so task 01 pins the cut in a test of its own.
- The fixture is pinned by the SHA-256 of its 109 lines, computed from
  section 5 of the backtest reference, so no test reads the Spec folder.
- The first example of `_prd.md` → User Experience is reproduced to the
  minute by the method over the fixture with Spec 0397 started at
  2026-11-16T11:30:44Z and read at 2026-11-16T13:36:00Z. Task 01 asserts
  its instants and task 03 its printed lines.
- `queueforecast.Row` holds no kind, so a Backlog Entry is no input of
  `Forecast`: the command passes one `Row` for each row that is a Spec and
  gives a Backlog Entry's row its two empty cells.
- An end behind a Spec in progress needs 11 times from start to merge, as
  every other end does; API Contract 6 names only the start gaps there.
- A row's `title` is the Backlog Entry's title and is empty for a Spec:
  `_prd.md` → Core Feature 12 names no title, and a Spec's slug already
  reads as one.
- The JSON form of the basis: a time or a pair of seconds that does not
  exist is `null`, as in a cell.
- A Task that changes a source file of a command names its guide
  (`SC-CLI-UNDOCUMENTED`), and task 02 changes what the command reads and
  nothing it prints. Task 02 therefore writes, in the Follow-up guide, the
  definitions of the Spec queue, of a Spec being authored and of the
  delivery history, and task 03 adds the table to the same section.
- The readers receive their Git runner in their input and not through the
  command dependencies, which gain only the clock, so a test counts or
  holds the Git calls of a reader without a seam in production code.
- `deliveryPlanIntent`, which the follow-up list reads Backlog Entries
  through, also reads every Finding's status and walks the repository
  inbox. The status view reads Backlog Entries only, so an unreadable
  Finding or inbox note never makes its `backlog` source unavailable.
- The tests of the table go in three new test files that the Build Order
  does not name: `internal/cli/follow_up_status_queue_test.go`,
  `skills/follow_up_queue_reference_test.go` and
  `skills/status_mode_queue_test.go`.

## Time in tests

A fixed clock from which a real deadline was computed broke a test of Spec
0397 on the day this graph was written. No test of this Spec compares a
time it seeded with the real clock. Each test names one instant as its
`now`, hands it to the code under test, and derives every other time from
that instant: the `created_at` of a Run, the committer date of a fixture
commit, the time of a phase entry and the name of its day file. No test
starts a Run, a Daemon or a Delivery Owner, and none waits on a timer. The
times of the replay fixture are data, and each replay case takes its `now`
from them.

## Spec 0397 merged while this graph was written

Spec 0397 reached the default branch while this graph was written, and the
planning branch was rebased onto it. Three things followed, and each is
part of this planning change:

- Task 02 creates `internal/phaserecord/read.go`, and the Skill Coverage Map
  now names the directory `internal/phaserecord/` as a source of the surface
  `command phase start`, whose only covering skill file is outside this
  Spec's authorization. The `phase start` command does not change, so
  `_prd.md` → Skills declares that surface unchanged, and task 02 declares
  the map and records the Coverage Review.
- The three links of `_idea.md` into the folder of Spec 0397, which its
  archive removed, became Git Locators at the commit that held the files.
- Every Verification command was probed on the tree that holds the code of
  Spec 0397, so each fails before the work because the work is absent.

The manifest keeps `requires`: the prerequisite is met by the Archive
Record of Spec 0397 on the default branch.
