// Suite: recorded Delivery Queue item worktrees.
// Invariant: readers list every recorded path in queue and item order without writing.
// Boundary IN: the public SQLite-backed store API, including OpenReader.
// Boundary OUT: Git worktree creation and command refusal formatting.
package store

import (
	"reflect"
	"testing"
)

func TestDeliveryQueueItemWorktreesListsEveryRecordedWorktree(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	home := t.TempDir()
	writer := openTestStore(t, ctx, home)
	// Insert queues and worktrees out of order; slug and path order differ from position.
	for _, root := range []string{"/repos/z", "/repos/a"} {
		if _, err := writer.CreateDeliveryQueue(ctx, root, []string{"second-alphabetically", "first-alphabetically", "no-worktree"}); err != nil {
			t.Fatal(err)
		}
		for _, slug := range []string{"first-alphabetically", "second-alphabetically"} {
			if _, _, _, err := writer.RecordDeliveryQueueItemWorktree(ctx, root, slug, "test/"+slug, root+"/"+slug); err != nil {
				t.Fatal(err)
			}
		}
		if err := writer.ClaimDeliveryQueueOwner(ctx, root, 4242, "recorded-owner"); err != nil {
			t.Fatal(err)
		}
	}
	before := make([]DeliveryQueue, 0, 2)
	for _, root := range []string{"/repos/a", "/repos/z"} {
		queue, found, err := writer.DeliveryQueue(ctx, root)
		if err != nil || !found {
			t.Fatalf("read queue: found=%t err=%v", found, err)
		}
		before = append(before, queue)
	}
	closeStore(t, writer)
	reader, err := OpenReader(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	defer closeStore(t, reader)
	got, err := reader.DeliveryQueueItemWorktrees(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := []DeliveryQueueItemWorktree{
		{GitRoot: "/repos/a", SpecSlug: "second-alphabetically", Worktree: "/repos/a/second-alphabetically"},
		{GitRoot: "/repos/a", SpecSlug: "first-alphabetically", Worktree: "/repos/a/first-alphabetically"},
		{GitRoot: "/repos/z", SpecSlug: "second-alphabetically", Worktree: "/repos/z/second-alphabetically"},
		{GitRoot: "/repos/z", SpecSlug: "first-alphabetically", Worktree: "/repos/z/first-alphabetically"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("worktrees=%+v, want %+v", got, want)
	}
	for _, queue := range before {
		after, found, err := reader.DeliveryQueue(ctx, queue.GitRoot)
		if err != nil || !found || !reflect.DeepEqual(after, queue) {
			t.Fatalf("query changed queue: found=%t err=%v got=%+v want=%+v", found, err, after, queue)
		}
	}
}
