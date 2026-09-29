package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

const reviewDispositionLockPollInterval = 10 * time.Millisecond

func reserveReviewFindingDisposition(
	ctx context.Context,
	ledgerPath string,
	record reviewRecord,
	finding reviewFinding,
	disposition reviewFindingDisposition,
) (line []byte, resultErr error) {
	lockPath := ledgerPath + ".lock"
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open review disposition ledger lock: %w", err)
	}
	if err := lockReviewDispositionFile(ctx, lockFile); err != nil {
		return nil, errors.Join(
			fmt.Errorf("acquire review disposition ledger lock: %w", err),
			lockFile.Close(),
		)
	}
	defer func() {
		releaseErr := errors.Join(unlockReviewDispositionFile(lockFile), lockFile.Close())
		if releaseErr == nil {
			return
		}
		releaseErr = fmt.Errorf("release review disposition ledger lock: %w", releaseErr)
		if resultErr == nil {
			resultErr = releaseErr
		} else {
			resultErr = errors.Join(resultErr, releaseErr)
		}
	}()

	dispositions, err := readReviewFindingDispositions(ledgerPath)
	if err != nil {
		return nil, err
	}
	if reviewFindingAlreadyDisposed(dispositions, record, finding) {
		return nil, fmt.Errorf("finding %q already has a disposition", finding.ID)
	}
	return appendReviewFindingDisposition(ledgerPath, disposition)
}
