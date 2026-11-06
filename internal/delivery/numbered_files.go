package delivery

import (
	"cmp"
	"fmt"
	"path"
	"slices"
	"strconv"
	"strings"
)

const BlockerNumberedFileCollision = "numbered-file-collision"

type NumberedFileClaims struct {
	Owner string   // the earlier item's Spec slug
	Paths []string // the Numbered Files that item adds
}

type NumberedFileCollision struct {
	Path  string // a Numbered File the item adds
	Other string // the colliding file
	Owner string // the earlier item's Spec slug, or "" for the default branch
}

type NumberedFileReport struct {
	Collisions []NumberedFileCollision
	NextFree   []string // "<directory>/<digits>", one per colliding directory
}

// FindNumberedFileCollisions compares paths already filtered by the caller's
// declaration. It neither modifies the claims nor reads configuration or Git.
func FindNumberedFileCollisions(defaultFiles, subject []string, ahead []NumberedFileClaims) NumberedFileReport {
	type numberKey struct {
		directory string
		number    uint64
	}
	type claim struct {
		path  string
		owner string
	}
	claims := make(map[numberKey][]claim)
	maxNumber := make(map[string]uint64)
	addClaims := func(paths []string, owner string) {
		for _, name := range paths {
			directory, number, _, ok := numberedFileNumber(name)
			if !ok {
				continue
			}
			key := numberKey{directory: directory, number: number}
			claims[key] = append(claims[key], claim{path: name, owner: owner})
			maxNumber[directory] = max(maxNumber[directory], number)
		}
	}
	addClaims(defaultFiles, "")
	for _, earlier := range ahead {
		addClaims(earlier.Paths, earlier.Owner)
	}

	var report NumberedFileReport
	widths := make(map[string]int)
	paths := slices.Clone(subject)
	slices.Sort(paths)
	for _, name := range slices.Compact(paths) {
		directory, number, width, ok := numberedFileNumber(name)
		if !ok {
			continue
		}
		for _, other := range claims[numberKey{directory: directory, number: number}] {
			report.Collisions = append(report.Collisions, NumberedFileCollision{Path: name, Other: other.path, Owner: other.owner})
			widths[directory] = max(widths[directory], width)
		}
	}
	slices.SortFunc(report.Collisions, func(a, b NumberedFileCollision) int {
		if order := cmp.Compare(a.Path, b.Path); order != 0 {
			return order
		}
		if order := cmp.Compare(a.Other, b.Other); order != 0 {
			return order
		}
		return cmp.Compare(numberedFileOwner(a.Owner), numberedFileOwner(b.Owner))
	})
	directories := make([]string, 0, len(widths))
	for directory := range widths {
		directories = append(directories, directory)
	}
	slices.Sort(directories)
	for _, directory := range directories {
		report.NextFree = append(report.NextFree, path.Join(directory, fmt.Sprintf("%0*d", widths[directory], maxNumber[directory]+1)))
	}
	return report
}

func numberedFileNumber(name string) (directory string, number uint64, width int, ok bool) {
	base := path.Base(name)
	for width < len(base) && base[width] >= '0' && base[width] <= '9' {
		width++
	}
	if width == 0 || width > 18 {
		return "", 0, 0, false
	}
	number, err := strconv.ParseUint(base[:width], 10, 64)
	if err != nil {
		return "", 0, 0, false
	}
	return path.Dir(name), number, width, true
}

func numberedFileOwner(owner string) string {
	if owner == "" {
		return "default branch"
	}
	return owner
}

func (report NumberedFileReport) Detail() string {
	if len(report.Collisions) == 0 {
		return ""
	}
	collisions := make([]string, 0, len(report.Collisions))
	for _, collision := range report.Collisions {
		collisions = append(collisions, fmt.Sprintf("%s collides with %s (%s)", collision.Path, collision.Other, numberedFileOwner(collision.Owner)))
	}
	return strings.Join(collisions, ", ") + "; next free number: " + strings.Join(report.NextFree, ", ")
}

func (report NumberedFileReport) Blocker() string {
	if len(report.Collisions) == 0 {
		return ""
	}
	return BlockerNumberedFileCollision + ": " + report.Detail()
}
