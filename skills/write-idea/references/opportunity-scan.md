# Opportunity scan

The strategist lens for `write-idea` step 6. The job is not to add features — it is to check, once, whether the idea on the table is the highest-leverage move before it becomes a spec. Think like an owner: what would make users unable to live without this product?

## 1. Assess the original idea

- High-leverage or incremental? A one-time improvement or something that compounds (data effects, habit formation, network effects)?
- What is the ceiling on its impact if it works perfectly?

## 2. Generate up to three alternatives

Span the scales deliberately:

- **More ambitious** — what if we thought bigger?
- **Simpler** — what if we stripped it to its essence? (Simple ≠ low-value; it often ships and validates faster.)
- **Adjacent** — what related problem could we solve instead?

Walk the categories to unstick thinking: Speed (what takes too long?) · Automation (what's repetitive?) · Intelligence (what could be smarter?) · Integration (what else do users use?) · Collaboration (how do users work together?) · Personalization (how is everyone different?) · Visibility (what's hidden that shouldn't be?) · Confidence (what creates anxiety? undo, previews) · Delight · Access (who can't use this yet?).

Prompts that work: "What would make a user tell a friend about this?" · "What do power users do manually that we could make native?" · "What would a competitor need to build to beat us?" · "What sounds crazy but might work?"

## 3. Score and recommend

Each alternative gets: **What** (one line — "better UX" is not an idea; "one-click rescheduling from the notification" is), **Why powerful**, **Scale** (small/medium/massive effort), and the same six-criteria score used in Feature Assessment (Impact, Reach, Frequency, Differentiation, Defensibility, Feasibility → Must do/Strong/Maybe/Pass).

Close with a recommendation the user can react to:

```markdown
### Recommendation

**Proceed with:** Original / Alternative N / Hybrid
**Rationale:** <why this is the highest-leverage move, citing research and scores>
```

Then ask the user to pick (A: original / B: alternative N / C: hybrid / D: other). Their choice — not the recommendation — sets the draft's scope.

## Behavioral traits

Challenge "obvious" ideas; prefer compounding over one-time; demand specificity; question assumptions ("users want X" may be wrong — what do they _need_?); cite evidence from the research and codebase findings, not vibes.
