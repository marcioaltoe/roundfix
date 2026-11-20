---
task: task_05
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec adds the Flow Command List with its Named Exclusions and Guard Entries in the package `internal/flowpermissions`, the command `roundfix permissions` that previews it for Claude Code, Codex and OpenCode and writes or removes it only with a preview's digest, the Doctor Command's check `flow-permissions`, the Update Check's part `flowPermissions`, one pattern of the governed set for the project permission files of the three harnesses, a documentation contract that keeps the list equal to what the owned skills name, the Roundfix skill's rule for applying the list, the command's guide and the Doctor, update and configuration guide parts, four glossary terms and the revision of Governed Path, with ADR-0444, ADR-0445 and ADR-0446. No existing command, flag, exit code or Project Config key changes.

Every row runs on this repository's own files, temporary homes, temporary Git repositories and temporary copies of the audited head. No row:

- starts an Agent harness, the real `claude`, `codex` or `opencode`, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a Run, a Daemon or a Delivery Owner;
- runs `roundfix doctor` or `roundfix update check` through the built binary, because the Doctor's other checks prove adapters and reach the forge and the Update Check reaches the release listing; the rows of the two reports run their entry points through the named tests;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- creates or changes a Harness Permission File of this repository: `.claude/settings.json`, `.claude/settings.local.json`, anything under `.codex/`, `opencode.json` or `opencode.jsonc`;
- sends a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`; their tests make their own temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`, call an installed formatter that needs the host's toolchains. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary Git repository for a product row is made outside the checkout with `git init` and one commit; a temporary copy of the audited head is made with `git archive` of that commit and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 5 by running, with `-v` over `./internal/flowpermissions`, `TestFlowCommandList`, `TestRenderingsHoldNoForbiddenRule`, `TestGuardEntriesMatchTheirOption`, `TestCodexRulesFile` and `TestRetiredSetKeepsEveryReleasedEntry`, and, with `-tags docscontract -v` over `./internal/docscontract`, `TestOwnedSkillsNameOnlyListedOrExcludedCommands` and `TestAnUnlistedSkillCommandIsReported`, recording every subtest by name; then the whole `./internal/flowpermissions` package, recording its pass.
3. MUST verify Success Metrics 3 and 4 by running, with `-v` over `./internal/flowpermissions`, `TestEditClaudeSettings`, `TestReadOpenCode`, `TestPlan` and `TestApply`, recording every subtest by name, and recording for each link subtest whether it passed or skipped and why.
4. MUST reproduce each Surface Transcript of `_techspec.md` → Surface Transcripts through `bin/roundfix`, in a fresh home prepared as above, in a temporary Git repository with one commit and no harness file, and record the whole output of each with its exit code: Surface Transcript 1 (`bin/roundfix permissions --help`, exit 0), Surface Transcript 2 (`bin/roundfix permissions --harness claude`, exit 0), Surface Transcript 3 (`bin/roundfix permissions --harness claude --apply 000000000000`, exit 3, nothing on standard output) and Surface Transcript 4 (`bin/roundfix permissions --harness opencode`, exit 0). Record that every line the transcripts fix appears in its order and that only the version and the digest vary. Run also, with `-v` over `./internal/cli`, `TestPermissionsSurfaceTranscripts`, `TestPermissionsCommand` and `TestRootHelpNamesThePermissionsCommand`, recording every subtest by name.
5. MUST reproduce the journey of Success Metrics 3 and 4 through `bin/roundfix`, in a fresh home prepared as above, in a temporary Git repository that holds `.claude/settings.json` with this content, kept byte for byte as `before.json`: a `permissions` object whose `allow` list holds `Bash(make test)` and whose `ask` list holds `Bash(rm *)`, a top-level `model` key, tab indentation and no final newline:
   - `bin/roundfix permissions --harness claude --format json` prints `state` `not_applied`, a non-empty `digest` and a `next` that names it; the counts of `add` and `addGuards` are recorded as the preview reports them, never as a number fixed in this plan; the file is unchanged;
   - `bin/roundfix permissions --harness claude --apply <that digest>` exits 0 and prints `Wrote <n> entries and <m> Guard Entries to .claude/settings.json.` with the recorded counts; the file then holds `Bash(make test)`, `Bash(rm *)` and the `model` key at their places, every entry of the preview under `permissions.allow` and every Guard Entry under `permissions.ask`, keeps its tabs and its lack of a final newline, and decodes as strict JSON;
   - a second preview reports `current` and `Nothing to write: the file holds the installed list.`, and exits 0;
   - `bin/roundfix permissions --harness claude --remove` names each entry it would remove and prints a digest; the removal with that digest exits 0 and leaves the file byte-identical to `before.json`;
   - after a new preview, one byte of the file is changed and the apply with the old digest exits 3, prints the two lines of `_techspec.md` → API Contract 3 on standard error and leaves the file as it was after the change;
   - `--apply` with no `--harness`, `--harness opencode --apply x`, `--harness vim` and `--format yaml` each exit 2; the command run in a directory that is not a Git repository exits 2.
