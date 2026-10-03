### spec judge

```bash
roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]
```

Raises advisory judgments about ADR attributions and goal-to-mechanism links
in one active Spec, resolved through the configured Spec Root. By default it
reads both the PRD and TechSpec. `--stage prd` judges PRD attributions;
`--stage techspec` judges TechSpec attributions and its Coverage Map against
the PRD's goals. These judgments read the PRD, TechSpec, and accepted ADRs.

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

The command reads `ROUNDFIX_OPENROUTER_API_KEY` for OpenRouter first, then
`ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly when the OpenRouter key is
absent. It does not read the generic `OPENROUTER_API_KEY`, so Roundfix's Jev
cost stays on its own key. Each response's versioned model ID is recorded;
thresholds apply only to the pinned Jev 1.13 model family.

Every request appends its answer, model, transport, and cost to the Judge Log
at `<home>/.roundfix/judge/<YYYY-MM>.jsonl`, using the UTC month. The monthly
ceiling is US$5.00 across both transports and every repository using that
Roundfix Home. The log contains no API key.

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
a missing PRD, or a missing TechSpec with `--stage techspec`.


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
