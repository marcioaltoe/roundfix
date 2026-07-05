---
task: task_05
spec: 0004-watch-merge-readiness
status: completed
type: docs
complexity: low
---

# Task 05: Sync docs and the Roundfix skill with the watch contract

## Overview

Document the shipped behavior: the merge-readiness Clean semantics, the
watch/resolve stdout report shape, and the `--no-agent-console` flag — in the
canonical Roundfix skill and command help, with the embedded copy regenerated.
Verifiable through the skills drift check inside the full gate.

## Requirements

1. MUST update the canonical Roundfix skill: watch's Clean now means the
   Review Source check succeeded on the final pushed commit (ADR-0019,
   including the `missing` note), the exact stdout report shapes for watch
   and resolve, and `--no-agent-console` on the operational commands;
   regenerate the embedded copy through the sync target.
2. MUST verify each documented flag and line shape against the built binary's
   actual output.
3. MUST verify every term against the glossary; call out gaps instead of
   inventing language (candidate gap to flag if felt: a term for
   merge-readiness).

## Subtasks

- [x] Skill updates + `make skills-sync`
- [x] Help-text cross-check against shipped output
- [x] Glossary pass

## Acceptance Criteria

- [x] Skill text matches shipped behavior exactly; drift check passes inside
      the full gate.
- [x] Every documented stdout line shape appears verbatim in a CLI test
      fixture.
- [x] No new un-glossaried term, or the gap is called out in the Result.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Features 2–4; User Experience. `_techspec.md` → Build Order
5. ADR-0019. Repo hard rule (canonical skill ships with CLI behavior
changes).

## Result

Status: completed.

Acceptance evidence:

- Skill sync: `.agents/skills/roundfix/SKILL.md` documents the Watch Run Clean
  contract, the missing-check stderr note, the watch/resolve stdout report
  examples, and `--no-agent-console`; `rtk make skills-sync` regenerated
  `skills/roundfix/SKILL.md`, and `diff -r .agents/skills/roundfix
  skills/roundfix` returned no drift.
- Help-text cross-check: the built binary output from `rtk ./bin/roundfix
  watch --help` includes `--until-clean  Repeat until no Unresolved Review
  Issues remain and Review Source check succeeds`; `resolve --help`,
  `watch --help`, and `implement --help` all include `--no-agent-console`.
  `TestRunCommandHelp` asserts the watch help wording.
- Stdout shape cross-check: the documented report examples
  `issue 001 resolved — major: handle test issue`, `Clean after 1 Round(s):
  1 resolved, 0 invalid, 0 failed, 0 unresolved.`, and `TimedOut after 0
  Round(s): 0 resolved, 0 invalid, 0 failed, 0 unresolved.` appear verbatim
  in `internal/cli/cli_test.go`.
- Glossary pass: the skill uses glossary terms such as Run, Watch Run, Clean,
  Review Source, Review Issue, Unresolved Review Issue, Final Push, Max
  Rounds, Daemon, Agent, Run Event Journal, Interactive Input, and Live Run
  View. Gap called out: `merge-readiness` / `merge-ready` appears in this
  Spec and ADR-0019, but not in `CONTEXT.md`; the skill avoids introducing
  that as a new product term.
- Verification passed: `rtk go run ./cmd/roundfix skills check` and
  `rtk make verify`.
