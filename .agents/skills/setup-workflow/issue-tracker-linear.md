# Issue tracker: Linear

Issues for Gesttione repos live in **Linear** (workspace `gesttione`, team **Development**, key
`DEV`). There is no CLI — every operation goes through the **Linear MCP** tools. Configure the
Linear MCP once (`https://mcp.linear.app/mcp`, OAuth) before using these skills.

## Structure — where things go

- **Initiative = product/system.** Visio, Tax, Factory (in repo `vortex`); Conexus; Fluxus; Argus;
  Ligare-Cards; Infra/DevOps; and **Lab** (PoCs/experiments — isolated, not product health).
- **Project = deliverable/epic.** Each product also has a permanent **`Maintenance`** project for
  small BAU work. A repo can hold several products (Vortex holds Visio/Tax/Factory); the product is
  the initiative, not the repo.
- **Issue = story**, **Sub-issue = subtask** (subtasks are checklist items, not separately pointed).
- **Statuses:** the skill places issues in **Triage** (needs routing) or **Backlog** (already
  scoped); the active flow is `Backlog → Todo → In Progress → Done`, plus `Deferred`/`Duplicate`
  (closed). Any PR-driven status changes are handled automatically by the GitHub integration — the
  skill does not set them.
- **Labels:** **Type** (conventional-commit values: feat, fix, refactor, docs, test, chore, build,
  ci, perf, style, revert), **Scope** (backend, frontend, worker, infra, devops — issue-level, when
  one layer is clear), and **Triage** (the five roles below).
- **Estimation:** Fibonacci `1, 2, 3, 5, 8, 13`, on the issue (split anything ≥ 13).

## Conventions (Linear MCP)

- **Create an issue:** `save_issue` with the team (`DEV`), the target `project`, a `title`, a
  `description`, `labels` (a Type value, an optional Scope value, and a triage role), an `estimate`
  (Fibonacci), a `priority`, and a `state`.
- **Read an issue:** `get_issue` with the identifier (`DEV-123`) or its URL.
- **List issues:** `list_issues` filtered by `project`, `state`, `label`, or `assignee`.
- **Comment:** `save_comment` with the `issueId`.
- **Apply/remove labels or change status:** `save_issue` (update) with `labels` / `state`.
- **Block/relate:** create blocking relations between issues (use the relations field), or
  reference `DEV-###` in the body.

Infer nothing about repos from `git remote` — the issue tracker is Linear regardless of which repo
you are in. The repo↔product link is the GitHub integration (PRs reflect into Linear; Linear never
writes issues back to GitHub).

## When a skill says "publish to the issue tracker"

Create a Linear **issue** under the relevant **Project** (the deliverable). Give it a Type label, a
Fibonacci estimate, a priority, and the appropriate triage role. New issues land in **Triage** if
they still need routing, or **Backlog** if already scoped. Output from `to-prd`/`to-issues` is
AFK-ready, so apply **`ready-for-agent`**.

## When a skill says "fetch the relevant ticket"

`get_issue DEV-123` (or by URL). For a whole feature, `list_issues` filtered by its `project`.

## PRD handling (`to-prd`)

`to-prd` **must publish the PRD to Linear** — Linear is the source of truth, not the repo. The PRD
maps to a **Project**, not an issue:

1. Draft the PRD. The skill may use `.scratch/<feature-slug>/PRD.md` as a working file, but
   `.scratch` is gitignored — it is transient, not the record.
2. Create or update the Linear **Project** for the deliverable, under the right product initiative,
   via `save_project`.
3. Put the **full PRD content** in the Project's **overview** — the `description` field of
   `save_project` is the overview markdown (`summary` is the short one-liner). Do **not** leave the
   PRD only as a local file or paste a link to one; copy the whole PRD in so the Project is
   self-contained.
4. Start the Project at state `Backlog` (or `Planned`); move to `In Progress` when work starts.

The deliverable is the Project; the slices are its issues. Once the overview holds the full PRD, the
`.scratch` working copy can be deleted.

## Issue decomposition (`to-issues`)

`to-issues` **must create each vertical slice as a Linear issue** (`save_issue`) under the PRD's
**Project** — not as local files:

- Title in the project's domain vocabulary; the **full slice body goes in the issue `description`**
  (the `to-issues` template: What to build / Acceptance criteria / Blocked by). Don't link to a
  local file — copy the content in.
- A **Type** label (`feat`, `fix`, …); a **Scope** label when the slice is clearly one layer; a
  Fibonacci **estimate**; a **priority**.
- **"Blocked by"** → use Linear issue **relations** (blocked-by) and/or reference `DEV-###` in the
  body. Create issues in dependency order so blocker identifiers already exist.
- Apply **`ready-for-agent`** unless told otherwise; already-scoped slices start in **Backlog**.

Sub-issues are checklist subtasks and are not pointed.

## Pull requests (linking back to the tracker)

PRs are created during implementation (not by `to-prd`/`to-issues`), but they must link back to the
issue so the status automation fires. The issue must already exist.

- **Link the PR to the issue** by either:
  - naming the branch with Linear's suggested name — `<user>/dev-123-slug`, via "Copy git branch
    name" on the card — so the PR auto-attaches to `DEV-123`; **or**
  - putting `DEV-123` (or `Fixes DEV-123` / `Closes DEV-123`) in the PR title or description.
- Once attached, the team's GitHub automations advance the **issue** status automatically; you don't
  set those by hand. The **project's** progress/health rolls up from its issues.
- **Requires** the GitHub integration connected. **Issue Sync stays OFF** — PR linking is a separate
  mechanism that keeps working with it off (Linear never writes issues back to GitHub).
- Use Conventional Commit style for the PR title; the issue's **Type** label (`feat`/`fix`/…) should
  match the commit type.

## Triage roles → Linear

The five canonical roles exist as a **Triage** label group (names map 1:1), with two native
fallbacks worth knowing:

| Role (mattpocock/skills) | In Linear |
| --- | --- |
| `needs-triage` | label `needs-triage` — or simply an issue sitting in the **Triage** status |
| `needs-info` | label `needs-info` |
| `ready-for-agent` | label `ready-for-agent` |
| `ready-for-human` | label `ready-for-human` |
| `wontfix` | label `wontfix` — also move the issue to the **Deferred** status |

(See also `triage-labels.md`; the right-hand column there already matches these names.)
