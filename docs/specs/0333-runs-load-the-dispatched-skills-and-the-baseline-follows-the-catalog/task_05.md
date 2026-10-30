---
task: task_05
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec moves the Setup Snapshots
to catalog commit `86df860` and retires seven skills the catalog removed and
the owned `implement-spec`; rewords the dispatch so person-started and polish
skills start only on request and splits product UI from public pages; and
makes Run Agents load only the skills the dispatch selects for their Task and
name them in their Result.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, reads `~/dev/skills` or another repository, reaches the
network, sends a request to a provider, the judge or a reviewer, or reads a
credential. Every `go test` it runs carries an explicit `-timeout`. Every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`, and passes
`GOMODCACHE="$(go env GOMODCACHE)"` and `GOCACHE="$(go env GOCACHE)"`; each
temporary directory is removed afterwards with `chmod -R u+w` then `rm -rf`.
The QA never writes into the checkout while a test suite runs. No row
compares one command's output with another's byte for byte; where a row
checks a printed line, temporary paths are placeholders.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification and `make verify-docs`, recording each result as a gate fact
   (Success Metric 5).
2. MUST verify Success Metric 7's mechanism and Invariant 2 through the
   built binary: in a temporary Git repository made from `git archive` of the
   audited head and committed, create `.agents/skills/tdd/SKILL.md` and
   `.agents/skills/implement-spec/SKILL.md` (a frontmatter `name` line
   only), run `bin/roundfix baseline update --repo <tmp> --no-history
   --format text`, and record exit `0`, `Baseline update: current`, `Skills
   retired: 2`, and one `- retired implement-spec:` and one `- retired tdd:`
   line, each with `no longer required by the Baseline; delete
   .agents/skills/<name>`; record that both trees still exist afterwards.
3. MUST verify Success Metric 4 through the built binary: `bin/roundfix
   skills list` names twelve Roundfix-owned skills and no `implement-spec`;
   in a temporary directory, `bin/roundfix skills install --target project
   --dir <tmp>` writes no `implement-spec` tree.
4. MUST verify Success Metrics 1 to 4 and Invariants 1, 3 and 4 by running
   `go test -count=1 -timeout 20m -v` over `./internal/baseline`, `./skills`,
   `./internal/agent`, `./internal/spec`, and over `./internal/daemon` and
   `./internal/cli` with the test names task_03 names in its Verification,
   and recording each pass of the tests tasks 01 to 04 name; and that
   `bin/roundfix baseline update --repo . --no-history --format text`, with a
   temporary `HOME`, reports `Skills differing: 0` and no `Skills retired`
   line on this repository.
5. MUST record, on the audited head: every `github` snapshot under
   `internal/baseline/assets/setups/` names
   `86df8608726979064a8fe46003dfaf403ffe0087`; no file under
   `internal/baseline/assets/modules/`, nor `skill-activations.json`, nor any
   snapshot names one of the ten Retired Skills of `_techspec.md` → Data
   Model 1 (search with `grep -rlwF` per name, excluding nothing else); and
   this repository's `docs/agents/skill-dispatch.md` names neither
   `implement-spec` nor a Retired Skill, and its `grill-with-docs` trigger
   starts with "When the person asks".
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - set the `typescript` snapshot's `activationBundles` back to the
     delivery base's (`git show <delivery base>:internal/baseline/assets/setups/typescript.json`):
     `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles` fails on
     `bundle.testing` and `TestNoCatalogEntryNamesARetiredSkill` fails on
     `tdd`;
   - put `tdd` back into `bundle.testing` of `skill-activations.json`:
     `TestNoCatalogEntryNamesARetiredSkill` fails;
   - reword `trigger.context-workflow.grill-with-docs` back to
     "Stress-testing a plan or decision against repository documents.":
     `TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks` fails;
   - add `frontend-design` to `bundle.frontend-ui-quality`:
     `TestTheFrontendBundlesSeparateProductUIFromPublicPages` fails;
   - drop the `skillDispatchGuidePath` entry from
     `assembleTaskContextBundle`:
     `TestAssembleTaskContextBundleListsTheSkillDispatchGuideWhenPresent`
     fails on its present-guide case;
   - drop the first Data Model 4 invariant from `taskExecutionInvariants`:
     `TestBuildTaskPromptTellsTheAgentToLoadOnlyTheDispatchedSkills` fails;
   - drop `docs/agents/skill-dispatch.md` from `CarryForwardInputs`:
     `TestCarryForwardRefusesATaskWhoseBindingTextChanged/skill_dispatch_guide`
     fails on the exit code;
   - add `implement-spec` back to `skillNames` with its tree restored from
     the delivery base: `TestImplementSpecIsNotARoundfixOwnedSkill` fails.
7. MUST execute the tests tasks 01 to 04 name in their Verification.
8. MUST record, as the outside-evidence row, the passages `_prd.md` →
   Acceptance evidence quotes, and compare them offline with this Spec's
   evidence:
   - Claude Code documentation, "Extend Claude with skills": "In a regular
     session, skill descriptions are loaded into context so Claude knows
     what's available, but full skill content only loads when invoked.",
     against the 0-of-108 frontend sessions of `_prd.md` and the passing
     `TestBuildTaskPromptTellsTheAgentToLoadOnlyTheDispatchedSkills`, whose
     prompt now tells the Agent to load the matching skills;
   - Anthropic, "Agent Skills" overview: "Progressive disclosure ensures only
     relevant content occupies the context window at any given time.",
     against the prompt's "never load a skill that matches nothing";
   - the catalog messages quoted in the two adopted Backlog Entries under
     `references/` (`01a18d69`, 2026-10-30 12:24 BRT, and `01a18da4`,
     2026-10-30 13:29 BRT), against the snapshots' revision and the Retired
     Skills.

   The row is never blocked; it names both documentation pages and both
   message ids with their dates as the sources.
9. MUST record Success Metric 6 and Success Metric 7 as one row each,
   `blocked (declared unreachable)`, against the PRD's two Unreachable
   Acceptance entries, and count two in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    `CONTEXT.md` defines **Loaded Skill List** citing ADR-0346, that **Spec
    Context Bundle** names the skill dispatch guide, and that no entry cites
    a Spec.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the CI workflows,
      `go.mod`, `go.sum`, `.roundfixrc.yml`, any `Makefile` line other than
      `OWNED_SKILLS` and its comment, any upstream skill other than
      `context7-cli`, `grilling` and `tui-design`, any owned skill other
      than `implement-task`, `qa-gate`, `roundfix`, `archive-spec`,
      `write-tasks` and the deleted `implement-spec`, any module other than
      `typescript`, `backend`, `frontend` and `context-workflow`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" TMPDIR="$(mktemp -d)" bin/roundfix archive 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, remove both
    temporary directories afterwards, and write the report's Outcome as one
    paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification and the documentation gate.
- [ ] Run the binary rows in temporary repositories and the named tests.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside sources, the declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] `baseline update` lists a held `tdd` and `implement-spec` as retired without deleting them or changing its state.
- [ ] The binary ships twelve owned skills and the snapshots follow `86df860` with no Retired Skill named.
- [ ] Prompts, bundle and carry-forward carry the load rule, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Data Models 1-4; API Contracts 1-9; Invariants; Testing Approach; Build Order 5
