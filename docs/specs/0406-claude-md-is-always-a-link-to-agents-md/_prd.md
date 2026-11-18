---
spec: 0406-claude-md-is-always-a-link-to-agents-md
status: active
created: 2026-11-18
surfaces: [backend, cli, docs]
type: fix
---

# CLAUDE.md is always a symbolic link to AGENTS.md

Every repository that uses the Baseline must keep one root instruction file, `AGENTS.md`, with `CLAUDE.md` a symbolic link to it, and the project has no option to do otherwise. Roundfix does not enforce this today. Adoption backs up a regular `CLAUDE.md` and leaves it in place, a managed refresh never looks at it, and `roundfix doctor` says nothing, so a repository can keep a `CLAUDE.md` that Claude Code reads and Codex and OpenCode never see. After this Spec, adoption makes `CLAUDE.md` the link, `roundfix baseline update` restores it whenever no content would be lost and otherwise names each passage to move first, `roundfix doctor` reports any other shape, and the Baseline guide states the rule.

## Sources

The maintainer's decision of 2026-11-18 is recorded in [`references/2026-11-18-maintainer-decision.md`](references/2026-11-18-maintainer-decision.md). The read-only inventory of the adopters and the prototype measurements are in [`references/2026-11-18-claude-link-inventory-and-measurements.md`](references/2026-11-18-claude-link-inventory-and-measurements.md).

No open Backlog Entry and no unresolved Finding shares this context. Spec 0405, which adds the Claude-only statement to the root block, is the prerequisite: its statement is what made one file enough.

## Prerequisites

`requires`: Spec 0405, whose root statement this standard relies on and whose Managed Refresh of this repository runs first, and Spec 0399, merged, whose Misplaced Repository Rule detector this Spec reuses and narrows. Authored against `eb888135`.

## Project Constraints

- Identifier strategy: not applicable — no persistent identifier is minted. The new names are a Doctor check name, two `DR-` codes, two warning codes and one managed-entry identity, all public CLI contract under `docs/agents/cli.md`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call. The Doctor line and the plan read local files and, for one file, Git's index. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0432 and ADR-0433 (new) record the decisions; the design amends ADR-0100 and ADR-0411 and keeps ADR-0070, ADR-0071, ADR-0203, ADR-0220, ADR-0341 and ADR-0285; ADR-0074, ADR-0078, ADR-0173, ADR-0250, ADR-0254, ADR-0399 and ADR-0421 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-18, quoted in the authorization record. It covers the Baseline module, template and profile pin this Spec changes, their formatter golden, this repository's managed guide and Setup Manifest, and the Roundfix and `setup-context-driven` skill files listed below. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0406-claude-md-is-always-a-link-to-agents-md/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/setup-context-driven/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `internal/baseline/assets/templates/index.json`, `skills/roundfix/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## Goals

- Adoption leaves every adopted repository with `CLAUDE.md` a relative symbolic link to `AGENTS.md`, whatever `CLAUDE.md` was before, and the root backup keeps any bytes it replaced.
- `roundfix baseline update` restores the link when `CLAUDE.md` is absent, or is a regular file or a link elsewhere that holds no passage of its own, and keeps the Baseline's backup of a replaced regular file. When it holds one, the update leaves it unchanged, names each passage and its destination, and refreshes everything else.
- `roundfix doctor` prints one `baseline-claude-link` line: `ok` for the link, `warn` with the next action for any other shape, `skipped` without a Setup Manifest. It never fails.
- A checkout on a host that does not create symbolic links is recognised and told what to do, never converted.
- An Agent that reads the Baseline guides learns that `AGENTS.md` is the root instruction file and `CLAUDE.md` its link, and the skills and user guide describe the new behaviour.

## Core Features

1. A read-only reader in the `baseline` package classifies the root `CLAUDE.md` of a repository: the link, absent, a regular file, a link elsewhere, or an unlinked checkout, with the passages of its own that it holds.
2. Every adoption plan writes `CLAUDE.md` as the link, and a managed refresh does so whenever the file holds no passage of its own; the file transaction replaces a regular file or a link at that path and restores it on rollback.
3. The Misplaced Repository Rule detector no longer reports the `@AGENTS.md` import line or the bytes of an unlinked checkout.
4. `roundfix baseline update` prints the withheld and unlinked-checkout warnings, and `roundfix doctor` prints the `baseline-claude-link` line.
5. The `agent-instructions.md` guide template, the Roundfix skill's `baseline`, `modes` and `setup` references, the `setup-context-driven` skill, the user guide pages for `baseline` and `doctor`, and the glossary state the rule.

## User Experience

In a repository whose `CLAUDE.md` is a regular file holding `@AGENTS.md`, today's plan says:

```text
Misplaced repository rules: 1
- CLAUDE.md:1-1: @AGENTS.md
  Move to: docs/agents/specific-repository.md
```

With this Spec, the plan replaces it:

```text
File changes: 2
- update CLAUDE.md (1 managed entries)
- create CLAUDE.<digest>.md (1 managed entries)
Misplaced repository rules: 0
```

When `CLAUDE.md` holds a rule of its own, the plan keeps it and says:

```text
Misplaced repository rules: 1
- CLAUDE.md:1-3: # Claude notes
  Move to: docs/agents/specific-repository.md
