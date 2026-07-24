<!-- source-baseline-entry: clause.backend.boundary-contracts -->
- MUST keep blocking, network, process, database, and daemon boundaries explicit about ownership, cancellation, timeouts, and error reporting. Test the lowest real boundary that proves the repository-authored contract; do not invent authentication, database, or transport policy.
<!-- /source-baseline-entry: clause.backend.boundary-contracts -->

<!-- source-baseline-entry: clause.backend.layered-architecture -->
- MUST organize backend behavior through domain, application, and infrastructure layers. Dependencies point inward toward domain behavior.
<!-- /source-baseline-entry: clause.backend.layered-architecture -->

<!-- source-baseline-entry: clause.backend.thin-http-handlers -->
- MUST keep HTTP handlers thin: validate and translate transport input, invoke one application use case, and translate the result into the repository's HTTP Contract.
<!-- /source-baseline-entry: clause.backend.thin-http-handlers -->

<!-- source-baseline-entry: clause.backend.http-independent-use-cases -->
- MUST keep application use cases independent of HTTP request, response, router, and middleware types.
<!-- /source-baseline-entry: clause.backend.http-independent-use-cases -->

<!-- source-baseline-entry: clause.backend.persistence-owner -->
- MUST keep persistence implementation in infrastructure and behind application-owned boundaries; schema and query definitions belong to the selected persistence capability.
<!-- /source-baseline-entry: clause.backend.persistence-owner -->

<!-- source-baseline-entry: clause.backend.prohibit-generic-layers -->
- MUST NOT introduce generic `modules` or `services` buckets as the normative backend architecture.
<!-- /source-baseline-entry: clause.backend.prohibit-generic-layers -->
