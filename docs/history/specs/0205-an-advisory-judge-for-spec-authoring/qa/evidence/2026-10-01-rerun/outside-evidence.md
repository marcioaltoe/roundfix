# Outside evidence read on 2026-10-01

Public documents only; no inference endpoint or authenticated service used. All eleven named pages were freshly reached through the web read tool during this rerun, plus the TypeSafe llms.txt index. The first curl-cffi attempt for the index was denied by the sandbox domain allowlist; web read reached the index and all pages. This resolved fetch constraint leaves no source blocked. The TypeSafe .md references were reached through their equivalent normal pages. No real key or inference endpoint was accessed.

## OpenRouter TypeSafe SDK guide

Source: [OpenRouter TypeSafe SDK guide](https://openrouter.ai/docs/guides/community/typesafe-sdk).

System One POST /api/v1/systemone with bearer authentication; model/state/questions requests and model/answers/usage responses, plus id/provider/usage.cost. Bare jev-1.13 maps into typesafe/jev-1.13; example returns the dated snapshot. Matches the shipped endpoint and metadata contract.

## OpenRouter Jev hub

Source: [OpenRouter Jev hub](https://openrouter.ai/docs/guides/community/jev).

Typed Choice/Noul/Score answers; OpenRouter billing, TypeSafe provider, 32,000-token state context, free output. Matches the transport and price basis. Hub now also promotes a Decisions API; this does not remove the SDK guide’s System One endpoint.

## OpenRouter Jev tutorial

Source: [OpenRouter Jev tutorial](https://openrouter.ai/docs/guides/community/jev-tutorial).

The response model names the dated snapshot serving the request; Choice confidence and Noul probability have distinct meanings. Matches accepted dated model shape. Current tutorial examples use /api/alpha/decisions, while the linked SDK path retains System One. No claim of live compatibility was inferred from tutorial examples.

## OpenRouter Jev 1.13 model page

Source: [OpenRouter Jev 1.13 model page](https://openrouter.ai/typesafe/jev-1.13).

Input US$0.042/million, free output, 32K context, TypeSafe provider. Matches the embedded price and primary model family.

## OpenRouter errors reference

Source: [OpenRouter errors reference](https://openrouter.ai/docs/api/reference/errors-and-debugging).

Account/key without credits remains an error cause; 402 handling distinguishes insufficient credits from a temporary in-flight budget and may carry Retry-After. Matches the specified fail-open stop on credit failure. Redirects to /docs/api_reference/errors-and-debugging; new transient-budget advice does not change this Spec’s authored stop policy.

## TypeSafe API reference

Source: [TypeSafe API reference](https://docs.typesafe.ai/api).

POST https://api.typesafe.ai/v1/systemone with bearer authorization and JSON; required state/model/questions and model/answers/usage, Choice choice/probabilities/confidence, Noul noul, input/output token counts; 401/422/429/529 errors. Matches direct endpoint and request/response contract.

## TypeSafe models

Source: [TypeSafe models](https://docs.typesafe.ai/models).

jev-1.13.0 at US$0.042/million; aliases move and response reports the actual model; pin the version once thresholds are tuned. English is the strongest language. Matches model, price, pinning and English restriction.

## TypeSafe confidence

Source: [TypeSafe confidence](https://docs.typesafe.ai/confidence).

Choice/Score confidence derives from the answer probability distribution; Noul has no confidence. Matches separate confidence/probability thresholds.

## TypeSafe citation-check cookbook

Source: [TypeSafe citation-check cookbook](https://docs.typesafe.ai/cookbooks/citation_check).

One Choice evaluates source support, with a 0.8 confidence boundary and human confirmation below it. Matches adopted support question and operating point; cookbook’s own older model example is illustrative, not the Spec’s pinned benchmark model.

## Agent Skills specification

Source: [Agent Skills specification](https://agentskills.io/specification).

Skill names require lowercase letters/digits/hyphens, maximum 64 characters; underscores are not allowed. Supports the decision not to create the_judge as a skill.

## Claude Code subagent reference

Source: [Claude Code subagent reference](https://code.claude.com/docs/en/sub-agents).

Frontmatter name cannot contain colon, reserved for plugin scopes. The current frontmatter table reserves colon for plugin scopes; the CLI-definition guidance additionally says not to start a name with a hyphen. PRD’s “only excludes colon” gloss omits this CLI-specific advice. No subagent is created and the_judge has no leading hyphen, so this does not contradict the shipped product contract.

## Independent pre-Spec measurements

Source: `docs/history/findings/2026-09-30-jev-judgments-measured-against-the-spec-archive.md`, an adopted finding predating this Spec. 5,240 direct requests, 5.19M input tokens, US$0.218, zero failures/retries, 171 archived Specs and 146 accepted ADRs; citation AUROC 0.89 versus overlap 0.80, goal-mechanism 0.86 versus 0.77. Citation operating point: non-support at confidence ≥0.8; goal point: probability <0.3. These values agree with the independently byte-compared configuration; no new calibration is claimed.

Supervised outside-session evidence, as recorded under TechSpec → Measured outside evidence (not rerun by QA): on 2026-10-01 TypeSafe refused jev-1.13 with 400 and answered jev-1.13.0; OpenRouter refused jev-1.13.0 with 400 and answered jev-1.13 as typesafe/jev-1.13-20260917, provider TypeSafe, reported cost equal to input tokens × US$0.042/million. Both answered 885/885 benchmark states, no retries; citation AUROC .8900/.8899, goal .8583/.8644, US$0.0269 each, p50 335/366 ms, p95 1067/697 ms. Advisory agreement 97.5–98.2%, mean absolute probability change .020. The original 1770-line log is outside this repository and unavailable here; these are maintainer-recorded measurements, not a fresh QA reproduction. Published primary documents independently support the request fields, model mapping and price.

## Comparison boundaries

Endpoint, body and metadata comparisons come from the API and SDK references; price comparisons come from both model pages; version-pinning advice comes from TypeSafe models. The hub/tutorial confirm billing and dated answer IDs, the confidence/cookbook pages support threshold semantics, and the two naming references address the no-skill/no-agent choice. Pages that do not specify a field or price are not credited with validating it. The API reference continues to list 401/422/429/529. The errors page now discusses a transient 402 in-flight budget; the authored policy deliberately stops on every 402 and remains advisory/fail-open. The tutorial’s Decisions API examples do not replace the System One path confirmed by the SDK guide. Model pages continue to publish Jev 1.13 and US$0.042 per million input tokens, with free output.

Origin distinction: the pre-Spec finding’s measurement was read directly in this rerun. The endpoint probes and two-transport benchmark remain maintainer-recorded supervised historical facts; their external scratch logs are unavailable here, and no live probe or fresh calibration is claimed. Published pages independently confirm the contract and price.
