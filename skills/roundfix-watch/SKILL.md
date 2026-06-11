---
name: roundfix-watch
description: Start and observe Roundfix review cleanup runs for CodeRabbit pull request feedback.
---

# Roundfix Watch

Use this skill when the user asks to resolve CodeRabbit comments, watch this
pull request, run Roundfix until clean, or clean up review bot feedback.

## Required Behavior

1. Prefer `roundfix` commands over manual GitHub scraping.
2. Inspect the current repository and Open Pull Request only when Roundfix needs
   missing command input.
3. Start the watched loop with:

   ```bash
   roundfix watch --source coderabbit --pr <number> --agent <agent> --until-clean
   ```

4. Let Roundfix own Review Source waits, CodeRabbit fetches, Round creation,
   Agent lifecycle, verification, Batch commits, Final Push, Review Source
   resolution, retries, timeouts, and Stop Request handling.
5. Report the Run ID, Open Pull Request, Review Source, Agent, and current Run state whenever you summarize progress.
6. Prefer the Roundfix Live Run View or daemon output for long waits.

## Do Not

- Do not manually scrape GitHub review comments when `roundfix fetch` or
  `roundfix watch` is available.
- Do not manually resolve CodeRabbit threads unless Roundfix is unavailable and
  the user explicitly asks for a manual fallback.
- Do not commit, push, or resolve Review Source threads outside Roundfix for a
  Run that Roundfix owns.

## Useful Commands

```bash
roundfix fetch --source coderabbit --pr <number>
roundfix resolve --pr <number> --agent <agent>
roundfix watch --source coderabbit --pr <number> --agent <agent> --until-clean
roundfix skills check
```
