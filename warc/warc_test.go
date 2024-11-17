package warc

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewWARC(t *testing.T) {
	data := "WARC/1.0\nWARC-Type: warcinfo\nWARC-Record-ID: <urn:uuid:1234>\nContent-Length: 0\nWARC-Date: 2023-10-10T10:10:10Z\n\n"
	r := strings.NewReader(data)
	warc := NewWARC(r)
	if warc == nil {
		t.Fatal("Expected non-nil WARC instance")
	}
}

func TestNewWARCFromFile(t *testing.T) {
	// Create a temporary file with WARC data
	file, err := os.CreateTemp("", "test.warc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	data := "WARC/1.0\nWARC-Type: warcinfo\nWARC-Record-ID: <urn:uuid:1234>\nContent-Length: 0\nWARC-Date: 2023-10-10T10:10:10Z\n\n"
	if _, err := file.Write([]byte(data)); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	warc, err := NewWARCFromFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if warc == nil {
		t.Fatal("Expected non-nil WARC instance")
	}
}

func TestNewWARCFromString(t *testing.T) {
	data := "WARC/1.0\nWARC-Type: warcinfo\nWARC-Record-ID: <urn:uuid:1234>\nContent-Length: 0\nWARC-Date: 2023-10-10T10:10:10Z\n\n"
	warc := NewWARCFromString(data)
	if warc == nil {
		t.Fatal("Expected non-nil WARC instance")
	}
}

func TestNewWARCFromBytes(t *testing.T) {
	data := []byte("WARC/1.0\nWARC-Type: warcinfo\nWARC-Record-ID: <urn:uuid:1234>\nContent-Length: 0\nWARC-Date: 2023-10-10T10:10:10Z\n\n")
	warc := NewWARCFromBytes(data)
	if warc == nil {
		t.Fatal("Expected non-nil WARC instance")
	}
}

func TestNextChunk(t *testing.T) {
	data := "WARC/1.0\nWARC-Type: warcinfo\nWARC-Record-ID: <urn:uuid:1234>\nContent-Length: 0\nWARC-Date: 2023-10-10T10:10:10Z\n\n"
	r := strings.NewReader(data)
	warc := NewWARC(r)

	chunk, err := warc.NextChunk()
	if err != nil {
		t.Fatal(err)
	}
	if chunk == nil {
		t.Fatal("Expected non-nil chunk")
	}
}

func TestNext(t *testing.T) {
	data := `WARC/1.0
WARC-Type: warcinfo
WARC-Record-ID: <urn:uuid:1234>
Content-Length: 114
WARC-Date: 2023-10-10T10:10:10Z

operator: test
software: test
robots: test
hostname: test
ip: test
description: test
useragent: test
format: test

WARC/1.0
WARC-Type: response
WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>
Content-Length: 13
Content-Type: text/plain

Hello, World!
`
	r := strings.NewReader(data)
	warc := NewWARC(r)

	record, kind, err := warc.Next()
	if err != nil {
		t.Fatal(err)
	}
	if record == nil {
		t.Fatal("Expected non-nil record")
	}
	if kind != WARCTypeWarcinfo {
		t.Fatalf("Expected kind to be %s, got %s", WARCTypeWarcinfo, kind)
	}

	warcRecord, ok := record.(WarcInfoRecord)
	if !ok {
		t.Fatal("Expected record to be of type *WARCRecord")
	}
	if warcRecord.Version != WARCVariant1_0 {
		t.Fatalf("Expected version to be %s, got %s", WARCVariant1_0, warcRecord.Version)
	}
	if warcRecord.RecordID != "<urn:uuid:1234>" {
		t.Fatalf("Expected RecordID to be %s, got %s", "<urn:uuid:1234>", warcRecord.RecordID)
	}
	if warcRecord.ContentLength != 0 {
		t.Fatalf("Expected ContentLength to be %d, got %d", 0, warcRecord.ContentLength)
	}
	expectedDate, _ := time.Parse(time.RFC3339, "2023-10-10T10:10:10Z")
	if !warcRecord.Date.Equal(expectedDate) {
		t.Fatalf("Expected Date to be %s, got %s", expectedDate, warcRecord.Date)
	}
}
