package baseline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
)

type orderedJSONKind uint8

const (
	orderedJSONNull orderedJSONKind = iota
	orderedJSONObject
	orderedJSONArray
	orderedJSONString
	orderedJSONNumber
	orderedJSONBool
)

type orderedJSONField struct {
	name  string
	value orderedJSONValue
}

type orderedJSONValue struct {
	kind   orderedJSONKind
	object []orderedJSONField
	array  []orderedJSONValue
	text   string
	flag   bool
}

type skillsLockDocument struct {
	before []byte
	root   orderedJSONValue
}

func loadSkillsLock(filename string) (skillsLockDocument, error) {
	info, err := os.Lstat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return skillsLockDocument{
			root: orderedObject(
				orderedJSONField{name: "version", value: orderedNumber("1")},
				orderedJSONField{name: "skills", value: orderedObject()},
			),
		}, nil
	}
	if err != nil {
		return skillsLockDocument{}, restoreError(
			SkillsRestoreExecution,
			"lock.read-failed",
			fmt.Sprintf("Could not read skills-lock.json: %v.", err),
			"Fix repository permissions and rerun the restoration preview.",
			err,
		)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return skillsLockDocument{}, restoreError(
			SkillsRestoreInvalid,
			"lock.unsafe-path",
			"skills-lock.json is not a regular repository file.",
			"Replace it with a regular repository file before restoring skills.",
			errors.New("skills-lock.json is unsafe"),
		)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return skillsLockDocument{}, restoreError(
			SkillsRestoreExecution,
			"lock.read-failed",
			fmt.Sprintf("Could not read skills-lock.json: %v.", err),
			"Fix repository permissions and rerun the restoration preview.",
			err,
		)
	}
	root, err := parseOrderedJSON(data)
	if err != nil {
		return skillsLockDocument{}, restoreError(
			SkillsRestoreInvalid,
			"lock.invalid",
			fmt.Sprintf("skills-lock.json is malformed: %v.", err),
			"Fix skills-lock.json before restoring external skills.",
			err,
		)
	}
	version, versionOK := root.field("version")
	skills, skillsOK := root.field("skills")
	if root.kind != orderedJSONObject ||
		!versionOK ||
		!version.isOne() ||
		!skillsOK ||
		skills.kind != orderedJSONObject {
		return skillsLockDocument{}, restoreError(
			SkillsRestoreInvalid,
			"lock.invalid",
			"skills-lock.json must use version 1 with a skills object.",
			"Fix skills-lock.json before restoring external skills.",
			errors.New("skills-lock.json schema is invalid"),
		)
	}
	return skillsLockDocument{before: data, root: root}, nil
}

func (document skillsLockDocument) clone() skillsLockDocument {
	return skillsLockDocument{
		before: append([]byte(nil), document.before...),
		root:   document.root.clone(),
	}
}

func (document skillsLockDocument) objectEntry(parent, name string) (map[string]any, bool) {
	parentValue, ok := document.root.field(parent)
	if !ok || parentValue.kind != orderedJSONObject {
		return nil, false
	}
	value, ok := parentValue.field(name)
	if !ok || value.kind != orderedJSONObject {
		return nil, false
	}
	result, ok := value.interfaceValue().(map[string]any)
	return result, ok
}

func (document *skillsLockDocument) setSkillsEntry(name string, entry map[string]any) {
	skills := document.root.fieldPointer("skills")
	if skills == nil {
		return
	}
	value := orderedObject(
		orderedJSONField{name: "source", value: orderedString(stringMapValue(entry, "source"))},
		orderedJSONField{name: "ref", value: orderedString(stringMapValue(entry, "ref"))},
		orderedJSONField{name: "sourceType", value: orderedString(stringMapValue(entry, "sourceType"))},
		orderedJSONField{name: "skillPath", value: orderedString(stringMapValue(entry, "skillPath"))},
		orderedJSONField{name: "computedHash", value: orderedString(stringMapValue(entry, "computedHash"))},
	)
	skills.setField(name, value)
}

func (document skillsLockDocument) marshalIndent() ([]byte, error) {
	var output bytes.Buffer
	if err := document.root.writeIndented(&output, 0); err != nil {
		return nil, err
	}
	output.WriteByte('\n')
	return output.Bytes(), nil
}

func parseOrderedJSON(data []byte) (orderedJSONValue, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := parseOrderedJSONValue(decoder)
	if err != nil {
		return orderedJSONValue{}, err
	}
	if err := requireJSONEOF(decoder); err != nil {
		return orderedJSONValue{}, err
	}
	return value, nil
}

func parseOrderedJSONValue(decoder *json.Decoder) (orderedJSONValue, error) {
	token, err := decoder.Token()
	if err != nil {
		return orderedJSONValue{}, err
	}
	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			object := orderedObject()
			for decoder.More() {
				nameToken, err := decoder.Token()
				if err != nil {
					return orderedJSONValue{}, err
				}
				name, ok := nameToken.(string)
				if !ok {
					return orderedJSONValue{}, errors.New("object key is not a string")
				}
				child, err := parseOrderedJSONValue(decoder)
				if err != nil {
					return orderedJSONValue{}, err
				}
				object.setField(name, child)
			}
			if _, err := decoder.Token(); err != nil {
				return orderedJSONValue{}, err
			}
			return object, nil
		case '[':
			array := orderedJSONValue{kind: orderedJSONArray}
			for decoder.More() {
				child, err := parseOrderedJSONValue(decoder)
				if err != nil {
					return orderedJSONValue{}, err
				}
				array.array = append(array.array, child)
			}
			if _, err := decoder.Token(); err != nil {
				return orderedJSONValue{}, err
			}
			return array, nil
		default:
			return orderedJSONValue{}, fmt.Errorf("unexpected delimiter %q", value)
		}
	case string:
		return orderedString(value), nil
	case json.Number:
		return orderedNumber(value.String()), nil
	case bool:
		return orderedJSONValue{kind: orderedJSONBool, flag: value}, nil
	case nil:
		return orderedJSONValue{kind: orderedJSONNull}, nil
	default:
		return orderedJSONValue{}, fmt.Errorf("unsupported JSON token %T", token)
	}
}

