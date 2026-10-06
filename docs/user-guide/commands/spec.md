### spec judge

```bash
roundfix spec judge <slug> [--stage <prd|techspec|tasks>] [--format <text|json>]
```

Raises advisory judgments about ADR attributions and goal-to-mechanism links
in one active Spec, resolved through the configured Spec Root. By default it
reads both the PRD and TechSpec. `--stage prd` judges PRD attributions;
`--stage techspec` judges TechSpec attributions and its Coverage Map against
the PRD's goals; `--stage tasks` asks for an advisory `light`, `standard`, or
`heavy` model-tier suggestion for each non-QA Task and prints
`suggested model-tier <task file>: <answer> at confidence <c>`. The suggestion
does not change any file and is never read at dispatch. These judgments read
the PRD, TechSpec, and accepted ADRs.

At every stage (`prd`, `techspec`, or both), the `source-grouping` question
also asks whether each Finding or Backlog Entry the Spec adopted belongs in
the same Spec as each open Backlog Entry (`open`) or unresolved Finding
(`pending` or `partial`). Adopted sources come from `references/_index.md`;
open sources come from this repository's `docs/backlog/` and `docs/findings/`.
For this question, only Findings and Backlog Entries are sent, with front
matter removed, Spec and ADR numbers scrubbed, and text cut to 1,500 Unicode
characters. Inbox Entries, terminal sources, and source code are excluded.

A pair with `P(same Spec)` at or above 0.3 prints
`suggested source-grouping <anchor> → <candidate>: P(same Spec) <probability>`,
with the probability to two decimal places. Answer it by adopting the open
source within the bound of four implementation Tasks plus its QA gate, or by
stating why it stays apart. A suggestion never gates. Its recall is low:
the measured question found 27% of true pairs, so no suggestion does not
mean no source fits.

Text output lists `advisory`, `suggested`, and `skipped` results followed by
`Judge: <a> advisory, <s> suggested, <c> clear, <k> skipped; …` and the cost
and transport. Clear judgments appear only in `--format json`, whose
schema is `roundfix/spec-judge/v1`. An advisory asks the author to correct the
artifact or explain why the text stands. A skipped result proves neither a
failure nor a clean result.

Keys, first set wins: `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, then the shared
`ROUNDFIX_OPENROUTER_API_KEY` (both OpenRouter), then
`ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly. Empty values count as unset.
The shared key keeps existing setups working when the judge key is absent.
The generic `OPENROUTER_API_KEY` is never read. The text summary names the
selected variable with `via <transport> on <variable>`; the JSON report and
every Judge Log line name it in `key_variable`, never by value. Without a
key, JSON has `key_variable: null` and the skip names the judge key first.
Each response's versioned model ID is recorded;
thresholds apply only to the pinned Jev 1.13 model family.

Every request appends its answer, model, transport, and cost to the Judge Log
at `<home>/.roundfix/judge/<YYYY-MM>.jsonl`, using the UTC month. The monthly
ceiling is `jev.monthly_ceiling_usd` in User Config, US$5 by default, across
both transports and every repository using that Roundfix Home. The log
contains no API key.

The whole run skips without a key, with an unreadable Judge Log, or when the
monthly ceiling is already reached. Non-English artifacts skip. Individual
judgments skip when their source does not qualify, the answering model does
not match the pin, the answer is unreadable, or another client error refuses
the request. A refused key (HTTP 401, 402, or 403), exhausted service retries,
a server failure, timeout, network error, reached ceiling, or unreadable or
unwritable log stops further requests. Judgments not asked count as skipped.

The command exits `0` whenever it ran, including advisory, suggested,
skipped, and stopped results. It never gates authoring or changes another
command's exit code. Exit `2` means invalid arguments, an unknown active Spec,
a missing PRD, or a missing TechSpec with `--stage techspec`. An unknown stage
fails with `unsupported --stage "<value>"; use prd, techspec or tasks`.


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

### Glossary Declaration

The Spec Consistency Check reads the `## Glossary` section of the PRD or
TechSpec. Each non-blank line declares a term added or changed, or explains
why a bold phrase is not a domain term:

```markdown
## Glossary

- adds: **New Term**
- changes: **Existing Term** — revise its definition
- not a term: **Example Label** — a test label
```

Use `None.` as the section's only entry when there are no declarations.
`adds` and `changes` allow an optional note after ` — `; `not a term` requires
a non-empty reason.

The glossary is each existing root `CONTEXT.md` and `GLOSSARY.md`, plus existing
files with either name reached by relative Markdown links in `CONTEXT-MAP.md`
or `GLOSSARY-MAP.md`. Definitions begin a line with `**Term**:`. Matching folds
case and collapses whitespace; bold phrases also match a defined term followed
by `s` or `es`.

| Code | Reported when |
| --- | --- |
| `SC-GLOSSARY-UNDECLARED` | A required declaration is absent, a declaration line is malformed, or a bold candidate is neither defined nor declared. Runs from the PRD stage, reading only the PRD at `--stage prd` and both artifacts after it. |
| `SC-GLOSSARY-UNPLANNED` | An added or changed term has no binding Task, or a changed term is not defined (declare it as added). Runs from the Tasks stage with a Task Graph. |
| `SC-GLOSSARY-MISSING` | Every binding Task has completed but the glossary still lacks the declared term. Runs from the Tasks stage with a Task Graph. |

A binding Task is a non-QA Task whose Context declares a glossary file under
`interface:` or `creates:` and whose `## Verification` text contains
`**Term**` exactly as declared. All three codes are errors, including without
`--strict`. Missing glossary files or a missing Task Graph record detector skips.

The bold rule reads single-line `**…**` spans outside fenced blocks and the
Glossary Declaration. A candidate has two to five words, no backtick, digit or
parenthesis, and does not end in `.`, `:`, `;`, `,`, `?` or `!`. Each word starts
with an uppercase letter, except `a`, `an`, `and`, `at`, `by`, `for`, `from`,
`in`, `of`, `on`, `or`, `per`, `the`, `to` and `with` after the first word.
Each uncovered phrase is reported once per artifact, at its first line.

The glossary horizon is the adding commit of
`.agents/skills/write-prd/references/glossary.md`. A PRD committed at or after
that commit must have the section in either artifact. An uncommitted PRD or
unreadable history is held to the rule once the guide exists. Earlier Specs
without a section record a skip; a Spec declaring the section is checked at
any age. The detector reads local files and Git history and makes no judge or
network call.
