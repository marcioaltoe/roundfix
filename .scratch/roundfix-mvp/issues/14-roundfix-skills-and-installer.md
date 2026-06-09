---
title: "Skills: ship Roundfix watch and resolve-round skills"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 59
  - 60
  - 61
blocked_by:
  - 08-agent-runtime-batch-resolution.md
  - 11-watch-review-round-loop.md
---

# Skills: ship Roundfix watch and resolve-round skills

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Ship the user-facing and child-agent skills that make Roundfix usable from local
agent workflows. The watch skill should steer users toward Roundfix commands for
review cleanup, while the resolve-round skill should encode the bounded Batch
contract for child Agents.

## Acceptance criteria

- [ ] The watch skill instructs agents to prefer `roundfix` commands over manual
      GitHub scraping and to report Run identity, pull request, Review Source,
      Agent, and state.
- [ ] The resolve-round skill instructs child Agents to read assigned issue
      files, triage issues, make valid fixes, update assigned statuses, and run
      verification.
- [ ] The resolve-round skill forbids commits, pushes, Review Source mutations,
      unassigned issue edits, and daemon-owned `duplicated` status changes.
- [ ] Skill installation supports Codex, Claude Code, and OpenCode-compatible
      directories where the MVP supports them.
- [ ] Skill text and generated artifacts use Roundfix names only and do not copy
      reference-project branding.
- [ ] Tests or scripted checks validate skill presence, required sections, and
      forbidden-action wording.

## Blocked by

- 08-agent-runtime-batch-resolution.md
- 11-watch-review-round-loop.md
