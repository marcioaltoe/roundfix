# Roundfix Skill contract

`rtk make skills-sync-check` exited 0 and its prerequisite Skill tests passed
4 tests. `rtk bin/roundfix skills check` exited 0 for every required owned
Skill. `rtk cmp -s .agents/skills/roundfix/SKILL.md
skills/roundfix/SKILL.md` exited 0.

Focused text inspection found the shipped contract in the canonical Skill:
both key names and defaults, one post-Final-Push request for enabled `watch`
and `resolve`, the `fetch` exemption, request-is-not-Evidence boundary, no
retry/backoff/capacity wait, and both refused coherence pairs with repair
actions.
