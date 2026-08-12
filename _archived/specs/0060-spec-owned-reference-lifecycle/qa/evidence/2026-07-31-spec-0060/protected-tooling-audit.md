# Protected-tooling audit

Build: `00ca18ee7c0fa2bbc31f00b98c41c4208170cf5f`

## Constraint sources

- Identifier strategy: not applicable; the PRD and TechSpec preserve basenames
  and Git history and create no Internal Identifier. Operative source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable; the only declared surface is docs.
  Operative source: `docs/agents/cli.md`.
- Active ADR: accepted ADR-0083 requires one history-preserving move to one
  primary owning Spec, post-adoption links for secondary Specs, and new
  promotions only.
- Tooling authority: the PRD at commit `397227f` and the TechSpec at commit
  `f6b284f` contain the exact bounded Skill and derived-pin authorization. Both
  commits are ancestors of Task 01 commit `01c80ea`.

## Commit path audit

`rtk git diff-tree --no-commit-id --name-only -r 01c80ea` returned the five
canonical authorial Skills, their five embedded counterparts, the seven
authorized derived digest artifacts, and `task_01.md`. No other path changed.

`rtk git diff-tree --no-commit-id --name-only -r 2b4e7af` returned only
`CONTEXT.md`, `docs/agents/docs-layout.md`, `docs/agents/spec-routing.md`, and
`task_02.md`.

`rtk git diff-tree --no-commit-id --name-only -r 00ca18e` returned only
`task_03.md`.

The chronological Task ancestry is `01c80ea` → `2b4e7af` → `00ca18e`.
There is no prerequisite-fix or consequent-fix commit to audit, and no fix is
folded into or ordered around the protected Task. Process-substitution
comparisons of every setup-owned marker block in the two edited agent guides
against base `37d57b2` exited 0.

## Deterministic regeneration

In disposable clone `/tmp/roundfix-qa0060-audit.utVK8D/repo`:

```text
rtk make baseline-digests
baseline-digests: no changes; derived artifacts already match their canonical sources
{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}

rtk git status --porcelain
<no output>
```

All five `rtk cmp -s` checks between `.agents/skills/<name>/SKILL.md` and
`skills/<name>/SKILL.md` exited 0. The authoritative gate's public
`roundfix skills check` also passed.
