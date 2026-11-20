---
spec: 0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set
status: active
created: 2026-11-20
surfaces: [backend, cli, docs]
type: fix
---

# A skill the skill dispatch names joins the Repository Skill Set

An adopter whose Baseline Profile selects the `typescript` module without the `frontend` module reads a skill dispatch guide that makes the skill `solid` mandatory for production code, while `roundfix baseline update`, `roundfix skills check`, Doctor and the restore command never require it: the restore command refuses it by name. A delivery in such an adopter stopped on 2026-11-20 for that skill. The cause is in the embedded catalog: the `typescript` module owns the production-code Skill Activation, and neither it nor the `core` module it depends on requires `solid`; only `frontend` does, and nothing checks an activation's bundle against its owner. After this Spec the `typescript` module requires `solid`, the catalog refuses any Skill Activation whose bundle names a skill its owner does not require, and one function reads the required skills for restore, the update and Doctor.

## Sources

The source is the maintainer's request of 2026-11-20, relayed by the operator's brief and copied to [`references/2026-11-20-request-and-characterization.md`](references/2026-11-20-request-and-characterization.md) with the measured cause, the reproduction through the released binary, the sweep for other gaps, the prototype and the sources consulted. That file is cited below as "the characterization". The maintainer's two answers of the same day, on the authorization and on a check at Run start, are copied whole to [`references/2026-11-20-maintainer-answers.md`](references/2026-11-20-maintainer-answers.md).

No open Backlog Entry or unresolved Finding shares this context. The seven open Backlog Entries concern a team's status report, the size of a session's context, the name of the guides folder, a `spec check` operand, the semantic analysis deadline, dead probe code and a vacuous Verification; none changes the Baseline catalog, the Repository Skill Set or the skill dispatch. `docs/findings/` is empty.

## Prerequisites

None. Authored against `851024fa` (v0.152.0). No Verification of this Spec depends on another Spec.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the catalog diagnostic reuses the activation identifier and the skill name the catalog already holds. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; catalog validation and the readers of the required skills read embedded assets and local files. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0448 (new) records the decision; ADR-0067, ADR-0191, ADR-0222 and ADR-0241 bind the design; ADR-0074, ADR-0204, ADR-0206, ADR-0256, ADR-0257, ADR-0322, ADR-0346, ADR-0350, ADR-0360 and ADR-0364 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Tudo, com patch v0.152.1" of 2026-11-20, quoted in the authorization record; it bounds one Baseline module edit and the Managed Refresh of this repository's Setup Manifest, grants a release at a `patch` Release Ceiling, and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`; Spec-contained authorization record: `docs/specs/0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set/_authorization.md`; bounded files: `docs/agents/setup-context.json`, `internal/baseline/assets/modules/typescript.json`.

## Goals

- A repository whose profile selects `typescript` without `frontend` has `solid` in its Repository Skill Set: the update preview and `roundfix skills check` name it as a Missing Required Skill with the restore command, and the approved update and the restore command accept it.
- No composition of the embedded catalog can render a skill in the skill dispatch guide that its Repository Skill Set lacks: the catalog fails to load when a Skill Activation's bundle names a skill its owner does not require.
- Restore, the update and Doctor read the required skills through one function.
- The required skills of every built-in Baseline Profile stay the same names.

## Core Features

1. The `typescript` module requires `solid`.
2. Catalog validation refuses a Skill Activation whose bundle names a skill that neither its owner nor a module the owner depends on requires, and names the activation and the skill.
3. One function returns the skills a list of modules requires; restore for a built-in profile, restore for a repository profile and the requirement Doctor, the update preview and `skills check` read all call it.
4. The glossary entry of Skill Activation states the rule, and ADR-0448 records the decision.

## Non-Goals / Out of Scope

- A check at Run start that refuses or warns about a missing required skill. None exists today, and the maintainer decided on 2026-11-20 that the refusal becomes its own Spec after this one (`## Decisions`).
- Changing Doctor's next action for a missing external skill, which names the skill and the catalog install command (`## Decisions`).
- Computing the Repository Skill Set from the skill dispatch instead of from the modules' `requiredSkills` (ADR-0448).
- A new check in `roundfix baseline profile validate`: a repository profile composes only embedded modules and must select every dependency, so the catalog rule covers it.
- Moving the production-code Skill Activation to another module, or removing `solid` from its bundle.
- Any change to `internal/baseline/assets/skill-activations.json`, the `frontend` module, a Setup Snapshot, a Baseline guide, an owned skill, a user-guide page, a command, a flag, an output line or an exit code.
- The adopter's repository, where another session installed `solid`.
- The change of the knowledge-base search in the Secondbrain guide, which is a separate Spec.

## Success Metrics

