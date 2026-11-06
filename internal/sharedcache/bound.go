package sharedcache

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"time"
)

// BoundBudget bounds measurement and trimming for the Build Cache Bound.
const BoundBudget = 30 * time.Second

// BoundMain parses -dir, -max-mb and -idle, trims dir as MeasureContext and
// ApplyContext plan, and returns 2 for a usage error and 0 otherwise.
func BoundMain(args []string, stderr io.Writer, now func() time.Time) int {
	flags := flag.NewFlagSet("gocache-bound", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	dir := flags.String("dir", "", "absolute Go build cache directory")
	maxMB := flags.Int64("max-mb", 0, "positive cache bound in MiB")
	idle := flags.Duration("idle", MinimumIdle, "minimum idle time before removal")
	usage := func(err error) int {
		fmt.Fprintf(stderr, "gocache-bound: %v\n", err)
		return 2
	}
	if err := flags.Parse(args); err != nil {
		return usage(err)
	}
	if flags.NArg() != 0 {
		return usage(fmt.Errorf("unexpected argument %q", flags.Arg(0)))
	}
	if !filepath.IsAbs(*dir) {
		return usage(errors.New("-dir must be an absolute path"))
	}
	if *maxMB <= 0 || *maxMB > math.MaxInt64/(1<<20) {
		return usage(errors.New("-max-mb must be a positive MiB count whose bytes fit an int64"))
	}
	if *idle < MinimumIdle {
		return usage(fmt.Errorf("-idle must be at least %s", MinimumIdle))
	}

	ctx, cancel := context.WithTimeout(context.Background(), BoundBudget)
	defer cancel()
	warning := func(err error) {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Fprintf(stderr, "gocache-bound: warning: %s: stopped after its 30s budget; the next make continues\n", *dir)
		} else {
			fmt.Fprintf(stderr, "gocache-bound: warning: %v; the build continues\n", err)
		}
	}
	plan, err := MeasureContext(ctx, *dir, *maxMB*(1<<20), *idle, now())
	if err != nil {
		warning(err)
		return 0
	}
	removed, bytes, err := ApplyContext(ctx, plan, *idle, now)
	if removed > 0 {
		fmt.Fprintf(stderr, "gocache-bound: %s: removed %d idle entries (%d bytes)\n", *dir, removed, bytes)
	}
	if err != nil {
		warning(err)
		return 0
	}
	if remaining := plan.Bytes - bytes; remaining > plan.Limit {
		fmt.Fprintf(stderr, "gocache-bound: warning: %s holds %d bytes, above its %d-byte bound, and no other entry is idle for %s; the build continues\n", *dir, remaining, plan.Limit, *idle)
	}
	return 0
}
