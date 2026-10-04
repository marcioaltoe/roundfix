---
task: task_01
spec: 0225-a-jev-ceiling-the-maintainer-sets
status: completed
type: docs
complexity: low
---

# Task 01: The guides and the Roundfix Skill describe the configured Jev ceiling

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship
the Roundfix Skill update. This Task describes, in the configuration and
`spec` guides and in the Roundfix Skill's `spec` and `runtime` references, the
behavior task_02 implements, as the TechSpec states it: the User Config value
`jev.monthly_ceiling_usd`, its US$5 default, that Project Config cannot set
it, that the judge and the Jev Router gate and its key-limit check read it,
and when the operator sets it. It removes every statement of US$5 as the
fixed ceiling from these files.

## Requirements

1. MUST describe in `docs/user-guide/configuration.md`, in a section of its
   own near "Active Implement Run ceiling", the key `jev.monthly_ceiling_usd`:
   a finite number of US dollars greater than zero, set only in User Config
   at `~/.roundfix/config.yml`, US$5 when unset; that a Project Config value
   is ignored with the warning of API Contract 2, quoted verbatim; the error
   of API Contract 3, quoted verbatim; that the judge, the Jev Router gate and
   its key-limit check read it and that the key's monthly limit must be at
   most that value; an example setting it to `50`; and that an older Roundfix
   binary refuses a User Config holding the key, so the operator sets it only
   after every Roundfix binary on the machine includes this change. MUST add a
   `jev.monthly_ceiling_usd` row to its Key reference.
2. MUST rewrite the Jev Router passages of `docs/user-guide/configuration.md`
   and `.agents/skills/roundfix/references/runtime.md` so the shared ceiling
   and the key's monthly limit are the configured ceiling, with US$5 named
   only as the default.
3. MUST rewrite the judge's ceiling sentence in
   `docs/user-guide/commands/spec.md` and
   `.agents/skills/roundfix/references/spec.md` so the monthly ceiling is
   `jev.monthly_ceiling_usd` in User Config, US$5 by default, across both
   transports and every repository using that Roundfix Home.
4. MUST leave none of these phrases in the four files above:
   `ceiling is US$5.00`, `US$5 monthly Jev ceiling`, `at most US$5 on the key`.
5. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running that record command is an implementation step of this Task, not
   part of its Verification: run it yourself after the last skill edit. Never
   write a digest by hand. The record flag never replaces a recorded digest, so
   if an entry for the raised version already exists with another digest,
   delete that entry and run the record command again.
6. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference or guide, `commands.md`, `.roundfixrc.yml`, the ADRs or
   `CONTEXT.md`, and MUST NOT write any file outside the repository.

## Subtasks

- [ ] Describe the key in the configuration guide and its Key reference.
- [ ] Rewrite the Jev Router passages in the configuration guide and the `runtime` reference.
- [ ] Rewrite the judge's ceiling sentence in the `spec` guide and reference.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The configuration guide carries `jev.monthly_ceiling_usd`, the API
      Contract 2 warning and the API Contract 3 error.
- [ ] The `spec` guide and both references carry `jev.monthly_ceiling_usd`;
      none of the four files states US$5 as the fixed ceiling.
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/spec.md`
- instruction: `docs/adr/0231-the-jev-ceiling-is-a-user-config-value.md`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "jev.monthly_ceiling_usd must be a finite number greater than 0" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "API Contract 3 error" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "config: jev.monthly_ceiling_usd in Project Config is ignored; set jev.monthly_ceiling_usd in User Config" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "API Contract 2 warning" >&2; exit 1; }; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md .agents/skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "jev.monthly_ceiling_usd" || { printf 'missing phrase in %s: %s\n' "$f" "jev.monthly_ceiling_usd" >&2; exit 1; }; done; for f in docs/user-guide/configuration.md docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md .agents/skills/roundfix/references/runtime.md; do text="$(tr -s '[:space:]' ' ' < "$f")"; case "$text" in *'ceiling is US$5.00'*|*'US$5 monthly Jev ceiling'*|*'at most US$5 on the key'*) printf 'fixed ceiling left in %s\n' "$f" >&2; exit 1;; esac; done` — expected: exit 0; before this Task no file names the key and all four state US$5 as the fixed ceiling, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/spec.md | grep -qF -- "jev.monthly_ceiling_usd" && tr -s '[:space:]' ' ' < skills/roundfix/references/runtime.md | grep -qF -- "jev.monthly_ceiling_usd" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirrors do not name the key, so the command fails; after it the mirrors equal their canonical files and the raised version and its content digest are recorded.

## References

- `_prd.md` → Core Feature 6; User Stories 1-5; Success Metric 3
- `_techspec.md` → Vocabulary Contract; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Build Order 1
- ADR-0187; ADR-0189; ADR-0231

## Result

Implemented the documentation slice for the configured Jev ceiling. The
configuration guide now documents `jev.monthly_ceiling_usd` as a User Config
value at `~/.roundfix/config.yml`, its US$5 unset default, Project Config
warning, invalid-value error, consumers, key-limit rule, `50` example, older
binary compatibility warning, and Key reference row. The `spec` guide and
Roundfix Skill `spec` and `runtime` references now describe the configured
ceiling and its default. The canonical Roundfix Skill version is `0.1.26` in
both front-matter fields; mirrors were synchronized and the recorded content
digest was generated by the required test command.

Focused implementation checks:

- `rtk proxy make skills-sync`: passed (exit 0).
- `env GOCACHE=/tmp/roundfix-0225-task01-gocache rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed (exit 0).
- `git diff --check`: passed (exit 0).
- Targeted scan for the three forbidden fixed-ceiling phrases: passed (none found).
- `diff -q` checks for the canonical and mirrored Skill files: passed (all equal).

Acceptance evidence:

- Configuration guide criterion: the new Jev monthly ceiling section contains
  both API Contract 2 and API Contract 3 messages verbatim and the Key
  reference row.
- Guide and reference criterion: all four target documents name
  `jev.monthly_ceiling_usd`, describe US$5 only as the default, and contain no
  forbidden fixed-ceiling phrase.
- Mirror and version criterion: canonical and mirror files compare equal, and
  `owned-skill-versions.json` records version `0.1.26` with the generated
  digest.

The Daemon must run the declared Verification commands and own Task status and
settlement.


Operator addendum (2026-10-04): after this Task settled, 0226 merged with the
Roundfix Skill at `0.1.26`; merging main into the item raised this Spec's
Skill change to `0.1.27` and recorded that version.

## Carry-forward provenance

- Source Run: `run_20261004T194332Z_e41f982109e26c75`
- Source commit: `16880f755621421b763ce11522b5a4162b963e75`
