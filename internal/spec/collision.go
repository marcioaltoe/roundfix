package spec

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// TaskTouchSet is the repository files one Task is known to touch and where
// each path was learned. The set is evidence from written artifacts, not an
// inference about what the Task intends to edit.
type TaskTouchSet struct {
	TaskID string
	Paths  map[string]TouchSource
}

// TouchSource names one artifact that made a Task's repository path known.
type TouchSource string

const (
	TouchFromVerification TouchSource = "verification command"
	TouchFromContext      TouchSource = "declared context"
	TouchFromPriorRun     TouchSource = "prior Run settlement commit"
)

// WaveCollision is two Tasks the Task Graph permits in one Wave that are
// known to touch at least one common repository file.
type WaveCollision struct {
	First  string
	Second string
	Paths  map[string]TouchSource
}

// Collisions reports every pair the Task Graph permits in one Wave that is
// known to touch a common repository file. It reads files and Git objects
// directly; it does not execute Verification, Git, or any other command.
func Collisions(repoRoot string, graph *Graph) ([]WaveCollision, error) {
	if graph == nil {
		return nil, errors.New("find Task Graph collisions: graph is required")
	}
	root, err := collisionRepoRoot(repoRoot)
	if err != nil {
		return nil, err
	}

	sets := make([]TaskTouchSet, len(graph.Tasks))
	for index, task := range graph.Tasks {
		paths, err := declaredTaskTouches(root, task)
		if err != nil {
			return nil, fmt.Errorf("find Task Graph collisions for Task %q: %w", task.ID, err)
		}
		sets[index] = TaskTouchSet{TaskID: task.ID, Paths: paths}
	}
	prior, err := priorRunTaskTouches(root, graph)
	if err != nil {
		return nil, fmt.Errorf("find Task Graph collisions from prior Run: %w", err)
	}
	for index := range sets {
		for path := range prior[sets[index].TaskID] {
			addTouchSource(sets[index].Paths, path, TouchFromPriorRun)
		}
	}

	ordered := taskDependencyClosure(graph.Tasks)
	var collisions []WaveCollision
	for first := 0; first < len(sets); first++ {
		for second := first + 1; second < len(sets); second++ {
			firstID, secondID := sets[first].TaskID, sets[second].TaskID
			if ordered[firstID][secondID] || ordered[secondID][firstID] {
				continue
			}
			shared := sharedTaskTouches(sets[first].Paths, sets[second].Paths)
			if len(shared) == 0 {
				continue
			}
			collisions = append(collisions, WaveCollision{
				First:  firstID,
				Second: secondID,
				Paths:  shared,
			})
		}
	}
	return collisions, nil
}

func collisionRepoRoot(repoRoot string) (string, error) {
	if strings.TrimSpace(repoRoot) == "" {
		return "", errors.New("find Task Graph collisions: repository root is required")
	}
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root %q: %w", repoRoot, err)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root %q: %w", repoRoot, err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("stat repository root %q: %w", root, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("repository root %q is not a directory", root)
	}
	return root, nil
}

func declaredTaskTouches(repoRoot string, task Task) (map[string]TouchSource, error) {
	paths := make(map[string]TouchSource)
	for _, command := range task.Verification {
		for _, candidate := range shellWords(command) {
			path, exists, err := repositoryFile(repoRoot, candidate)
			if err != nil {
				return nil, fmt.Errorf("inspect Verification path %q: %w", candidate, err)
			}
			if exists {
				addTouchSource(paths, path, TouchFromVerification)
			}
		}
	}
	for _, ref := range task.Context {
		path, exists, err := repositoryFile(repoRoot, ref.Path)
		if err != nil {
			return nil, fmt.Errorf("inspect Context path %q: %w", ref.Path, err)
		}
		if exists {
			addTouchSource(paths, path, TouchFromContext)
		}
	}
	return paths, nil
}

