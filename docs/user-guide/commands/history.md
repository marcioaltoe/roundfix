# `history`

`roundfix history sanitize` converts each Legacy Archive Folder to one Archive
Record and removes retired Findings, Rollups, Backlog Entries, ADRs, Review
Artifacts and handoffs. A Reduced History Entry is pending like a full file.
The `superseded` unit removes uncited superseded Archive Records without
regeneration grants. Delivered and cancelled records stay. A folder without
QA, an override or supersession receives `no-qa`.

```bash
roundfix history sanitize
roundfix history sanitize --format json
roundfix history sanitize --batch 40 --advise
roundfix history sanitize --apply --batch 40
roundfix history sanitize --apply --batch 1 --promote docs/history/specs/example/references/lesson.md
```

Units are ordered by Legacy Archive Folder slug, then `findings`, `backlog`,
`adr`, `reviews`, `handoffs` and `superseded`, then `record-links` and
`live-links`. Legacy Archive Folders use Lenient Legacy Reading: the
Task Graph tolerates projection rows outside the graph and retired Task types,
and the plan names each tolerance as `tolerates <folder>: <tolerance>`. Active
Specs keep strict reading. Each kind counts as one unit while files remain
pending. `--batch <n>` selects the next positive number of units it can convert;
Refused Units do not count toward `n`. Without it the plan shows all pending
units. A later batch starts from the next remaining unit. No new Reduced History Entry is written.

Lenient Legacy Reading also accepts a legacy list of maps in `unproven`. Each
map becomes one stable text line in the Archive Record; active Specs continue to
require the string-only list. A legacy `qa_override: true` with a missing, empty
or blank approval or reason records `legacy: not recorded` in that field;
recorded values and the other override fields are preserved. When
`references/_index.md` is absent or lists no path, `sources` lists, sorted by
file name, the regular Markdown files directly under `references/` whose front
matter identifies a Backlog Entry (`type`, `status`, `created`) or a Finding
(`status`, `created_at`). An empty PRD `created` takes the TechSpec's valid
`created` date, then, in `history sanitize` and `baseline update`, the author
date of the oldest commit reachable from `HEAD` that added the PRD under the
Spec Root or archive root. The folder reader stops at the TechSpec date. The
summary is read from the PRD's title paragraph, or from the first paragraph
after its first `## ` heading when that paragraph is empty; block quote markers
are removed before its lines are joined. Legacy Link Rebase rewrites the
relative destinations in the outcome and in each promoted copy: a destination
inside the folder points to its promoted copy or becomes a Git Locator, and
every other destination keeps the same repository path. A Refused Unit is listed as `refused <unit>: <reason>`,
with the reason on one line, and remains untouched.

`roundfix baseline update` plans these same units at once as Pending History and
binds the selected units into its Baseline Plan Digest. Use
`roundfix baseline update --no-history` when the operator wants the History
Sanitize Command's reviewed-batch path. The History Sanitize Command remains the
path for applying `--batch <n>` batches.

The plan writes nothing under the repository. Its first line is:

```text
history sanitize plan: <u> unit(s) pending; <f> file(s) (<b> bytes) leave docs/history
```

Each folder line reports removed file count and bytes, record path and size,
disposition and delivery (`delivery <12-hex> #<pr>`, `delivery <12-hex>` or
`delivery unknown`). Candidate lines list files outside core artifacts and
`qa/evidence/`, using the archive plan's grouping. Removal units print `<unit>: removes <n> file(s) (<b> bytes)`.
`cites <path>:<line> names <target>` reports tracked Markdown, including Archive
Records. A citation in an Archive Record, a Live Document or `docs/user-guide/`
appends `; rewrites to <replacement>`; apply rewrites it with the Citation
Rewrite. A removed link becomes its label followed by the Git Locator
`` `<path>@<commit>` ``; a prose path becomes the locator. Other citations are
listed and left unchanged. A Reduced History Entry uses its original provenance
path and commit when that commit qualifies. The rewrite changes only the
planned citing files and skips files removed by the same plan.

