---
status: accepted
knowledge: true
created_at: 2026-11-19T15:00:00Z
updated_at: 2026-11-19T15:00:00Z
deprecated_at: null
superseded_by: null
---

# A repository publishes a Standard Export List the Baseline keeps, and varies it only in its own lines

A repository that has the Secondbrain turned on publishes what
`.secondbrain-export` lists, and the Baseline owns part of that file: the
Standard Export List, between two marker lines, which names knowledge and
never process. The repository keeps its own lines outside the markers, each
with its reason, and no Baseline operation changes them.

## Context

Until this decision no Baseline module named `.secondbrain-export`, and no
code of Roundfix read or wrote it apart from one test of this repository's
own file. Each project copied a template that lists `docs/` whole and
`AGENTS.md`. A count on 2026-11-18 found 3,534 documentation files in fifteen
project mirrors, 1,509 of them process: copies of the Baseline's guides,
Archive Records, Backlog Entries, Spec folders, Findings and handoffs. The
same day the project that owns the Secondbrain began to filter process out of
what it sends on to the search index, which is a second list that no
repository sees and that `roundfix doctor` cannot check.

The maintainer decided on 2026-11-18 that Roundfix defines, through the
Baseline, what each project exports, and that "no other project creates these
rules on its own". On how a project may differ he chose "Padrão + variação
declarada (Recommended)", and on `docs/design/` he chose "Sim, na lista padrão
(Recommended)".

Two published conventions bear on the form. A tool that keeps a block inside a
file a person also edits finds its block by a pair of marker lines: the
Ansible `blockinfile` documentation says that a marker which does not tell
the begin line from the end line "may result in the block being repeatedly
inserted on subsequent playbook runs". And a curated list for
machine readers kept at the source is an established convention: "Agents are
best served by concise, expert-level information gathered in a single,
accessible location" (the llms.txt proposal).

## Decision

**The file has two parts.** The Standard Export List is every line between
`# roundfix:standard-export:begin` and `# roundfix:standard-export:end`. The
Export Variation is every line outside them.

**The Standard Export List names knowledge.** It holds, one path per line and
in this order: `README.md`; the glossary files of the repository's declared
layout; `docs/user-guide/`, `docs/references/`, `docs/design/` and
`docs/lessons/`, each listed whether or not the folder exists yet; and the
path of each ADR whose Knowledge Mark is true
([ADR-0439](0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md)).
It is a function of the repository's files and of nothing else, so two runs
on the same tree write the same lines.

**Process is never published.** `AGENTS.md`, `CLAUDE.md`, `docs/agents/`, the
Spec Root, `docs/backlog/`, `docs/findings/`, `docs/handoffs/` and the History
Root are in neither part.

**The Baseline keeps the list, and only the list.** `roundfix baseline update`
rewrites the lines between the markers and no other line. The one edit a
person or an Agent makes there is the path of an ADR written with the mark,
added in order, which is what the update would write.

**The first update converts the file once.** A file without the markers is
read as the list a project copied: a line that names a standard item or a
process path, and an exclusion that lies under a process path, are removed;
every other line stays, in its order, below the new block. The plan names
each path that stops being exported before anything is written.

**A repository varies in its own lines.** It may add a file or a folder of
knowledge the standard lacks, or drop a standard item with an exclusion line.
Each line carries its reason in a comment directly above it. The Doctor
Command reports a line without a reason, a process path, and a Standard
Export List that differs from what the update would write; it changes
nothing, and none of the three fails it.

**A missing or doubled marker is refused, never guessed.** With one marker
missing or either one repeated, the update writes no line of the file, says
which marker is wrong and leaves the rest of its plan to apply.

**A repository with the Secondbrain turned off gets no list.** The update
writes no file there and the Doctor Command skips the check.

**This repository publishes less than the standard.** Its own lines drop the
glossary, `docs/references/`, `docs/design/` and every ADR, by the
maintainer's answer "Sim, mais docs/lessons/ (Recommended)": it publishes
`README.md`, `docs/user-guide/` and `docs/lessons/`.

## Considered Options

- **Keep the filter where the mirrors are read.** Rejected by the
  maintainer's decision: the list would stay outside the repository, where no
  check reads it and no project sees it.
- **A list the repository edits whole, checked against the standard.**
  Rejected: without a block of its own the Baseline cannot tell its old lines
  from the repository's, so it could neither add a standard item later nor
  repair a list.
- **A second file for the repository's lines.** Rejected: the sync belongs to
  another project and reads one file.
- **The block recorded as a managed artifact of the Setup Manifest.**
  Rejected: its content follows the repository's ADRs, so the recorded digest
  would go stale with each marked ADR, and the check every settlement makes of
  the Setup Manifest would stop a Run for a line of documentation.
- **List only the standard folders that exist.** Rejected: a folder created
  later, such as `docs/lessons/` with its first lesson, would stay unpublished
  until someone ran the update. The sync skips a listed path that does not
  exist and says so in one line.
- **`docs/adr/` whole in the standard list.** Rejected by the maintainer, who
  chose "Só as marcadas".

## Consequences

- A mirror shrinks at the first push to the default branch that follows a
  repository's update. A page elsewhere that cites a mirror path under a
  process folder loses its source unless it pins the citation first.
- A folder under `docs/` that the layout does not name stops being exported at
  that update, until the repository adds it to its own lines or moves its
  content.
- Until its ADRs are marked, a repository exports no ADR.
- The sync prints one warning for each standard folder a repository does not
  have.
- The Release Plan counts `.secondbrain-export` as maintenance-only, so the
  commit of an update asks for no release.
- This decision changes what [ADR-0248](0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md)
  says of the export: it no longer mirrors `docs/` whole, and the history
  exclusion it kept is gone with the History Root.
