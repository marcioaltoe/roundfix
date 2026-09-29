---
task: task_04
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: completed
type: backend
complexity: high
---

# Task 04: The report names the user-flow binary

## Overview

The QA Agent completes the report the Daemon seeded, and it runs public-CLI
rows with whichever `roundfix` it finds. Four QA Reports of 2026-09-25 to
2026-09-28 rewrote the Daemon's `auditing_binary` to name the Agent's own build,
so the field named two binaries, and nothing checks which binary the public rows
exercised. The report is written by the QA Agent in the Spec's `qa/` directory
and read by `settleQAVerdict` in `internal/daemon/task_engine.go`, which
settles the QA Task, and by the maintainer who ships the Spec.

## Requirements

1. MUST add `UserFlowBinary` to `spec.QAReport` in `internal/spec/qa.go`, read
   from the optional front matter key `user_flow_binary`; a report without it
   stays readable with an empty value.
2. MUST add `SelfAudit bool` to `spec.AuditorEvidence` in
   `internal/spec/auditor_evidence.go`, true when the binary's build commit,
   with `-dirty` removed, is a commit object in the audited repository, and
   false for a binary with no build commit.
3. MUST make the QA step, after `writeMechanicalQAReport`, read the seeded
   report back and keep its `AuditingBinary` and `AuditorStaleness`, and pass
   them, the evidence and the audited head (`git rev-parse HEAD` in
   `plan.WorkDir`) to `settleQAVerdict`, whose only caller is the QA step.
4. MUST make `settleQAVerdict` refuse a `pass` or `partial` that
   `QAReportEligibility` accepts when its `auditing_binary` or
   `auditor_staleness` is present and differs from the seeded value, with a
   cause containing `auditor fields are Daemon-owned`.
5. MUST make `settleQAVerdict`, in a self-audit only, refuse a `pass` or
   `partial` that `QAReportEligibility` accepts when `user_flow_binary` is
   missing, names no build commit, or names a build commit that is not a prefix
   of the audited head, with a cause containing `user_flow_binary`. The build
   commit is the text after the first `(` up to the first `,` or `)`, with
   `-dirty` removed, and it MUST be at least seven hexadecimal characters.
   Outside a self-audit the key MUST NOT be checked.
6. MUST settle each refusal with the existing reason form
   `QA verdict <verdict> not accepted: <cause>`, and MUST leave
   `QAReportEligibility`, `roundfix archive`, `roundfix settle` and
   `roundfix qa-report accept` unchanged, so
   `TestArchivedPassCorpusRemainsArchiveEligible` stays green and unchanged.
7. MUST append to the QA prompt, in a self-audit only, the line
   `Self-audit: build roundfix from this Run Worktree with make build, run every public-CLI row with ./bin/roundfix and never a roundfix found on PATH, and record its --version line as user_flow_binary.`
8. MUST put new spec tests in `internal/spec/qa_user_flow_binary_test.go`:
   - `TestReadQAReportReadsTheUserFlowBinary`;
   - `TestReadQAReportWithoutAUserFlowBinaryStaysReadable`;
   - `TestAuditorEvidenceMarksASelfAudit`: a build commit of the fixture
     repository, with and without `-dirty`;
   - `TestAuditorEvidenceOutsideTheBuildRepositoryIsNotASelfAudit`: a build
     commit absent from the repository, and an empty one.
9. MUST put new daemon tests in `internal/daemon/qa_user_flow_binary_test.go`,
   injecting `Dependencies.Auditor`; the self-audit cases use a build commit of
   the fixture repository, the others one absent from it:
   - `TestQASettlementAcceptsTheSeededAuditorFields`: the Agent turns the seed
     into a `pass` with a Results row and keeps both lines; `completed`;
   - `TestQASettlementRefusesARewrittenAuditingBinary`;
   - `TestQASettlementRefusesARewrittenAuditorStaleness`;
   - `TestSelfAuditSettlementAcceptsAUserFlowBinaryBuiltFromTheAuditedHead`: a
     `user_flow_binary` of `roundfix 0.17.0 (<short head>-dirty, built <time>)`
     settles `completed`;
   - `TestSelfAuditSettlementRefusesAMissingUserFlowBinary`;
   - `TestSelfAuditSettlementRefusesAUserFlowBinaryFromAnotherCommit`;
   - `TestSelfAuditSettlementRefusesAUserFlowBinaryWithoutABuildCommit`: a
     released-style `roundfix 0.17.0`;
   - `TestSettlementOutsideASelfAuditIgnoresTheUserFlowBinary`: no
     `user_flow_binary` and a foreign one both settle `completed`;
   - `TestSelfAuditQAPromptNamesTheUserFlowBinary`, and
     `TestQAPromptOutsideASelfAuditOmitsTheUserFlowBinary` for the mirror case.
   Each refused case asserts the QA Task settles `failed` with a reason
   containing its cause.
10. MUST state in `.agents/skills/qa-gate/SKILL.md` that the auditor fields are
    Daemon-owned (that phrase), so the gate keeps the seeded `auditing_binary`
    and `auditor_staleness` lines; that the gate records the `--version` line
    of the binary its public rows ran as `user_flow_binary` (add the key to the
    section 6 template); and that in a Roundfix self-audit that binary is built
    from the audited tree and never taken from PATH. MUST regenerate
    `skills/qa-gate/SKILL.md` with `make skills-sync`.
11. MUST name `user_flow_binary` in the QA Report entry of `CONTEXT.md` and say
    there that the auditor fields are Daemon-owned.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A report that rewrites a seeded auditor field is not accepted by Daemon
      settlement; one that keeps them is.
