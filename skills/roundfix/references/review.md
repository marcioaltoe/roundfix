## Pre-PR review

```bash
roundfix review [--base <ref>]
roundfix review dispose <finding-id> --dismiss --evidence <text>
roundfix review dispose <finding-id> --fixed-by <commit>
```

Runs the configured pre-Pull-Request reviewer over the current candidate. The
workflow computes the candidate diff from the merge base of the current head
and the selected base, then hands that diff to the configured reviewer in a
read-only session; it does not ask the reviewer to discover the candidate. The
resulting record names the repository, that merge base as `baseCommit`, the
resolved base tip as `baseTipCommit`, the head commit, effective provider, and
policy source. A head with no shared history with the selected base exits `2`
before any reviewer call or readiness probe.

`--base <ref>` selects the base Git ref. When omitted, Roundfix uses the
repository's default branch.

The configured policy accepts `codex`, `claude`, `coderabbit`, and `none`.
Explicit `none` performs no reviewer call and no readiness probe, records a
configured omission, and exits `0`. `claude` runs through the same read-only
review path as `codex`. `coderabbit` remains refused because no supported local
CodeRabbit review surface is installed or specified; it exits `2` as blocked
instead of being treated as an omission.

Roundfix classifies the reviewer's answer by substance rather than exact
formatting. It recognizes a `No findings` verdict after case folding,
surrounding Markdown emphasis, and trailing punctuation are normalized, or a
`Findings:` verdict with its findings text. A pass must be the whole answer:
after normalization, the answer is `No findings`, or its only content line is
`Findings:` with only `none`, `n/a`, or `no findings` after the colon. Any other
headerless answer blocks as ambiguous. A line that normalizes to `Findings`
after an optional trailing ASCII or fullwidth colon is removed starts the
findings body. A verdict-shaped line after that header is findings text and does
not create a conflict. Exactly one verdict must be present; both verdicts or
neither verdict block the review. Roundfix reads the verdict from the
reviewer's final message. The answer file keeps every message, separated by a
blank line. The record `pre-pr-review.json` and answer
`pre-pr-review-answer.txt` live under `pre-pr-review/<checkout key>/` in the
Artifact Directory, and the review record's `answerPath` names that answer
file. Roundfix never reads another checkout's record. Roundfix sets
`answerPath` only when the prompt reached a reviewer; a pre-prompt failure has
no answer path or answer file.

A findings record keeps the reviewer's original `findings` text and also lists
each finding as `F1`, `F2`, and so on in `findingItems`. The reviewer prompt
asks for one `- ` list item per finding. Start each item with `path:line` or
`path:start-end`, optionally backticked, naming a line of the candidate diff.
For a comma-listed anchor, only the first range counts. Describe what breaks
in a clause beginning with `Failure:`. An older record without `findingItems` derives the same identities from its findings text when
Roundfix reads it.

The prompt carries Delivery Conventions version
`roundfix/delivery-conventions/v1`. These describe what a delivery writes by
design:

- C1. A Spec's QA Report records the head it audited and is committed after
  that head, so it never names the commit that records it.
- C2. The Daemon writes a Task file's status and its `## Result`,
  `## Recorded paths` and `## Carry-forward provenance` sections after the
  Task's Verification passes; a Result that calls status Daemon-owned agrees
  with a `completed` status.
- C3. The archive commit moves a completed Spec's directory to the archive
  root and stamps its archive front matter.
- C4. A planning candidate authors a Spec whose Tasks are all pending and
  which has no QA Report; that Spec's own delivery implements it and is
  reviewed then.

Roundfix checks each anchor against the same candidate diff supplied in the
prompt. Context and added lines count, as does the new-side start line of a
pure deletion. Deleted files and files changed without a hunk, such as binary
or rename-only files, hold every line. A range counts when any line overlaps
a held line. A missing or out-of-diff anchor records
`dismissed-by-validation`, rule `unanchored`, with its reason, and prints a
diagnostic on stderr.

Validation has three dismissal rules. `unanchored` is the mechanical anchor
check above. `convention:C1` through `convention:C4` requires the whole anchor
to lie in that convention's mechanically computed region and a sealed
validator judgment that the finding merely restates it. `no-failure` requires
that the finding has no `Failure:` clause and the validator judges that it
states no failure. A finding with a failure clause cannot receive that rule.
Other anchored findings stand. Each item's validation records `stands` or
`dismissed-by-validation`, with the dismissal's rule and reason.

Only anchored findings inside a convention region or lacking a failure clause
are sent to the validator. It uses the review's selected runtime in a sealed
session and may use no tools. The validator fails closed: an unavailable
runner, timeout, tool use, malformed answer, missing or duplicate finding ID,
ineligible rule or blank reason leaves every asked finding standing. The
record's `validation` identifies the conventions version and reports
`validator` as `not-needed`, `ran` or `unavailable`, with the unavailable
reason. Each validation dismissal prints its ID, rule and reason on stderr.