func repositoryFile(repoRoot, candidate string) (string, bool, error) {
	candidate = strings.TrimSpace(candidate)
	if candidate == "" || strings.HasPrefix(candidate, "-") ||
		strings.ContainsAny(candidate, "$`*?[]{}") {
		return "", false, nil
	}

	path := candidate
	if !filepath.IsAbs(path) {
		path = filepath.Join(repoRoot, filepath.FromSlash(path))
	}
	path = filepath.Clean(path)
	relative, err := filepath.Rel(repoRoot, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false, nil
	}
	relative = filepath.ToSlash(relative)
	if relative == ".git" || strings.HasPrefix(relative, ".git/") {
		return "", false, nil
	}

	resolved, err := filepath.EvalSymlinks(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	resolvedRelative, err := filepath.Rel(repoRoot, resolved)
	if err != nil || resolvedRelative == ".." || strings.HasPrefix(resolvedRelative, ".."+string(filepath.Separator)) {
		return "", false, nil
	}
	info, err := os.Stat(resolved)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if !info.Mode().IsRegular() {
		return "", false, nil
	}
	return relative, true, nil
}

// shellWords is deliberately smaller than a shell parser. It separates the
// literal words a Verification command writes down while refusing expansion;
// repositoryFile then accepts only words that already resolve to files.
func shellWords(command string) []string {
	var words []string
	var word strings.Builder
	var quote rune
	escaped := false
	flush := func() {
		if word.Len() > 0 {
			words = append(words, word.String())
			word.Reset()
		}
	}
	for _, char := range command {
		if escaped {
			word.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' && quote != '\'' {
			escaped = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				word.WriteRune(char)
			}
			continue
		}
		switch char {
		case '\'', '"':
			quote = char
		case ' ', '\t', '\r', '\n', ';', '|', '&', '(', ')', '<', '>':
			flush()
		default:
			word.WriteRune(char)
		}
	}
	if escaped {
		word.WriteRune('\\')
	}
	flush()
	return words
}

func addTouchSource(paths map[string]TouchSource, path string, source TouchSource) {
	if _, exists := paths[path]; !exists {
		paths[path] = source
	}
}

func sharedTaskTouches(first, second map[string]TouchSource) map[string]TouchSource {
	shared := make(map[string]TouchSource)
	for path, firstSource := range first {
		secondSource, exists := second[path]
		if !exists {
			continue
		}
		shared[path] = mergedTouchSource(firstSource, secondSource)
	}
	return shared
}

func mergedTouchSource(first, second TouchSource) TouchSource {
	if first == second {
		return first
	}
	sources := []string{string(first), string(second)}
	sort.Strings(sources)
	return TouchSource(strings.Join(sources, "; "))
}

func taskDependencyClosure(tasks []Task) map[string]map[string]bool {
	needs := make(map[string][]string, len(tasks))
	for _, task := range tasks {
		needs[task.ID] = append([]string(nil), task.Needs...)
	}
	closure := make(map[string]map[string]bool, len(tasks))
	for _, task := range tasks {
		seen := make(map[string]bool)
		stack := append([]string(nil), needs[task.ID]...)
		for len(stack) > 0 {
			last := len(stack) - 1
			need := stack[last]
			stack = stack[:last]
			if seen[need] {
				continue
			}
			seen[need] = true
			stack = append(stack, needs[need]...)
		}
		closure[task.ID] = seen
	}
	return closure
}

type gitHash [20]byte

func parseGitHash(value string) (gitHash, bool) {
	var hash gitHash
	value = strings.TrimSpace(value)
	if len(value) != 40 {
		return hash, false
	}
	for index := range hash {
		part, err := strconv.ParseUint(value[index*2:index*2+2], 16, 8)
		if err != nil {
			return gitHash{}, false
		}
		hash[index] = byte(part)
	}
	return hash, true
}

func (hash gitHash) String() string {
	const digits = "0123456789abcdef"
	encoded := make([]byte, 40)
	for index, value := range hash {
		encoded[index*2] = digits[value>>4]
		encoded[index*2+1] = digits[value&0x0f]
	}
	return string(encoded)
}

type gitObject struct {
	kind string
	data []byte
}

type packedGitObject struct {
	packPath string
	offset   uint64
}

type gitObjectStore struct {
	gitDir     string
	commonDir  string
	objectDirs []string
	packed     map[gitHash]packedGitObject
	packFiles  map[string]*os.File
	objects    map[gitHash]gitObject
	offsets    map[string]gitObject
}

func openGitObjectStore(repoRoot string) (*gitObjectStore, error) {
	gitPath := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("stat Git directory %q: %w", gitPath, err)
	}
	gitDir := gitPath
	if !info.IsDir() {
		content, err := os.ReadFile(gitPath)
		if err != nil {
			return nil, fmt.Errorf("read Git directory pointer %q: %w", gitPath, err)
		}
		value := strings.TrimSpace(string(content))
		if !strings.HasPrefix(value, "gitdir:") {
			return nil, fmt.Errorf("Git directory pointer %q is malformed", gitPath)
		}
		gitDir = strings.TrimSpace(strings.TrimPrefix(value, "gitdir:"))
		if !filepath.IsAbs(gitDir) {
			gitDir = filepath.Join(repoRoot, gitDir)
		}
		gitDir = filepath.Clean(gitDir)
	}
	commonDir := gitDir
	if content, err := os.ReadFile(filepath.Join(gitDir, "commondir")); err == nil {
		commonDir = filepath.Clean(filepath.Join(gitDir, strings.TrimSpace(string(content))))
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read Git common directory: %w", err)
	}
	store := &gitObjectStore{
		gitDir:    gitDir,
		commonDir: commonDir,
		packed:    make(map[gitHash]packedGitObject),
		packFiles: make(map[string]*os.File),
		objects:   make(map[gitHash]gitObject),
		offsets:   make(map[string]gitObject),
	}
	if err := store.addObjectDir(filepath.Join(commonDir, "objects"), map[string]bool{}); err != nil {
		store.close()
		return nil, err
	}
	return store, nil
}