Warning: baseline.root.claude-link.withheld: CLAUDE.md: CLAUDE.md holds 1 passage of its own (CLAUDE.md:1-3); move each into docs/agents/specific-repository.md, then run roundfix baseline update to replace it with a symbolic link to AGENTS.md
```

and `roundfix doctor` prints:

```text
baseline-claude-link: warn (DR-BASELINE-CLAUDE-LINK: CLAUDE.md is a regular file with 1 passage of its own: CLAUDE.md:1-3; next: move each passage into docs/agents/specific-repository.md, then run roundfix baseline update)
```

## Non-Goals / Out of Scope

- Moving a passage out of `CLAUDE.md`. A person or a supervised Agent moves it, as ADR-0411 decides.
- A repository option to keep a regular `CLAUDE.md`, or a copy of `AGENTS.md` where links do not work.
- Nested `CLAUDE.md` files, which stay warning-only carriers.
- Changing an `AGENTS.md` that is itself a link, or the block on an unsafe `CLAUDE.md` alias.
- Editing any adopter's repository or sending it a message. Every adopter on the host already keeps the link; the maintainer updates the fleet after the release.

## Success Metrics

1. Success Metric: a greenfield and a preservation adoption of fixture repositories whose `CLAUDE.md` is absent, a regular file, the `@AGENTS.md` import or a link to another file each end with `CLAUDE.md` a symbolic link whose target is exactly `AGENTS.md`, and every replaced regular file has its backup.
2. Success Metric: on adopted fixtures, a managed refresh creates the link when `CLAUDE.md` is absent; converts a regular file that is empty, holds only `@AGENTS.md`, only Managed Regions or a copy of `AGENTS.md`, writing its backup, and a link to a file without a passage of its own; and leaves byte-identical, with the withheld warning naming each passage, a regular file or a linked file that holds one. A second refresh after each conversion reports `current`.
3. Success Metric: an unlinked checkout reproduced with `core.symlinks=false` is left byte-identical by a plan, with the unlinked-checkout warning, and the detector reports no passage for it.
4. Success Metric: `roundfix doctor` prints the `ok`, `warn` and `skipped` forms of `baseline-claude-link` for each shape and exits as it did before.
5. Success Metric: a fault injected after the link replaces a regular `CLAUDE.md` rolls the apply back to the original bytes and mode.
6. Success Metric, outside evidence: run read-only with the built binary in each adopted checkout of the inventory, `roundfix baseline update --no-skills --no-history --format json` plans no change and no warning for `CLAUDE.md`, and every `git status --porcelain` stays unchanged.

## Acceptance evidence

The outside-evidence row rests on repositories this Spec did not build: the six distinct adopted repositories of the inventory of 2026-11-18, seven checkouts in all, each recorded with its commit. The QA gate re-reads them read-only. For a checkout whose commit moved, it compares with a fresh read made by the inventory's method.

## Unreachable Acceptance

None. The adopted repositories are local checkouts. When the QA host lacks one, the QA gate runs the same method on the checkouts it has and names the missing one.

## Validation

The decisions were checked on 2026-11-18 against these sources:

- The seven adopted checkouts, read as the inventory describes.
- A prototype on a scratch worktree, which adopted a fixture, converted and created the link in a managed refresh, withheld a file with a rule of its own, and ran the `internal/baseline` and `internal/cli` suites.
- The code of `internal/baseline/plan.go`, `apply.go`, `transaction.go`, `repository.go` and `misplaced_rules.go`, and of `internal/cli/doctor.go`.
- ADR-0070, ADR-0100, ADR-0341 and ADR-0411.

## Glossary

- adds: **Claude Instruction Link**
- changes: **Managed Refresh**
- changes: **Misplaced Repository Rule**

## Decisions

- 2026-11-18: What standard do the root instruction files follow? — "Se não for necessário e formos manter arquivo AGENTS.md e o CLAUDE.md com symlink para o AGENTS.md, vai ser sempre assim. Sem opção ao usuário de mudar o baseline no projeto. O único caminho para personalização do baseline é o arquivo que contem regras especificas do repositorio, hoje chamado de specific-repository.md"; after the Supervisor's answer that the Claude-only statement is scoped "In Claude Code,", one `AGENTS.md` and `CLAUDE.md` always a link.
- 2026-11-18: Where does the rule that makes the link mandatory go? — "Spec 0406 depois da 0405 (Recommended)".
- 2026-11-18: Which authorization applies? — "Mesmas do BE (Recommended)": implement, commit, push, pull request, merge and release in the single release at the end of the queue, up to five corrections; Makefile, CI, `go.mod` and lint only if a correction requires one.
- 2026-11-18: How does an update treat a regular `CLAUDE.md`? — "Só depois de mover (Recommended)": when the regular `CLAUDE.md` holds passages of its own, the update plan lists each one to move to `docs/agents/specific-repository.md` and does not convert; with no passage of its own it converts it to the symbolic link, with the Baseline's backup; no content is ever lost.
- 2026-11-18: Who decides the remaining product questions of Spec 0406? — "deve parar e perguntar.": a product question stops the authoring and goes back to the maintainer; technical choices stay with the author and are recorded as ADRs.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for Spec 0406 of 2026-11-18, the three answers above and the inventory.

- A managed refresh that withholds the conversion still refreshes everything else and never refuses the update; the conversion of a regular file writes the root backup `CLAUDE.<digest>.md` (ADR-0432).
- An unlinked checkout is left as it is and reported (ADR-0433).
- The rule is stated in the intro of the `agent-instructions.md` guide template, not in a new Normative Clause, so no Source Baseline row is needed.
- `type: fix` with a `minor` Release Ceiling, because Doctor gains a line and the update gains two warnings.

## Open Questions

None.
