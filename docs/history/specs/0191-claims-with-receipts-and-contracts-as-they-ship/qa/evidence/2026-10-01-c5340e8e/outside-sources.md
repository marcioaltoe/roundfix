# Published evidence observed 2026-10-01

Shell retrieval with `rtk curl-cffi get <URL> --impersonate chrome --timeout 30 --body` failed for both sources: CONNECT tunnel response 403. The web tool successfully opened both exact URLs, so this transport restriction does not block the source row. No model or paid API was called.

## Claude documentation

Requested: https://docs.anthropic.com/en/docs/test-and-evaluate/strengthen-guardrails/reduce-hallucinations
Resolved: https://platform.claude.com/docs/en/test-and-evaluate/strengthen-guardrails/reduce-hallucinations
Section: Basic hallucination minimization strategies, Verify with citations; web-open lines 80 and 91.

> For each claim, find a direct quote from the documents that supports it.

The same guidance requires withdrawal when supporting evidence cannot be found: “it must retract the claim.” This supports evidence-first grounding. This checker proves quote presence; it does not establish semantic support.

## Cram documentation

https://pypi.org/project/cram/ — Project description; web-open lines 35 and 71–89.

> Cram tests look like snippets of interactive shell sessions.

The format describes commands following a dollar sign, subsequent expected output, and comparison of expected versus actual output. The page does not explicitly document bracketed nonzero status syntax; that detail is not credited by this row. This independently supplies prior art for recording and comparing a CLI session. Roundfix separately records stdout, stderr and exit.
