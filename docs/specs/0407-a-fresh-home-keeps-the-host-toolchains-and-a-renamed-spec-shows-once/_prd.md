---
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
status: active
created: 2026-11-18
surfaces: [backend, cli, docs]
type: fix
---

# A fresh home keeps the host toolchains, and a renamed Spec shows once

Two defects of the tooling around a Spec's own lifecycle, reported together by the maintainer. The QA gates of Specs 0399 to 0402 ran test suites under a fresh `HOME` that hid the host's Rust toolchain, so `TestGuidanceCompositionJourney/rust-cli` failed and each Spec needed a plan amendment or a QA Archive Override; neither the `qa-gate` skill nor the `write-tasks` skill says what a fresh home must keep. And `roundfix follow-up --status` lists a Spec whose slug was renamed twice, once as the Spec and once as the old slug still being authored, because it joins the Phase Record to the checkout by slug. After this Spec, both skills state that suites keep the session's own home and that every command under a fresh home receives the toolchain locations of the repository's stacks, and the status view counts a phase entry only while no other slug of the same Spec Number holds the Spec.

## Sources

The adopted Backlog Entry [`references/2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md`](references/2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md) is the maintainer's report of both defects. The maintainer's answers that concern this Spec are in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The characterization, the reproduction and the prototype are in [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md).

The two defects share no component: one is the text of two owned skills, the other the status view's row reader. The maintainer, told of that measurement, chose to keep them in one Spec (`## Decisions`): they are one Backlog Entry the maintainer authorized as one Spec, a source keeps exactly one owning Spec, and each half alone would be a one-Task Spec paying a whole Run, QA gate, review and release. No other open Backlog Entry or unresolved Finding changes the QA environment or the status view's authoring rows.

## Prerequisites

