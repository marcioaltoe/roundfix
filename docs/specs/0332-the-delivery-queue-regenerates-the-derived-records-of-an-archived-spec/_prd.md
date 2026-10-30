---
spec: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
status: active
created: 2026-10-30
surfaces: [backend]
---

# The Delivery Queue regenerates the derived records of an archived Spec

When a delivered item's Pull Request conflicts with `main` only in derived
records, the Delivery Queue merges the default branch and reruns each record's
regeneration command. The coverage record command always failed there. It
rewrote `docs/references/coverage-record.json` and then exited 1, because the
suite guard refused the write. The item parked `delivery-error`, and the
operator merged by hand. The cause is not the archive. The item's Archive
Record still sanctions the command, and the guard reads it. The cause is that
the coverage record writer never names its command to the guard, so no record
can sanction its write. After this Spec the writer names its command, and the
queue resolves such a conflict without the operator. A write that no record in
the tree sanctions is still refused.

## Sources

One Backlog Entry, adopted from `docs/backlog/` by this Spec's planning
change (see `references/_index.md`):

1. [`references/2026-10-27-the-queue-cannot-regenerate-the-coverage-record-of-an-archived-spec.md`](references/2026-10-27-the-queue-cannot-regenerate-the-coverage-record-of-an-archived-spec.md),
   `type: fix`: the park of Spec 0329's item, Pull Request #709, on
   2026-10-27 at about 18:10 BRT. Its `Where` section blames the guard's
   reader for "active Specs" only. The characterization under Validation
   corrects that diagnosis, and this Spec fixes the cause it found.

At the authoring base `docs/backlog/` held only this entry and
`docs/findings/` held none.

## Prerequisites

None in `requires`. Authored against `35de7f74` (v0.126.0) with no active
Spec. Specs 0333, 0334 and 0335 are authored after this one and deliver after
it. This Spec goes first, so that the queue can resolve, without the operator,
derived-record conflicts that arise while they deliver.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the fixture
  Archive Records use the existing `NNNN-<slug>.md` record name. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is touched; every test and QA row runs offline in temporary
  repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0149 makes the guard exempt only declared outputs under the declared command, ADR-0149: "exempts only those paths under only those commands"; ADR-0192 makes the queue regenerate derived paths on a conflict without granting authority, ADR-0192: "It grants no authority"; ADR-0081 keeps the authorization drawn around the edit that causes a regeneration, ADR-0081: "the authorization stays drawn around the cause". ADR-0345, new, records the decision and relates it to the Archive Record's regeneration list, ADR-0345: "Every test that writes a derived record names its command." Listed because they cite one of those ADRs, and not applying: ADR-0233 (skill version lines at merge, unchanged), ADR-0250 (the coverage record's platforms, unchanged), ADR-0256 (release waits for skills, unchanged), ADR-0197, ADR-0272, ADR-0279 and ADR-0286 (pre-PR review rounds, records and dismissals, not touched). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AS decisions of 2026-10-30 quoted in the authorization record, over
  the Governed Paths this Spec declares; no Makefile, workflow, `go.mod`,
  `go.sum`, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec/_authorization.md`;
  bounded files:
  `internal/spec/coverage_test.go`,
  `docs/agents/specific-repository.md`,
  `docs/references/coverage-record.json`.

## Goals

- A conflict that is only in derived records resolves in the queue, the
  coverage record included, after the item is archived.
- The sanction for that regeneration is the delivered item's own Archive
  Record. Without it the write is still refused.
- Every derived record the queue regenerates is accounted for. Each writer in
  a guarded package names its command, and the others are shown to need no
  name.

## Core Features

1. The coverage record writer names
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   to the suite guard, only while `-update-coverage-record` is set.
2. The regeneration contract (`TestRegenerationIsDeclared`) runs every
   Project Config regeneration in a repository copy. That copy is
   sanctioned only by an archived item's Archive Record, its coverage record
   is stale, and no declaration is injected. A second copy, whose record does
   not list the coverage command, must have its coverage regeneration refused.
3. A queue-level fixture is an archived item whose Pull Request conflicts only
   in the coverage record. The queue's conflict recovery merges it and commits
   the regenerated record. If the item's record lacks the command, the
   recovery fails and leaves the item unchanged.
4. The repository rule and the glossary record the decision:
   `docs/agents/specific-repository.md` and **Archive Record** in
   `CONTEXT.md`, citing ADR-0345.

## Skills

- unchanged: config keys — only the regeneration contract test in `internal/config` changes; no configuration key, default, validation or message changes, and the Behavior Surface fingerprint reads no test file.

## Non-Goals / Out of Scope

- Narrowing the guard's sanction set to the delivered item's record alone.
  The guard keeps reading every record in the tree. ADR-0345: "The guard
  keeps reading every record the tree holds."
- Changing the queue's merge, its park reasons or its log names
  (`derived-regeneration-<n>.log`).
- Guarding `./skills` or `./internal/docscontract`, or changing their record
  commands.
- Changing `.roundfixrc.yml`, the Makefile or any authorization reader.

## Success Metrics

1. Success Metric: an archived item's record alone sanctions the coverage
   command in a repository copy. Run on a stale coverage record, the command
   exits 0 and rewrites the record.
2. Success Metric: with a record that does not list the command, the same run
   exits non-zero and prints
   `suiteguard: modified: docs/references/coverage-record.json`.
3. Success Metric: the queue's conflict recovery resolves the archived item's
   coverage-record-only conflict with one merge commit carrying the
   `Roundfix-Delivery: derived-merge` trailer. Without the sanction, it
   returns `regenerate derived paths` and leaves the worktree clean at the
   item head.
4. Success Metric: every Project Config regeneration in the copy exits 0 and
   changes only declared paths. Removing the writer's declaration makes the
   contract fail.
5. Success Metric: after release, the next Delivery Queue item whose Pull
   Request conflicts only in derived records, the coverage record among them,
   merges without the operator.

## Acceptance evidence

The outside-evidence row rests on two sources this Spec did not write.

The first is the park of Pull Request #709 that the adopted Backlog Entry
quotes. The queue ran the coverage record command after the archive commit,
the suite guard printed `suiteguard: modified:
docs/references/coverage-record.json`, and the operator merged by hand.

The second is the Go command documentation for `go test` (`cmd/go`,
https://pkg.go.dev/cmd/go/internal/test, retrieved 2026-10-30), quoted by
the QA row: "Tests that open files within the package's module or that
consult environment variables only match future runs in which the files and
environment variables are unchanged." This passage is why the fixtures copy
the repository from inside the test process. A copy made by a child process
would leave the cached result blind to a later change in the writer.

## Unreachable Acceptance

- criterion: Success Metric 5, after release, the next Delivery Queue item
  whose Pull Request conflicts only in derived records, the coverage record
  among them, merges without the operator.
  reason: it needs a real Pull Request conflict on GitHub after release,
  which no Verification or QA row may create.
  satisfied-by: the queue log of the first such item in cycle AS (Specs 0333
  to 0335), recorded by the operator in the queue interventions log.

## Validation

Characterized on 2026-10-30 at `35de7f74` in a scratch clone. Commands and
logs are in the operator's scratch directory `spec-0332/`.

- `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
  over a coverage record with one test entry removed rewrote the record and
  exited 1 with `suiteguard: modified: docs/references/coverage-record.json`.
  The tree held no active Spec, and seven Archive Records, 0329's among them,
  list that command with that output.
