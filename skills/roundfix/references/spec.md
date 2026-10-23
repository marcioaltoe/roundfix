## Spec Consistency Check

Use the read-only Spec Consistency Check before a Run to compare written Spec
citations, declarations, and cross-references without editing artifacts or
emitting a QA verdict:

```bash
roundfix spec check [<slug> ...] [--format <text|json>] [--strict]
```

With no slug, the command checks every active Spec in the Spec Root. Findings
are `error` when the check locates both sides of a contradiction and `gap` when
it surfaces a candidate it cannot settle; `--strict` promotes gaps to errors.
The authoring-honesty contract includes these stable identifiers:

- `SC-QA-BYTES-PINNED` — a requirement of a pending QA Task names a byte
  invariant and a commit-shaped hexadecimal token together; compare with the
  Delivery Base instead. A gap that runs from the Tasks stage.

These four gaps hold every active Spec, whatever its age. Before an active
Spec's next Run, add the sections and lines each finding names; the fix of every
finding names the section or line and the skill rule that requires it.

`SC-HISTORY-RETAINED` is an error from the PRD stage once the manifest-recorded
`guide.docs-layout` region carries “A retirement deletes the file” after
whitespace collapses. Until then the detector is listed as skipped with its
reason. It reports files under the History Root's Findings, Backlog, ADR,
Review and Handoff families, terminal Findings left active, broken relative
links in Archive Records, Live Documents and the user guide, record self-links,
and Git Locators whose commit lacks the path or is unreachable from the default
branch and tags. Each diagnostic names the file and line and the deletion or
replacement Git Locator. For a broken Archive Record link to a file the record
promoted, the fix names the relative link to its `docs/references/` copy instead
(ADR-0308). Outside Git, locator resolution and validation are
listed as skipped; file checks still run. It also reports a citation that cites a removed history file without a Git Locator in any Citation Source outside docs/history, and skips a link to a file a pending Task declares in creates:. The check reads only local files and Git, including tracked and untracked, unignored Citation Sources, and sends nothing anywhere.

Claim Receipts resolve `ADR-NNNN` in `docs/adr/`, then read the newest addition
or modification reachable from `HEAD` in Git, including removed ADRs. Rollup
members resolve as unresolved Findings, existing archived files or deletions
under the active or history Findings family in Git. `SC-BACKLOG-UNMOVED` tells
a terminal Backlog Entry to be deleted, with Git retaining its text.

The Glossary Declaration is a `## Glossary` section in the PRD or TechSpec.
It lists each term the Spec adds or changes and each bolded phrase declared not
a domain term. The glossary horizon begins with the commit that adds
`.agents/skills/write-prd/references/glossary.md`; a Spec before that horizon
is checked under the historical skip, while a declaration is checked at any
age. The glossary findings are `SC-GLOSSARY-UNDECLARED` for an uncovered bold
term, `SC-GLOSSARY-UNPLANNED` for a declared term without a binding Task or a
changed term missing from the glossary, and `SC-GLOSSARY-MISSING` for a declared
term still absent after its binding Tasks complete.

The optional Skills Declaration lives under the PRD's `## Skills` section.
Each entry is `- unchanged: <surface id> — <reason>` with a known id and a
non-blank reason. A non-QA Task changing a covered Behavior Surface's source
needs a non-QA skills Task declaring a covering skill file under `interface:`
or `creates:`, or an unchanged entry plus a non-QA Task declaring the Skill
Coverage Map to record the Coverage Review. Both detectors skip a repository
without `docs/references/skill-coverage.json` and a PRD committed before the
oldest commit that added it; an uncommitted PRD is held.

- `SC-SKILLS-MALFORMED` — a Skills Declaration line has another shape, names
  an unknown surface or gives a blank reason. Runs from the PRD stage.
- `SC-SKILLS-UNTASKED` — a non-QA Task declares a covered surface's source
  without a skills Task or an unchanged entry with its map Task. Runs from
  the Tasks stage.

