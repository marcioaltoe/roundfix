package skillhash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ErrUnsupportedLockVersion identifies a lock whose hashing contract is unknown.
var ErrUnsupportedLockVersion = errors.New("unsupported skills-lock.json version")

// UnsupportedLockVersionError preserves the version text written in the lock.
type UnsupportedLockVersionError struct{ Version string }

func (err UnsupportedLockVersionError) Error() string {
	return fmt.Sprintf("skills-lock.json version %s is not supported; this Roundfix reads versions 1 and 2", err.Version)
}

func (err UnsupportedLockVersionError) Is(target error) bool {
	return target == ErrUnsupportedLockVersion
}

// ParseLockVersion reads the supported numeric lock versions.
func ParseLockVersion(text string) (int, error) {
	// Compare the exact numeric value so a fractional version cannot round to 1 or 2.
	if json.Valid([]byte(text)) {
		number, ok := new(big.Rat).SetString(text)
		if ok && number.IsInt() && number.Num().IsInt64() {
			version := number.Num().Int64()
			if version == 1 || version == 2 {
				return int(version), nil
			}
		}
	}
	return 0, UnsupportedLockVersionError{Version: text}
}

// SumForLockVersion computes the folder digest without changing caller-owned files.
func SumForLockVersion(version int, files []File) (string, error) {
	switch version {
	case 1:
		return Sum(files), nil
	case 2:
		ordered := make([]File, 0, len(files))
		for _, file := range files {
			file.Path = filepath.ToSlash(file.Path)
			parts := strings.Split(file.Path, "/")
			excluded := false
			for index, part := range parts {
				if part == "metadata.json" || (index < len(parts)-1 && (part == ".git" || part == "__pycache__" || part == "__pypackages__")) {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
			if !bytes.ContainsRune(file.Content, 0) {
				file.Content = bytes.ReplaceAll(file.Content, []byte("\r\n"), []byte("\n"))
			}
			ordered = append(ordered, file)
		}
		sort.Slice(ordered, func(i, j int) bool { return ordered[i].Path < ordered[j].Path })
		digest := sha256.New()
		for _, file := range ordered {
			_, _ = digest.Write([]byte(file.Path))
			_, _ = digest.Write([]byte{0})
			_, _ = digest.Write([]byte(strconv.Itoa(len(file.Content))))
			_, _ = digest.Write([]byte{0})
			_, _ = digest.Write(file.Content)
		}
		return hex.EncodeToString(digest.Sum(nil)), nil
	default:
		return "", UnsupportedLockVersionError{Version: strconv.Itoa(version)}
	}
}
