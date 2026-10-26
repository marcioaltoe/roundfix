---
name: archive-spec
description: Archive a completed spec — verify every task completed, QA passed, and indexed references are self-contained, then write <slug>.md under the resolved archive root (<spec-root>/_archived/, or docs/history/specs/ for the built-in docs/specs root) and remove the Spec folder, which stays in Git at the record's source_revision. Runs automatically at the end of the implement-spec loop after a QA pass, or whenever the user asks to archive a spec.
argument-hint: "<spec slug>"
metadata:
  category: delivery
  tags: [workflow, documentation, process]
  version: 0.0.15
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
version: 0.0.15
---

# Archive Spec

Archive a completed spec by writing `<slug>.md` under the resolved archive root and removing `<resolved-spec-root>/<slug>/`, which stays in Git at the record's `source_revision`. The source is the configured Spec Root; the destination is its resolved archive root — `docs/history/specs/` for the built-in `docs/specs` root, or `<spec-root>/_archived/` for an external or non-default root, matching the `roundfix archive` destination. A normal archive means _implemented, verified, and self-contained_ — every task done, QA passed, and every indexed reference owned by the Spec. A QA Archive Override records the explicit exception without claiming verification. After either disposition, one `ls <spec-root>/` separates live work from history, and the archive stays greppable as the record of what was built and why.

The trigger is spec completion, not publication: run this automatically at the end of the `implement-spec` loop once the QA gate passes, or whenever the user asks. Merge and release are separate, user-driven steps — the archive commit simply travels with the branch and ships inside the feature's own PR.

## Archive Relink

Before writing the Archive Record, the archive performs an Archive Relink. It
automatically promotes each cited non-core file that a Live Document references
from the owner Spec to
`docs/references/<basename>`, then rewrites citations from each Live Document
to the promoted copy; citations of the Spec folder or a core artifact point to
the Archive Record. The record lists rewritten files under `relinked` when any
exist. The source keeps one owner, so a secondary Spec needs no copy (ADR-0260,
ADR-0083).

The Spec folder retires by deletion. In the copied outcome and each promoted
copy, a destination inside the removed folder names a promoted copy or becomes
a Git Locator, `` `<path>@<commit>` ``; it never points at the record itself.
The commit holds the path and is reachable from the default branch or a tag.
When no such commit exists, the locator names the path and the record's pull
request. Other destinations retain their repository target (ADR-0278).

The archive refuses to write when a promotion destination already exists, when
two promotions share a basename, or when a Live Document to be rewritten has
changes not committed at `HEAD`. Citations in files that are not Live
Documents, or citations that remain after rewriting, are warnings naming the
file, line and target; the archive leaves those files unchanged.

### Squash merges

A squash merge keeps `source_revision` only in the pull request's head ref; the
record names that pull request in `pull_request`.

### QA settlement

The same outcome settles the authored `qa` Task and determines what the archive leaves:

| Outcome | Settles | Archives |
| --- | --- | --- |
| `pass` | Settles the QA Task as `completed` and makes the Spec archive-eligible when the report has no disallowed blocked rows. | The Archive Record, which names the QA Report and verdict; the Spec, its QA report and evidence stay in Git at the record's `source_revision`. |
| qualifying declared `partial` | Settles the QA Task as `completed` when no row failed, was skipped or is finding-blocked, every declared-blocked row is covered by a matching `## Unreachable Acceptance` declaration, and every environment-blocked row is the pre-PR Pull Request row, recorded as `blocked (environment: no open Pull Request)` with the Pull Request row named in its provenance, or an outside-evidence row the Run sandbox could not reach, recorded as `blocked (environment: network denied: <host>)` with the outside-evidence row named in its provenance. Neither row needs an Unreachable Acceptance declaration, and a partial whose only unmet rows are such rows qualifies. | The Archive Record, which names the QA Report and verdict and carries the declarations' `satisfied-by` record as `unproven`; the Spec, its QA report and evidence stay in Git at the record's `source_revision`. |
| `environment-blocked` | Leaves the row blocked; the report can still settle as `pass` when equivalent evidence satisfies the environment policy. | Nothing by itself; a qualifying report can archive the Spec. |
| `failed` | Leaves the QA Task unresolved and refuses archive unless an authorized override applies. | Nothing. |
| `missing` | Leaves the QA Task unresolved and refuses archive unless an authorized override applies. | Nothing. |
| `override` | Does not change the QA Task status or report verdict; settles archive as explicitly authorized despite failed or missing QA. | The Archive Record with `qa_override`, `qa_override_approval`, `qa_override_reason`, `qa_override_qa_outcome`, `qa_override_qa_task_status` when the QA Task is incomplete, and `qa_override_revision`; the QA Task and Reports stay unchanged in Git at the record's `source_revision`. |