func (store *gitObjectStore) close() {
	for _, file := range store.packFiles {
		_ = file.Close()
	}
}

func (store *gitObjectStore) addObjectDir(objectDir string, seen map[string]bool) error {
	objectDir = filepath.Clean(objectDir)
	if seen[objectDir] {
		return nil
	}
	seen[objectDir] = true
	store.objectDirs = append(store.objectDirs, objectDir)

	indexes, err := filepath.Glob(filepath.Join(objectDir, "pack", "*.idx"))
	if err != nil {
		return fmt.Errorf("list Git pack indexes in %q: %w", objectDir, err)
	}
	sort.Strings(indexes)
	for _, indexPath := range indexes {
		if err := store.readPackIndex(indexPath); err != nil {
			return err
		}
	}
	alternates, err := os.Open(filepath.Join(objectDir, "info", "alternates"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open Git object alternates: %w", err)
	}
	defer alternates.Close()
	scanner := bufio.NewScanner(alternates)
	for scanner.Scan() {
		alternate := strings.TrimSpace(scanner.Text())
		if alternate == "" {
			continue
		}
		if !filepath.IsAbs(alternate) {
			alternate = filepath.Join(objectDir, alternate)
		}
		if err := store.addObjectDir(alternate, seen); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read Git object alternates: %w", err)
	}
	return nil
}

func (store *gitObjectStore) readPackIndex(indexPath string) error {
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("read Git pack index %q: %w", indexPath, err)
	}
	if len(content) < 8+256*4 || !bytes.Equal(content[:4], []byte{0xff, 't', 'O', 'c'}) || binary.BigEndian.Uint32(content[4:8]) != 2 {
		return fmt.Errorf("Git pack index %q is not version 2", indexPath)
	}
	count := int(binary.BigEndian.Uint32(content[8+255*4 : 8+256*4]))
	hashesAt := 8 + 256*4
	offsetsAt := hashesAt + count*20 + count*4
	if count < 0 || offsetsAt+count*4+40 > len(content) {
		return fmt.Errorf("Git pack index %q is truncated", indexPath)
	}
	largeOffsetsAt := offsetsAt + count*4
	packPath := strings.TrimSuffix(indexPath, ".idx") + ".pack"
	for index := 0; index < count; index++ {
		var hash gitHash
		copy(hash[:], content[hashesAt+index*20:hashesAt+(index+1)*20])
		rawOffset := binary.BigEndian.Uint32(content[offsetsAt+index*4 : offsetsAt+(index+1)*4])
		var offset uint64
		if rawOffset&0x80000000 == 0 {
			offset = uint64(rawOffset)
		} else {
			largeIndex := int(rawOffset & 0x7fffffff)
			position := largeOffsetsAt + largeIndex*8
			if position+8 > len(content)-40 {
				return fmt.Errorf("Git pack index %q has an invalid large offset", indexPath)
			}
			offset = binary.BigEndian.Uint64(content[position : position+8])
		}
		if _, exists := store.packed[hash]; !exists {
			store.packed[hash] = packedGitObject{packPath: packPath, offset: offset}
		}
	}
	return nil
}

