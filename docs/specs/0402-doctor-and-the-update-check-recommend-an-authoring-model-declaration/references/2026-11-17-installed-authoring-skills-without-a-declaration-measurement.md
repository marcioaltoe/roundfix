# Installed authoring skills without a declaration: measurement of 2026-11-17

The author read the Roundfix adopter checkouts under `~/dev` on the Mac mini on 2026-11-17, read-only: no file, ref or branch of those repositories changed, and `git status --porcelain` was empty in Argus before and after. A checkout counts as an adopter when it holds `docs/agents/setup-context.json` or `.roundfixrc.yml`.

## What was read

For each adopter, at its checked-out `HEAD`:

- which of `write-idea`, `write-prd`, `write-techspec` and `write-tasks` has `.agents/skills/<name>/SKILL.md`;
- whether `.roundfixrc.yml` holds the line `  authoring:` under `skills:`;
- which of the four the repository publishes under `skills/<name>/SKILL.md`;
- how many `roundfix-*` files `.claude/agents/` holds.

## Results

| checkout | HEAD | installed authoring skills | `skills.authoring` declared | published authoring skills | generated subagents |
| --- | --- | --- | --- | --- | --- |
| argus | `fc53eb7` (`main`) | all four | no | none | 0 |
| quaestor | `e188944` (`main`) | all four | yes | none | 5 |
| roundfix | `bcd716df` (`main`) | all four | yes | all four | 5 |
| skills | `095b79c` (`main`) | none | no | none | 0 |

Argus is the case this Spec answers. Its `.roundfixrc.yml` has no `skills:` key, `.claude/agents/` does not exist, and none of its four installed authoring skills holds a `model:` or `effort:` line. Their versions are `write-idea` 0.0.7, `write-prd` 0.0.18, `write-techspec` 0.0.20 and `write-tasks` 0.0.32.

With the installed Roundfix 0.147.1, `roundfix update check --format json` run in Argus printed, reduced with `jq` to the fields that matter here:

```json
{"installed":"0.147.1","state":"current","baseline":"plan_ready","authoringModels":null,"next":"apply"}
```

`authoringModels` is absent, so the report says nothing about the four authoring skills that run on whatever model the session uses. `roundfix doctor` was not run in Argus, because it probes the ACP adapters and the agent runtimes; its code (`internal/cli/doctor.go`, `loadedAuthoringState`) prints the `authoring-models:` line under the same condition that adds `authoringModels`.

## Limits

- Only the checkouts on the Mac mini were read. The other adopters live on the MacBook and were not reached.
- The measurement records the state on 2026-11-17. Argus may declare the section before this Spec ships; the QA gate reproduces the recorded state in a fixture instead of reading Argus.
