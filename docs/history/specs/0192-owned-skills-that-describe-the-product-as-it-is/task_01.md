---
task: task_01
spec: 0192-owned-skills-that-describe-the-product-as-it-is
status: completed
type: test
complexity: medium
---

# Task 01: Every command the CLI lists is named in the Roundfix skill

## Overview

The Roundfix skill is what an Agent reads to drive the CLI, and the binary's embedded copy is authoritative for adopters. It has no text for seven command forms the root help lists, omits the `agent-selection` event category, and describes a delivery order without the pre-PR review. This Task adds a repository test that reads the command list from the CLI itself and fails for a command the skill does not name. It then brings the skill to the point where that test passes. Review this Task by the test first: it is the part that keeps the skill true after this Spec.

## Requirements

1. MUST add `internal/docscontract/command_documentation_test.go`, in the package and under the `docscontract` build tag the directory's other tests use, with the two helpers the TechSpec's Interfaces section names:
   - `commandPaths(help string) []string` follows "The command path rule" in `_techspec.md` exactly;
   - `undocumentedCommands(paths []string, text string) []string` returns each path whose literal `roundfix <path>` is absent from `text`.
2. MUST add these tests to that file:
   - `TestCommandPathsReadSubcommandsAndAlternatives` feeds a literal help text. It expects `window set`, `window show` and `window clear` from `roundfix window <set|show|clear>`, `archive` from `roundfix archive <slug>`, `baseline capabilities check` from its usage line, one entry for a path two usage lines repeat, and nothing from `roundfix --help`.
   - `TestEveryCommandIsNamedInTheRoundfixSkill` reads the root help through `cli.Run([]string{"--help"}, …)`. It fails when `commandPaths` returns fewer than forty paths, or when `deliver retry`, `window clear` or `baseline capabilities check` is not among them. It then fails with the list of paths `undocumentedCommands` returns for `.agents/skills/roundfix/SKILL.md`.
   - `TestAnUndocumentedCommandIsReported` removes every `roundfix reopen` from a copy of the skill text and expects `undocumentedCommands` to return exactly `reopen`.
3. MUST add text to `.agents/skills/roundfix/SKILL.md` for each command form the new test reports, written from the command's own help output and its section in `docs/user-guide/commands.md`. On the starting main these are `roundfix window set`, `roundfix window show`, `roundfix window clear`, `roundfix qa-report accept`, `roundfix baseline capabilities check`, `roundfix baseline profile init` and `roundfix init`. Each gets its synopsis, what it reads and writes, and its exit codes where the help states them. `roundfix qa-report accept` has no `--help`; its usage is printed by `roundfix qa-report --help`.
4. MUST list `agent-selection` with the other four public categories in the skill's Supervisor Run Event Stream text, as `roundfix events --help` does.
5. MUST replace the sentence that starts `Follow one order per Spec:` in the skill with this order: implement the graph including its authored gate, run the configured pre-PR review, archive on the branch, pass the repository gate, open the Pull Request, verify current-head checks, and merge. The sentence MUST contain the words `run the configured pre-PR review, archive on the branch`. It MUST say that with pre-PR review set to `none` the review is recorded as omitted.
6. MUST make the numbered steps under "Autonomous Spec delivery" follow that order:
   - a step titled `**Review the candidate.**` comes before `**Archive and publish.**` and names `roundfix review` and `roundfix review dispose`;
   - the step that runs `roundfix watch … --until-clean` is titled `**Resolve Pull Request feedback.**` and says it applies only when the repository's Review Source gives feedback on the Pull Request.
7. MUST state, in the step about corrective Tasks, that a Task added as a dependency of a QA gate that already settled `completed` needs the gate reopened with `roundfix reopen --spec <slug>` before the next Run, and that a failed or pending gate needs no reopen. The text MUST contain the words `reopen the settled gate first`.
8. MUST leave both version fields of the Roundfix skill at `0.0.2`. Three existing tests in `./skills` locate the owned-skill floor through that skill's version line and fail when it moves; the Spec that makes the floor follow the bundle raises it.
9. MUST regenerate the mirror with `make skills-sync`, then run `make baseline-digests`, and name in the Result every path either command rewrote.
10. MUST NOT change the `### QA settlement` section, `.agents/skills/roundfix/agents/openai.yaml`, any existing test, or any CLI source. The adapter version that file names belongs to another Spec.
11. MUST keep every phrase the existing contract tests require from the skill. `TestBaselineSkillContract` in `./skills` and the documentation contract tests in `./internal/docscontract` name them.