func (store *gitObjectStore) object(hash gitHash) (gitObject, error) {
	if object, exists := store.objects[hash]; exists {
		return object, nil
	}
	for _, objectDir := range store.objectDirs {
		path := filepath.Join(objectDir, hash.String()[:2], hash.String()[2:])
		file, err := os.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return gitObject{}, fmt.Errorf("open loose Git object %s: %w", hash, err)
		}
		reader, err := zlib.NewReader(file)
		if err != nil {
			_ = file.Close()
			return gitObject{}, fmt.Errorf("decompress loose Git object %s: %w", hash, err)
		}
		content, readErr := io.ReadAll(reader)
		closeErr := reader.Close()
		fileErr := file.Close()
		if readErr != nil {
			return gitObject{}, fmt.Errorf("read loose Git object %s: %w", hash, readErr)
		}
		if closeErr != nil {
			return gitObject{}, fmt.Errorf("close loose Git object %s: %w", hash, closeErr)
		}
		if fileErr != nil {
			return gitObject{}, fmt.Errorf("close loose Git object file %s: %w", hash, fileErr)
		}
		header, data, found := bytes.Cut(content, []byte{0})
		kind, _, valid := strings.Cut(string(header), " ")
		if !found || !valid {
			return gitObject{}, fmt.Errorf("loose Git object %s has a malformed header", hash)
		}
		object := gitObject{kind: kind, data: data}
		store.objects[hash] = object
		return object, nil
	}
	packed, exists := store.packed[hash]
	if !exists {
		return gitObject{}, fmt.Errorf("Git object %s is missing", hash)
	}
	object, err := store.packObject(packed.packPath, packed.offset)
	if err != nil {
		return gitObject{}, fmt.Errorf("read packed Git object %s: %w", hash, err)
	}
	store.objects[hash] = object
	return object, nil
}

func (store *gitObjectStore) packObject(packPath string, offset uint64) (gitObject, error) {
	key := packPath + "\x00" + strconv.FormatUint(offset, 10)
	if object, exists := store.offsets[key]; exists {
		return object, nil
	}
	file := store.packFiles[packPath]
	if file == nil {
		opened, err := os.Open(packPath)
		if err != nil {
			return gitObject{}, fmt.Errorf("open Git pack %q: %w", packPath, err)
		}
		file = opened
		store.packFiles[packPath] = file
	}
	kind, dataAt, baseOffset, baseHash, err := readPackObjectHeader(file, offset)
	if err != nil {
		return gitObject{}, err
	}
	reader, err := zlib.NewReader(io.NewSectionReader(file, int64(dataAt), 1<<63-1-int64(dataAt)))
	if err != nil {
		return gitObject{}, fmt.Errorf("decompress Git pack object at %d: %w", offset, err)
	}
	data, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr != nil {
		return gitObject{}, fmt.Errorf("read Git pack object at %d: %w", offset, readErr)
	}
	if closeErr != nil {
		return gitObject{}, fmt.Errorf("close Git pack object at %d: %w", offset, closeErr)
	}
	object := gitObject{kind: kind, data: data}
	if kind == "ofs-delta" || kind == "ref-delta" {
		var base gitObject
		if kind == "ofs-delta" {
			base, err = store.packObject(packPath, baseOffset)
		} else {
			base, err = store.object(baseHash)
		}
		if err != nil {
			return gitObject{}, err
		}
		data, err = applyGitDelta(base.data, data)
		if err != nil {
			return gitObject{}, fmt.Errorf("apply Git delta at %d: %w", offset, err)
		}
		object = gitObject{kind: base.kind, data: data}
	}
	store.offsets[key] = object
	return object, nil
}

