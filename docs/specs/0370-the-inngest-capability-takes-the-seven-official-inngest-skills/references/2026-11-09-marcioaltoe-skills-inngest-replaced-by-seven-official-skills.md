# ai-memory message 01a1c01f-b6ad from marcioaltoe-skills, 2026-11-09 11:44Z — reported by marcioaltoe-skills, nature: measured

Subject: skills catalog: inngest replaced by 7 official inngest-* skills — adjust Baseline

The marcioaltoe/skills catalog removed the skill `inngest` and replaced it with seven official skills from inngest/inngest-skills (Apache-2.0). The change is published on main: PR #156, squash commit e45466d2b9e15480eb286a4c19b513494a867bce.

Slug changes:
- `inngest` → removed. Replacements, all at skills/05-implementation-loop/<slug>, update mode auto: `inngest-setup`, `inngest-durable-functions`, `inngest-steps`, `inngest-events`, `inngest-flow-control`, `inngest-cli`, `inngest-middleware`.
- Presets that named it: `setups/typescript.txt`. It now lists the seven slugs in place of `inngest`.

Unrelated change in the same window, with no slug change: `pdf` now updates automatically from anthropics/skills (PR #154, commit 7835ffc).

Impact found in the Roundfix repo, read-only: `inngest` is still named in internal/baseline/assets/modules/backend.json (the skill list and the `trigger.backend.inngest` dispatch trigger), internal/baseline/assets/setups/typescript.json and go-cli-typescript-bun.json, the golden fixture internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md, and parity-corpus fixtures under internal/baseline/testdata/parity-corpus/v1/fixtures/.

Request: update the Baseline setup snapshots and the backend module so the optional Inngest capability installs and dispatches the new slugs instead of `inngest`. Suggested trigger mapping: building or changing functions and steps → inngest-durable-functions and inngest-steps; events and fan-out → inngest-events; concurrency, throttling and debounce → inngest-flow-control; project setup → inngest-setup; local Dev Server and CLI → inngest-cli; logging, tracing and dependency injection → inngest-middleware. Then refresh the fixtures and validate. Marcio decides when to publish it there. No reply expected.

Evidence field: marcioaltoe/skills PR #156, squash commit e45466d2b9e15480eb286a4c19b513494a867bce on main (merged 2026-11-09T11:44:13Z); setups/typescript.txt; skills-registry.json
