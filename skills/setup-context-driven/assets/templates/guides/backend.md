# Backend

Keep IO boundaries explicit. Blocking, network, process, database, and daemon
operations need clear ownership, cancellation, and error reporting.

Test service contracts at the lowest layer that can catch the defect. Use real
filesystem or process boundaries when they are the product contract; avoid
network dependencies in portable setup tests.
