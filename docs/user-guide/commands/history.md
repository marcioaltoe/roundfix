# `history`

`roundfix history sanitize` plans the existing history. It converts a Legacy
Archive Folder to one Archive Record, reduces retired Findings and Backlog
Entries to their front matter, title, first paragraph and Git provenance, and
removes retired Review Artifacts and handoffs. Retired ADRs and existing records
stay whole. A folder without QA, an override or supersession receives `no-qa`;
that disposition records absent QA rather than inventing a passing verdict.

```bash
roundfix history sanitize
roundfix history sanitize --batch 40 --advise
roundfix history sanitize --apply --batch 40
roundfix history sanitize --apply --batch 1 --promote docs/history/specs/example/references/lesson.md
```

Units are ordered by Legacy Archive Folder slug, then `findings`, `backlog`,
`reviews` and `handoffs`. Each kind counts as one unit while files remain pending.
`--batch <n>` selects the next positive number of units; without it the plan shows
all pending units. A later batch starts from the next remaining unit. Already
reduced entries are skipped.

The plan writes nothing under the repository. Its first line is:

```text
history sanitize plan: <u> unit(s) pending; <f> file(s) (<b> bytes) leave docs/history
```

Each folder line reports removed file count and bytes, record path and size,
disposition and delivery (`delivery <12-hex> #<pr>`, `delivery <12-hex>` or
`delivery unknown`). Candidate lines list files outside core artifacts and
`qa/evidence/`, using the archive plan's grouping. Kind lines report reductions
or removals. `cites <path>:<line> names <target>` reports tracked Markdown outside
the History Root that cites a removed path; citations never refuse the batch.

`--advise` requires `--batch` and cannot accompany `--apply`. It adds Jev's advice
for candidate files in selected folders, using the archive judge's process key
variables, configured monthly ceiling and log under Roundfix Home. Missing keys,
reached ceilings and service failures print `no advice (<reason>)` and send no
request when skipped. Advice never changes selection or the command's success.
The judge log is the only advice write and stays outside the repository.

`--apply` requires `--batch`. It validates the entire batch before writing and
uses `HEAD` for records' `source_revision` and reduced entries' full-text revision.
`--promote` is repeatable, requires `--apply`, and names repository-relative files
inside selected folders. It copies each file byte-identically to
`docs/references/<basename>` and lists the destination in the record before
removing the folder.

The confirmation is one line:

```text
history sanitize applied <n> unit(s): wrote <r> Archive Record(s), reduced <f> file(s), removed <d> file(s) (<b> bytes) kept in Git at <12-hex> and tag history-full; promoted <p> file(s) to docs/references/; <remaining> unit(s) remain
```

Refusals use Preflight Validation on stderr and exit `2`: invalid or missing
batch values, unknown flags, extra arguments or unknown subcommands; advice
without a batch or with apply; promotion without apply; an external Spec Root
or no Git repository; a dirty tree including staged and untracked files; a
missing, lightweight or non-ancestor `history-full` tag; a tag missing any path
removed or rewritten; and any batch planning error. Promotions outside the batch,
non-regular files, core artifacts, duplicate destination basenames, existing
destinations and unsafe paths are refused before writing. `history` alone and
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

   Apply requires that tag to be an ancestor of `HEAD` and to hold every path
   the batch removes or rewrites. Recover a full file with
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