### Rows the Spec's authorization keeps offline

The qualifying declared `partial` above also accepts an outside-evidence row
the gate did not look up because the Spec's own authorization forbids that
lookup, recorded as `blocked (environment: <subject> prohibited by Spec authorization)`
with the outside-evidence row named in its provenance. The subject names what
the authorization forbids, in its words: `network` for a network lookup, or
`another repository`; it is 1 to 80 characters and holds no parenthesis. Like
the network-denied row, it never decides a qualifying partial and needs no
Unreachable Acceptance declaration, a row counts toward one exempt kind at
most, and settlement, `roundfix qa-report accept` and `roundfix archive` read
it through the same policy (ADR-0276, ADR-0335).

## Cancelling a Spec

A Spec that will not be delivered is cancelled with `roundfix cancel <slug>
--reason <text>`, never deleted by hand. The command writes an Archive Record
with disposition `cancelled`, returns each adopted Backlog Entry to
`docs/backlog/` with `status: open` and each adopted Finding to
`docs/findings/` with `status: pending`, then removes the Spec folder while its
bytes stay in Git at `source_revision`. A cancelled record never meets a
prerequisite, never supersedes another Spec and is refused by Delivery Retry.

## Ignored files and a failed removal

The Archive Command and the Cancel Command refuse while Git ignores any file
under the Spec folder. Each command lists every ignored file and writes
nothing. The operator must remove each file or commit it with git add -f, then
retry.

If the final removal of the Spec folder fails, both commands keep every write,
exit `1`, and print the command that finishes the removal:
`chmod -R u+w '<folder>' && rm -rf '<folder>'`.

## Preconditions — verify, don't trust

Check all three with fresh command evidence before touching anything:

1. **Required tasks completed.** For a normal archive, read each `task_NN.md`
   listed in `_tasks.md`; every `status` must be `completed`. For a QA Archive
   Override, every non-QA Task must be `completed`; the QA Task may remain
   `pending` or `failed`. If the QA Task is also `completed`, the newest report
   must remain ineligible or the override is unnecessary and refused.

   **Command:** run `grep -n '^status:' <each-task-file-listed-in-_tasks.md>` and
   retain its output. Any value other than `status: completed` blocks a normal
   archive. During an override, only a non-QA Task with another status blocks
   the archive; name that Task file.

2. **QA passed.** The newest report in `qa/` must pass the repository's QA
   verifier. Do not substitute a line grep for structured validation. In
   Roundfix, use the same `internal/spec.QAVerdict` contract as the Archive
   Command: select the newest report by the `qa-report-YYYY-MM-DD[-NN].md`
   filename contract, parse its YAML frontmatter, require a supported
   `verdict`, require both blocked-row fields to be non-negative integers when
   present, and reject `verdict: pass` when `rows_blocked_finding` is nonzero.
   Retain the verifier's report path and result as evidence. A missing `qa/`
   directory, malformed newest report, or non-passing verdict blocks the
   normal archive. Proceed only under an explicit QA Archive Override performed
   through `roundfix archive <slug> --qa-override --approval <source> --reason
   <text>`. Never record the override by editing the Archive Record. A
   qualifying newest report does not make a failed or pending QA Task completed;
   the override is refused only when the report qualifies and every Task is
   completed.