func orderedObject(fields ...orderedJSONField) orderedJSONValue {
	return orderedJSONValue{kind: orderedJSONObject, object: fields}
}

func orderedString(value string) orderedJSONValue {
	return orderedJSONValue{kind: orderedJSONString, text: value}
}

func orderedNumber(value string) orderedJSONValue {
	return orderedJSONValue{kind: orderedJSONNumber, text: value}
}

func (value orderedJSONValue) clone() orderedJSONValue {
	cloned := value
	cloned.object = make([]orderedJSONField, len(value.object))
	for index, field := range value.object {
		cloned.object[index] = orderedJSONField{name: field.name, value: field.value.clone()}
	}
	cloned.array = make([]orderedJSONValue, len(value.array))
	for index, item := range value.array {
		cloned.array[index] = item.clone()
	}
	return cloned
}

func (value orderedJSONValue) field(name string) (orderedJSONValue, bool) {
	if value.kind != orderedJSONObject {
		return orderedJSONValue{}, false
	}
	for _, field := range value.object {
		if field.name == name {
			return field.value, true
		}
	}
	return orderedJSONValue{}, false
}

func (value *orderedJSONValue) fieldPointer(name string) *orderedJSONValue {
	if value.kind != orderedJSONObject {
		return nil
	}
	for index := range value.object {
		if value.object[index].name == name {
			return &value.object[index].value
		}
	}
	return nil
}

func (value *orderedJSONValue) setField(name string, child orderedJSONValue) {
	if value.kind != orderedJSONObject {
		return
	}
	for index := range value.object {
		if value.object[index].name == name {
			value.object[index].value = child
			return
		}
	}
	value.object = append(value.object, orderedJSONField{name: name, value: child})
}

func (value orderedJSONValue) isOne() bool {
	if value.kind != orderedJSONNumber {
		return false
	}
	number, err := strconv.ParseFloat(value.text, 64)
	return err == nil && number == 1
}

func (value orderedJSONValue) interfaceValue() any {
	switch value.kind {
	case orderedJSONObject:
		result := make(map[string]any, len(value.object))
		for _, field := range value.object {
			result[field.name] = field.value.interfaceValue()
		}
		return result
	case orderedJSONArray:
		result := make([]any, len(value.array))
		for index, item := range value.array {
			result[index] = item.interfaceValue()
		}
		return result
	case orderedJSONString:
		return value.text
	case orderedJSONNumber:
		number, err := strconv.ParseInt(value.text, 10, 64)
		if err == nil {
			return number
		}
		floating, _ := strconv.ParseFloat(value.text, 64)
		return floating
	case orderedJSONBool:
		return value.flag
	default:
		return nil
	}
}

func (value orderedJSONValue) writeIndented(output io.Writer, indent int) error {
	write := func(text string) error {
		_, err := io.WriteString(output, text)
		return err
	}
	switch value.kind {
	case orderedJSONObject:
		if len(value.object) == 0 {
			return write("{}")
		}
		if err := write("{\n"); err != nil {
			return err
		}
		for index, field := range value.object {
			if err := write(spaces(indent + 2)); err != nil {
				return err
			}
			name, err := json.Marshal(field.name)
			if err != nil {
				return err
			}
			if _, err := output.Write(name); err != nil {
				return err
			}
			if err := write(": "); err != nil {
				return err
			}
			if err := field.value.writeIndented(output, indent+2); err != nil {
				return err
			}
			if index+1 != len(value.object) {
				if err := write(","); err != nil {
					return err
				}
			}
			if err := write("\n"); err != nil {
				return err
			}
		}
		if err := write(spaces(indent)); err != nil {
			return err
		}
		return write("}")
	case orderedJSONArray:
		if len(value.array) == 0 {
			return write("[]")
		}
		if err := write("[\n"); err != nil {
			return err
		}
		for index, item := range value.array {
			if err := write(spaces(indent + 2)); err != nil {
				return err
			}
			if err := item.writeIndented(output, indent+2); err != nil {
				return err
			}
			if index+1 != len(value.array) {
				if err := write(","); err != nil {
					return err
				}
			}
			if err := write("\n"); err != nil {
				return err
			}
		}
		if err := write(spaces(indent)); err != nil {
			return err
		}
		return write("]")
	case orderedJSONString:
		data, err := json.Marshal(value.text)
		if err != nil {
			return err
		}
		_, err = output.Write(data)
		return err
	case orderedJSONNumber:
		return write(value.text)
	case orderedJSONBool:
		return write(strconv.FormatBool(value.flag))
	case orderedJSONNull:
		return write("null")
	default:
		return errors.New("unknown ordered JSON value")
	}
}

func spaces(count int) string {
	return string(bytes.Repeat([]byte{' '}, count))
}

func stringMapValue(values map[string]any, key string) string {
	value, _ := values[key].(string)
	return value
}
