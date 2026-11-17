---
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
status: active
created: 2026-11-16
surfaces: [backend, cli, docs]
type: fix
---

# Adopters keep their own rules in the specific-repository guide

Adopters put their own rules in files the Baseline owns. They append sections to `docs/agents/skill-dispatch.md`, write overrides above the Managed Region of `docs/agents/agent-instructions.md`, and add rules between the Managed Regions of the root `AGENTS.md`. They follow the text they received: the guide templates say repository authors "own project-specific extensions outside setup markers". A managed refresh preserves those bytes and says nothing about them, so the fleet's rules end up in many places and drift from the Baseline next to them. After this Spec, the Baseline says in the guides an Agent reads before editing that a repository's own rules and variations go only in `docs/agents/specific-repository.md`, and it shows how to write one there. `roundfix baseline update` and `roundfix doctor` name each passage to move and where it goes, and neither moves anything.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see `references/_index.md`):

1. [`references/2026-11-16-adopters-put-their-own-rules-in-baseline-owned-files.md`](references/2026-11-16-adopters-put-their-own-rules-in-baseline-owned-files.md), recorded by the operator from the maintainer's report in the Supervisor session of 2026-11-16.

The maintainer's answers are in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The inventory his answer asked for is in [`references/2026-11-16-adopter-rule-placement-inventory.md`](references/2026-11-16-adopter-rule-placement-inventory.md).

No other open Backlog Entry and no unresolved Finding shares this context:

- The docs-folders refactor entry is parked by the maintainer.
- The Delivery Queue and release-guide entry changes delivery, not rule placement.
- The large-context entry waits on measurements from Spec 0397.
- The Conexus status entry is about the status view.

## Prerequisites

None in `requires`. Authored against `bcf0e675` (v0.147.0). No active Spec shares a file with this one.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The JSON field, the Doctor line and its `DR-` code are named in the TechSpec. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call. The new Doctor line reads local files only. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0411 (new) records the decision; the design keeps ADR-0145 (an Unrecorded Managed Region is refreshed and named), ADR-0203 (the repository owns its rules file), ADR-0220 (only a definite answer fails Doctor), ADR-0186 and ADR-0257 (adopter-neutral clauses keep their identities when reworded); ADR-0222, ADR-0071, ADR-0258 and ADR-0342 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16, quoted in the authorization record. It covers the Baseline modules and templates, their derived goldens, this repository's managed guides and Setup Manifest, and the owned skills this Spec changes. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0399-adopters-keep-their-own-rules-in-the-specific-repository-guide/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/setup-context-driven/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `docs/agents/skill-dispatch.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/repository-extension.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/templates/extensions/specific-repository.md`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `internal/baseline/assets/templates/guides/skill-dispatch.md`, `internal/baseline/assets/templates/index.json`, `internal/baseline/assets/templates/root/repository-extension.md`, `skills/roundfix/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## Goals

- An Agent that reads the Baseline guides, or the root instruction file, before it edits learns that a repository's own rules and variations go only in `docs/agents/specific-repository.md`. It also learns how to write one there.
- `roundfix baseline update` names every Misplaced Repository Rule by path, line range and first line, with `docs/agents/specific-repository.md` as its destination. It does so in text and JSON, in plan, current and applied results, without adding to the Plan Document, changing the exit code or writing any byte.
- `roundfix doctor` prints one `baseline-rules` line: `ok`, `warn` naming the passages and the move, or `skipped` without a Setup Manifest. The line never fails.
- The five adopted repositories on the authoring host get exactly the passages the inventory lists. The Turborepo block, the clean repositories and the `CLAUDE.md` symbolic link report nothing.

## Core Features

1. The Baseline text states the rule:
   - The guide templates of `agent-instructions.md` and `skill-dispatch.md`.
   - The root pointer to the repository's own rules.
   - Three reworded clauses that keep their identities and force: `clause.core.keep-root-compact`, `clause.core.activate-matching-skills` and `clause.context.docs-one-job-per-directory`.
   - The starter text of a new `docs/agents/specific-repository.md`, which shows how a rule is written there.
2. A read-only detector in the `baseline` package lists the Misplaced Repository Rules of a repository from its Setup Manifest.
3. `roundfix baseline update` prints them in every result it gives once it reads a valid Setup Manifest.
4. `roundfix doctor` prints one `baseline-rules` line for them.
5. The user guide pages for `baseline` and `doctor`, the Roundfix skill's `baseline` and `setup` references, the `setup-context-driven` skill and the glossary describe them.

## User Experience

Today, in the quaestor checkout, the plan says nothing about the three rules it added after the Managed Region of `skill-dispatch.md`:

```text
File changes: 2
- update docs/agents/autonomous-work.md (1 managed entries)
- update docs/agents/setup-context.json (1 managed entries)
Retention evidence: 0
```

With this Spec, the same plan adds:

