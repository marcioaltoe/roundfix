### archive

```bash
roundfix archive <slug>
roundfix archive <slug> --qa-override --approval <source> --reason <text>
roundfix archive <slug> --plan [--offline]
roundfix archive <slug> --pull-request <n>
roundfix archive <slug> --promote <path> [--promote <path> ...]
```

Non-interactive; creates no Run and never pushes. Verifies every Task is
`completed` and accepts either `verdict: pass` or a declared-only `partial`
verdict. Declared-only means every unmet row is declared unreachable and fully
covered by the Spec's `## Unreachable Acceptance` declarations. For that case,
archive records the declarations' `satisfied-by` actions under `unproven` in
the Archive Record, so it names what was never verified. Each entry is read as
text of the Spec's `_prd.md`, so its relative links are rewritten like the outcome's, and an
entry without one keeps its bytes (ADR-0308). It writes
`docs/history/specs/<slug>.md` for the built-in Spec Root, or
`<specs.root>/_archived/<slug>.md` for another configured root, then removes
`<specs.root>/<slug>/`. The archive performs an Archive Relink before writing:
Live Document citations are rewritten to promoted copies or the Archive Record.

In the outcome, each `unproven` entry and each promoted copy, a citation of a
removed folder file points at its promoted copy or becomes a Git Locator: `path@commit`, written as
code. The locator prefers `source_revision` when the Default Branch Head or a
tag reaches it; otherwise it names the newest commit the Default Branch Head
reaches that holds the file. If no such commit exists, it names the path and
the record's pull request, or says `not reachable from any branch or tag`.
The Default Branch Head is the fetched `origin/<default>`, the local branch only
without one, with `origin/main` or `origin/master` naming the default when
origin/HEAD is unset and exactly one exists.
A removed folder destination never becomes a link from the record to itself.
Outward relative destinations are rebased to keep naming the same file.
The archive reads every path relative to the checkout it runs in, so it
writes the same record and relinks the same files from the main checkout,
a linked worktree such as a Delivery Queue item worktree, a subdirectory
of either, or a path that crosses a symbolic link.

A path whose resolved target lies outside the repository, such as a
`docs/history` or `docs/references` symbolic link to another directory, is
refused with exit `2` before anything is written, and the refusal names the
path. A symbolic link whose target stays inside the repository, a checkout
reached through a symbolic link and an external Spec Root keep working.

The original folder's bytes stay in Git at the record's `source_revision`. The
Archive Record is the small `<slug>.md` file under the archive root; it names
the disposition, QA Report and verdict, source revision, and promoted files.
Its summary is the QA Report's Outcome paragraph or the PRD's title paragraph.
A summary paragraph that ends with a colon carries the list or paragraph that follows it, on the same line.
After a successful archive, the command warns when `source_revision` is not an
ancestor of the Default Branch Head. `roundfix archive <slug>` with
`--pull-request <n>` stores the delivery pull request in the Archive Record,
including a record the current branch already added.

Use `roundfix archive <slug> --plan` to list the core artifacts, QA evidence,
and candidate files that the cut removes. Candidate files receive optional
Archive Advice from Jev. Advice is advisory and is skipped when the judge key
is absent or the monthly ceiling is reached. The plan leaves repository files
unchanged. Use repeatable `--promote <path>` to copy a confirmed candidate to
`docs/references/` during the archive change.
When `ROUNDFIX_RUN_ID` is non-empty in the command's environment, every
candidate has `no advice (inside Run <run-id>, named by ROUNDFIX_RUN_ID)`.
Archive Advice selects no transport, reads no key or Judge Log, writes no Judge
Log and sends no request, even with a judge key present.
Every Agent Session and Verification command of a Run receives this marker
with the Run's ID. Roundfix also withholds the Judge Key Variables, including
`TYPESAFE_API_KEY`, from those children; the marker keeps the judge silent even
if a runtime restores a key from the user's shell startup.

For a plan that contacts no judge, use `roundfix archive <slug> --plan --offline`.
It skips loading the judge and reading Judge Key Variables, prints
`judge: skipped (offline)` after the evidence line, and lists each candidate in
file order as `candidate <file> <bytes> bytes: no advice (offline)`.
A QA Task always uses this offline form. `--offline` without `--plan` exits
`2` with `--offline requires --plan`. `--plan` alone keeps Archive Advice outside a Run.

Existing archived folders retain their legacy link semantics.

A Spec that will not be delivered is cancelled with `roundfix cancel`.

A plan prints:

```text
archive plan for <slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md
core <k> file(s): _prd.md, _techspec.md, _tasks.md, _authorization.md, <task files>, <QA Reports>, references/_index.md and adopted sources
evidence <e> file(s) (<eb> bytes) under qa/evidence/
candidate references/<file>.md <bytes> bytes: no advice (<KEY_VARIABLE> is not set)
promote with: roundfix archive <slug> --promote <path>
```

The Archive Record lists changed citing files under its optional `relinked`
key. A cited non-core file that a Live Document references is promoted
automatically to `docs/references/<basename>`; citations of the folder or a
core artifact target the Archive Record. Citations in non-Live Documents are
left unchanged and reported as unrelinked.

The archive refuses without changing files when a promotion destination exists,
when two promotions share a basename, or when a Live Document to be rewritten
has changes not committed at `HEAD`. A Git planning failure in a work tree
also refuses. Its `--plan` output adds these lines before `promote with:`:

```
promotes <file> to docs/references/<basename>
relinks <path>: <n> citation(s)
leaves <path>:<line> naming <target> unrelinked
```

An unrelinked citation writes
`roundfix: warning: archive leaves <path>:<line> naming removed <target> unrelinked`
to stderr. When the archive changes one or more Live Documents, the success
line appends `; relinked <m> file(s)` after the promotion suffix.

The Spec folder must match its repository's `HEAD`. Modified, staged, deleted
or untracked Spec files cause exit `2`; commit them before archiving so the
recorded revision holds every removed file. For an external Spec Root, this
revision comes from that root's repository.

The archive also refuses with exit `2`, writing nothing, while Git ignores a
file under the Spec folder. It lists every ignored file by its sorted,
repository-relative path, including files inside ignored directories. To
continue, remove each file or commit it with git add -f, then retry.
Uncommitted changes keep their existing refusal and are reported first.

If the final Spec folder removal fails, the command keeps the Archive Record
and every other write, prints nothing on standard output, and exits `1`.
Standard error names the record, the remaining Spec folder and the cause, then
prints the exact command that finishes the removal:

```text
roundfix: archive wrote <record> but could not remove the Spec folder <folder>: <cause>
roundfix: finish the removal with: chmod -R u+w '<absolute-folder>' && rm -rf '<absolute-folder>'
```

Run the printed command to restore write permission and finish removing the
committed folder. Its absolute path is quoted for a POSIX shell, including
spaces and embedded single quotes. Promotions and relinked documents stay
written, and `--pull-request` is recorded even when removal fails. Retrying
`roundfix archive` refuses because the record already exists.

A normal archive prints:

```text
archived <slug> -> docs/history/specs/<slug>.md; removed <n> file(s) (<b> bytes) kept in Git at <12-hex>
```

When files are promoted, the successful confirmation appends
`; promoted <n> file(s) to docs/references/`; relinked files add the suffix
`; relinked <m> file(s)`.

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
outside-evidence row recorded as `blocked (environment: network denied: <host>)` or
`blocked (environment: network prohibited by Spec authorization)`, a declared count not covered by the Spec's
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
