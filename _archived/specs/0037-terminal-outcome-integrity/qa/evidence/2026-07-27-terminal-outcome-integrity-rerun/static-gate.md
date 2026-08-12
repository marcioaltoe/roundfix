# Static gate

Build: `ef6eb44ad8951112b1c3641bb7fd21793b440f95`

Command:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0037-rerun-gocache make verify
```

Result: exit `0`.

```text
Go test: 2477 passed in 23 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: roundfix, write-idea, write-prd,
write-techspec, write-tasks, setup-context-driven, implement-task,
implement-spec, brainstorming, council, business-analyst, archive-spec,
qa-gate, evidence-gate
rtk go build -buildvcs=false ... -o bin/roundfix ./cmd/roundfix
```

Post-gate checks:

- `rtk git -c core.fsmonitor=false diff --check`: exit `0`, no output.
- `rtk git -c core.fsmonitor=false status --short`: only the QA report and its
  evidence directory are changed; verification introduced no unrelated delta.
- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`: exit
  `0`, no output.
- `rtk ./bin/roundfix skills check`: exit `0`; every shipped Skill contract
  passed.

Verdict: pass.
