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

Use `roundfix archive <slug>` after a Spec's Tasks are completed and the newest
QA Report is acceptable under the one declared-acceptance eligibility policy:
either `verdict: pass` with no disallowed blocked rows, or a `partial` verdict
whose only unmet rows other than the pre-PR Pull Request row are declared
unreachable and fully covered by the Spec's `## Unreachable Acceptance`
declarations. The pre-PR Pull Request row, recorded as `blocked (environment:
no open Pull Request)` with the Pull Request row named in its provenance,
never decides a qualifying partial and needs no Unreachable Acceptance
declaration. Settlement and archive both apply this same policy. A `fail`, an
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

For a passing verdict, archive stamps `_prd.md` with `status: archived`,
`archived`, and `source_slug`. For the declared-only `partial` case, it also
stamps the declarations' `satisfied-by` actions under `unproven`, so a reader
of the archived record learns what was never verified. It then moves
`<specs.root>/<slug>/` to its resolved archive directory. With the default
Spec Root, stdout carries the deterministic report:

```text
archived <slug> -> docs/history/specs/<slug>
```

Refusals exit `2` through Preflight Validation, name the first unmet condition
on stderr, and leave the active Spec folder in place. Every refusal outside the
one declared-only case is unchanged: a finding-blocked row, an
environment-blocked row other than the pre-PR Pull Request row, a declared count not covered by the Spec's
declarations, `verdict: fail`, missing QA, and any non-completed Task all
refuse. `qa_override` keeps its existing meaning for explicitly authorized
archival of genuinely failed or missing evidence; declared unreachability does
not use or weaken that override.

To archive despite failed, missing or otherwise ineligible QA, run:

```bash
roundfix archive <slug> --qa-override --approval <source> --reason <text>
```

The command requires both approval and reason, keeps every non-QA Task
`completed`, and is refused only when a normal archive would succeed. It stamps
the approval source, reason, observed QA outcome and archived revision. When the
QA Task is not completed, it also stamps `qa_override_qa_task_status`. It does
not change the QA Task or report verdict.

