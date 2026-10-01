---
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
prd: _prd.md
created: 2026-09-30
---

# A skill and a command guide read one command at a time — Technical Spec

## Executive Summary

Two documents are split by command: the Roundfix Skill
(`.agents/skills/roundfix/SKILL.md`) and the command reference
(`docs/user-guide/commands.md`). Each becomes an entry file plus companion
files, with no generator. The split moves whole sections and rewords nothing,
and a line-count proof run by the Daemon shows that nothing was lost.

Every contract that pinned a phrase in either document reads the entry file
and its companions as one text, through one small reader. That change lands
first, while the documents are still single files, so the split itself touches
no test logic.

The trade-off this design accepts is that a reference file starts with a
second-level heading and the skill can no longer be pasted as one file. The
alternative, a generated single file, would keep one shared file that every
CLI Task edits, which is the collision this Spec removes.

## Project Constraints

- Identifier strategy: not applicable — no identifier is added. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0187 (this Spec), ADR-0130,
  ADR-0081, ADR-0149, ADR-0166, ADR-0167 and ADR-0176 hold as the PRD states. The gate
  is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117,
  ADR-0155 and ADR-0156, and by ADR-0093 and ADR-0094. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer's authorization of
  2026-09-30 for every skill and the standing grant of 2026-09-21 for the four
  governed test files, recorded in [_authorization.md](_authorization.md);
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/release.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/review-runs.md`,
  `.agents/skills/roundfix/references/runs.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/references/spec-delivery.md`,
  `.agents/skills/roundfix/references/stop.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `internal/docscontract/publicdocs_test.go`, `internal/cli/cli_test.go`,
  `internal/cli/baseline_documentation_contract_test.go`,
  `skills/baseline_skill_contract_test.go`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

One new package, `internal/mdtree`, holds the reader. It is new because three
packages need it (`skills`, `internal/cli` tests and `internal/docscontract`
tests) and none of them may import another's test helpers.

Nothing else is new. The bundle already handles nested files:

- `skills/skills.go` embeds each owned skill directory recursively, and
  `Files`, `Install` and `SkillFolderHash` walk the whole tree. Other owned
  skills already ship a `references/` directory.
- `make skills-sync` copies each skill directory with `cp -R`, and
  `make skills-sync-check` compares with `diff -r`. No `Makefile` change is
  needed.
- Owned skills carry no content pin, so `make baseline-digests` has nothing to
  rewrite for a skill's content. It still runs, as the sanctioned step after a
  skill edit.
- `SC-CLI-UNDOCUMENTED` accepts any declared path under `.agents/skills/`,
  `skills/`, `docs/user-guide/` or `docs/agents/`, so a command file already
  satisfies it.
- The test-set selector and the release classifier decide by directory, so
  the new files fall in the same sets as the old ones.

## Implementation Design

### Interfaces

```go
// Package mdtree reads a Markdown entry file with the companion files it
// routes to, as one text.
package mdtree

// Text returns the bytes of entry followed by every .md file directly inside
// companion, in lexical order, each preceded by one newline. A missing
// companion directory yields entry alone. A missing entry is an error.
func Text(fsys fs.FS, entry, companion string) (string, error)
```

Callers pass `os.DirFS(repoRoot)` or the embedded bundle. The two pairs are
`roundfix/SKILL.md` with `roundfix/references`, and
`docs/user-guide/commands.md` with `docs/user-guide/commands`.

### Data Models

None. No schema, record or payload changes.

### Pins read the tree

`CheckReadiness` in `skills/skills.go` checks required wording and banned
branding in `roundfix/SKILL.md`. It moves that check into an unexported
function over an `fs.FS`, which reads `mdtree.Text` for the skill. Diagnostics
keep their text and their `roundfix/SKILL.md` path.

Every Go reader of either document's content switches to `mdtree.Text`. On
`9e439dbb` they are:

| Test | File |
| --- | --- |
| `TestBaselineDocumentationContract` (command reference and canonical skill rows) | `internal/docscontract/publicdocs_test.go` |
| `TestProfilesDocumentationContractMatchesPublicGuidance` | `internal/docscontract/publicdocs_test.go` |
| `TestReleasePlanDocumentationContract` (canonical and mirror) | `internal/docscontract/publicdocs_test.go` |
| `TestEventsHelpDocumentsAgentSelectionFilter` | `internal/cli/cli_test.go` |
| `TestBaselineExamplesParse` | `internal/cli/baseline_documentation_contract_test.go` |
| `TestNoPythonBaselineRuntime` | `skills/baseline_skill_contract_test.go` |

Earlier Specs of the wave may add readers. The sweep that finds them is
`grep -rn --include='*.go' -e 'commands\.md' -e '"roundfix", "SKILL.md"' -e 'roundfix/SKILL\.md' .`.
Three readers stay on the entry file on purpose:

- `TestSettlementGuidanceIsOneTable` reads the section the entry file keeps.
- `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical` appends to the entry
  file.
- Tests that only check that an installed `roundfix/SKILL.md` exists.

### The skill layout