The family units include `docs/history/<kind>/`, the legacy trees
`docs/findings/_archived/` and `_archived/{findings,backlog,adr}/`, terminal
Findings and Backlog Entries directly in their active directories, retired
ADRs directly in `docs/adr/`, and finished Review Artifacts under the Spec
Root's review roots. A superseded ADR naming an accepted successor leaves even
when the successor does not name it yet: the plan appends its lineage line.
An ADR without an accepted successor is refused by its repository-relative
path while every other retired ADR leaves. Unknown or missing lifecycle
statuses in active directories remain untouched. Active directory roots stay.

The `superseded` citation check reads every tracked Markdown file surviving the
plan, outside `.agents/skills/`, `skills/` and `testdata` directories. Relative
links, bounded record paths and bounded slugs all count, including front matter
values. Text already inside a Git Locator does not count.

After removal units, `record-links` repairs each surviving Archive Record's
relative links to missing files and to itself. It tries the record-relative
path, then the record's original `source` folder; the first reachable commit
holding the target wins. A self-link becomes the Git Locator of the source
folder. When neither path has a qualifying commit, the record-relative path
and its pull request are used. Code spans, fences, URLs and Git Locators are
not relative links. A record removed by the same plan is skipped.


A superseded ADR's accepted successor receives a missing lineage line exactly
as `Supersedes ADR-<number> (<date>); Git keeps its text at <locator>.` The date
is HEAD's committer date in UTC (`YYYY-MM-DD`); the locator is the removed
path's Git Locator. Apply appends the line after a blank line, preserves the
successor's existing bytes and never adds the same line twice. Several lines
for one successor follow removal-path order. Each line is a History Citation:
`cites <successor>:<n> names <removed>; rewrites to <line>` in history text,
and under `citations` in JSON. Baseline text prefixes the citation with `- `;
its Plan Digest binds the exact line, and an uncommitted successor refuses
the ADR unit. Successors are cited, never listed as removal paths.

A superseded ADR with no accepted successor stays with this reason:
`superseded ADR has no accepted successor; name the accepted ADR that replaces it in superseded_by, or set its status to deprecated`.
History text prints `refused <path>: <reason>`; baseline text prefixes it with
`- `. JSON lists `{"unit":"<path>","reason":"<reason>"}` under `refused`.
The refusal names that ADR's repository-relative path, leaves the other
retired ADRs eligible and does not change the command's exit code.

After `record-links`, `live-links` repairs missing relative links in Live
Documents and `docs/user-guide/` pages, including the configured Spec Root.
It writes the Git Locator printed by `SC-HISTORY-RETAINED`, changes only link
spans and skips Archive Records, removed pages, self-links, code spans,
fences, URLs and existing Git Locators. A second plan after apply has no
`live-links` unit.

The repair line is `live-links: rewrites <l> link(s) in <n> file(s)`.
JSON lists unit `live-links`, action `rewrite`, sorted page `paths`, `files`
and `links`; its citations and replacement text join the Plan Digest.

The repair line is `record-links: rewrites <l> link(s) in <n> Archive Record(s)`.
Rewritten records count toward units, while the file and byte totals count only
files that leave the tree. Apply reports rewrites alongside removal citations.

When a unit cannot be converted, the plan continues and lists
`refused <unit>: <reason>`. Its first line adds `; <r> unit(s) refused`. A failed
newest QA Report without an override is recorded as `failed-qa`, beside `no-qa`
for a folder with no QA Report, override or supersession.

`--advise` requires `--batch` and cannot accompany `--apply`. It adds Jev's advice
for candidate files in selected folders, using the archive judge's command-environment key
variables, configured monthly ceiling and log under Roundfix Home. Missing keys,
reached ceilings and service failures print `no advice (<reason>)` and send no
request when skipped. Advice never changes selection or the command's success.
The judge log is the only advice write and stays outside the repository.

When `ROUNDFIX_RUN_ID` is non-empty in the command's environment, every
candidate has `no advice (inside Run <run-id>, named by ROUNDFIX_RUN_ID)`.
Archive Advice selects no transport, reads no key or Judge Log, writes no Judge
Log and sends no request, even with a judge key present.
Every Agent Session and Verification command of a Run receives this marker
with the Run's ID. Roundfix also withholds the Judge Key Variables, including
`TYPESAFE_API_KEY`, from those children; the marker keeps the judge silent even
if a runtime restores a key from the user's shell startup.

