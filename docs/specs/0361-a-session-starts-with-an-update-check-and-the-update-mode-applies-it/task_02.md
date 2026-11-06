---
task: task_02
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
status: pending
type: backend
complexity: high
---

# Task 02: roundfix update check reports a newer release and the Baseline plan

## Overview

A new read-only command, `roundfix update check`, is the Update Check. It
reads one `gh api` release listing, cached for 24 hours in Roundfix Home,
compares the installed Roundfix with the latest release, prints the `### `
titles of each newer release, summarizes this repository's `roundfix baseline
update` preview run in process, and ends with a `Next:` line that names the
`--update` Skill Mode and the configured `update.mode`. It always exits `0` for
a report. The command guide, the Roundfix `setup` reference, the command
indexes, the Skill Coverage Map and the glossary describe it.

## Requirements

1. MUST add to `internal/app/release.go` the `Release` type,
   `ReleaseListingPageSize` and `ListReleases` of `_techspec.md` →
   Interfaces and API Contract 2, through the existing `runGitHubAPI`, and
   add `TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases` and
   `TestListReleasesMapsGitHub404ToNoReleases` to
   `internal/app/release_test.go` with the existing `withGitHubAPIRunner`
   helper, asserting the endpoint string, the order, the titles, a body
   without a title and the skipped draft and prerelease.
2. MUST add `internal/cli/update_check.go` with the `update` group and its
   `check` subcommand, modeled on the `storage` group (`runStorageCommand`):
   `--refresh` and `--format <text|json>`, the cache of API Contract 3 in
   `<home>/.roundfix/update-check.json` under a five-second lookup timeout,
   the text report of API Contract 4, the JSON document of API Contract 5,
   the Baseline line of API Contract 6 through `runBaselineUpdateCommand` in
   process with `--repo <git root> --format json` and its standard error
   discarded, and the usage and failures of API Contract 7. Every line MUST
   be written word for word as those contracts and Surface Transcripts 1 to 5
   give it. The preview MUST NOT receive `--yes`, `--confirm-plan`,
   `--adopt-suggested` or `--decision` (Invariant 6), and the command MUST NOT
   write any file other than the cache (Invariant 1).
3. MUST add `updateCheck updateCheckDependencies` to `commandDependencies` in
   `internal/cli/cli.go`, defaulted like `versionFreshness`, dispatch
   `update` to the new command, add the `update` and `update check` usage of
   Surface Transcript 1 to `commandUsage`, and add the root usage line and the
   command row of API Contract 7. MUST add an `update check` row to
   `TestRunCommandHelp` in `internal/cli/cli_test.go` and update every other
   test that pins the root usage on the default branch when the Task runs.
4. MUST add `internal/cli/update_check_test.go` with
   `TestUpdateCheckTranscripts`, `TestUpdateCheckReadsTheNetworkAtMostOncePerDay`,
   `TestUpdateCheckSummarizesTheBaselinePreview` and `TestUpdateCheckUsage`
   with every row `_techspec.md` → Testing Approach names, a fake release
   lister and a fixed clock injected through the command-dependency context,
   temporary homes, and `newBaselineUpdateRepository` and
   `staleBaselineUpdateRepository` for the Baseline rows, following the
   package's rule on `t.Parallel()`. `TestUpdateCheckTranscripts` asserts the
   whole standard output of each row, with Surface Transcripts 2 to 4 among
   them, empty standard error and exit `0`; no test reaches the network or
   runs `gh`.
5. MUST create `docs/user-guide/commands/update.md` and add the
   `docs/user-guide/commands.md` and `docs/user-guide/usage.md` rows as API
   Contract 8 says, the hook paragraph word for word; and add to the
   canonical `.agents/skills/roundfix/references/setup.md`, before `##
   Pending History after upgrade`, the `## Update check` section of API
   Contract 8 word for word, and `` `update` `` after `` `upgrade` `` in the
   `setup` row of the index table in `.agents/skills/roundfix/SKILL.md`.
6. MUST add two entries to `docs/references/skill-coverage.json`:
   `command update check` with skills
   `.agents/skills/roundfix/references/setup.md` and sources
   `docs/user-guide/commands/update.md` and `internal/cli/update_check.go`,
   and `guide docs/user-guide/commands/update.md` with the same skill and its
   own guide as source.
