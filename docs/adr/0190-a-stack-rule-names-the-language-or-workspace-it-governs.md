---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A stack rule names the language or workspace it governs

The TypeScript, Bun, backend and frontend modules were written for a repository
with one language. A repository-owned Baseline Profile can compose them with
the Go or Rust modules. In that repository "use Bun-owned commands", "do not
substitute another package manager" and "the repository's declared package
manager" bind the other language's toolchain, and the trigger "Writing or
changing tests." dispatches a Vitest skill for a Go test. An audit on
2026-09-30 also found three adopters that had each written the same correction
by hand: tests run through the package script, never through the bare Bun
runner.

A stack rule now names what it governs, in one of two ways. A clause whose
literal reading would bind another language's toolchain carries its scope in
its own sentence. Every stack guide opens with one sentence that names the
language or workspace its rules govern. A Skill Activation whose bundle holds a
skill for one technology names that technology in its trigger. When a
dispatched skill's default conflicts with a Baseline rule or a
Repository-Specific Normative Rule, the rule governs, and the Baseline says so
once, in the clause that requires the skills.

Two alternatives were rejected. Rewording every backend and frontend clause
would change text the Source Baseline retains by its bytes, so an upgrade could
no longer prove those clauses survived. Rendering each workspace path into the
guide needs a profile field and a render token that do not exist, which is new
behavior for a later wave.

## Consequences

A repository with one language reads the same obligations as before, with one
added sentence per stack guide. A repository that composes stacks can adopt the
TypeScript modules without a rule that contradicts its Go or Rust toolchain.
The scope is stated in words: no check reads a workspace path, so a rule cannot
yet be bound to a directory.
