## Supersede Command

Use `roundfix supersede --spec <slug> --by <slug> --reason <text>` when another
Spec delivered the content of a Spec that was never decomposed. This records a
supersession amendment in the superseded Spec. Do not edit that Spec's status or
write a note instead: the amendment is the lifecycle record that archive reads.

The `--spec` value names the superseded Spec, `--by` names the active or
archived Spec that delivered its content, and `--reason` records the
explanation. The command writes only `_supersession.md`; it creates no Run,
writes no Run Event Journal entry, and never commits or pushes.

The command refuses an unknown superseded Spec, an unknown superseding Spec, a
self-supersession, or a Spec that already carries a supersession. Exit `0`
means the amendment was recorded, exit `1` means the write failed, and exit
`2` means Preflight Validation failed.

A Spec cannot archive while another file names its active directory. The
command exits `2`, lists each file and line, and leaves every file in place.
This includes tracked and untracked non-ignored files other than Markdown,
outside the Spec Root, its archive root and `docs/history`. Replace a code or
test dependency on the Spec's files with a fixture or an exported constant,
then retry the archive. This refusal also applies with `--qa-override`.

### Authorization audit

The mechanical authorization audit reads each governed Task commit's grant at
its fork point first. When that grant does not cover the commit, the audit can
use the grant the Task ran under: the record in the Task commit's parent, but
only while the delivery target carries byte-identical content at the same
path. The audit reports the latest delivery-target commit that established
that content as the authorizing revision. A parent-only record, an older
record that the delivery target later narrowed or revoked, and a Task commit
that edits its own record remain refusals.

## Archive Command

Archive refuses a Spec with a Glossary Gap, also under a QA Archive Override.

Use `roundfix archive <slug>` after a Spec's Tasks are completed and the newest
QA Report is acceptable under the one declared-acceptance eligibility policy:
either `verdict: pass` with no disallowed blocked rows, or a `partial` verdict
whose only unmet rows are the pre-PR Pull Request row or an outside-evidence row
the Run sandbox could not reach, with any declared-blocked rows covered by the
Spec's `## Unreachable Acceptance` declarations. The pre-PR Pull Request row,
recorded as `blocked (environment: no open Pull Request)` with the Pull Request
row named in its provenance, and the outside-evidence row, recorded as
`blocked (environment: network denied: <host>)` with the outside-evidence row
named in its provenance, never decide a qualifying partial and need no
Unreachable Acceptance declaration. Any other blocked outside-evidence row
still blocks Pull Request preparation. Settlement and archive both apply this
same policy. A `fail`, an
undeclared `partial`, a missing or unparseable
report, and a `pass` carrying finding-, declared-, or precondition-blocked rows
all refuse; an environment-blocked row remains acceptable under the existing
policy. For a Spec without a Task Graph, a
recorded supersession is accepted in place of those completed Tasks and the
passing gate. Every other archive precondition still applies, and a Spec with
a Task Graph keeps its existing evidence rules. The command is
non-interactive, creates no Run, and never pushes. Before touching the
filesystem, it verifies every Task in the Spec's Task Graph has
`status: completed` and that the newest QA Report meets that eligibility
contract.

For a passing verdict, archive writes `<slug>.md` under the resolved archive
root and removes the Spec folder, which stays in Git at the record's
`source_revision`. For a qualifying `partial`, the record also carries the
declarations' `satisfied-by` actions under `unproven`, so a reader learns what
was never verified. With the default Spec Root, stdout carries the
deterministic report:

```text
archived <slug> -> docs/history/specs/<slug>.md; removed <n> file(s) (<b> bytes) kept in Git at <12-hex>
```

When files are promoted, the successful confirmation appends
`; promoted <n> file(s) to docs/references/`.
Existing archived folders retain their legacy link semantics.

Refusals exit `2` through Preflight Validation, name the first unmet condition
on stderr, and leave the active Spec folder in place. Every refusal outside the
one declared-only case is unchanged: a finding-blocked row, an
environment-blocked row other than the pre-PR Pull Request row or a
network-denied outside-evidence row, a declared count not covered by the Spec's
declarations, `verdict: fail`, missing QA, and any non-completed Task all
refuse. `qa_override` keeps its existing meaning for explicitly authorized
archival of genuinely failed or missing evidence; declared unreachability does
not use or weaken that override.

To archive despite failed, missing or otherwise ineligible QA, run:

```bash
roundfix archive <slug> --qa-override --approval <source> --reason <text>
```

The command requires both approval and reason, keeps every non-QA Task
`completed`, and is refused only when a normal archive would succeed. The record carries
`qa_override`, `qa_override_approval`, `qa_override_reason`,
`qa_override_qa_outcome` and `qa_override_revision`. When the QA Task is not
completed, it also carries `qa_override_qa_task_status`. It does
not change the QA Task or report verdict.
### Archive Record

The command writes `<slug>.md` under the resolved archive root and removes the
Spec folder. The record names the QA Report and verdict and carries its
`source_revision`, the Git revision where the removed folder remains
recoverable. Run `roundfix archive <slug> --plan` to inspect the cut and
Archive Advice before it happens. Use repeatable `--promote <path>` to copy a
confirmed candidate to `docs/references/` in the archive change.

## History Sanitize Command

Use `roundfix history sanitize` to inspect the existing history left by
archives before Archive Records. It is a dry run unless invoked with
`--apply --batch <n>`; `--advise` requests advisory classifications for the
next batch's candidate files, and `--promote <path>` copies an explicitly
selected file to `docs/references/` while applying a batch.

The command converts each Legacy Archive Folder into an Archive Record,
reduces retired Findings and Backlog Entries to Reduced History Entries, and
removes retired Review Artifacts and handoffs. A folder without a QA Report,
override or supersession receives the `no-qa` disposition. A failed QA without
an override is recorded as `failed-qa`, keeping its verdict and report name.
A Legacy Archive Folder is read with Lenient Legacy Reading: tolerated rows are
named, while an active Spec is read strictly. Lenient Legacy Reading accepts a
legacy list of maps in `unproven`, with one stable text line per map. A Refused
Unit is listed with its one-line reason and skipped, so it remains untouched and
does not consume a batch slot. `roundfix baseline update` plans these same units
at once; pass `roundfix baseline update --no-history` to keep this command's
reviewed-batch path.
Apply requires the annotated `history-full` tag on an ancestor of `HEAD` and
holding every path the batch rewrites or removes, so the full history remains
reachable in Git.
Each `--apply --batch <n>` is one reviewable, revertible Pull Request; the
command never commits, tags, pushes or opens one.
