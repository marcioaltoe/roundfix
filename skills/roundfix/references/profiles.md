## Agent selection

Roundfix routes Agent work through Agent Selection Profiles. A profile is one
Preferred Selection plus a required ordered Fallback Chain. Project Config wins
over User Config, which wins over built-ins; a higher-scope profile replaces a
lower-scope profile as one object. Roundfix never reads or mutates
runtime-owned model configuration, credentials, or adapter settings.

Required built-ins:

- `general`, `backend`, and `qa`: preferred
  `codex / gpt-6.1-sol / high`, fallback `claude / opus / high`.
- `frontend`: preferred `claude / opus / high`, fallback
  `codex / gpt-6.1-sol / xhigh`.
- `review`: preferred `codex / gpt-5.6-luna / max`, fallback
  `codex / gpt-6.1-sol / high`.

Optional Task Type categories `data`, `infra`, `docs`, `test`, and `chore`
inherit the effective `general` profile when absent. If configured, they must
be complete. The Model Catalog recognizes `gpt-6.1-sol`, `gpt-6-astra`, `gpt-6-sol`,
`gpt-6-luna`, `gpt-5.6-sol`, `gpt-5.6-terra`, and `gpt-5.6-luna` as official Codex identifiers, plus the Claude identifiers the
adapter advertises: `opus`, `sonnet`, `claude-fable-5-1`, `haiku`, and `default`.
The adapter advertises Opus 5.5 as `opus[1m]`; the capability parser removes the
bracketed context suffix, so `opus` is the catalog value. Catalog validity is distinct
from advisory recommendation rank and from operational availability: exact
proof in the effective environment is the only readiness authority. Explicit
custom model strings, including adapter aliases, are sent to the ACP Runtime
verbatim for the same proof and do not enter an allowlist.

When an adapter advertises an independent reasoning control, Roundfix treats
every advertised Agent Model identifier as opaque. A bracketed identifier such
as `opus[1m]` is selectable exactly as printed, and its canonical prefix
`opus` remains selectable with a separate reasoning effort. The `[1m]` suffix
is a context-window annotation, not a reasoning effort; an explicit
`reasoning_effort: 1m` is rejected against the adapter's advertised efforts.
Adapters without an independent reasoning control keep the existing
`canonical[effort]` variant encoding.

Project Config and User Config use the profile structure:

```yaml
profiles:
  general:
    preferred:
      runtime: codex
      model: gpt-6.1-sol
      reasoning_effort: "high"
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: "high"
  backend:
    preferred:
      runtime: codex
      model: gpt-6.1-sol
      reasoning_effort: "high"
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: "high"
  frontend:
    preferred:
      runtime: claude
      model: opus
      reasoning_effort: "high"
    fallbacks:
      - runtime: codex
        model: gpt-6.1-sol
        reasoning_effort: "xhigh"
  qa:
    preferred:
      runtime: codex
      model: gpt-6.1-sol
      reasoning_effort: "high"
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: "high"
  review:
    preferred:
      runtime: codex
      model: gpt-5.6-luna
      reasoning_effort: "max"
    fallbacks:
      - runtime: codex
        model: gpt-6.1-sol
        reasoning_effort: "high"
```

Use the profile management commands for inspection, writes, and disposable
proof:

```bash
roundfix profiles show --category backend --json
roundfix profiles configure --scope project --file profiles.yml --dry-run --json
roundfix profiles validate --json
```

`profiles show` is read-only and returns `roundfix/profiles/v2` JSON with the
effective source, inherited source, Preferred Selection, ordered fallbacks, and
the Recommended Profile. Each of the ten Agent Work Categories has one dated
`2026-09-30`: its Preferred Selection at rank 1 with role `preferred`, then its
Fallback Chain with role `fallback`. Rows include the selection, source date,
and rationale. Interactive configure prints the same advisory rows. The
Recommended Profile never selects, routes, proves availability, or writes config.

`profiles configure` prepares the candidate in memory, validates it, and
exact-proves each distinct Preferred Selection and fallback before
confirmation. `--file` reads a strict profile fragment, Interactive Input
collects one complete profile, `--dry-run` performs proof without writing and
reports `changed: false`, and `--json` returns
`roundfix/profiles-configure/v1`. Proof failure, cleanup failure, decline, and
output failure preserve target bytes. It preserves unrelated config and never
edits runtime-owned settings or credentials.

`profiles validate` is read-only proof through disposable ACP Sessions. It
deduplicates exact tuples, reports every category reference, closes every
disposable session on success or error, sends no prompt, creates no Run, and
returns `roundfix/profiles-validate/v1` JSON with tuple-level status.

Non-interactive Agent-starting commands (`resolve`, `watch`, and `implement`)
accept exactly two selection forms: omit every selection flag to use profiles,
or provide a complete one-Run Preferred Selection override:

```bash
roundfix watch --source coderabbit --pr 123 --until-clean
roundfix resolve --pr 123 --agent codex --model gpt-5.6-sol --reasoning-effort high --no-input
roundfix implement --spec example-spec --agent claude --model opus --reasoning-effort xhigh --detach
```