1. Success Metric: on a repository profile made from `standard-typescript-monorepo` and cut to `core`, `context-workflow`, `typescript`, `bun`, `autonomous-work`, `spec-workflow`, `secondbrain` and `repository-extension`, the required skills hold `solid` with the contract the TypeScript Setup Snapshot pins, and a restore that selects `solid` passes the selection.
2. Success Metric: a copy of the embedded catalog whose `typescript` module no longer requires `solid` fails to load with one `catalog.skill.activation.skill.unrequired` diagnostic that names `trigger.production-code` and `solid`; the embedded catalog loads without it.
3. Success Metric: for every module of the embedded catalog taken with its dependencies, and for every built-in profile, each skill the rendered skill dispatch guide names is among the required skills.
4. Success Metric: through the built binary, in the temporary repository of the characterization, section 3, `roundfix skills check` prints the line for `.agents/skills/solid` with `roundfix baseline skills restore --profile ts-backend --skill solid`, and that restore command with an absent `--source-dir` reports `restore.source-dir-invalid`, never `restore.skill-invalid`.
5. Success Metric: the required skills of the five built-in profiles are the names measured before the change: 41, 105, 38, 34 and 95 of them (the characterization, section 5).
6. Success Metric, outside evidence: the catalog rule is the declared-dependency rule of a published build system, quoted under `## Acceptance evidence`, applied to a Skill Activation and its owner.

## Acceptance evidence

The outside-evidence row of Success Metric 6 rests on published literature this Spec did not write. The Bazel documentation, "Dependencies" (https://bazel.build/concepts/dependencies), states: "For correct builds, the graph of actual dependencies A must be a subgraph of the graph of declared dependencies D." Its example is this defect: a target uses what it never declared, keeps building while another target brings the dependency in, and breaks when that target leaves. Here the `typescript` module used `solid` through its activation, never declared it, and broke when a profile left `frontend` out. The QA gate reads the quoted passage from this file, offline, and checks that the diagnostic of Success Metric 2 enforces it.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-20 against the code and the embedded catalog at `851024fa`, a reproduction in the package and through the released binary, a sweep of all 11 activations and 17 modules, and a prototype on a scratch clone, all recorded in the characterization. The Secondbrain file `projects/sql-library/mirror/docs/adr/0004-o-baseline-foi-adotado-com-profile-proprio.md` showed a second adopter whose profile is a built-in profile with modules deleted, which is why the rule holds per owner and not per built-in profile. The Bazel documentation supports declaring at the owner and checking the declaration; the Go modules reference describes the rejected alternative, computing the list from the references. The `qmd` query of the Secondbrain was attempted three times and failed with a locked database; the Indexus search answered instead.

## Skills

- unchanged: command doctor — Doctor reads the required external skills through one Catalog method and keeps every check name, status, message, next action, output line and exit code; the Roundfix setup reference describes Doctor's checks as they stay and needs no change.

## Glossary

- changes: **Skill Activation**

## Decisions

- 2026-11-20: "Autoriza implementar, commit, push, PR, merge e soltar a versão da Spec 0415?" — option 1, "Tudo, com patch v0.152.1": the operations `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling, the two Governed Paths and the regeneration by the repository's commands, as `_authorization.md` quotes it from ai-memory message `01a1f976-6dbf-7a43-93ee-89cc6806e754`.

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0415 of 2026-11-20, the maintainer's request it relays and the characterization.

- The brief left the design to the author. The author chose to declare and check: the owner of a Skill Activation requires every skill of its bundle, and the catalog refuses anything else (ADR-0448). Computing the set from the dispatch was rejected.
- `solid` joins the `typescript` module and stays in the `frontend` module, which dispatches it under its own trigger.
- The three readers of the required skills become one function; no output of restore, the update, `skills check` or Doctor changes for a built-in profile.
- `roundfix baseline profile validate` gains no check: the catalog rule covers every repository profile.
- The update preview, the approved update and `skills check` already name a Missing Required Skill with the restore command; once the set is right they name `solid`, and nothing is added to them.
- `type: fix`, released at a `patch` Release Ceiling by the maintainer's answer above; no command, flag, field or output line is added.

- delegated: 2026-11-20, techspec
  - question: Does this fix add a check at Run start that refuses or warns when a required skill is not installed?
  - options: no check at Run start, the update preview, the approved update, skills check and Doctor keep reporting it; a warning in the Run's log; a refusal before the first Task
  - adopted: no check at Run start — no Run-start skill check exists to correct, the brief asks to fill only the gap of what exists, and a new refusal is product behavior that the feature route decides
  - under: Project Config
  - answer: 2026-11-20, kept — "Depois, Spec própria"
- delegated: 2026-11-20, techspec
  - question: Does Doctor's next action for a missing external skill change from the catalog install command to the restore command?
  - options: keep the catalog install command the Doctor guide documents; change it to the restore command the update preview prints
  - adopted: keep the catalog install command the Doctor guide documents — Doctor already names the missing skill and a command that installs it, and changing a documented output line is outside a fix that changes no command output
  - under: Project Config
  - answer: none yet

## Open Questions

None. The maintainer answered the authorization on 2026-11-20 (`## Decisions`). The refusal at Run start when a required skill is missing is not part of this Spec: by the maintainer's answer of the same day, the session opens it as its own Spec after this one, with its own authorization request.