func readPackObjectHeader(file *os.File, offset uint64) (string, uint64, uint64, gitHash, error) {
	var baseHash gitHash
	position := offset
	value, err := readByteAt(file, &position)
	if err != nil {
		return "", 0, 0, baseHash, err
	}
	typeCode := (value >> 4) & 7
	for value&0x80 != 0 {
		value, err = readByteAt(file, &position)
		if err != nil {
			return "", 0, 0, baseHash, err
		}
	}
	kinds := map[byte]string{1: "commit", 2: "tree", 3: "blob", 4: "tag", 6: "ofs-delta", 7: "ref-delta"}
	kind := kinds[typeCode]
	if kind == "" {
		return "", 0, 0, baseHash, fmt.Errorf("Git pack object at %d has unsupported type %d", offset, typeCode)
	}
	var baseOffset uint64
	if typeCode == 6 {
		value, err = readByteAt(file, &position)
		if err != nil {
			return "", 0, 0, baseHash, err
		}
		distance := uint64(value & 0x7f)
		for value&0x80 != 0 {
			value, err = readByteAt(file, &position)
			if err != nil {
				return "", 0, 0, baseHash, err
			}
			distance = ((distance + 1) << 7) | uint64(value&0x7f)
		}
		if distance > offset {
			return "", 0, 0, baseHash, fmt.Errorf("Git pack object at %d has invalid base distance", offset)
		}
		baseOffset = offset - distance
	} else if typeCode == 7 {
		if _, err := file.ReadAt(baseHash[:], int64(position)); err != nil {
			return "", 0, 0, baseHash, fmt.Errorf("read Git delta base at %d: %w", offset, err)
		}
		position += uint64(len(baseHash))
	}
	return kind, position, baseOffset, baseHash, nil
}

func readByteAt(file *os.File, position *uint64) (byte, error) {
	var buffer [1]byte
	if _, err := file.ReadAt(buffer[:], int64(*position)); err != nil {
		return 0, err
	}
	(*position)++
	return buffer[0], nil
}

func applyGitDelta(base, delta []byte) ([]byte, error) {
	baseSize, rest, ok := gitDeltaSize(delta)
	if !ok || baseSize != uint64(len(base)) {
		return nil, errors.New("delta has an invalid base size")
	}
	resultSize, rest, ok := gitDeltaSize(rest)
	if !ok {
		return nil, errors.New("delta has an invalid result size")
	}
	result := make([]byte, 0, resultSize)
	for len(rest) > 0 {
		opcode := rest[0]
		rest = rest[1:]
		if opcode&0x80 == 0 {
			if opcode == 0 || int(opcode) > len(rest) {
				return nil, errors.New("delta has an invalid insert")
			}
			result = append(result, rest[:opcode]...)
			rest = rest[opcode:]
			continue
		}
		var copyOffset, copySize uint32
		for index, mask := range []byte{0x01, 0x02, 0x04, 0x08} {
			if opcode&mask != 0 {
				if len(rest) == 0 {
					return nil, errors.New("delta has a truncated copy offset")
				}
				copyOffset |= uint32(rest[0]) << (index * 8)
				rest = rest[1:]
			}
		}
		for index, mask := range []byte{0x10, 0x20, 0x40} {
			if opcode&mask != 0 {
				if len(rest) == 0 {
					return nil, errors.New("delta has a truncated copy size")
				}
				copySize |= uint32(rest[0]) << (index * 8)
				rest = rest[1:]
			}
		}
		if copySize == 0 {
			copySize = 0x10000
		}
		end := uint64(copyOffset) + uint64(copySize)
		if end > uint64(len(base)) {
			return nil, errors.New("delta copies beyond its base")
		}
		result = append(result, base[copyOffset:end]...)
	}
	if uint64(len(result)) != resultSize {
		return nil, errors.New("delta result size does not match its header")
	}
	return result, nil
}

func gitDeltaSize(content []byte) (uint64, []byte, bool) {
	var size uint64
	for shift := 0; shift < 64 && len(content) > 0; shift += 7 {
		value := content[0]
		content = content[1:]
		size |= uint64(value&0x7f) << shift
		if value&0x80 == 0 {
			return size, content, true
		}
	}
	return 0, nil, false
}

