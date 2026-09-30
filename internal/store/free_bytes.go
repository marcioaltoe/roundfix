package store

import (
	"context"
	"fmt"
)

// FreeBytes reads the bytes held by free pages in the Run Database without a
// transaction or a pragma write.
func (store *Store) FreeBytes(ctx context.Context) (int64, error) {
	pageSize, err := storagePragmaInt64(ctx, store.db, "page_size")
	if err != nil {
		return 0, err
	}
	freePages, err := storagePragmaInt64(ctx, store.db, "freelist_count")
	if err != nil {
		return 0, err
	}
	freeBytes, err := checkedMultiply(freePages, pageSize)
	if err != nil {
		return 0, fmt.Errorf("calculate Run Database free bytes: %w", err)
	}
	return freeBytes, nil
}
