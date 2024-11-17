package cdx

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

const CDXTimestampFormat = "20060102150405"

// Unmarshal parses CDX formatted data and stores the result in v
func Unmarshal[T any](data []byte, v T) error {
	// Create scanner to read lines
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// Read header line
	if !scanner.Scan() {
		return errors.New("empty CDX file")
	}

	// Parse header
	header := scanner.Text()
	if !strings.HasPrefix(header, "CDX") {
		return fmt.Errorf("invalid CDX header: %s", header)
	}

	// Parse format from header
	fields := strings.Fields(header)
	format := make(CDXFormat, len(fields)-1)
	for i := range format {
		format[i] = CDXField(fields[i+1][0])
	}

	// Create CDX file
	cdxFile := NewCDXFile(format)

	// Parse records
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Split line into fields
		parts := strings.Fields(line)
		if len(parts) != len(format) {
			return fmt.Errorf("invalid record length: got %d, want %d", len(parts), len(format))
		}

		// Create new record
		record := CDXRecord{}

		// Parse each field
		for i, field := range format {
			value := parts[i]
			if err := setField(&record, field, value); err != nil {
				return fmt.Errorf("error parsing field %c: %v", field, err)
			}
		}

		cdxFile.Records = append(cdxFile.Records, record)
	}

	// Copy result to v using reflection
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}
	rv = rv.Elem()
	rv.Set(reflect.ValueOf(*cdxFile))

	return scanner.Err()
}

// setField sets a field in the CDX record based on its type
func setField(record *CDXRecord, field CDXField, value string) error {
	rv := reflect.ValueOf(record).Elem()

	// Find struct field with matching cdx tag
	for i := 0; i < rv.NumField(); i++ {
		tag := rv.Type().Field(i).Tag.Get("cdx")
		if tag == string(field) {
			f := rv.Field(i)

			switch f.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Int, reflect.Int64:
				n, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				f.SetInt(n)
			case reflect.Struct:
				if f.Type() == reflect.TypeOf(time.Time{}) {
					t, err := time.Parse(CDXTimestampFormat, value)
					if err != nil {
						return err
					}
					f.Set(reflect.ValueOf(t))
				}
			}
			return nil
		}
	}
	return fmt.Errorf("unknown field: %c", field)
}
