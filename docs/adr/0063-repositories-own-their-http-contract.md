# ADR-0063: Repositories own their HTTP contract

Status: Accepted

The Standard TypeScript Monorepo Profile requires each repository to persist an HTTP Contract Decision instead of imposing universal REST or POST-only semantics. Setup reuses an existing supported contract or asks for the mode and its typed exceptions, because framework capability cannot safely infer application API policy.