6. MUST reproduce the Codex journey through `bin/roundfix`, in a fresh home prepared as above, in a temporary Git repository holding `.codex/rules/default.rules` with one rule of its own, kept byte for byte: the preview names `.codex/rules/roundfix-flow.rules`, prints `No option of a listed read can be kept asked about in Codex.` and the three Codex lines of `What this does not do:` of `_techspec.md` → API Contract 2; the apply with its digest writes that file with the header line `# To remove it: roundfix permissions --harness codex --remove`, one `prefix_rule(` block per group with `decision = "allow",`, no quoted alternative `"permissions"`, `"baseline"`, `"profiles"`, `"release"`, `"skills"`, `"push"`, `"merge"` or `"auth"`, and leaves `default.rules` byte-identical; the removal with its preview's digest deletes `roundfix-flow.rules` only and keeps `.codex/rules/` and `default.rules`.
7. MUST measure offline what the PRD and the TechSpec left unmeasured, and record each measurement with what stays unmeasured:
   - the short-option form of a Guard Entry: record from the JSON preview of requirement 5 that `addGuards` holds `Bash(gh api* -X*)`, `Bash(gh api* -F*)`, `Bash(gh api* -f*)`, `Bash(git grep* -O*)` and `Bash(git archive* -o*)` for every guarded entry still on the list, and that `TestGuardEntriesMatchTheirOption/short_option_form`, `after_another_argument` and `right_after_the_command` pass, which prove the form under the documented reading of `*` as any text; record that how Claude Code itself matches the short form in a real session stays unmeasured and is the first entry of `_prd.md` → Unreachable Acceptance;
   - the file name of the Codex rules: record that the written file is `.codex/rules/roundfix-flow.rules`, inside the folder and with the extension the Codex documentation names, as `references/2026-11-19-harness-rule-measurements.md` quotes it, that its `prefix_rule` blocks have the form whose decision that record measured with Codex's offline checker, and that `TestCodexRulesFile/reads_back_every_entry` passes; record that whether a Codex session loads the file under that name stays unmeasured and is the second entry of `_prd.md` → Unreachable Acceptance. The gate does not run `codex`, `claude` or `opencode`, which the Spec's authorization keeps out of every test and Verification.
