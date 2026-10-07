### archive

```bash
roundfix archive <slug>
roundfix archive <slug> --qa-override --approval <source> --reason <text>
roundfix archive <slug> --plan
roundfix archive <slug> --promote <path> [--promote <path> ...]
```

Non-interactive; creates no Run and never pushes. Verifies every Task is
`completed` and accepts either `verdict: pass` or a declared-only `partial`
verdict. Declared-only means every unmet row is declared unreachable and fully
covered by the Spec's `## Unreachable Acceptance` declarations. For that case,
archive records the declarations' `satisfied-by` actions under `unproven` in
the Archive Record, so it names what was never verified. It writes
`docs/history/specs/<slug>.md` for the built-in Spec Root, or
`<specs.root>/_archived/<slug>.md` for another configured root, then removes
`<specs.root>/<slug>/`. The PRD and relative links are never rewritten.
The original folder's bytes stay in Git at the record's `source_revision`. The
Archive Record is the small `<slug>.md` file under the archive root; it names
the disposition, QA Report and verdict, source revision, and promoted files.

Use `roundfix archive <slug> --plan` to list the core artifacts, QA evidence,
and candidate files that the cut removes. Candidate files receive optional
Archive Advice from Jev. Advice is advisory and is skipped when the judge key
is absent or the monthly ceiling is reached. The plan leaves repository files
unchanged. Use repeatable `--promote <path>` to copy a confirmed candidate to
`docs/references/` during the archive change.
Existing archived folders retain their legacy link semantics.

The Spec folder must match its repository's `HEAD`. Modified, staged, deleted
or untracked Spec files cause exit `2`; commit them before archiving so the
recorded revision holds every removed file. For an external Spec Root, this
revision comes from that root's repository.

A normal archive prints:

```text
archived <slug> -> docs/history/specs/<slug>.md; removed <n> file(s) (<b> bytes) kept in Git at <12-hex>
```

An override inserts ` with QA override` after the slug. An existing Archive
Record or legacy folder for the same slug refuses the archive before any file
changes. Broken outward Markdown links retain their existing preflight refusal.

A Spec cannot archive while another file names its active directory. The
command exits `2`, lists each file and line, and leaves every file in place.
This includes tracked and untracked non-ignored files other than Markdown,
outside the Spec Root, its archive root and `docs/history`. Replace a code or
test dependency on the Spec's files with a fixture or an exported constant,
then retry the archive. This refusal also applies with `--qa-override`.

A `pending` verdict is never accepted. A `pass` or otherwise-eligible `partial`
that records no QA row is refused before archive changes the Spec.
A report whose front matter is empty or duplicated is unreadable and refused;
archive leaves the Spec and report in place.

Every other refusal is unchanged: a finding-blocked row, an
environment-blocked row other than the pre-PR Pull Request row or an
outside-evidence row recorded as `blocked (environment: network denied: <host>)`, a declared count not covered by the Spec's
declarations, or `verdict: fail` exits `2` and names the first unmet condition.
`qa_override` keeps its existing meaning for explicitly authorized archival
when the normal QA prerequisite is unmet; declared unreachability does not use
or weaken that override.

Use `--qa-override` only with explicit authorization to archive despite failed,
missing or otherwise ineligible QA. `--approval <source>` records who or what
authorized the override, and `--reason <text>` records why. The command still
requires every non-QA Task to be `completed`. It accepts a failed or pending QA
Task regardless of the newest report's verdict, and refuses only when every
Task is `completed` and that report qualifies because the same Spec can archive
normally. It records the approval source, reason, observed QA outcome and
archived revision without changing the QA Task or report verdict. When the QA
Task is not completed, it also records `qa_override_qa_task_status` with that
status; a completed QA Task records an empty value. When the newest report is
unreadable, the recorded outcome names it relative to the Spec folder and never
stores an absolute machine path.

After checking active-directory pins, archive refuses a Spec with a Glossary
Gap: an undeclared bold term or malformed declaration, an added or changed
term without a binding glossary-writing Task, or a declared term still missing
after every binding Task completed. It exits `2`, changes no file, and reports
`Spec "<slug>" cannot archive with a Glossary Gap: <code>: <summary>`, joining
multiple findings with `; `. This refusal also applies under a QA Archive
Override (`--qa-override`); the override waives only the QA prerequisite.
See [Glossary Declaration](spec.md#glossary-declaration) for the codes,
matching rules and glossary horizon.
