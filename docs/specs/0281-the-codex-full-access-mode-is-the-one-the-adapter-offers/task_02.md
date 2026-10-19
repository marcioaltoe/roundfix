---
task: task_02
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
status: completed
type: docs
complexity: medium
---

# Task 02: Describe the full-access selection and the sandboxed browser driver

## Overview

task_01 makes a full-access request select the mode the Agent Session
advertises and describes it in the Roundfix `profiles` reference. This Task
records it in the decision trail and the remaining readers' documents:
ADR-0011, the configuration guide, the glossary and the `.roundfixrc.yml`
comment, which still cites codex-acp 1.1.9. It also gives the `qa-gate` skill the browser guidance of
the adopted source's addendum. It is verifiable on its own through phrase
checks, the skill version record and the Behavior Surface Record.

## Requirements

1. MUST answer the Secondbrain Inbox Entry of 2026-10-18,
   `references/2026-10-18-agent-full-access-pede-um-modo-que-o-codex-acp-recusa.md`
   (`_prd.md` → Sources 1, including its addendum of 2026-10-18 16:15 BRT), by
   writing API Contract 7 of `_techspec.md` and the documentation of API
   Contracts 1 to 5 (ADR-0293).
2. MUST append to `docs/adr/0011-acp-runtimes-run-batches-with-full-access.md`
   the paragraph `**Amended by [ADR-0293](0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md) (2026-10-19).**`.
   It says the runtime's full-access mode is now the first candidate the
   Agent Session advertises (`agent-full-access`, then `full-access`, for
   Codex), that a session advertising none is reported degraded instead of
   failing, and that the Codex sandbox preset is sent when the session
   advertises `sandbox_mode`. Set its `updated_at` to the edit's time and
   change nothing else in it.
3. MUST rewrite "Access policy readiness" in
   `docs/user-guide/configuration.md` to describe API Contracts 3 to 5. It
   names the candidates `agent-full-access`, `full-access` and
   `bypassPermissions`, the degraded text `no advertised mode grants full
   access`, the status `full_access_mode_not_advertised`, and the remedy:
   disable full access or select a runtime that advertises a full-access
   mode. The `defaults.agent_full_access` table row MUST say that `true`
   requests the full-access mode the session advertises.
4. MUST replace the comment above `agent_full_access` in `.roundfixrc.yml`, as
   API Contract 7 states. It no longer cites codex-acp 1.1.9 or says that a
   true value fails every Task. It says the repository keeps the Agent sandbox
   on purpose, that a true value selects the mode the session advertises
   (`agent-full-access` on codex-acp 2.x), and that the QA gate drives
   browsers through the chrome-devtools MCP server. The line
   `agent_full_access: false` MUST stay byte-identical.
5. MUST add to `.agents/skills/qa-gate/SKILL.md` the heading `### Browser
   driver inside a sandbox`, under "4. Exercise real flows" after "Frontend
   sweep", with the text of API Contract 7. It names agent-browser as the
   example daemon-backed tool and the chrome-devtools MCP server as the
   driver to prefer. No text inside any `### QA settlement` section may
   change.
6. MUST add **Full-Access Mode** to `CONTEXT.md` through domain-modeling,
   linking ADR-0293. It is the adapter mode an Agent Session is switched to
   when full access is requested: the first of the runtime's full-access
   candidates that the session advertises in its `mode` config option.
   `_Avoid_: sandbox mode, YOLO mode`.
7. MUST search the guides, the agent guides, the ADRs outside
   `docs/history/` and the skill copies for text that names `full-access` as
   the only Codex mode or cites codex-acp 1.1.9, and qualify each.
8. MUST raise both version fields of `.agents/skills/qa-gate/SKILL.md`
   (`metadata.version` and the top-level `version`, today `0.0.17`) to
   `0.0.18` by hand, because the record command refuses changed content
   under an unchanged version. Then the Agent itself MUST run these
   commands, in this order, before handing back: the Daemon only runs the
   Verification, it never regenerates a mirror or a record. No record or
   mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised `qa-gate` skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Amend ADR-0011 and rewrite the configuration guide's access policy section.
