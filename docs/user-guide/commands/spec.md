### spec judge

```bash
roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]
```

Raises advisory judgments about ADR attributions and goal-to-mechanism links
in one active Spec, resolved through the configured Spec Root. By default it
reads both the PRD and TechSpec. `--stage prd` judges PRD attributions;
`--stage techspec` judges TechSpec attributions and its Coverage Map against
the PRD's goals. It reads only the PRD, TechSpec, and accepted ADRs.

Text output lists `advisory` and `skipped` results followed by a cost and
transport summary. Clear judgments appear only in `--format json`, whose
schema is `roundfix/spec-judge/v1`. An advisory asks the author to correct the
artifact or explain why the text stands. A skipped result proves neither a
failure nor a clean result.

The command reads `ROUNDFIX_OPENROUTER_API_KEY` for OpenRouter first, then
`ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly when the OpenRouter key is
absent. It does not read the generic `OPENROUTER_API_KEY`, so Roundfix's Jev
cost stays on its own key. Each response's versioned model ID is recorded;
thresholds apply only to the pinned Jev 1.13 model family.

Every request appends its answer, model, transport, and cost to the Judge Log
at `<home>/.roundfix/judge/<YYYY-MM>.jsonl`, using the UTC month. The monthly
ceiling is US$5.00 across both transports and every repository using that
Roundfix Home. The log contains no API key.

The whole run skips without a key, with an unreadable Judge Log, or when the
monthly ceiling is already reached. Non-English artifacts skip. Individual
judgments skip when their source does not qualify, the answering model does
not match the pin, the answer is unreadable, or another client error refuses
the request. A refused key (HTTP 401, 402, or 403), exhausted service retries,
a server failure, timeout, network error, reached ceiling, or unreadable or
unwritable log stops further requests. Judgments not asked count as skipped.

The command exits `0` whenever it ran, including advisory, skipped, and
stopped results. It never gates authoring or changes another command's exit
code. Exit `2` means invalid arguments, an unknown active Spec, a missing
PRD, or a missing TechSpec with `--stage techspec`.
