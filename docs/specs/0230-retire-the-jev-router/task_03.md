---
task: task_03
spec: 0230-retire-the-jev-router
status: pending
type: docs
complexity: medium
---

# Task 03: The guides, the Roundfix Skill and the Project Config comment state the subscription rule instead of the router

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship
the Roundfix Skill update. This Task replaces, in the configuration guide and
the Roundfix Skill's `runtime` reference, every passage that describes the
Jev Router as current (its selection, key, gate, credit floor, refusal codes,
relay and Judge Log lines) with the subscription rule and its refusals, as
the TechSpec states them; adds the rule to the model reference and to the
repository's Project Config comment; and raises the skill's version. The
ADR lifecycle (ADR-0218 and ADR-0234 retired, ADR-0231 noted) was done when
the Spec was authored and is not this Task's work.

## Requirements

1. MUST replace the "Jev Router" section and the "Jev Router credit floor"
   section of `docs/user-guide/configuration.md` with one section, after the
   Agent selection profiles introduction, that states the rule
   "OpenAI and Anthropic models run only through the codex and claude
   subscriptions", which `opencode` selections are refused (the retired
   `roundfix-openrouter` provider; under `openrouter`, the authors `openai`
   and `anthropic`, the router authors `openrouter` and `typesafe`, and `@`
   presets), that the refusal holds in every scope including
   `roundfix profiles configure`, quotes the message of API Contract 2
   verbatim with `<field>` and `<model>` placeholders, and names the runner
   refusal `subscription_only`; and that other OpenRouter models stay
   selectable.
2. MUST, in the same guide, add `jev.router_min_credit_usd` to the list of
   deprecated keys, saying it was removed with the Jev Router, quoting the
   warning of API Contract 5 verbatim; remove its Key reference row; say in
   the Jev monthly ceiling section that the judge reads the ceiling, without
   naming a router gate or key-limit check; and say that `router-prompt`
   lines written before the retirement stay in the Judge Log and count toward
   their month's ceiling.
3. MUST replace the "Jev Router" section of
   `.agents/skills/roundfix/references/runtime.md` with a section stating the
   rule, the refused selections including
   `roundfix-openrouter/typesafe/jev-router`, the reason `subscription_only`
   and that a refusal before Agent work activates the fallback, and that
   `jev.router_min_credit_usd` is a deprecated key.
4. MUST add to `docs/references/model-selection.md`, right after the table
   that maps OpenRouter identifiers to `opencode` selections, a dated note
   whose text includes "Since ADR-0235, Roundfix refuses the OpenAI and
   Anthropic rows of this table" and names the rule.
5. MUST extend the comment of `.roundfixrc.yml` that says never to route a
   `gpt-*-sol` model through OpenRouter with the sentence "Roundfix refuses
   it: OpenAI and Anthropic models run only through the codex and claude
   subscriptions (ADR-0235)." and MUST NOT change any configuration value in
   that file.
6. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running that record command is an implementation step of this Task, not
   part of its Verification: run it after the last skill edit and never write
   a digest by hand. If an entry for the raised version already exists with
   another digest, delete that entry and run the record command again.
7. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference or guide, any ADR, `CONTEXT.md`, `CHANGELOG.md`, the
   judge's documentation or any file outside the repository.

## Subtasks

- [ ] Replace the router sections of the configuration guide with the rule.
- [ ] Register the retired key in the guide's deprecated keys.
- [ ] Replace the router section of the `runtime` reference.
- [ ] Add the note to the model reference and the sentence to the Project Config comment.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The configuration guide and the `runtime` reference state the rule,
      the refused selections and `subscription_only`, and no longer describe
      the router's gate, floor, relay or refusal codes.
- [ ] The guide quotes API Contracts 2 and 5 verbatim.
- [ ] The model reference and the Project Config comment state the rule, and
      no Project Config value changed.
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- instruction: `docs/adr/0235-openai-and-anthropic-models-run-only-through-the-codex-and-claude-subscriptions.md`
- instruction: `docs/adr/0187-the-roundfix-skill-and-the-command-reference-are-read-one-command-at-a-time.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/model-selection.md`
- interface: `.roundfixrc.yml`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for phrase in "OpenAI and Anthropic models run only through the codex and claude subscriptions" subscription_only jev.router_min_credit_usd; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/runtime.md "$phrase" >&2; exit 1; }; done; ! grep -qE 'openrouter_credit_low|openrouter_credit_refused|jev_router_key_unbounded|jev_router_key_missing|api/v1/credits|^#+[[:space:]]Jev[[:space:]]Router' docs/user-guide/configuration.md .agents/skills/roundfix/references/runtime.md || { printf 'router passage left in the guide or the runtime reference\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "can reach OpenAI or Anthropic models through OpenRouter; OpenAI and Anthropic models run only through the codex and claude subscriptions" || { printf 'missing API Contract 2 message\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "config: jev.router_min_credit_usd is deprecated and ignored; the Jev Router was retired, so remove it" || { printf 'missing API Contract 5 warning\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "roundfix-openrouter/typesafe/jev-router" || { printf 'runtime reference does not name the refused router selection\n' >&2; exit 1; }` — expected: exit 0; before this Task neither file states the rule or `subscription_only`, and both describe the router's refusal codes, so the command fails.
- `tr -s '[:space:]' ' ' < docs/references/model-selection.md | grep -qF -- "Since ADR-0235, Roundfix refuses the OpenAI and Anthropic rows of this table" || { printf 'missing model reference note\n' >&2; exit 1; }; tr -d '#' < .roundfixrc.yml | tr -s '[:space:]' ' ' | grep -qF -- "Roundfix refuses it: OpenAI and Anthropic models run only through the codex and claude subscriptions (ADR-0235)." || { printf 'missing Project Config comment\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/roundfix/references/runtime.md | grep -qF -- "subscription_only" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task the model reference has no note and the Project Config comment no rule, so the command fails; after it the mirrors equal their canonical files and the raised version is recorded.

## References

- `_prd.md` → User Story 5; Core Feature 6; Success Metric 5
- `_techspec.md` → API Contract 2; API Contract 5; Vocabulary Contract; Build Order 1
- ADR-0235; ADR-0187; ADR-0189; ADR-0027
