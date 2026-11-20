---
task: task_02
spec: 0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set
status: completed
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec makes the `typescript` Baseline module require `solid`, adds the catalog diagnostic `catalog.skill.activation.skill.unrequired` that refuses a Skill Activation whose bundle names a skill its owner does not require, moves restore, the update preview, `skills check` and Doctor onto one `Catalog.RequiredSkills` method, regenerates the module version, the digests, the goldens and this repository's Setup Manifest, and adds one sentence to the glossary entry of **Skill Activation**, with ADR-0448. No command, flag, output line, field, exit code or Project Config key is added or changed.

Every row runs on this repository's own files, temporary homes, temporary Git repositories and temporary copies of the audited head. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a Run, a Daemon or a Delivery Owner;
- runs an approved `roundfix baseline update` or a restore that reaches the skills' upstream repository; every restore row stops at the selection with an absent `--source-dir`;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- sends a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary Git repository for a product row is made outside the checkout with `git init` and one commit; a temporary copy of the audited head is made with `git archive` of that commit and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 2 by running, with `-v` over `./internal/baseline`, `TestAnActivationOwnerRequiresEverySkillOfItsBundle`, recording every subtest by name, and recording that `a_sibling_module_does_not_satisfy_the_owner` expects exactly one `catalog.skill.activation.skill.unrequired` diagnostic with path `trigger.production-code` and info `solid` while the `frontend` module still requires `solid`.
3. MUST verify Success Metrics 3 and 5 by running, with `-v` over `./internal/baseline`, `TestTheRequiredSkillsCoverTheRenderedDispatch`, recording every subtest by name and the five counts `the_built_in_profiles_keep_their_names` asserts: 41, 105, 38, 34 and 95 for `go-cli-tui`, `go-cli-typescript-monorepo`, `python`, `rust-cli` and `standard-typescript-monorepo`, with `solid` only in the two TypeScript profiles; then `TestSkillSnapshotProfileKeepsEveryBuiltInProfile`, `TestSkillSnapshotProfileResolvesARepositoryProfile` and `TestEveryModuleRequiresEverySkillItDispatches`, recording each pass, and the whole `./internal/baseline` package, recording its pass.
4. MUST verify Success Metric 1 by running, with `-v` over `./internal/baseline`, `TestARepositoryProfileWithoutFrontendRequiresSolid`, recording its pass and that it asserts `solid` among the required skills with the contract the `typescript` Setup Snapshot pins and a restore selection that reports `restore.source-dir-invalid`.
5. MUST verify Success Metric 4 through `bin/roundfix`, in a fresh home prepared as above, in the temporary repository of `references/2026-11-20-request-and-characterization.md`, section 3, made outside the checkout: `git init` and one commit; `bin/roundfix baseline profile init --id ts-backend --from standard-typescript-monorepo`; the profile's `modules` cut to `core`, `context-workflow`, `typescript`, `bun`, `autonomous-work`, `spec-workflow`, `secondbrain` and `repository-extension`; `docs/agents/setup-context.json` holding `{"profile": "ts-backend", "modules": [those eight]}`; `bin/roundfix skills install --target project`; and a `skills-lock.json` holding `{"version": 1, "skills": {}}`. Record the whole output and exit code of each of these:
   - `bin/roundfix baseline profile validate ts-backend` reports the profile valid as a repository profile and exits 0;
   - `bin/roundfix skills check` exits 1 and prints, among the lines "required by the Baseline and not installed", the line `.agents/skills/solid: required by the Baseline and not installed; restore with roundfix baseline skills restore --profile ts-backend --skill solid`; record the number of such lines as the command prints it, never as a number fixed in this plan;
   - `bin/roundfix baseline skills restore --profile ts-backend --skill solid --source-dir <a path that does not exist>` exits 2 and reports `restore.source-dir-invalid`, never `restore.skill-invalid`.
   The update preview reads the complete Setup Manifest before it inspects skills, so the two-field `docs/agents/setup-context.json` above cannot reach it. Run it in a second temporary repository made outside the checkout, in the same fresh home, whose Setup Manifest the binary writes by adopting the Baseline:
   - `git init`; `bin/roundfix baseline profile init --id ts-backend --from standard-typescript-monorepo`; the profile's `modules` cut to the same eight, its `capabilities` set to `[]`, and `auth.provider` and `auth.base_path` removed from its `decisions`, since a decision whose capability is not selected fails validation; a `Makefile` with the targets `verify` and `verify-incremental`, each running `@true`; `.agents/skills/context7-cli/SKILL.md` and `.agents/skills/exa-web-search/SKILL.md`, each a fixture file with `name` and `description` frontmatter, which satisfy the two universal installed-skill capabilities without a restore; then one commit;
   - `bin/roundfix baseline profile validate ts-backend` reports the profile valid as a repository profile;
   - `bin/roundfix baseline plan --profile ts-backend --repo . --format json` with these answers, each passed as `--decision`: `preservation.mode=greenfield`, `language.generated=English`, `verification.gate=make verify`, `verification.incremental=make verify-incremental`, `spec.scaffold=true`, `domain.layout=single-context`, `triage.external=false`, `autonomous.enabled=true`, `runtime.backend=codex gpt-5.5 xhigh`, `runtime.design=claude opus xhigh`, `secondbrain.enabled=true`, `repository.extension.enabled=true`, `identifier.strategy={"kind":"uuid-v7"}` and `http.contract={"mode":"REST"}`; it exits 0 and prints a `planDigest`; the plan is kept outside the repository;
   - `bin/roundfix baseline apply --plan <that plan> --confirm-plan <its planDigest> --repo . --format text` reports `Baseline apply: verified` and writes `docs/agents/setup-context.json` for `ts-backend`;
   - `bin/roundfix skills install --target project` and a `skills-lock.json` holding `{"version": 1, "skills": {}}`; then one commit.

   Then run `bin/roundfix baseline update --repo <that second repository> --no-history --format text`, without `--yes`, so that it writes nothing and reaches no network, and record its exit code, that its preview names `solid` as a Missing Required Skill with the restore command, `- missing solid: required by the Baseline and not installed; restore with roundfix baseline skills restore --profile ts-backend --skill solid`, and that `git status --short` in that repository is empty afterwards. When that preview needs the network for another reason, record the row as blocked by the environment with that reason and keep the three rows above as its evidence.
