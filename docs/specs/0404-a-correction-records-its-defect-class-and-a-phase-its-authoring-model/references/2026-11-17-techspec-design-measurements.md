# What the TechSpec of Spec 0404 measured and read (2026-11-17)

Recorded by the author of the TechSpec on the Mac mini on 2026-11-17, between about 15:55 and 17:10 BRT, in a clone of this repository on the branch `docs/plan-0404`, based on `a8236a87`. Every command was read-only. The Research Subagent read the code in four sweeps and the web in one; the author checked the statements this file keeps against the files and ran every command named here.

## 1. What a command run by the TechSpec's Author Subagent reads

The Author Subagent ran in Claude Code; its `AI_AGENT` variable held `claude-code_2-1-325_agent`. Only variable names were listed, and two values were printed.

- `CLAUDE_CODE_SESSION_ID` was set, and `roundfix phase start techspec` recorded its entry under it.
- `CLAUDE_EFFORT` held `xhigh`, the effort this repository declares for the `techspec` phase.
- No variable's value held `fable`, `opus`, `sonnet` or `haiku`.
- The subagent's own instructions state its model by name and by exact identifier.

This is the third session with the same result, after the two the Idea and the PRD recorded.

## 2. The identity of a repository by its root commits

- The clone was not shallow (`git rev-parse --is-shallow-repository` printed `false`), and its `origin` was an HTTPS address without a user name.
- `git rev-list --max-parents=0 HEAD` printed one commit and took 11 ms.
- `internal/baseline/repository.go` derives the identity in the unexported `inspectRepositoryIdentity`: the sorted root commits of `HEAD` and the object format, hashed with SHA-256 under `roundfix/repository-identity/v1`. The exported `InspectRepository` also reads the bounded repository snapshot, which this Spec does not need.

## 3. Where the plan of a delivered Spec is readable

- The commit `a8236a87` added the Archive Record of Spec 0400 to the default branch and holds no folder of that Spec.
- Its first parent holds `docs/specs/0400-planning-references-and-release-preparation-need-no-manual-classification/_tasks.md`, whose manifest names four Tasks and `qa: task_04`.
- An Archive Record holds `qa_task` and `source_revision` and no list of Tasks. `spec.EnsureSourceRevision` fetches `refs/pull/<n>/head` from the remote when the source revision is missing from the local objects.

## 4. The Governed Paths the design changes

A `go test -overlay` probe in `internal/speccheck`, which wrote nothing to the repository, ran `speccheck.GovernedPath` over 109 paths in two runs: every source, test, guide, skill, record and document file the TechSpec names or expects the Tasks to create, the seven new ADRs' directory, and controls. `git status --porcelain` printed nothing after each run.

