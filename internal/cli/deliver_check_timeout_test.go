// Suite: Delivery Queue check wait wiring.
// Boundary IN: loaded config and the production dependency builder.
// Boundary OUT: no forge, process or host configuration is accessed.
package cli

import (
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
)

func TestTheDeliveryEngineReceivesTheConfiguredCheckWait(t *testing.T) {
	t.Parallel()
	loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), HomeDir: t.TempDir(), GitRoot: t.TempDir()}
	dependencies := commandDeliveryEngineDependencies(nil, loaded)
	if dependencies.CheckTimeout != roundconfig.DefaultDeliveryCheckTimeout {
		t.Fatalf("default check timeout = %s", dependencies.CheckTimeout)
	}
	loaded.Config.Delivery.CheckTimeout = 45 * time.Minute
	loaded.Config.Runs.SlotWait = 90 * time.Minute
	dependencies = commandDeliveryEngineDependencies(nil, loaded)
	if dependencies.CheckTimeout != 45*time.Minute || dependencies.RunSlotWait != 90*time.Minute {
		t.Fatalf("check timeout = %s, slot wait = %s", dependencies.CheckTimeout, dependencies.RunSlotWait)
	}
}
