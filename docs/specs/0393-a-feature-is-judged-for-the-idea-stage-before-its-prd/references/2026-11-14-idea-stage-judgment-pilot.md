# Pilot of the Idea-stage judgment (recorded by the operator, 2026-11-14 ~07:10–07:20 BRT)

Five yes/no questions asked to `jev-latest` (it answered as `jev-1.14.0`) through the TypeSafe HTTP API, over a state of three named parts: a two-sentence description of the product, a description of the Idea stage, and the request text. Thresholds are the TypeSafe documentation's starting points, not tuned values: 0.7 and above reads as yes, 0.3 and below as no, and anything between goes to a person. The Idea stage was "indicated" when any of the first three questions read yes, "undecided" when none read yes and one fell between, and "not indicated" otherwise.

## The questions

| id | instructions | yes means | no means |
| --- | --- | --- | --- |
| `shape_open` | Is the right shape of the solution for `request` still open? | At least two materially different product designs could satisfy the request, and the request does not settle which one to build. | The request already names one concrete design, or only one reasonable design exists. |
| `unverified_assumption` | Does `request` rest on an assumption about its users, the business or the market that it gives no evidence for? | A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong. | The need is shown by evidence in the request itself, such as a measured defect, a reported failure or a decision the product owner already made. |
| `market_question` | Does deciding whether to build `request` depend on a market or competitive question? | The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it. | The decision depends only on the product's own behavior and its current users. |
| `new_capability` | Does `request` ask for a kind of capability that `product` does not have today? | It adds a new capability, surface or concept to the product. | It corrects or adjusts behavior the product already has. |
| `behavior_specified` | Does `request` already specify the behavior to build precisely enough to write its requirements without first asking whether to build it or in what shape? | The inputs, the outputs and the rules of the wanted behavior are stated, and what remains is detail. | What exactly should be built, for whom, or whether it is worth building is still an open question. |

## The answers

| case | language | shape_open | unverified_assumption | market_question | new_capability | behavior_specified | reading |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Control written by the operator: a defect report with a reproduction (a command prints a header twice) | English | 0.07 | 0.06 | 0.03 | 0.04 | 0.98 | not indicated |
| Control written by the operator: a loose product idea (a hosted dashboard with paid plans) | English | 0.79 | 0.97 | 0.97 | 0.97 | 0.04 | indicated |
| The Indexus request that became Spec 0392 (ai-memory message 01a1d68d) | Portuguese | 0.39 | 0.36 | 0.22 | 0.92 | 0.88 | undecided |
| The maintainer's request that became this Spec | Portuguese | 0.23 | 0.61 | 0.46 | 0.81 | 0.63 | undecided |

Input sizes were 876 to 2083 tokens and each call returned 98 output tokens.

## What followed

Both undecided requests went to the maintainer, who sent both through the Idea stage. For Spec 0392 the stage's research measured that Claude Code ignores a skill's `model` header when the agent calls the skill, and the shape of that Spec changed before its PRD. The judge did not predict that: it was undecided.

## Limits

Four cases, two of them synthetic; one run each; no labeled set; the two real requests are in Portuguese and the questions in English. Nothing here is a measured threshold.