`--apply` requires `--batch`. It validates the entire batch before writing and
uses `HEAD` for records' `source_revision`.
`--promote` is repeatable, requires `--apply`, and names repository-relative files
inside selected folders. It writes each file to
`docs/references/<basename>` after applying Legacy Link Rebase to its relative
destinations and path text, and lists the destination in the record before
removing the folder.

The confirmation is one line:

```text
history sanitize applied <n> unit(s): wrote <r> Archive Record(s), removed <d> file(s) (<b> bytes) kept in Git at <12-hex> and tag history-full; rewrote <c> citation(s) in <m> file(s); promoted <p> file(s) to docs/references/; <remaining> unit(s) remain
```

Apply prints each `refused <unit>: <reason>` line before the confirmation. When
there are refusals, the confirmation includes `; <r> unit(s) refused` before the
remaining-unit suffix. If every examined unit is refused, apply writes nothing,
prints the refusal lines, refuses with exit `2`, and reports the preflight reason
on stderr.

Refusals use Preflight Validation on stderr and exit `2`: invalid or missing
batch values, unknown flags, extra arguments or unknown subcommands; advice
without a batch or with apply; promotion without apply; an external Spec Root
or no Git repository; a dirty tree including staged and untracked files; a
missing, lightweight or non-ancestor `history-full` tag; a removed path neither the tag nor a commit the default branch reaches holds; and any Git failure while reading delivery history.
Promotions outside the batch or inside a Refused Unit, duplicate promotion
destinations and promotion errors are also refusals before writing. An individual
unit conversion error is listed as a Refused Unit instead of stopping the plan.
non-regular files, core artifacts, duplicate destination basenames, existing
destinations and unsafe paths are refused before writing. These whole-command
refusals are the cases covered by the TechSpec's Invariants 9 to 13 and Surface
Transcripts 1 to 3. `history` alone and
`--help` print usage and exit `0`.

A write failure exits `1` and names the unit. Restore the partial tree with
`git restore` and `git clean` before retrying. When nothing remains, the plan
prints `history sanitize plan: 0 unit(s) pending; nothing pending`; apply exits
`0` without writing. The command never commits, tags, pushes or opens a Pull
Request.

## Operator batch procedure

Run this after the implementation merges, outside a Task.

1. Build `bin/roundfix` from `main`. On clean `main`, create and push the
   annotated History Full Tag on the last commit before the first batch:

   ```bash
   git tag -a history-full -m "docs/history before the sanitize batches (ADR-0248)"
   git push origin history-full
   ```

   Apply requires that tag to be an ancestor of `HEAD`. Each removed path must
   be held by the tag or a commit the default branch reaches; paths written
   after the tag can therefore be removed. Otherwise apply refuses by name:
   `history-full does not hold <path> and no commit the default branch reaches holds it`. Recover a full file with
   `git show history-full:docs/history/<path>` or use its record's revision.

2. Save `roundfix history sanitize`'s plan. Request `--batch 40 --advise` per
   batch for Jev's advice on each text candidate; it uses the judge key and
   skips when the monthly ceiling is reached.
3. For each batch, create `chore/history-sanitize-<k>` from `main`, run
   `roundfix history sanitize --apply --batch 40`, stage it with
   `git add -A docs` (the repository contract inspects tracked files), then
   run `make verify` and `make verify-docs`. Commit `chore: sanitize history batch <k>`, open one
   Pull Request, merge it and update `main`. Revert a batch by reverting its
   commit.
4. Promote lessons in the batch that holds them, with `--promote <path>`.
   `--promote` needs a clean tree like `--apply`; a lesson a batch missed is
   copied from `history-full` into `docs/references/` and named in its record's
   `promoted` field. This repository completed its sanitize on 2026-10-07 in
   six batches and promoted, among others,
   `docs/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`
   and `docs/references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md`.

   Send cross-project lessons to the Secondbrain inbox following
   [Secondbrain guidance](../../agents/secondbrain.md). The 0035 skill analysis
   and 0079 pilot report are already there.
5. Six batches of 40 cover the measured 223 folders; the sixth includes the
   four kind units. In that last batch the operator also drops the
   `docs/history` exclusions from `.secondbrain-export`, as required by the
   export contract. The sanitize command itself does not edit that file.
6. Confirm the final plan reports `0 unit(s) pending`.
