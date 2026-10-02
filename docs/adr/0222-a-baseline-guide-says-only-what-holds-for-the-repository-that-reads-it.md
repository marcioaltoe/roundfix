---
status: accepted
created_at: 2026-10-02T00:00:00Z
updated_at: 2026-10-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Baseline guide says only what holds for the repository that reads it

After the stack wave, four sentences the Baseline renders still did not hold
for every repository that reads them. The backend and frontend guides named
"the workspace" without its path, although the built-in profiles declare
`packages/backend` and `packages/frontend`. The clause "Do not introduce
generic `modules` or `services` buckets as the normative backend architecture"
had no scope and read as a ban on domain services. Four modules dispatched nine
triggers to skills whose metadata turns off model invocation, so an Agent was
told to activate a skill it cannot load. And the production-code and debugging
Skill Activations were owned by the TypeScript module but worded for every
language, which left their owner looking accidental.

Four rules settle them:

- A built-in profile binds each declared workspace to the guide that governs
  it, and that guide renders the path in its scope sentence. A guide with no
  bound workspace, which includes every repository-owned profile, keeps its
  generic sentence byte for byte.
- A clause whose wording changes takes a new identity that declares
  `replaces` for the old one, so an adopter's update records the old clause
  `replaced` instead of letting retained text drift under a stable identity.
- The skill guide tells the Agent to ask the person to run a matching skill
  that only a person can start, instead of activating it or carrying out its
  workflow itself. One core clause covers every such skill, present or future;
  no per-trigger marker is kept, because the skill's own metadata already
  says who may start it.
- A Skill Activation belongs to a module that only profiles whose setups
  install every skill of its bundle select. The catalog already refuses any
  other owner (`catalog.profile.skill.dispatch-outside-setup`), and the
  production-code and debugging bundles hold skills that only the TypeScript
  setups install, so `typescript` owns them. Their wording stays
  language-neutral, as ADR-0190 asks of a bundle that holds no skill for one
  technology: in a repository that composes Go with TypeScript they apply to
  both languages.

Two alternatives were rejected. Removing the person-only triggers would drop
the pointer that tells the Agent when the person should run the skill. Moving
the two activations to `core` would dispatch skills the Go and Rust setups do
not install.

## Consequences

Standard TypeScript Monorepo and composed-profile adopters read their
workspace paths and the scoped bucket clause on their next update, and the
Source Baseline gains a row for each new clause. Every adopter's skill guide
gains one mandatory clause. A repository-owned profile cannot yet bind a
workspace, so its guides keep the generic scope sentence.
