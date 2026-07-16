---
name: setup-context-driven
description: Configure a repo for CONTEXT-driven development (the method explained in docs/user-guide/context-driven-development.md) — scaffold the full docs/ layout (_inbox, adr, agents, design, findings, handoffs, references, specs, user-guide) and the CONTEXT.md glossary, and seed the docs/agents/ usage guides (docs layout with the findings template, issue tracker, spec routing, domain docs, triage labels, autonomous work model, and optional Secondbrain guidance). Run when preparing a repo for the write-prd/write-tasks/implement pipeline or for Roundfix-driven autonomous work; re-run to audit and refresh only setup-owned managed content.
disable-model-invocation: true
metadata:
  category: setup
  tags: [workflow, prd, issues, planning, triage, repository-context, agents]
  version: 0.8.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Setup Context-Driven

Configure a repository for CONTEXT-driven development through the portable asset catalog and `scripts/context_setup.py`. The script is the source of truth for audit, apply, setup snapshots, managed markers, decisions, and finding codes.

## Asset map

- `assets/profiles/` selects a supported profile and canonical skill setup snapshot: `typescript-bun-monorepo`, `go-cli-tui`, or `rust-cli`.
- `assets/modules/` owns compact root pointers, supporting guides, rule IDs, required decisions, and required skills.
- `assets/templates/` stores generated repository content. Root blocks must stay short and point to `docs/agents/` guides.
- `assets/setups/` stores bundled canonical skill setup snapshots. Normal audit/apply uses only these bundled files; it never needs `~/dev/skills`, the network, or third-party Python packages.
- `references/` is workflow guidance for agents, not generated output.

## Audit-first workflow

1. Inspect the repository just enough to pick the likely profile. Use local files only: root instructions, `CONTEXT.md`, `docs/`, package files, and language/runtime markers.
2. Run audit before asking setup questions:

   ```bash
   python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --format json
   ```

   If a profile is known but the manifest is missing or stale, rerun with `--profile <profile-id>` to preview the selected composition.
3. Read the JSON result. Present only:
   - selected profile and ordered modules;
   - selected canonical skill setup name;
   - active or conditional modules and their trigger decisions;
   - blocking findings by code/path/action;
   - `plannedChanges`, including `state` and `condition` for conditional operations;
   - optional cleanup information only when audit was run with `--show-extra-skills`.
4. Do not dump generated Markdown by default. Mention that full templates live under `assets/templates/` if the user asks to inspect them.

## Decisions

Use stored compatible decisions from `docs/agents/setup-context.json` first. Ask only for `decision.required` findings, one decision code at a time, in the order the CLI reports them. This includes dependent questions introduced after an enabled capability changes the Decision Plan, such as `runtime.backend`, `runtime.design`, and `verification.gate` after `autonomous.enabled=true`. Do not ask again for a stored compatible value.

Question routing:

- `spec.scaffold` — confirm local `docs/specs/<feature-slug>/` is the planning source.
- `domain.layout` — ask whether the repository is `single-context` or `multi-context`.
- `triage.external` — ask only when external forge issues are relevant.
- `autonomous.enabled` — ask whether Supervisor-to-ACP Runtime delegation applies.
- `runtime.backend` — ask for the backend/default implementation runtime and model.
- `runtime.design` — ask for the design, UI, UX, or frontend runtime and model.
- `verification.gate` — ask for the command agents must run before completion claims.
- `language.generated` — generated repository content must be `English`.
- `secondbrain.enabled` — ask whether read-only local Secondbrain guidance must be generated.
- `adoption.*` — ask only after showing the existing unmarked file that would become setup-owned.

Record each answer by passing `--decision ID=VALUE` to apply. For newly introduced decisions, ask the new code once, then let the manifest carry it on later runs.

## Preview and apply

Apply is explicit and never runs before the user sees the managed change summary.

Before applying, tell the user:

- the profile, ordered modules, and canonical skill setup;
- every managed file or block that will be created, refreshed, or removed;
- that only `docs/agents/setup-context.json` and declared setup-owned Markdown boundaries can change;
- that repository-authored bytes outside managed markers remain untouched;
- that extra installed skills are informational only and this workflow never removes skills.

Ask for confirmation. After confirmation, run:

```bash
python3 .agents/skills/setup-context-driven/scripts/context_setup.py apply --repo <repo> --format json --profile <profile-id> --decision <id=value> ...
```

If apply returns `decision.required`, stop, present the unchanged selection and preview from the response, and ask only the next unresolved decision. If it returns blocking findings, report the code, path, and action; do not patch around the finding manually.

## Optional reports

Use `--show-extra-skills` only when the user asks to review installed skills outside the selected setup:

```bash
python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --format json --show-extra-skills
```

Report `skills.extra.installed` and `skills.local.untracked` as review information. Never suggest a removal command.

Use `sync-setups` only as a maintainer operation with an explicit canonical setups directory:

```bash
python3 .agents/skills/setup-context-driven/scripts/context_setup.py sync-setups --source-dir <canonical-setups> --check --format json
```

Normal repository setup does not require this checkout.

## Secondbrain

Secondbrain is opt-in through `secondbrain.enabled=true`. When enabled, apply creates one compact root pointer and `docs/agents/secondbrain.md`. When disabled, apply creates neither; if a previous setup owned those artifacts, apply removes only the marked managed Secondbrain block or guide content.

Generated Secondbrain guidance must remain read-only. It must require index-first lookup, `qmd query`, project-mirror caution, file citations, Hermes escalation for durable updates, and secret safety. It must forbid writes to the Secondbrain, `raw/`, and `projects/*/mirror/`.

## Completion

After apply, rerun audit in JSON mode. A clean setup has exit code `0`. Report remaining findings by code and path. Do not claim setup is complete without fresh audit evidence.