8. MUST reproduce the link refusal and the unreadable file through `bin/roundfix`, in a fresh home prepared as above: in a temporary Git repository whose `.claude` is a link to a directory outside it, the apply with the preview's digest exits 1, prints `.claude/settings.json leaves the repository through a link. Nothing was written.` on standard error and writes nothing outside the repository; in a temporary Git repository whose `.claude/settings.json` holds a `//` comment, the preview prints `State: not readable: ` with a reason and no digest, and `bin/roundfix permissions --harness claude --apply 000000000000` exits 1, prints `.claude/settings.json cannot be read as Claude Code documents it: <reason>. Nothing was written.` on standard error and leaves the file byte-identical. Record the platform when the link cannot be created.
9. MUST verify Success Metric 6 by running, with `-v` over `./internal/cli`, `TestDoctorReportsFlowPermissions` and `TestUpdateCheckReportsFlowPermissions`, recording every subtest by name, then `-run 'Doctor|UpdateCheck|Permissions' ./internal/cli`, recording its pass. Record that the Doctor rows run the command's entry point, that `never_fails` and `exit_unchanged` hold, that `after_the_baseline_checks_before_residue` holds, and that each warning names `roundfix permissions --harness <name>`.
10. MUST verify Success Metric 7 by running `TestOnlyTheCommandAndTheReportsUseTheList` with `-v` over `./internal/flowpermissions`, recording both subtests, and `TestPermissionsCommand/home_files_untouched` and the `writes_nothing` subtests of the two report tests.
11. MUST verify Success Metric 8 by running, with `-v` over `./internal/speccheck`, `TestGovernedPath`, and with `-tags repocontract -v`, `TestGovernedSetOnlyGrows`, recording every new row and subtest by name, then the whole `./internal/speccheck` package with and without `-tags repocontract`, recording each pass.
12. MUST verify Success Metrics 9 and 10 by running, with `-v`, `TestTheRoundfixSkillDescribesFlowPermissions` and `TestTheRoundfixSkillNamesTheFlowPermissionsReports` over `./skills`, recording every subtest by name; then `TestEveryOwnedSkillVersionIsRecorded`, the whole `./skills` package, `make skills-sync-check`, the whole `-tags docscontract ./internal/docscontract` package and, with `-v`, `TestEveryCommandIsNamedInTheRoundfixSkill`, `TestTheSkillCoverageMapIsCurrent` and `TestGlossaryPhrasesAreDefinedTerms`, and `-run '^TestCoverageEquivalence$' ./internal/spec`, recording each pass. Record that `skills/testdata/owned-skill-versions.json` holds the version each copy of the Roundfix `SKILL.md` declares.
13. MUST run each sabotage below on its own temporary copy of the audited head. Each sabotage keeps the package compiling and regenerates or synchronizes nothing, so the named test fails on its assertion; record that it fails on the named subtest or test.
    1. The entry `spec` is removed from the `roundfix` group of the list: `TestOwnedSkillsNameOnlyListedOrExcludedCommands` fails naming `roundfix spec`.
    2. An entry whose words are `python3` and `-c` is added to the list data: `TestRenderingsHoldNoForbiddenRule/claude`.
    3. The short-option Guard Entry of `gh api` for `-X` renders as `Bash(gh api -X*)` instead of `Bash(gh api* -X*)`: `TestGuardEntriesMatchTheirOption/after_another_argument`.
    4. The settings editor writes its result through `json.MarshalIndent` of the decoded file instead of spans: `TestEditClaudeSettings/tab_indentation_kept`.
    5. `Apply` skips the comparison of the given digest with the plan's: `TestApply/stale_digest_writes_nothing`.
    6. The governed pattern loses its `opencode` alternative: `TestGovernedPath/OpenCode_configuration`.
    7. The Doctor check returns the status `failed` for a behind harness: `TestDoctorReportsFlowPermissions/never_fails`.
    8. The sentence "A recommendation, silence or a missing answer is never that yes." is removed from the mirror `skills/roundfix/references/setup.md`: `TestTheRoundfixSkillDescribesFlowPermissions/silence_is_never_that_yes`.
    9. `printPermissionsExclusions` skips again every catalog entry whose second word begins with a hyphen: `TestPermissionsTextPreviewNamesEveryExclusion/claude` and `TestPermissionsTextPreviewNamesEveryExclusion/codex`, each naming `git -C`.
14. MUST execute the tests task_01 to task_04 and the corrective task_06 name in their Verification.
15. MUST search and record:
    - every rendering of the list printed by the JSON previews of requirements 5 and 6 holds no entry whose first word after `Bash(` is `bash`, `sh`, `zsh`, `python`, `python3`, `node`, `go`, `make`, `npm`, `bun` or `npx`, no `*` in the place of a command or a subcommand, and no rule for a Named Exclusion;
    - every preview names the Named Exclusions with their reasons and prints "It is not a security boundary: a roundfix command runs what this repository declares."; the text preview of Claude Code and of Codex names, entry by entry, the same Named Exclusions with the same reasons as the JSON preview of the same plan, `git -C` included, and its Git line reads `  git -C, config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration`, the line corrective task_06 fixes in place of the example of `_techspec.md` → API Contract 2, which predates the classification of `git -C`;
    - both copies of the Roundfix `SKILL.md`, `setup` and `modes` references carry the texts of `_techspec.md` → API Contract 7 and of task_04, and `.agents/skills` and `skills` hold the same bytes for each changed file; no `### QA settlement` section changed;
    - `docs/user-guide/commands/permissions.md`, `docs/user-guide/commands.md`, `docs/user-guide/usage.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/update.md` and `docs/user-guide/configuration.md` carry their texts, and the Doctor and Update Check help carry their sentences;
    - ADR-0444, ADR-0445 and ADR-0446 are `accepted` and carry the Knowledge Mark;
    - no existing command, flag, exit code or Project Config key of Roundfix was added, changed or removed other than the new command `roundfix permissions`.