6. MUST verify Goal 3 and invariant 4 by running, over `./internal/cli`, `-run 'Doctor|SkillsCheck|MissingRequiredSkill|ExternalSkill|SkillsRestore|BaselineUpdate'`, recording its pass, and, with `-v`, `TestDoctorNamesAMissingRequiredSkill`, `TestBaselineUpdatePreviewNamesAMissingRequiredSkill`, `TestSkillsCheckNamesAMissingRequiredSkill` and `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent`, recording each pass; then the whole `./internal/cli` package, recording its pass. Record by reading the source that `loadRestoreProfile`, `loadSkillSnapshotProfile` and `resolveExternalSkillRequirement` take the required skills only from `Catalog.RequiredSkills`, and that the Doctor and restore messages for an unknown module are the ones `_techspec.md` → API Contract 2 quotes.
7. MUST verify that the catalog changed only as `_techspec.md` → API Contract 3 states: `internal/baseline/assets/modules/typescript.json` differs from the Delivery Base only in `"solid"` after `"onioncry"` in `requiredSkills` and its `version` line; `internal/baseline/module-versions.json` gains only the new version of `typescript` with its digest; and running, on a temporary copy of the audited head, the module record command, `make baseline-digests` and `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` changes no file, which proves that the version line, the digests, the goldens and `docs/agents/setup-context.json` are the output of the repository's commands.
8. MUST run each sabotage below on its own temporary copy of the audited head. Each sabotage keeps the package compiling and regenerates nothing, so the named test fails on its assertion; record that it fails on the named subtest or test.
   1. `"solid"` is removed from the `requiredSkills` of the `typescript` module: `LoadEmbeddedCatalog` fails with `catalog.skill.activation.skill.unrequired` for `trigger.production-code` and `solid`, and `TestAnActivationOwnerRequiresEverySkillOfItsBundle/the_embedded_catalog_loads` fails.
   2. The new rule in `validateSkillActivations` is deleted: `TestAnActivationOwnerRequiresEverySkillOfItsBundle/a_sibling_module_does_not_satisfy_the_owner` fails.
   3. The rule also accepts a skill that any module of the catalog requires, not only the owner and its dependencies: `TestAnActivationOwnerRequiresEverySkillOfItsBundle/a_sibling_module_does_not_satisfy_the_owner` fails.
   4. `Catalog.RequiredSkills` skips the last module it is given: `TestTheRequiredSkillsCoverTheRenderedDispatch` fails.