## Subtasks

- [ ] Add the command path helpers and their three tests.
- [ ] Document the seven missing command forms in the skill.
- [ ] Correct the event categories, the delivery order and its steps.
- [ ] State when a corrective Task needs the gate reopened.
- [ ] Regenerate the mirror and digests.

## Acceptance Criteria

- [ ] `commandPaths` returns the paths the rule defines for a literal help text, including alternatives and excluding flag-only lines.
- [ ] Every command path the root help lists is named as `roundfix <path>` in the canonical Roundfix skill.
- [ ] Removing one command's name from the skill text makes `undocumentedCommands` return exactly that path.
- [ ] A root help that yields fewer than forty paths fails the contract instead of passing it.
- [ ] The skill lists `agent-selection`, states the order with the pre-PR review before the archive, and says when `roundfix reopen` is needed.
- [ ] The `### QA settlement` section is byte-identical in `qa-gate`, `archive-spec` and the Roundfix skill.
- [ ] Both version fields still read `0.0.2`, and the mirror is byte-identical to the canonical skill.

## Context

- instruction: `docs/user-guide/commands.md`
- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `internal/docscontract/command_documentation_test.go`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestCommandPathsReadSubcommandsAndAlternatives|TestEveryCommandIsNamedInTheRoundfixSkill|TestAnUndocumentedCommandIsReported|TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract)$" ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCommandPathsReadSubcommandsAndAlternatives TestEveryCommandIsNamedInTheRoundfixSkill TestAnUndocumentedCommandIsReported TestBaselineDocumentationContract TestProfilesDocumentationContractMatchesPublicGuidance TestReleasePlanDocumentationContract; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for pair in ".agents/skills/roundfix/SKILL.md|agent-selection" ".agents/skills/roundfix/SKILL.md|run the configured pre-PR review, archive on the branch" ".agents/skills/roundfix/SKILL.md|**Review the candidate.**" ".agents/skills/roundfix/SKILL.md|**Resolve Pull Request feedback.**" ".agents/skills/roundfix/SKILL.md|reopen the settled gate first"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; for pair in ".agents/skills/roundfix/SKILL.md|archive, open the Pull Request, watch until Clean"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done; grep -q "^version: 0.0.2$" .agents/skills/roundfix/SKILL.md && grep -q "^  version: 0.0.2$" .agents/skills/roundfix/SKILL.md || { printf 'version fields are not 0.0.2 in %s\n' .agents/skills/roundfix/SKILL.md >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { printf 'mirror differs: %s\n' skills/roundfix >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestSettlementGuidanceIsOneTable|TestBaselineSkillContract|TestAuthorialSkillSync)$" ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettlementGuidanceIsOneTable TestBaselineSkillContract TestAuthorialSkillSync; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task none of the three new tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Interfaces; The command path rule; The Roundfix skill
- `_prd.md` → Goal 2; Goal 4; Core Feature 1; Core Feature 6; Core Feature 7; Success Metric 1; Success Metric 4; Success Metric 5
- `_techspec.md` → Testing Approach 1; Testing Approach 3; Testing Approach 4; Build Order 1

## Result

Implemented the command-name contract against real root help through `cli.Run`,
with first-seen deduplication, alternative expansion, and literal documentation
matching. The canonical Roundfix skill now describes the seven missing command
forms from the built CLI's help and the commands guide. It includes
`agent-selection`, orders pre-PR review before archive and publication, limits
Pull Request feedback watching to the applicable Review Source, and explains
reopening a completed QA gate before a corrective dependency's next Run.

### Focused evidence

All Go commands below used `GOCACHE=/private/tmp/roundfix-task01-go-cache`
because the sandbox denied the default macOS Go cache. These are implementation
checks, not the Daemon's declared Verification.

- Built the current CLI with `rtk proxy go build -o
  /private/tmp/roundfix-task01 ./cmd/roundfix`. Its help for `window set`,
  `window show`, `window clear`, `qa-report`, `baseline capabilities check`,
  `baseline profile init`, `init`, and `events` exited `0` and supplied the
  documented usage, effects, and stated exit codes. No state-changing command
  was invoked.