16. MUST record, as outside-evidence rows, the evidence this Spec did not author, each with its date and none blocked, through what `_prd.md` → Acceptance evidence and `references/2026-11-19-harness-rule-measurements.md` record, looking nothing up:
    - the Claude Code permissions documentation (read 2026-11-19): compare "Rules are evaluated in order: deny, then ask, then allow." with the Guard Entries written under `permissions.ask` in requirement 5;
    - the Claude Code sub-agents documentation (read 2026-11-19): compare "Each inherits the parent conversation's permission rules" with the list written into the project's shared settings file, which a subagent of the session reads;
    - the Codex rules documentation (read 2026-11-19): compare "Run the command outside the sandbox without prompting." with the Codex preview line "In Codex an allowed command runs outside the sandbox with nobody asked.";
    - Pillar Security, "GitPwned: Allowlist to RCE", 2026-07-20 (read 2026-11-19): compare "Command-name allowlists don't work for commands with complex option semantics." with the Guard Entries of `git show`, `git diff`, `git grep`, `git archive` and `gh api`.
17. MUST record the three entries of `_prd.md` → Unreachable Acceptance, one row each, as `blocked` with `declared: ` followed by the entry's criterion text: the walk in an interactive Claude Code session of a trusted folder; the walk in an interactive Codex session that loads the rules file of Roundfix's own; and Success Metric 12 over the 30 days after a repository applies the list.
18. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD declares `adds:` for **Flow Command List**, **Named Exclusion**, **Guard Entry** and **Harness Permission File** and `changes:` for **Governed Path**, and the TechSpec declares none; that `CONTEXT.md` defines the four terms right after the entry of **Undeclared Governed Path** with the texts of task_03 and task_04, that the entry of **Governed Path** gains only its one sentence, and that no other entry changed; and that `TestGlossaryPhrasesAreDefinedTerms` passes.
19. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `internal/cli/cli_test.go`, `docs/references/coverage-record.json`, any file under `docs/agents/` or `internal/baseline/assets/`, any file under `.claude/` or `.codex/`, `opencode.json`, `opencode.jsonc`, any Vendored Skill under `.agents/skills`, and any skill other than `roundfix`, of which only `SKILL.md`, the `setup` and `modes` references and their mirrors;
    - no accepted ADR other than ADR-0444, ADR-0445 and ADR-0446 changed;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
20. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
21. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, `claude`, `codex` or `opencode`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - create or change a Harness Permission File of this repository;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the four Tasks and their packages.
- [ ] Reproduce the four Surface Transcripts and the Claude Code and Codex journeys through the binary in temporary repositories.
- [ ] Measure offline the short-option Guard Entries and the Codex file, and record what stays unmeasured.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside evidence, the three declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A preview writes nothing, an apply writes exactly what the preview with that digest showed, a stale digest writes nothing, and a removal returns a settings file to its first bytes.
- [ ] No rendering holds a rule for a shell, an interpreter, a wildcard command or a Named Exclusion, the Codex file allows every entry, and each Guard Entry is under the ask list.
- [ ] The Doctor Command and the Update Check report every state without failing or writing, the harness permission files are Governed Paths, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-11; Core Features 1-16; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Success Metric 10; Success Metric 11; Success Metric 12; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; Invariants 1-10; API Contracts 1-7; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Testing Approach; Build Order 4
- [`references/2026-11-19-harness-rule-measurements.md`](references/2026-11-19-harness-rule-measurements.md)
- [`references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md`](references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md)
- task_06 → Requirements 1-5
- ADR-0444, ADR-0445, ADR-0446