- `SC-SURFACE-RECORD-UNTASKED` — a pending non-QA Task declares a user
  guide page without declaring `docs/references/behavior-surfaces.json`
  under `interface:` or `creates:` and naming `-record-skill-coverage` in a
  requirement, or an authorization Limits line names the record. The guide
  Task itself must declare both; another Task cannot cover it. Runs from the
  Tasks stage as a gap, with no horizon; skips repositories without the record
  and Specs without a Task Graph.

- `SC-SKILL-VERSION-UNTASKED` — a pending non-QA Task declares a file of
  an owned skill under `.agents/skills/<name>/` without its mirror, both
  `SKILL.md` files, `skills/testdata/owned-skill-versions.json`, or a
  requirement naming `-record-skill-versions`. Mirror deletions use
  `deletes:`; other declarations use `interface:` or `creates:`. The Task
  itself declares the work; another Task cannot cover it. A gap with no
  horizon from the Tasks stage; skips repositories without the version record
  and Specs without a Task Graph.

- `SC-VERIFY-PATH-UNDECLARED` — a pending non-QA Task's Verification
  names a literal repository path the Task does not declare, no Task in its
  dependency closure creates, and the checked tree does not hold. Declarations
  include `creates:`, `interface:`, `deletes:`, `instruction:` and
  `## Recorded paths`; directory prefixes count too. Excludes grep/rg patterns,
  printf/echo operands, negated test operands, variables, globs, ellipses and
  paths under the Spec Root. An error with no horizon from the Tasks stage;
  Specs without a Task Graph skip with the manifest path. Reads files only
  and executes no command.

- `SC-VERIFY-SUBSUMED-SELECTOR` — a pending non-QA Task's Verification
  passes a pytest node id together with a collection argument of the same
  invocation that contains it (its file, a directory above it, `.`, or a
  node id that prefixes it), so pytest 9 runs that argument whole and the
  command passes before the test exists. `--keep-duplicates` exempts the
  invocation. An error with no horizon from the Tasks stage; reads files only
  and executes no command. Name the selector alone and run its containing path
  in a separate command (ADR-0315).

- `SC-DECISIONS-UNRECORDED` — a PRD's `## Decisions` section lacks both a
  dated answer and the autonomous-route line naming its decision source. A gap
  from the PRD stage that reads shape only. Its horizon is the oldest commit
  that added `.agents/skills/write-prd/references/clarification-record.md`: a
  PRD whose adding commit is at or after it, or an uncommitted PRD, is held;
  earlier PRDs and repositories without the guide skip the code. A missing
  PRD skips it with its path.

- `SC-RECEIPT-UNPROVEN`: a written Claim Receipt has an unresolved source, fewer than three words, or a quote absent from its source.
- `SC-RECEIPT-MISSING`: a held Spec attributes a claim to an accepted ADR without a receipt for that record in the same paragraph.
- `SC-TRANSCRIPT-UNDECLARED`: a held TechSpec has no Surface Transcripts declaration.
- `SC-TRANSCRIPT-MALFORMED`: a declared Surface Transcript breaks the command, stdout, stderr or exit block rules.
- `SC-TRANSCRIPT-UNGATED`: a declared Surface Transcript is named in no Requirement of the pending QA Task.

- `SC-VERIFY-WORK-INDEPENDENT` — a Task's Verification contains only
  repository-wide gates and working-tree cleanliness checks, so it cannot
  distinguish Task work from no work.
- `SC-REQUIREMENT-CONTRADICTORY` — declared `MUST` and `MUST NOT` clauses
  require and forbid the same named state.
- `SC-REHEARSAL-UNDECLARED` — a Task that rehearses or proves a gate lacks a
  complete `## Rehearsal Cases` declaration with
  `- Case: <case>; Observation: <observation>` entries.
- `SC-TOOLING-UNDECLARED` — a pending non-QA Task declares a Governed Path
  that its authorization record or a present Tooling authority row omits.
- `SC-CLI-UNDOCUMENTED` — a pending non-QA Task names a CLI surface without
  naming a guide in that Task or its transitive dependencies.