None. Authored against `1a8e08b4` (v0.148.0). Spec 0408, authored beside this one, delivers after it and also changes `CONTEXT.md`, `docs/references/behavior-surfaces.json` and `skills/testdata/owned-skill-versions.json`; each Spec regenerates the two records on the tree it runs on.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the status view reads the existing Spec Number of a slug. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; the status view reads local files and the Phase Record, and the skills describe local environment variables. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0434 (new) records the status decision and amends ADR-0410; ADR-0080, ADR-0353, ADR-0407, ADR-0418 and ADR-0421 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-18, quoted in the authorization record. It covers the `qa-gate` and `write-tasks` skill files and their mirrors, and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/qa-gate/SKILL.md`, `skills/write-tasks/SKILL.md`.

## Goals

- A QA gate keeps the session's own `HOME` and `TMPDIR` for test suites and the repository Verification, and passes the host toolchain locations to every command it runs under a fresh `HOME`, so a fresh home never hides an installed toolchain, in this repository or in an adopter of any stack.
- A `qa` Task author is told the same default, so plans stop prescribing a fresh `HOME` that hides a toolchain.
- The status view shows one row per Spec Number among the Specs being authored: the phase entries of a renamed Spec's earlier slug make no row.

## Core Features

1. The `qa-gate` skill gains a section that states which commands keep the session's home, which toolchain locations a command under a fresh `HOME` receives, read before `HOME` changes, and that a failure caused only by a hidden toolchain is rerun, never recorded as an environment block.
2. The `write-tasks` skill gains a section that tells a `qa` Task author the same default and names the toolchain locations as the exception to a requirement that forbids reading the real home directory.
3. `followUpQueueRows` drops a phase entry's slug when another slug of the same Spec Number has a folder under the Spec Root, an Archive Record, a Delivery Queue item or a later counted phase entry.
4. The `follow-up` user guide, the glossary entry **Spec Queue Table** and ADR-0410 state the rule, and the Skill Coverage Map records the Coverage Review of the two Behavior Surfaces it touches.

## Non-Goals / Out of Scope

- Making the Daemon or any command set or isolate `HOME`; no Roundfix code isolates the QA gate's home today, and the gate's environment stays the skill's instruction.
- Changing `testfixture.RequireTool` to skip a tool that is present but cannot run; that would let a gate pass without the test it names.
- Rewriting the Phase Record, or reading renames from Git history.
- Showing a marked row for a renamed slug.
- Any Baseline guide or Baseline module; the skills reach adopters through the owned-skill installation.

## Success Metrics

1. Success Metric: both copies of the `qa-gate` skill hold the section on a fresh home with each of its guarantees, and both copies of the `write-tasks` skill hold the QA environment section, pinned phrase by phrase by a contract test that fails when any phrase is removed.
2. Success Metric: on temporary repositories, the status view drops the old slug's authoring row when the renamed Spec is planned, has only an Idea folder, is a queue item, is a merged queue item, has an Archive Record, or has a later phase entry, and keeps a row for another Spec Number and for a folder under the Spec Root.
3. Success Metric: through the built binary, the reproduction of `references/2026-11-18-characterization.md` section 2 prints one `0401` row, and a second Spec Number still prints its own row.
4. Success Metric, outside evidence: on a QA host with rustup, under a fresh temporary `HOME`, `TestGuidanceCompositionJourney/rust-cli` fails when the toolchain locations are not passed and passes when the gate passes the ones the `qa-gate` section names, and the whole `./internal/cli` package passes that way.

## Acceptance evidence

The outside-evidence row of Success Metric 4 rests on a measurement this Spec did not design: the QA Reports of Specs 0399 and 0402 and the QA Archive Overrides of Specs 0400 and 0401 recorded the failure, at the source revisions their Archive Records name under `docs/history/specs/`, and the published rustup book states that `RUSTUP_HOME` "Sets the root `rustup` folder, used for storing installed toolchains and configuration options" and defaults to `~/.rustup`. The QA gate reads those records offline and reproduces the failure and its fix on its own host.

## Unreachable Acceptance

None. A QA host without rustup records Success Metric 4 with the passage the rustup book gives and the Archive Records, and names that rustup is absent.

## Validation

The decisions were checked on 2026-11-18 against the QA Reports and Archive Records of Specs 0399 to 0402, the code of `internal/cli/follow_up_queue.go` and `internal/phaserecord`, the `qa-gate`, `write-tasks` and `write-prd` skills, ADR-0080, ADR-0407 and ADR-0410, a reproduction with the built binary and a prototype on a scratch worktree, all recorded in `references/2026-11-18-characterization.md`.

## Skills

- unchanged: command follow-up — the status view keeps every command, flag, field, output line and exit code; only a renamed Spec's earlier slug stops making a row, and the Roundfix follow-up reference already relays the rows as the command gives them.
- unchanged: guide docs/user-guide/commands/follow-up.md — the guide gains one sentence on renamed Specs; the Roundfix follow-up reference does not restate when a Spec is being authored and needs no change.

## Glossary

- changes: **Spec Queue Table**

## Decisions

- 2026-11-18: When does the corrective Spec for the two defects start, and with which authorization? — "Spec após a release"; authorization "Mesmas do BE (Recommended)".
- 2026-11-18: In which order do the four Specs of the cycle deliver? — "Ordem: \"Corretiva → regras → Indexus (Recommended)\"; o fix da fila entra logo após a corretiva."
- 2026-11-18, about 22:35 BRT: The two defects share no component, as the author measured; is the Spec split? — "Uma Spec só (Recommended)": both corrections stay in one Spec and ship in one patch release.

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0407 of 2026-11-18, the maintainer's answers above and the characterization.

- Both defects stay in one Spec, by the maintainer's answer below to the measured finding that they share no component.
- The QA environment rule lives in the `qa-gate` and `write-tasks` skills, which every adopter installs, and names Go and Rust as examples of a general rule, so a repository without Rust passes nothing extra.
- Suites and the repository Verification keep the session's own home; a fresh home is for rows that run the built product.
- A failure caused only by a fresh home hiding a toolchain is a gate error to rerun, never an environment block.
- The status view drops a renamed Spec's earlier slug, identified by the Spec Number (ADR-0434).
- No ADR records the QA environment rule: the skill text is its durable record and no code enforces it.
- `type: fix` with a `patch` Release Ceiling: no command, flag, field or output line is added.

## Open Questions

None.
