---
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
status: active
created: 2026-11-19
surfaces: [backend, cli, docs]
type: fix
---

# The status counts a delivering Spec's Tasks, and a Roundfix update needs no release

Two defects an adopter met in one status report, reported by the maintainer on 2026-11-19 from the Fiscus. `roundfix follow-up --status` showed Spec 0097, running in the Delivery Queue, at 0 of 16 Tasks while 15 were done on its Item Branch: the status view counts every Spec's Tasks in the checkout's Spec Root, where a delivering Spec's plan still holds every Task pending. And the next release needed a classification by hand for four of six commits, three of which changed only files Roundfix itself writes into the adopter or a `CONTEXT.md`: the Release Plan counts only planning evidence as maintenance-only. After this Spec, a Delivery Queue item's Tasks are counted on its Item Branch as Git holds it, and in an adopter a commit confined to the files Roundfix manages, to `CONTEXT.md` files and to planning evidence needs no release, while Roundfix's own classification does not change.

## Sources

The source is the maintainer's report and answers of 2026-11-19, copied verbatim to [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md) and cited as "the maintainer, 2026-11-19": a status capture of the Fiscus, the operator's first measurement and the two answers. The reproduction, the prototype and the measurements on both histories are in [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md), cited as "measured by the operator".

The two defects share no component: one is the status view's Task count, the other the Release Plan's maintenance boundary. The maintainer chose one Spec for both, as for Spec 0407 ("Spec fix 0414 logo após a 0411 (Recommended)" — uma Spec com as duas correções (como a 0407)): both surfaced in one status report of one adopter, a source keeps one owning Spec, and each half alone would be a one-Task Spec paying a whole Run, QA gate, review and release. No open Backlog Entry or unresolved Finding changes the status view's Task count or the Release Plan's boundary.

## Prerequisites