- `SC-LOOP-ORDER-DIVERGENT` — the shipped clause, repository guide, and
  Baseline module asset declare different Spec loop orders.
- `SC-ABSENCE-UNPROVEN` — a TechSpec has no Absence Evidence declaration or an
  item lacks its search command or result. Runs from the TechSpec stage.
- `SC-ADR-ALTERNATIVES-UNRELATED` — a listed ADR has no line relating the
  design to that ADR's rejected alternatives. Runs from the TechSpec stage.
- `SC-COVERAGE-UNREAD` — a Coverage Map lacks its `Deviations from the PRD:`
  line or has a deviation that names no PRD unit. Runs from the TechSpec stage.
- `SC-CONSUMERS-UNPROVEN` — a pending non-QA Task lacks a valid Consumer
  Declaration, or names a suite no Verification or valid checkpoint runs.
  Runs from the Tasks stage.
- `SC-WAVE-COLLISION` — two Tasks the current Task Graph lets run in one Wave
  touch a common repository file: a Verification operand, a declared
  `creates:`, `interface:` or `deletes:` path the tree holds, or a path the
  Task's newest settlement commit on any branch changed. The summary ends
  with the one edge that orders the pair, ``add a `needs` edge from Task
  <later> to Task <earlier>``. The QA gate's precondition, the settlement
  checks and the pre-dispatch refusal read the same working-tree Task Graph
  and settlement history, so a corrective Task given that edge passes all of
  them.

The missing-receipt gap starts at the oldest commit that added
`.agents/skills/write-techspec/references/concrete-contracts.md`. A PRD
committed before that guide, or a repository with no committed guide, lists
`SC-RECEIPT-MISSING` and, when a TechSpec is present,
`SC-TRANSCRIPT-UNDECLARED` as skipped with the horizon reason. An uncommitted PRD
or unreadable history is held once the guide exists. Every written receipt
is proved regardless of the horizon.

Exit `0` means no errors, including a non-strict gaps-only result. Exit `1`
means at least one error, and exit `2` means a usage error or unreadable Spec
Root. Text is the default output; JSON uses the `roundfix-speccheck/v1` schema.

### Advisory judge

```bash
roundfix spec judge <slug> --stage prd
roundfix spec judge <slug> --stage techspec
roundfix spec judge <slug> --stage tasks
roundfix spec judge <slug> [--stage <prd|techspec|tasks>] [--format <text|json>]
```

The Advisory Judge reads one active Spec from the configured Spec Root. With
no stage it judges both artifacts. With `--stage tasks` it reads `_tasks.md`,
asks one model-tier question for each non-QA Task, and prints
`suggested model-tier <task file>: <answer> at confidence <c>`. The suggestion
is advisory, writes no Spec file, and is never read at dispatch. It never gates
and never changes another command's exit code: advisory, skipped, and stopped
runs exit `0`.

Answer each `advisory` line by correcting the artifact or stating why the
text stands. A `skipped` result is neither a failure nor a clean result;
record the reason instead of claiming the artifact was judged clean.

At every stage, the `source-grouping` question pairs each Finding or Backlog
Entry the Spec adopted with every open Backlog Entry or unresolved Finding.
A `suggested` line at `P(same Spec)` of 0.3 or more is answered by adopting
the open source within the grouping bound or by stating why it stays apart.
The grouping bound is four implementation Tasks plus its QA gate. A
suggestion never gates. Only Findings and Backlog Entries are sent for this
question; its recall is low, so no suggestion does not prove no source fits.

An unknown stage fails with `unsupported --stage "<value>"; use prd, techspec
or tasks`.

When `ROUNDFIX_RUN_ID` is non-empty in the command's environment, the judge
selects no transport, reads no key or Judge Log, writes no Judge Log and sends
no request. The summary is
`Judge: skipped: inside Run <run-id>, named by ROUNDFIX_RUN_ID; <n> judgment(s) not asked`.
Every Agent Session and Verification command of a Run receives this marker
with the Run's ID. Roundfix also withholds the Judge Key Variables, including
`TYPESAFE_API_KEY`, from those children; the marker keeps the judge silent even
if a runtime restores a key from the user's shell startup.

