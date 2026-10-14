---
task: task_04
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: completed
type: docs
complexity: low
---

# Task 04: The skills, the QA prompt, the guides and the glossary describe it

## Overview

task_01 to task_03 refuse a QA Report commit and a push that carry a Secret
Match and keep the Judge Key Variables out of a Run. This Task writes the
behavioral rule a QA Agent needs, that evidence records environment variable
names and never values, into the `qa-gate` skill and the Daemon's QA prompt,
and describes the refusals, the `secret-found` park and the withheld variables
in the `roundfix` skill, the user guides and `CONTEXT.md`. It is verifiable by
phrase checks, the mirror check, the skill version record and the Behavior
Surface Record.

## Requirements

1. MUST answer both Backlog Entries of 2026-10-13,
   [QA evidence records environment values, secrets included](references/2026-10-13-qa-evidence-records-environment-values.md)
   and
   [A QA Run can send Spec content to the judge provider](references/2026-10-13-a-qa-run-can-send-spec-content-to-the-judge-provider.md),
   for the rule the QA protocol lacked: the `qa-execution` protocol the first
   entry names is the `qa-gate` skill and the QA prompt.
2. MUST add, in `.agents/skills/qa-gate/SKILL.md`, a `### Evidence carries no
   secret` heading directly after `### Runnable evidence` and before
   `### QA settlement`, saying: record environment variable names, never
   their values, and record whether a variable is set as `<NAME> set` or
   `<NAME> not set`; never persist the output of a command that prints the
   environment, naming `env`, `printenv`, `set`, `export -p`, `declare -x`,
   `make -p` and `make --print-data-base`; never write a secret-shaped sample
   into evidence, recording a Secret Match by path, line and kind only; never
   reload a shell profile to obtain a key; and that the Daemon refuses to
   commit a QA Report whose files match a Secret Pattern. The bytes of every
   `### QA settlement` section MUST NOT change.
3. MUST add to `qaGateContract` in `internal/agent/spec_prompt.go`, before its
   last line, the line `- Record environment variable names, never their
   values, and never persist the output of a command that prints the
   environment (env, printenv, set, export -p, declare -x, make -p); the
   Daemon refuses a QA Report whose files match a Secret Pattern.`, keeping
   every other line.
4. MUST describe, in `.agents/skills/roundfix/references/implement.md` at the
   QA Report commit boundary, the refusal of `_techspec.md` → API Contract 2
   with the phrase `QA evidence carries a secret` and the Final Push refusal;
   in `.agents/skills/roundfix/references/deliver.md` the `secret-found`
   blocker in the Park Class table under `finding` and its next action from
   API Contract 4; and in `.agents/skills/roundfix/references/runtime.md`
   that no Agent Session and no Verification command receives the Judge Key
   Variables and that the light session reads its key under
   `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` (API Contract 5), replacing the
   sentence that says it receives only the selected variable's name.
5. MUST make the same changes in `docs/user-guide/commands/implement.md`,
   `docs/user-guide/commands/deliver.md` (its Park Class table) and
   `docs/user-guide/configuration.md` (its light-tier key paragraph), with the
   same three phrases, and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
6. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise and record the `qa-gate` and `roundfix` versions; no mirror or
   record is hand-edited.
7. MUST add **Secret Pattern**, **Secret Match**, **Secret Scan** and
   **Judge Key Variables** to `CONTEXT.md` through `domain-modeling`, near
   **Archive Advice** and **Stage Key**, each linked to ADR-0273: a Secret
   Pattern is one reviewed credential shape with a kind name, kept in one
   table; a Secret Match is one line that matches a Secret Pattern, reported
   by path, line and kind and never by value; a Secret Scan checks the files
   a QA Report commit would stage and the commits a Roundfix push would send,
   and refuses on a Secret Match; the Judge Key Variables are the variables
   the judge reads its credentials from, which no Agent Session and no
   Verification command of a Run receives.