type priorTaskCommit struct {
	hash      gitHash
	tree      gitHash
	parent    gitHash
	hasParent bool
	time      int64
}

func priorRunTaskTouches(repoRoot string, graph *Graph) (map[string]map[string]bool, error) {
	result := make(map[string]map[string]bool, len(graph.Tasks))
	for _, task := range graph.Tasks {
		result[task.ID] = make(map[string]bool)
	}
	if graph.Spec.Slug == "" || len(graph.Tasks) == 0 {
		return result, nil
	}
	store, err := openGitObjectStore(repoRoot)
	if err != nil || store == nil {
		return result, err
	}
	defer store.close()
	roots, err := store.references()
	if err != nil {
		return nil, err
	}
	wanted := make(map[string]bool, len(graph.Tasks))
	for _, task := range graph.Tasks {
		wanted[task.ID] = true
	}
	candidates := make(map[string]priorTaskCommit)
	seen := make(map[gitHash]bool)
	stack := append([]gitHash(nil), roots...)
	for len(stack) > 0 {
		last := len(stack) - 1
		hash := stack[last]
		stack = stack[:last]
		if seen[hash] {
			continue
		}
		seen[hash] = true
		object, err := store.object(hash)
		if err != nil {
			return nil, err
		}
		if object.kind != "commit" {
			continue
		}
		commit, parents, specSlug, taskID, err := parsePriorTaskCommit(hash, object.data)
		if err != nil {
			return nil, err
		}
		stack = append(stack, parents...)
		if specSlug != graph.Spec.Slug || !wanted[taskID] {
			continue
		}
		previous, exists := candidates[taskID]
		if !exists || commit.time > previous.time || (commit.time == previous.time && hash.String() > previous.hash.String()) {
			candidates[taskID] = commit
		}
	}
	for taskID, commit := range candidates {
		changed, err := store.changedPaths(commit)
		if err != nil {
			return nil, fmt.Errorf("read Task %q settlement commit %s: %w", taskID, commit.hash, err)
		}
		for _, candidate := range changed {
			path, exists, err := repositoryFile(repoRoot, candidate)
			if err != nil {
				return nil, err
			}
			if exists {
				result[taskID][path] = true
			}
		}
	}
	return result, nil
}

func parsePriorTaskCommit(hash gitHash, content []byte) (priorTaskCommit, []gitHash, string, string, error) {
	commit := priorTaskCommit{hash: hash}
	var parents []gitHash
	var specSlug, taskID string
	header, message, _ := bytes.Cut(content, []byte("\n\n"))
	for _, line := range strings.Split(string(header), "\n") {
		key, value, found := strings.Cut(line, " ")
		if !found {
			continue
		}
		switch key {
		case "tree":
			parsed, ok := parseGitHash(value)
			if !ok {
				return commit, nil, "", "", fmt.Errorf("Git commit %s has an invalid tree", hash)
			}
			commit.tree = parsed
		case "parent":
			parsed, ok := parseGitHash(value)
			if !ok {
				return commit, nil, "", "", fmt.Errorf("Git commit %s has an invalid parent", hash)
			}
			parents = append(parents, parsed)
		case "committer":
			fields := strings.Fields(value)
			if len(fields) >= 2 {
				commit.time, _ = strconv.ParseInt(fields[len(fields)-2], 10, 64)
			}
		}
	}
	if len(parents) > 0 {
		commit.parent = parents[0]
		commit.hasParent = true
	}
	for _, line := range strings.Split(string(message), "\n") {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Roundfix-Spec":
			specSlug = strings.TrimSpace(value)
		case "Roundfix-Task":
			taskID = strings.TrimSpace(value)
		}
	}
	return commit, parents, specSlug, taskID, nil
}

