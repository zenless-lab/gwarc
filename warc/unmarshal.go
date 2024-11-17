package warc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Unmarshal parses WARC formatted data and stores the result in the value pointed to by v.
func Unmarshal[T any](data []byte, v T) error {
	// Check if v implements WARCRecordUnmarshaler
	if unmarshaler, ok := any(v).(WARCRecordUnmarshaler); ok {
		return unmarshaler.UnmarshalWARCRecord(data)
	}

	reader := bufio.NewReader(bytes.NewReader(data))

	versionLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read version: %v", err)
	}

	version := strings.TrimSpace(strings.TrimPrefix(versionLine, "WARC/"))

	warcVersion := WARCVariant(version)
	if warcVersion != WARCVariant1_0 && warcVersion != WARCVariant1_1 {
		return fmt.Errorf("unsupported WARC version: %s", version)
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}
	elem := val.Elem()

	versionField := elem.FieldByName("Version")
	if versionField.IsValid() && versionField.CanSet() {
		versionField.Set(reflect.ValueOf(warcVersion))
	}

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read header: %v", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s", line)
		}

		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	contentLength, _ := strconv.ParseInt(headers["Content-Length"], 10, 64)
	content := make([]byte, contentLength)
	_, err = reader.Read(content)
	if err != nil {
		return fmt.Errorf("failed to read content: %v", err)
	}
	contentField := elem.FieldByName("Content")
	if contentField.IsValid() && contentField.CanSet() {
		contentField.SetBytes(content)
	}

	if val.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}

	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("warc")
		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, ",")
		headerName := tagParts[0]

		value, exists := headers[headerName]
		if !exists && len(tagParts) > 1 && tagParts[1] == "omitempty" {
			continue
		}

		if err := setField(field, value); err != nil {
			return fmt.Errorf("failed to set field %s: %v", fieldType.Name, err)
		}
	}

	return nil
}

func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int64:
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(v)
		}
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if t, err := time.Parse(time.RFC3339, value); err == nil {
				field.Set(reflect.ValueOf(t))
			}
		}
	}
	return nil
}

// Valid checks if the provided data is a valid WARC formatted data.
func Valid(data []byte) error {
	reader := bufio.NewReader(bytes.NewReader(data))

	versionLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read version: %v", err)
	}

	version := strings.TrimSpace(strings.TrimPrefix(versionLine, "WARC/"))

	warcVersion := WARCVariant(version)
	if warcVersion != WARCVariant1_0 && warcVersion != WARCVariant1_1 {
		return fmt.Errorf("unsupported WARC version: %s", version)
	}

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read header: %v", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s", line)
		}

		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	if _, exists := headers["Content-Length"]; !exists {
		return errors.New("missing Content-Length header")
	}

	contentLength, err := strconv.ParseInt(headers["Content-Length"], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Content-Length value: %v", err)
	}

	content := make([]byte, contentLength)
	_, err = reader.Read(content)
	if err != nil {
		return fmt.Errorf("failed to read content: %v", err)
	}

	return nil
}
