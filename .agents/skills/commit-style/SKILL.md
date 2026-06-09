---
name: commit-style
description: Creates Conventional Commits messages. Use when the user asks for "commit", "git commit", or after changes are ready to be recorded in git history.
metadata:
  category: git
  tags: [git, commits, conventional-commits]
  version: 0.1.0
  author: Marcio Altoé
---

# Commit Style

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/). Focus on **why** the change exists, not **what** changed; the diff already shows what changed.

## When To Use

- The user asks for "commit", "git commit", or "commit this".
- A sequence of changes is ready to be recorded.
- Git history needs a new descriptive entry.

## How To Apply

1. **Identify the type** based on the nature of the change:
   - `feat` - new functionality
   - `fix` - bug fix
   - `refactor` - code change without behavior change
   - `docs` - documentation only
   - `test` - tests only
   - `chore` - maintenance, dependencies, or config without runtime impact
   - `perf` - performance improvement
   - `style` - formatting only, not visual UI changes

2. **Choose the scope** when useful. It is the affected area: `feat(auth):`, `fix(api):`, `refactor(cli):`. Use the primary folder or module name changed.

3. **Write the subject line**:
   - English only for this repository.
   - Imperative mood: "add", "fix", "remove"; not "added" or "fixed".
   - No trailing period.
   - Maximum 72 characters.

4. **Add a body when the reason is not obvious**:
   - Blank line after the subject.
   - Focus on **why**, not **what**.
   - Wrap lines at 72 characters.

5. **Add a footer when needed**:
   - `BREAKING CHANGE: <description>` for breaking compatibility changes.
   - `Closes #123` or `Refs #123` for related issues.

## Examples

Good:

```text
feat(auth): add token refresh flow

Sessions were expiring mid-task and forcing users to restart.
The refresh runs silently 30s before expiry.

Closes #142
```

Good short form:

```text
fix(cli): handle empty skills directory without crashing
```

Bad:

```text
update files
```

Bad because it describes the diff instead of the reason:

```text
fix: change line 42 of auth.ts to use Date.now()
```

## Anti-Patterns

- Generic messages: "update", "fix stuff", "wip".
- Mixing multiple change types in one commit. If you added a feature and fixed a bug, use two commits.
- Past-tense subjects: "added", "fixed". Use imperative mood.
- Bodies that explain the diff line by line; the diff is already available.
- Using `feat` for changes that are not features. Renaming a file is `refactor`, not `feat`.
