// Suite: Delivery Engine review capacity wiring.
// Invariant: the Queue Owner uses the operator's configured Review Capacity.
// Boundary IN: commandDeliveryEngineDependencies with a loaded configuration.
// Boundary OUT: Engine.Run, providers and external processes are not started.
package cli

import (
	"io"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestTheDeliveryEngineGetsTheReviewCapacity(t *testing.T) {
	t.Parallel()
	var loaded roundconfig.Loaded
	loaded.Config.Delivery.ReviewConcurrency = 3
	dependencies := commandDeliveryEngineDependencies(nil, loaded, io.Discard)
	if dependencies.ReviewCapacity != 3 {
		t.Fatalf("ReviewCapacity = %d, want 3", dependencies.ReviewCapacity)
	}
}
