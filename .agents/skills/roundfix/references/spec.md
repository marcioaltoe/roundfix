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
roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]
```

The Advisory Judge reads one active Spec from the configured Spec Root. With
no stage it judges both artifacts. It never gates and never changes another
command's exit code: advisory, skipped, and stopped runs exit `0`.

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

Set `ROUNDFIX_OPENROUTER_API_KEY` for OpenRouter first, or
`ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly when the first key is
absent. The generic `OPENROUTER_API_KEY` is not read. Every request is
recorded in `<home>/.roundfix/judge/<YYYY-MM>.jsonl` (UTC month); the monthly
ceiling is US$5.00 across both transports. Missing keys, non-English
artifacts, source or answer skips, and stopped service or log failures remain
advisory information. JSON includes clear judgments and the chosen transport.

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