3. **The Spec is self-contained.** Apply this precondition when
   `docs/specs/<slug>/references/_index.md` exists or is a symbolic link; a
   symbolic link, including a broken one, is invalid index state rather than a
   legacy Spec. Only a Spec where that path neither exists nor is a symbolic
   link predates this contract and passes without retrofitting historical
   artifacts. For an indexed Spec, every indexed `path` must exist relative to
   `_index.md`, every never-updated `source` path must be absent, and no
   Markdown link destination inside the Spec may point into `docs/_inbox/` or
   `docs/findings/`.

   **Commands:** first run this link-destination check; its syntax deliberately
   matches inline or reference-style Markdown links, not prose that merely
   names either tree:

   ```bash
   spec_dir=docs/specs/<slug>
   link_hits=$(grep --include='*.md' -RInE '(\]\([^)]*(docs/)?(_inbox|findings)/[^)]*\)|^[[:space:]]*\[[^]]+\]:[[:space:]]*<?[^[:space:]>]*(docs/)?(_inbox|findings)/)' "$spec_dir")
   link_status=$?
   if test "$link_status" -eq 0; then
     printf '%s\n' "$link_hits"
     echo "self-containment failed: rewrite each listed link at adoption step 8"
     exit 1
   fi
   test "$link_status" -eq 1 || exit "$link_status"
   ```

   Then parse and validate every data row. The index belongs to the current
   Spec: `owner` must equal its four-digit prefix, `type` must be `inbox`,
   `finding`, or `backlog`, and each `source` and `path` must appear only once. A `path` must
   be one basename relative to `_index.md`; reject absolute paths, `.`, `..`,
   path separators, and symbolic links instead of allowing traversal or a link
   outside `references/`. Run the following from the repository root and
   retain the normalized rows plus any diagnostic as evidence:

   ```bash
   slug=<slug>
   index="docs/specs/$slug/references/_index.md"
   expected_owner=${slug%%-*}
   parsed_index=$(mktemp) || exit 1
   trap 'rm -f "$parsed_index"' EXIT HUP INT TERM

   if test -L "$(dirname "$index")"; then
     printf 'self-containment failed: references/ must not be a symbolic link\n' >&2
     exit 1
   fi

   if test -L "$index" || test ! -f "$index"; then
     printf 'self-containment failed: references/_index.md must be a regular, non-symbolic-link file\n' >&2
     exit 1
   fi

   awk -F '|' -v expected_owner="$expected_owner" '
   function trim(value) {
     gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
     return value
   }
   function reject(message) {
     print "self-containment failed: " message > "/dev/stderr"
     invalid = 1
   }
   function separator(value) {
     value = trim(value)
     return value ~ /^:?-{3,}:?$/
   }
   /^[[:space:]]*$/ || /^[[:space:]]*#/ { next }
   !header {
     if (NF != 7 || trim($2) != "source" || trim($3) != "type" ||
         trim($4) != "owner" || trim($5) != "adopted date" ||
         trim($6) != "path") {
       reject("_index.md must use the fixed source | type | owner | adopted date | path header")
     } else {
       header = 1
     }
     next
   }
   !divider {
     if (NF != 7 || !separator($2) || !separator($3) || !separator($4) ||
         !separator($5) || !separator($6)) {
       reject("_index.md has an invalid table separator")
     } else {
       divider = 1
     }
     next
   }
   {
     source = trim($2)
     type = trim($3)
     owner = trim($4)
     adopted = trim($5)
     path = trim($6)
     if (NF != 7 || source == "" || type == "" || owner == "" ||
         adopted == "" || path == "") {
       reject("invalid index row at line " NR ": " $0)
       next
     }
     if (type != "inbox" && type != "finding" && type != "backlog") {
       reject("type must be `inbox`, `finding`, or `backlog` at line " NR ": " type)
     }
     if (owner != expected_owner) {
       reject("owner must be " expected_owner " at line " NR ": " owner)
     }
     if (seen_source[source]++) {
       reject("duplicate source at line " NR ": " source)
     }
     if (seen_path[path]++) {
       reject("duplicate path at line " NR ": " path)
     }
     if (adopted !~ /^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$/) {
       reject("adopted date must be YYYY-MM-DD at line " NR ": " adopted)
     }
     if ((type == "inbox" && source !~ /^docs\/_inbox\/[^\/]+\.md$/) ||
         (type == "finding" && source !~ /^docs\/findings\/[^\/]+\.md$/) ||
         (type == "backlog" && source !~ /^docs\/backlog\/[^\/]+\.md$/)) {
       reject("source does not match type at line " NR ": " source)
     }
     print source "|" type "|" owner "|" adopted "|" path
   }
   END {
     if (!header || !divider) {
       reject("_index.md is missing its fixed table header")
     }
     if (invalid) {
       exit 1
     }
   }
   ' "$index" > "$parsed_index" || exit $?

   while IFS='|' read -r source type owner adopted path; do
     case "$path" in
       ""|.|..|/*|*/*|*\\*)
         printf 'self-containment failed: path must be one basename relative to `_index.md`: %s\n' "$path" >&2
         exit 1
         ;;
     esac
     source_basename=${source##*/}
     if test "$path" != "$source_basename"; then
       printf 'self-containment failed: path must equal source basename for %s: %s != %s\n' "$source" "$path" "$source_basename" >&2
       exit 1
     fi
     current="$(dirname "$index")/$path"
     if test -L "$current" || test ! -f "$current"; then
       printf 'self-containment failed: invalid or missing path %s; finish adoption step 7\n' "$path" >&2
       exit 1
     fi
     if test -e "$source" || test -L "$source"; then
       printf 'self-containment failed: source still exists at %s; finish adoption step 7\n' "$source" >&2
       exit 1
     fi
   done < "$parsed_index"
   cat "$parsed_index"
   ```

