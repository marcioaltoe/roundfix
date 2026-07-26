package skillhash

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// File is one already collected relative path and its content.
type File struct {
	Path    string
	Content []byte
}

// Sum returns the external skills CLI digest for files without changing their
// caller-owned order.
func Sum(files []File) string {
	ordered := append([]File(nil), files...)
	for index := range ordered {
		ordered[index].Path = filepath.ToSlash(ordered[index].Path)
	}

	collator := collate.New(language.AmericanEnglish)
	sort.SliceStable(ordered, func(i, j int) bool {
		return collator.CompareString(ordered[i].Path, ordered[j].Path) < 0
	})

	digest := sha256.New()
	for _, file := range ordered {
		_, _ = digest.Write([]byte(file.Path))
		_, _ = digest.Write(file.Content)
	}
	return hex.EncodeToString(digest.Sum(nil))
}
