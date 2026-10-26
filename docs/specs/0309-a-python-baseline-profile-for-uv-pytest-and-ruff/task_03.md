---
task: task_03
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
status: completed
type: docs
complexity: low
---

# Task 03: The user guide lists the Python profile and the skills it dispatches

## Overview

The user guide's profile list names four built-in profiles. This Task adds
`python` with a paragraph on what it governs, that it fixes no application
surface, where its skills come from and which it dispatches only when the
repository has chosen them, records the Coverage Review the PRD's Skills
section plans, and re-records the Behavior Surface Record last (`_prd.md` →
Core Feature 4; ADR-0322).

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
   `docs/agents/specific-repository.md`. Its `python` setup is the Python
   preset the upstream skills repository publishes, imported by the asset
   sync like every other setup. Its guide makes the gate check the lockfile
   offline and type-check with ty besides ruff and pytest. It dispatches the
   uv, ruff, ty, pytest and Python-services skills whenever their work
   arises, and the Typer, FastAPI, file-pipeline and Streamlit skills only
   when the repository has chosen that tool. Like every profile, its alignment needs the repository's
   declared `make` Verification targets before a plan is ready." The section
   MUST name no command the binary lacks, and the guide MUST NOT contain
   `python3` or `Python fallback`.
3. MUST record the Coverage Review of `_prd.md` → Skills by adding to the
   `guide docs/user-guide/context-driven-development.md` entry of
   `docs/references/skill-coverage.json` the key `review` with the value
   "2026-10-24 — " followed by that surface's reason in `_prd.md` → Skills,
   changing no other entry.
4. MUST, as its LAST step and after the last edit of the user guide,
   re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   which rewrites `docs/references/behavior-surfaces.json`; re-record again if
   the guide changes after it.
5. MUST NOT change any skill, any other user guide page, `README.md`,
   `CONTEXT.md`, any Go file or any Baseline asset.

## Subtasks

- [ ] Update the profile list and add the paragraph.
- [ ] Record the Coverage Review, then re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The user guide names `python`, says it fixes no application surface, where its setup comes from, that its gate runs the offline lock check and ty, and which skills it dispatches only when chosen, without a retired runtime marker.
- [ ] The Skill Coverage Map records the review and the Behavior Surface Record is current.

## Consumers

- the user guide — consumers: `./internal/docscontract`

## Context

- instruction: `docs/adr/0322-a-python-profile-imports-the-published-upstream-python-preset.md`
- instruction: `skills/baseline_skill_contract_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `g="$(tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md)"; printf '%s\n' "$g" | grep -qE -- 'Roundfix ships .go-cli-tui., .go-cli-typescript-monorepo., .python., .rust-cli., and .standard-typescript-monorepo.[.]' || { echo 'profile list lacks python' >&2; exit 1; }; for p in 'Profile governs a Python project run with uv, ruff and pytest through' 'whether the project is a command-line tool, a backend service or a worker. It fixes no application surface' 'preset the upstream skills repository publishes, imported by the asset sync like every other setup.' 'FastAPI, file-pipeline and Streamlit skills only when the repository has chosen that tool.' 'Its guide makes the gate check the lockfile offline and type-check with ty besides ruff and pytest.'; do printf '%s\n' "$g" | grep -qF -- "$p" || { printf 'missing in guide: %s\n' "$p" >&2; exit 1; }; done; grep -qF -- '"review": "2026-10-24 — the only change names the new built-in' docs/references/skill-coverage.json || { echo 'missing coverage review' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract && go test -count=1 -run '^TestNoPythonBaselineRuntime$' ./skills` — expected: exit 0. Before this Task the guide names four profiles, so the command fails at its first check. Afterwards the guide and the review are present, the Behavior Surface Record is current and the retired-runtime guard passes.

## References

- `_prd.md` → Core Feature 4; Success Metric 6; Skills; User Story 1
- `_techspec.md` → Current behavior 10-11; Build Order 3
- ADR-0322

## Result

Implementation-ready. Three files changed; no Go file, skill, other guide page,
`README.md`, `CONTEXT.md` or Baseline asset was touched.

### What changed

- `docs/user-guide/context-driven-development.md` (`### Profiles`): the first
  sentence now reads "Roundfix ships `go-cli-tui`,
  `go-cli-typescript-monorepo`, `python`, `rust-cli`, and
  `standard-typescript-monorepo`." and a new `python` paragraph follows the
  composed profile's root-gate paragraph, naming what it governs, that it fixes
  no application surface, that its `python` setup is the upstream preset
  imported by the asset sync, that its gate runs the offline lock check and ty,
  the always-dispatched and chosen-only skills, and the `make` Verification
  dependency.
- `docs/references/skill-coverage.json`: the
  `guide docs/user-guide/context-driven-development.md` entry gained a `review`
  key with `2026-10-24 — ` followed by that surface's reason in `_prd.md` →
  Skills; no other entry changed.
- `docs/references/behavior-surfaces.json`: re-recorded as the LAST step after
  the last guide edit; the only diff is the guide's fingerprint.

### Acceptance criterion 1 — the guide names python and the retired markers are absent

Focused check: normalized-whitespace grep of the guide confirms the profile list
`Roundfix ships `go-cli-tui`, `go-cli-typescript-monorepo`, `python`,
`rust-cli`, and `standard-typescript-monorepo`.` matches, and each required
paragraph fragment is present verbatim:

- `Profile governs a Python project run with uv, ruff and pytest through`
- `whether the project is a command-line tool, a backend service or a worker. It fixes no application surface`
- `preset the upstream skills repository publishes, imported by the asset sync like every other setup.`
- `FastAPI, file-pipeline and Streamlit skills only when the repository has chosen that tool.`
- `Its guide makes the gate check the lockfile offline and type-check with ty besides ruff and pytest.`

`grep -n 'python3\|Python fallback' docs/user-guide/context-driven-development.md`
returns no match, so the guide carries no retired runtime marker. The section
names no command the binary lacks (only the existing `roundfix baseline profile
show/validate/init` and `roundfix baseline skills restore` examples it already
had).

### Acceptance criterion 2 — the Coverage Review is recorded and the Record is current

Focused check: `grep -qF '"review": "2026-10-24 — the only change names the new built-in' docs/references/skill-coverage.json` matches, and the file
remains valid JSON. After the last guide edit, the Behavior Surface Record was
re-recorded with:

```
go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
```

which exited 0 and rewrote only the guide's fingerprint. `git diff
docs/references/behavior-surfaces.json` shows a single-line change to
`guide docs/user-guide/context-driven-development.md`.

### Focused checks run

- `go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract` — ok (Record current).
- `go test -count=1 -run '^TestNoPythonBaselineRuntime$' ./skills` — ok (retired-runtime guard passes).
- `go test -count=1 -tags docscontract ./internal/docscontract` — ok (full docscontract suite).
- `go test -count=1 ./skills` — ok (full skills suite).
- The task's full `## Verification` command, run as a focused check, exited 0.

### Follow-ups (out of this Task's slice)

- None. The remaining `python` profile behavior lives in task_01/task_02 and the
  Starter files in Spec 0310.