- Before the skill edit, `rtk proxy go test -tags docscontract
  ./internal/docscontract -run 'Test(CommandPaths|EveryCommand|AnUndocumented)'
  -count=1` exited `1`, reporting exactly `window set`, `window show`,
  `window clear`, `qa-report accept`, `baseline capabilities check`,
  `baseline profile init`, and `init` as missing. After the edit it exited `0`.
- After the final skill edit, `rtk proxy go test -tags docscontract
  ./internal/docscontract -run
  'Test(CommandPaths|EveryCommand|AnUndocumented|BaselineDocumentation|ProfilesDocumentation|ReleasePlanDocumentation)'
  -count=1` exited `0`, covering the three new tests and the existing related
  documentation contracts without editing them.
- A temporary Go overlay outside the repository replaced the help supplied to
  the coverage test's parser with `roundfix --help` alone. Running that test
  through the overlay exited `1` with `root help yielded 0 command paths, want
  at least 40`. The repository test source was unchanged by this probe.
- A read-only Python check against the built binary and repository files found
  51 unique advertised command paths and no missing literal names. It also
  checked lifecycle phrases and step ordering, both version fields, every file
  in the Roundfix skill mirror, and the settlement section against both HEAD
  and the two other canonical skills.
- `rtk proxy make verify-incremental` initially exited `2`. Two force-stop
  integration tests could not read the process table inside the sandbox. Its
  suiteguard also detected this Result being written while the suite ran.
  `TestBaselineExamplesParse` treated the new Baseline help synopses as
  executable Bash recipes; the synopses now use `text` fences because their
  optional-argument notation is reference syntax. A focused
  `rtk proxy go test ./internal/cli -run '^TestBaselineExamplesParse$' -count=1`
  exited `0`, checking that distinction. The existing example-parser dispatcher does not
  cover `baseline capabilities check`; extending it is a follow-up outside
  this Task's prohibition on existing-test edits. The new contract still
  checks the literal command and real root help, and the command's own help
  was exercised above.
- With the worktree held stable and process-table access permitted,
  `rtk proxy make verify-incremental` exited `0`: formatting, vet, the Go
  suite (including the force-stop tests), mirror checks, skill readiness,
  and build passed. The retry log is
  `/private/tmp/roundfix-task01-incremental-retry.log`. The focused
  documentation-contract command above also exited `0` again after the final
  synopsis-fence correction.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Command path rule | `TestCommandPathsReadSubcommandsAndAlternatives` passes with alternatives, positional stops, flag and parenthesis stops, repeated paths, and an ignored line outside the required prefix. |
| Every root-help path is named | `TestEveryCommandIsNamedInTheRoundfixSkill` passes against `cli.Run`; the independent binary check found all 51 names. |
| Removing a name reports exactly that path | `TestAnUndocumentedCommandIsReported` passes after removing every `roundfix reopen` literal from an in-memory skill copy. |
| Fewer than forty paths refuses | The coverage test checks the floor before documentation matching; the temporary flag-only-help overlay produced the expected floor failure. |
| Event category, review order, and reopen guidance | Read-only checks confirmed `agent-selection`, the required order sentence and step titles, review before archive, and `reopen the settled gate first`; prose also records `none` as omitted and no reopen for failed or pending gates. |
| Shared QA settlement remains identical | Read-only byte comparison confirmed the Roundfix section is unchanged from HEAD and identical to `qa-gate` and `archive-spec`. |
| Versions and mirror | Both version fields remain `0.0.2`; read-only comparison confirmed the entire Roundfix mirror matches the canonical directory. |

### Regeneration and scope

`rtk make skills-sync` exited `0` and changed only
`skills/roundfix/SKILL.md`. After synchronization,
`rtk proxy make baseline-digests` exited `0` with `changed: false`; no derived
path changed. Both commands were repeated after the final prose adjustment
with the same outcome.

Task-owned changed paths are `.agents/skills/roundfix/SKILL.md`,
`skills/roundfix/SKILL.md`,
`internal/docscontract/command_documentation_test.go`, and this Task's Result.
The Task file was already changed on arrival by its Daemon-owned status.
No other Task, manifest, existing test, CLI source, adapter file, or settlement
section was edited. Status and declared Verification remain Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20260930T144240Z_33c18e8a522f7217`
- Source commit: `c5f014945f862a47e1f21270a53f5b0617464c7e`