The body of `SKILL.md`, after its front matter, is cut at every line that
starts with `## ` outside a fenced code block. Each section moves whole, with
its heading, to the file the table names. A file that receives several
sections keeps them in their original order.

| Section heading | File under `references/` |
| --- | --- |
| acpx dependency | `runtime.md` |
| Setup, doctor, and upgrade; Config compatibility | `setup.md` |
| Pre-PR review | `review.md` |
| Delivery queue | `deliver.md` |
| Release planning | `release.md` |
| Context-Driven Baseline | `baseline.md` |
| Agent selection | `profiles.md` |
| Run storage retention | `storage.md` |
| Run Worktree reconciliation | `reconcile.md` |
| Spec Consistency Check; Spec close audit | `spec.md` |
| Stopping Runs | `stop.md` |
| Detached Runs; Run discovery and Attach; Live Run View | `runs.md` |
| Supervisor Run Event Stream | `events.md` |
| A review only happens when it is asked for; User-Facing Review Runs; Review checkout and spec worktree isolation; Review Artifact Storage | `review-runs.md` |
| User-Facing Spec Runs | `implement.md` |
| Driving a Spec implementation loop; Autonomous Spec delivery | `spec-delivery.md` |
| Reopen Command; Settle Command | `settle.md` |
| Supersede Command; Archive Command | `archive.md` |

`SKILL.md` keeps, in this order:

1. its front matter, with the version raised;
2. the `# Roundfix` introduction;
3. the reference index, between the lines
   `<!-- roundfix:reference-index:begin -->` and
   `<!-- roundfix:reference-index:end -->`;
4. the `### QA settlement` block, from its heading to the next heading, moved
   out of whatever section holds it, with every byte unchanged;
5. the sections `Context-Efficient Evidence Boundaries`,
   `Assigned Review Issue Batches`, `Assigned Task Batches`,
   `Forbidden Actions` and `Completion Report`.

The index is the only new text. It holds one table row per reference file: a
relative link to the file, the commands it covers, and when to read it. It
also states that a change to a command edits that command's reference file.

An earlier Spec of the wave may have added a section the table does not name.
Its file is decided by this rule, in order:

1. the file of the top-level command its heading names;
2. else the file of the first `roundfix <command>` it shows;
3. else it stays in `SKILL.md`.

| Top-level commands | File |
| --- | --- |
| `init`, `setup`, `migrate`, `doctor`, `upgrade`, `skills` | `setup.md` |
| `fetch`, `resolve`, `watch` | `review-runs.md` |
| `review` | `review.md` |
| `implement` | `implement.md` |
| `deliver`, `window` | `deliver.md` |
| `settle`, `reopen`, `qa-report` | `settle.md` |
| `archive`, `supersede` | `archive.md` |
| `reconcile` | `reconcile.md` |
| `release` | `release.md` |
| `spec` | `spec.md` |
| `baseline` | `baseline.md` |
| `profiles` | `profiles.md` |
| `stop` | `stop.md` |
| `gc`, `storage` | `storage.md` |
| `runs`, `attach` | `runs.md` |
| `events` | `events.md` |

The Task records each placement the rule decided in its Result.

### The command reference layout

`docs/user-guide/commands.md` is cut at every `### ` heading outside a fenced
code block. Each such section moves whole to
`docs/user-guide/commands/<name>.md`. The name is the heading's first word when
that word is a top-level command, and otherwise the heading in lower case with
spaces replaced by hyphens. Two sections with the same name share a file, in
their original order. On `9e439dbb` the names are `setup`, `doctor`, `migrate`,
`review`, `deliver`, `upgrade`, `init`, `gc`, `skills`, `profiles`, `baseline`,
`fetch`, `resolve`, `watch`, `review-report-shape`, `window`, `implement`,
`reopen`, `settle`, `archive`, `qa-report`, `supersede`, `reconcile`, `runs`,
`attach` and `events`.

A second-level section with no `### ` child moves the same way, under the same
naming rule (`stop`, `detached-runs`), except `Global contract` and
`Agent boundaries`, which stay.

`commands.md` keeps its title and introduction, `Global contract`, each group
heading with the text before its first `### `, `Agent boundaries`, and a
command index between `<!-- roundfix:command-index:begin -->` and
`<!-- roundfix:command-index:end -->` with one row per file.

A relative link inside a moved section is rewritten to reach the same target
from the new directory: a target gains one leading `../`, and an in-page
anchor to a section that moved becomes that section's file. Links into a moved
section from elsewhere in the repository (three in
`docs/user-guide/usage.md`) are rewritten to the command's file. Nothing else
on a line changes, and no paragraph is re-wrapped.

### The move proof

Each splitting Task's Verification compares the document on `HEAD`, which is
the Task's starting commit while the Daemon verifies, with the working tree:

- it takes the non-blank lines of the old file's body and of the new files,
  leaving out the front matter and the marked index block;
- for the command reference it replaces every link target by an empty one on
  both sides;
- it sorts both lists and requires them to be identical.

A lost, duplicated or reworded line fails the comparison. The QA gate repeats
it on the splitting commit and its parent.

### API Contracts

1. API Contract: `roundfix skills check` — exits `0` for the shipped bundle.
   It reports `missing required wording "<phrase>"` for `roundfix/SKILL.md`
   only when no file of the Roundfix skill holds the phrase.
