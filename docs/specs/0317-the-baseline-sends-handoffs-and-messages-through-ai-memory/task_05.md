---
task: task_05
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec moves every upstream Setup
Snapshot to the upstream commit that replaced `handoff` with six
`ai-memory-*` skills, makes the core module dispatch `ai-memory-handoff` and
`ai-memory-messaging`, replaces the eight Secondbrain inbox clauses with
ai-memory clauses of the same force, keeps the Secondbrain read-only, removes
the follow-up command's `secondbrain-inbox` source, deprecates its two keys,
and moves the authoring skills, the glossary and the history guide to the new
channel.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, reads `~/dev/skills` or the Secondbrain, calls ai-memory,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or reads another repository. Every `go test` it runs
carries an explicit `-timeout`, and every row that runs `bin/roundfix` sets
`HOME` and `TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Invariant 1: run
   `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
   `TestTheCoreModuleDispatchesTheAiMemorySkills`,
   `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill` and
   `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract`, and record
   from the audited `internal/baseline/assets/setups/*.json` that every
   `github` snapshot names `65d5f365cace27c4b468db1723213723b10c870b`, and
   that the `go` snapshot's names match `_prd.md` → Acceptance evidence
   (the six `ai-memory-*` entries, no `handoff`).
3. MUST verify Success Metrics 2 and 3 and Invariants 2 and 3: run
   `TestEveryInboxClauseAnAdopterHoldsIsReplaced` and
   `TestTheGuidesSendCrossProjectWorkThroughAiMemory`, and record that
   `grep` finds none of `Inbox Entr`, `inbox/<destination>` and `_triaged`
   in this repository's `docs/agents/*.md`. Exclude the Spec's own test
   files and planning files from any such search, because they name the
   retired text on purpose.
4. MUST verify Success Metric 4, Invariants 4 and 5 and Surface Transcript 1
   end to end: with
   `HOME` and `TMPDIR` at fresh temporary directories, write
   `$HOME/.roundfix/config.yml` naming `secondbrain.root` (pointing at a
   temporary directory that holds `inbox/<name>/pending.md`) and
   `secondbrain.inbox_namespace`, create a fresh temporary Git repository
   with one commit and `docs/specs/`, and run `bin/roundfix follow-up
   --format json` there; record exit `0`, the two warning lines of
   `_techspec.md` → API Contract 6 on standard error, the source names in
   API Contract 5's order without `secondbrain-inbox`, no item for the
   pending file, and `"items":[]`. Record the `pull-requests` source as it
   comes; it never decides the row.
5. MUST verify Success Metric 5: run `go test -timeout 20m -count=1
   ./internal/baseline ./skills ./internal/config ./internal/cli` and the
   documentation contracts with `-tags docscontract`, and, with a binary
   built from the audited head and fresh `HOME` and `TMPDIR`, run
   `bin/roundfix baseline update --repo <clone> --no-skills --no-history
   --yes --format text` in a temporary `git clone --no-local` of this
   checkout detached at the audited head, after the suites finish, and record
   `File changes: 0`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - remove `"clause.secondbrain.inbox-write-permission"` from the
     `replaces` list of `clause.core.ai-memory-is-optional` in
     `internal/baseline/assets/modules/core.json`:
     `TestEveryInboxClauseAnAdopterHoldsIsReplaced` reports it unaccounted
     (it is the only clause naming it);
   - in `internal/baseline/assets/modules/secondbrain.json`, put back "is an
     Inbox Entry addressed to the Baseline's owner;" in
     `clause.secondbrain.baseline-owned-guidance`:
     `TestTheGuidesSendCrossProjectWorkThroughAiMemory` reports
     `docs/agents/secondbrain.md` (its input enables the Secondbrain);
   - replace ", when the session exposes the ai-memory MCP tools." with "
     when needed." in both core triggers:
     `TestTheCoreModuleDispatchesTheAiMemorySkills`;
   - remove the `secondbrain.inbox_namespace` entry from
     `deprecatedConfigKeys` in `internal/config/config.go`:
     `TestTheRetiredSecondbrainKeysLoadWithAWarning/user_root_and_namespace`
     and `/project_namespace` fail;
   - remove the empty-`secondbrain`-mapping removal from
     `internal/config/config.go`: every case of
     `TestTheRetiredSecondbrainKeysLoadWithAWarning` fails;
   - restore `followUpInbox` and its call in
     `internal/cli/follow_up_inventory.go` from the Delivery Base with a
     User Config `secondbrain.root` read again:
     `TestFollowUpReadsNoSecondbrainInbox` fails on the source list.
   Record any sabotage that the catalog's digest check stops before the
   named test runs as failing at catalog load, naming the message.
7. MUST execute the tests task_01 to task_04 name in their Verification.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - because this Spec's authorization forbids reading `~/dev/skills` in QA
     and network access, record the upstream commit's setup lists as
     `blocked (environment: upstream checkout and network prohibited by Spec
     authorization)`, after comparing the audited snapshots with the
     description the PRD quotes; such a row never decides a qualifying
     `partial`;
   - record the ai-memory architecture page as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the passage the PRD
     cites, and record that the audited clause
     `clause.core.ai-memory-is-optional` covers a refused send.
9. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary and documentation checks:
    - `CONTEXT.md` defines **Cross-project Message**, **Triage** and
      **Follow-up Item** as the Tasks state, defines neither the Inbox Entry
      nor the Secondbrain Root, and links no deleted file
      (`docs/agents/domain.md`);
    - `docs/adr/0095-the-secondbrain-inbox-is-the-single-capture-door.md` is
      absent, ADR-0330 is `accepted` and names ADR-0095 as superseded;
    - search `docs/user-guide/`, `.agents/skills/roundfix/`, the three
      authoring skills and `docs/agents/` for `secondbrain-inbox`,
      `secondbrain.root` (outside the deprecation text), `Inbox Entr` and the
      `handoff` skill, and record what it finds.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every file under
      `internal/baseline/assets/source-baselines/` and
      `internal/baseline/assets/retention/`, every profile source other than
      `standard-typescript-monorepo.json`'s formatter digest, and every
      upstream skill outside the twelve restored and `handoff`;
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
    `HOME="$(mktemp -d)" TMPDIR="$(mktemp -d)" bin/roundfix archive 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end follow-up row and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every upstream snapshot follows the ai-memory commit, the core module dispatches the ai-memory skills, and every retired inbox clause is `replaced`.
- [ ] `roundfix follow-up` reads no Secondbrain inbox and loads a config naming the retired keys with one warning each.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0317-the-baseline-sends-handoffs-and-messages-through-ai-memory/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Surface Transcript 1; Invariants 1-6; Testing Approach; Build Order 5
