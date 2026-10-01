---
status: approved
granted: 2026-09-30
action: tell the pre-PR reviewer the Delivery Conventions, validate every review finding against the candidate diff and those conventions before it can park a delivery, keep each dismissal in the record, review a second round as the delta of one Reviewer Lineage with a continued session, enforce the two-round ceiling, and describe it in the Roundfix skill
consuming: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/review.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0203

On 2026-09-30 the maintainer approved the program this Spec belongs to
("Três ondas"). Its Wave 8 names "revisor com sessão contínua + achados
validados": a reviewer with a continuous session and validated findings. The
same day the maintainer expressly authorized the skill files: "considere
autorizado a ajustar todas as skills se necessário". The design this Spec
records is the maintainer's: the delivery's commit conventions in the prompt,
a validation step that records dismissals with their reason, a session that
persists across the rounds of one item, and the two-round ceiling kept.

The set was measured with `GovernedPath` on `30e8504f`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/review.md` — after Spec 0194 this is
  where an agent learns `roundfix review`. It must describe the Delivery
  Conventions, the finding grammar, validation, the Reviewer Lineage, the
  continued session and the ceiling (task_04). `skills/roundfix/references/review.md`
  is its generated mirror and is not governed.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — an owned
  skill's content changes only with its version, and both version fields
  live in these two files (task_04). No other line of them changes.

## What is not governed

The following are ordinary:

- `internal/cli/review.go`, `internal/cli/deliver_workflow.go`,
  `internal/agent/agent.go`, `internal/agent/acpx_runner.go` and
  `internal/agent/acp_stream.go`;
- the new files `internal/cli/review_conventions.go`,
  `internal/cli/review_validation.go`,
  `internal/cli/review_convention_validator.go` and
  `internal/cli/review_lineage.go`;
- the existing tests `internal/cli/review_test.go`,
  `internal/cli/review_disposition_test.go`,
  `internal/cli/review_head_bound_test.go`,
  `internal/cli/review_merge_base_test.go`,
  `internal/cli/review_record_checkout_test.go`,
  `internal/cli/review_archived_spec_test.go` and
  `internal/cli/review_final_message_test.go`;
- the new test files `internal/cli/review_validation_test.go`,
  `internal/cli/review_convention_validator_test.go`,
  `internal/cli/review_lineage_test.go`,
  `internal/cli/review_session_test.go` and
  `internal/agent/agent_session_id_test.go`;
- `skills/roundfix/references/review.md`,
  `skills/testdata/owned-skill-versions.json` and
  `docs/user-guide/commands/review.md`.

`internal/cli/cli_test.go` and `docs/references/coverage-record.json` are
governed and are not touched. No help text, flag or exit-code table outside
`roundfix review` changes.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No change to the Pre-PR Review Policy, the providers, the review profile,
  `.roundfixrc.yml`, archived Specs, existing QA Reports or the `### QA
  settlement` section of any skill.
- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, or the CI workflows.
- The live Artifact Directories under `~/.roundfix` and the acpx session
  store under `~/.acpx` are never written by a Task or the gate. The QA gate
  reads the disposition ledger read-only.
- No test reaches a real reviewer, a real ACP adapter, a provider, the
  TypeSafe API or the network.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
