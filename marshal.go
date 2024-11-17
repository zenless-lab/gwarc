package gwarc

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Marshal converts a WARCRecord into WARC format bytes
func Marshal(v any) ([]byte, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("v must be a struct or pointer to struct")
	}

	if validator, ok := v.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer

	versionField := val.FieldByName("Version")
	if !versionField.IsValid() {
		return nil, fmt.Errorf("Version field is required")
	}
	fmt.Fprintf(&buf, "WARC/%s\r\n", versionField.Interface())

	headers := make(map[string]string)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("warc")
		if tag == "" {
			continue
		}

		tagParts := parseTag(tag)
		headerName := tagParts.name
		if headerName == "" {
			continue
		}

		value := formatFieldValue(field)
		if value == "" && tagParts.omitempty {
			continue
		}

		headers[headerName] = value
	}

	for name, value := range headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", name, value)
	}

	contentField := val.FieldByName("Content")
	if contentField.IsValid() {
		content := contentField.Bytes()
		headers["Content-Length"] = strconv.Itoa(len(content))
		fmt.Fprintf(&buf, "Content-Length: %d\r\n", len(content))

		buf.WriteString("\r\n")

		buf.Write(content)
	} else {
		buf.WriteString("Content-Length: 0\r\n\r\n")
	}

	return buf.Bytes(), nil
}

type tagInfo struct {
	name      string
	omitempty bool
}

func parseTag(tag string) tagInfo {
	parts := strings.Split(tag, ",")
	info := tagInfo{name: parts[0]}
	if len(parts) > 1 {
		info.omitempty = parts[1] == "omitempty"
	}
	return info
}

func formatFieldValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			var values []string
			for i := 0; i < field.Len(); i++ {
				values = append(values, field.Index(i).String())
			}
			return strings.Join(values, ", ")
		}
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t := field.Interface().(time.Time)
			return t.Format(time.RFC3339)
		}
	}
	return fmt.Sprint(field.Interface())
}