- Governed and changed by the design: the fourteen paths the authorization record already bounded, and two more, `.agents/skills/roundfix/references/corrections.md`, which is new, and `.agents/skills/roundfix/references/storage.md`.
- Governed and not changed: `.roundfixrc.yml`, `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, `internal/docscontract/testdata/corpus-golden.json`, `internal/baseline/repository_test.go`, `docs/references/coverage-record.json`, `.agents/skills/roundfix/references/runs.md`, the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.
- Not governed: every Go source file the design changes or creates, the new test files, `internal/docscontract/command_documentation_test.go`, `internal/docscontract/user_guide_contract_test.go`, the test files under `skills/`, the mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json`, `CONTEXT.md`, the ADRs, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json` and the user guide.

## 5. The existing assertions the design moves

Read at `a8236a87`, with the changes Spec 0401 plans for the phase command taken from its Task files.

- `internal/cli/phase_test.go` compares the whole standard output, the whole standard error and the whole entry line of each call, and its first transcript the whole help text; `TestPhaseStartRecordsEachOfTheSevenPhases` requires exactly nine keys for every phase.
- `internal/phaserecord/record_test.go` holds the literal line of a `prd` entry and decodes entries against the nine-key list; `internal/phaserecord/read_test.go` compares the entries read back.
- `skills/phase_reference_test.go` requires the command line `roundfix phase start <phase> [--spec <slug>]` exactly once in the phase reference; `skills/phase_boundary_test.go` requires the phrase that runs the command in each authoring skill.
- `internal/cli/review_disposition_test.go` requires that the help of `roundfix review dispose` contain its two present usage lines, and `internal/cli/review_dismiss_after_fix_test.go` holds the bytes of a ledger line.
- `internal/delivery/park_class_test.go` requires a Park Class for every exported blocker constant; `internal/docscontract/command_documentation_test.go` holds the list of blockers of the `finding` row for the delivery guide and the delivery reference.
- `internal/cli/worktree_build_cache_test.go` and `internal/cli/phase_record_test.go` hold the Run Retention lines of the GC report in order.
- More than a dozen tests hold lines of `roundfix deliver status`, listed in the Research Subagent's sweep by file: the park, owner, turn, lane, run-slot, token-ceiling, limits, usage, lineage, numbered-files and revalidation tests of `internal/cli`.
- `internal/cli/cli_test.go` checks the root help and thirteen command helps by substring, and none of them is the help of `review dispose`, `phase`, `gc` or `corrections`; it holds the help of `deliver`.
- `internal/docscontract` requires every command of the root help in the Roundfix skill and in the top-level user guide, one index row for each command page, and one reference for each command family of the skill's reference index.

## 6. Sources read outside this repository

Secondbrain, read on 2026-11-17 (`wiki/index.md`, then two `qmd query` runs):

- The file `raw/files/books/designing-data-intensive-applications-early-release-2016/18-pages-447-491.md` holds chapter 11 of Designing Data-Intensive Applications, on the append-only ledger in which a mistake is answered by one more entry. It supports a record in which a later declaration is a new line and the earlier ones stay.
- The page `wiki/sources/digest-roundfix-raiz-git-valida-saida-humana-e-maquina-e-id-curto-2026-11-05.md` records that Git names an object by a leading substring that is unique in the repository. It supports accepting a commit by a prefix.
- The other hits were mirrors of this repository's ADRs.

External sources, read through Exa by the Research Subagent on 2026-11-17 and quoted as it returned them:

- Git, `git-rev-list` documentation, https://git-scm.com/docs/git-rev-list — on `--max-parents=0`, that it "gives all root commits".
- Git, `git-clone` documentation, https://git-scm.com/docs/git-clone — on `--depth`, "Create a shallow clone with a history truncated to the specified number of commits." It is why a shallow clone has another identity.
- Git, `git-merge` documentation, https://git-scm.com/docs/git-merge — `--allow-unrelated-histories` merges "histories of two projects that started their lives independently", so a repository may have several roots.
- Two tools that name a repository by its root commit and one that prefers the remote address, read in their public source: one keeps the root commit because it "survives remote renames, host migrations, and SSH/HTTPS URL changes"; the other because "root commits are shared across unrelated forks that clone the same template". This is the counter-evidence ADR-0421 answers with the Spec slug.
- Go, `net/url` documentation, https://pkg.go.dev/net/url — a password in a URL "has proven to be a security risk in almost every case where it has been used".
- The Open Group, `write()`, https://pubs.opengroup.org/onlinepubs/9799919799/functions/write.html — with `O_APPEND`, "the file offset shall be set to the end of the file prior to each write". The Linux manual of `open(2)` adds that appends from several processes may corrupt a file on NFS.
- JSON Lines, https://jsonlines.org/ — one valid JSON value per line, UTF-8, ended by a line feed.
- Fowler, "Event Sourcing", https://martinfowler.com/eaaDev/EventSourcing.html — "all changes to application state are stored as a sequence of events".
- GitHub Docs, "About pull request merges", https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges — on a squash merge, "the pull request's commits are squashed into a single commit".
- Git, `git-notes` documentation, https://git-scm.com/docs/git-notes — notes are "attached to objects, without touching the objects themselves".
- Claude Code, environment variables, https://code.claude.com/docs/en/env-vars — `CLAUDE_CODE_SESSION_ID` and `CLAUDE_EFFORT` are set in the subprocesses of the Bash tool; `ANTHROPIC_MODEL` is the "Name of the model setting to use"; the page lists no variable that names the running model. The page on subagents says a subagent's model comes from its definition, an invocation parameter, a variable or the conversation, and does not say that a subagent is told its model.
- On a required cause field before a merge the sweep found no study. The nearest, an Atlassian support page on a required root-cause field, applies it to two priorities and offers "NA" for the rest; a Lucent study reports that "different analysts often came up with different RCA results". Both support keeping `unclear` and the later declaration, as the PRD already decided.

A digest of the sources this phase added went to the project that owns the Secondbrain, after a duplicate check; the report of the phase names the message.

## Limits

- Section 1 is one session of one harness version.
- Section 4 covers the paths a TechSpec can name. The Tasks stage measures the paths the Tasks declare and every file a regeneration command rewrites.
- Section 5 is a reading, not a run: the Task that changes a package runs that package before and after, and declares every test whose expectation it changes.