- [ ] Replace the `.roundfixrc.yml` comment.
- [ ] Add the `qa-gate` heading and the glossary term.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] ADR-0011 and the configuration guide describe the advertised-mode selection and the degraded report.
- [ ] `.roundfixrc.yml` no longer cites codex-acp 1.1.9 and keeps `agent_full_access: false`.
- [ ] The `qa-gate` skill names the chrome-devtools MCP server for a sandboxed browser, and `CONTEXT.md` defines **Full-Access Mode**.

## Consumers

None. This Task changes no Shared Contract: it edits documentation, a
configuration comment and skill text, and the suites that read them run in
its Verification.

## Context

- instruction: `docs/adr/0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `docs/adr/0011-acp-runtimes-run-batches-with-full-access.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.roundfixrc.yml`
- interface: `CONTEXT.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `tr -s '[:space:]' ' ' < docs/adr/0011-acp-runtimes-run-batches-with-full-access.md | grep -qF -- '**Amended by [ADR-0293](0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md) (2026-10-19).**' || { printf 'missing amendment note in ADR-0011\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'agent-full-access' || { printf 'missing agent-full-access in the configuration guide\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'full_access_mode_not_advertised' || { printf 'missing full_access_mode_not_advertised in the configuration guide\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'no advertised mode grants full access' || { printf 'missing degraded text in the configuration guide\n' >&2; exit 1; }; ! grep -qF -- '1.1.9' .roundfixrc.yml || { printf '.roundfixrc.yml still cites codex-acp 1.1.9\n' >&2; exit 1; }; grep -qF -- 'agent-full-access' .roundfixrc.yml || { printf 'missing agent-full-access in .roundfixrc.yml\n' >&2; exit 1; }; grep -qxF -- '  agent_full_access: false' .roundfixrc.yml || { printf 'agent_full_access: false changed\n' >&2; exit 1; }; grep -qF -- '### Browser driver inside a sandbox' .agents/skills/qa-gate/SKILL.md || { printf 'missing qa-gate heading\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- 'chrome-devtools MCP server' || { printf 'missing chrome-devtools MCP server in qa-gate\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Full-Access Mode**:' || { printf 'missing glossary term Full-Access Mode\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/config || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task ADR-0011 has no amendment note, so the command fails at its first check. Afterwards the phrases, the unchanged setting, the glossary term, the skill mirrors and versions, the configuration tests and the documentation contracts, including the Behavior Surface Record, pass.

## References

- `_prd.md` → Sources; User Story 3; Core Feature 4; Glossary
- `_techspec.md` → API Contract 7; API Contract 5; Current behavior; Build Order 2
- ADR-0293; ADR-0011; ADR-0189

## Result

- Updated ADR-0011 with the ADR-0293 amendment and the edit timestamp; rewrote
  the configuration guide's access policy readiness text and default row to
  describe advertised-mode selection, the degraded policy, status, and remedy.
- Replaced the `.roundfixrc.yml` comment while preserving the byte-exact
  `agent_full_access: false` setting; added the sandbox browser-driver
  guidance to both skill copies and added the linked **Full-Access Mode** term
  to `CONTEXT.md`.
- Ran `make skills-sync` successfully; the mirror is synchronized. Raised both
  `qa-gate` version fields to `0.0.18` and ran
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` successfully (`1 passed`), recording the digest.
  Ran the Behavior Surface Record command successfully:
  `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  (`1 passed`).
- Focused checks: `git diff --check` passed; the post-edit diff contains only
  the assigned documentation, configuration, skill, generated-record, and
  Daemon-owned task paths; the QA settlement section has no diff; the search
  for the obsolete `codex-acp 1.1.9` citation is clear in the scoped guides,
  agent guides, ADRs, skills, `.roundfixrc.yml`, and `CONTEXT.md`.

Acceptance evidence:

- ADR-0011 and `docs/user-guide/configuration.md` contain the advertised
  candidates and `full_access_mode_not_advertised` degraded report.
- `.roundfixrc.yml` contains `agent-full-access`, no `1.1.9`, and the exact
  `  agent_full_access: false` line.
- `.agents/skills/qa-gate/SKILL.md` and its synchronized mirror name the
  chrome-devtools MCP server under `### Browser driver inside a sandbox`, and
  `CONTEXT.md` defines **Full-Access Mode** with the ADR-0293 link.

## Carry-forward provenance

- Source Run: `run_20261019T205953Z_8c49586070d0475f`
- Source commit: `7fcd40d38d0bfcdfa667e051fd7c5b3b43a829c8`
