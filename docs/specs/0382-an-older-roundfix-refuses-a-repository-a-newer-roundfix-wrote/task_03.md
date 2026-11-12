---
task: task_03
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
status: completed
type: docs
complexity: medium
---

# Task 03: The Roundfix baseline and setup references describe the record and the refusal

## Overview

Both copies of the Roundfix skill's `baseline` and `setup` references tell an
Agent session what the Recorded Roundfix Version is, what an older binary
prints, and that the remedy is `roundfix upgrade`, so a session that meets the
refusal runs the upgrade instead of working around Roundfix. The skill version
rises with the record command. This answers the message reported by conexus on
2026-11-11
(`references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md`).

## Requirements

1. MUST change the canonical `.agents/skills/roundfix/references/baseline.md`
   as `_techspec.md` → API Contract 11 says: right before the heading `##
   Vendored skills`, the heading `### Recorded Roundfix Version` and one
   paragraph whose first sentence is word for word "The Setup Manifest records
   as `roundfixVersion` the release of the Roundfix that last applied a
   Baseline change to the repository." and which says that a plan records the
   running release only when it changes a file and that release is higher,
   that an older binary refuses every command that loads the repository's
   configuration, `roundfix init --scope project` and the six Baseline
   commands with the refusal of API Contract 2, and which ends, word for word,
   "The remedy is `roundfix upgrade`; never edit the field."
2. MUST change the canonical `.agents/skills/roundfix/references/setup.md`:
   right after the paragraph under `## Update check` that begins "`roundfix
   update check [--refresh]", one paragraph whose first sentence is word for
   word "When a newer Roundfix wrote the repository's Baseline, the check
   warns on stderr and reports `Baseline: failed:` with the refusal; run
   `roundfix upgrade` before anything else in that repository." and which
   names the `baseline-version: failed` Doctor line and that the Run readers
   and `roundfix stop` warn and continue without Project Config. MUST NOT add
   or change text inside any `### QA settlement` section, and MUST add no
   capitalized phrase that `CONTEXT.md` does not define.
3. MUST run `make skills-sync` to copy the two references to
   `skills/roundfix/references/baseline.md` and
   `skills/roundfix/references/setup.md`, then raise the Roundfix skill version
   with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   never hand-edit a mirror, a version field or the record. Of the Governed
   Paths the Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/baseline.md`,
   `.agents/skills/roundfix/references/setup.md` and `skills/roundfix/SKILL.md`.
4. MUST, as its last step, re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the `baseline` and `setup` reference paragraphs.
- [ ] Sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the two references carry the prescribed headings and sentences.
- [ ] The Roundfix skill version is recorded, and the skills, documentation and glossary contracts pass.

## Context

- instruction: `docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`
- instruction: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `baseline` and `setup` references and the skill version record — consumers: `./skills`
- The Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'The Setup Manifest records as .roundfixVersion. the release of the Roundfix that last applied a Baseline change to the repository' || { echo "baseline reference unchanged: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'The remedy is .roundfix upgrade.; never edit the field' || { echo "baseline reference remedy missing: $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'When a newer Roundfix wrote the repository.s Baseline, the check warns on stderr' || { echo "setup reference unchanged: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither reference describes the record; afterwards both copies do, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record and the glossary phrase contract pass.

## References

- `_prd.md` → User Story 1; User Story 3; Core Feature 5; Success Metric 7
- `_techspec.md` → API Contract 11; Build Order 3
- ADR-0392

## Result

Skills loaded: implement-task, tech-writer, evidence-gate

- Added the prescribed Recorded Roundfix Version paragraph immediately before `## Vendored skills` in the canonical `baseline` reference, including the higher-version plan rule, the configuration-loading command refusal, `roundfix init --scope project`, all six Baseline commands, the API refusal text, and the exact upgrade remedy.
- Added the prescribed Update Check paragraph immediately after the Update Check paragraph in the canonical `setup` reference, including the `Baseline: failed:` report, `baseline-version: failed` Doctor line, and the Run-reader and `roundfix stop` recovery behavior.
- Focused evidence: `git diff --check` passed; `make skills-sync` completed; `cmp` confirmed both canonical/mirror reference pairs are identical; the required text checks found the heading and sentences in both copies.
- Focused evidence: `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` passed (`1 passed in 1 packages`) and recorded Roundfix skill version `0.1.222` in both skill copies and the owned-skill version record.
- The required Behavior Surface Record command is the final step after this Result update; its fresh outcome is reported in the handoff.