None for Verification. Authored against `145d6c76` (v0.148.1). Spec 0411 delivers before this Spec and Spec 0409 after it; Spec 0411 also changes `CONTEXT.md`, `docs/references/behavior-surfaces.json`, the Roundfix skill's `SKILL.md` in both copies and `skills/testdata/owned-skill-versions.json`, and Spec 0409, a change of the authoring skills, is expected to touch the same records. Each Spec raises the skill version and regenerates the two records on the tree it runs on; this Spec writes "the next version", never a number.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the status view reads the existing Item Branch name and Archive Record of a Delivery Queue item, and the Release Plan reads existing skill names and paths. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; both corrections read local Git objects of the repository. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0447 (new) records the release decision; ADR-0412, ADR-0349 and ADR-0359 constrain it; ADR-0410, ADR-0434, ADR-0426, ADR-0057 and ADR-0071 are listed to record how they bear on the status count; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Spec fix 0414 logo após a 0411 (Recommended)" of 2026-11-19 with the cycle BE grants ("Mesmas do BE"), quoted in the authorization record. It covers the Roundfix skill's `SKILL.md` and `references/release.md` and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/release.md`, `skills/roundfix/SKILL.md`.

## Goals

- The status view shows a Delivery Queue item's completed and total Tasks as its Item Branch holds them, in the Spec queue and under `Pending Specs`, so a delivering Spec never reads as not started.
- In an adopter, the Release Plan needs no classification by hand for a Roundfix update, for the adoption of Roundfix recommendations, or for a plan that touches a `CONTEXT.md`, and still asks for one when a commit changes any other released path.
- Roundfix's own Release Plan classifies every range exactly as before.

## Core Features

1. The Task Graph loader can read a Spec folder through a given reader, so the status view loads a Spec as a Git revision holds it without writing a file or reading another working tree.
2. The status view counts a Delivery Queue item's Tasks on its Item Branch: the Spec folder at the branch head; once the branch holds the Spec's Archive Record, the folder at the record's source revision; with no branch yet, or a branch holding neither, the Spec Root of the checkout, as today.
3. The Release Plan recognises an adopter by its `skills-lock.json` at the target recording the `roundfix` skill, and there counts every Adopter Maintenance Path as maintenance-only, reading the Setup Manifest and `skills-lock.json` at the range's base and target.
4. The `follow-up` and release guides, the `release plan` help, the Roundfix skill's release reference, the glossary and ADR-0447 state both rules; the Skill Coverage Map records the Coverage Review of the status surfaces.

## Non-Goals / Out of Scope

- Reading a Delivery Queue item's worktree files, its uncommitted Task statuses or a Task Worktree.
- A new field, line or exit code in the status view or the Release Plan; a marker for an unreadable Item Branch.
- Counting other paths of an adopter as maintenance-only: a guide or skill the records do not name, `CONTEXT-MAP.md`, `docs/user-guide/`, nested `AGENTS.md` files, or any product file.
- Any change to Roundfix's own classification, to the Release Ceiling's authority or to the release-preparation commit rule.
- A Project Config key, or any read of Roundfix configuration by the Release Plan Command.
- Any Baseline guide or Baseline module; the Roundfix skill reaches adopters through the owned-skill installation.

## Success Metrics

1. Success Metric: on temporary repositories, a running Delivery Queue item whose Item Branch settled one of two Tasks shows `running 1/2 Tasks` where the unchanged tree shows `running 0/2 Tasks`; an item archived on its branch after a corrective Task shows the count of the source revision; an item with no branch, a missing branch or a branch without the Spec shows the Spec Root's count; an uncommitted status in the item's worktree is not counted; the `Pending Specs` line shows the same count; the reader writes nothing.
2. Success Metric: on temporary repositories, a release range of an adopter whose commits change only Adopter Maintenance Paths and planning evidence ends without a classification by hand, a commit that also changes another path still blocks, the same range in a repository whose lock records no `roundfix` skill blocks as before, and the plan ignores an invalid `.roundfixrc.yml`.
3. Success Metric: both copies of the Roundfix release reference, the release runbook, the usage guide and the `release plan` help name each kind of Adopter Maintenance Path and the adopter test, pinned phrase by phrase by contract tests that fail when any phrase is removed.
4. Success Metric, outside evidence: on the shape of the Fiscus range `v0.9.0..ad1381b` measured by the operator, the built binary leaves only the commit that changes `docs/user-guide/` and `.env.example` blocking, where the binary of the Delivery Base blocks four commits; and on this repository's own ranges from `v0.140.0`, `v0.147.0` and `v0.148.0`, both binaries print the same decision, blocking commits and crossing flags.

## Acceptance evidence

The outside-evidence row of Success Metric 4 rests on a repository this Spec did not build and a measurement it did not design: the Fiscus commits `4b83839`, `ec8481d`, `f923d37` and `ad1381b`, whose subjects and changed paths `references/2026-11-19-characterization.md` section 2 records as measured by the operator, and this repository's own release history, read offline from its Git objects. The QA gate rebuilds the Fiscus commits' path shape in a temporary repository from that record, because it reads no other repository, and runs both binaries on it. Published literature states the same path rule: googleapis/release-please#1875, "Idea is to exclude commits if all files in those commits are from specific paths."

## Unreachable Acceptance

None. The Roundfix ranges need the stable tags `v0.140.0`, `v0.147.0` and `v0.148.0` in the audited clone; a clone without them records that and runs the ranges from the earliest stable tag it holds.

## Validation

The decisions were checked on 2026-11-19 against the code of `internal/cli/follow_up_queue.go`, `internal/cli/follow_up_status.go`, `internal/spec/spec.go`, `internal/daemon/task_engine.go`, `internal/releaseplan` and `internal/cli/releaseplan_git_source.go`, the Setup Manifest and `skills-lock.json` of this repository and of the Fiscus, ADR-0349, ADR-0359, ADR-0410, ADR-0412 and ADR-0426, a reproduction of each defect and a prototype of both fixes on a scratch worktree, all recorded in `references/2026-11-19-characterization.md`. The Secondbrain (`wiki/index.md`, `wiki/entities/roundfix.md`; `qmd` was locked) and Exa (release-please `exclude-paths`) were consulted; section 4 of that file records what each gave.

## Skills

- unchanged: command follow-up — the status view keeps every command, flag, field, output line and exit code; only a Delivery Queue item's counts change source, and the Roundfix follow-up reference already relays the rows as the command gives them.
- unchanged: guide docs/user-guide/commands/follow-up.md — the guide gains one sentence on where a Delivery Queue item's Tasks are counted; the Roundfix follow-up reference does not restate how counts are read and needs no change.

## Glossary

- adds: **Adopter Maintenance Path**
- changes: **Spec Queue Table**

## Decisions

- 2026-11-19: When is the correction of the two defects planned, in how many Specs, with which authorization and release? — "Spec fix 0414 logo após a 0411 (Recommended)" — uma Spec com as duas correções (como a 0407), TechSpec no Opus, sem Ideia, "Mesmas do BE", release patch, entregue antes da 0409.
- 2026-11-19: Which files count as no-release maintenance in an adopter? — "Arquivos geridos pelo Roundfix + CONTEXT.md (Recommended)" — os arquivos que o Baseline e o skills install escrevem (Setup Manifest e skills-lock), .claude/agents gerado, .roundfixrc.yml e todo CONTEXT.md; o resto como hoje. Para o próprio Roundfix nada muda (skills e Baseline são produto dele).
- 2026-11-19: How are open questions settled while no person is present? — "Seguir o padrão do autor (Recommended)".

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0414 of 2026-11-19, the maintainer's answers above and the characterization.

Each item below is an operator decision delegated by the maintainer (2026-11-19), with its question and options; the first option is the one taken.

- Which source counts a delivering Spec's Tasks? Options: (a) the Item Branch through Git objects, then the Spec Root (taken: the branch holds every settled Task, the Daemon commits each Task with its status, and the `follow-up` guide promises to read no file of another checkout); (b) the item worktree's files, then the Item Branch, then the Spec Root (refused: reads another checkout and can read a Task file the Daemon is writing); (c) the Run's events (refused: a second source of Task status beside the Task files the Daemon owns, ADR-0057).
- What does a row show when the Item Branch cannot be read? Options: (a) the Spec Root's count, as today, with no new field (taken: keeps text and JSON compatible and the case needs a branch that lost both the folder and the record); (b) `?/?` in the status (refused: changes the text of every unreadable row, pinned by existing transcripts); (c) a new source entry naming the failure (refused: a new JSON entry for a case no Delivery Queue produces).
- How does the Release Plan know it runs in an adopter? Options: (a) the target's `skills-lock.json` records the `roundfix` skill (taken: an adopter installs that skill, Roundfix authors it and records only Vendored Skills, and no Roundfix configuration is read; ADR-0447); (b) a Project Config key (refused: breaks "reads no Roundfix configuration"); (c) a Roundfix source path such as `internal/baseline/assets/` (refused: ties the classifier to Roundfix's layout); (d) no test, the rule everywhere (refused: changes Roundfix's own classification, measured on `1a8e08b4`).
- Where does the managed set come from? Options: (a) the Setup Manifest and `skills-lock.json` at the base and the target, plus the fixed Baseline paths `CLAUDE.md` and `docs/agents/specific-repository.md`, the generated `.claude/agents/roundfix-*.md`, `.roundfixrc.yml` and every `CONTEXT.md` (taken: the records say exactly what the Baseline and the installation wrote, and a repository-local skill with scripts stays a released change); (b) a fixed path list (refused: counts a repository-local skill and a repository's own guides); (c) the records at every commit (refused: two reads per record cover what the range adds or retires).
- `type: fix` with a `patch` Release Ceiling: no command, flag, field or exit code is added.

## Open Questions

None.