A QA Archive Override performed only through `roundfix archive <slug>
--qa-override --approval <source> --reason <text>` overrides the unmet normal QA
prerequisite: an incomplete QA Task, an ineligible newest report, or missing or
unreadable QA evidence. The command is refused only when every Task is completed
and the newest report qualifies, because that Spec can archive normally. It
never overrides self-containment: verification can be overridden by the
maintainer, but self-containment is a property of the artifact and must be
repaired by finishing adoption.

A merged PR or release tag is **not** a precondition, and the archive never waits for one.

If any check fails, stop and report the offending Task, report, source, or link
and the adoption step that fixes a self-containment failure — the Spec stays
active.

To archive despite failed, missing or otherwise ineligible QA, use the explicit
approval and reason with the archive command:

```bash
roundfix archive <slug> --qa-override --approval <source> --reason <text>
```

The override still requires every non-QA Task to be `completed`. It accepts a
failed or pending QA Task regardless of the newest report's verdict and refuses
only when every Task is completed and the newest report qualifies, because the
same Spec can archive normally. It records the approval source, reason, observed
QA outcome and archived revision without changing the QA Task or report. When
the QA Task is not completed, it also records `qa_override_qa_task_status`. An
unreadable report's recorded outcome names its path relative to the Spec folder,
never the machine's absolute path.

## Steps

For a QA Archive Override, perform the archive only through
`roundfix archive <slug> --qa-override --approval <source> --reason <text>`.
For a normal archive, run `roundfix archive <slug>`; it verifies the
preconditions, writes the Archive Record and removes the Spec folder. Never
hand-edit the Archive Record; the command owns its refusals, provenance, and
archive lifecycle. The Delivery Queue writes the commit subject `docs: archive
<slug>` (Conventional Commits), or, when that would pass 100 characters,
`docs: archive the <spec number> Spec` (`docs: archive the Spec` for a slug
without a Spec number). Do not push unless asked.

Report the new path and anything carried over (open follow-ups from task
`## Result` sections belong in new specs, not in the archive). When the work
isn't merged yet, suggest opening the PR via `github-pr-workflow`; opening the
PR is the user's call.

## Archive Record

An archive writes `<slug>.md` under the resolved archive root. The record names
the QA Report and verdict, its disposition, and `source_revision`, the Git
revision where the removed Spec folder remains recoverable. Before archiving,
use `roundfix archive <slug> --plan` to see the cut and Archive Advice for
candidate files. Promote durable knowledge with repeatable
`--promote <path>`; promotion copies that file to `docs/references/` in the
archive change.

## Unarchive

Rare, explicit, reversed: Restore the Spec folder from Git at the record's `source_revision` and delete the record. Reopening usually means new work — prefer a fresh spec that references the archived one.

## Anti-patterns

- Archiving with failing or missing QA silently — the override must be the user's word, on the record.
- Blocking the archive on a merged PR or release tag — completion (tasks + QA) is the gate; publishing is a separate, user-driven step.
- Leaving a completed spec in the active set "until the PR merges" — the active folder is for live work only.
- Editing an archived spec — it is a record; new requirements get a new spec that links back.
- Deleting instead of archiving — the graveyard is where "didn't we already try this?" gets answered.