Set `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` for OpenRouter first; the shared
`ROUNDFIX_OPENROUTER_API_KEY` is its fallback, followed by
`ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly. Empty values count as unset.
The generic `OPENROUTER_API_KEY` is never read. The text summary ends
`via <transport> on <variable>`; JSON and every Judge Log line carry the
selected name in `key_variable`, never its value. Without a key, JSON has
`key_variable: null` and the skip names the judge key first. Every request is
recorded in `<home>/.roundfix/judge/<YYYY-MM>.jsonl` (UTC month); the monthly
ceiling is `jev.monthly_ceiling_usd` in User Config, US$5 by default, across
both transports and every repository using that Roundfix Home. Missing keys,
non-English artifacts, source or answer skips, and stopped service or log
failures remain advisory information. JSON includes clear judgments and the
chosen transport.




## Spec close audit

Use the read-only Spec Audit Command after merge and sync to inspect one active
or archived Spec by slug:

```bash
roundfix spec audit <slug>
roundfix spec audit <slug> --format json
```

The `<slug>` argument is required. `--format text` is the default and reports
surviving branches and worktrees with their classification evidence, residue
reclaim commands, and any undelivered artifacts with the branch that holds
them. `--format json` emits one `roundfix-specaudit/v1` object with the same
result.

Every survivor has one of four kinds:

| Kind | Meaning |
| --- | --- |
| `pull-request` | The survivor backs an open Pull Request. |
| `pending` | The survivor holds unintegrated work. |
| `residue` | Git evidence says the survivor is safe to reclaim; the report includes the exact reclaim command. |
| `preserved` | The survivor must remain intact: an Active Run owns it, its state is unpushed or shared, or the audit could not classify it. |

The audit reports and never reclaims: it does not change Git state, the Run
Database, or Spec artifacts. A reclaim command in the report is an operator
action, not an action the audit performs.

Exit `0` means no residue or undelivered work was found. Exit `1` means residue
or undelivered work was found, or the audit could not run. Exit `2` means a
usage error or an unknown Spec slug.


The full Spec Consistency Check reports `SC-SPEC-PATH-PINNED` as an error for
each line in a non-Markdown file that names the checked Spec's active
directory. It searches tracked and untracked non-ignored files outside the
Spec Root, its archive root and `docs/history`, regardless of Task status.
Read a fixture or an exported constant instead of the Spec's file; a Spec
archives and may be deleted. The Daemon's Settlement Check refuses a pin
introduced by the Task. Staged PRD and TechSpec checks do not run this detector.

A Task Context entry whose path under the Spec Root is missing resolves to
the same relative path under the archive root when that file exists. The
built-in `docs/specs` root resolves through `docs/history/specs`; other roots
resolve through `<spec-root>/_archived`. Existing active paths keep resolving,
and a path missing in both places still reports `SC-REF-UNRESOLVED`.


A non-completed Task reports one `SC-VERIFY-TRUNCATED` error for each
Verification bullet whose first code span ends in a backslash or whose text
after that span has an odd number of backticks. Completed Tasks retain their
historical evidence. Repair the span before starting a Run.

With `--run-verification`, the check runs authored commands from the
working-tree Spec in a disposable checkout of `HEAD`. The shared prober first
runs `sh -n -c <command>` without executing the command. A rejected command
reports `malformed` with the shell's parser message, is never executed, and
makes the check exit `1`. A parsable command that fails remains `honest`;
a command that passes remains `vacuous`. The Daemon refuses malformed commands
before opening an Agent Session. If the parser cannot start, the existing
command probe still runs.

After `Verification tree: HEAD`, the text report names each Task Graph or Task
file listed by Git status as an uncommitted Verification source:

```text
Uncommitted Verification source: <path> (<untracked|modified>)
```

These lines disclose provenance and do not refuse execution. JSON includes
`verification.uncommitted`, an array of `{path, state}` objects, empty when
all sources are committed, and reports `malformed` plus the parser message in
`verification.commands` through the `verdict` and `cause` fields.
