---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A stack rule carries force, and a preference names what licenses the exception

The Go, Rust, CLI and TUI modules shipped their rules as unlabelled paragraphs,
while every other module ships Normative Clauses with a stated enforcement
level. An Agent read "prefer the standard library" and "use stdlib `testing`"
as advice in one repository and as a ban in another, and Roundfix's own
repository had to restate its Go rules by hand. Every rule of those four
modules is now a clause that carries `mandatory`, `prohibited` or
`stop-and-ask`, and each clause states one obligation.

A preference is written as an obligation about its exception. "Stdlib first"
is mandatory as a preference: a third-party Go module is allowed when the
change that adds it records its reason in an ADR or another recorded
repository decision. The shipped Go module prohibits no library by name. A
repository that forbids one, as Roundfix forbids Cobra and testify, says so in
its own Repository-Specific Normative Rules, and that rule and the Baseline
rule both govern over a dispatched skill's default. The Rust module states a
policy, not a preference: library and domain code return typed errors, a
type-erased error such as `anyhow` is allowed only at a binary's entry point,
and `unwrap`, `expect` and `panic!` are prohibited on any path user input, the
environment, a file or an IO result can reach.

Two alternatives were rejected. A prohibition of every third-party Go module
would bind adopters to one repository's choice and turn each legitimate
dependency into a rule violation. A clause that forbids named libraries would
contradict the vendored Go skills without giving an adopter a way to record a
reasoned exception.

## Consequences

An adopter of the Go CLI/TUI or Rust CLI profile reads its stack rules with
force after its next Baseline update. A Go adopter that adds a module now
records why. A Rust adopter whose library code returns `anyhow::Result` has a
rule to meet. The vendored skills that recommend Cobra, testify or `anyhow`
are not edited; the rule governs where they disagree.
