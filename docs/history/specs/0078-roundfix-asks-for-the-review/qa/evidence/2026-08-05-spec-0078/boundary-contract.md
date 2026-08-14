# Configuration and Review Source boundary

Fresh command `rtk go test ./internal/config -count=1 -run
'Test(BuiltinReviewRequestDefaults|LoadAppliesReviewRequestHierarchy|ValidateRejectsEmptyReviewRequestCommand|LoadRejectsNonBooleanReviewRequest)$'
-v` exited 0 with 7 tests passing.

Fresh request-boundary command `rtk go test
./internal/reviewsource/coderabbit -count=1 -run '^TestClientRequestReview' -v`
exited 0 with 11 tests passing.

Together they observed:

- built-in asking `false` and command `@coderabbitai review`;
- User Config then Project Config precedence;
- typed rejection of non-boolean asking and rejection of a blank command;
- exact same-head maintainer/Roundfix marker deduplication;
- different-head, command-only, embedded-marker, and Review Source-authored
  near misses do not suppress a request;
- page-two marker discovery through the real paginated `GHClient` command;
- list failure returns after one attempt with zero post/event;
- publish failure returns its wrapped cause and emits no success event.

No Evidence is read or returned by the request interface.
