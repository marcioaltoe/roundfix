# CLI readiness journeys

The current built binary exercised the equal-to-minimum state:

```text
rtk ./bin/roundfix skills check
exit 0
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec, write-tasks, setup-context-driven, implement-task, implement-spec, brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate

rtk ./bin/roundfix doctor
skills: ok (38 required: 14 Roundfix-owned, 24 external)
```

A scratch clone supplied independently built binaries and matching installed
Skill fixtures for the other public states.

One minor above the `0.0.2` minimum:

```text
roundfix-above skills check
exit 0
Roundfix skill check passed: ...

roundfix-above doctor
skills: ok (38 required: 14 Roundfix-owned, 24 external)
```

Below minimum:

```text
roundfix-below skills check
exit 1
roundfix/SKILL.md: below minimum: skill "roundfix" requires 0.0.2, found 0.0.1; upgrade: roundfix skills install --target project

roundfix-below doctor
skills: failed (below minimum: skill "roundfix" requires 0.0.2, found 0.0.1; next: roundfix skills install --target project)
```

No top-level version:

```text
roundfix-unversioned skills check
exit 0
Roundfix skill check unversioned: roundfix

roundfix-unversioned doctor
skills: unversioned (38 required: 14 Roundfix-owned, 24 external; unversioned: roundfix)
```

An unreadable installed `roundfix/SKILL.md` produced the same distinct
`unversioned` state. Moving the complete installed `roundfix` Skill directory
out of the repository produced:

```text
skills: failed (missing: roundfix; next: roundfix skills install --target project)
```

Doctor also reported unrelated adapter/profile readiness failures in this
sandbox, but it completed the independent `skills:` check in every run. No QA
row depends on the unavailable Agent Session readiness surface.

