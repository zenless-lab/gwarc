package cdx

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Marshal converts a CDX file struct into a byte array
func Marshal(v interface{}) ([]byte, error) {
    // Check if input is nil
    if v == nil {
        return nil, errors.New("input is nil")
    }

    // Check if input is a pointer to CDXFile
    cdxFile, ok := v.(*CDXFile)
    if !ok {
        return nil, errors.New("input must be a pointer to CDXFile")
    }

    var buf bytes.Buffer

    // Write header
    if len(cdxFile.Header.Format) > 0 {
        buf.WriteString(cdxFile.Header.Format.String())
        buf.WriteString("\n")
    }

    // Write records
    for _, record := range cdxFile.Records {
        line, err := marshalRecord(record, cdxFile.Header.Format, cdxFile.Header.Delimiter)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal record: %w", err)
        }
        buf.WriteString(line)
        buf.WriteString("\n")
    }

    return buf.Bytes(), nil
}

// marshalRecord converts a single CDX record to a string
func marshalRecord(record CDXRecord, fields []CDXField, delimiter rune) (string, error) {
    v := reflect.ValueOf(record)
    t := v.Type()

    values := make([]string, len(fields))
    for i, field := range fields {
        fieldTag := string(field)
        
        // Find struct field with matching cdx tag
        for j := 0; j < t.NumField(); j++ {
            if t.Field(j).Tag.Get("cdx") == fieldTag {
                fieldValue := v.Field(j)
                values[i] = formatField(fieldValue)
                break
            }
        }
    }

    return joinFields(values, delimiter), nil
}

// formatField converts a reflect.Value to its string representation
func formatField(v reflect.Value) string {
    switch v.Kind() {
    case reflect.String:
        if v.String() == "" {
            return "-"
        }
        return v.String()
    case reflect.Int, reflect.Int64:
        if v.Int() == 0 {
            return "-"
        }
        return strconv.FormatInt(v.Int(), 10)
    case reflect.Int32:
        if v.Int() == 0 {
            return "-"
        }
        return strconv.FormatInt(int64(v.Int()), 10)
    case reflect.Struct:
        // Special handling for time.Time
        if t, ok := v.Interface().(time.Time); ok {
            if t.IsZero() {
                return "-"
            }
            return t.Format(CDXTimestampFormat)
        }
    }
    return "-"
}

// joinFields joins field values with the specified delimiter
func joinFields(fields []string, delimiter rune) string {
    return strings.Join(fields, string(delimiter))
}