- A scratch test file that only called
  `suiteguard.DeclareSanctionedRegeneration` with that command made the same
  run exit 0. So the Archive Records already sanction the write, and the only
  missing piece is the writer's declaration. The file was removed afterwards,
  and the tree was clean.
- `rg -n "DeclareSanctionedRegeneration" -g '*.go'` finds declarations in
  `internal/baseline` (module versions, Baseline digests), in the suite
  guard's own tests, and in the regeneration contract's fixture. That fixture
  writes an `init` declaring the coverage command into the copied
  `internal/spec`, which hid the defect. `internal/spec` itself declares
  nothing.
- Three variants were run in repository copies. In each, every copied
  Archive Record was removed and one archived item's record was written, a
  module was edited and the coverage record was made stale, and all seven
  `.roundfixrc.yml` regenerations were run:
  - without a declaration, the coverage command exited 1;
  - with a declaration and the command in the record, all seven exited 0;
  - with a declaration and the command missing from the record, the
    coverage command exited 1 again.

  The Baseline commands passed on the item's record alone.
- `internal/suiteguardcontract.readArchiveRecordRegenerations` reads
  `regeneration` from every `docs/history/specs/*.md` record. Spec 0329's
  record lists the coverage command with
  `docs/references/coverage-record.json`.
- `TestRegenerationIsDeclared` takes 58 s on this machine.
- Secondbrain, read only: `wiki/index.md`, and
  `qmd query "roundfix suite guard sanctioned regeneration archive record derived paths conflict" --all --files --min-score 0.3`.
  It found mirrors of this repository's backlog entry, archive guide and
  Archive Records. It also found two wiki digests of this repository's
  history: historical authorizations and regeneration ownership
  (2026-09-09), and external validation of Archive Records (2026-10-16).
  None of them records the writer-side declaration.
- Exa: the `go test` caching passage quoted under Acceptance evidence.

## Glossary

- changes: **Archive Record**

## Decisions

Maintainer answers to the cycle AS questions (AskUserQuestion, 2026-10-30),
verbatim:

- 2026-10-30: Plan — "Aprovar como está (Recommended)", option text: "Quatro
  Specs na ordem validada pelo Jev; a 0332 entra primeiro para que a fila
  resolva sozinha conflitos de registros gerados durante a 0333."
- 2026-10-30: Authorization — "Mesmas regras do AR (Recommended)", option
  text: "Push de branches, PR dos planos, entrega e merge pela Delivery
  Queue, eu escolho o ponto de release e aviso os adotantes por ai-memory."
- 2026-10-30: Release — "Uma no fim do ciclo (Recommended)", option text:
  "v0.127.0 depois das quatro Specs; os adotantes recebem a revisão de skills
  e as correções juntas."

Autonomous route: no person present for the technical choices; decisions taken from the operator's cycle AS brief for Spec 0332 of 2026-10-30 and the characterization above.

Taken from those sources, with no further interview:

- The cause is the writer's missing declaration, not the archive (ADR-0345).
- After archive, the item's Archive Record is the sanction. That is not the
  authorization at `source_revision`, and not the queue's declarations
  (ADR-0345).
- The guard keeps its rule, and an unsanctioned write is still refused
  (ADR-0345, ADR-0149).

## Open Questions

None.
