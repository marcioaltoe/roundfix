# What marks each commit of an Item Branch, and what an Author Subagent's command can read (2026-11-17)

Recorded by the author of the PRD of Spec 0404 on the Mac mini on 2026-11-17, between about 15:10 and 15:50 BRT, in a clone of this repository at `74d0e3a9`. Every command was read-only. Section 1 is a reading of the source by the Research Subagent, which the author checked line by line with a search of the same files; sections 2 and 3 are the author's own measurements.

## 1. The commits Roundfix makes on an Item Branch

A Delivery Queue item moves through the stages `queued`, `running`, `reviewing`, `archiving`, `gating`, `publishing`, `checking`, `merging` and `merged`, with `parked` beside them (`internal/store/delivery.go`).

| commit | who makes it | subject | what marks it |
| --- | --- | --- | --- |
| Task commit, a Corrective Task's included | the Daemon | `<type>: <title>` | the trailers `Roundfix-Spec` and `Roundfix-Task` (`internal/daemon/task_engine.go`, line 2924) |
| QA Report commit | the Daemon | `docs: qa report for <slug> (<verdict>)` | the trailer `Roundfix-Spec` and no `Roundfix-Task`; readers match the whole message (`internal/worktree/worktree.go`, line 887) |
| archive commit | the Delivery Queue | `docs: archive <slug>`, or `docs: archive the <NNNN> Spec` past 100 characters | no trailer; the queue accepts an archive commit by its change (`internal/cli/deliver_workflow.go`, line 1386) |
| Pull Request Record Commit | the Delivery Queue | `docs: record pull request #<n> in the <NNNN> Archive Record` | no trailer; accepted by proof of its change (`internal/cli/deliver_workflow.go`, line 1374) |
| merge of the default branch with regenerated derived paths | the Delivery Queue | `chore: merge default branch and regenerate derived paths` | the trailer `Roundfix-Delivery: derived-merge` (`internal/cli/deliver_workflow.go`, line 2864) |
| a correction by the Supervisor or a person | whoever commits by hand | any | nothing |

- The pre-PR review writes no commit: its record, the per-head copies and the ledger `pre-pr-review-dispositions.jsonl` live in the Artifact Directory.
- A queue item stores the heads it accepted as candidates and its merge commit, and no other commit. A Delivery Retry that accepts a moved head appends it to the item's candidates and returns the item to `reviewing`; for an item parked `corrective-spec-required` it first asks for the correction proof of ADR-0272. It records nothing else about the new commits (`internal/delivery/engine.go`, lines 679 to 788).
- No code reads a commit's author or committer to tell Roundfix's commits from a person's.
- A new Delivery Queue replaces a terminal one with its items (`internal/store/delivery.go`, line 172), and Merge Evidence releases an Item Branch once no live item names it.
- Run Retention never touches the Delivery Queue tables, the review records or the disposition ledger (ADR-0255; ADR-0272 says of the per-head records "Nothing reclaims them yet").

Read against the commits of pull request #892 in [`2026-11-17-phase-record-and-causes-view-measurement.md`](2026-11-17-phase-record-and-causes-view-measurement.md): the three Task commits, the QA Report commit, the archive commit and the Pull Request Record Commit each carry one of the marks above, and the three commits the Supervisor made, `f7a132d4`, `d23cf4ba` and `a8b702c0`, carry none.

## 2. What a command run by the PRD's Author Subagent reads

The Author Subagent of this PRD ran in Claude Code; its `AI_AGENT` variable held `claude-code_2-1-325_agent`. Only variable names were listed, and two values were printed.

- `CLAUDE_CODE_SESSION_ID` was set, and `roundfix phase start prd` recorded the entry under it.
- `CLAUDE_EFFORT` held `xhigh`, the effort this repository's Project Config declared for the `prd` phase at `74d0e3a9`.
- No variable's value held `fable`, `opus`, `sonnet` or `haiku`.
- The subagent's own instructions state its model by name and by exact identifier.

This repeats, in a second session, what the Idea's author measured in the first.

## 3. The Governed Paths the design is expected to change

A `go test -overlay` probe in `internal/speccheck`, which wrote nothing to the repository, ran `speccheck.GovernedPath` over 64 paths: the skill files, guides, records, source files and documents this PRD expects the Tasks to change, this Spec's planning files, the three new ADRs, ADR-0407 and controls. `git status --porcelain` printed nothing afterwards.

- Governed: `.roundfixrc.yml`; every file under `.agents/skills/`; `skills/<name>/SKILL.md`; `docs/references/coverage-record.json`; `internal/cli/cli_test.go`; `internal/docscontract/publicdocs_test.go`; and the controls `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.
- Not governed: the mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json`, the generated subagents under `.claude/agents/`, `skills-lock.json`, `CONTEXT.md`, the ADRs, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, the user guide and every Go source file probed.

## Limits

- Section 1 reads one revision. The TechSpec confirms it and lists every commit the Delivery Queue or a Run can make before it defines the commit that is a Correction Sign.
- Section 2 is one session of one harness version, on one machine.
- Section 3 covers the paths a PRD can name. The TechSpec and the Tasks stage measure the paths the Tasks declare.