2. API Contract: `roundfix skills install` — writes every file of the Roundfix
   skill to the target, including each `roundfix/references/<name>.md`.

## Coverage Map

- Goal 1 → The skill layout; API Contract 2.
- Goal 2 → The skill layout; The command reference layout; Testing Approach 5.
- Goal 3 → The move proof.
- Goal 4 → Pins read the tree; API Contract 1.
- Core Feature 1 → The skill layout.
- Core Feature 2 → The command reference layout.
- Core Feature 3 → Interfaces; Pins read the tree; API Contract 1.
- Core Feature 4 → Testing Approach 5.
- Success Metric 1 → Testing Approach 3.
- Success Metric 2 → The move proof; Testing Approach 3, Testing Approach 4.
- Success Metric 3 → Testing Approach 1, Testing Approach 2.
- Success Metric 4 → Testing Approach 5.
- API Contract 1 → Pins read the tree; Testing Approach 2.
- API Contract 2 → System Architecture; Testing Approach 3.

## Integration Points

- **Agent runtimes.** Codex, Claude Code and OpenCode load `SKILL.md` and read
  a reference through their own file tools, from the link in the index.
  `agents/openai.yaml` keeps `entrypoint: SKILL.md`.
- **Adopting repositories.** `roundfix skills install` and the Baseline skill
  restoration copy the whole skill directory, so the references arrive with
  the next install. Nothing is removed from a target, and the old layout has
  no file the new one lacks.

## Testing Approach

1. **Reader.** A new `internal/mdtree/mdtree_test.go` over `fstest.MapFS`
   covers: the entry alone when the companion directory is absent; companion
   files appended in lexical order; nested directories and non-Markdown files
   ignored; a missing entry reported as an error.
2. **Wording check.** A new `skills/roundfix_wording_test.go` over
   `fstest.MapFS` covers: a required phrase that lives only in a reference
   satisfies the check; a phrase missing from every file is reported; banned
   branding in a reference is reported. `TestCheckValidatesRoundfixSkillArtifacts`
   stays green on the real bundle before and after the split.
3. **Skill layout.** A new `skills/roundfix_layout_test.go` reads the canonical
   skill and covers: the index names every file under `references/`; every
   index link resolves; `SKILL.md` stays within the byte budget, read from one
   named constant; `Install` writes each reference. The Verification also runs
   the move proof and every pin listed above.
4. **Command reference layout.** A new
   `internal/docscontract/commands_index_test.go`, under the `docscontract`
   tag, covers: the index names every file under `commands/`; every relative
   link in `commands.md`, `commands/*.md` and `usage.md` whose target lies
   under `docs/user-guide/` resolves. The Verification also runs the move
   proof and the pins that read the command reference.
5. **Authoring rule.** A new
   `internal/speccheck/wave_collision_command_files_test.go` covers: two
   same-Wave Tasks that declare different command files raise no
   `SC-WAVE-COLLISION`; two that declare the same command file raise it.
   Phrase checks cover the `write-tasks` skill, its mirror and its template,
   plus `make skills-sync-check`.

## Build Order

1. The reader, the wording check over it, and every pin switched to it,
   task_01 (depends on: none).
2. The Roundfix skill split and its layout tests, task_02 (depends on: 1).
3. The command reference split and its index tests, task_03 (depends on: 1).
4. The `write-tasks` rule and template, and the collision tests, task_04
   (depends on: 2, 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

Steps 2 and 3 share no file and can run in one Wave.

## Risks & Considerations

- **A pin that silently checks less.** A negative pin that kept reading only
  the entry file would pass on a forbidden phrase in a reference. Step 1
  switches every reader before anything moves, and the QA gate repeats the
  sweep.
- **A reworded line.** The move proof fails on any changed line, so a
  well-meant heading fix or re-wrap cannot land here.
- **A reader added by an earlier Spec of the wave.** The sweep command finds
  it, and task_01 names every file it changes in its Result.
- **A dead link stays dead.** Some links in the command reference already
  point at paths that no longer exist. The rewrite keeps each target, so this
  Spec neither repairs nor adds a broken link, and the link test covers only
  targets inside the user guide.
- **The version floor.** Raising the skill's version must follow the rule in
  force on the Task's starting commit. On `9e439dbb` two tests in
  `skills/skills_test.go` find the Roundfix skill's version line by the
  floor's literal, and a rehearsal of the split showed both fail once the
  version moves. task_02 changes them to find the line by the declared
  version.
- **Rehearsed before authoring.** The split, both move proofs, the mirror copy
  and `make baseline-digests` were rehearsed in a scratch clone at `9e439dbb`.
  The entry file came to about 12 KB, the proofs passed on the clean split and
  failed on a lost line, a duplicated section, a reworded line and an
  unchanged version, and the digest command reported no change.

## Decisions

- One file per command family for the skill and per command for the guide,
  with no generator. See ADR-0187.
- Contracts read the entry file with its companions as one text.
- The `### QA settlement` block and the assigned-Batch contract stay in the
  entry file, because every reader needs them.