8. MUST NOT write any string a Secret Pattern matches into any file.

## Subtasks

- [ ] Add the evidence rule to the `qa-gate` skill and the QA prompt.
- [ ] Describe the refusals, the park and the withheld variables in the
      `roundfix` references and the user guides.
- [ ] Sync the mirrors, record both skill versions and the Behavior Surface
      Record.
- [ ] Add the four glossary terms.

## Acceptance Criteria

- [ ] The `qa-gate` skill and the QA prompt carry the evidence rule, and every
      `### QA settlement` section is unchanged.
- [ ] Both copies of the three `roundfix` references and the three guides
      carry their phrases, and the mirrors match.
- [ ] Both skill versions, the Behavior Surface Record and `CONTEXT.md` are
      current.

## Consumers

- The QA prompt contract and the owned skill versions — consumers: `./internal/agent`, `./skills`

## Context

- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `internal/agent/spec_prompt.go`
- interface: `docs/user-guide/commands/implement.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do grep -q '^### Evidence carries no secret$' "$f" || { printf 'missing heading in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'record environment variable names, never their values' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/agent/spec_prompt.go | grep -qF -- 'Record environment variable names, never their values' || { printf 'the QA prompt lacks the rule\n' >&2; exit 1; }; for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md docs/user-guide/commands/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'QA evidence carries a secret' || { printf 'missing refusal in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do grep -qF -- 'secret-found' "$f" || { printf 'missing park in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md docs/user-guide/configuration.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Judge Key Variables' || { printf 'missing variables in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Secret Pattern**:' || { printf 'CONTEXT.md does not define **Secret Pattern**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Secret Match**:' || { printf 'CONTEXT.md does not define **Secret Match**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Secret Scan**:' || { printf 'CONTEXT.md does not define **Secret Scan**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Judge Key Variables**:' || { printf 'CONTEXT.md does not define **Judge Key Variables**\n' >&2; exit 1; }` — expected: exit 0; before this Task none of the phrases and terms exists.
- `grep -q '^### Evidence carries no secret$' .agents/skills/qa-gate/SKILL.md || exit 1; make skills-sync-check || exit 1; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills || exit 1; go test -count=1 -run '^TestBuildQAPrompt' ./internal/agent || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0; before this Task the heading is absent.

## References

- [_prd.md](_prd.md) — Goals; User Story 5; Core Feature 1; Core Feature 6; Glossary
- [_techspec.md](_techspec.md) — API Contract 2; API Contract 4; API Contract 5; Coverage Map; Build Order 4
- ADR-0273; ADR-0256; ADR-0189; ADR-0233; ADR-0244

## Result

Acceptance evidence:

- The canonical `qa-gate` skill and QA prompt carry the evidence rule. A focused byte comparison confirmed that the `### QA settlement` section in the canonical skill is unchanged.
- `make skills-sync` passed, and the canonical and shipped copies of the three Roundfix references plus both owned skill files compare equal. The references and guides carry the QA refusal, `secret-found` park, and Judge Key Variables wording.
- The owned-skill recorder passed and added `qa-gate` `0.0.14` and `roundfix` `0.1.73` to `skills/testdata/owned-skill-versions.json`. `CONTEXT.md` contains the four ADR-0273 glossary definitions. The Behavior Surface Record recorder passed and updated only the three changed guide fingerprints in `docs/references/behavior-surfaces.json`.

Focused checks:

- `GOCACHE=/tmp/roundfix-task04-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed.
- `gofmt -d internal/agent/spec_prompt.go`, `git diff --check`, mirror comparisons, and the QA settlement byte comparison — passed.
- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run '^$' ./internal/agent ./skills` — passed with no tests selected.
- The first composite phrase check exposed the required case-sensitive lowercase `record` wording; after correction and regeneration, the final scoped contract, mirror, settlement, formatting, and path checks passed.
- Verification feedback identified stale fingerprints for the three changed guides; `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — passed.
