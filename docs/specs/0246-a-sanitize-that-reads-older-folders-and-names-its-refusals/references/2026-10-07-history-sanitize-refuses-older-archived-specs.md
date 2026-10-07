---
type: fix
status: promoted
created: 2026-10-07
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
---

# History sanitize refuses older archived Specs and aborts the whole plan

## Problem

In Oraculum, with Roundfix 0.55.0, `roundfix history sanitize` fails Preflight
at the first Legacy Archive Folder and prints no plan. The cause is that
`BuildArchiveRecord` validates each legacy folder against today's schema.
Four archived folders failed it, each for its own reason:

1. **0001.** The `_tasks.md` projection table names `task_12` and `task_13`,
   which were moved to 0004 and commented out of the graph. The refusal reads
   `projection table row names unknown Task "task_12"`.
2. **0004.** The projection table lists `task_03` to `task_05`. They were left
   out of the graph on purpose.
3. **0005.** `task_01` has type `refactor`, which is no longer allowed.
4. **0036.** It was archived by maintainer decision with a failing QA report
   and no `qa_override` stamp. The sanitize copies the `fail` verdict into the
   record, and `ParseArchiveRecord` refuses it with
   `unknown archive disposition "fail"`.

The adopter cannot repair these folders. The docs layout forbids hand edits
to archived Spec metadata, and `archive --qa-override` works only on an
active Spec. With local patches the plan measures 71 units: 67 folders,
2,849 files and 31.4 MB.

## Direction

- **Read legacy folders leniently.** Tolerate projection rows outside the
  graph and retired Task types. The current schema stays enforced for active
  Specs.
- **Record archive-time failures honestly.** Give a folder archived with a
  failing QA and no override a valid record, for example a disposition such as
  `archived-with-failed-qa` that keeps the report name. Never turn it into a
  pass.
- **Report every refusal.** The plan lists each refused unit with its reason
  and still plans the others. A batch skips refused units, or stops before
  them, but never aborts at the first one.

## Sources

Secondbrain inbox, triaged 2026-10-07:
`inbox/roundfix/_triaged/2026-10-07-history-sanitize-recusa-specs-arquivadas-antigas.md`
(Oraculum).
