package cli

import (
	"context"
	"fmt"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
)

type commandFreeSpaceChecker struct {
	loaded roundconfig.Loaded
}

var _ delivery.FreeSpaceChecker = commandFreeSpaceChecker{}

func (checker commandFreeSpaceChecker) LowFreeSpace(ctx context.Context) []string {
	var low []string
	for _, location := range inspectFreeSpace(ctx, checker.loaded) {
		if location.err == nil && location.available < doctorDiskSpaceThreshold {
			low = append(low, fmt.Sprintf("%s %s has %d MiB free, below 10240 MiB", location.label, location.path, location.available/(1024*1024)))
		}
	}
	return low
}
