---
status: accepted
created_at: 2026-10-20T09:00:00Z
updated_at: 2026-10-20T09:00:00Z
deprecated_at: null
superseded_by: null
---

# The history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator

On 2026-10-19 an adopter's `roundfix baseline update` (0.91.0) removed 27
retired history files, listed 45 History Citations and rewrote 30 of them.
It left 18 citations pointing at removed files in ADRs, agent guides, READMEs
and a Makefile comment. It also named a date-named ADR `ADR-2026`, dated with
the day of the update, and turned links to a file a pending Task will create,
and links that were already broken, into `(not reachable from any branch or
tag)`. `roundfix spec check --strict` then exited 0. The maintainer decided on
2026-10-20 to fix all four and to have the check report what is left.

## Decision

**Every listed citation is rewritten.** The history plan looks for History
Citations in every Citation Source: a tracked regular text file that holds no
NUL byte in its first 8,000 bytes, outside the skill trees `.agents/skills/`
and `skills/` and outside every `testdata` directory. A citation is a
repository-relative path, a relative link destination, or a prose path that
starts with `./` or `../`, resolved against the citing file's directory. The
plan rewrites each citation it lists to the Git Locator of the removed file,
or to the surviving record or promoted copy. ADR-0278's narrower set (Archive
Records, Live Documents, the user guide and the root instruction carriers) is
dropped. In Markdown a link becomes its label followed by the locator in
code; a path inside code gets no second pair of delimiters. In any other
source only path text changes, and the locator is written as plain
`<path>@<commit>`.

**A file a pending Task creates is not missing.** The `live-links` repair and
`SC-HISTORY-RETAINED` leave alone a relative link whose target a Task of an
active Spec declares in `creates:` and has not completed.

**A broken link is resolved by its file name.** When no reachable commit holds
the path a missing relative link names, the repair and the check look for the
one path with the same file name among the files the plan removes and the
Removed History Paths. A single match gives its Git Locator. With none or
several, the path keeps today's `(not reachable from any branch or tag)` form.

**A lineage line names the ADR the way the repository does.** A retired ADR
named `NNNN-<slug>.md` is `ADR-NNNN`. Any other name, such as
`YYYY-MM-DD-<slug>.md`, is the file name without `.md`, the form
`superseded_by` already uses to name such an ADR. The numeric forms `ADR-NNNN` and
`NNNN` select only numbered ADRs. The date is the day the ADR was superseded:
its `updated_at` in UTC, else the committer date in UTC of the newest commit
that changed it, else the committer date of `HEAD`.

**The check reports what is left.** Once the retirement rule is present,
`SC-HISTORY-RETAINED` also reports, as an error, each citation in a Citation
Source of a Removed History Path that is not written as a Git Locator. The
fix line prints the locator to use. A Removed History Path is a path under
`docs/history/findings/`, `docs/history/backlog/` or `docs/history/adr/` that
`HEAD` or its history holds and the working tree lacks. Files under
`docs/history/` are exempt, Archive Records included. So are `creates:`
targets of pending Tasks and a link the check already reports as missing.

## Considered Options

- Reporting the citations and leaving them unrewritten was the inbox's
  second option; the maintainer chose to rewrite every listed citation and to
  report as well.
- Markdown-only sources leave the Makefile comment the adopter found.
- Including the skill trees would edit upstream-managed skills, which the
  Baseline forbids, and a sync would undo the edit. Including `testdata`
  would rewrite fixtures that cite such paths on purpose.
- Resolving every relative-looking word, such as `findings/x.md`, would
  rewrite prose; only `./` and `../` say that a path is relative.
- A new rule ID was rejected. The finding shares `SC-HISTORY-RETAINED`'s gate,
  its scan and its fix line, and it reports the same thing: history the rule
  removed and something still references.
- Choosing among several files with the same name would guess; the form that
  says the path is not reachable stays honest.
- Keeping `HEAD`'s date dates the line with the day the update ran, not the
  day of the decision the line records.

## Consequences

- An adopter's next update rewrites citations in ADRs, agent guides outside
  setup markers, READMEs and Makefiles, and lists each one first, bound into
  the Plan Digest.
- The text and JSON formats keep their fields; more citations carry
  `rewrites to`, and a lineage line of a numbered ADR now carries its
  `updated_at` date.
- A repository whose earlier update left such citations fails
  `spec check --strict` until it applies the printed locators.

This amends ADR-0282's lineage line and extends ADR-0278's rewrite and
check; the archive's relink of ADR-0260 is unchanged.
