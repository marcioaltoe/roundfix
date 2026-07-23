# ADR-0068: Baseline Command uses one confirmation-gated workflow

Status: Accepted

The public `roundfix baseline` command owns one preflight-to-verification
workflow for initial adoption and later updates. Human callers use the root
command and review one consolidated Change Plan through linear prompts.
Non-interactive callers use `roundfix baseline plan` to emit a portable
`roundfix/baseline-plan/v1` document and `roundfix baseline apply` to apply
only that document and its exact approved digest; results use
`roundfix/baseline-result/v1`. Separate human audit and apply journeys were
rejected because they expose engine phases instead of the maintainer's task,
while explicit automation stages make approval transportable and auditable.
