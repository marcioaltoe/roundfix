package baseline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type document map[string]any

func decodeDocument(data []byte, assetPath string) (document, []Diagnostic) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	value, diagnostics, err := decodeJSONValue(decoder, assetPath)
	if err != nil {
		return nil, append(diagnostics, Diagnostic{
			Code: "catalog.json.invalid",
			Path: assetPath,
			Info: err.Error(),
		})
	}
	if token, err := decoder.Token(); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("unexpected trailing token %v", token)
		}
		return nil, append(diagnostics, Diagnostic{
			Code: "catalog.json.trailing",
			Path: assetPath,
			Info: err.Error(),
		})
	}
	object, ok := value.(map[string]any)
	if !ok {
		return nil, append(diagnostics, Diagnostic{
			Code: "catalog.json.object.required",
			Path: assetPath,
		})
	}
	return document(object), diagnostics
}

func decodeJSONValue(decoder *json.Decoder, assetPath string) (any, []Diagnostic, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, nil, err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return token, nil, nil
	}

	switch delimiter {
	case '{':
		object := make(map[string]any)
		var diagnostics []Diagnostic
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return nil, diagnostics, err
			}
			key, ok := keyToken.(string)
			if !ok {
				return nil, diagnostics, fmt.Errorf("object key is %T", keyToken)
			}
			value, nested, err := decodeJSONValue(decoder, assetPath)
			diagnostics = append(diagnostics, nested...)
			if err != nil {
				return nil, diagnostics, err
			}
			if _, exists := object[key]; exists {
				diagnostics = append(diagnostics, Diagnostic{
					Code: "catalog.json.key.duplicate",
					Path: assetPath,
					Info: key,
				})
			}
			object[key] = value
		}
		end, err := decoder.Token()
		if err != nil {
			return nil, diagnostics, err
		}
		if end != json.Delim('}') {
			return nil, diagnostics, fmt.Errorf("object ended with %v", end)
		}
		return object, diagnostics, nil
	case '[':
		var values []any
		var diagnostics []Diagnostic
		for decoder.More() {
			value, nested, err := decodeJSONValue(decoder, assetPath)
			diagnostics = append(diagnostics, nested...)
			if err != nil {
				return nil, diagnostics, err
			}
			values = append(values, value)
		}
		end, err := decoder.Token()
		if err != nil {
			return nil, diagnostics, err
		}
		if end != json.Delim(']') {
			return nil, diagnostics, fmt.Errorf("array ended with %v", end)
		}
		return values, diagnostics, nil
	default:
		return nil, nil, fmt.Errorf("unexpected delimiter %q", delimiter)
	}
}

func objectValue(value any) (document, bool) {
	object, ok := value.(map[string]any)
	return document(object), ok
}

func objectList(value any) ([]document, bool) {
	values, ok := value.([]any)
	if !ok {
		return nil, false
	}
	result := make([]document, 0, len(values))
	for _, value := range values {
		object, ok := objectValue(value)
		if !ok {
			return nil, false
		}
		result = append(result, object)
	}
	return result, true
}

func stringList(value any) ([]string, bool) {
	values, ok := value.([]any)
	if !ok {
		return nil, false
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		text, ok := value.(string)
		if !ok || text == "" {
			return nil, false
		}
		result = append(result, text)
	}
	return result, true
}

func stringValue(doc document, key string) (string, bool) {
	value, ok := doc[key].(string)
	return value, ok && value != ""
}

func integerValue(doc document, key string) (int64, bool) {
	number, ok := doc[key].(json.Number)
	if !ok {
		return 0, false
	}
	value, err := number.Int64()
	return value, err == nil
}
