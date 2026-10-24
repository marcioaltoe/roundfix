---
task: task_03
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: docs
complexity: low
---

# Task 03: The user guide lists the Python profile and the glossary defines the Common Setup Snapshot

## Overview

The user guide's profile list names four built-in profiles. This Task adds
`python` with a paragraph on what it governs, that it fixes no application
surface, which skills it takes and
why no Python skill is dispatched yet, defines **Common Setup Snapshot** in the
glossary, records the Coverage Review the PRD's Skills section plans, and
re-records the Behavior Surface Record last (`_prd.md` → Core Feature 5;
ADR-0322).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected item 1 for the user guide.
2. MUST change the first sentence of the `### Profiles` section of
   `docs/user-guide/context-driven-development.md` to "Roundfix ships
   `go-cli-tui`, `go-cli-typescript-monorepo`, `python`, `rust-cli`, and
   `standard-typescript-monorepo`." and add, after the paragraph about the
   composed profile's root gate, this paragraph: "The `python` Profile
   governs a Python project run with uv, ruff and pytest through
   `docs/agents/python.md`, whether the project is a command-line tool, a
   backend service or a worker. It fixes no application surface: rules for
   one surface belong in the repository's own
   `docs/agents/specific-repository.md`. It takes the `python` setup, a Common Setup
   Snapshot of the `go` and `rust` setups that holds only the skills both
   carry, so no Go, TUI or Rust skill is restored. It dispatches no
   Python-specific skill because the upstream skills repository has none yet.
   Like every profile, its alignment needs the repository's declared `make`
   Verification targets before a plan is ready." The section MUST name no
   command the binary lacks, and the guide MUST NOT contain `python3` or
   `Python fallback`.
3. MUST add **Common Setup Snapshot** to `CONTEXT.md` through
   `domain-modeling`, right after the **Composed Setup Snapshot** entry, as
   "**Common Setup Snapshot**:" followed on the next line by "A setup snapshot
   derived by name from two or more setups that keeps, in the first one's
   order, only the skills and activation bundles every one of them carries
   with equal entries, so a built-in profile can take the language-neutral
   skills of toolchains without its own upstream list", followed by a space
   and a link to ADR-0322 written as the neighbouring entries write theirs
   (link text `ADR-0322`, target
   `docs/adr/0322-a-python-profile-takes-the-skills-the-go-and-rust-setups-share.md`,
   in parentheses) and a closing period,
   and the line "_Avoid_: Intersection setup, shared setup".
4. MUST record the Coverage Review of `_prd.md` → Skills by adding to the
   `guide docs/user-guide/context-driven-development.md` entry of
   `docs/references/skill-coverage.json` the key `review` with the value
   "2026-10-23 — " followed by that surface's reason in `_prd.md` → Skills,
   changing no other entry.
5. MUST, as its LAST step and after the last edit of the user guide,
   re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   which rewrites `docs/references/behavior-surfaces.json`; re-record again if
   the guide changes after it.
6. MUST NOT change any skill, any other user guide page, `README.md`, any Go
   file or any Baseline asset.

## Subtasks

- [ ] Update the profile list and add the paragraph.
- [ ] Add the glossary entry.
- [ ] Record the Coverage Review, then re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The user guide names `python`, says it fixes no application surface, and describes its setup and skills without a retired runtime marker.
- [ ] The glossary defines **Common Setup Snapshot** with its ADR link.
- [ ] The Skill Coverage Map records the review and the Behavior Surface Record is current.

## Consumers

- the user guide and glossary — consumers: `./internal/docscontract`

## Context

- instruction: `docs/adr/0322-a-python-profile-takes-the-skills-the-go-and-rust-setups-share.md`
- instruction: `skills/baseline_skill_contract_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `CONTEXT.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `g="$(tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md)"; printf '%s\n' "$g" | grep -qE -- 'Roundfix ships .go-cli-tui., .go-cli-typescript-monorepo., .python., .rust-cli., and .standard-typescript-monorepo.[.]' || { echo 'profile list lacks python' >&2; exit 1; }; for p in 'Profile governs a Python project run with uv, ruff and pytest through' 'whether the project is a command-line tool, a backend service or a worker. It fixes no application surface' 'setups that holds only the skills both carry, so no Go, TUI or Rust skill is restored.' 'It dispatches no Python-specific skill because the upstream skills repository has none yet.'; do printf '%s\n' "$g" | grep -qF -- "$p" || { printf 'missing in guide: %s\n' "$p" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Common Setup Snapshot**:' || { echo 'missing glossary term' >&2; exit 1; }; grep -qF -- '"review": "2026-10-23 — the only change names the new built-in' docs/references/skill-coverage.json || { echo 'missing coverage review' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract && go test -count=1 -run '^TestNoPythonBaselineRuntime$' ./skills` — expected: exit 0. Before this Task the guide names four profiles, so the command fails at its first check. Afterwards the guide, the glossary and the review are present, the Behavior Surface Record is current and the retired-runtime guard passes.

## References

- `_prd.md` → Core Feature 5; Success Metric 5; Skills; Glossary; User Story 1
- `_techspec.md` → Current behavior 8-9; Build Order 3
- ADR-0322