`--agent`, `--model`, and `--reasoning-effort` are all-or-none. A partial subset
exits `2` before config load, adapter or profile proof, Session creation,
worktree or artifact creation, or Run persistence. A complete override replaces
only the Preferred Selection for every relevant category and preserves each
configured Fallback Chain. If one override applies across multiple Task or QA
categories, Roundfix emits a warning. An explicit empty
`--reasoning-effort ""` counts as present and requests model-managed reasoning;
an explicit empty `--model ""` is invalid and exits `2`. Never reinterpret a
rejected explicit `high` as model-managed reasoning.

Before an operational Run mutates state, Roundfix validates Task Types, resolves
the relevant profiles, deduplicates exact preferred/fallback tuples, proves
them sequentially through disposable sessions, and closes those sessions.
`fetch` remains Agent-free. `resolve` and `watch` use only `review`; `implement`
derives its categories, including `qa`, from Task Graph metadata and accepts no
per-run QA selection.

After Run creation, automatic fallback is notification-first and pre-prompt
only. If selection start fails before the first prompt, Roundfix records the
failed attempt, publishes `agent_selection_fallback`, renders the same notice on
stderr/TUI/Attach/Run Event Stream, and only then activates the next configured
fallback in order. Once `agent_work_started` is recorded, there is no fallback
for prompt, tool, verification, cancellation or rate-limit failure; they keep
their normal failure semantics. A Lost Rollout is the one exception and gets no repair: the Task continues in a new Agent
Session, using the next Fallback Selection before the First Handoff or while a
QA report is pending, and the same selection otherwise. An Agent Session owner
gets at most two recoveries; the third settles the Task as runtime
infrastructure.

Legacy `defaults.agent` and `runtimes.<runtime>.model` /
`runtimes.<runtime>.reasoning_effort` remain readable only for scopes without a
`profiles` section. A same-scope mix fails with migration guidance. Migrate by
removing `defaults.agent` and `runtimes`, writing complete profiles with
`roundfix profiles configure --scope user|project --file <path>`, then running
`roundfix profiles validate`.

Legacy profile migration is separate from adapter migration. If the effective
Codex or Claude command fails official lineage proof, use `roundfix setup` to
diagnose it and authorize the applicable official pinned override. Setup and
Doctor use the same Adapter Readiness contract; unknown and below-pin lineages
fail with the applicable official install action, and neither command treats a
same-name executable as proof.

Initial progress and the Live Run View show the concrete stored selection:

```text
Agent: Codex
Agent Model: gpt-5.6-sol
Default Reasoning Effort: high
```

Attach reads compatibility summary values and per-scope selection history from
the Run Database, not from current User Config or Project Config. Legacy Runs
that predate per-scope selection history render it as unavailable instead of
inventing records.

`specs.root` is a User Config and Project Config key that defaults to
`docs/specs`. Project Config overrides User Config, which overrides the
built-in default. Relative values resolve against the repository root of the
user's checkout; absolute values are used as-is. Roundfix resolves the Spec
Root once per command and carries that absolute path into Run and Task
Worktrees, so Worktrees read and write the same Spec artifacts as the checkout.
Validation rejects an empty root, a missing root, or a root that is not a
directory, and the Preflight Validation message names the resolved path. When
the resolved root is outside the repository working tree after symlink
evaluation, Spec artifacts are external and stay out of code-repository
commits. A non-default root at Implement Run startup prints one stderr line:

```text
Spec Root: <path>
```

`logs.agent` is a User Config and Project Config key that defaults to `false`.
When it is false, Roundfix writes no per-Batch Agent log files. The Run Event
Journal still records every Agent payload, and `--no-agent-console` only hides
Agent-source console events from non-TTY stderr. Set the key to `true` for
development or debugging when file logs are useful:

```yaml
logs:
  agent: true
```

With `logs.agent: true`, per-Batch Agent log files use
`<artifact_dir>/runs/<run-id>/agent/batch-<nnn>.log`. The Detached Run console
log remains unconditional and is not controlled by `logs.agent` (ADR-0030).

`notify.enabled` is a User Config and Project Config key that defaults to
`true`; `notify.command` defaults to the empty string. Project Config overrides
User Config, which overrides the built-in defaults. With the default empty
command, Roundfix uses the native desktop notifier when available: `osascript`
on macOS, `notify-send` on Linux, and a silent no-op on other platforms or when
the native tool is missing. A non-empty `notify.command` replaces the native
path and runs through the shell with a 30s timeout. The command receives the
completed Run context in `ROUNDFIX_RUN_ID`, `ROUNDFIX_OUTCOME`,
`ROUNDFIX_KIND`, and `ROUNDFIX_TARGET`; targets are `pr:<number>` for review
Runs and `spec:<slug>` for Spec Runs. Terminal context adds
`ROUNDFIX_REASON`, `ROUNDFIX_CONSOLE_LOG`, `ROUNDFIX_ATTACH_COMMAND`,
`ROUNDFIX_REVIEW_ISSUES_KNOWN`, and `ROUNDFIX_NEXT_ACTION`. Set
`notify.enabled: false` to disable outcome notifications entirely.
