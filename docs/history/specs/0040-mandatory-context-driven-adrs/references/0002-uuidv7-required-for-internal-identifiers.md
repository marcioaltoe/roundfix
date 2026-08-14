# UUIDv7 required for Internal Identifiers

## Status

Accepted

## Context

Mixing UUID versions, shape-only validation, sentinel values, and
library-specific defaults makes identity contracts depend on which component
created a record. A producer can then persist an identifier that a consumer
rejects even though both systems describe the value as a UUID.

A portable Context-Driven repository needs one invariant for technical
identities it creates and controls while preserving identities whose format is
owned by a business rule, protocol, or external system.

## Decision

Every Internal Identifier created and controlled by the project must be a valid
UUIDv7 as defined by RFC 9562. The rule applies to primary entity or resource
identities, references to those identities, records created through internal
adapters or libraries, deterministic seeds, and fixtures that represent
project-owned entities or resources.

The project maintains one canonical UUIDv7 generation and validation contract.
Language-specific implementations may adapt that contract, but every generator
must produce UUIDv7 and every validator must verify textual format, version 7,
and the RFC variant. Validation that checks only the generic UUID shape does
not satisfy the contract.

Persistence defaults, frameworks, authentication systems, and other libraries
that create project-owned records must use the canonical generation contract
or receive an already generated UUIDv7. Deterministic seed and fixture values
must also be valid UUIDv7 values.

Identifiers originating outside the project preserve the format defined by
their source. Natural keys, business codes, regulatory identifiers,
protocol-defined identities, and values that do not represent persistent
identity are not Internal Identifiers and are not converted to UUIDv7.

When existing data violates this invariant, the project must define an
explicit migration and rollback plan appropriate to its production state.
Adopting this decision never implies that destructive reset or reseeding is
safe.

## Alternatives Considered

### Accept any UUID version

This reduces migration work but preserves version-dependent behavior and
prevents strict contracts between producers and consumers.

### Require UUIDv7 only in new components

This limits immediate changes but leaves identity validation dependent on the
component that created each record.

### Allow permanent sentinel exceptions

This keeps familiar bootstrap values but forces every contract to carry an
exception that deterministic valid UUIDv7 fixtures can avoid.

## Consequences

- Domain, API, persistence, integration, seed, and fixture contracts can
  validate the same identifier invariant.
- Project-controlled libraries and adapters cannot silently fall back to
  UUIDv4 or shape-only UUID validation.
- Existing schemas, defaults, data, fixtures, and consumers may require an
  explicit migration before they conform.
- UUIDv7 provides approximate temporal ordering and exposes generation-time
  information to recipients of the identifier.