func (store *gitObjectStore) references() ([]gitHash, error) {
	seen := make(map[gitHash]bool)
	var hashes []gitHash
	add := func(value string) {
		if hash, ok := parseGitHash(value); ok && !seen[hash] {
			seen[hash] = true
			hashes = append(hashes, hash)
		}
	}
	for _, root := range []string{store.gitDir, store.commonDir} {
		head, err := store.resolveReference(root, "HEAD", map[string]bool{})
		if err != nil {
			return nil, err
		}
		add(head)
		refsRoot := filepath.Join(root, "refs")
		err = filepath.WalkDir(refsRoot, func(path string, entry os.DirEntry, walkErr error) error {
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			add(strings.TrimSpace(string(content)))
			return nil
		})
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read Git references: %w", err)
		}
		packed, err := os.Open(filepath.Join(root, "packed-refs"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("open packed Git references: %w", err)
		}
		scanner := bufio.NewScanner(packed)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || line[0] == '#' || line[0] == '^' {
				continue
			}
			value, _, _ := strings.Cut(line, " ")
			add(value)
		}
		scanErr := scanner.Err()
		closeErr := packed.Close()
		if scanErr != nil {
			return nil, fmt.Errorf("read packed Git references: %w", scanErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close packed Git references: %w", closeErr)
		}
	}
	return hashes, nil
}

func (store *gitObjectStore) resolveReference(root, name string, seen map[string]bool) (string, error) {
	path := filepath.Join(root, filepath.FromSlash(name))
	if seen[path] {
		return "", nil
	}
	seen[path] = true
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read Git reference %q: %w", name, err)
	}
	value := strings.TrimSpace(string(content))
	if strings.HasPrefix(value, "ref:") {
		ref := strings.TrimSpace(strings.TrimPrefix(value, "ref:"))
		for _, candidateRoot := range []string{store.gitDir, store.commonDir} {
			resolved, err := store.resolveReference(candidateRoot, ref, seen)
			if err != nil || resolved != "" {
				return resolved, err
			}
		}
		return "", nil
	}
	return value, nil
}

func (store *gitObjectStore) changedPaths(commit priorTaskCommit) ([]string, error) {
	current, err := store.treeFiles(commit.tree)
	if err != nil {
		return nil, err
	}
	previous := make(map[string]gitHash)
	if commit.hasParent {
		parentObject, err := store.object(commit.parent)
		if err != nil {
			return nil, err
		}
		if parentObject.kind != "commit" {
			return nil, fmt.Errorf("parent %s is not a commit", commit.parent)
		}
		parent, _, _, _, err := parsePriorTaskCommit(commit.parent, parentObject.data)
		if err != nil {
			return nil, err
		}
		previous, err = store.treeFiles(parent.tree)
		if err != nil {
			return nil, err
		}
	}
	changed := make(map[string]bool)
	for path, hash := range current {
		if previousHash, exists := previous[path]; !exists || previousHash != hash {
			changed[path] = true
		}
	}
	for path := range previous {
		if _, exists := current[path]; !exists {
			changed[path] = true
		}
	}
	paths := make([]string, 0, len(changed))
	for path := range changed {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func (store *gitObjectStore) treeFiles(root gitHash) (map[string]gitHash, error) {
	files := make(map[string]gitHash)
	var walk func(gitHash, string) error
	walk = func(hash gitHash, prefix string) error {
		object, err := store.object(hash)
		if err != nil {
			return err
		}
		if object.kind != "tree" {
			return fmt.Errorf("Git object %s is %s, not a tree", hash, object.kind)
		}
		content := object.data
		for len(content) > 0 {
			headerEnd := bytes.IndexByte(content, 0)
			if headerEnd < 0 || headerEnd+1+20 > len(content) {
				return fmt.Errorf("Git tree %s is malformed", hash)
			}
			header := string(content[:headerEnd])
			mode, name, found := strings.Cut(header, " ")
			if !found || name == "" {
				return fmt.Errorf("Git tree %s has a malformed entry", hash)
			}
			var child gitHash
			copy(child[:], content[headerEnd+1:headerEnd+1+20])
			content = content[headerEnd+1+20:]
			path := name
			if prefix != "" {
				path = prefix + "/" + name
			}
			if mode == "40000" || mode == "040000" {
				if err := walk(child, path); err != nil {
					return err
				}
				continue
			}
			files[path] = child
		}
		return nil
	}
	if err := walk(root, ""); err != nil {
		return nil, err
	}
	return files, nil
}
