---
task: task_02
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
status: pending
type: docs
complexity: medium
---

# Task 02: A fresh home keeps the host toolchains in the QA skills

## Overview

The QA gates of Specs 0399 to 0402 ran test suites under a fresh `HOME` that hid the host's Rust toolchain, because nothing tells a gate or a plan author what a fresh home must keep. This Task adds one section to the `qa-gate` skill and one to the `write-tasks` skill, word for word as `_techspec.md` → API Contract 2 gives them, and a contract test that pins each guarantee in both copies. It answers the maintainer's Backlog Entry of 2026-11-18, [`references/2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md`](references/2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md), symptom 1.

## Requirements

1. MUST add to `.agents/skills/qa-gate/SKILL.md`, right before the heading `### Evidence links resolve in the committed tree`, the section `### A fresh home keeps the host toolchains` of API Contract 2, word for word, as a plain Markdown heading and its three paragraphs. The section stays outside `### QA settlement`, and no text inside `### QA settlement` changes.
2. MUST add to `.agents/skills/write-tasks/SKILL.md`, right before the heading `## Verification paths`, the section `## QA environment` of API Contract 2, word for word, as a plain Markdown heading and its paragraph.
3. MUST NOT change any other sentence of either skill, any other owned skill, or any Vendored Skill.
4. MUST add `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome` in the new file `skills/qa_fresh_home_test.go`, calling `t.Parallel()` first. It reads both copies through `roundfixLayoutSources` and runs, for each copy and skill, a subtest `<copy>/<skill>` that requires the section's heading exactly once and cuts the section at the next heading of the same level. Inside it, one subtest per guarantee compares whitespace-normalized text and requires the phrase exactly once:
   - for `qa-gate`, each of the seven sentences of the section's second paragraph, the two sentences of its first paragraph and the sentence of its third paragraph, named `suites_keep_the_session_home`, `a_product_row_may_get_a_fresh_home`, `locations_read_before_home_changes`, `go_locations`, `rust_locations`, `another_stack`, `never_installs_a_toolchain`, `a_home_prohibition_lets_them_pass`, `created_outside_the_checkout` and `a_hidden_toolchain_is_no_environment_block`;
   - for `write-tasks`, each of its five sentences, named `the_gate_keeps_the_toolchains`, `suites_keep_the_session_home`, `the_locations_are_passed`, `a_home_prohibition_names_its_exception` and `a_suite_states_why`;
   - `before_the_qa_settlement_section`, which for `qa-gate` requires the section to start before `### QA settlement`.
5. MUST regenerate, never by hand: the mirrors under `skills/` with `make skills-sync`, then the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of each changed skill, then `make skills-sync` again so each mirror equals its source.
6. MUST NOT change `docs/references/behavior-surfaces.json`, `CONTEXT.md`, `internal/`, the Makefile, a workflow, `go.mod` or a lint configuration. Of the Governed Paths it changes only the four `SKILL.md` files below, all listed in `_authorization.md`.

## Subtasks

- [ ] Add the `qa-gate` section and the `write-tasks` section.
- [ ] Add the contract test.
- [ ] Sync the mirrors and record the owned-skill versions.

## Acceptance Criteria

- [ ] Both copies of both skills hold their section exactly once, with every guarantee, and the `qa-gate` section sits before `### QA settlement`.
- [ ] Removing any pinned sentence fails its named subtest.
- [ ] The mirrors equal their sources and the new versions are recorded.

## Context

- instruction: `skills/phase_boundary_test.go`
- instruction: `skills/roundfix_layout_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/qa_fresh_home_test.go`

## Consumers

- The `qa-gate` and `write-tasks` skill text and the owned-skill version record — consumers: `./skills`
- The `write-tasks` skill text, read by the public documentation contracts — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/suites_keep_the_session_home TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/a_product_row_may_get_a_fresh_home TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/locations_read_before_home_changes TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/go_locations TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/rust_locations TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/another_stack TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/never_installs_a_toolchain TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/a_home_prohibition_lets_them_pass TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/created_outside_the_checkout TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/a_hidden_toolchain_is_no_environment_block TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/before_the_qa_settlement_section TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/write-tasks/the_gate_keeps_the_toolchains TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/write-tasks/suites_keep_the_session_home TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/write-tasks/the_locations_are_passed TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/write-tasks/a_home_prohibition_names_its_exception TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/write-tasks/a_suite_states_why TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/suites_keep_the_session_home TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/a_product_row_may_get_a_fresh_home TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/locations_read_before_home_changes TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/go_locations TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/rust_locations TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/another_stack TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/never_installs_a_toolchain TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/a_home_prohibition_lets_them_pass TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/created_outside_the_checkout TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/a_hidden_toolchain_is_no_environment_block TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/before_the_qa_settlement_section TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/write-tasks/the_gate_keeps_the_toolchains TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/write-tasks/suites_keep_the_session_home TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/write-tasks/the_locations_are_passed TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/write-tasks/a_home_prohibition_names_its_exception TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/write-tasks/a_suite_states_why; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; make skills-sync-check && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the test and both sections do not exist, so the pass lines are missing; afterwards every guarantee passes in both copies, the mirrors match, the whole skills suite, with the recorded versions, passes, and the documentation contracts, which forbid profile policy terms in both copies of `write-tasks`, pass.

## References

- `_prd.md` → Goals (first and second bullets); Core Features 1 and 2; Success Metric 1
- `_techspec.md` → API Contract 2; Coverage Map; Testing Approach; Build Order 2
- [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md), section 1
- ADR-0080
