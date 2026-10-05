---
task: task_01
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
status: pending
type: docs
complexity: low
---

# Task 01: The guides and the Roundfix Skill describe the credit floor, the named refusal and the relay

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship
the Roundfix Skill update. This Task describes, in the configuration guide
and in the Roundfix Skill's `runtime` reference, the behavior task_02 to
task_04 implement, as the TechSpec states it: the User Config value
`jev.router_min_credit_usd` with its US$15 default, the gate's account-credit
read and its refusal `openrouter_credit_low`, the named OpenRouter refusal
`openrouter_credit_refused`, the loopback relay, and the model, provider and
response id on each `router-prompt` Judge Log line. It also adds to ADR-0218
the note that ADR-0234 supersedes it in part.

## Requirements

1. MUST describe in `docs/user-guide/configuration.md`, in a section of its
   own after "Jev monthly ceiling", the key `jev.router_min_credit_usd`: a
   finite number of US dollars greater than zero, set only in User Config at
   `~/.roundfix/config.yml`, US$15 when unset; that the Jev Router gate
   compares it with the lower of the OpenRouter account's balance (total
   credits less total usage) and the key's remaining limit; the Project
   Config warning of API Contract 2 and the error of API Contract 3, each
   quoted verbatim; an example setting it to `20`; and that an older Roundfix
   binary refuses a User Config holding the key, so the operator sets it only
   after every Roundfix binary on the machine includes this change. MUST add
   a `jev.router_min_credit_usd` row with default `15` to its Key reference.
2. MUST extend the Jev Router passages of `docs/user-guide/configuration.md`
   and `.agents/skills/roundfix/references/runtime.md` to state: that the
   gate also reads `GET /api/v1/credits` before each routed prompt, naming
   its fields `total_credits` and `total_usage`, and refuses
   below the floor with API Contract 4, quoting its code
   `openrouter_credit_low`; that an unreadable credits answer is
   `jev_spend_unreadable`; that an OpenRouter HTTP 402 on a routed request is
   `openrouter_credit_refused` (API Contract 5), a failed selection before
   Agent work begins and a failed Work Item after; that routed sessions reach
   OpenRouter through a relay Roundfix runs on the loopback interface, which
   forwards requests unchanged and never logs or stores the key; and that
   each `router-prompt` line records the models and providers the router
   reported and the last response id (API Contract 6).
3. MUST append to `docs/adr/0218-the-jev-router-is-a-project-selected-opencode-model-under-the-shared-jev-ceiling.md`,
   after its existing ceiling note, a paragraph beginning
   `**Relay and credit (2026-10-04).** Superseded in part by ADR-0234:` that
   says the inline provider now points at the loopback relay and the gate
   also refuses below the account credit floor, and that every other part of
   the decision stands; and MUST update its `updated_at`.
4. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running that record command is an implementation step of this Task, not
   part of its Verification: run it yourself after the last skill edit. Never
   write a digest by hand. The record flag never replaces a recorded digest, so
   if an entry for the raised version already exists with another digest,
   delete that entry and run the record command again.
5. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference or guide, `.roundfixrc.yml`, ADR-0234, `CONTEXT.md`, or
   any file outside the repository.

## Subtasks

- [ ] Describe the floor in the configuration guide and its Key reference.
- [ ] Extend the Jev Router passages in the configuration guide and the `runtime` reference.
- [ ] Add the supersession note to ADR-0218.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The configuration guide carries `jev.router_min_credit_usd`, the API
      Contract 2 warning, the API Contract 3 error, `openrouter_credit_low`
      and `openrouter_credit_refused`.
- [ ] The `runtime` reference names the floor's key, both reasons and the
      relay; ADR-0218 names ADR-0234.
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/adr/0218-the-jev-router-is-a-project-selected-opencode-model-under-the-shared-jev-ceiling.md`
- instruction: `docs/adr/0234-the-jev-router-gate-reads-the-account-credit-and-a-relay-names-the-routed-model.md`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "jev.router_min_credit_usd must be a finite number greater than 0" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "API Contract 3 error" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "config: jev.router_min_credit_usd in Project Config is ignored; set jev.router_min_credit_usd in User Config" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "API Contract 2 warning" >&2; exit 1; }; for phrase in jev.router_min_credit_usd openrouter_credit_low openrouter_credit_refused total_credits loopback; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/runtime.md "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/adr/0218-the-jev-router-is-a-project-selected-opencode-model-under-the-shared-jev-ceiling.md | grep -qF -- "Superseded in part by ADR-0234:" || { printf 'missing ADR-0234 note in ADR-0218\n' >&2; exit 1; }` — expected: exit 0; before this Task no file names the floor, either reason, the credits fields or the relay, and ADR-0218 does not name ADR-0234, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/runtime.md | grep -qF -- "openrouter_credit_refused" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirror does not name the refusal, so the command fails; after it the mirrors equal their canonical files and the raised version and its content digest are recorded.

## References

- `_prd.md` → Core Feature 6; User Stories 1-5; Success Metric 4
- `_techspec.md` → Vocabulary Contract; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Build Order 1
- ADR-0187; ADR-0189; ADR-0218; ADR-0234
