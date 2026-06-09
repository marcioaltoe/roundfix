---
description: Commit and push all changes in BOTH the main repo and the .knowledge repo (which backs the symlinked docs/, .scratch/, .compozy/, CONTEXT.md), using the commit-style skill for Conventional Commit messages.
argument-hint: "[optional hint to bias the commit messages, e.g. 'production readiness CD split']"
---

**Input:** `$ARGUMENTS` = optional free-text hint to bias the commit messages (scope, theme, ticket).
Empty is fine — derive the message from the diff.

## Why this command has two repos

This repo uses an external knowledge workspace. `docs/`, `.scratch/`, `.compozy/`, and `CONTEXT.md`
at the repo root are **symlinks** into `.knowledge/projects/<project>/…`, and `.knowledge/` is a
**separate git repository** (`gesttione-solutions/knowledge.git`, branch `main`) — see
`scripts/knowledge-bootstrap.sh`. Editing anything under those paths changes the **knowledge** repo,
not the main repo. So a complete "commit and push everything" must operate on **both** repos:

- **Main repo** (`gesttione-solutions/fluxus.git`) — code, workflows, config. On a feature branch.
- **Knowledge repo** (`.knowledge/`) — docs, PRDs, `.scratch/` issues, `CONTEXT.md`. On `main`.

The root symlinks are tracked-and-unchanged in the main repo, so `git add -A` there stages only real
main-repo changes (it never pulls knowledge content into the main repo).

## 1. Survey both repos first

Run and read the output before staging anything:

```bash
git status --short                 # main repo
git -C .knowledge status --short   # knowledge repo (docs/.scratch/CONTEXT.md changes)
```

Decide which repos actually have changes. Skip any repo that is clean. If **both** are clean, report
"nothing to commit" and stop.

Never stage or commit secrets, tokens, private URLs, `.env*` files, SSH keys, or raw logs. If the diff
contains anything secret-like, stop and ask the user.

## 2. Commit + push the knowledge repo (docs / .scratch / CONTEXT.md)

Only if `git -C .knowledge status --short` shows changes. The knowledge repo is **shared** and lives on
`main`, so rebase before pushing to avoid non-fast-forward rejects. Never force-push it.

1. Sync to avoid rejects:
   ```bash
   git -C .knowledge fetch origin main
   git -C .knowledge pull --rebase --autostash origin main
   ```
2. Stage everything:
   ```bash
   git -C .knowledge add -A
   ```
3. Generate the message with the **`commit-style`** skill from the staged knowledge diff
   (`git -C .knowledge diff --staged`). These are docs/tracker changes, so the type is usually
   `docs:` (or `chore:` for tracker bookkeeping). Weave in `$ARGUMENTS` if given. Then commit:
   ```bash
   git -C .knowledge commit -m "<commit-style message>"
   ```
4. Push:
   ```bash
   git -C .knowledge push origin main
   ```
   If the push is rejected, re-run the fetch + `pull --rebase --autostash`, resolve any conflict, and
   push again. Never `--force`.

## 3. Commit + push the main repo (code / workflows)

Only if `git status --short` shows changes.

1. **Quality gate (blocking).** If the main-repo changes touch anything beyond pure docs (code,
   workflows, config — which is the normal case), run the project gate and read the real output:
   ```bash
   rtk make verify
   ```
   It must end with **zero errors and zero warnings**. If it fails, stop, report the failure, and do
   not commit. (A main-repo change set that is exclusively markdown/docs may skip this.)
2. Stage everything (this also records file deletions/renames, e.g. removed workflows):
   ```bash
   git add -A
   ```
3. Generate the message with the **`commit-style`** skill from the staged main diff
   (`git diff --staged`). Pick the correct Conventional Commit type from the actual change
   (`feat`/`fix`/`ci`/`chore`/`refactor`/…). Weave in `$ARGUMENTS` if given. Then commit:
   ```bash
   git commit -m "<commit-style message>"
   ```
4. Push the **current branch** (never assume `main`; set upstream if missing):
   ```bash
   git push -u origin HEAD
   ```
   Never push to a protected branch directly and never `--force`.

## 4. Report

Summarize concisely (pt-BR for domain content):

- For each repo touched: the commit type/subject, the short SHA, and the branch it was pushed to.
- The `rtk make verify` result for the main repo (with observed numbers), or that it was skipped
  (docs-only / no main changes).
- Any repo skipped because it was clean.
- If a push needed a rebase, note it. If anything was left uncommitted on purpose (e.g. a suspected
  secret), call it out explicitly.