If every finding is dismissed by validation, the outcome is
`findings-dismissed` with exit `0`, without writing a ledger disposition. If
none carries a readable anchor, the review blocks with exit `2` and reason
`findings name no file and line`. Records without validation fields remain
readable and their findings count as standing.

Use `roundfix review dispose` to record one disposition for one finding. A
dismissal requires non-blank evidence and an unchanged reviewed `HEAD`. A fix
requires a resolving commit that differs from and descends from the reviewed
head and is reachable from the current `HEAD`. Evidence is copied as text and
is never executed.

Successful dispositions append one JSON line to
`pre-pr-review-dispositions.jsonl` in the Artifact Directory and print that
same line. The line ties the finding's identity and text to its repository and
reviewed head, and records either `evidence` or `fixedBy` with an RFC 3339 UTC
timestamp. The ledger is append-only. Roundfix refuses a missing or mismatched
findings record, an unknown identity, invalid or blank forms, moved-head
dismissals, invalid fixing commits, and a second disposition. It also refuses
a finding already `dismissed-by-validation` with
`finding "<ID>" was dismissed by validation (<rule>)`. Every refusal
exits `2`, starts stderr with `roundfix: review dispose refused:`, and appends
nothing.

For `codex` and `claude`, a findings verdict stands for its repository,
`baseCommit`, head commit, and provider while the base branch moves. When the
Artifact Directory already holds a matching `findings` or
`findings-dismissed` record, Roundfix reuses it before preparing, probing, or
prompting an Agent session. The reused record sets `reused` and carries the
ledger entries whose repository, head, finding identity, and text match its
`findingItems` in `dispositions`.

When every finding is dismissed by validation or has one evidence-backed
`dismissed` disposition, the reused record reports `findings-dismissed` and
exits `0`. Otherwise it remains `findings`, exits `1`, and stderr names every
standing finding identity without a disposition. A `fixed` disposition never
clears that head; the fix belongs to a changed candidate. The `none` and
`coderabbit` policies keep their behavior above.

Reviews in one checkout form a **Reviewer Lineage** when the provider and
merge base stay the same and each reviewed head descends from the last.
Round 1 reviews the full candidate diff. Round 2 reviews only the delta from
the round-1 head to the current head; its prompt always includes every round-1
finding with its validation and operator disposition. A fixed or dismissed
finding is raised again only if the delta reintroduces it. Anchors are still
validated against the full merge-base candidate diff in either round.
A rebase, changed provider or merge base, or a head that does not descend
starts a new lineage at round 1. A blocked review retries its own round with
its recorded previous data.

When round 1 leaves standing findings, its reviewer session stays open. The
checkout record stores the random session name, the selection index into the
review profile's preferred selection and fallbacks, and `sessionOpen: true`.
Round 2 prepares that name on that selection through acpx. If preparation
fails, Roundfix ends that name and tries the normal selection loop with fresh
random names. The round-2 prompt carries the recorded findings even when
resume is unavailable.

Every round records the ACP session ID reported by its prompt stream in
`lineage.acpSessionIds`, or an empty string when absent. `continued` is true
only when the recorded session was prepared and both rounds report the same
non-empty ID. Reusing a name alone does not prove continuation. Round 2 always
ends its session, including on a blocked or findings verdict. Every other
round-1 outcome ends its session; a lineage change ends the prior open session
before the new review, and a same-head reuse that resolves to
`findings-dismissed` closes it too. A removed checkout may leave an open
session record; Roundfix does not clean up another checkout's session.

A third review in a lineage never calls the reviewer. It closes as
`ceiling-closed`, with exit `0`, only when every standing round-2 finding has
exactly one operator disposition: dismissed with evidence, or fixed by a
commit the current head contains. Otherwise it prints a blocked result with
exit `2` naming the findings to dispose. That blocked result is not persisted,
so `roundfix review dispose` can still read the round-2 record. The closed
record names the round-2 head as `lineage.reviewedHead`, and Delivery advances
on `ceiling-closed`.

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
`specContextTruncated` to true. The reviewer judges the delivery against the
carried decisions and the alternatives they reject. Specs dropped by the total
context bound are named in the record's `skippedSpecs`; the record's `specs`
names the Specs whose context was carried.

Exit codes:

- `0` — one substantive no-findings verdict, configured omission, or a
  `findings-dismissed` or `ceiling-closed` verdict.
- `1` — the validated reviewer verdict or a reused verdict still has standing
  findings; the record carries them.
- `2` — preflight failed or the review was blocked.

Runtime failure, timeout, transport anomaly, empty output, and unclassifiable
output each block the selected mode, with a reason in the record. None can
become a pass or an omission. A configured selection fallback is eligible only
when selection fails before the prompt is sent; failures after the prompt are
review failures.