- [ ] In a self-audit a `pass` settles only with a `user_flow_binary` built from
      the audited head, and the QA prompt says so; outside a self-audit the key
      is not checked.
- [ ] Archive, settle and `qa-report accept` keep today's acceptance.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `internal/spec/qa_user_flow_binary_test.go`
- creates: `internal/daemon/qa_user_flow_binary_test.go`
- interface: `internal/spec/qa.go`
- interface: `internal/spec/auditor_evidence.go`
- interface: `internal/daemon/task_engine.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReadQAReportReadsTheUserFlowBinary|TestReadQAReportWithoutAUserFlowBinaryStaysReadable|TestAuditorEvidenceMarksASelfAudit|TestAuditorEvidenceOutsideTheBuildRepositoryIsNotASelfAudit|TestArchivedPassCorpusRemainsArchiveEligible|TestQASettlementAcceptsTheSeededAuditorFields|TestQASettlementRefusesARewrittenAuditingBinary|TestQASettlementRefusesARewrittenAuditorStaleness|TestSelfAuditSettlementAcceptsAUserFlowBinaryBuiltFromTheAuditedHead|TestSelfAuditSettlementRefusesAMissingUserFlowBinary|TestSelfAuditSettlementRefusesAUserFlowBinaryFromAnotherCommit|TestSelfAuditSettlementRefusesAUserFlowBinaryWithoutABuildCommit|TestSettlementOutsideASelfAuditIgnoresTheUserFlowBinary|TestSelfAuditQAPromptNamesTheUserFlowBinary|TestQAPromptOutsideASelfAuditOmitsTheUserFlowBinary)$" ./internal/spec ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReadQAReportReadsTheUserFlowBinary TestReadQAReportWithoutAUserFlowBinaryStaysReadable TestAuditorEvidenceMarksASelfAudit TestAuditorEvidenceOutsideTheBuildRepositoryIsNotASelfAudit TestArchivedPassCorpusRemainsArchiveEligible TestQASettlementAcceptsTheSeededAuditorFields TestQASettlementRefusesARewrittenAuditingBinary TestQASettlementRefusesARewrittenAuditorStaleness TestSelfAuditSettlementAcceptsAUserFlowBinaryBuiltFromTheAuditedHead TestSelfAuditSettlementRefusesAMissingUserFlowBinary TestSelfAuditSettlementRefusesAUserFlowBinaryFromAnotherCommit TestSelfAuditSettlementRefusesAUserFlowBinaryWithoutABuildCommit TestSettlementOutsideASelfAuditIgnoresTheUserFlowBinary TestSelfAuditQAPromptNamesTheUserFlowBinary TestQAPromptOutsideASelfAuditOmitsTheUserFlowBinary; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "auditor fields are Daemon-owned" && grep -q "user_flow_binary" .agents/skills/qa-gate/SKILL.md && grep -q "user_flow_binary" CONTEXT.md && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "auditor fields are Daemon-owned" && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the fourteen new tests exists and neither the skill nor the glossary names `user_flow_binary`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The user-flow binary

## Result

### Implementation

- `QAReport` now reads the optional `user_flow_binary`, and
  `AuditorEvidence.SelfAudit` records whether the auditing binary's normalized
  build commit resolves to a commit object in the audited repository.
- The QA step retains the seeded auditor fields, auditor evidence, and audited
  head for its one settlement call. Eligible `pass` and `partial` reports now
  refuse rewritten daemon-owned fields and, only for self-audits, require a
  hexadecimal `user_flow_binary` build commit that prefixes the audited head.
- Self-audit prompts name the worktree build and `./bin/roundfix` contract. The
  canonical qa-gate skill, its generated mirror, and the QA Report glossary
  entry document the same ownership and `user_flow_binary` contract.
- Added the fourteen named spec and daemon tests, with separate cases for each
  rewritten field, missing/foreign/released-style user-flow identities, the
  non-self-audit exemption, and both prompt branches.

### Focused checks

- Red signal: the focused daemon run for rewritten auditor metadata, missing
  self-audit metadata, and the self-audit prompt failed because both reports
  settled `completed` and the prompt omitted the required line.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run
  'Test(ReadQAReport.*UserFlowBinary|AuditorEvidence.*SelfAudit)'
  ./internal/spec` — passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run
  'Test(QASettlement|SelfAuditSettlement|SettlementOutsideASelfAudit|SelfAuditQAPrompt|QAPromptOutsideASelfAudit)'
  ./internal/daemon` — passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1
  ./internal/spec ./internal/daemon` — passed.
- `make skills-sync` — passed and regenerated `skills/qa-gate/SKILL.md` from
  the canonical skill.
- `GOCACHE=/private/tmp/roundfix-task04-gocache make verify-incremental` — the
  sandboxed run reached the suite but two Unix force-stop integration tests
  could not enumerate the process table; the permission-enabled rerun passed,
  including `go vet`, all Go packages, skill tests, skill checks, and the build.
- The Task's declared `## Verification` command was not run; the Daemon owns
  that command and terminal settlement.

### Acceptance evidence

- Seed ownership: `TestQASettlementAcceptsTheSeededAuditorFields`,
  `TestQASettlementRefusesARewrittenAuditingBinary`, and
  `TestQASettlementRefusesARewrittenAuditorStaleness` passed in the focused
  daemon suite.
- Self-audit binary: the valid-head, missing, foreign-commit, no-build-commit,
  outside-self-audit, and prompt mirror cases passed in the focused daemon
  suite.
- Existing acceptance remains isolated in `QAReportEligibility`, which was not
  changed; the full `internal/spec` and `internal/daemon` suites and the
  permission-enabled incremental repository gate passed.
