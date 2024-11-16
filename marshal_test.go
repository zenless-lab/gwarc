package gwarc_test

import (
	"strings"
	"testing"
	"time"

	. "github.com/zenless-lab/gwarc"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []string
		wantErr bool
	}{
		{
			name: "basic record",
			input: &WARCRecord{
				Version:     WARCVariant1_0,
				Type:        WARCTypeResponse,
				Date:        time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				RecordID:    "<urn:uuid:12345678-1234-1234-1234-123456789012>",
				ContentType: "text/plain",
				Content:     []byte("Hello, World!"),
			},
			want: []string{
				"WARC/1.0",
				"Content-Length: 13",
				"WARC-Type: response",
				"WARC-Date: 2024-01-01T10:00:00Z",
				"WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>",
				"Content-Type: text/plain",
				"Hello, World!",
			},
			wantErr: false,
		},
		{
			name: "with optional fields",
			input: &WARCRecord{
				Version:     WARCVariant1_1,
				Type:        WARCTypeResponse,
				Date:        time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				RecordID:    "<urn:uuid:12345678-1234-1234-1234-123456789012>",
				IPAddress:   "192.168.1.1",
				TargetURI:   "http://example.com",
				ContentType: "text/html",
			},
			want: []string{
				"WARC/1.1",
				"WARC-Type: response",
				"WARC-Date: 2024-01-01T10:00:00Z",
				"WARC-Record-ID: <urn:uuid:12345678-1234-1234-1234-123456789012>",
				"Content-Type: text/html",
				"WARC-IP-Address: 192.168.1.1",
				"WARC-Target-URI: http://example.com",
				"Content-Length: 0",
			},
			wantErr: false,
		},
		{
			name: "with concurrent records",
			input: &WARCRecord{
				Version:  WARCVariant1_0,
				Type:     WARCTypeResource,
				Date:     time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				RecordID: "<urn:uuid:12345678>",
				ConcurrentTo: []string{
					"<urn:uuid:concurrent1>",
					"<urn:uuid:concurrent2>",
				},
			},
			want: []string{
				"WARC/1.0",
				"WARC-Type: resource",
				"WARC-Date: 2024-01-01T10:00:00Z",
				"WARC-Record-ID: <urn:uuid:12345678>",
				"WARC-Concurrent-To: <urn:uuid:concurrent1>, <urn:uuid:concurrent2>",
				"Content-Length: 0",
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "non-struct input",
			input:   "invalid",
			want:    []string{},
			wantErr: true,
		},
		{
			name: "missing required fields",
			input: &WARCRecord{
				Version: WARCVariant1_0,
			},
			want:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotStr := string(got)
			for _, want := range tt.want {
				if !strings.Contains(gotStr, want) {
					t.Errorf("Marshal() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
