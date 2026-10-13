---
spec: 0263-archive-and-cancel-keep-what-git-ignores
status: active
created: 2026-10-12
surfaces: [backend, cli, docs]
---

# Archive and cancel keep what Git ignores

`roundfix archive` and `roundfix cancel` remove the Spec folder and say its
bytes stay in Git at the record's `source_revision`. Files Git ignores break
that promise: both commands remove them, count them among the files kept in
Git, and no commit holds them. When the final folder removal fails, both
commands have already written the record and every other output, yet archive
reports "No side effects" with exit 2 and neither command says how to finish
([the Backlog Entry](references/2026-10-12-archive-and-cancel-remove-ignored-files-in-the-spec-folder.md)).
This bug fix makes both commands refuse while Git ignores a file under the
Spec folder, list each such file and write nothing; and makes a failed folder
removal keep what was written and print the exact command that finishes it.

The maintainer decided on 2026-10-12, finally: refuse, list each ignored file
from `git status` with `--ignored`, and write nothing; recover forward from a
failed removal and print the finishing command; keep accepting symbolic-link
destination parents, because the external Specs Root of a knowledge workspace
is one. One shared helper serves both commands.

## Prerequisites

None. The other Specs of cycle AF (0262 and 0264) and the unmerged Spec 0258
are authored at the same time. Spec 0258 adds structural `spec check` gaps for
every active Spec; this Spec already carries the sections it reads. No
Verification here depends on another Spec's artifacts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or Archive
  Record field changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; tests run on temporary repositories and homes. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0269 (this Spec) decides the
  refusal and the forward recovery. ADR-0247: "so the recorded revision holds
  exactly what was removed", which this Spec makes true for ignored files.
  ADR-0248: "Existing history is sanitized in batches after a history-full
  tag"; History Sanitize converts only Legacy Archive Folders and is unchanged.
  ADR-0251: "A Legacy Archive Folder is read leniently" and ADR-0254:
  "Baseline update sanitizes pending history in the change it plans"; neither
  reader nor sanitizer changes.
  ADR-0262: "refuses while a Run or a Delivery Queue item is active for the
  Spec"; the Cancel Command gains one more refusal. ADR-0261: "An operator
  records the number with" the archive's `--pull-request`, which still records
  after a failed removal. ADR-0260: "An archive relinks live citations", and
  the relinked documents stay written after a failed removal. ADR-0184: "A
  TechSpec states a command surface as a transcript", so the four outputs are
  Surface Transcripts. ADR-0256: "A release waits for the skills that describe
  a changed surface". ADR-0244: "A Spec declares the domain terms it
  introduces". ADR-0257: "Inside a Run the Daemon is the only full gate".
  ADR-0258: "A Spec is authored against the QA rerun classes". ADR-0259: "A
  test runs in parallel unless it names why it cannot". ADR-0271: "A test bounds by real time only what it asserts"; the new tests give no child, goroutine or Git a real-time allowance they do not assert. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AF grant of 2026-10-12 ("Sim, como no AE", carrying the cycle AE "Concedo"
  for the Governed Paths each Spec declares) and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"); bounded
  files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`, `internal/spec/archive.go`,
  `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`. No Makefile,
  workflow, lint, formatter or module file change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0263-archive-and-cancel-keep-what-git-ignores/_authorization.md`.

## Goals

- No archive or cancellation removes a file that no commit holds.
- A refused archive or cancellation names every file Git ignores under the
  Spec folder and changes no byte.
- A failed folder removal keeps every write, exits 1 and prints the command
  that finishes it.

## User Stories

1. As an operator, I want `roundfix archive` and `roundfix cancel` to refuse
   while Git ignores a file in the Spec folder and to name each one, so that a
   scratch log or evidence file is never deleted without a copy.
2. As an operator whose folder removal failed, I want the command to tell me
   what it already wrote and the exact command that finishes the removal, so
   that I can complete the archive or cancellation instead of untangling a
   half-written tree.

## Core Features

1. Ignored file refusal. Both commands read the Spec folder with
   `git status --porcelain=v1 -z --untracked-files=all --ignored`, in the
   repository that holds the Specs Root, through one shared helper. While Git
   ignores any file under the folder they refuse with exit 2, list each file
   by its repository-relative path, sorted, one by one even inside an ignored
   directory, and write nothing. Uncommitted changes keep their existing
   refusal and are reported first.
2. Forward recovery. When the final folder removal fails, the record, the
   promoted, returned and relinked files stay written, archive still records
   `--pull-request`, and both commands exit 1 with two standard-error lines:
   what was written and the finishing command
   `chmod -R u+w '<folder>' && rm -rf '<folder>'`.
3. Documentation. The help of both commands, their user guides, the Roundfix
   skill's archive reference and the archive-spec skill describe the refusal
   and the recovery; `CONTEXT.md` defines the two new terms.

## Non-Goals / Out of Scope

- Accepting a symbolic-link destination parent stays as it is (maintainer
  decision of 2026-10-12).
- `roundfix archive --plan` output, the Archive Record format and the
  Delivery Queue's park and retry rules do not change.
- Copying, committing or warning about ignored files instead of refusing.
- Rolling back after a failed folder removal.

## Success Metrics

1. Success Metric: with one ignored file at the folder root and two in an
   ignored directory, each command exits 2, lists all three paths, and leaves
   the tree byte-identical, ignored files and Git status included.
2. Success Metric: with a write-protected subdirectory, each command exits 1,
   keeps its record and other writes, prints the finishing command, and running
   that command exactly as printed leaves the tree a successful run would
   leave.
3. Success Metric: the four Surface Transcripts of the TechSpec match the
   built binary line for line.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Git manual for `git status` (https://git-scm.com/docs/git-status, read
  through Exa on 2026-10-12, Git 2.56.0 installed): it displays "paths in the
  working tree that are not tracked by Git (and are not ignored by
  gitignore[5])", and `--ignored` in `traditional` mode shows ignored files
  "unless --untracked-files=all is specified, in which case individual files in
  ignored directories are displayed". This is why the old guard missed ignored
  files and why the new one lists each file.
- The Go standard library, read offline with `go doc os.RemoveAll`: RemoveAll
  "removes everything it can but returns the first error it encounters". This
  is why a failed removal leaves a partly removed folder and why the recovery
  goes forward rather than back.
- The Secondbrain, read during authoring through `wiki/index.md` and two
  `qmd query` searches: the mirror of ADR-0247,
  `projects/roundfix/mirror/docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`,
  confirms the promise this Spec completes, and
  `wiki/sources/digest-roundfix-git-add-f-apaga-arquivo-2026-09-17.md` records
  that Roundfix already stages ignored Task output with `git add -f`, which is
  the commit path the refusal tells the operator to use. No page described
  ignored files lost at archive. The Git manual and Go documentation passages
  are standard references, not research, so no Inbox Entry was captured.

## Glossary

- adds: **Ignored Spec File**
- adds: **Folder Removal Recovery**

## Decisions

- Both commands refuse while Git ignores a file under the Spec folder, list
  each file and write nothing; see ADR-0269.
- A failed folder removal keeps every write, exits 1 and prints
  `chmod -R u+w '<folder>' && rm -rf '<folder>'`; see ADR-0269.
- Symbolic-link destination parents stay accepted; see ADR-0269.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the characterization, the contracts
and the build order.
