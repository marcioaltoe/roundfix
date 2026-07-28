---
status: done
created_at: 2026-07-28
updated_at: 2026-07-28
---

# profiles configure — a one-category fragment deletes every other configured profile (2026-07-28)

`roundfix profiles configure --scope project --file <fragment>` replaces the
entire `profiles:` map with the fragment instead of merging the named
categories into it. Configuring one category therefore **silently deletes every
other configured profile**.

This contradicts the command's documented contract, which states it "preserves
unrelated config and never edits runtime-owned settings or credentials".

## Reproduction

Project Config before — five configured categories (`general`, `backend`,
`frontend`, `qa`, `review`), each with a Preferred Selection and Fallback
Chain. Fragment naming only `frontend`:

```yaml
frontend:
  preferred:
    runtime: claude
    model: opus
    reasoning_effort: xhigh
  fallbacks:
    - runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
```

```console
$ roundfix profiles configure --scope project --file frontend.yml --yes --json
{"schema":"roundfix/profiles-configure/v1","changed":true,…}
$ git diff --stat -- .roundfixrc.yml
 1 file changed, 30 insertions(+), 73 deletions(-)
```

After: `profiles:` contains `frontend` only. `general`, `backend`, `qa`, and
`review` are gone, along with every Fallback Chain they carried. The report says
`changed: true` and names only the category you asked for; nothing announces the
deletion.

Restoring from a pre-change copy and editing the single value by hand produces
the honest diff: **3 lines, not 103.**

## Why it is easy to miss

The obvious verification does not catch it. A reviewer checking that top-level
and second-level keys survived sees `profiles:` still present and concludes
nothing was lost — the destruction is one level deeper, in that key's children.
A `fluxus` session on 2026-07-27 performed exactly this check, reported "no
loss: first and second level keys identical, 20 comments preserved, the only
value change is claude-opus-5 → opus", and moved on. **Any repository whose
Project Config was written by this command with a partial fragment should be
audited for missing profile categories.**

The 103-line diff also disguises the loss: it reads as reformatting (the command
also reindents 2 → 4 spaces), so a reviewer attributes the churn to style and
stops reading.

## Secondary defect — a declined write exits 0

Without `--yes` and without a TTY, the command prints its confirmation prompt,
writes nothing, reports `"changed": false`, and **exits `0`**. Automation that
checks the exit code reads a silent no-op as success. `changed: false` is honest
but only in the JSON body; the exit code does not distinguish "declined" from
"applied". For a CLI meant to be driven by Agents, refusing to act deserves a
non-zero exit or an explicit refusal channel.

## Suggested resolution

1. Merge the fragment's named categories into the existing map instead of
   replacing it. Deleting a category should require naming it explicitly.
2. Render the effective change before confirmation as a per-category summary —
   added, replaced, and **removed** — so a destructive write cannot be mistaken
   for a formatting diff.
3. Preserve the file's existing indentation and key order, writing only what
   changed. A config-writing command that reformats its target makes every diff
   unreviewable.
4. Exit non-zero when a confirmation is declined, or require `--yes` in
   non-interactive contexts rather than defaulting to a silent no-op.

## Suggested acceptance checks

- Configuring one category leaves every other configured category byte-identical.
- A fragment that would remove a category is refused, or reports the removal in
  the confirmation summary.
- A single-value change produces a single-value diff.
- A declined non-interactive write is distinguishable from a successful one by
  exit code alone.

## What worked — keep

- The proof-before-write behavior is the reason this command is still the right
  surface: it caught an invalid Agent Selection before it could fail a Run. Keep
  that; the defect is in what it writes, not in what it validates.

## Addendum — 2026-07-28 — Routed to Spec 0056

All four suggested resolutions are owned by
[Spec 0056 — Profiles configure merge semantics](../specs/0056-profiles-configure-merge-semantics/_prd.md).