```text
Misplaced repository rules: 1
- docs/agents/skill-dispatch.md:109-115: ## Repository-specific rules
  Move to: docs/agents/specific-repository.md
```

`roundfix doctor` in that checkout prints:

```text
baseline-rules: warn (DR-BASELINE-RULE-MISPLACED: 1 passage belongs in docs/agents/specific-repository.md: docs/agents/skill-dispatch.md:109-115; next: move each passage into docs/agents/specific-repository.md; roundfix baseline update lists them)
```

## Non-Goals / Out of Scope

- Moving, deleting or rewriting a passage. A person or a supervised Agent moves it.
- Refusing an update, or failing Doctor, because of a passage.
- Inspecting `docs/agents/specific-repository.md`, nested carriers, or guides the Setup Manifest does not list.
- Judging whether a passage is a rule. Every passage outside a Managed Region is named.
- Changing a Managed Refresh, the Plan Document, the Plan Digest or an exit code.
- Editing any adopter's repository or sending it a message. The maintainer updates the fleet after the release.
- The docs-folders refactor, which the maintainer parked.

## Success Metrics

1. Success Metric: the rendered `docs/agents/agent-instructions.md` and `docs/agents/skill-dispatch.md` of this repository, and the formatter golden of the Standard TypeScript profile, carry the placement rule and its example. The root `AGENTS.md` pointer says that rules are written only in `docs/agents/specific-repository.md`. A new `docs/agents/specific-repository.md` starts with the how-to text. The three clauses keep their identities and their `mandatory` force.
2. Success Metric: for fixture repositories built from the inventory's shapes, the detector returns exactly the expected path, line range and first line of each passage. These shapes are:
   - a section after a guide's Managed Region;
   - one before it;
   - one between two root Managed Regions, next to a foreign `BEGIN:`/`END:` block;
   - several sections after one Managed Region;
   - an Unrecorded Managed Region.

   For a clean carrier, a symbolic-link `CLAUDE.md` and a foreign block, it returns nothing.
3. Success Metric: `roundfix baseline update --format json` carries `misplacedRepositoryRules` in plan, current and applied results. The text prints the `Misplaced repository rules: <n>` section. The Plan Document gains no field, and the exit code equals the one for the same repository without the passage.
4. Success Metric: `roundfix doctor` prints `baseline-rules: ok`, `warn` or `skipped` as above and exits as it did before.
5. Success Metric, outside evidence: run read-only with the built binary, `roundfix baseline update --no-skills --no-history --format json` reports, for quaestor, indexus and fiscus, exactly the passages the inventory lists for their recorded commits. For argus it reports none. It leaves each `git status --porcelain` unchanged.

## Acceptance evidence

The outside-evidence row rests on repositories this Spec did not build. These are the five adopted repositories on the authoring host: argus, quaestor, fiscus, indexus and this repository. The inventory of 2026-11-16 records each one with its commit. The QA gate re-reads them read-only. For a repository whose commit moved, it compares the report with a fresh read made by the inventory's method.

## Unreachable Acceptance

None. The adopted repositories are local checkouts. When the QA host lacks one, the QA gate runs the same method on the checkouts it has and names the missing one.

## Validation

The decisions were checked on 2026-11-16 against these sources:

- The five checkouts, read as the inventory describes.
- The installed Roundfix 0.147.0 plan in four of them, which reported none of the passages.
- The code of `internal/baseline/preservation.go`: `classifyManagedRegions` and `partitionRootSource`.
- The code of `internal/cli/baseline_update.go` and `internal/cli/doctor.go`.
- ADR-0145, ADR-0203, ADR-0220 and ADR-0257.

## Glossary

- adds: **Misplaced Repository Rule**

## Decisions

- 2026-11-16: What scope does the fix take? — "Texto + detecção (Recommended)": the rule and an example in agent-instructions, skill-dispatch and specific-repository. `baseline update` and `doctor`, on finding a local edit in a Baseline-owned file or a rule in the root `AGENTS.md`, say which passage to move to `specific-repository.md`.
- 2026-11-16: How is the adopters' current state inventoried? — "Ler os checkouts (Recommended)": read-only reads of the adopters' checkouts on this host. No commit and no message to another project. The inventory becomes a reference of the Spec and the source of its test cases.
- 2026-11-16: Which authorization applies? — "Mesmas do BE (Recommended)": push, pull request, merge and release right after the Spec, up to three corrections. Makefile, CI, `go.mod` and lint only if a correction requires one.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for the TechSpec of Spec 0399 of 2026-11-16, the three answers above and the inventory.

- No automatic move, from the operator's brief.
- `type: fix` with a `minor` Release Ceiling, because the JSON result gains a field and Doctor gains a line.
- Detection outside the Plan Document (ADR-0411).
- A foreign `BEGIN:`/`END:` block is not reported (ADR-0411).

## Open Questions

None.
