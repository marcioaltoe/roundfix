---
name: conventional-commits
description: Create Conventional Commits and PR titles that pass Cocogitto validation. Use when the user asks to commit, stage changes, split commits, write a commit message, open or update a PR title, or prepare squash-merge history.
metadata:
  category: git
  tags: [git, github, workflow, communication, documentation]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Conventional Commits

Create commits and PR titles that are easy to review and pass `cog verify`.

## Required format

Use Conventional Commits:

```text
type(scope): imperative subject
```

The scope is optional:

```text
fix: handle empty setup files
feat(setups): add rust CLI preset
docs(readme): clarify setup install commands
```

## Allowed types

- `feat` - user-facing feature or capability
- `fix` - bug fix
- `refactor` - code change without behavior change
- `docs` - documentation only
- `test` - tests only
- `chore` - maintenance without runtime behavior change
- `build` - build tooling, dependencies, or packaging
- `ci` - CI workflow changes
- `perf` - performance improvement
- `style` - formatting only, not visual UI changes
- `revert` - revert a previous change

## Subject rules

- Use imperative mood: `add`, `fix`, `remove`, `rename`; not `added`, `fixed`, `removes`.
- Start lowercase after the colon.
- Do not end with a period.
- Keep the first line under 72 characters when possible.
- Describe the intent, not the file names.

## Commit workflow

1. Inspect the working tree:
   - `git status --short`
   - `git diff --stat`
   - `git diff`

2. Decide commit boundaries.
   - Split unrelated changes into separate commits.
   - Split behavior, tests, docs, formatting, dependencies, and generated files when they are reviewable separately.
   - If one file has mixed changes, use patch staging.

3. Stage only intended changes.
   - Use `git add <path>` for clean file boundaries.
   - Use `git add -p` for mixed hunks.
   - Review staged content with `git diff --cached`.

4. Write the message.
   - Choose the type from the allowed list.
   - Add a scope when it helps review, usually the package, module, collection, or feature.
   - Add a body only when the reason is not obvious.
   - Use footers for `BREAKING CHANGE:`, `Refs #123`, or `Closes #123`.

5. Validate when Cocogitto is available.
   - `cog verify "$(git log -1 --pretty=%B)"`
   - For a proposed title: `cog verify "$PR_TITLE"`

6. Run the smallest relevant verification before committing, or use the repo's required verification command when one exists.

## PR title rules

Pull request titles must also use Conventional Commits because squash merge often uses the PR title as the final commit message.

Before opening or editing a PR:

```bash
PR_TITLE='feat(setups): add Go CLI preset'
cog verify "$PR_TITLE"
```

If `cog` is not installed, still write the title in the exact Conventional Commit format and expect CI to validate it.

## Examples

Good:

```text
feat(auth): add token refresh flow
```

Good with body:

```text
fix(cli): handle empty setup files

The installer previously treated an empty setup file as a successful install.
This now fails early with a clear error.
```

Bad:

```text
update files
```

Bad:

```text
fix: change setup line 42
```

## Output

When asked to commit, report:

- commit boundaries chosen
- staged files for each commit
- final commit message
- verification command and result
