## Pre-PR review

Run the configured reviewer over the current candidate with:

```bash
roundfix review [--base <ref>]
roundfix review dispose <finding-id> --dismiss --evidence <text>
roundfix review dispose <finding-id> --fixed-by <commit>
```

`--base` selects the base ref; when omitted, Roundfix uses the repository's
default branch. The workflow computes the candidate diff from the merge base
of the current head and the selected base, then hands that diff to the
configured reviewer in a read-only session. The review record names the
repository, that merge base as `baseCommit`, the resolved base tip as
`baseTipCommit`, the head commit, effective provider, and policy source, so it
is bound to the candidate that was examined. A head with no shared history
with the selected base exits `2` before any reviewer call or readiness probe.

The policy values are `codex`, `claude`, `coderabbit`, and `none`. Explicit
`none` performs no reviewer call and no readiness probe, records a configured
omission for the candidate, and exits `0`. `claude` runs through the same
read-only review path as `codex`. `coderabbit` remains refused because no
supported local CodeRabbit review surface is installed or specified, and exits
`2`.

Roundfix classifies the reviewer's answer by substance, not exact formatting.
It recognizes a `No findings` verdict after case folding, surrounding Markdown
emphasis, and trailing punctuation are normalized, or a `Findings:` verdict
with its findings text. A pass must be the whole answer: after normalization,
the answer is `No findings`, or its only content line is `Findings:` with only
`none`, `n/a`, or `no findings` after the colon. Any other headerless answer
blocks as ambiguous. A line that normalizes to `Findings` after an optional
trailing ASCII or fullwidth colon is removed starts the findings body. A
verdict-shaped line after that header is findings text and does not create a
conflict. Exactly one verdict must be present; both verdicts or neither verdict
block the review. Roundfix reads the verdict from the reviewer's final message.
The answer file keeps every message, separated by a blank line. The record and
the answer live under `pre-pr-review/<checkout key>/` in the Artifact
Directory, and the review record's `answerPath` names the answer file. Roundfix
never reads another checkout's record. Roundfix sets `answerPath` only after it
sends the prompt to a reviewer; a pre-prompt failure has no answer path or
answer file.

A findings record preserves the original `findings` text and assigns `F1`,
`F2`, and so on to the entries in `findingItems`. The reviewer prompt requires
one `- ` list item per finding with its file and line. Roundfix derives the
same identities when it reads an older findings record without items.

Use `roundfix review dispose` to append one disposition for one finding. Use
`--dismiss --evidence <text>` only at the reviewed `HEAD`, with non-blank
evidence. Use `--fixed-by <commit>` only for a resolving commit that differs
from and descends from the reviewed head and is reachable from the current
`HEAD`. Evidence is copied verbatim and never executed.

Each success appends one JSON line to
`pre-pr-review-dispositions.jsonl` in the Artifact Directory and prints the
same line. The append-only ledger ties the finding identity and text to the
repository and reviewed head, and records either the evidence or fixing commit
with an RFC 3339 UTC timestamp. A missing or mismatched findings record,
unknown identity, invalid or blank form, moved-head dismissal, invalid fixing
commit, or second disposition exits `2` with
`roundfix: review dispose refused:` on stderr and appends nothing.

For `codex` and `claude`, a findings verdict stands for its repository,
`baseCommit`, head commit, and provider while the base branch moves. When the
Artifact Directory already holds a matching `findings` or
`findings-dismissed` record, Roundfix reuses it before preparing, probing, or
prompting an Agent session. The reused record sets `reused` and carries exact
repository, head, finding-identity, and finding-text ledger matches in
`dispositions`.

When every finding has one evidence-backed `dismissed` disposition, the reused
record reports `findings-dismissed` and exits `0`. Otherwise it remains
`findings`, exits `1`, and stderr names every finding identity without a
disposition. A `fixed` disposition never clears that head; the changed
candidate needs a fresh review. A different repository, base, head, or provider
also gets a fresh review, and Roundfix never reuses `reviewed`, `blocked`, or
`omitted`. The `none` and `coderabbit` policies keep their behavior above.

When the candidate adds or changes a Spec folder under the configured Spec
Root, or under its resolved archive root, Roundfix discovers that folder from
the candidate diff. Specs archived within the candidate are read from the
archive root at `HEAD`. For each usable Spec, the review prompt carries the
PRD `Decisions` section and TechSpec. A changed Spec without a `## Decisions`
section, a PRD, or a TechSpec is skipped, its slug is listed in the record's
`skippedSpecs`, and the review proceeds with the remaining context. The
record's `specs` names the Specs whose context was carried.

The record's `archivedSpecs` field always lists the sorted slugs of changed
Spec folders under the resolved archive root, including folders whose context
was skipped; it is always present and is `[]` when the candidate archives no
Spec. When a findings verdict has archived Specs, stderr names those slugs,
says that an archived Spec is never corrected in place, and tells the operator
to author a corrective Spec with its own authorization and QA gate. Delivery
parks that review as `corrective-spec-required`; a no-findings verdict prints
no corrective-Spec line.

Spec context is bounded at 32 KiB per Spec and 64 KiB in total. When context is
truncated, the prompt includes `[Spec context truncated]` and the record sets
`specContextTruncated` to true. The reviewer must judge the delivery against
the carried decisions and the alternatives they reject. Specs dropped by the
total context bound are named in the record's `skippedSpecs`; the record's
`specs` names the Specs whose context was carried.

Exit codes are:

- Exit `0` — the reviewer returned one substantive no-findings verdict,
  explicit `none` recorded its configured omission, or a reused verdict is
  `findings-dismissed`.
- Exit `1` — the reviewer returned findings or a reused verdict still has
  standing findings; the record carries them.
- Exit `2` — preflight failed or the selected review is blocked.

A runtime failure, timeout, transport anomaly, empty output, or
unclassifiable output blocks the selected mode. Each is recorded as blocked
with its reason and never becomes a pass or a configured omission. A provider
selection failure may activate the next configured review fallback only before
the prompt is sent; failures after the prompt remain blocked review failures.

