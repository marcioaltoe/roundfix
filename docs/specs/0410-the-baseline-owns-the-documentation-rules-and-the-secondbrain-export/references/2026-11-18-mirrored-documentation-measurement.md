# Mirrored documentation, ADR status and rename cost — measurement of 2026-11-18

The author of the Idea of Spec 0410 took these counts on 2026-11-18, read-only, from two sources: this repository at `d1c43f9ef26c64bf75000c71699e8bdaae784cf6`, and the project mirrors of the Secondbrain at its commit `a2b713574` (2026-11-18, 23:01 BRT). A mirror holds what a project's `.secondbrain-export` published at its last push to `main`, so a file a project does not export is invisible here. Every count is of `.md` files unless it says otherwise.

## 1. What the mirrors hold

The 15 mirrors hold 4,397 files: 3,534 under `docs/`, 50 at the root of a mirror, and 813 elsewhere, almost all of them the skill texts of the skills repository.

| Folder | Files | Kind in the request |
| --- | --- | --- |
| `docs/history/` | 1,155 | process |
| `docs/agents/` | 146 | process |
| `docs/backlog/` | 125 | process |
| `docs/specs/` | 65 | process |
| `docs/findings/` | 16 | process |
| `docs/handoffs/` | 2 | process |
| `docs/adr/` | 1,106 | knowledge, in part |
| `docs/references/` | 774 | knowledge |
| `docs/user-guide/` | 79 | knowledge |
| `docs/design/` | 15 | knowledge |
| other folders and loose files under `docs/` | 51 | not classified |

Process files under `docs/` add up to 1,509, and the four knowledge folders to 1,974. The split is the one the request proposes; no file was judged one by one. Every one of the 15 mirrors holds `AGENTS.md`, and 3 hold `CHANGELOG.md`.

This repository's mirror holds 364 ADRs, 23 files under `docs/references/` and 43 under `docs/user-guide/`.

The request quotes the Indexus's own count of the same day: 4,150 files of 11 projects. The two counts cover different sets of mirrors and agree on the proportion.

## 2. The status of the ADRs

Thirteen mirrors hold ADRs, 1,106 files in all.

| Status line | Files | Where |
| --- | --- | --- |
| `status: accepted` | 1,083 | all thirteen |
| `status: superseded`, with `superseded_by` filled | 4 | 3 in the Indexus since 2026-11-16, 1 in another project since 2026-08-27 |
| no `status:` line | 19 | 15 in one project, 3 in a second, 1 in a third |

The docs-layout guide of both repositories that keep a superseded ADR carries the clause that deletes it.

In this repository all 364 ADRs carry `status: accepted`, and none is superseded, deprecated, rejected or proposed. An ADR changed in part by a later one keeps that status and says so in a line of its body:

- 64 ADRs carry an "Amended by" or a "Narrowed by" line, 83 lines in all;
- 73 ADRs carry one of "Amended by", "Narrowed by", "Extended by" or "Superseded in part";
- in 1 of the 64 the first such line is within the first 12 lines of the file, and in 37 of the 64 it sits in the last quarter of the file.

Two examples: the mark of `docs/adr/0092-triage-intent-lives-in-a-typed-backlog.md` is line 40 of 40, and `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md` carries five, on lines 104 to 112.

## 3. What a folder rename would touch

Files of this repository that name each path literally:

| Path | Go files, not tests | Go test files | Baseline modules and profiles | Files under `.agents/skills/` |
| --- | --- | --- | --- | --- |
| `docs/agents` | 18 | 104 | 18 | 17 |
| `docs/references` | 16 | 49 | 2 | 12 |
| `docs/backlog` | 10 | 40 | 2 | 4 |
| `docs/findings` | 9 | 30 | 2 | 3 |
| `docs/handoffs` | 2 | 7 | 1 | 1 |
| `docs/design` | 0 | 2 | 1 | 0 |

The count is `grep -rl` of the path over `internal/` and `cmd/` for the Go files, over `internal/baseline/assets/modules/` and `internal/baseline/assets/profiles/` for the Baseline, and over `.agents/skills/`. It leaves out the derived Baseline files, which repeat every path, and every adopter's own files. Only the Spec Root can be configured (`internal/config/config.go:51`).

## 4. Layout the directory clause does not name

Seen in the mirrors:

- reference files at the root: `DESIGN.md` in 4 of the 15, `DATABASE.md` in 1, `CONTEXT-MAP.md` in 2, and one each of `RELEASE.md` and `CONTRIBUTING.md`;
- folders under `docs/` that the clause does not name: `reference`, `_refacs`, `prompts` and `schema`, each in one project;
- `docs/_inbox/` in 4 projects;
- 3 loose files directly under `docs/`, in 2 projects.

The Baseline itself places `DESIGN.md` at the root: `internal/baseline/assets/modules/frontend.json:297` declares the path and `:326` tells an Agent to read it before UI work.

## 5. The export files

- No Go code of this repository reads or writes `.secondbrain-export`. The only reader is the test `internal/docscontract/secondbrain_export_test.go`, which checks this repository's own file.
- No Baseline module names the file.
- The three export files read on the author's machine each list `docs/` whole. The template the projects copy, which the Secondbrain keeps, lists `CONTEXT.md`, `docs/`, `README.md`, `AGENTS.md` and `DESIGN.md`.
- The mirror is replaced on every sync, so a path that leaves an export list leaves the mirror at the project's next push to `main`.

## 6. Lessons

- `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md:11` records the one earlier count: of 383 files in `docs/history` on 2026-10-14, 7 held a lesson that no ADR states.
- No mirror holds a folder for lessons.
- The Archive Record keeps no count of retries, Corrections or QA rounds (`internal/spec/archive_record.go:43`).