9. MUST execute the tests task_01 names in its Verification.
10. MUST record the outside-evidence row of Success Metric 6, with its date and not blocked, through what `_prd.md` → Acceptance evidence records, looking nothing up: the Bazel documentation "Dependencies", read 2026-11-20, "For correct builds, the graph of actual dependencies A must be a subgraph of the graph of declared dependencies D." Compare it with the diagnostic of requirement 2 and sabotage 1: the skills an activation uses (its bundle) must be a subset of the skills its owner declares (the `requiredSkills` of the owner and its dependencies), and a skill that only a sibling module declares, as `frontend` declared `solid`, does not count.
11. MUST record that `_prd.md` → Unreachable Acceptance declares none, so no row is blocked by declaration.
12. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD and the TechSpec declare `changes:` for **Skill Activation** and no `adds:`; that `CONTEXT.md` gives the entry **Skill Activation**, after its first sentence and before its `_Avoid_` line, the sentence of `_techspec.md` → API Contract 4 word for word with ADR-0448 linked; that no other entry changed; and that `TestGlossaryPhrasesAreDefinedTerms` passes with `-tags docscontract` over `./internal/docscontract`.
13. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`, and the only Governed Paths changed are `internal/baseline/assets/modules/typescript.json` and `docs/agents/setup-context.json`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/behavior-surfaces.json`, `docs/references/coverage-record.json`, any entry of `docs/references/skill-coverage.json` other than the `review` of `command doctor`, which holds the Coverage Review of task_01 requirement 11 word for word, any file under `docs/agents/` other than `docs/agents/setup-context.json`, any file under `internal/baseline/assets/` other than `internal/baseline/assets/modules/typescript.json`, any user-guide page, and any owned or Vendored Skill;
    - no accepted ADR other than ADR-0448 changed, and ADR-0448 is `accepted`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
14. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
15. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of task_01 and their packages.
- [ ] Reproduce the repository profile without `frontend` through the binary in a temporary repository.
- [ ] Prove the regenerated records are the commands' output on a temporary copy.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The embedded catalog loads, every activation's bundle is required by its owner or the owner's dependencies, and a catalog that breaks the rule does not load and names the activation and the skill.
- [ ] A repository profile with `typescript` and without `frontend` reports `solid` as a Missing Required Skill with the restore command, and the restore command accepts the selection.
- [ ] The five built-in profiles keep their required skills, every sabotage makes its named test fail, and the scope audit finds no change outside the declared files and the two granted Governed Paths.

## Verification

- `newest="$(find 'docs/specs/0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; Invariants 1-5; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Testing Approach; Risks & Considerations; Build Order 2
- [`references/2026-11-20-request-and-characterization.md`](references/2026-11-20-request-and-characterization.md), sections 3, 5 and 7
- ADR-0448
