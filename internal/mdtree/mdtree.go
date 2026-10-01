// Package mdtree reads a Markdown entry file with its companion files as one text.
package mdtree

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// Text returns the bytes of entry followed by every .md file directly inside
// companion, in lexical order, each preceded by one newline. A missing
// companion directory yields entry alone. A missing entry is an error.
func Text(fsys fs.FS, entry, companion string) (string, error) {
	data, err := fs.ReadFile(fsys, entry)
	if err != nil {
		return "", fmt.Errorf("read markdown entry %s: %w", entry, err)
	}
	files, err := fs.ReadDir(fsys, companion)
	if errors.Is(err, fs.ErrNotExist) {
		return string(data), nil
	}
	if err != nil {
		return "", fmt.Errorf("read markdown companions %s: %w", companion, err)
	}
	var text strings.Builder
	text.Write(data)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		name := path.Join(companion, file.Name())
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return "", fmt.Errorf("read markdown companion %s: %w", name, err)
		}
		text.WriteByte('\n')
		text.Write(data)
	}
	return text.String(), nil
}
