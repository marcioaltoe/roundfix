### archive

```bash
roundfix archive <slug>
roundfix archive <slug> --qa-override --approval <source> --reason <text>
```

Non-interactive; creates no Run and never pushes. Verifies every Task is
`completed` and accepts either `verdict: pass` or a declared-only `partial`
verdict. Declared-only means every unmet row is declared unreachable and fully
covered by the Spec's `## Unreachable Acceptance` declarations. For that case,
archive stamps the declarations' `satisfied-by` actions under `unproven` in
`_prd.md`, so the archived record names what was never verified. It then moves
`<specs.root>/<slug>/` to `docs/history/specs/<slug>/` for the built-in Spec
Root, or to `<specs.root>/_archived/<slug>/` for any other configured root.

Before the move, the archive rewrites each relative Markdown link that leaves
the Spec — including inline links, images and reference definitions outside
code blocks and code spans — so it reaches the same path from the archived
location. It keeps the link's fragment, query and angle-bracket form. It keeps
a link whose target was already archived when its unchanged destination reaches
that target from the archived location. If any other relative links that leave
the Spec and do not resolve remain, the command exits `2` before changing any
file and names each link with its file, line and destination. When links were
rewritten, the successful confirmation appends `; rewrote <n> relative link(s)`.
Destinations inside the Spec, absolute destinations, URLs, HTML anchors and
non-Markdown files are not rewritten, including evidence scripts that climb a
fixed number of directories.


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
normally. It stamps the approval source, reason, observed QA outcome and
archived revision without changing the QA Task or report verdict. When the QA
Task is not completed, it also stamps `qa_override_qa_task_status` with that
status; a completed QA Task omits the field. When the newest report is
unreadable, the recorded outcome names it relative to the Spec folder and never
stores an absolute machine path.
