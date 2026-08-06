---
status: done
created_at: 2026-08-05
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-agent-selection-and-execution-environments.md
---

# 2026-08-05 — `agent_full_access` passes config validation and fails every Task

status: pending

## What was observed

Spec 0014's QA gate blocked 10 of its 22 rows and its own evidence file named
the remedy:

```
Unblocking action: rerun QA in a full-access session that permits localhost TCP
connections to the already provisioned scratch PostgreSQL service.
```

Setting `defaults.agent_full_access: true` in `.roundfixrc.yml` is the only
documented way to ask for that. Config validation accepted it. Every Task in
the next Run then failed instantly, before any Agent work, with the same
reason:

```
Agent failed: set acpx Agent Session mode "full-access": acpx infrastructure
error after exit code 1: Agent rejected session/set_mode for mode
"full-access": Invalid params (ACP -32602). The adapter may not implement
session/set_mode, or the requested value is not supported.
```

The Run settled 0 completed, 3 failed, 4 skipped. Measured on codex-acp 1.1.9,
with `roundfix doctor` reporting every check green immediately before.

Reverting the setting and relaunching produced 6 completed Tasks and a QA
report from the same graph, unchanged otherwise.

## Root cause

The setting is validated for shape but never proved against the runtime that
has to honour it. Doctor proves Node, acpx, adapter lineage and version, every
Agent Selection Profile tuple through disposable sessions, the Repository Skill
Set, and codex hygiene — an unusually thorough readiness contract — and none of
it touches session mode.

The failure then lands at the worst possible moment: after Run creation, per
Task, with no fallback, because `agent_full_access` is not part of the Agent
Selection lifecycle that fallback covers. A configuration typo is caught before
a Run exists; a configuration value the adapter cannot honour costs the whole
Run.

## What would settle it

Prove the mode the same way profiles are proved. Doctor already opens
disposable ACP Sessions for every distinct tuple; when `agent_full_access` is
true, one of those sessions can attempt `session/set_mode` and report:

```
full-access: failed (adapter codex-acp 1.1.9 rejected session/set_mode with
ACP -32602); next: set defaults.agent_full_access to false, or use an adapter
that implements session/set_mode
```

Two smaller improvements stand alone:

- **Fail at preflight, not per Task.** If the mode cannot be set, no Task will
  ever run; refusing before Run creation turns a dead Run into an actionable
  exit `2`.
- **Say what full-access is for, and whether it works.** The QA evidence
  contract tells operators to rerun in a full-access session. If no supported
  path to one exists for the configured adapter, the prescription sends them
  into exactly this wall.

## Related

[[2026-08-04-a-spec-archives-with-pass-while-a-user-story-was-never-exercised]]
is the cost this setting was meant to remove: rows that need a running service
stay unverified without it, and it cannot currently be turned on.

## Spec pointer

None yet.
