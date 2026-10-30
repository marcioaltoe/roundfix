---
task: task_03
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the follow-up
command accept every documented message type from one table that its error,
its ordering fallback and both follow-up guides share; reads an absent Spec
Root inside the repository as an empty one in every command that enumerates
Specs while named-Spec commands keep their refusal (ADR-0348); and holds the
follow-up guides to the JSON field names the command emits.

Every row runs on this repository's own files, temporary repositories and
homes, the follow-up fixture's fake `gh` and fake judge, and temporary copies
of the audited head. No row reads or writes the real `~/.roundfix`, reads
another repository or `~/dev/skills`, reaches the network, sends a request to
a provider, the judge or a reviewer, or reads a credential. Every `go test`
it runs carries an explicit `-timeout`. Every row that runs `bin/roundfix`
sets `HOME` and `TMPDIR` to fresh temporary directories, because the CLI
loads `$HOME/.roundfix/config.yml`, passes `GOMODCACHE` and `GOCACHE` from
`go env` so no module cache lands in them, and removes them afterwards. The
QA never writes into the checkout while a test suite runs.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 3 and 4 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` over `./internal/config` and
   `./internal/cli` with the test names task_01 names in its Verification,
   and recording each pass.
3. MUST verify Success Metrics 1, 2 and 5 and Invariants 5 and 6 by running
   the test names task_02 names in its Verification, the guide contract test
   with `-tags docscontract`, and recording each pass.
4. MUST verify through the built binary, in one temporary repository with
   `docs/history/specs/` and no `docs/specs`, that `bin/roundfix follow-up
   --format json` reports the `specs` source `checked` with count `0`,
   that `bin/roundfix spec check` and `bin/roundfix history sanitize` exit
   `0`, that `bin/roundfix spec check 0001-missing` exits `2` naming
   `which does not exist`, and that a messages file with type `other` is
   listed with exit `0` while one with type `chore` exits `2` with the
   accepted list. No row compares one binary's output with another's byte
   for byte.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - make `ResolveSpecsRootOrAbsent` return the not-exist error again:
     `TestAnAbsentSpecRootReadsLikeAnEmptyOne` fails on its `follow-up`
     subtest, whose first run reports `specs` `unavailable`;
   - read absence with `os.Stat` instead of `os.Lstat`:
     `TestResolveSpecsRootOrAbsentStillRefusesADanglingLink` fails, because
     its `docs/specs` link to a missing directory becomes absent;
   - drop the in-repository ancestor check:
     `TestResolveSpecsRootOrAbsentStillRefusesAMissingExternalRoot` fails on
     its absolute root under another temporary directory;
   - remove the `RequirePresent` call from `runSpecCheckCommand` with
     slugs: `TestANamedSpecInAnAbsentSpecRootKeepsTheRefusal` fails on its
     `spec check <slug>` case, whose text no longer names `which does not
     exist`;
   - remove `other` from `followUpMessageTypes`:
     `TestFollowUpAcceptsEveryDocumentedMessageType` fails on its `other`
     message, and the guide contract test fails on the type sentence;
   - write the accepted list of the error by hand without `information`:
     `TestFollowUpRefusesAnUnknownMessageTypeNamingTheAcceptedOnes` fails
     on its `chore` case;
   - give `feature` the class `other` in the table:
     `TestFollowUpAcceptsEveryDocumentedMessageType` fails on the
     `feature` message's literal `priority_name`;
   - rename the JSON tag `requested_by` of `followUpItem`:
     `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` fails on the
     item sentence;
   - delete `cost_usd` from the judge sentence of the skill reference:
     `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` fails on the
     judge sentence.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST search `.agents/skills/roundfix/references/follow-up.md`,
   `skills/roundfix/references/follow-up.md` and
   `docs/user-guide/commands/follow-up.md` for `requester` and
   `blocking project`, and record that neither appears.
8. MUST record, as the outside-evidence row, the Git SCM Wiki GitFaq passages
   `_prd.md` → Acceptance evidence quotes (section "Can I add empty
   directories?", archived at
   `https://web.archive.org/web/20120723224559/https:/git.wiki.kernel.org/index.php/GitFaq`,
   retrieved 2026-10-30): "Directories are added automatically when adding
   files inside them. That is, directories never have to be added to the
   repository, and are not tracked on their own." and "If you really need a
   directory to exist in checkouts you should create a file in it.", and
   Quaestor's report kept in the adopted Backlog Entry under `references/`
   (ai-memory message `01a18c1a`, 2026-10-30). Compare them offline with
   Git's own behavior: in a temporary repository, commit
   `docs/specs/0001-x/_prd.md`, commit its removal, clone the repository
   locally, record that the clone has no `docs/specs`, and run the row 4
   follow-up check in that clone. The row is never blocked; name the Git SCM
   Wiki and the message id with its date as the source.
9. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    `CONTEXT.md` **Spec Root** cites ADR-0348 and cites no Spec.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      anything under `internal/baseline/assets/` or `docs/agents/`, any
      skill other than the Roundfix files the Tasks declare,
      `internal/cli/archive.go` and `internal/cli/cancel.go`;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and the binary rows.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside source, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every enumerating command reads an absent Spec Root as an empty one, and named Specs, dangling links and missing external roots stay refused.
- [ ] Every documented message type is accepted with its class; an undocumented one is refused with the list.
- [ ] Both follow-up guides name the JSON fields and types the command uses.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 3
- ADR-0348; ADR-0319
