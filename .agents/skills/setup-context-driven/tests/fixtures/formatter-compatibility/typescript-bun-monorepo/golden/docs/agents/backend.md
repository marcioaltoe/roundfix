<!-- setup-context-driven:begin id=guide.backend version=0.0.1 -->

# Backend

This setup-owned guide defines portable backend rules. Repository-authored
architecture and service contracts remain authoritative.

- MUST keep blocking, network, process, database, and daemon boundaries explicit about ownership, cancellation, timeouts, and error reporting. Test the lowest real boundary that proves the repository-authored contract; do not invent authentication, database, or transport policy.

- MUST organize backend behavior through domain, application, and infrastructure layers. Dependencies point inward toward domain behavior.

- MUST keep HTTP handlers thin: validate and translate transport input, invoke one application use case, and translate the result into the repository's HTTP Contract.

- MUST keep application use cases independent of HTTP request, response, router, and middleware types.

- MUST keep persistence implementation in infrastructure and behind application-owned boundaries; schema and query definitions belong to the selected persistence capability.

- MUST NOT introduce generic `modules` or `services` buckets as the normative backend architecture.

<!-- setup-context-driven:end id=guide.backend -->