7. MUST add through `domain-modeling` the `CONTEXT.md` entry **Update Check**
   of API Contract 10, word for word, after **Upgrade Command**, adding no
   other capitalized phrase.
8. MUST run `make skills-sync`, then raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide, help and skill edit; never
   hand-edit a mirror, a version field or a record. Of the Governed Paths the
   Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/setup.md`, `internal/cli/cli_test.go`
   and `skills/roundfix/SKILL.md`. Spec 0360 changes other sections of the
   `setup` reference and `SKILL.md`; keep these edits to the new section and
   the one index cell.

## Subtasks

- [ ] Add `ListReleases` and its tests.
- [ ] Add the `update check` command, its cache, reports and in-process preview.
- [ ] Wire the dispatch, usage, root usage and help test.
- [ ] Add the command tests.
- [ ] Write the guide, the index rows, the `setup` section, the coverage entries and the glossary entry; sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] With a newer release the report lists every newer release newest first with its titles, the Baseline line and the `Next:` line for the configured mode, and exits `0`.
- [ ] Two checks within 24 hours make one lookup unless `--refresh`; a failed lookup is one line, exits `0` and is cached.
- [ ] The Baseline line reports current, a plan with its counts, another state, or not checked outside a repository, and the check writes no repository file.
- [ ] The guide, the indexes, both copies of the `setup` reference, the coverage map and `CONTEXT.md` describe the command.

## Context

- instruction: `docs/adr/0371-a-session-starts-with-an-update-check-and-the-skill-applies-the-update.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/cli/upgrade.go`
- instruction: `internal/cli/gc.go`
- instruction: `internal/cli/baseline_update.go`
- instruction: `internal/cli/baseline_update_test.go`
- interface: `internal/app/release.go`
- interface: `internal/app/release_test.go`
- creates: `internal/cli/update_check.go`
- creates: `internal/cli/update_check_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- creates: `docs/user-guide/commands/update.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/skill-coverage.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `app.ListReleases` and `app.Release` — consumers: `./internal/app`, `./internal/cli`
- The `update check` command, its root usage line and help — consumers: `./internal/cli`
- The guide, the `setup` reference, the Roundfix skill version, the coverage map and the glossary — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases|TestListReleasesMapsGitHub404ToNoReleases)$' ./internal/app 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases TestListReleasesMapsGitHub404ToNoReleases; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/app && GOOS=windows go build -buildvcs=false ./internal/app` — expected: exit 0. Before this Task the tests do not exist; afterwards the release listing is parsed as specified and the app suite passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestUpdateCheckTranscripts|TestUpdateCheckReadsTheNetworkAtMostOncePerDay|TestUpdateCheckSummarizesTheBaselinePreview|TestUpdateCheckUsage|TestRunCommandHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUpdateCheckTranscripts TestUpdateCheckReadsTheNetworkAtMostOncePerDay TestUpdateCheckSummarizesTheBaselinePreview TestUpdateCheckUsage; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/cli` — expected: exit 0. Before this Task the command tests do not exist, so their pass lines are missing; afterwards the command, its cache, reports and usage pass and the CLI suite passes and builds for Windows.

- `test -f docs/user-guide/commands/update.md || { echo 'missing update guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qF -- 'A hook is repository tooling configuration, so add it only with the maintainer' || { echo 'missing hook paragraph' >&2; exit 1; }; for f in .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do grep -q '^## Update check$' "$f" || { echo "missing section: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -q -- 'listing is one .gh api. call cached for 24 hours in' || { echo "missing cache sentence: $f" >&2; exit 1; }; done; grep -q '^| .update. | \[command guide\](commands/update.md) |$' docs/user-guide/commands.md || { echo 'missing commands row' >&2; exit 1; }; grep -F -- '**Update Check**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Update Check\*\*(:| —( |$))' || { echo 'missing term: Update Check' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no guide, section, row or glossary entry exists; afterwards they do, and the documentation contracts, the coverage map, the Behavior Surface fingerprint, the mirrors and the recorded skill version pass.

## References

- `_prd.md` → Goals 1 and 2; Core Features 2 and 5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 7; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior 1, 2 and 4; Interfaces; Data Models; API Contracts 2 to 7; API Contract 8; API Contract 10; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Invariants 1 to 4 and 6; Declared breaks; Testing Approach; Build Order 2
- ADR-0371
